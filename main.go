package main

import (
	"fmt"

	"github.com/LightBounded/lab-site-scraper/collectors"
)

func main() {
	fmt.Print("What would you like to scrape? (1) Publications (2) People (3) All: ")

	var choice int
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		collectors.GetPublications()
	case 2:
	case 3:	
	}
}
