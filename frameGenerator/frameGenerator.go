package frameGenerator

import (
	"image"
	"io/ioutil"
	"math/rand"
	"strconv"
	"time"
	config "youtube-slideshow/configuration"

	"github.com/anthonynsimon/bild/imgio"
)

type (
	FrameGenerator struct {
		framesCount int
		FramesCh    chan image.Image
		StopCh      chan bool
	}
)

func (fg *FrameGenerator) InitConfig() {
	fg.framesCount = config.Fps * config.SlideDur

}

func (fg *FrameGenerator) GenerateFrames(photosPath string) error {
	var err error
	png, err = chooseRandomPng()
	if err != nil {
		return err
	}

	imgProc := ImageProc{
		img:         img{},
		framesCount: fg.framesCount,
		FramesCh:    fg.FramesCh,
	}

	files, err := ioutil.ReadDir(photosPath)
	if err != nil {
		return err
	}

	for _, f := range files {
		imgProc.img = img{
			imgPath: photosPath + "/" + f.Name(),
			imgName: f.Name(),
		}
		imgProc.img.open()

		imgProc.generateFrames()
	}

	fg.StopCh <- true

	return nil
}

func chooseRandomPng() (image.Image, error) {
	//Read all png files
	files, err := ioutil.ReadDir(config.PngDir)
	if err != nil {
		return nil, err
	}
	//Generate random picture number
	rand.Seed(time.Now().Unix())
	picNum := strconv.Itoa(rand.Intn(len(files)) + 1)

	path := config.PngDir + picNum + ".png"

	png, err := imgio.Open(path)
	if err != nil {
		return nil, err
	}

	return png, nil
}
