package main

import (
	"context"
	"flag"
	"fmt"
	"image/gif"
	"io"
	"net/http"
	"os"
	"strings"

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

	gif := getGIF(*img)
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

func getGIF(p string) *gif.GIF {
	var r io.Reader
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		resp, err := http.Get(p)
		fatal(err)
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("invalid http status code %d", resp.StatusCode)
			fatal(err)
		}
		r = resp.Body
	} else {
		f, err := os.Open(p)
		fatal(err)
		defer f.Close()
		r = f
	}
	gif, err := gif.DecodeAll(r)
	fatal(err)
	return gif
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
