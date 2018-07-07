package main

import (
	"context"
	"fmt"
	"image/gif"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/peterhellberg/giphy"
	"github.com/post-l/hw/examples"
	"github.com/post-l/hw/matrix/toolkit"
)

func main() {
	examples.Main(run)
}

func run(m toolkit.Matrix) error {
	tk := toolkit.New(m)
	client := giphy.NewClient()
	client.Limit = 100
	// Artists: @robindavey @mrdiv @patakk @64-x-64
	res, err := client.Search([]string{"art neon trippy"})
	if err != nil {
		return err
	}
	if res.Meta.Status != http.StatusOK {
		return fmt.Errorf("invalid status %d: %s", res.Meta.Status, res.Meta.Msg)
	}
	for {
		i := rand.Intn(len(res.Data))
		item := res.Data[i]
		if err := playGIFFromURL(tk, item.Images.FixedWidth.URL); err != nil {
			fmt.Fprintf(os.Stderr, "could not play gif: %v\n", err)
		}
	}
}

func playGIFFromURL(tk *toolkit.ToolKit, urlStr string) error {
	gif, err := getGIF(urlStr)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := tk.PlayGIF(ctx, gif); err != nil && err != context.DeadlineExceeded {
		return err
	}
	return nil
}

func getGIF(p string) (*gif.GIF, error) {
	resp, err := http.Get(p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid http status code %d", resp.StatusCode)
	}
	return gif.DecodeAll(resp.Body)
}
