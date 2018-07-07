package matrix

type ScanMode int

const (
	Progressive = ScanMode(iota + 1)
	Interlaced
)

// DefaultHardwareConfig default WS281x configuration
var DefaultHardwareConfig = HardwareConfig{
	Rows:            64,
	Cols:            64,
	PWMBits:         5,
	Brightness:      100,
	ScanMode:        Interlaced,
	Mapping:         DefaultHardwareMapping,
	ShowRefreshRate: true,
}

// HardwareConfig rgb-led-matrix configuration
type HardwareConfig struct {
	// Rows the number of rows supported by the display, so 32 or 16.
	Rows int
	// Cols the number of columns supported by the display, so 32 or 64 .
	Cols int
	// PWMBits sets PWM bits used for output. Default is 11, but if you only deal with
	// limited comic-colors, 1 might be sufficient. Lower require less CPU and
	// increases refresh-rate.
	PWMBits int
	// Brightness is the initial brightness of the panel in percent. Valid range
	// is 1..100
	Brightness int
	// ScanMode progressive or interlaced
	ScanMode        ScanMode // strip color layout
	Mapping         HardwareMapping
	ShowRefreshRate bool
}
