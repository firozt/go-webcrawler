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

func main() {
	var DEV_MODE = false
	// LOAD ENV
	var err error = godotenv.Load()

	// load devmode if .env has it to 1 else false (even if no .env)
	if err == nil && os.Getenv("DEVMODE") == "1" {
		DEV_MODE = true
	}

	// config variables
	fmt.Printf("STARTING MAIN, DEV_MODE: %v\n", DEV_MODE)

	var allowedOrigins map[string]bool = map[string]bool{
		"https://domainsearch.ramizabdulla.me/": true,
		"https://domainsearch.ramizabdulla.me":  true,
	}
	HOSTNAME := "0.0.0.0"
	if DEV_MODE {
		allowedOrigins["*"] = true
		HOSTNAME = "127.0.0.1"
	}
	var PORT string = "8080"
	var MAX_ADDED_LINKS_PER_PAGE uint8 = 255
	var NUM_OF_WORKERS uint8 = 2
	var MAX_UNIQUE_CRAWLED_PAGES uint64 = 25

	// starting everything
	db := repository.InitDB() // creates db conn and obj
	defer db.Close()
	pagesRepo := repository.NewPagesRepository(db)                                                                                // creates pagesRepo API using DB
	graphRepo := repository.NewGraphRepository(db)                                                                                // creates graphRepo API using DB
	webcrawler := webcrawler.NewCrawler(pagesRepo, graphRepo, MAX_ADDED_LINKS_PER_PAGE, NUM_OF_WORKERS, MAX_UNIQUE_CRAWLED_PAGES) // crawls sites and saves to DB
	server := server.NewServer(webcrawler, HOSTNAME, PORT, &allowedOrigins)                                                       // creates webserver instance
	server.Run()
	println("APPLICATION TERMINATED")
}
