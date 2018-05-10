package main

import (
	"context"
	"flag"
	"image"
	"image/color"
	"time"

	"github.com/fogleman/gg"
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
}

func run(m toolkit.Matrix) {
	tk := toolkit.New(m)
	ctx := context.Background()
	sz := m.Bounds().Size()
	a := NewAnimation(sz)
	tk.PlayAnimation(ctx, a)
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}

type Animation struct {
	ctx      *gg.Context
	position image.Point
	dir      image.Point
	stroke   int
}

func NewAnimation(sz image.Point) *Animation {
	return &Animation{
		ctx:    gg.NewContext(sz.X, sz.Y),
		dir:    image.Point{1, 1},
		stroke: 5,
	}
}

func (a *Animation) Image() image.Image {
	a.ctx.SetColor(color.Black)
	a.ctx.Clear()
	a.ctx.DrawCircle(float64(a.position.X), float64(a.position.Y), float64(a.stroke))
	a.ctx.SetColor(color.RGBA{255, 0, 0, 255})
	a.ctx.Fill()
	return a.ctx.Image()
}

func (a *Animation) Delay() time.Duration {
	return time.Second / 30
}

func (a *Animation) Next() error {
	a.position.X += 1 * a.dir.X
	a.position.Y += 1 * a.dir.Y
	if a.position.Y+a.stroke >= a.ctx.Height() {
		a.dir.Y = -1
	} else if a.position.Y-a.stroke < 0 {
		a.dir.Y = 1
	}
	if a.position.X+a.stroke > a.ctx.Width() {
		a.dir.X = -1
	} else if a.position.X-a.stroke < 0 {
		a.dir.X = 1
	}
	return nil
}
