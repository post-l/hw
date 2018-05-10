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
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			m.Set(x, y, color.RGBA{255, 0, 0, 255})
			m.Render()
		}
	}
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
