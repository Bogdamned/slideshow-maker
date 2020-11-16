package parser

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	config "youtube-slideshow/configuration"

	"github.com/PuerkitoBio/goquery"
)

type (
	parserConfig struct {
		bookingUrl string
		cityUrlFmt string
		limit      int
	}
	HotelParser struct {
		parserConfig
		hotelPage  hotelPage
		searchPage searchPage
	}

	page struct {
		url string
	}
)

func (p *HotelParser) InitConfig() {
	p.bookingUrl = config.BookingUrl
	p.cityUrlFmt = config.CityUrlFmt
	p.limit = config.Limit
}

//Scraps all hotels by city or country
func (p *HotelParser) ParseHotelsByCountry(city string) {
	offsetStep := 0
	urlMask := fmt.Sprintf(p.cityUrlFmt, city) + "&offset="
	p.searchPage.url = urlMask + strconv.Itoa(offsetStep)

	fmt.Println(p.searchPage.url)

	for p.searchPage.pageHasHotels() {
		p.searchPage.url = urlMask + strconv.Itoa(offsetStep)

		stop := p.parseHotelsByPage()
		// If limit was reached stop the parser
		if stop {
			log.Println("Reached limit of hotels to parse.")
			config.Quit <- 1
			break
		}
		//Step to the next page
		offsetStep += 15
	}

	config.Quit <- 1
}

//Downloads all hotel photos from a searchpage with hotels list
func (p *HotelParser) parseHotelsByPage() bool {
	limitReached := false
	response, err := http.Get(p.searchPage.url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Goquery: error loading HTTP response body. ", err)
	}

	document.Find("a").EachWithBreak(func(index int, element *goquery.Selection) bool {

		if element.HasClass("hotel_name_link url") {
			hotelHref, exists := element.Attr("href")
			if exists {
				p.hotelPage.url = p.bookingUrl + formatHotelHref(hotelHref)
				parsed := p.hotelPage.parseHotel()

				if parsed {
					p.limit--
				}

			}
		}

		limitReached = (p.limit == 0)
		loopAgain := !limitReached

		return loopAgain
	})

	return limitReached
}

//Formats proper hotel link from scrapped one
func formatHotelHref(href string) string {
	strLim := strings.Index(href, "?")

	fmtHref := href[1:strLim]

	return fmtHref
}
