package weathertool

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"

	weathertool "github.com/Callumogr18/Gocast/weather_tool"
)

func cleanString(b []byte) []byte {
	/*
		After adding this -> b = bytes.ReplaceAll(b, []byte("\n"), []byte(""))
		The output was still outputting escape sequence chars. Claude suggested
		adding the subsequent two lines, works!
	*/
	b = bytes.ReplaceAll(b, []byte("\n"), []byte(""))  // strip real newline bytes (pretty-printed XML)
	b = bytes.ReplaceAll(b, []byte(`\n`), []byte(" ")) // replace literal "\n" text markers with a space
	b = bytes.Join(bytes.Fields(b), []byte(" "))       // collapse any resulting runs of whitespace

	return b
}

// fetchAndRead performs an HTTP GET request and returns the response body as bytes.
// It handles error checking and body closing automatically.
func fetchAndRead(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, &weathertool.UpstreamError{URL: url, Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &weathertool.UpstreamError{URL: url, StatusCode: resp.StatusCode, Err: fmt.Errorf("unexpected status %s", resp.Status)}
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &weathertool.UpstreamError{URL: url, Err: err}
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

// This function is used to parse the JSON endpoint for weather
// warning. Leagacy weather warning endpoint was in traditional
// XML format, new endpoint is JSON format
func fetchAndParseJSON[T any](url string, kind string) (T, error) {
	var data T

	body, err := fetchAndRead(url)
	if err != nil {
		return data, fmt.Errorf("failed to parse %s forecast: %w", kind, err)
	}

	if err := json.Unmarshal(body, &data); err != nil {
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
