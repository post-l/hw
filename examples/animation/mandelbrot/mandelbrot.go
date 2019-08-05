// Package mandelbrot image renderer.
// Inspired from github.com/esimov/gobrot
package mandelbrot

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"time"
)

type Animation struct {
	image              *image.RGBA
	colors             []color.RGBA
	maxIteration       int
	escapeRadius, x, y float64
}

func NewAnimation(sz image.Point) *Animation {
	a := &Animation{
		image:        image.NewRGBA(image.Rect(0, 0, sz.X, sz.Y)),
		colors:       interpolateColors(4000),
		maxIteration: 800,
		escapeRadius: .02401245,
		x:            -0.0091275,
		y:            0.7899912,
	}
	fmt.Println(a.colors)
	a.render()
	return a
}

func (*Animation) Delay() time.Duration {
	return 1 * time.Second
}

func (a *Animation) Image() image.Image {
	return a.image
}

func (a *Animation) Next() error {
	a.escapeRadius -= 0.001
	a.render()
	return nil
}

func interpolateColors(nbColors int) []color.RGBA {
	factor := 1.0 / float64(nbColors)
	steps := []float64{}
	cols := []uint32{}
	interpolated := []uint32{}
	interpolatedColors := []color.RGBA{}
	colors := []color.RGBA{
		{0x00, 0x04, 0x0f, 0xff},
		{0x03, 0x26, 0x28, 0xff},
		{0x07, 0x3e, 0x1e, 0xff},
		{0x18, 0x55, 0x08, 0xff},
		{0x5f, 0x6e, 0x0f, 0xff},
		{0x84, 0x50, 0x19, 0xff},
		{0x9b, 0x30, 0x22, 0xff},
		{0xb4, 0x92, 0x2f, 0xff},
		{0x94, 0xca, 0x3d, 0xff},
		{0x4f, 0xd5, 0x51, 0xff},
		{0x66, 0xff, 0xb3, 0xff},
		{0x82, 0xc9, 0xe5, 0xff},
		{0x9d, 0xa3, 0xeb, 0xff},
		{0xd7, 0xb5, 0xf3, 0xff},
		{0xfd, 0xd6, 0xf6, 0xff},
		{0xff, 0xf0, 0xf2, 0xff},
	}

	for index, col := range colors {
		if index != 0 {
			stepRatio := float64(index+1) / float64(len(colors))
			step := float64(int(stepRatio*100)) / 100 // truncate to 2 decimal precision
			steps = append(steps, step)
		} else {
			steps = append(steps, 0)
		}
		uintColor := uint32(col.R)<<24 | uint32(col.G)<<16 | uint32(col.B)<<8 | uint32(col.A)
		cols = append(cols, uintColor)
	}

	var min, max, minColor, maxColor float64
	if len(colors) == len(steps) && len(colors) == len(cols) {
		for i := 0.0; i <= 1; i += factor {
			for j := 0; j < len(colors)-1; j++ {
				if i >= steps[j] && i < steps[j+1] {
					min = steps[j]
					max = steps[j+1]
					minColor = float64(cols[j])
					maxColor = float64(cols[j+1])
					uintColor := cosineInterpolation(maxColor, minColor, (i-min)/(max-min))
					interpolated = append(interpolated, uint32(uintColor))
				}
			}
		}
	}

	for _, pixelValue := range interpolated {
		interpolatedColors = append(interpolatedColors, uint32ToRgba(pixelValue))
	}

	return interpolatedColors
}

func (a *Animation) render() {
	sz := a.image.Rect.Size()
	width := sz.X
	height := sz.Y
	ratio := float64(width) / float64(height)
	xmin, xmax := a.x-a.escapeRadius/2.0, math.Abs(a.x+a.escapeRadius/2.0)
	ymin, ymax := a.y-a.escapeRadius*ratio/2.0, math.Abs(a.y+a.escapeRadius*ratio/2.0)
	xsize, ysize := xmax-xmin, ymax-ymin

	for iy := 0; iy < height; iy++ {
		for ix := 0; ix < width; ix++ {
			var x = xmin + xsize*float64(ix)/float64(width-1)
			var y = ymin + ysize*float64(iy)/float64(height-1)
			norm, it := mandelIteration(x, y, a.maxIteration)
			iteration := float64(a.maxIteration-it) + math.Log(norm)
			itAbs := int(math.Abs(iteration))
			if itAbs < len(a.colors)-1 {
				color1 := a.colors[itAbs]
				color2 := a.colors[itAbs+1]
				color := linearInterpolation(rgbaToUint(color1), rgbaToUint(color2), uint32(iteration))

				a.image.Set(ix, iy, uint32ToRgba(color))
			}
		}
	}
}

func cosineInterpolation(c1, c2, mu float64) float64 {
	mu2 := (1 - math.Cos(mu*math.Pi)) / 2.0
	return c1*(1-mu2) + c2*mu2
}

func linearInterpolation(c1, c2, mu uint32) uint32 {
	return c1*(1-mu) + c2*mu
}

func mandelIteration(cx, cy float64, maxIter int) (float64, int) {
	x, y, xx, yy := 0.0, 0.0, 0.0, 0.0
	for i := 0; i < maxIter; i++ {
		xy := x * y
		xx = x * x
		yy = y * y
		if xx+yy > 4 {
			return xx + yy, i
		}
		x = xx - yy + cx
		y = 2*xy + cy
	}
	logZn := (x*x + y*y) / 2
	return logZn, maxIter
}

func rgbaToUint(color color.RGBA) uint32 {
	return uint32(color.R)<<24 | uint32(color.G)<<16 | uint32(color.B)<<8 | uint32(color.A)
}

func uint32ToRgba(col uint32) color.RGBA {
	r := col >> 24 & 0xff
	g := col >> 16 & 0xff
	b := col >> 8 & 0xff
	a := 0xff
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}
