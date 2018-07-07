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

	"github.com/post-l/hw/examples"
	"github.com/post-l/hw/matrix/toolkit"
)

var (
	gifPath = flag.String("path", "", "gif path or url")
)

func main() {
	examples.Main(run)
}

func run(m toolkit.Matrix) error {
	gif, err := getGIF(*gifPath)
	if err != nil {
		return err
	}
	tk := toolkit.New(m)
	ctx := context.Background()
	return tk.PlayGIF(ctx, gif)
}

func getGIF(p string) (*gif.GIF, error) {
	var r io.Reader
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		resp, err := http.Get(p)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("invalid http status code %d", resp.StatusCode)
		}
		r = resp.Body
	} else {
		f, err := os.Open(p)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		r = f
	}
	return gif.DecodeAll(r)
}
