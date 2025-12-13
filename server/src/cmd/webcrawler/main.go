/*
This file is the main entrypoint of the crawler, calls other packages
to handle webcrawling given parameters
*/

package main

import (
	"fmt"
	"os"

	repository "github.com/firozt/crawler/src/internal/Repository"
	server "github.com/firozt/crawler/src/internal/Server"
	webcrawler "github.com/firozt/crawler/src/internal/WebCrawler"
	"github.com/joho/godotenv"
)

var DEV_MODE = false

func main() {
	// LOAD ENV
	var err error = godotenv.Load()
	if err == nil && os.Getenv("PORT") == "1" {
		DEV_MODE = true
	}

	// config variables
	fmt.Printf("STARTING MAIN, DEV_MODE: %v\n", DEV_MODE)

	var allowedOrigins map[string]bool = map[string]bool{
		"https://domainsearch.ramizabdulla.me/": true,
	}
	if DEV_MODE {
		allowedOrigins["http://localhost:5173"] = true
	}
	var HOSTNAME string = "0.0.0.0"
	var PORT string = "8080"
	var MAX_ADDED_LINKS_PER_PAGE uint8 = 255
	var NUM_OF_WORKERS uint8 = 2
	var MAX_UNIQUE_CRAWLED_PAGES uint64 = 25

	// starting everything
	db := repository.InitDB() // creates db conn and obj
	defer db.Close()
	pagesRepo := repository.NewPagesRepository(db)                                                                     // creates pagesRepo API using DB
	webcrawler := webcrawler.NewCrawler(pagesRepo, MAX_ADDED_LINKS_PER_PAGE, NUM_OF_WORKERS, MAX_UNIQUE_CRAWLED_PAGES) // crawls sites and saves to DB
	server := server.NewServer(webcrawler, HOSTNAME, PORT, &allowedOrigins)                                            // creates webserver instance
	server.Run()
	println("APPLICATION TERMINATED")
	// runs server
}
