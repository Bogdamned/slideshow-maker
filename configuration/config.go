package config

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
)

//Struct to read from json configuration file
type cfg struct {
	Width           int     `json: width`
	Height          int     `json: height`
	Fps             int     `json: fps`
	SlideDur        int     `json: slideDur`
	Percent         float64 `json: percent`
	MoviesDirectory string  `json: movies_directory`
	Limit           int     `json: limit` //set -1 to parse unlimited
}

var Config = cfg{}
var FramesCh chan image.Image
var EndCh chan int
var Quit = make(chan int)
var ParsedHotels = make(chan string)

//var MovieGenerated = make(chan bool, 1)
var Parsed = false

//Directories config
const Root = "./files/"

var BookingDir = Root + "booking/"
var MoviesDir = Root + "movies/"
var PngDir = Root + "png/"
var MusicDir = Root + "music/"

//Parser config
const BookingUrl = "https://www.booking.com"
const CityUrlFmt = "https://www.booking.com/searchresults.ru.html?ss=%s&nflt=ht_id%%3D204%%3B&shw_aparth=0"

var Limit int = 10

//Replace below with config
//Image configuration
var Width int
var Height int

//Movie configuration
var Fps int
var SlideDur int

//Horizontal effet amount
var Percent float64

func InitConfiguration(configPath string) {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Println("Unable to open", configPath, " ", err)
		log.Println("Reading default configuration")

		jsonFile, err = os.Open("./configuration/config.json")
		if err != nil {
			log.Println("Default configuration file config.json is missing. Going to use hardcoded settings")
		} else {
			log.Println("Successfully opened default config.json")
		}
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &Config)

	log.Println("Initializing configuration:")
	if Config.Width != 0 {
		Width = Config.Width
	} else {
		Width = 1280
	}
	log.Println("Movie width: ", Width)

	if Config.Height != 0 {
		Height = Config.Height
	} else {
		Height = 720
	}
	log.Println("Movie height: ", Height)

	if Config.Fps != 0 {
		Fps = Config.Fps
	} else {
		Fps = 35
	}
	log.Println("Movie fps: ", Fps)

	if Config.SlideDur != 0 {
		SlideDur = Config.SlideDur
	} else {
		SlideDur = 3
	}
	log.Println("Slide duration: ", SlideDur)

	if Config.Percent != 0.0 {
		Percent = Config.Percent
	} else {
		Percent = 0.1
	}
	log.Println("Zooming effects amount: ", Percent)

	if Config.MoviesDirectory != "" {
		MoviesDir = Config.MoviesDirectory
	}
	log.Println("Directory where movies will be saved: ", MoviesDir)

	if Config.Limit != 0 {
		Limit = Config.Limit
	}
	log.Println("Numbers of hotels to parse: ", Limit)

	fmt.Println("\n")

	FramesCh = make(chan image.Image, Config.Fps*Config.SlideDur)

}
