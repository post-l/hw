// An implementation of Conway's Game of Life.
package life

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"time"
)

type Animation struct {
	l *Life
}

func NewAnimation(sz image.Point) *Animation {
	return &Animation{
		l: New(sz.X, sz.Y),
	}
}

func (*Animation) Delay() time.Duration {
	return time.Second / 10
}

func (a *Animation) Image() image.Image {
	return a.l.a
}

func (a *Animation) Next() error {
	a.l.Step()
	return nil
}

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]color.RGBA
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]color.RGBA, h)
	for i := range s {
		s[i] = make([]color.RGBA, w)
	}
	return &Field{s: s, w: w, h: h}
}

// ColorModel returns the canvas' color model, always color.RGBAModel
func (*Field) ColorModel() color.Model { return color.RGBAModel }

// Bounds return the topology of the Canvas
func (f *Field) Bounds() image.Rectangle { return image.Rect(0, 0, f.w, f.h) }

func (f *Field) At(x, y int) color.Color {
	return f.s[y][x]
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, c color.RGBA) {
	f.s[y][x] = c
}

// Get returns the specified cell.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Get(x, y int) color.RGBA {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) color.RGBA {
	// Count the adjacent cells that are alive.
	alive := 0
	var r, g, b uint32
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if j != 0 || i != 0 {
				c := f.Get(x+i, y+j)
				if c != (color.RGBA{}) {
					alive++
					r += uint32(c.R)
					g += uint32(c.G)
					b += uint32(c.B)
				}
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	if alive == 3 {
		return color.RGBA{
			A: 255,
			R: uint8(math.Min(float64(r)*0.3336666667, 255)),
			G: uint8(math.Min(float64(g)*0.3336666667, 255)),
			B: uint8(math.Min(float64(b)*0.3336666667, 255)),
		}
	} else if alive == 2 {
		return f.s[y][x]
	}
	return color.RGBA{}
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b *Field
	w, h int
}

// New returns a new Life game state with a random initial state.
func New(w, h int) *Life {
	a := NewField(w, h)
	c := color.RGBA{A: 255, R: 231, G: 76, B: 60}
	nbCells := w * h / 4
	nbCCells := nbCells / 3
	for i := 0; i < nbCells; i++ {
		if i == nbCCells {
			c = color.RGBA{A: 255, R: 46, G: 204, B: 113}
		} else if i == nbCCells*2 {
			c = color.RGBA{A: 255, R: 52, G: 152, B: 219}
		}
		a.Set(rand.Intn(w), rand.Intn(h), c)
	}
	return &Life{
		a: a, b: NewField(w, h),
		w: w, h: h,
	}
}

// Step advances the game by one instant, recomputing and updating all cells.
func (l *Life) Step() {
	// Update the state of the next field (b) from the current field (a).
	for y := 0; y < l.h; y++ {
		for x := 0; x < l.w; x++ {
			l.b.Set(x, y, l.a.Next(x, y))
		}
	}
	// Swap fields a and b.
	l.a, l.b = l.b, l.a
}
