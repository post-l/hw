package matrix

type ScanMode int

const (
	Progressive = ScanMode(iota + 1)
	Interlaced
)

// DefaultHardwareConfig default WS281x configuration
var DefaultHardwareConfig = HardwareConfig{
	Rows:       64,
	Cols:       64,
	Brightness: 100,
	ScanMode:   Interlaced,
	Mapping:    DefaultHardwareMapping,
}

// HardwareConfig rgb-led-matrix configuration
type HardwareConfig struct {
	// Rows the number of rows supported by the display, so 32 or 16.
	Rows int
	// Cols the number of columns supported by the display, so 32 or 64 .
	Cols int
	// Brightness is the initial brightness of the panel in percent. Valid range
	// is 1..100
	Brightness int
	// ScanMode progressive or interlaced
	ScanMode ScanMode // strip color layout
	Mapping  HardwareMapping
}
