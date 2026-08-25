package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Callumogr18/Gocast/display"
	"github.com/Callumogr18/Gocast/htmx/server"
	XMLprocessing "github.com/Callumogr18/Gocast/weather_tool/XML-Processing"
)

func main() {
	var web_choice string

	fmt.Print("WebUI Testing (y/n) > ")
	fmt.Scan(&web_choice)

	if web_choice == "y" {
		router := server.NewRouter()

		fmt.Println("Serving on http://localhost:8080")
		if err := http.ListenAndServe(":8080", router); err != nil {
			fmt.Println("Error running server:", err)
			os.Exit(1)
		}
		return
	} else {
		choice, err := display.RunSelector()
		if err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

		if choice == "" {
			fmt.Println("No selection made.")
			return
		}

		XMLprocessing.FetchXML(choice)
	}
}
