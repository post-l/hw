package circle

import (
	"image"
	"image/color"
	"math/rand"
	"time"

	"github.com/fogleman/gg"
)

type Animation struct {
	ctx     *gg.Context
	circles []*Circle
}

func NewAnimation(sz image.Point) *Animation {
	ctx := gg.NewContext(sz.X, sz.Y)
	circles := make([]*Circle, 15)
	for i := range circles {
		circles[i] = NewCircle(ctx)
	}
	return &Animation{
		ctx:     ctx,
		circles: circles,
	}
}

func (a *Animation) Image() image.Image {
	a.ctx.SetColor(color.Black)
	a.ctx.Clear()
	for _, c := range a.circles {
		c.Draw()
	}
	return a.ctx.Image()
}

func (a *Animation) Delay() time.Duration {
	return time.Second / 10
}

func (a *Animation) Next() error {
	for _, c := range a.circles {
		c.Next()
	}
	return nil
}

type Circle struct {
	ctx      *gg.Context
	position Point
	dir      Point
	color    color.RGBA
	stroke   float64
}

type Point struct {
	X float64
	Y float64
}

func NewCircle(ctx *gg.Context) *Circle {
	c := &Circle{
		ctx:    ctx,
		stroke: float64(rand.Intn(5) + 5),
	}
	c.position.X = float64(rand.Intn(ctx.Width()))
	c.position.Y = float64(rand.Intn(ctx.Height()))
	c.dir.X = initRandDir()
	c.dir.Y = initRandDir()
	c.randColor()
	return c
}

func (c *Circle) Next() {
	c.position.X += c.dir.X
	c.position.Y += c.dir.Y
	update := false
	if c.position.Y+c.stroke >= float64(c.ctx.Height()) {
		update = true
		c.dir.Y = -randDir()
	} else if c.position.Y-c.stroke < 0 {
		update = true
		c.dir.Y = randDir()
	}
	if c.position.X+c.stroke > float64(c.ctx.Width()) {
		update = true
		c.dir.X = -randDir()
	} else if c.position.X-c.stroke < 0 {
		update = true
		c.dir.X = randDir()
	}
	if update {
		c.randColor()
	}
}

func (c *Circle) randColor() {
	c.color.R = uint8(rand.Intn(256))
	c.color.G = uint8(rand.Intn(256))
	c.color.B = uint8(rand.Intn(256))
}

func (c *Circle) Draw() {
	c.ctx.DrawCircle(c.position.X, c.position.Y, c.stroke)
	c.ctx.SetColor(c.color)
	c.ctx.Fill()
}

func initRandDir() float64 {
	dir := rand.Float64()*2 - 1
	if dir < 0 {
		return dir - 1
	}
	return dir + 1
}

func randDir() float64 {
	return (rand.Float64() / 2) + 2
}
