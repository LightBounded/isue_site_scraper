package collectors

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/LightBounded/lab-site-scraper/utils"

	"github.com/gocolly/colly"
)

type Publication struct {
	Title   string   `json:"title"`
	Url     string   `json:"url"`
	Authors []string `json:"authors"`
	Type    string   `json:"type"`
	Date    string   `json:"date"`
}

func cleanAuthors(authors []string) []string {
	newAuthors := []string{}

	for _, author := range authors {
		if len(author) == 0 {
			continue
		}

		author = strings.Trim(author, " ")
		author = strings.ReplaceAll(author, "and ", "")
		author = strings.ReplaceAll(author, ".", "")

		newAuthors = append(newAuthors, author)
	}
	return newAuthors
}

// TODO: Clean up this function
func GetPublications() {
	c := colly.NewCollector(
		colly.AllowedDomains("www.eecs.ucf.edu"),
	)

	publications := []Publication{}

	c.OnHTML("div#publications>div.row:not(.year)", func(e *colly.HTMLElement) {
		selection := e.DOM

		publicationType := strings.TrimSpace(selection.Find("div.slot-6").Text())
		ul := selection.Find("div.slot-7-8-9>ul").Children()
		for i := 0; i < ul.Length(); i++ {
			li := ul.Eq(i)
			publicationText := utils.TrimExcessiveSpaces(li.Text())

			var publication Publication

			var publicationTitle string
			var publicationURL string

			authorsString := strings.Split(publicationText, "\"")[0]
			publicationAuthors := strings.Split(authorsString, "., ")
			publicationAuthors = cleanAuthors(publicationAuthors)

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

	// Create JSON file of publications
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)

		// Encode to JSON
		p, err := json.Marshal(publications)
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err := os.Stat("data"); os.IsNotExist(err) {
			os.Mkdir("data", 0755)
		}

		if err := os.WriteFile("data/publications.json", p, 0644); err != nil {
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
