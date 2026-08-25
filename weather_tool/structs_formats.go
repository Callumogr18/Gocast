package weathertool

/*
structs_formats.go:
Data formats for the XML data extracted from the URLS
*/

import "strings"

/*
=========================
National forecast structs
=========================
*/
type NationalData struct {
	Issued struct {
		Time string `xml:"issued,attr"`
	} `xml:"issued"`
	Today   string `xml:"today"`
	Outlook string `xml:"outlook"`
}

type ProvinceData struct {
	Issued struct {
		Time string `xml:"issued,attr"`
	} `xml:"issued"`
	Tomorrow string `xml:"tomorrow"`
	Pollen   string `xml:"pollen"`
}

/*
=======================
County forecast structs
=======================
*/
type CountyForecast struct {
	Issued   string   `xml:"issued,attr"`
	Counties []County `xml:"county"`
}

type County struct {
	Name string `xml:"name"`
	Days []Day  `xml:"day"`
}

type Day struct {
	DayNum        int    `xml:"day_num"`
	Date          string `xml:"date"`
	MinTemp       string `xml:"min_temp"`
	MaxTemp       string `xml:"max_temp"`
	Weather       string `xml:"weather"`
	WindSpeed     string `xml:"wind_speed"`
	WindDir       string `xml:"wind_dir"`
	Rainfall6To18 string `xml:"rainfall_6_18"`
	Rainfall18To6 string `xml:"rainfall_18_6"`
}

func (d Day) Wind() string {
	return strings.TrimSpace(d.WindSpeed) + " km/h"
}

// Replaces sun_with_grey_clouds -> sun with grey clouds
func (d Day) WeatherLabel() string {
	return strings.ReplaceAll(d.Weather, "_", " ")
}

/*
=======================
Weather Warning structs
=======================
*/
type Warning struct {
	ID          int      `json:"id"`
	CapID       string   `json:"capId"`
	Type        string   `json:"type"`
	Severity    string   `json:"severity"`
	Certainty   string   `json:"certainty"`
	Level       string   `json:"level"`
	Issued      string   `json:"issued"`
	Updated     string   `json:"updated"`
	Onset       string   `json:"onset"`
	Expiry      string   `json:"expiry"`
	Headline    string   `json:"headline"`
	Description string   `json:"description"`
	Regions     []string `json:"regions"`
	Status      string   `json:"status"`
}

// LevelClass returns Level lowercased for use as a CSS class
// (warning-yellow / warning-orange / warning-red).
func (w Warning) LevelClass() string {
	return strings.ToLower(w.Level)
}
