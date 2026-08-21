package weathertool

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// fetchAndRead performs an HTTP GET request and returns the response body as bytes.
// It handles error checking and body closing automatically.
func fetchAndRead(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %s", resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func FetchXML(choice int) {
	var location string

	switch choice {
	case 1:
		bytes, err := fetchAndRead("https://www.met.ie/Open_Data/xml/xNational.xml")
		if err != nil {
			fmt.Println("failed to fetch national forecast:", err)
			return
		}

		var data NationalData
		if err := xml.Unmarshal(bytes, &data); err != nil {
			fmt.Println("failed to parse national forecast:", err)
			return
		}
		fmt.Printf("%+v\n", data)

	case 2:
		fmt.Print("Enter Province: ")
		fmt.Scan(&location)

		url := fmt.Sprintf("https://www.met.ie/Open_Data/xml/x%s.xml", location)
		bytes, err := fetchAndRead(url)
		if err != nil {
			fmt.Println("failed to fetch regional forecast:", err)
			return
		}

		var data ProvinceData
		if err := xml.Unmarshal(bytes, &data); err != nil {
			fmt.Println("failed to parse regional forecast:", err)
			return
		}
		fmt.Printf("%+v\n", data)

	case 3:
		fmt.Print("Enter County: ")
		fmt.Scan(&location)

		bytes, err := fetchAndRead("https://www.met.ie/Open_Data/xml/county_forecast.xml")
		if err != nil {
			fmt.Println("failed to fetch county forecast:", err)
			return
		}

		var data CountyForecast
		if err := xml.Unmarshal(bytes, &data); err != nil {
			fmt.Println("failed to parse county forecast:", err)
			return
		}

		for _, county := range data.Data {
			if strings.EqualFold(county.Name, location) {
				fmt.Printf("%+v\n", county)
				return
			}
		}
		fmt.Printf("county %q not found\n", location)
	}
}
