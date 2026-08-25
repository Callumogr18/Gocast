package server

/*
handlers.go:
Requests made to the backend (/national, /regional, /county) have their corresponding
handler functions fetch the weather forecast data and then render the .html file
for the GET request made
	* homehandler(w http.ResponseWriter, r *http.Request)
	* nationalHandler(w http.ResponseWriter, r *http.Request)
	* regionalHandler(w http.ResponseWriter, r *http.Request)
	* countyHandler(w http.ResponseWriter, r *http.Request)

A helper function 'renderTemplate()' parses the path and executes it against data into a buffer
first, so a template error surfaces as a clean 500 instead of a partial
body written straight to w
*/

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	weathertool "github.com/Callumogr18/Gocast/weather_tool"
	XMLprocessing "github.com/Callumogr18/Gocast/weather_tool/XML-Processing"
)

// renderTemplate parses path and executes it against data into a buffer
// first, so a template error surfaces as a clean 500 instead of a partial
// body written straight to w.
func renderTemplate(w http.ResponseWriter, path string, data any) {
	tmpl := template.Must(template.ParseFiles(path))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(buf.Bytes())
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "htmx/html/home.html", nil)
}

func nationalHandler(w http.ResponseWriter, r *http.Request) {
	data, err := XMLprocessing.FetchNational()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	renderTemplate(w, "htmx/html/national.html", data)
}

func regionalHandler(w http.ResponseWriter, r *http.Request) {
	province := r.URL.Query().Get("province")

	data, err := XMLprocessing.FetchRegional(province)
	if err != nil {
		if errors.Is(err, weathertool.ErrUnknownLocation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	renderTemplate(w, "htmx/html/regional.html", data)
}

func countyHandler(w http.ResponseWriter, r *http.Request) {
	county := r.URL.Query().Get("county")

	data, err := XMLprocessing.FetchCounty(county)
	if err != nil {
		if errors.Is(err, weathertool.ErrUnknownLocation) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	renderTemplate(w, "htmx/html/county.html", data)
}

func warningHandler(w http.ResponseWriter, r *http.Request) {
	data, err := XMLprocessing.FetchWarning()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	renderTemplate(w, "htmx/html/warning.html", data)
}
