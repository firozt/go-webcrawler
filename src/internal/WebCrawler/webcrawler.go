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

var MAX_ADDED_LINKS uint8 = 5
var NUM_OF_WORKERS uint8 = 5

type WebCrawler struct {
	repo *repository.PagesRepository
}

func NewCrawler(repo *repository.PagesRepository) *WebCrawler {
	return &WebCrawler{repo: repo}
}

// starts the crawling proces on a url
func (c *WebCrawler) StartCrawl(url string) error {
	println("STARTING CRAWL")

	q := TSQ.NewThreadSafeQueue[string]()

	err := c.handlePage(url, q)
	if err != nil {
		fmt.Printf("ERROR: Could not scrape %v\nError: %v\n", url, err)
	}
	println("ALL CRAWLING DONE")

	return nil
	// var wg sync.WaitGroup

	// var i uint8 = 0
	// for i < NUM_OF_WORKERS || q.Len() < 1 {
	// 	wg.Add(1)
	// 	go workerAction(c, q, &wg)
	// }

	// wg.Wait()
	// println("ALL CRAWLING DONE")
	// return nil
}

// function that keeps parsing and saving the start of the queue
func workerAction(c *WebCrawler, q *TSQ.ThreadSafeQueue[string], wg *sync.WaitGroup) {
	for i := 0; i < 2; i++ {
		url, ok := q.Dequeue()
		if !ok {
			break
		}
		c.handlePage(url, q)
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

func (c *WebCrawler) handlePage(url string, q *TSQ.ThreadSafeQueue[string]) error {
	q.Enqueue(url)
	htmlBody, err := parser.ParseSite(url)
	if err != nil {
		return err
	}

	text, links := parser.GetTextAndLinks(htmlBody)
	// println()
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

	// enqueue links found
	for _, link := range links {
		q.Enqueue(link)
	}

	return nil
}
