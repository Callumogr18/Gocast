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

TO-DO:
Add a helper function to fetch the weather data and parse the neccessary .html file
*/

import (
	"bytes"
	"errors"
	"html/template"
	"net/http"

	weathertool "github.com/Callumogr18/Gocast/weather_tool"
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
	data, err := weathertool.FetchNational()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	renderTemplate(w, "htmx/html/national.html", data)
}

func regionalHandler(w http.ResponseWriter, r *http.Request) {
	province := r.URL.Query().Get("province")

	data, err := weathertool.FetchRegional(province)
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

	data, err := weathertool.FetchCounty(county)
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
