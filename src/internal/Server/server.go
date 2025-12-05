package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	webcrawler "github.com/firozt/crawler/src/internal/WebCrawler"
)

type Server struct {
	hostname string
	port     string
	crawler  *webcrawler.WebCrawler
}

func NewServer(crawler *webcrawler.WebCrawler, hostname string, port string) *Server {
	return &Server{
		crawler:  crawler,
		hostname: hostname,
		port:     port,
	}
}

func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.HandleRoot)
	mux.HandleFunc("/api/v1/crawl", s.StartCrawl)

	fmt.Printf("Server listening to %v:%v\n", s.hostname, s.port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", s.hostname, s.port), mux)
}

func (s *Server) HandleRoot(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello World")
}

type StartCrawlBody struct {
	URL            string `json:"url"`
	MaxDepth       uint8  `json:"maxDepth"`
	FollowExternal bool   `json:"followExternal"`
}

func (s *Server) StartCrawl(resp http.ResponseWriter, req *http.Request) {
	// read body config
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(resp, "failed to read body of request", http.StatusBadRequest)
		return
	}

	var config StartCrawlBody
	if err := json.Unmarshal(body, &config); err != nil {
		http.Error(resp, "invalid json", http.StatusBadRequest)
		return
	}

	err = s.crawler.StartCrawl(config.URL)
	if err != nil {
		http.Error(resp, fmt.Sprintf("crawl failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte("Crawl finished"))
}
