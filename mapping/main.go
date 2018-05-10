package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"text/tabwriter"
)

type HardwareMapping struct {
	outputEnable int
	clock        int
	strobe       int

	a, b, c, d, e int
	r1, g1, b1    int
	r2, g2, b2    int
}

var hw = HardwareMapping{
	outputEnable: 4,
	clock:        17,
	strobe:       21,

	a: 22,
	b: 26,
	c: 27,
	d: 20,
	e: 24,

	r1: 5,
	g1: 13,
	b1: 6,

	r2: 12,
	g2: 16,
	b2: 23,
}

var physToGpioR2 = []int{
	-1,     // 0
	-1, -1, // 1, 2
	2, -1,
	3, -1,
	4, 14,
	-1, 15,
	17, 18,
	27, -1,
	22, 23,
	-1, 24,
	10, -1,
	9, 25,
	11, 8,
	-1, 7, // 25, 26

	// B+

	0, 1,
	5, -1,
	6, 12,
	13, -1,
	19, 16,
	26, 20,
	-1, 21,

	// the P5 connector on the Rev 2 boards:

	-1, -1,
	-1, -1,
	-1, -1,
	-1, -1,
	-1, -1,
	28, 29,
	30, 31,
	-1, -1,
	-1, -1,
	-1, -1,
	-1, -1,
}

var gpioToPin = map[string]int{
	"GPIO5_B0": (8 + 152),
	"GPIO5_B1": (9 + 152),
	"GPIO5_B2": (10 + 152),
	"GPIO5_B3": (11 + 152),
	"GPIO5_B4": (12 + 152),
	"GPIO5_B5": (13 + 152),
	"GPIO5_B6": (14 + 152),
	"GPIO5_B7": (15 + 152),
	"GPIO5_C0": (16 + 152),
	"GPIO5_C1": (17 + 152),
	"GPIO5_C2": (18 + 152),
	"GPIO5_C3": (19 + 152),

	"GPIO6_A0": (184),
	"GPIO6_A1": (1 + 184),
	"GPIO6_A2": (2 + 184),
	"GPIO6_A3": (3 + 184),
	"GPIO6_A4": (4 + 184),

	"GPIO7_A0": (0 + 216),
	"GPIO7_A7": (7 + 216),
	"GPIO7_B0": (8 + 216),
	"GPIO7_B1": (9 + 216),
	"GPIO7_B2": (10 + 216),
	"GPIO7_C1": (17 + 216),
	"GPIO7_C2": (18 + 216),
	"GPIO7_C6": (22 + 216),
	"GPIO7_C7": (23 + 216),

	"GPIO8_A3": (3 + 248),
	"GPIO8_A4": (4 + 248),
	"GPIO8_A5": (5 + 248),
	"GPIO8_A6": (6 + 248),
	"GPIO8_A7": (7 + 248),
	"GPIO8_B0": (8 + 248),
	"GPIO8_B1": (9 + 248),
}

var pinToGpioStr = []string{
	"-1",       // 0
	"-1", "-1", //1, 2
	"GPIO8_A4", "-1", //3, 4
	"GPIO8_A5", "-1", //5, 6
	"GPIO0_C1", "GPIO5_B1", //7, 8
	"-1", "GPIO5_B0", //9, 10
	"GPIO5_B4", "GPIO6_A0", //11, 12
	"GPIO5_B6", "-1", //13, 14
	"GPIO5_B7", "GPIO5_B2", //15, 16
	"-1", "GPIO5_B3", //17, 18
	"GPIO8_B1", "-1", //19, 20
	"GPIO8_B0", "GPIO5_C3", //21, 22
	"GPIO8_A6", "GPIO8_A7", //23, 24
	"-1", "GPIO8_A3", //25, 26
	"GPIO7_C1", "GPIO7_C2", //27, 28
	"GPIO5_B5", "-1", //29, 30
	"GPIO5_C0", "GPIO7_C7", //31, 32
	"GPIO7_C6", "-1", //33, 34
	"GPIO6_A1", "GPIO7_A7", //35, 36
	"GPIO7_B0", "GPIO6_A3", //37, 38
	"-1", "GPIO6_A4", //39, 40
}

func gpioToBank(gpio int) int {
	if gpio < 24 {
		return 0
	}
	return ((gpio - 24) / 32) + 1
}

func gpioToBankPin(gpio int) int {
	if gpio < 24 {
		return gpio
	}
	return (gpio - 24) % 32
}

type FieldMapping struct {
	name    string
	piPin   int
	pin     int
	gpio    int
	bank    int
	bankPin int
	gpioStr string
}

type FieldsMapping []*FieldMapping

func (f FieldsMapping) Len() int           { return len(f) }
func (f FieldsMapping) Less(i, j int) bool { return f[i].gpioStr < f[j].gpioStr }
func (f FieldsMapping) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func main() {
	piPinToPhysPin := make(map[int]int)
	for physPin, piPin := range physToGpioR2 {
		if piPin != -1 {
			piPinToPhysPin[piPin] = physPin
		}
	}
	t := reflect.TypeOf(hw)
	v := reflect.ValueOf(hw)
	var fms FieldsMapping
	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Name
		piPin := int(v.Field(i).Int())
		pin := piPinToPhysPin[piPin]
		gpioStr := pinToGpioStr[pin]
		gpio := gpioToPin[gpioStr]
		bank := gpioToBank(gpio)
		bankPin := gpioToBankPin(gpio)
		fms = append(fms, &FieldMapping{
			name:    name,
			piPin:   piPin,
			pin:     pin,
			gpioStr: gpioStr,
			gpio:    gpio,
			bank:    bank,
			bankPin: bankPin,
		})
	}
	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Name\tPi Pin\tPhys Pin\tGpio\tBank\tBank Pin\tGpio Name\n")
	for _, fm := range fms {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\t%d\t%s\n", fm.name, fm.piPin, fm.pin, fm.gpio, fm.bank, fm.bankPin, fm.gpioStr)
	}
	tw.Flush()
	fmt.Fprintf(os.Stdout, "\n\n")
	sort.Sort(fms)
	tw = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "Name\tPi Pin\tPhys Pin\tGpio\tBank\tBank Pin\tGpio Name\n")
	for _, fm := range fms {
		fmt.Fprintf(tw, "%s\t%d\t%d\t%d\t%d\t%d\t%s\n", fm.name, fm.piPin, fm.pin, fm.gpio, fm.bank, fm.bankPin, fm.gpioStr)
	}
	tw.Flush()
}
