package main

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"gosim"
	"image/color"
	"math/rand"
)

func main() {
	rand := rand.New(rand.NewSource(0))

	b := gosim.NewBrownian(rand)

	const T = 1
	const N = 1000
	const output = "brownian.png"

	points := make(plotter.XYs, N)
	for i := 0; i < N; i++ {
		t := float64(i) * T / (N - 1)
		points[i].X = t
		points[i].Y = b.At(t)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Brownian Motion"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Value"

	l, err := plotter.NewLine(points)
	if err != nil {
		panic(err)
	}

	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Color = color.RGBA{B: 255, R: 255, A: 255}
	p.Add(l)

	if err := p.Save(6, 6, output); err != nil {
		panic(err)
	}

	fmt.Printf("Output written to %s\n", output)
}
