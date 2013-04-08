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
	g := gosim.NewGeometricBrownian(rand, 1, 1)

	const T = 1
	const N = 1000
	const output = "brownian.png"

	std_points := make(plotter.XYs, N)
	geo_points := make(plotter.XYs, N)
	for i := 0; i < N; i++ {
		t := float64(i) * T / (N - 1)
		std_points[i].X = t
		std_points[i].Y = b.At(t)
		geo_points[i].X = t
		geo_points[i].Y = g.At(t)
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Brownian Motion"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Value"

	std, err := plotter.NewLine(std_points)
	if err != nil {
		panic(err)
	}

	geo, err := plotter.NewLine(geo_points)
	if err != nil {
		panic(err)
	}

	std.LineStyle.Width = vg.Points(1)
	std.LineStyle.Color = color.RGBA{B: 255, R: 255, A: 255}
	geo.LineStyle.Width = vg.Points(1)
	geo.LineStyle.Color = color.RGBA{B: 255, G: 255, A: 255}
	p.Add(std, geo)

	p.Legend.Add("Standard Brownian Motion", std)
	p.Legend.Add("Geometric Brownian Motion, drift 1", geo)
	p.Legend.Top = true
	p.Legend.Left = true

	if err := p.Save(6, 6, output); err != nil {
		panic(err)
	}

	fmt.Printf("Output written to %s\n", output)
}
