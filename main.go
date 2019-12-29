package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"net"
	"time"

	"image/color"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width       = 1920 / 10
	height      = 1080 / 10
	numConns    = 4
	serverWand  = "151.217.111.34:1234"
	serverBühne = "151.217.176.193:1234"
)

type myPixel struct {
	x, y       int
	r, g, b, a uint8
}

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

	pixelChan := make(chan *myPixel)
	for i := 0; i < numConns; i++ {
		go func() {
			for {
				conn, err := net.Dial("tcp", serverBühne)
				if err != nil {
					log.Printf("dial: %v", err)
					time.Sleep(50 * time.Millisecond)
					continue
				}
				line := make([]byte, 2048)
				for {
					foo := ""
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))
					foo = foo + fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))

					_, _ = conn.Write([]byte(foo))
					if err != nil {
						log.Printf("write: %v", err)
						conn.Close()
						break
					}
					_, err = conn.Read(line)
					if err != nil {
						log.Printf("unable to read: %v", err)
					}

					px := myPixel{}
					_, err = fmt.Sscanf(string(line),
						"PX %d %d %02x%02x%02x%02x",
						&px.x, &px.y,
						&px.r, &px.g, &px.b, &px.a)
					if err != nil {
						log.Printf("unable to parse `%s`: %v", line, err)
					}
					pixelChan <- &px
				}
			}
		}()
	}

	for !win.Closed() {
		for i := 0; i < 10; i++ {
			px := <-pixelChan
			img.Set(px.x, px.y, color.RGBA{px.r, px.g, px.b, px.a})
		}

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
