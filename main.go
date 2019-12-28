package main

import (
	"image"

	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width  = 1920
	height = 1080
)

func run() {
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Pixel Grab",
		Bounds: pixel.R(0, 0, width, height),
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}

	// image data structure to set pixels on
	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{width, height},
	})

	for !win.Closed() {

		// code from michi
		// for each pixel
		img.Set(10, 30, color.White) // example

		// update window with current pixels
		pic := pixel.PictureDataFromImage(img)
		sprite := pixel.NewSprite(pic, pic.Bounds())
		sprite.Draw(win, pixel.IM.Moved(win.Bounds().Center()))
		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
