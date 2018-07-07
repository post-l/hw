package examples

import (
	"flag"
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

func Main(run func(toolkit.Matrix) error) {
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	if *emFlag {
		m := emulator.NewEmulator(&matrix.DefaultHardwareConfig)
		go func() {
			if err := run(m); err != nil {
				log.Println("run:", err)
			}
			m.Close()
		}()
		m.Run()
	} else {
		b, err := tinkerboard.New()
		if err != nil {
			log.Fatal("board:", err)
		}
		m := matrix.New(b, &matrix.DefaultHardwareConfig)
		defer m.Close()
		if err := run(m); err != nil {
			log.Fatal("run:", err)
		}
	}
}
