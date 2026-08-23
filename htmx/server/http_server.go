package server

import (
	"html/template"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("htmx/html/home.html"))

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewRouter builds the HTTP router for the web UI.
func NewRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("GET /{$}", homeHandler)
	return router
}
