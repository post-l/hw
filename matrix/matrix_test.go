package matrix_test

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"testing"
	"time"

	"github.com/disintegration/imaging"

	"github.com/post-l/hw/board/tinkerboard"
	"github.com/post-l/hw/matrix"
)

func TestMatrix(t *testing.T) {
	b, err := tinkerboard.New()
	if err != nil {
		t.Fatal("board:", err)
	}
	m := matrix.New(b, &matrix.DefaultHardwareConfig)
	defer m.Close()
	f, err := os.Open("testdata/gopher.jpg")
	if err != nil {
		t.Fatal("open:", err)
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatal("jpeg:", err)
	}
	w, h := m.Geometry()
	img = imaging.Resize(img, w, h, imaging.Lanczos)
	draw.Draw(m, m.Bounds(), img, image.ZP, draw.Src)
	m.Render()
	time.Sleep(10 * time.Second)
}
