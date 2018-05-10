// An implementation of Conway's Game of Life.
package main

import (
	"bytes"
	"context"
	"flag"
	"image"
	"image/color"
	"log"
	"math/rand"
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
	rand.Seed(time.Now().UnixNano())

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
	a := NewAnimation(m.Bounds())
	tk.PlayAnimation(ctx, a)
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Animation struct {
	l   *Life
	img *image.RGBA
}

func NewAnimation(r image.Rectangle) *Animation {
	sz := r.Size()
	return &Animation{
		l:   NewLife(sz.X, sz.Y),
		img: image.NewRGBA(r),
	}
}

func (*Animation) Delay() time.Duration {
	return time.Second / 20
}

func (a *Animation) Image() image.Image {
	for y := 0; y < a.l.h; y++ {
		for x := 0; x < a.l.w; x++ {
			var c color.Color
			if a.l.a.Alive(x, y) {
				c = color.RGBA{255, 0, 0, 255}
			} else {
				c = color.Black
			}
			a.img.Set(x, y, c)
		}
	}
	return a.img
}

func (a *Animation) Next() error {
	a.l.Step()
	return nil
}

// Field represents a two-dimensional field of cells.
type Field struct {
	s    [][]bool
	w, h int
}

// NewField returns an empty field of the specified width and height.
func NewField(w, h int) *Field {
	s := make([][]bool, h)
	for i := range s {
		s[i] = make([]bool, w)
	}
	return &Field{s: s, w: w, h: h}
}

// Set sets the state of the specified cell to the given value.
func (f *Field) Set(x, y int, b bool) {
	f.s[y][x] = b
}

// Alive reports whether the specified cell is alive.
// If the x or y coordinates are outside the field boundaries they are wrapped
// toroidally. For instance, an x value of -1 is treated as width-1.
func (f *Field) Alive(x, y int) bool {
	x += f.w
	x %= f.w
	y += f.h
	y %= f.h
	return f.s[y][x]
}

// Next returns the state of the specified cell at the next time step.
func (f *Field) Next(x, y int) bool {
	// Count the adjacent cells that are alive.
	alive := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (j != 0 || i != 0) && f.Alive(x+i, y+j) {
				alive++
			}
		}
	}
	// Return next state according to the game rules:
	//   exactly 3 neighbors: on,
	//   exactly 2 neighbors: maintain current state,
	//   otherwise: off.
	return alive == 3 || alive == 2 && f.Alive(x, y)
}

// Life stores the state of a round of Conway's Game of Life.
type Life struct {
	a, b *Field
	w, h int
}

// NewLife returns a new Life game state with a random initial state.
func NewLife(w, h int) *Life {
	a := NewField(w, h)
	for i := 0; i < (w * h / 4); i++ {
		a.Set(rand.Intn(w), rand.Intn(h), true)
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

// String returns the game board as a string.
func (l *Life) String() string {
	var buf bytes.Buffer
	for x := 0; x < l.w+2; x++ {
		buf.WriteByte('-')
	}
	buf.WriteByte('\n')
	for y := 0; y < l.h; y++ {
		buf.WriteByte('|')
		for x := 0; x < l.w; x++ {
			b := byte(' ')
			if l.a.Alive(x, y) {
				b = '*'
			}
			buf.WriteByte(b)
		}
		buf.WriteByte('|')
		buf.WriteByte('\n')
	}
	for x := 0; x < l.w+2; x++ {
		buf.WriteByte('-')
	}
	buf.WriteByte('\n')
	return buf.String()
}
