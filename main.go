package main

import (
	"fmt"

	weathertool "github.com/Callumogr18/Gocast/weather_tool"
)

func main() {
	var workflow int

	//fmt.Println("Hello World")
	fmt.Println("Enter choice\n1. National\n2.Province\n3.County")
	fmt.Print(">>> ")

	fmt.Scan(&workflow)
	weathertool.FetchXML(workflow)

}
