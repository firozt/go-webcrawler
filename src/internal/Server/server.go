/*
This file handles all endpoint logic, and exposing endpoints for http requests
*/

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	parser "github.com/firozt/crawler/src/internal/Parser"
)

// initialise the http web server
func InitWebServer(HOSTNAME string, PORT string) {
	// initialise new server mux to handle traffic flow
	mux := http.NewServeMux()

	// ============= ENDPOINT MAPPINGS ============= //

	mux.HandleFunc("/", HandleRoot)
	mux.HandleFunc("/api/v1/crawl", StartCrawl)

	// ============================================= //

	fmt.Printf("Server listening to %v:%v\n", HOSTNAME, PORT)
	http.ListenAndServe(fmt.Sprintf("%v:%v", HOSTNAME, PORT), mux)
}

func HandleRoot(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello Word")
}

type StartCrawlBody struct {
	URL            string `json:"url"`
	MaxDepth       uint8  `json:"maxDepth"`
	FollowExternal bool   `json:"followExternal"`
}

func StartCrawl(resp http.ResponseWriter, req *http.Request) {
	// obtain config from body of request
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(resp, "failed to read body", http.StatusBadRequest)
		return
	}
	var config StartCrawlBody
	if err := json.Unmarshal(body, &config); err != nil {
		http.Error(resp, "invalid json", http.StatusBadRequest)
		return
	}

	var htmlResponse string = parser.ParseSite(config.URL)

}
