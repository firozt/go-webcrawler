/*
This file is the main entrypoint of the crawler, calls other packages
to handle webcrawling given parameters
*/

package main

import (
	repository "github.com/firozt/crawler/src/internal/Repository"
	server "github.com/firozt/crawler/src/internal/Server"
	webcrawler "github.com/firozt/crawler/src/internal/WebCrawler"
)

func main() {
	// config variables
	var HOSTNAME string = "localhost"
	var PORT string = "8080"
	var MAX_ADDED_LINKS_PER_PAGE uint8 = 3
	var NUM_OF_WORKERS uint8 = 5

	// starting everything
	db := repository.InitDB() // creates db conn and obj
	defer db.Close()
	pagesRepo := repository.NewPagesRepository(db)                                           // creates pagesRepo API using DB
	webcrawler := webcrawler.NewCrawler(pagesRepo, MAX_ADDED_LINKS_PER_PAGE, NUM_OF_WORKERS) // crawls sites and saves to DB
	if true {
		server := server.NewServer(webcrawler, HOSTNAME, PORT) // creates webserver instance
		server.Run()                                           // runs server
	} else {
		webcrawler.StartCrawl("https://books.toscrape.com/index.html")
	}
}
