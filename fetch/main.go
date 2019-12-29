package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net"
	"os"
	"time"
)

const numConns = 4
const width = 1920
const height = 1080
const serverWand = "151.217.111.34:1234"
const serverBühne = "151.217.176.193:1234"

type pixel struct {
	x, y int
	color uint
}

func main() {
	var conns []net.Conn

	for i := 0; i < numConns; i++ {
		conn, err := net.Dial("tcp", serverBühne)
		if err != nil {
			log.Fatalf("dial: %v", err)
		}
		log.Printf("connection %d established\n", i)
		conns = append(conns, conn)
	}

	start := time.Now()

	go func() {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				fmt.Fprintf(conns[y%numConns], "PX %d %d\n", x, y)
			}
		}
		log.Println("requested all pixels")
	}()

	pixelChan := make(chan *pixel)

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

				pixelChan <- &pixel{x, y, c}
			}
		}(conns[i])
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < width*height; i++ {
		px := <-pixelChan
		if i % 1000 == 0 {
			log.Printf("Progress %02.f%%\n", 100*float64(i)/float64(width*height))
		}
		r := uint8(px.color & (0xFF<<3*8))
		g := uint8(px.color & (0xFF<<2*8))
		b := uint8(px.color & (0xFF<<1*8))
		a := uint8(px.color & (0xFF<<0*8))

		c := color.RGBA{r, g, b, a}
		img.Set(px.x, px.y, c)
	}

	if err := png.Encode(os.Stdout, img); err != nil {
		log.Fatalf("png: %v\n", err)
	}

	log.Printf("done after %v\n", time.Since(start))
}
