package collectors

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Person struct {
	Name           string
	ImageUrl       string
	Type           string
	AdditionalInfo string
}

func GetPeople() []Person {
	people := []Person{}

	c := colly.NewCollector(
		colly.AllowedDomains("www.eecs.ucf.edu"),
	)

	c.OnHTML("div.slot-0-1-2>div", func(e *colly.HTMLElement) {
		selection := e.DOM

		personType := strings.TrimSpace(selection.Find("h2").Text())

		selection.Find("ul.people>li").Each(func(i int, s *goquery.Selection) {
			var person Person

			personImage := s.Find("img")
			personInfo := s.Find("div.invisible")

			// The name of the person is in the alt attribute of the image
			// but for some the alt attribute is misppelled
			// so we get the name from the span with the class nametext of the
			// invisible div containing additional information about the person
			if person.Name = personImage.AttrOr("alt", ""); person.Name == "" {
				person.Name = strings.TrimSpace(personInfo.Find("span.nametext").Text())
			}

			person.ImageUrl = personImage.AttrOr("src", "")
			person.Type = personType
			person.AdditionalInfo = personInfo.Find("span:last-of-type").Text()

			people = append(people, person)
		})
	})

	// Create JSON file of people
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)

		// Encode to JSON
		p, err := json.Marshal(people)
		if err != nil {
			fmt.Println(err)
			return
		}

		if _, err := os.Stat("data"); os.IsNotExist(err) {
			os.Mkdir("data", 0755)
		}

		if err := os.WriteFile("data/people.json", p, 0644); err != nil {
			log.Fatal(err)
		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.Visit("https://www.eecs.ucf.edu/isuelab/people/")

	return people
}
