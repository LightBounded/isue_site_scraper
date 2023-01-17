package main

import (
	"fmt"
	"strings"
	"encoding/json"
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
			var publication Publication

			li := ul.Eq(i)
			publicationText := standardizeSpaces(li.Text())
			authorsString := strings.Split(publicationText, "\"")[0]
			publicationAuthors := strings.Split(authorsString, "., ")
			publicationYear := strings.Split(publicationText, " ")[len(strings.Split(publicationText, " "))-1]
			publicationMonth := strings.Split(publicationText, " ")[len(strings.Split(publicationText, " "))-2]
			publicationDate := publicationMonth + " " + publicationYear

			var publicationTitle string

			if len(strings.Split(authorsString, "\"")) != 1 {
				publicationTitle = strings.Split(authorsString, "\"")[1]
			}

			var publicationURL string
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

		fmt.Println(string(p))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit("https://www.eecs.ucf.edu/isuelab/publications/")

}
