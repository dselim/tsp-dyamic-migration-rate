package main

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var (
	palette = color.Palette{
		color.RGBA{R: 255, G: 0, B: 0},
		color.RGBA{R: 0, G: 255, B: 0},
		color.RGBA{R: 0, G: 0, B: 255},
		color.RGBA{R: 0, G: 255, B: 255},
		color.RGBA{R: 255, G: 0, B: 255},
		color.RGBA{R: 255, G: 255, B: 0},
		color.RGBA{R: 127, G: 205, B: 187},
		color.RGBA{R: 117, G: 107, B: 177},
		color.RGBA{R: 217, G: 95, B: 14},
		color.RGBA{R: 0, G: 0, B: 0},
	}
)

func plotLines(lines []plotter.XYs, legends []string, title, xlabel, ylabel, filepath string) {
	if len(lines) > len(palette) {
		panic("not enough colors !")
	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel

	p.Add(plotter.NewGrid())

	for i, line := range lines {
		li, po, err := plotter.NewLinePoints(line)
		if err != nil {
			panic(err)
		}
		li.LineStyle.Width = vg.Points(1)
		li.Color = palette[i]
		po.Color = palette[i]

		p.Add(li, po)
		p.Legend.Add(legends[i], li, po)
	}

	if err := p.Save(40*vg.Centimeter, 40*vg.Centimeter, filepath); err != nil {
		panic(err)
	}
}
