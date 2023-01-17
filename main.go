package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

// Trims excessive spaces from a string
func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

type Publication struct {
	Title   string   `json:"title"`
	Url     string   `json:"url"`
	Authors []string `json:"authors"`
	Type    string   `json:"type"`
	Date    string   `json:"date"`
}

func main() {
	publications := []Publication{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.eecs.ucf.edu"),
	)

	c.OnHTML("div#publications>div.row:not(.year)", func(e *colly.HTMLElement) {
		selection := e.DOM

		publicationType := strings.TrimSpace(selection.Find("div.slot-6").Text())
		ul := selection.Find("div.slot-7-8-9>ul").Children()
		for i := 0; i < ul.Length(); i++ {
			li := ul.Eq(i)
			publicationText := standardizeSpaces(li.Text())

			var publication Publication

			var publicationTitle string
			var publicationURL string

			authorsString := strings.Split(publicationText, "\"")[0]
			publicationAuthors := strings.Split(authorsString, "., ")

			publicationYear := strings.Split(publicationText, " ")[len(strings.Split(publicationText, " "))-1]
			publicationMonth := strings.Split(publicationText, " ")[len(strings.Split(publicationText, " "))-2]
			publicationDate := publicationMonth + " " + publicationYear

			if len(strings.Split(publicationText, "\"")) > 2 {
				publicationTitle = strings.Split(publicationText, "\"")[1]
			}

			value, exists := li.Find("a").Attr("href")
			if exists {
				publicationURL = value
			}

			publication = Publication{
				Type:    publicationType,
				Url:     publicationURL,
				Authors: publicationAuthors,
				Title:   publicationTitle,
				Date:    publicationDate,
			}

			publications = append(publications, publication)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)

		p, err := json.Marshal(publications)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = os.WriteFile("publications.json", p, 0644)
		if err != nil {
			log.Fatal(err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit("https://www.eecs.ucf.edu/isuelab/publications/")

}
