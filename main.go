package main

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

type Publication struct {
	title   string
	authors []string
	date    string
}

// Trims excessive spaces from a string
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func getPublicationFromString(s string) Publication {
	authorsString := strings.Split(s, "\"")[0]
	authors := strings.Split(authorsString, "., ")

	title := strings.Split(s, "\"")[1]
	year := strings.Split(s, " ")[len(strings.Split(s, " "))-1]
	month := strings.Split(s, " ")[len(strings.Split(s, " "))-2]
	
	date := month + " " + year
	return Publication{
		title:   title,
		authors: authors,
		date:    date,
	}
}

func main() {
	publications := []Publication{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.eecs.ucf.edu"),
	)

	c.OnHTML("ul[class='noicon loose']>li", func(e *colly.HTMLElement) {
		sanitizedText := standardizeSpaces(e.Text)

		publication := getPublicationFromString(sanitizedText)
		println(publication.title)
		publications = append(publications, publication)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit("https://www.eecs.ucf.edu/isuelab/publications/")
}
