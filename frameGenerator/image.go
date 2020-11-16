package frameGenerator

import (
	"image"
	"log"
	"strings"
	config "youtube-slideshow/configuration"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type img struct {
	imgPath  string
	imgName  string
	imgImage image.Image
	imgOrig  image.Image
	rgbaCopy *image.RGBA
}

func (i *img) open() {
	img, err := imgio.Open(i.imgPath)
	if err != nil {
		log.Panic(err)
	}
	i.imgOrig = img
}

func (i img) isHorizontal() bool {
	isHorizontal := false

	if i.imgOrig.Bounds().Dx() > i.imgOrig.Bounds().Dy() {
		isHorizontal = true
	}

	return isHorizontal
}

func (i img) nameWithoutFormat() string {
	name := ""

	name = i.imgName[:strings.Index(i.imgName, ".")]

	return name
}

func (i *img) prepareHorizontal() {
	i.imgResize(config.Width, config.Height)
}

func (i *img) prepareVertical() {
	w := float64(i.imgOrig.Bounds().Dx())
	h := float64(i.imgOrig.Bounds().Dy())

	k := float64(config.Width) / w

	newWidth := int(w * k)
	newHeight := int(h * k)

	i.imgResize(newWidth, newHeight)
}

func (i *img) stretch() {
	newWidth := int(float64(i.imgImage.Bounds().Dx()) * (1.00 + config.Percent))
	newHeight := int(float64(i.imgImage.Bounds().Dy()) * (1.00 + config.Percent))

	i.imgResize(newWidth, newHeight)
}

func (i *img) imgResize(width, height int) {
	resized := transform.Resize(i.imgOrig, width, height, transform.Linear)
	i.imgImage = resized
}
