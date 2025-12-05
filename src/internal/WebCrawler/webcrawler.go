package webcrawler

import (
	"fmt"
	"strings"
	"sync"
	"time"

	parser "github.com/firozt/crawler/src/internal/Parser"
	repository "github.com/firozt/crawler/src/internal/Repository"
	TSQ "github.com/firozt/crawler/src/internal/ThreadSafeQueue"
)

type WebCrawler struct {
	repo                     *repository.PagesRepository
	MAX_ADDED_LINKS_PER_PAGE uint8
	NUM_OF_WORKERS           uint8
}

func NewCrawler(repo *repository.PagesRepository, MAX_ADDED_LINKS_PER_PAGE uint8, NUM_OF_WORKERS uint8) *WebCrawler {
	return &WebCrawler{
		repo:                     repo,
		MAX_ADDED_LINKS_PER_PAGE: MAX_ADDED_LINKS_PER_PAGE,
		NUM_OF_WORKERS:           NUM_OF_WORKERS,
	}
}

// starts the crawling proces on a url
func (c *WebCrawler) StartCrawl(url string) error {
	println("STARTING CRAWL")

	q := TSQ.NewThreadSafeQueue[string]()
	//TEMP
	q.Enqueue(url)
	q.Dequeue()
	links, err := c.handlePage(url)

	if err != nil {
		fmt.Printf("ERROR: Could not scrape %v\nError: %v\n", url, err)
	}

	for _, link := range links {
		q.Enqueue(link)
	}
	var wg sync.WaitGroup

	var i uint8 = 0
	for i < c.NUM_OF_WORKERS || q.Len() < 1 {
		wg.Add(1)
		go workerAction(c, q, &wg)
		i++
	}

	wg.Wait()
	println("ALL CRAWLING DONE")
	return nil
}

// function that keeps parsing and saving the start of the queue
func workerAction(c *WebCrawler, q *TSQ.ThreadSafeQueue[string], wg *sync.WaitGroup) {
	for i := 0; i < 5; i++ {
		url, ok := q.Dequeue()
		println("CHECKING ", q.Len(), url)

		if !ok {
			break
		}
		links, err := c.handlePage(url)
		if err != nil {
			continue
		}

		unseenLinks := 0
		// keep adding from links until we enq N unseen links or we reached the end of the link list
		for i := 0; unseenLinks < int(c.MAX_ADDED_LINKS_PER_PAGE) && i >= len(links); i++ {
			if q.Enqueue(links[i]) {
				unseenLinks++
			}
		}

		time.Sleep(1 * time.Second) // wait a second so i dont get banned lol
	}
	// for ok := true; ok; {
	// 	url, ok := q.Dequeue()
	// 	if !ok {
	// 		break
	// 	}
	// 	c.handlePage(url, q)
	// 	time.Sleep(1 * time.Second) // wait a second so i dont get banned lol
	// }
	wg.Done()
}

func (c *WebCrawler) handlePage(url string) ([]string, error) {
	htmlBody, err := parser.ParseSite(url)
	var links []string
	if err != nil {
		return links, err
	}

	text, links := parser.GetTextAndLinks(htmlBody)

	cleaned_text := parser.CleanText(strings.Join(text, " "))
	page := repository.Page{
		Title:   "TODO GET TITLE",
		URL:     url,
		Content: cleaned_text,
	}
	// save to database
	if err := c.repo.InsertPage(page); err != nil {
		return links, err
	}
	links = parser.ValidateLinks(links, url)

	return links, err
}

// simple passthrough, sqlite does the heavy lifting here
func (c *WebCrawler) SearchCrawled(phrase string, limit int) []repository.Page {
	pages := c.repo.SearchPages(phrase, limit)
	return pages
}
