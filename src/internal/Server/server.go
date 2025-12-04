/*
This file handles all endpoint logic, and exposing endpoints for http requests
*/

package server

import (
	"fmt"
	"net/http"
)

// initialise the http web server
func InitWebServer(HOSTNAME string, PORT string) {
	// initialise new server mux to handle traffic flow
	mux := http.NewServeMux()

	// ============= ENDPOINT MAPPINGS ============= //

	mux.HandleFunc("/", HandleRoot)

	// ============================================= //

	fmt.Printf("Server listening to %v:%v\n", HOSTNAME, PORT)
	http.ListenAndServe(fmt.Sprintf("%v:%v", HOSTNAME, PORT), mux)
}

func HandleRoot(resp http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(resp, "Hello Word")
}
