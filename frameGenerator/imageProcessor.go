package frameGenerator

import (
	"image"
)

type ImageProc struct {
	img         img
	framesCount int
	framesImgs  []*image.RGBA
	FramesCh    chan image.Image
}

//Generates frames for one slide from the source image
func (ip *ImageProc) generateFrames() {
	if ip.img.isHorizontal() {
		ip.img.prepareHorizontal()
		randomHorizontalEffect().applyEffect(ip.img, 0.75, ip.framesCount, ip.FramesCh)
	} else {
		ip.img.prepareVertical()
		randomVerticalEffect().applyEffect(ip.img, 0.75, ip.framesCount, ip.FramesCh)
	}
}
