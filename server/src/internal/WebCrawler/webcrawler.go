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
	// TSQ "github.com/firozt/crawler/src/internal/ThreadSafeQueue"
)

type WebCrawler struct {
	repo                     *repository.PagesRepository
	MAX_ADDED_LINKS_PER_PAGE uint8
	NUM_OF_WORKERS           uint8
}

type FetchedWebData struct {
	HTMLBody string
	Domain   string
	URL      string
}

func NewCrawler(repo *repository.PagesRepository, MAX_ADDED_LINKS_PER_PAGE uint8, NUM_OF_WORKERS uint8) *WebCrawler {
	return &WebCrawler{
		repo:                     repo,
		MAX_ADDED_LINKS_PER_PAGE: MAX_ADDED_LINKS_PER_PAGE,
		NUM_OF_WORKERS:           NUM_OF_WORKERS,
	}
}

func (c *WebCrawler) StartCrawl(seedURL string, allowExternal bool) error {
	println("============ STARTED CRAWLING ============")

	seedURL = parser.NormalizeURL(seedURL)
	if !parser.IsValidURL(seedURL) {
		return errors.New("initial URL is not valid")
	}

	var wg sync.WaitGroup
	var seenUrlCache sync.Map

	const MAX_QUEUE_SIZE = 1000
	urlQueue := make(chan string, MAX_QUEUE_SIZE)
	toParseQueue := make(chan *FetchedWebData, MAX_QUEUE_SIZE)
	pageQueue := make(chan *repository.Page, MAX_QUEUE_SIZE)

	var pending int64            // stores pending actions, when this is 0 we can safely end all goroutines
	atomic.AddInt64(&pending, 1) // seed URL counts as pending
	urlQueue <- seedURL

	// start DB worker
	wg.Add(1)
	go databaseInteractionWorker(&wg, pageQueue, c.repo)

	// start fetcher work pool
	for i := uint8(0); i < c.NUM_OF_WORKERS; i++ {
		wg.Add(1)
		go fetcherWorker(i, &seenUrlCache, &wg, urlQueue, toParseQueue, &pending)
	}

	// start parser work pool
	for i := uint8(0); i < c.NUM_OF_WORKERS; i++ {
		wg.Add(1)
		go parserWorker(i, &seenUrlCache, &wg, urlQueue, toParseQueue, pageQueue, allowExternal, &pending)
	}

	// start manager routine
	go func() {
		for {
			if atomic.LoadInt64(&pending) == 0 {
				close(urlQueue)
				close(toParseQueue)
				close(pageQueue)
				return
			}
			time.Sleep(500 * time.Millisecond) // polling rate, twice a second
		}
	}()

	wg.Wait() // wait for all routines to finish
	println("============ FINISHED CRAWLING ============")
	return nil
}

// fetcher worker: gets URLs from urlQueue, fetches HTML, sends to toParseQueue
func fetcherWorker(workerId uint8, seenUrlCache *sync.Map, wg *sync.WaitGroup, urlQueue chan string, toParseQueue chan *FetchedWebData, pending *int64) {
	defer wg.Done()

	for url := range urlQueue {
		fmt.Printf("fetcher-%d: acquired url: %s\n", workerId, url)

		if _, seen := seenUrlCache.LoadOrStore(url, true); seen {
			atomic.AddInt64(pending, -1) // done with this URL, already seen
			continue
		}

		htmlBody, err := parser.ParseSite(url)
		if err != nil {
			atomic.AddInt64(pending, -1) // cant parse, end
			continue
		}

		toParseQueue <- &FetchedWebData{
			URL:      url,
			Domain:   parser.GetDomain(url),
			HTMLBody: htmlBody,
		}
		time.Sleep(1 * time.Second)
	}
}

// parser worker: parses HTML, produces pages and new URLs
func parserWorker(workerId uint8, seenUrlCache *sync.Map, wg *sync.WaitGroup, urlQueue chan string, toParseQueue chan *FetchedWebData,
	pageQueue chan *repository.Page, allowExternal bool, pending *int64) {
	defer wg.Done()

	for fetchedData := range toParseQueue {
		fmt.Printf("parser-%d: acquired fetchedData\n", workerId)

		textData, links, title, err := parser.GetTextAndLinks(fetchedData.HTMLBody, fetchedData.Domain)
		if err != nil {
			atomic.AddInt64(pending, -1)
			continue
		}

		cleanedTextData := parser.CleanText(strings.Join(textData, " "))
		page := repository.Page{
			Title:   title,
			URL:     fetchedData.URL,
			Content: cleanedTextData,
		}

		pageQueue <- &page

		// Add validated links to the queue
		validatedLinks := parser.ValidateLinks(links, fetchedData.URL, fetchedData.Domain, allowExternal)
		for _, link := range validatedLinks {
			if _, seen := seenUrlCache.LoadOrStore(link, true); !seen {
				atomic.AddInt64(pending, 1)
				urlQueue <- link
			}
		}

		// Finished processing this page
		atomic.AddInt64(pending, -1)
	}
}

// DB worker: single goroutine inserts pages into the database
func databaseInteractionWorker(wg *sync.WaitGroup, pageQueue chan *repository.Page, repo *repository.PagesRepository) {
	defer wg.Done()
	for page := range pageQueue {
		if err := repo.InsertPage(*page); err != nil {
			println("Could not insert into page")
		}
	}
}

// simple passthrough, sqlite does the heavy lifting here
func (c *WebCrawler) SearchCrawled(phrase string, limit int) []repository.Page {
	pages := c.repo.SearchPages(phrase, limit)
	fmt.Printf("query returned %d number of rows", len(pages))
	return pages
}
