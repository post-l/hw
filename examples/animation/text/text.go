package text

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

type Animation struct {
	ctx     *gg.Context
	pages   [][]string
	pageIdx int
	alpha   int
	state   state
	time    time.Time
}

type state int

const (
	fadeIn = iota
	show
	fadeOut
)

func NewAnimation(sz image.Point) (*Animation, error) {
	ctx := gg.NewContext(sz.X, sz.Y)
	ttf, err := getRobotoFont()
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(ttf)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 12,
	})
	ctx.SetFontFace(face)
	return &Animation{
		ctx: ctx,
	}, nil
}

func getRobotoFont() ([]byte, error) {
	resp, err := http.Get("https://github.com/google/fonts/blob/master/apache/roboto/Roboto-Medium.ttf?raw=true")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func (a *Animation) RandQuote() error {
	resp, err := http.Get("https://talaikis.com/api/quotes/random/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var res struct {
		Quote    string `json:"quote"`
		Author   string `json:"author"`
		Category string `json:"cat"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if res.Quote == "" {
		return errors.New("got empty quote")
	}
	text := fmt.Sprintf("%s -- %s (%s)", res.Quote, res.Author, res.Category)
	lines := a.ctx.WordWrap(text, float64(a.ctx.Width()))
	nbScreenLines := int(float64(a.ctx.Height()) / a.ctx.FontHeight())
	nbPages := len(lines) / nbScreenLines
	a.pages = make([][]string, nbPages)
	prev := 0
	for i := range a.pages {
		cur := prev + nbScreenLines
		if cur > len(lines) {
			cur = len(lines)
		}
		a.pages[i] = lines[prev:cur]
		prev = cur
	}
	a.pageIdx = 0
	a.alpha = 0
	a.state = fadeIn
	return nil
}

func (a *Animation) Image() image.Image {
	return a.ctx.Image()
}

func (a *Animation) Delay() time.Duration {
	return 100 * time.Millisecond
}

func (a *Animation) Next() error {
	switch a.state {
	case fadeIn:
		a.alpha += 15
		if a.alpha == 255 {
			a.state = show
			a.time = time.Now()
		}
	case show:
		if time.Since(a.time) >= 2*time.Second {
			a.state = fadeOut
		}
		return nil
	case fadeOut:
		a.alpha -= 15
		if a.alpha == 0 {
			a.pageIdx = (a.pageIdx + 1) % len(a.pages)
			a.state = fadeIn
		}
	}

	a.ctx.SetRGB255(0, 0, 0)
	a.ctx.Clear()
	a.ctx.SetRGBA255(0xff, 0xff, 0xff, a.alpha)
	y := float64(0)
	for _, line := range a.pages[a.pageIdx] {
		a.ctx.DrawStringAnchored(line, 0, y, 0, 1)
		y += a.ctx.FontHeight()
	}
	return nil
}
