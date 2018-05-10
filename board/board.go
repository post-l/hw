package board

import (
	"io"
)

type Board interface {
	io.Closer
	SetPinMode(pin int, mode PinMode)
	DigitalRead(pin int) bool
	DigitalWrite(pin int, v bool)
	DigitalWrites([]PinValue)
}

type PinMode int

const (
	Input = PinMode(iota + 1)
	Output
)

type PinValue struct {
	Pin   int
	Value bool
}
