package tinkerboard_test

import (
	"testing"

	"github.com/post-l/hw/board"
	"github.com/post-l/hw/board/tinkerboard"
)

func TestBoard(t *testing.T) {
	tb, err := tinkerboard.New()
	if err != nil {
		t.Fatalf("expect New to return no error: %v", err)
	}
	tt := []struct {
		name string
		pin  int
	}{
		{"GPIO0_C1", tinkerboard.GPIO0_C1},
		{"GPIO5_B2", tinkerboard.GPIO5_B2},
		{"GPIO5_B3", tinkerboard.GPIO5_B3},
		{"GPIO5_B4", tinkerboard.GPIO5_B4},
		{"GPIO5_B5", tinkerboard.GPIO5_B5},
		{"GPIO5_B6", tinkerboard.GPIO5_B6},
		{"GPIO5_B7", tinkerboard.GPIO5_B7},
		{"GPIO5_C0", tinkerboard.GPIO5_C0},
		{"GPIO6_A3", tinkerboard.GPIO6_A3},
		{"GPIO6_A4", tinkerboard.GPIO6_A4},
		{"GPIO7_A7", tinkerboard.GPIO7_A7},
		{"GPIO7_B0", tinkerboard.GPIO7_B0},
		{"GPIO7_C6", tinkerboard.GPIO7_C6},
		{"GPIO7_C7", tinkerboard.GPIO7_C7},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			pin := tc.pin
			tb.SetPinMode(pin, board.Output)
			tb.DigitalWrite(pin, true)
			if got, want := tb.DigitalRead(pin), true; got != want {
				t.Errorf("invalid digital read value: got %v; want %v", got, want)
			}
			tb.DigitalWrite(pin, false)
			if got, want := tb.DigitalRead(pin), false; got != want {
				t.Errorf("invalid digital read value: got %v; want %v", got, want)
			}
		})
	}
}
