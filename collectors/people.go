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

type People struct {
	Director              Person   `json:"director"`
	AssociateDirector     Person   `json:"associateDirector"`
	AffiliatedFaculty     []Person `json:"affiliatedFaculty"`
	ResearchStaff         []Person `json:"researchStaff"`
	PhdStudents           []Person `json:"phdStudents"`
	MastersStudents       []Person `json:"mastersStudents"`
	UndergraduateStudents []Person `json:"undergraduateStudents"`
	Alumni                []Person `json:"alumni"`
	FormerVisitors        []Person `json:"formerVisitors"`
}

type Person struct {
	Name           string `json:"name"`
	PageURL        string `json:"pageURL"`
	ImageURL       string `json:"imageURL"`
	Type           string `json:"type"`
	AdditionalInfo string `json:"additionalInfo"`
}

func GetPeople() People {
	people := People{}

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

			personPageURL := s.Find("a").AttrOr("href", "")

			// If the url is prefixed with "http" it is a link to another page
			// otherwise it is a route
			if !strings.HasPrefix(personPageURL, "http") {
				personPageURL = strings.Split(personPageURL, ".")[0]
			}
			
			// The name of the person is in the alt attribute of the image
			// but for some, the alt attribute is misppelled
			// so we get the name from the span with the class nametext of the
			// invisible div containing additional information about the person
			if person.Name = personImage.AttrOr("alt", ""); person.Name == "" {
				person.Name = strings.TrimSpace(personInfo.Find("span.nametext").Text())
			}

			person.ImageURL = personImage.AttrOr("src", "")
			person.Type = personType
			person.AdditionalInfo = personInfo.Find("span:last-of-type").Text()
			person.PageURL = strings.ToLower(personPageURL)

			switch personType {
			case "Director":
				people.Director = person
			case "Associate Director":
				people.AssociateDirector = person
			case "Affiliated Faculty":
				people.AffiliatedFaculty = append(people.AffiliatedFaculty, person)
			case "Research Staff":
				people.ResearchStaff = append(people.ResearchStaff, person)
			case "Ph.D. Students":
				people.PhdStudents = append(people.PhdStudents, person)
			case "Masters Students":
				people.MastersStudents = append(people.MastersStudents, person)
			case "Undergraduate Students":
				people.UndergraduateStudents = append(people.UndergraduateStudents, person)
			case "Alumni":
				people.Alumni = append(people.Alumni, person)
			case "Former Visitors":
				people.FormerVisitors = append(people.FormerVisitors, person)
			}
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
