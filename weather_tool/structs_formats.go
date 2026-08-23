package weathertool

/*
structs_formats.go:
Data formats for the XML data extracted from the URLS
*/

import "strings"

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
