package main

import (
	"context"
	"fmt"
	"image/gif"
	"net/http"
	"os"
	"time"

	"github.com/peterhellberg/giphy"
	"github.com/post-l/hw/examples"
	"github.com/post-l/hw/examples/animation/circle"
	"github.com/post-l/hw/examples/animation/life"
	"github.com/post-l/hw/examples/animation/mandelbrot"
	"github.com/post-l/hw/examples/animation/text"
	"github.com/post-l/hw/matrix/toolkit"
)

type pwmBitser interface {
	PWMBits() int
	SetPWMBits(pwmBits int)
}

func main() {
	examples.Main(run)
}

func run(m toolkit.Matrix) error {
	tk := toolkit.New(m)
	sz := m.Bounds().Size()

	ca := circle.NewAnimation(sz)

	ta, err := text.NewAnimation(sz)
	if err != nil {
		return err
	}

	ma := mandelbrot.NewAnimation(sz)

	animations := []func(context.Context){
		func(ctx context.Context) { tk.PlayAnimation(ctx, ma) },
		func(ctx context.Context) { randTextAnim(ctx, m, tk, ta) },
		func(ctx context.Context) { tk.PlayAnimation(ctx, ca) },
		func(ctx context.Context) { tk.PlayAnimation(ctx, life.NewAnimation(sz)) },
		func(ctx context.Context) { randGIFFromGiphy(ctx, tk) },
	}

	for {
		for _, anim := range animations {
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			anim(ctx)
			cancel()
		}
	}
}

func randTextAnim(ctx context.Context, m toolkit.Matrix, tk *toolkit.ToolKit, ta *text.Animation) {
	if err := ta.RandQuote(); err != nil {
		fmt.Fprintf(os.Stderr, "could not get random quote: %v\n", err)
		return
	}
	if v, ok := m.(pwmBitser); ok {
		defer v.SetPWMBits(v.PWMBits())
		v.SetPWMBits(3)
	}
	tk.PlayAnimation(ctx, ta)
}

func randGIFFromGiphy(ctx context.Context, tk *toolkit.ToolKit) {
	res, err := giphy.DefaultClient.Random([]string{"art neon trippy"})
	if err != nil || res.Meta.Status != http.StatusOK {
		fmt.Fprintf(os.Stderr, "could not query giphy: %v\n", err)
		return
	}
	gif, err := getGIF(ctx, res.Data.ImageURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not get gif: %v\n", err)
		return
	}
	if err := tk.PlayGIF(ctx, gif); err != nil && err != context.DeadlineExceeded {
		fmt.Fprintf(os.Stderr, "could not play gif: %v\n", err)
	}
}

func getGIF(ctx context.Context, urlStr string) (*gif.GIF, error) {
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid http status code %d", resp.StatusCode)
	}
	return gif.DecodeAll(resp.Body)
}
