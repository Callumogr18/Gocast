package weathertool

// TO-DO:
// =======================
// FILE NEEDS RE-FACTORING
// =======================

/*
fetch_xml.go:
This file has a legacy function which still serves the CLI in 'FetchXML()',
this parses the user choice in the UI and reflects the corresponsding options.

Three functions to fetch different data such as national, regional and county:
	* FetchNational(ctx) (NationalData, error)
	* FetchRegional(ctx, s) (ProvinceData, error)
	* FetchCounty(ctx, s)   (County, error)

These share a common fetch/parse shape via three helpers:
	* fetchAndRead(ctx, url string) ([]byte, error)
	* fetchAndParse[T any](ctx, url, kind string) (T, error)
	* fetchLocation[T any](ctx, loc, kind string, parse func(string) (string, error)) (T, error)

Every fetch takes a context.Context so a caller that goes away (a browser that
disconnects, a cancelled CLI run) cancels the upstream request instead of
leaving it in flight.
*/

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const xmlBaseURL = "https://www.met.ie/Open_Data/xml/x%s.xml"
const countyForecastURL = "https://www.met.ie/Open_Data/xml/county_forecast.xml"

var httpClient = &http.Client{Timeout: 10 * time.Second}

// fetchAndRead performs an HTTP GET request and returns the response body as bytes.
// It handles error checking and body closing automatically.
func fetchAndRead(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, &UpstreamError{URL: url, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &UpstreamError{URL: url, StatusCode: resp.StatusCode, Err: fmt.Errorf("unexpected status %s", resp.Status)}
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &UpstreamError{URL: url, Err: err}
	}

	return cleanString(bytes), nil
}

// fetchAndParse fetches the XML document at url and unmarshals it into T.
// kind names the forecast ("national", "regional", …) for error messages.
func fetchAndParse[T any](url, kind string) (T, error) {
	var data T

	body, err := fetchAndRead(url)
	if err != nil {
		return data, fmt.Errorf("failed to fetch %s forecast: %w", kind, err)
	}

	if err := xml.Unmarshal(body, &data); err != nil {
		return data, fmt.Errorf("failed to parse %s forecast: %w", kind, err)
	}

	return data, nil
}

// fetchLocation validates loc with parse, then fetches and parses the
// matching feed.
func fetchLocation[T any](loc, kind string, parse func(string) (string, error)) (T, error) {
	name, err := parse(loc)
	if err != nil {
		var zero T
		return zero, err
	}

	return fetchAndParse[T](fmt.Sprintf(xmlBaseURL, name), kind)
}

// FetchNational fetches and parses the national forecast.
func FetchNational() (NationalData, error) {
	return fetchAndParse[NationalData](fmt.Sprintf(xmlBaseURL, "National"), "national")
}

// FetchRegional fetches and parses the regional forecast
func FetchRegional(s string) (ProvinceData, error) {
	return fetchLocation[ProvinceData](s, "regional", ParseProvince)
}

// FetchCounty fetches the all-counties forecast and returns the entry
// matching s. Matching is case-insensitive since the feed names are
// UPPERCASE (e.g. "DUBLIN").
func FetchCounty(s string) (County, error) {
	doc, err := fetchAndParse[CountyForecast](countyForecastURL, "county")
	if err != nil {
		return County{}, err
	}

	s = strings.TrimSpace(s)
	for _, c := range doc.Counties {
		if strings.EqualFold(c.Name, s) {
			return c, nil
		}
	}

	return County{}, fmt.Errorf("%w: %q", ErrUnknownLocation, s)
}

/*
FetchXML is a function in place for the CLI option as of now, an
alternative will probably have to be made eventually for clean-up
*/
func FetchXML(choice string) {
	var location string

	switch choice {
	case "National":
		data, err := FetchNational()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", data)

	case "Regional":
		fmt.Print("Enter Province: ")
		fmt.Scan(&location)

		data, err := FetchRegional(location)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", data)

	case "County":
		fmt.Print("Enter County: ")
		fmt.Scan(&location)

		data, err := FetchCounty(location)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%+v\n", data)
	}

}
