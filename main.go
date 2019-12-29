package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width       = 2048 // actual 4096 but doesn't fit my laptop's screen
	height      = 768
	numConns    = 1
	serverWand  = "151.217.111.34:1234"
	serverBühne = "151.217.176.193:1234" // 4096x768
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

	//pixelChan := make(chan *myPixel, 100)
	for i := 0; i < numConns; i++ {
		go func() {
			for {
				conn, err := net.Dial("tcp", serverBühne)
				if err != nil {
					log.Printf("dial: %v", err)
					time.Sleep(200 * time.Millisecond) // give me a break :)
					continue
				}
				go func() {
					defer conn.Close()
					scanner := bufio.NewScanner(conn)
					for scanner.Scan() {
						line := scanner.Text()
						px := myPixel{}
						_, err = fmt.Sscanf(line,
							"PX %d %d %02x%02x%02x%02x",
							&px.x, &px.y,
							&px.r, &px.g, &px.b, &px.a)
						if err != nil {
							log.Printf("unable to parse `%s`: %v", line, err)
						}
						img.Set(px.x, px.y, color.RGBA{px.r, px.g, px.b, px.a})
					}
					if err := scanner.Err(); err != nil {
						log.Printf("scan error: %v", err)
						return
					}
				}()
				for {
					_, _ = conn.Write([]byte(fmt.Sprintf("PX %d %d\n", rand.Int31n(width), rand.Int31n(height))))
					if err != nil {
						log.Printf("write: %v", err)
						conn.Close()
						break
					}

				}
			}
		}()
	}

	for !win.Closed() {
		time.Sleep(100 * time.Millisecond) // poor man's frame rate XoXo
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
