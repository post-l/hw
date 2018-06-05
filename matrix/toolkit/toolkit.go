package toolkit

import (
	"context"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"io"
	"time"

	"github.com/disintegration/imaging"
)

type Animation interface {
	Delay() time.Duration
	Image() image.Image
	Next() error
}

type Frame struct {
	Image image.Image
	Delay time.Duration
}

// Matrix is an interface that represent any RGB matrix, very useful for testing.
type Matrix interface {
	draw.Image
	Render()
}

// ToolKit is a convinient set of function to operate with a led of Matrix.
type ToolKit struct {
	m Matrix
}

// New returns a new ToolKit wrapping the given Matrix.
func New(m Matrix) *ToolKit {
	return &ToolKit{
		m: m,
	}
}

// DrawImage draws the given image.
func (tk *ToolKit) DrawImage(img image.Image) {
	draw.Draw(tk.m, tk.m.Bounds(), img, image.ZP, draw.Src)
	tk.m.Render()
}

// PlayAnimation play the image during the delay returned by Next, until an err
// is returned, if io.EOF is returned, PlayAnimation finish without an error.
func (tk *ToolKit) PlayAnimation(ctx context.Context, a Animation) error {
	delay := a.Delay()
	t := time.Now()
	var dt time.Duration
	for {
		steps := int(dt / delay)
		for step := 0; step < steps; step++ {
			if err := a.Next(); err != nil {
				if err == io.EOF {
					return nil
				}
				return err
			}
		}
		img := a.Image()
		tk.DrawImage(img)
		dt = dt % delay
		d := delay - dt
		select {
		case <-ctx.Done():
			return ctx.Err()
		case now := <-time.After(d):
			dt += now.Sub(t)
			t = now
		}
	}
}

// PlayFrames draws a sequence of frames.
func (tk *ToolKit) PlayFrames(ctx context.Context, frames []Frame, loopCount int) error {
	l := len(frames)
	i := 0
	loop := 0
	for {
		f := frames[i]
		tk.DrawImage(f.Image)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(f.Delay):
		}
		i++
		if i >= l {
			if loopCount > 0 {
				loop++
				if loop >= loopCount {
					return nil
				}
			}
			i = 0
		}
	}
}

// PlayGIF draws a gif. It use the contained images and
// delays and loops over it.
func (tk *ToolKit) PlayGIF(ctx context.Context, gif *gif.GIF) error {
	if len(gif.Image) == 0 {
		return errors.New("no image in the gif")
	}
	frames := make([]Frame, len(gif.Image))
	sz := tk.m.Bounds()
	w, h := sz.Dx(), sz.Dy()
	buf := image.NewRGBA(gif.Image[0].Bounds())
	for i, img := range gif.Image {
		b := img.Bounds()
		draw.Draw(buf, b, img, b.Min, draw.Over)
		frames[i] = Frame{
			Image: imaging.Resize(buf, w, h, imaging.Lanczos),
			Delay: time.Duration(gif.Delay[i]*10) * time.Millisecond,
		}
	}
	return tk.PlayFrames(ctx, frames, gif.LoopCount)
}
