package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync/atomic"

	webcrawler "github.com/firozt/crawler/src/internal/WebCrawler"
)

type ServerStats struct {
	totalCrawledEndpointsServed uint64
	totalCrawledPages           uint64
	totalSearchEndpointServed   uint64
}

type Server struct {
	hostname       string
	port           string
	crawler        *webcrawler.WebCrawler
	allowedOrigins *map[string]bool // domain allowed set
	stats          *ServerStats
}

func NewServer(crawler *webcrawler.WebCrawler, hostname string, port string, allowedOrigins *map[string]bool) *Server {
	return &Server{
		crawler:        crawler,
		hostname:       hostname,
		port:           port,
		allowedOrigins: allowedOrigins,
		stats: &ServerStats{
			totalCrawledEndpointsServed: *new(uint64),
			totalCrawledPages:           *new(uint64),
			totalSearchEndpointServed:   *new(uint64),
		},
	}
}

// Middleware wrapper to handle endpoint logging on each request made,
// also handles CORS header checks
// returns handlerfunction (endpoint function)
func (s *Server) MiddleWare(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		originHeader := r.Header.Get("Origin")

		// set CORS headers for all responses if origin exists
		if originHeader != "" {
			if !(*s.allowedOrigins)[originHeader] {
				fmt.Println("Reqeust was blocked due to CORS")
				w.Header().Set("Access-Control-Allow-Origin", originHeader)
				http.Error(w, fmt.Sprintf("origin %s not allowed", originHeader), http.StatusMethodNotAllowed)
				return
			}
			w.Header().Set("Access-Control-Allow-Origin", originHeader)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		} else {
			// for non-browser requests allow everything, maybe turn off?
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		// handle OPTIONS preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		// check method
		if r.Method != method {
			http.Error(w, fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}

// function to start the server, running on given host and port
func (s *Server) Run() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/crawl", s.MiddleWare("POST", s.StartCrawl))
	mux.HandleFunc("/api/v1/search", s.MiddleWare("GET", s.SearchCrawled))
	mux.HandleFunc("/api/v1/health", s.MiddleWare("GET", s.Health))
	mux.HandleFunc("/", s.CatchAll)
	fmt.Printf("Server listening to %v:%v\n", s.hostname, s.port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", s.hostname, s.port), mux)
}

// ==================== ENDPOINTS ==================== //

func (s *Server) CatchAll(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusNotFound)

	json.NewEncoder(resp).Encode(map[string]any{
		"error": map[string]any{
			"code":    "ENDPOINT_NOT_FOUND",
			"message": "The requested endpoint does not exist.",
			"method":  req.Method,
			"path":    req.URL.Path,
		},
	})
}

func (s *Server) StartCrawl(resp http.ResponseWriter, req *http.Request) {
	if atomic.LoadUint64(&s.stats.totalCrawledPages) > 150000 { // roughly 3rps for 15 hours worth of straight requests
		http.Error(resp, "The server is currently under too much load. Please try again at a later date", 429)
	}
	type StartCrawlBody struct {
		URL            string `json:"url"`
		MaxDepth       uint8  `json:"maxDepth"`
		FollowExternal bool   `json:"followExternal"`
	}

	// read body config
	body, err := io.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		http.Error(resp, "failed to read body of request", http.StatusBadRequest)
		return
	}

	var config StartCrawlBody
	if err := json.Unmarshal(body, &config); err != nil {
		http.Error(resp, "Malfored body", http.StatusBadRequest)
		return
	}

	pagesCrawled, err := s.crawler.StartCrawl(config.URL, config.FollowExternal)
	if err != nil {
		fmt.Println("Crawl failed: ", err)
		http.Error(resp, fmt.Sprintf("crawl failed: %v", err), http.StatusInternalServerError)
		return
	}

	// add stats
	atomic.AddUint64(&s.stats.totalCrawledEndpointsServed, 1)
	atomic.AddUint64(&s.stats.totalCrawledPages, pagesCrawled)

	// produce return
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(map[string]any{
		"status":       "completed",
		"pagesCrawled": strconv.Itoa(int(pagesCrawled)),
	})

}

func (s *Server) SearchCrawled(resp http.ResponseWriter, req *http.Request) {
	// add stats
	atomic.AddUint64(&s.stats.totalSearchEndpointServed, 1)

	// parse query parameters
	query := req.URL.Query().Get("q")
	limitStr := req.URL.Query().Get("limit")
	domain := req.URL.Query().Get("domain")
	if query == "" {
		http.Error(resp, "missing query parameter 'q'", http.StatusBadRequest)
		return
	}

	if domain == "" {
		http.Error(resp, "wildcard domain not allowed", http.StatusBadRequest)
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
	results := s.crawler.SearchCrawled(query, domain, limit)

	// encode results as JSON
	resp.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(resp).Encode(results); err != nil {
		http.Error(resp, "failed to encode response", http.StatusInternalServerError)
	}
}

// health check for the server, shows stats for server
func (s *Server) Health(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	json.NewEncoder(resp).Encode(map[string]any{
		"health":                     "OK",
		"totalPagesCrawled":          s.stats.totalCrawledPages,
		"totalCrawlEndpointsServed":  s.stats.totalCrawledEndpointsServed,
		"totalSearchEndpointsServed": s.stats.totalSearchEndpointServed,
	})
}
