package webcrawler

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	parser "github.com/firozt/crawler/src/internal/Parser"
	repository "github.com/firozt/crawler/src/internal/Repository"
)

type WebCrawler struct {
	repo                     *repository.PagesRepository
	graphRepo                *repository.GraphRepository
	MAX_ADDED_LINKS_PER_PAGE uint8
	MAX_UNIQUE_CRAWLED_PAGES uint64
	NUM_OF_WORKERS           uint8
}

type FetchedWebData struct {
	HTMLBody string
	Domain   string
	URL      string
}

func NewCrawler(repo *repository.PagesRepository, graphRepo *repository.GraphRepository, MAX_ADDED_LINKS_PER_PAGE uint8, NUM_OF_WORKERS uint8, MAX_UNIQUE_CRAWLED_PAGES uint64) *WebCrawler {
	return &WebCrawler{
		repo:                     repo,
		graphRepo:                graphRepo,
		MAX_ADDED_LINKS_PER_PAGE: MAX_ADDED_LINKS_PER_PAGE,
		NUM_OF_WORKERS:           NUM_OF_WORKERS,
		MAX_UNIQUE_CRAWLED_PAGES: MAX_UNIQUE_CRAWLED_PAGES,
	}
}

func (c *WebCrawler) StartCrawl(seedURL string, allowExternal bool) (uint64, error) {
	println("============ STARTED CRAWLING ============")

	seedURL = parser.NormalizeURL(seedURL)
	if !parser.IsValidURL(seedURL) {
		return 0, errors.New("initial URL is not valid")
	}

	var wg sync.WaitGroup
	var seenUrlCache sync.Map

	const MAX_QUEUE_SIZE = 1000
	urlQueue := make(chan string, MAX_QUEUE_SIZE)
	toParseQueue := make(chan *FetchedWebData, MAX_QUEUE_SIZE)
	pageQueue := make(chan *repository.Page, MAX_QUEUE_SIZE)

	var pending int64 // stores pending actions, when this is 0 we can safely end all goroutines
	var numCrawledPages uint64
	atomic.AddInt64(&pending, 1) // seed URL counts as pending

	seenUrlCache.LoadOrStore(seedURL, true)
	urlQueue <- seedURL

	// start DB worker
	wg.Add(1)
	go c.databaseInteractionWorker(&wg, pageQueue)

	// start fetcher work pool
	for i := uint8(0); i < c.NUM_OF_WORKERS; i++ {
		wg.Add(1)
		go c.fetcherWorker(i, &wg, urlQueue, toParseQueue, &pending, &numCrawledPages)
	}
	c.graphRepo.InsertPageNode(seedURL, "NO_TITLE")
	// start parser work pool
	for i := uint8(0); i < c.NUM_OF_WORKERS; i++ {
		wg.Add(1)
		go c.parserWorker(i, &seenUrlCache, &wg, urlQueue, toParseQueue, pageQueue, allowExternal, &pending, &numCrawledPages)
	}

	// start manager routine
	go func() {
		for {
			// close all chanels (stop all goroutines) when conditions met
			if atomic.LoadInt64(&pending) == 0 {
				close(urlQueue)
				close(toParseQueue)
				for len(pageQueue) > 0 {
					time.Sleep(500 * time.Millisecond)
				}
				close(pageQueue)
				return
			}
			time.Sleep(500 * time.Millisecond) // polling rate, twice a second
		}
	}()

	wg.Wait() // wait for all routines to finish
	println("============ FINISHED CRAWLING ============")
	return numCrawledPages, nil
}

// fetcher worker: gets URLs from urlQueue, fetches HTML, sends to toParseQueue
func (c *WebCrawler) fetcherWorker(workerId uint8, wg *sync.WaitGroup, urlQueue chan string, toParseQueue chan *FetchedWebData, pending *int64, numCrawledPages *uint64) {
	defer wg.Done()

	for url := range urlQueue {
		fmt.Printf("fetcher-%d: acquired url: %s\n", workerId, url)

		htmlBody, err := parser.ParseSite(url)
		if err != nil {
			fmt.Printf("--fetcher-%d couldnt parse URL, dropping - err: %s\n", workerId, err)
			atomic.AddInt64(pending, -1) // cant parse, end
			continue
		}
		// if atomic.LoadUint64(numCrawledPages) < c.MAX_UNIQUE_CRAWLED_PAGES { // channel still open
		toParseQueue <- &FetchedWebData{
			URL:      url,
			Domain:   parser.GetDomain(url),
			HTMLBody: htmlBody,
		}
		time.Sleep(1 * time.Second)
	}
}

