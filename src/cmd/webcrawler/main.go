/*
This file is the main entrypoint of the crawler, calls other packages
to handle webcrawling given parameters
*/

package main

import server "github.com/firozt/crawler/src/internal/Server"

func main() {
	var HOSTNAME string = "localhost"
	var PORT string = "8080"
	// crawlSite("https://en.wikipedia.org/wiki/Chair")
	// db := repository.InitDB() // creates db conn
	// defer db.Close()
	// crawlSite("https://books.toscrape.com/index.html")
	// pageRepo := repository.NewPagesRepository(db) // creates DAO for pages table
	server.InitWebServer(HOSTNAME, PORT)

}

// }

// func crawlSite(url string) {
// 	// initial link queue
// 	queue := ThreadSafeQueue.NewThreadSafeQueue[string]()
// 	body := parser.ParseSite(url)
// 	text, links := parser.GetTextAndLinks(body)
// 	fmt.Println(len(text))
// 	for _, link := range links {
// 		queue.Enqueue(link)
// 	}
// 	fmt.Println("Links:\n", links)

// 	// setup worker pool
// 	var wg sync.WaitGroup
// 	numWorkers := 5

// 	worker := func(id int) {
// 		defer wg.Done()
// 		for j := 0; j < 2; j++ {
// 			link, ok := queue.Dequeue()
// 			if !ok {
// 				// q empty
// 				// TODO, wait for links
// 				return
// 			}
// 			fmt.Printf("Worker %d processing %s\n", id, link)
// 			body := parser.ParseSite(link)
// 			_, newlinks := parser.GetTextAndLinks(body)
// 			fmt.Println("Links:\n", newlinks)
// 		}
// 	}

// 	// Start workers
// 	for i := 0; i < numWorkers; i++ {
// 		wg.Add(1)
// 		go worker(i)
// 		time.Sleep(500000000)
// 	}
// 	wg.Wait()

// }
