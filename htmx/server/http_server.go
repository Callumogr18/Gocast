package server

/*
http_server.go:
Handles the routing for the application with the 'NewRouter()' function
which returns a HTTP request multiplexer
*/

import (
	"net/http"
)

// NewRouter builds the HTTP router for the web UI.
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("GET /{$}", homeHandler)
	router.HandleFunc("GET /national", nationalHandler)
	router.HandleFunc("GET /regional", regionalHandler)
	router.HandleFunc("GET /county", countyHandler)
	router.HandleFunc("GET /warning", warningHandler)

	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("htmx/static"))))

	return router
}