// parser worker: parses HTML, produces pages and new URLs
func (c *WebCrawler) parserWorker(workerId uint8, seenUrlCache *sync.Map, wg *sync.WaitGroup, urlQueue chan string, toParseQueue chan *FetchedWebData, pageQueue chan *repository.Page, allowExternal bool, pending *int64, numCrawledPages *uint64) {
	defer wg.Done()

	for fetchedData := range toParseQueue {

		fmt.Printf("parser-%d: acquired fetchedData\n", workerId)

		textData, links, title, err := parser.GetTextAndLinks(fetchedData.HTMLBody, fetchedData.Domain)
		if title == "" {
			split := strings.Split(fetchedData.URL, "/")
			lastResource := split[len(split)-1]
			c.graphRepo.AlterPageNodeTitle(fetchedData.URL, strings.ReplaceAll(lastResource, "%20", " ")) // file resource instead
		} else {
			c.graphRepo.AlterPageNodeTitle(fetchedData.URL, title) // update with new title
		}

		if err != nil {
			fmt.Printf("--parser-%d: Could not get text and links from html. dropping \n", workerId)
			atomic.AddInt64(pending, -1)
			continue
		}

		cleanedTextData := parser.CleanText(strings.Join(textData, " "))
		page := repository.Page{
			Title:   title,
			URL:     fetchedData.URL,
			Content: cleanedTextData,
		}

		// tell db worker to insert this
		pageQueue <- &page

		// this site is done
		atomic.AddUint64(numCrawledPages, 1)
		fmt.Printf("parser-%d finished parsing\n", workerId)

		// Add validated links to the queue
		validatedLinks := parser.ValidateLinks(links, fetchedData.URL, fetchedData.Domain, allowExternal)
		successfullAdds := uint8(0) // holds number of actual adds to queue
		for _, link := range validatedLinks {
			// check config upper bounds
			if successfullAdds >= c.MAX_ADDED_LINKS_PER_PAGE || atomic.LoadUint64(numCrawledPages) >= c.MAX_UNIQUE_CRAWLED_PAGES {
				break
			}
			_, seen := seenUrlCache.LoadOrStore(link, true)
			// is seen ignore
			if !seen && atomic.LoadUint64(numCrawledPages)+uint64(atomic.LoadInt64(pending)) < c.MAX_UNIQUE_CRAWLED_PAGES {
				// add new valids
				c.graphRepo.InsertPageNode(link, "")              // add node (no name for now)
				c.graphRepo.InsertPageEdge(fetchedData.URL, link) // add edge
				successfullAdds++
				atomic.AddInt64(pending, 1)
				urlQueue <- link
			}

		}

		// finished processing this page and adding all links
		fmt.Printf("parser-%d finished parsing\n", workerId)
		atomic.AddInt64(pending, -1)
	}
}

// DB worker single goroutine inserts pages into the database
func (c *WebCrawler) databaseInteractionWorker(wg *sync.WaitGroup, pageQueue chan *repository.Page) {
	defer wg.Done()
	for page := range pageQueue {
		if err := c.repo.InsertPage(*page); err != nil {
			fmt.Printf("--Could not insert into page:  %s\n ", err)
		} else {
			fmt.Printf("inserted page %s  to db\n", page.URL)
		}
	}
}

// simple passthrough, sqlite does the heavy lifting here
func (c *WebCrawler) SearchCrawled(phrase string, url string, limit int) []repository.Page {
	pages := c.repo.SearchPages(phrase, parser.GetDomain(url), limit)
	return pages
}

func (c *WebCrawler) GetGraphData(url string) repository.LinkGraph {
	graph := c.graphRepo.GetAllLinkRelations(parser.GetDomain(url))
	return *graph
}
