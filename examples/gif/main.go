package main

import (
	"context"
	"flag"
	"image/gif"
	"os"

	"github.com/post-l/hw/board/tinkerboard"
	"github.com/post-l/hw/matrix"
	"github.com/post-l/hw/matrix/emulator"
	"github.com/post-l/hw/matrix/toolkit"
)

var (
	img    = flag.String("gif", "gopher-dance-long-3x.gif", "gif path")
	emFlag = flag.Bool("emulator", false, "use emulator")
)

func main() {
	flag.Parse()

	f, err := os.Open(*img)
	fatal(err)
	gif, err := gif.DecodeAll(f)
	fatal(err)
	f.Close()

	if *emFlag {
		m := emulator.NewEmulator(&matrix.DefaultHardwareConfig)
		go func() {
			<-m.Ready()
			run(gif, m)
		}()
		m.Init()
	}
	b, err := tinkerboard.New()
	fatal(err)
	m := matrix.New(b, &matrix.DefaultHardwareConfig)
	defer m.Close()
	run(gif, m)
}

func run(gif *gif.GIF, m toolkit.Matrix) {
	tk := toolkit.New(m)
	ctx := context.Background()
	err := tk.PlayGIF(ctx, gif)
	fatal(err)
}

func fatal(err error) {
	if err != nil {
		panic(err)
	}
}
