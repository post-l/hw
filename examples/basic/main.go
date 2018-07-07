package main

import (
	"image/color"
	"time"

	"github.com/post-l/hw/examples"
	"github.com/post-l/hw/matrix/toolkit"
)

func main() {
	examples.Main(run)
}

func run(m toolkit.Matrix) error {
	bounds := m.Bounds()
	c := color.RGBA{0, 0, 255, 255}
	thirdX := (bounds.Min.X + bounds.Max.X) / 3
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		t := time.Now()
		if x == thirdX-1 {
			c = color.RGBA{255, 255, 255, 255}
		} else if x == thirdX*2 {
			c = color.RGBA{255, 0, 0, 255}
		}
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			m.Set(x, y, c)
		}
		m.Render()
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			m.Set(x, y, c)
		}
		m.Render()
		if d := 150*time.Millisecond - time.Since(t); d > 0 {
			time.Sleep(d)
		}
	}
	time.Sleep(10 * time.Second)
	return nil
}
