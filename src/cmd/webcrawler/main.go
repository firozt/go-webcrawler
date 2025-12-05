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
	var HOSTNAME string = "localhost"
	var PORT string = "8080"
	db := repository.InitDB() // creates db conn and obj
	defer db.Close()
	pagesRepo := repository.NewPagesRepository(db) // creates pagesRepo API using DB
	webcrawler := webcrawler.NewCrawler(pagesRepo) // crawls sites and saves to DB
	webcrawler.StartCrawl("https://books.toscrape.com/index.html")
	if false {
		server := server.NewServer(webcrawler, HOSTNAME, PORT) // creates webserver instance
		server.Run()                                           // runs server
	}
}
