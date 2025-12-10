package webcrawler

import (
	"errors"
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
func (c *WebCrawler) StartCrawl(url string, allowExternal bool) error {
	url = parser.NormalizeURL(url)
	if _, err := parser.ParseSite(url); err != nil {
		return errors.New("initial url is not a valid site")
	}

	println("STARTING CRAWL")
	q := TSQ.NewThreadSafeQueue[string]()
	q.Enqueue(url)
	q.Dequeue()
	links := c.handlePage(url, false) // never do external on first page, allows more related links
	fmt.Printf("INTIAL LINKS %s", links)
	for _, link := range links {
		q.Enqueue(parser.NormalizeURL(parser.RemoveFragment(link)))
	}
	var wg sync.WaitGroup

	var i uint8 = 0
	for i < c.NUM_OF_WORKERS && !(q.Len() < 1) {
		wg.Add(1)
		go workerAction(c, q, &wg, allowExternal)
		i++
	}

	wg.Wait()
	println("ALL CRAWLING DONE")
	return nil
}

// function that keeps parsing and saving the start of the queue
func workerAction(c *WebCrawler, q *TSQ.ThreadSafeQueue[string], wg *sync.WaitGroup, allowExternal bool) {
	for i := 0; i < 5; i++ {
		url, ok := q.Dequeue()
		println("CHECKING ", q.Len(), url)

		if !ok {
			break
		}
		links := c.handlePage(url, allowExternal)

		unseenLinks := 0
		// keep adding from links until we enq N unseen links or we reached the end of the link list
		for i := 0; unseenLinks < int(c.MAX_ADDED_LINKS_PER_PAGE) && i < len(links); i++ {
			if q.Enqueue(parser.RemoveFragment(links[i])) {
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

func (c *WebCrawler) handlePage(url string, allowExternal bool) []string {
	htmlBody, err := parser.ParseSite(url)
	var links []string
	if err != nil {
		return links
	}
	domain := parser.GetDomain(url)
	text, links, title, err := parser.GetTextAndLinks(htmlBody, domain)
	if err != nil {
		return []string{}
	}

	cleaned_text := parser.CleanText(strings.Join(text, " "))
	page := repository.Page{
		Title:   title,
		URL:     url,
		Content: cleaned_text,
	}
	// save to database
	if err := c.repo.InsertPage(page); err != nil {
		return links
	}
	links = parser.ValidateLinks(links, url, domain, allowExternal)

	return links
}

// simple passthrough, sqlite does the heavy lifting here
func (c *WebCrawler) SearchCrawled(phrase string, limit int) []repository.Page {
	pages := c.repo.SearchPages(phrase, limit)
	fmt.Printf("query returned %d number of rows", len(pages))
	return pages
}
