package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io"
	"math/rand"
	"time"

	"github.com/nsf/termbox-go"
	"github.com/post-l/hw/examples"
	"github.com/post-l/hw/matrix/toolkit"
)

func main() {
	examples.Main(run)
}

func run(m toolkit.Matrix) error {
	if err := termbox.Init(); err != nil {
		return err
	}
	go updateInputChan()

	tk := toolkit.New(m)
	sz := m.Bounds().Size()
	a := NewAnimation(sz)
	ctx := context.Background()
	tk.PlayAnimation(ctx, a)

	termbox.Close()
	fmt.Println("score:", len(a.snake.body))
	return nil
}

var input = make(chan termbox.Event, 1)

func updateInputChan() {
	var nextEv termbox.Event
	for {
		ev := termbox.PollEvent()
		switch ev.Type {
		case termbox.EventKey:
			select {
			case input <- ev:
			default:
				nextEv = ev
			}
		case termbox.EventInterrupt:
			if nextEv.Type == termbox.EventKey {
				input <- nextEv
				nextEv = ev
			}
		}
	}
}

type Animation struct {
	screen draw.Image
	sz     image.Point
	snake  *Snake
	foods  []image.Point
}

func NewAnimation(sz image.Point) *Animation {
	a := &Animation{
		screen: image.NewRGBA(image.Rect(0, 0, sz.X, sz.Y)),
		sz:     sz,
		foods:  make([]image.Point, 10),
	}
	a.snake = NewSnake(a)
	for i := range a.foods {
		a.foods[i] = randPoint(sz)
	}
	return a
}

func (a *Animation) Image() image.Image {
	draw.Draw(a.screen, a.screen.Bounds(), image.Black, image.ZP, draw.Src)
	for _, p := range a.foods {
		a.screen.Set(p.X, p.Y, color.RGBA{G: 255, A: 255})
	}
	a.snake.Draw()
	return a.screen
}

func (a *Animation) Delay() time.Duration {
	return 150 * time.Millisecond
}

func (a *Animation) Next() error {
	termbox.Interrupt()
	select {
	case ev := <-input:
		switch ev.Key {
		case termbox.KeyArrowUp:
			if a.snake.dir != Down {
				a.snake.dir = Up
			}
		case termbox.KeyArrowDown:
			if a.snake.dir != Up {
				a.snake.dir = Down
			}
		case termbox.KeyArrowLeft:
			if a.snake.dir != Right {
				a.snake.dir = Left
			}
		case termbox.KeyArrowRight:
			if a.snake.dir != Left {
				a.snake.dir = Right
			}
		case termbox.KeyEnter:
			return io.EOF
		}
	default:
	}
	return a.snake.Next()
}

type Direction int

const (
	Up = Direction(iota)
	Down
	Left
	Right
)

type Snake struct {
	a     *Animation
	dir   Direction
	body  []image.Point
	color color.RGBA
}

func NewSnake(a *Animation) *Snake {
	return &Snake{
		a:    a,
		dir:  Right,
		body: []image.Point{randPoint(a.screen.Bounds().Size())},
		color: color.RGBA{
			R: 255,
			G: 0,
			B: 0,
			A: 255,
		},
	}
}

func (s *Snake) Next() error {
	p := s.body[0]
	switch s.dir {
	case Up:
		p.Y--
		if p.Y < 0 {
			p.Y = s.a.sz.Y - 1
		}
	case Down:
		p.Y++
		if p.Y >= s.a.sz.Y {
			p.Y = 0
		}
	case Left:
		p.X--
		if p.X < 0 {
			p.X = s.a.sz.X - 1
		}
	case Right:
		p.X++
		if p.X >= s.a.sz.X {
			p.X = 0
		}
	}
	tail := s.body[len(s.body)-1]
	for i := len(s.body) - 1; i >= 1; i-- {
		s.body[i] = s.body[i-1]
		if p.Eq(s.body[i]) {
			return io.EOF
		}
	}
	s.body[0] = p
	for i, food := range s.a.foods {
		if p.Eq(food) {
			s.body = append(s.body, tail)
			s.a.foods[i] = randPoint(s.a.sz)
			fmt.Println("score:", len(s.body))
		}
	}
	return nil
}

func (s *Snake) Draw() {
	for _, p := range s.body {
		s.a.screen.Set(p.X, p.Y, s.color)
	}
}

func randPoint(sz image.Point) image.Point {
	return image.Point{
		X: rand.Intn(sz.X),
		Y: rand.Intn(sz.Y),
	}
}
