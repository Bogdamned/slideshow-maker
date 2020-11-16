package fileManager

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	config "youtube-slideshow/configuration"
)

func init() {
	//Check if root directory exists if not then create it
	if !DirExists(config.Root) {
		os.Mkdir(config.Root, 0700)
	}
}

func DownloadPhotos(saveDir string, imgUrls []string) {
	//Check if dir already exists
	if !DirExists(saveDir) {
		//create directory
		err := os.MkdirAll(saveDir, 0700)

		if err == nil {
			log.Println("Created directory: ", saveDir)
		} else {
			log.Fatal("Error occured while creating the directory: ", saveDir, " Err: ", err.Error())
		}

		//Loop over photos
		for i, url := range imgUrls {
			DownloadJPG(url, saveDir, strconv.Itoa(i+1))
		}
	} else {
		log.Println("Directory already exists, stopped downloading proccess: ", saveDir)
	}
}

func DownloadJPG(imgURL, imgPath, imgName string) {
	imgName += ".jpg"

	response, err := http.Get(imgURL)
	if err != nil {
		log.Println("An error occured requesting image url: ", imgURL, ". Error: ", err)
	}
	defer response.Body.Close()

	//Upload photo
	outFile, err := os.Create(imgPath + imgName)
	if err != nil {
		log.Println("An error occured on creating a file ", imgPath+imgName, " Error: ", err)
		return
	}
	defer outFile.Close()

	// Copy data from HTTP response to file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		log.Println("An error occured while saving ", imgName, "to file. Error: ", err)
		return
	} else {
		log.Println("Successfully downloaded img: ", imgName)
	}
}

func DirExists(path string) bool {
	exists := false

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		exists = true
	}

	return exists
}
