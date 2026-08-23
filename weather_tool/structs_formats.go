package weathertool

import (
	"encoding/xml"
)

type NationalData struct {
	Issued  string `xml:"issued"`
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
	Name xml.Name     `xml:"forecast"`
	Data []CountyData `xml:"county"`
}

type CountyData struct {
	Name string      `xml:"name"`
	Days []DayResult `xml:"day"`
}

type DayResult struct {
	Date    string `xml:"date"`
	MinTemp int    `xml:"min_temp"`
}
