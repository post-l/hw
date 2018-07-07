package main

import (
	"context"
	"time"

	"github.com/post-l/hw/examples"
	"github.com/post-l/hw/examples/animation/circle"
	"github.com/post-l/hw/examples/animation/life"
	"github.com/post-l/hw/matrix/toolkit"
)

func main() {
	examples.Main(run)
}

func run(m toolkit.Matrix) error {
	tk := toolkit.New(m)
	sz := m.Bounds().Size()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		ca := circle.NewAnimation(sz)
		tk.PlayAnimation(ctx, ca)
		cancel()

		ctx, cancel = context.WithTimeout(context.Background(), 1*time.Minute)
		la := life.NewAnimation(sz)
		tk.PlayAnimation(ctx, la)
		cancel()
	}
}
