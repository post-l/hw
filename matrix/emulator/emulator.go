package emulator

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/post-l/hw/matrix"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
)

type Emulator struct {
	PixelPitch              int
	Gutter                  int
	Width                   int
	Height                  int
	GutterColor             color.Color
	PixelPitchToGutterRatio int
	Margin                  int

	leds    []color.RGBA
	w       screen.Window
	s       screen.Screen
	sz      size.Event
	readyCh chan struct{}
}

func NewEmulator(hc *matrix.HardwareConfig) *Emulator {
	e := &Emulator{
		Width:                   hc.Cols,
		Height:                  hc.Rows,
		GutterColor:             color.Gray{Y: 20},
		PixelPitchToGutterRatio: 2,
		Margin:                  10,
		readyCh:                 make(chan struct{}),

		leds: make([]color.RGBA, hc.Cols*hc.Rows),
	}
	pixelPitch := 12
	e.updatePixelPitchForGutter(pixelPitch / e.PixelPitchToGutterRatio)
	return e
}

// Init initialize the emulator, creating a new Window and waiting until is
// painted. If something goes wrong the function panics
func (e *Emulator) Init() {
	driver.Main(e.mainWindowLoop)
}

func (e *Emulator) Close() error {
	return nil
}

func (e *Emulator) Ready() <-chan struct{} { return e.readyCh }

// ColorModel returns the canvas' color model, always color.RGBAModel
func (e *Emulator) ColorModel() color.Model { return color.RGBAModel }

// Bounds return the topology of the Canvas
func (e *Emulator) Bounds() image.Rectangle { return image.Rect(0, 0, e.Width, e.Height) }

func (e *Emulator) At(x, y int) color.Color {
	pos := x + (y * e.Width)
	return e.leds[pos]
}

func (e *Emulator) Set(x, y int, c color.Color) {
	pos := x + (y * e.Width)
	e.leds[pos] = color.RGBAModel.Convert(c).(color.RGBA)
}

func (e *Emulator) Render() {
	gutterWidth := e.calculateGutterForViewableArea(e.sz.Size())
	e.updatePixelPitchForGutter(gutterWidth)
	e.w.Fill(e.sz.Bounds(), e.GutterColor, screen.Src)
	for col := 0; col < e.Width; col++ {
		for row := 0; row < e.Height; row++ {
			dr := e.ledRect(col, row)
			c := e.At(col, row)
			e.w.Fill(dr, c, screen.Src)
		}
	}
	e.w.Publish()
}

func (e *Emulator) mainWindowLoop(s screen.Screen) {
	e.s = s
	// Calculate initial window size based on whatever our gutter/pixel pitch currently is.
	dims := e.matrixWithMarginsRect()
	wopts := &screen.NewWindowOptions{
		Title:  "RGB LED Matrix Emulator",
		Width:  dims.Max.X,
		Height: dims.Max.Y,
	}
	w, err := s.NewWindow(wopts)
	if err != nil {
		panic(err)
	}
	e.w = w
	defer e.w.Release()
	isReady := false
	for {
		evn := w.NextEvent()
		switch evn := evn.(type) {
		case key.Event:
			if evn.Code == key.CodeEscape {
				os.Exit(42)
			}
		case paint.Event:
			e.Render()
			if !isReady {
				isReady = true
				close(e.readyCh)
			}
		case size.Event:
			e.sz = evn
		case error:
			fmt.Println("render:", err)
		}
	}
}

// Some formulas that allowed me to better understand the drawable area. I found that the math was
// easiest when put in terms of the Gutter width, hence the addition of PixelPitchToGutterRatio.
//
// PixelPitch = PixelPitchToGutterRatio * Gutter
// DisplayWidth = (PixelPitch * LEDColumns) + (Gutter * (LEDColumns - 1)) + (2 * Margin)
// Gutter = (DisplayWidth - (2 * Margin)) / (PixelPitchToGutterRatio * LEDColumns + LEDColumns - 1)
//
//  MMMMMMMMMMMMMMMM.....MMMM
//  MGGGGGGGGGGGGGGG.....GGGM
//  MGLGLGLGLGLGLGLG.....GLGM
//  MGGGGGGGGGGGGGGG.....GGGM
//  MGLGLGLGLGLGLGLG.....GLGM
//  MGGGGGGGGGGGGGGG.....GGGM
//  .........................
//  MGGGGGGGGGGGGGGG.....GGGM
//  MGLGLGLGLGLGLGLG.....GLGM
//  MGGGGGGGGGGGGGGG.....GGGM
//  MMMMMMMMMMMMMMMM.....MMMM
//
//  where:
//    M = Margin
//    G = Gutter
//    L = LED

// matrixWithMarginsRect Returns a Rectangle that describes entire emulated RGB Matrix, including margins.
func (e *Emulator) matrixWithMarginsRect() image.Rectangle {
	upperLeftLED := e.ledRect(0, 0)
	lowerRightLED := e.ledRect(e.Width-1, e.Height-1)
	return image.Rect(upperLeftLED.Min.X-e.Margin, upperLeftLED.Min.Y-e.Margin, lowerRightLED.Max.X+e.Margin, lowerRightLED.Max.Y+e.Margin)
}

// ledRect Returns a Rectangle for the LED at col and row.
func (e *Emulator) ledRect(col int, row int) image.Rectangle {
	x := (col * (e.PixelPitch + e.Gutter)) + e.Margin
	y := (row * (e.PixelPitch + e.Gutter)) + e.Margin
	return image.Rect(x, y, x+e.PixelPitch, y+e.PixelPitch)
}

// calculateGutterForViewableArea As the name states, calculates the size of the gutter for a given viewable area.
// It's easier to understand the geometry of the matrix on screen when put in terms of the gutter,
// hence the shift toward calculating the gutter size.
func (e *Emulator) calculateGutterForViewableArea(size image.Point) int {
	maxGutterInX := (size.X - 2*e.Margin) / (e.PixelPitchToGutterRatio*e.Width + e.Width - 1)
	maxGutterInY := (size.Y - 2*e.Margin) / (e.PixelPitchToGutterRatio*e.Height + e.Height - 1)
	if maxGutterInX < maxGutterInY {
		return maxGutterInX
	}
	return maxGutterInY
}

func (e *Emulator) updatePixelPitchForGutter(gutterWidth int) {
	e.PixelPitch = e.PixelPitchToGutterRatio * gutterWidth
	e.Gutter = gutterWidth
}
