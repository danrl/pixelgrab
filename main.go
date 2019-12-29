package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"
	"net"
	"sync/atomic"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	width               = 4096
	height              = 768
	scale       float64 = 0.5 // Use 0.25 or 0.5
	numConns            = 4
	serverWand          = "151.217.111.34:1234"
	serverBühne         = "151.217.176.193:1234" // 4096x768
)

var showEveryPixel = int(1.0 / scale)
var currentPixel uint64

type myPixel struct {
	x, y       int
	r, g, b, a uint8
}

func (p myPixel) getX() int {
	return p.x / showEveryPixel
}
func (p myPixel) getY() int {
	return p.y / showEveryPixel
}

func run() {
	// Populate work array which will act as a queue
	var queue []myPixel
	for x := 0; x < width; x += showEveryPixel {
		for y := 0; y < height; y += showEveryPixel {
			queue = append(queue, myPixel{x: x, y: y})
		}
	}

	// Randomize so that the image builds up quicker
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(queue), func(i, j int) { queue[i], queue[j] = queue[j], queue[i] })
	queueLen := uint64(len(queue)) - 2*numConns

	// Create our canvas
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  "Pixel Grab",
		Bounds: pixel.R(0, 0, width*scale, height*scale),
		VSync:  true,
	})
	if err != nil {
		panic(err)
	}

	// image data structure to set pixels on
	img := image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{int(width * scale), int(height * scale)},
	})

	// Use multiple connections for faster buildup
	for i := 0; i < numConns; i++ {
		go func() { // Do connection in a separate goroutine
			for { // Reconnect handling
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

						// When setting the pixel, we need to scale accordingly.
						img.Set(px.getX(), px.getY(), color.RGBA{px.r, px.g, px.b, px.a})
					}
					if err := scanner.Err(); err != nil {
						log.Printf("scan error: %v", err)
						return
					}
				}()
				for {
					readNext := atomic.AddUint64(&currentPixel, 1) - 1 // Increase first and get the previous pixel position
					currentPixel = currentPixel % queueLen             // Synchronization yolo
					cmd := fmt.Sprintf("PX %d %d\n", queue[readNext].x, queue[readNext].y)
					_, _ = conn.Write([]byte(cmd))
					if err != nil {
						log.Printf("write: %v", err)
						conn.Close()
						break
					}

				}
			}
		}()
	}

	// Render image periodically onto our canvas
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
