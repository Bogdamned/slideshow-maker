package parser

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	config "youtube-slideshow/configuration"
	"youtube-slideshow/fileManager"

	"github.com/PuerkitoBio/goquery"
)

type hotelPage struct {
	page
	document         *goquery.Document
	imgUrls          []string
	hotelId          string
	hotelName        string
	hotelDescription string
	country          string
	city             string
}

//Scraps all photos and info related to a hotel
func (p *hotelPage) parseHotel() bool {
	fmt.Print(fmt.Sprintf("Parsing %v \n", p.url))

	p.parseHotelId()

	downloadDest := config.BookingDir + p.hotelId + "/"
	if p.isParsed() {
		log.Println(fmt.Sprintf("Hotel %v was already parsed.\n", p.hotelName))
		return false
	}

	p.createGoQueryDocument()
	p.parseImgUrls()

	if len(p.imgUrls) < 5 {
		log.Println(fmt.Sprintf("Hotel %v has less than 5 photos, so it shouldnt be parsed.\n", p.hotelName))
		return false
	}

	p.parseName()
	p.parseInfo()

	fileManager.DownloadPhotos(downloadDest, p.imgUrls)
	fmt.Println("Parse and sending to chan: " + config.BookingDir + p.hotelId)

	p.saveToTxt()

	config.ParsedHotels <- config.BookingDir + p.hotelId

	return true
}

// Extracts hotel id from url
func (p *hotelPage) parseHotelId() {
	parsedUrl, err := url.Parse(p.url)
	if err != nil {
		log.Fatal("Unable to parse hotel's URL:", err.Error())
	}

	slashPos := strings.LastIndex(parsedUrl.Path, "/") + 1
	dotPos := strings.Index(parsedUrl.Path, ".")

	p.hotelId = parsedUrl.Path[slashPos:dotPos]
}

// Creates goquery document for a hotel page
func (p *hotelPage) createGoQueryDocument() {
	response, err := http.Get(p.url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Goquery: error loading HTTP response body. ", err)
	}

	p.document = document
}

// Gets links to all hotel photos.
func (p *hotelPage) parseImgUrls() {
	p.imgUrls = []string{}

	response, err := http.Get(p.url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Goquery: error loading HTTP response body. ", err)
	}

	document.Find("img").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("data-highres")
		if exists {
			p.imgUrls = append(p.imgUrls, imgSrc)
		}
	})
}

//Gets hotel name from hotel url
func (p *hotelPage) parseName() {

	p.document.Find(".hp__hotel-name").Each(func(index int, element *goquery.Selection) {
		hotel := strings.TrimSpace(element.Text())
		p.hotelName = hotel
	})
}

// Gets hotel info: city, country, description
func (p *hotelPage) parseInfo() {

	// Parse city
	p.city, _ = p.document.Find("#ss").First().Attr("value")
	descr := ""

	// Parse country
	p.document.Find("#breadcrumb > ol > li:nth-child(2) > div > a").Each(func(index int, element *goquery.Selection) {
		p.country = strings.TrimSpace(element.Text())
	})

	// Parse description
	p.document.Find("#summary").Children().Each(func(index int, element *goquery.Selection) {
		if element.HasClass("chain-content chain-content-break-line") {
			return
		} else {
			descr = descr + strings.TrimSpace(element.Text())
		}

	})
	p.hotelDescription = descr

}

//Checks if hotel was already parsed
func (p *hotelPage) isParsed() bool {
	isParsed := false
	hotelDir := config.BookingDir + p.hotelId

	if fileManager.DirExists(hotelDir) {
		isParsed = true
	}

	return isParsed
}

// Saves hotel information to txt file
func (p *hotelPage) saveToTxt() error {
	err := os.MkdirAll(config.MoviesDir+p.hotelId, 0700)
	if err == nil {
		log.Println("Created directory: ", config.MoviesDir+p.hotelId)
	} else {
		log.Fatal("Error occured while creating the directory: ", config.MoviesDir+p.hotelId, " Err: ", err)
		return err
	}

	f, err := os.Create(config.MoviesDir + p.hotelId + "/" + p.hotelId + ".txt")
	defer f.Close()

	_, err = f.WriteString(p.hotelName + "\n" + "Country: " + p.country + "\n" + "City: " + p.city + "\n" + "Description: " + p.hotelDescription + "\n")
	if err != nil {
		return err
	}

	return nil
}
