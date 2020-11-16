package frameGenerator

import (
	"image"
	"math/rand"
	"time"
	config "youtube-slideshow/configuration"

	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/transform"
)

type effector interface {
	applyEffect(i img, percent float64, frames int, fCh chan image.Image)
}

//Horizontal effects
type leftTopFadeOut bool
type rightBotFadeOut bool
type rightTopFadeOut bool
type leftBotFadeOut bool
type zoomIn bool
type zoomOut bool

//Vertical effects
type topBotSlide bool
type botTopSlide bool

var png image.Image

// Sends frame to movie maker routine through channnel
func sendFrame(img image.Image, rect image.Rectangle, fCh chan image.Image) {
	fCh <- applyPng(transform.Crop(img, rect))
}

func applyPng(img image.Image) image.Image {
	if png.Bounds().Dx() != config.Width || png.Bounds().Dy() != config.Height {
		png = transform.Resize(png, config.Width, config.Height, transform.NearestNeighbor)
	}
	return blend.Normal(img, png)
}

func (zi zoomIn) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	iWidth := float64(i.imgImage.Bounds().Dx())
	iHeight := float64(i.imgImage.Bounds().Dy())

	stepWidth := (float64(config.Width) / float64(frames)) * config.Percent
	stepHeight := (float64(config.Height) / float64(frames)) * config.Percent

	for it := 0; it < frames; it++ {
		iWidth += stepWidth
		iHeight += stepHeight

		i.imgResize(int(iWidth), int(iHeight))

		rect := image.Rect(int((iWidth-float64(config.Width))/2), int((iHeight-float64(config.Height))/2), int((iWidth+float64(config.Width))/2), int((iHeight+float64(config.Height))/2))

		sendFrame(i.imgImage, rect, fCh)
	}
}

func (zo zoomOut) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	i.stretch()

	iWidth := float64(i.imgImage.Bounds().Dx())
	iHeight := float64(i.imgImage.Bounds().Dy())

	stepWidth := (float64(config.Width) / float64(frames)) * config.Percent
	stepHeight := (float64(config.Height) / float64(frames)) * config.Percent

	for it := 0; it < frames; it++ {
		iWidth -= stepWidth
		iHeight -= stepHeight

		i.imgResize(int(iWidth), int(iHeight))

		rect := image.Rect(int((iWidth-float64(config.Width))/2), int((iHeight-float64(config.Height))/2), int((iWidth+float64(config.Width))/2), int((iHeight+float64(config.Height))/2))

		sendFrame(i.imgImage, rect, fCh)
	}
}

func (rfo rightBotFadeOut) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	i.stretch()

	iWidth := float64(i.imgImage.Bounds().Dx())
	iHeight := float64(i.imgImage.Bounds().Dy())

	stepWidth := (float64(config.Width) / float64(frames)) * config.Percent
	stepHeight := (float64(config.Height) / float64(frames)) * config.Percent

	for it := 0; it < frames; it++ {
		iWidth -= stepWidth
		iHeight -= stepHeight

		i.imgResize(int(iWidth), int(iHeight))
		rect := image.Rect(int(iWidth-float64(config.Width)), int(iHeight-float64(config.Height)), int(iWidth), int(iHeight))

		sendFrame(i.imgImage, rect, fCh)
	}
}

func (rfo rightTopFadeOut) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	i.stretch()

	iWidth := float64(i.imgImage.Bounds().Dx())
	iHeight := float64(i.imgImage.Bounds().Dy())

	stepWidth := (float64(config.Width) / float64(frames)) * config.Percent
	stepHeight := (float64(config.Height) / float64(frames)) * config.Percent

	for it := 0; it < frames; it++ {
		iWidth -= stepWidth
		iHeight -= stepHeight

		i.imgResize(int(iWidth), int(iHeight))
		rect := image.Rect(int(iWidth)-config.Width, 0, int(iWidth), config.Height)

		sendFrame(i.imgImage, rect, fCh)
	}
}

func (lfo leftTopFadeOut) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	i.stretch()
	rect := image.Rect(0, 0, config.Width, config.Height)

	iWidth := float64(i.imgImage.Bounds().Dx())
	iHeight := float64(i.imgImage.Bounds().Dy())

	stepWidth := (float64(config.Width) / float64(frames)) * config.Percent
	stepHeight := (float64(config.Height) / float64(frames)) * config.Percent

	for it := 0; it < frames; it++ {
		iWidth -= stepWidth
		iHeight -= stepHeight

		i.imgResize(int(iWidth), int(iHeight))

		sendFrame(i.imgImage, rect, fCh)
	}
}

func (lfo leftBotFadeOut) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	i.stretch()

	iWidth := float64(i.imgImage.Bounds().Dx())
	iHeight := float64(i.imgImage.Bounds().Dy())

	stepWidth := (float64(config.Width) / float64(frames)) * config.Percent
	stepHeight := (float64(config.Height) / float64(frames)) * config.Percent

	for it := 0; it < frames; it++ {
		iWidth -= stepWidth
		iHeight -= stepHeight

		i.imgResize(int(iWidth), int(iHeight))
		rect := image.Rect(0, int(iHeight)-config.Height, config.Width, int(iHeight))

		sendFrame(i.imgImage, rect, fCh)
	}
}

func (ud topBotSlide) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	var rect = image.Rectangle{}
	imgHeight := float64(i.imgImage.Bounds().Dy())

	//set cur height to conastant value of height
	curHeight := float64(config.Height)
	step := 1.0

	for it := 0; it < frames; it++ {
		curHeight += step
		if curHeight < imgHeight {
			rect = image.Rect(0, int(curHeight)-config.Height, config.Width, int(curHeight))
		} else {
			rect = image.Rect(0, int(imgHeight)-config.Height, config.Width, int(imgHeight))
		}

		sendFrame(i.imgImage, rect, fCh)
	}
}

func (ud botTopSlide) applyEffect(i img, percent float64, frames int, fCh chan image.Image) {
	var rect = image.Rectangle{}

	//set cur height to conastant value of height
	curHeight := float64(i.imgImage.Bounds().Dy())
	step := 1.0

	for it := 0; it < frames; it++ {
		curHeight -= step
		if int(curHeight) > config.Height {
			rect = image.Rect(0, int(curHeight)-config.Height, config.Width, int(curHeight))
		} else {
			rect = image.Rect(0, 0, config.Width, config.Height)
		}

		sendFrame(i.imgImage, rect, fCh)
	}
}

func randomHorizontalEffect() effector {
	effects := []effector{rightTopFadeOut(false), leftBotFadeOut(false), rightBotFadeOut(false)} //zoomIn(false)} , zoomOut(false)
	rand.Seed(time.Now().Unix())

	return effects[rand.Intn(len(effects))]
}

func randomVerticalEffect() effector {
	effects := []effector{topBotSlide(false), botTopSlide(false)}
	rand.Seed(time.Now().Unix())

	return effects[rand.Intn(len(effects))]
}
