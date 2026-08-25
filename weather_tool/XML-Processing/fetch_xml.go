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
	"fmt"
	"net/http"
	"strings"
	"time"

	weathertool "github.com/Callumogr18/Gocast/weather_tool"
)

const xmlBaseURL = "https://www.met.ie/Open_Data/xml/x%s.xml"
const countyForecastURL = "https://www.met.ie/Open_Data/xml/county_forecast.xml"
const warningURL = "https://www.met.ie/Open_Data/json/warning_IRELAND.json"

var httpClient = &http.Client{Timeout: 10 * time.Second}

// FetchNational fetches and parses the national forecast.
func FetchNational() (weathertool.NationalData, error) {
	return fetchAndParse[weathertool.NationalData](fmt.Sprintf(xmlBaseURL, "National"), "national")
}

// FetchRegional fetches and parses the regional forecast
func FetchRegional(s string) (weathertool.ProvinceData, error) {
	return fetchLocation[weathertool.ProvinceData](s, "regional", weathertool.ParseProvince)
}

// FetchCounty fetches the all-counties forecast and returns the entry
// matching s. Matching is case-insensitive since the feed names are
// UPPERCASE (e.g. "DUBLIN").
func FetchCounty(s string) (weathertool.County, error) {
	doc, err := fetchAndParse[weathertool.CountyForecast](countyForecastURL, "county")
	if err != nil {
		return weathertool.County{}, err
	}

	s = strings.TrimSpace(s)
	for _, c := range doc.Counties {
		if strings.EqualFold(c.Name, s) {
			return c, nil
		}
	}

	return weathertool.County{}, fmt.Errorf("%w: %q", weathertool.ErrUnknownLocation, s)
}

// FetchWarning extracts the warning forecast alerts
func FetchWarning() ([]weathertool.Warning, error) {
	return fetchAndParseJSON[[]weathertool.Warning](warningURL, "warning")
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
