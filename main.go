package main

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"log"
	"net"

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
	x, y  int
	color uint
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

	var conns []net.Conn

	for i := 0; i < numConns; i++ {
		conn, err := net.Dial("tcp", serverBühne)
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		log.Printf("connection %d established\n", i)
		conns = append(conns, conn)
	}

	go func() {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				fmt.Fprintf(conns[y%numConns], "PX %d %d\n", x, y)
			}
		}
		log.Println("requested all pixels")
	}()

	pixelChan := make(chan *myPixel)

	for i := 0; i < numConns; i++ {
		go func(conn net.Conn) {
			reader := bufio.NewReader(conn)
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					log.Fatalf("read: %v", err)
				}

				var x, y int
				var c uint
				_, err = fmt.Sscanf(line, "PX %d %d %06x", &x, &y, &c)
				if err != nil {
					log.Printf("unable to parse `%s`: %v", line, err)
				}

				pixelChan <- &myPixel{x, y, c}
			}
		}(conns[i])
	}

	for !win.Closed() {
		for i := 0; i < 10; i++ {
			px := <-pixelChan
			r := uint8(px.color & (0xFF << 3 * 8))
			g := uint8(px.color & (0xFF << 2 * 8))
			b := uint8(px.color & (0xFF << 1 * 8))
			a := uint8(px.color & (0xFF << 0 * 8))
			img.Set(px.x, px.y, color.RGBA{r, g, b, a})
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
