package matrix

import (
	"context"
	"image"
	"image/color"
	"time"

	"github.com/post-l/hw/board"
)

const bitPlanes = 11

type Matrix struct {
	b  board.Board
	hc *HardwareConfig

	leds      []led
	dRows     int
	dRowAddrs [][]board.PinValue

	ctx    context.Context
	cancel context.CancelFunc
}

func New(b board.Board, hc *HardwareConfig) *Matrix {
	hm := hc.Mapping
	for _, pin := range hm.pins() {
		b.SetPinMode(pin, board.Output)
	}

	dRows := hc.Rows / 2
	dRowAddrs := make([][]board.PinValue, dRows)
	for i := 0; i < dRows; i++ {
		dRowAddrs[i] = []board.PinValue{
			{Pin: hm.a, Value: (i & 1) != 0},
			{Pin: hm.b, Value: (i & 2) != 0},
			{Pin: hm.c, Value: (i & 4) != 0},
			{Pin: hm.d, Value: (i & 8) != 0},
			{Pin: hm.e, Value: (i & 16) != 0},
		}
	}

	ctx, cancel := context.WithCancel(context.Background())

	m := &Matrix{
		b:  b,
		hc: hc,

		leds:      make([]led, hc.Cols*hc.Rows),
		dRows:     dRows,
		dRowAddrs: dRowAddrs,

		ctx:    ctx,
		cancel: cancel,
	}
	go m.run()
	return m
}

func (m *Matrix) Close() error {
	m.cancel()
	return nil
}

// ColorModel returns the canvas' color model, always color.RGBAModel
func (m *Matrix) ColorModel() color.Model { return color.RGBAModel }

// Bounds return the topology of the Canvas
func (m *Matrix) Bounds() image.Rectangle { return image.Rect(0, 0, m.hc.Cols, m.hc.Rows) }

func (m *Matrix) At(x, y int) color.Color {
	pos := x + (y * m.hc.Cols)
	return m.leds[pos]
}

func (m *Matrix) Set(x, y int, c color.Color) {
	pos := x + (y * m.hc.Cols)
	co := color.RGBAModel.Convert(c).(color.RGBA)
	m.leds[pos] = led{
		r: cie[co.R],
		g: cie[co.G],
		b: cie[co.B],
	}
}

func (m *Matrix) Render() {
	// TODO: Sync/Swap image buffer.
}

func (m *Matrix) run() {
	for {
		m.render()
		select {
		case <-m.ctx.Done():
			return
		default:
		}
	}
}

func (m *Matrix) render() {
	hm := m.hc.Mapping
	hdRows := m.dRows / 2
	padding := m.dRows * m.hc.Cols
	for row := 0; row < m.dRows; row++ {
		drow := row
		if m.hc.ScanMode == Interlaced {
			if row < hdRows {
				drow = row << 1
			} else {
				drow = ((row - hdRows) << 1) + 1
			}
		}
		drowAddr := m.dRowAddrs[drow]
		m.b.DigitalWrites(drowAddr)

		for x := uint(8); x < bitPlanes; x++ {
			curBit := uint16(1 << x)
			for col := 0; col < m.hc.Cols; col++ {
				p1 := col + drow*m.hc.Rows
				l1 := m.leds[p1]
				p2 := p1 + padding
				l2 := m.leds[p2]
				data := []board.PinValue{
					{hm.r1, (l1.r & curBit) == curBit},
					{hm.g1, (l1.g & curBit) == curBit},
					{hm.b1, (l1.b & curBit) == curBit},

					{hm.r2, (l2.r & curBit) == curBit},
					{hm.g2, (l2.g & curBit) == curBit},
					{hm.b2, (l2.b & curBit) == curBit},

					{hm.clock, false},
				}
				m.b.DigitalWrites(data)
				m.b.DigitalWrite(hm.clock, true)
			}
			data := []board.PinValue{
				{hm.r1, false},
				{hm.g1, false},
				{hm.b1, false},

				{hm.r2, false},
				{hm.g2, false},
				{hm.b2, false},

				{hm.clock, false},
			}
			m.b.DigitalWrites(data)
			m.b.DigitalWrite(hm.strobe, true)
			m.b.DigitalWrite(hm.strobe, false)
			m.b.DigitalWrite(hm.outputEnable, false)
			t := 130 << uint(x)
			if t > 28000 {
				d := time.Duration(t)
				time.Sleep(d / 2)
			} else {
				for i := t >> 3; i != 0; i-- {
				}
			}
			m.b.DigitalWrite(hm.outputEnable, true)
		}
	}
}

type led struct {
	r, g, b uint16
}

func (c led) RGBA() (r, g, b, a uint32) {
	return uint32(c.r), uint32(c.g), uint32(c.b), 0
}
