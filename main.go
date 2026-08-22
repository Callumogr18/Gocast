package main

import (
	"fmt"
	"os"

	"github.com/Callumogr18/Gocast/display"
	weathertool "github.com/Callumogr18/Gocast/weather_tool"
)

func main() {
	choice, err := display.RunSelector()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if choice == "" {
		fmt.Println("No selection made.")
		return
	}

	weathertool.FetchXML(choice)
}
