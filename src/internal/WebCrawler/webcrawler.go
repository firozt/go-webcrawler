package webcrawler

import (
	"fmt"
	"strings"
	"sync"

	parser "github.com/firozt/crawler/src/internal/Parser"
	repository "github.com/firozt/crawler/src/internal/Repository"
	TSQ "github.com/firozt/crawler/src/internal/ThreadSafeQueue"
)

var MAX_ADDED_LINKS uint8 = 5
var NUM_OF_WORKERS uint8 = 5

type WebCrawler struct {
	repo *repository.PagesRepository
}

func NewCrawler(repo *repository.PagesRepository) *WebCrawler {
	return &WebCrawler{repo: repo}
}

func (c *WebCrawler) StartCrawl(url string) error {
	q := TSQ.NewThreadSafeQueue[string]()

	err := c.handlePage(url, q)
	if err != nil {
		fmt.Printf("ERROR: Could not scrape %v\nError: %v\n", url, err)
	}

	var wg sync.WaitGroup

	var i uint8 = 0
	for i < NUM_OF_WORKERS || q.Len() < 1 {
		wg.Add(1)
		go workerAction(c, q, &wg)
	}

	wg.Wait()
	return nil
}

func workerAction(c *WebCrawler, q *TSQ.ThreadSafeQueue[string], wg *sync.WaitGroup) {
	for ok := true; ok; {
		url, ok := q.Dequeue()
		if !ok {
			break
		}
		c.handlePage(url, q)
	}
	wg.Done()
}

func (c *WebCrawler) handlePage(url string, q *TSQ.ThreadSafeQueue[string]) error {
	q.Enqueue(url)
	htmlBody, err := parser.ParseSite(url)
	if err != nil {
		return err
	}

	text, links := parser.GetTextAndLinks(htmlBody)

	page := repository.Page{
		Title:   "TODO GET TITLE",
		URL:     url,
		Content: strings.Join(text, " "),
	}
	// save to database
	if err := c.repo.InsertPage(page); err != nil {
		return err
	}
	if len(links) > int(MAX_ADDED_LINKS) {
		links = links[:MAX_ADDED_LINKS]
	}

	for _, link := range links {
		q.Enqueue(link)
	}

	return nil
}
