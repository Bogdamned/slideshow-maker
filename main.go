package main

import (
	"fmt"
	"image"
	"strings"
	config "youtube-slideshow/configuration"
	fg "youtube-slideshow/frameGenerator"
	movmk "youtube-slideshow/movieMaker"
	prsr "youtube-slideshow/parser"
)

func main() {

	city := "Италия" //os.Args[1]
	cfg := ""        //os.Args[2]

	config.InitConfiguration(cfg)
	go generateMoviesForAllHotelsByCountry(city)

	<-config.EndCh

}

func createParser() PhotoParser {
	var parser PhotoParser
	hotelPrsr := new(prsr.HotelParser)
	hotelPrsr.InitConfig()
	parser = hotelPrsr

	return parser
}

func createFrameGenerator() FrameGenerator {
	var frameGen FrameGenerator
	frGen := new(fg.FrameGenerator)
	frGen.InitConfig()
	frameGen = frGen

	return frameGen
}

func createMovieMaker() MovieMaker {
	var mmk MovieMaker
	mm := new(movmk.MovieMaker)
	mmk = mm

	return mmk
}

func generateMoviesForAllHotelsByCountry(country string) {
	parser := createParser()
	go parser.ParseHotelsByCountry(country)

	go func() {
		var photosPath string
		for {
			select {
			case photosPath = <-config.ParsedHotels:
				framesCh := make(chan image.Image)
				stopCh := make(chan bool)
				movieMaker := new(movmk.MovieMaker)
				movieMaker.FramesCh = framesCh
				movieMaker.Stop = stopCh
				go movieMaker.MakeMovie(photosPath[strings.LastIndex(photosPath, "/")+1:])

				frGen := new(fg.FrameGenerator)
				frGen.InitConfig()
				frGen.FramesCh = framesCh
				frGen.StopCh = stopCh

				frGen.GenerateFrames(photosPath)

			case <-config.Quit:
				config.Parsed = true
				fmt.Println("Everything is parsed")
				return
			}
		}
	}()
}
