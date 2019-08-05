package matrix

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"math"
	"time"

	"github.com/post-l/hw/board"
	"github.com/post-l/hw/board/tinkerboard"
)

const pwmBitsLen = 11

var vreal = []int{
	53,    // 130
	120,   // 260
	250,   // 520
	510,   // 1040
	1000,  // 2080
	2100,  // 4160
	4800,  // 8320
	10000, // 16640
	30000, // 33280
	60000, // 66560
	37000, // 133120
}

type Matrix struct {
	b  *tinkerboard.TinkerBoard
	hc *HardwareConfig

	buf       []uint8
	bbuf      []uint8
	dRows     int
	dRowAddrs []*tinkerboard.BankWriter

	colorClkMask *tinkerboard.BankWriter
	data         *tinkerboard.BankWriter

	pwmStartBit int

	cie [256]uint16

	swapc chan struct{}

	ctx    context.Context
	cancel context.CancelFunc
}

func New(b board.Board, hc *HardwareConfig) *Matrix {
	hm := hc.Mapping
	for _, pin := range hm.pins() {
		b.SetPinMode(pin, board.Output)
	}

	dRows := hc.Rows / 2
	dRowAddrs := make([]*tinkerboard.BankWriter, dRows)
	addrPins := []int{hm.a, hm.b, hm.c, hm.d, hm.e}
	for i := uint32(0); i < uint32(dRows); i++ {
		bw := tinkerboard.NewBankWriter(addrPins)
		bw.Set(i)
		dRowAddrs[i] = bw
	}

	ctx, cancel := context.WithCancel(context.Background())

	colorPins := []int{hm.r1, hm.g1, hm.b1, hm.r2, hm.g2, hm.b2, hm.clock}
	bufSize := hc.PWMBits * hc.Cols * dRows

	m := &Matrix{
		b:  b.(*tinkerboard.TinkerBoard),
		hc: hc,

		buf:       make([]uint8, bufSize),
		bbuf:      make([]uint8, bufSize),
		dRows:     dRows,
		dRowAddrs: dRowAddrs,

		colorClkMask: tinkerboard.NewBankWriter(colorPins),
		data:         tinkerboard.NewBankWriter(colorPins),

		pwmStartBit: pwmBitsLen - hc.PWMBits,

		swapc: make(chan struct{}),

		ctx:    ctx,
		cancel: cancel,
	}
	m.createLuminanceCIETable(hc.Brightness, hc.PWMBits)
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

func (m *Matrix) At(x, y int) color.Color { return color.RGBA{} }

func (m *Matrix) Set(x, y int, c color.Color) {
	var colorMask, roffset, goffset, boffset uint8
	if y >= m.dRows {
		colorMask = 7 // 0b000111
		roffset, goffset, boffset = 8, 16, 32
		y -= m.dRows
	} else {
		colorMask = 56 // 0b111000
		roffset, goffset, boffset = 1, 2, 4
	}
	i := x + y*m.hc.Cols*m.hc.PWMBits
	co := color.RGBAModel.Convert(c).(color.RGBA)
	r := m.cie[co.R]
	g := m.cie[co.G]
	b := m.cie[co.B]
	for bit := uint(0); bit < uint(m.hc.PWMBits); bit++ {
		colorBits := m.bbuf[i] & colorMask
		mask := uint16(1 << bit)
		if r&mask != 0 {
			colorBits |= roffset
		}
		if g&mask != 0 {
			colorBits |= goffset
		}
		if b&mask != 0 {
			colorBits |= boffset
		}
		m.bbuf[i] = colorBits
		i += m.hc.Cols
	}
}

func (m *Matrix) PWMBits() int { return pwmBitsLen - m.pwmStartBit }

// SetPWMBits sets PWM bits used for output. Default is 11, but if you only deal with
// limited comic-colors, 1 might be sufficient. Lower require less CPU and
// increases refresh-rate.
func (m *Matrix) SetPWMBits(pwmBits int) {
	m.pwmStartBit = pwmBitsLen - pwmBits
	m.createLuminanceCIETable(m.hc.Brightness, pwmBits)
}

// Render renders the back buffer. It waits to the next VSync and
// swaps the active buffer with the back buffer one.
func (m *Matrix) Render() {
	m.swapc <- struct{}{}
}

func (m *Matrix) run() {
	i := 0
	var tc <-chan time.Time
	if m.hc.ShowRefreshRate {
		t := time.NewTicker(1 * time.Second)
		defer t.Stop()
		tc = t.C
	}
	for {
		m.render()
		i++
		select {
		case <-m.swapc:
			m.buf, m.bbuf = m.bbuf, m.buf
		case <-tc:
			fmt.Println(i, "fps")
			i = 0
		case <-m.ctx.Done():
			return
		default:
		}
	}
}

func (m *Matrix) render() {
	hm := m.hc.Mapping
	hdRows := m.dRows / 2
	colSize := m.hc.Cols * m.hc.PWMBits
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
		m.b.PerfWrites(drowAddr)

		i := drow * colSize
		for x := m.pwmStartBit; x < pwmBitsLen; x++ {
			for col := 0; col < m.hc.Cols; col++ {
				v := uint32(m.buf[i])
				m.data.Set(v)
				m.b.PerfWrites(m.data)
				m.b.DigitalWrite(hm.clock, true)
				i++
			}

			m.b.PerfWrites(m.colorClkMask)

			m.b.DigitalWrite(hm.strobe, true)
			m.b.DigitalWrite(hm.strobe, false)

			m.b.DigitalWrite(hm.outputEnable, false)
			if x == pwmBitsLen-1 {
				d := time.Duration(vreal[x])
				time.Sleep(d)
			} else {
				for i := vreal[x]; i != 0; i-- {
				}
			}
			m.b.DigitalWrite(hm.outputEnable, true)
		}
	}
}

func (m *Matrix) createLuminanceCIETable(brightness, pwmBits int) {
	outFactor := (1 << uint(pwmBits)) - 1
	for c := range m.cie {
		m.cie[c] = uint16(float64(outFactor) * luminanceCIE1931(uint8(c), uint8(brightness)))
	}
}

// Do CIE1931 luminance correction and scale to output bitplanes.
func luminanceCIE1931(c, brightness uint8) float64 {
	v := float64(c) * float64(brightness) / 255.0
	if v <= 8 {
		return v / 902.3
	}
	return math.Pow((v+16)/116.0, 3)
}
