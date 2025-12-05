package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

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

// Middleware wrapper to handle endpoint logging on each request made,
// returns handlerfunction (endpoint function)
func (s *Server) MiddleWare(method string, next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			fmt.Printf("Invalid method type expected %v got %v", method, r.Method)
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		fmt.Printf("[%s] %s %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
		next(w, r)
	}
}

// function to start the server, running on given host and port
func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.MiddleWare("GET", s.HandleRoot))
	mux.HandleFunc("/api/v1/crawl", s.MiddleWare("POST", s.StartCrawl))
	mux.HandleFunc("/api/v1/search", s.MiddleWare("GET", s.SearchCrawled))

	fmt.Printf("Server listening to %v:%v\n", s.hostname, s.port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", s.hostname, s.port), mux)
}

// ==================== ENDPOINTS ==================== //

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

func (s *Server) SearchCrawled(resp http.ResponseWriter, req *http.Request) {
	// parse query parameters
	query := req.URL.Query().Get("q")
	limitStr := req.URL.Query().Get("limit")

	if query == "" {
		http.Error(resp, "missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	// default limit
	limit := 10
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		} else {
			http.Error(resp, "invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	// fetch results from crawler repository
	results := s.crawler.SearchCrawled(query, limit)

	// encode results as JSON
	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(results); err != nil {
		http.Error(resp, "failed to encode response", http.StatusInternalServerError)
	}
}
