// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"
	"math"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
)

// ExampleFunction draws some functions.
func ExampleFunction() {
	quad := hplot.NewFunction(func(x float64) float64 { return x * x })
	quad.Color = color.RGBA{B: 255, A: 255}

	exp := hplot.NewFunction(func(x float64) float64 { return math.Pow(2, x) })
	exp.Dashes = []vg.Length{vg.Points(2), vg.Points(2)}
	exp.Width = vg.Points(2)
	exp.Color = color.RGBA{G: 255, A: 255}

	sin := hplot.NewFunction(func(x float64) float64 { return 10*math.Sin(x) + 50 })
	sin.Dashes = []vg.Length{vg.Points(4), vg.Points(5)}
	sin.Width = vg.Points(4)
	sin.Color = color.RGBA{R: 255, A: 255}

	p := hplot.New()
	p.Title.Text = "Functions"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	p.Add(quad, exp, sin)
	p.Legend.Add("x^2", quad)
	p.Legend.Add("2^x", exp)
	p.Legend.Add("10*sin(x)+50", sin)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 100

	err := p.Save(200, 200, "testdata/functions.png")
	if err != nil {
		log.Panic(err)
	}
}

// ExampleFunction_logY draws a function with a Log-Y axis.
func ExampleFunction_logY() {
	quad := hplot.NewFunction(func(x float64) float64 { return x * x })
	quad.Color = color.RGBA{B: 255, A: 255}

	fun := hplot.NewFunction(func(x float64) float64 {
		switch {
		case x < 6:
			return 20
		case 6 <= x && x < 7:
			return 0
		case 7 <= x && x < 7.5:
			return 30
		case 7.5 <= x && x < 9:
			return 0
		case 9 <= x:
			return 40
		}
		return 50
	})
	fun.LogY = true
	fun.Color = color.RGBA{R: 255, A: 255}

	p := hplot.New()
	p.Title.Text = "Functions - Log-Y scale"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	p.Y.Scale = plot.LogScale{}
	p.Y.Tick.Marker = plot.LogTicks{}

	p.Add(fun)
	p.Add(quad)
	p.Add(hplot.NewGrid())
	p.Legend.Add("x^2", quad)
	p.Legend.Add("fct", fun)
	p.Legend.ThumbnailWidth = 0.5 * vg.Inch

	p.X.Min = 5
	p.X.Max = 10
	p.Y.Min = 10
	p.Y.Max = 100

	err := p.Save(200, 200, "testdata/functions_logy.png")
	if err != nil {
		log.Panic(err)
	}
}
