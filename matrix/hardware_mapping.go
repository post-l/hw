package matrix

import "github.com/post-l/hw/board/tinkerboard"

var DefaultHardwareMapping = HardwareMapping{
	outputEnable: tinkerboard.GPIO0_C1,
	clock:        tinkerboard.GPIO5_B4,
	strobe:       tinkerboard.GPIO6_A4,

	a: tinkerboard.GPIO5_B7,
	b: tinkerboard.GPIO7_B0,
	c: tinkerboard.GPIO5_B6,
	d: tinkerboard.GPIO6_A3,
	e: tinkerboard.GPIO5_B3,

	r1: tinkerboard.GPIO5_B5,
	g1: tinkerboard.GPIO5_C0,
	b1: tinkerboard.GPIO7_C6,

	r2: tinkerboard.GPIO7_C7,
	g2: tinkerboard.GPIO5_B2,
	b2: tinkerboard.GPIO7_A7,
}

type HardwareMapping struct {
	outputEnable int
	clock        int
	strobe       int

	a, b, c, d, e int
	r1, g1, b1    int
	r2, g2, b2    int
}

func (hm *HardwareMapping) pins() []int {
	return []int{
		hm.outputEnable, hm.clock, hm.strobe,
		hm.a, hm.b, hm.c, hm.d, hm.e,
		hm.r1, hm.g1, hm.b1,
		hm.r2, hm.g2, hm.b2,
	}
}
