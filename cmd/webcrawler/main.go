package main

import (
	"fmt"
	"sync"
	"time"

	parser "github.com/firozt/crawler/internal/Parser"
	"github.com/firozt/crawler/internal/ThreadSafeQueue"
)

func main() {
	begin("https://en.wikipedia.org/wiki/Chair")
}

func begin(url string) {
	// initial link queue
	queue := ThreadSafeQueue.NewThreadSafeQueue[string]()
	body := parser.ParseSite(url)
	text, links := parser.GetTextAndLinks(body)
	fmt.Println(len(text))
	for _, link := range links {
		queue.Enqueue(link)
	}

	// setup worker pool
	var wg sync.WaitGroup
	numWorkers := 5

	worker := func(id int) {
		defer wg.Done()
		for j := 0; j < 2; j++ {
			link, ok := queue.Dequeue()
			if !ok {
				// q empty
				// TODO, wait for links
				return
			}
			fmt.Printf("Worker %d processing %s\n", id, link)
			parser.ParseSite(link)
		}
	}

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i)
		time.Sleep(1000000)
	}
	wg.Wait()
	queue.All()
}
