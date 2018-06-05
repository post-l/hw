package main

import (
	"flag"
	"image/color"
	"time"

	"github.com/post-l/hw/board/tinkerboard"
	"github.com/post-l/hw/matrix"
	"github.com/post-l/hw/matrix/emulator"
	"github.com/post-l/hw/matrix/toolkit"
)

var (
	emFlag = flag.Bool("emulator", false, "use emulator")
)

func main() {
	flag.Parse()

	if *emFlag {
		m := emulator.NewEmulator(&matrix.DefaultHardwareConfig)
		go func() {
			<-m.Ready()
			run(m)
		}()
		m.Init()
	}
	b, err := tinkerboard.New()
	fatal(err)
	m := matrix.New(b, &matrix.DefaultHardwareConfig)
	defer m.Close()
	run(m)
	time.Sleep(10 * time.Second)
}

func run(m toolkit.Matrix) {
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
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
