package parser

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type searchPage struct {
	page
}

//Checks if there are any hotel links on the page
func (p *searchPage) pageHasHotels() bool {
	hasHotels := false

	response, err := http.Get(p.url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Goquery: error loading HTTP response body. ", err)
	}

	document.Find("a").EachWithBreak(func(index int, element *goquery.Selection) bool {
		if element.HasClass("hotel_name_link url") {
			_, exists := element.Attr("href")
			if exists {
				hasHotels = true
			}
		}

		loopAgain := !hasHotels
		return loopAgain
	})

	return hasHotels
}
