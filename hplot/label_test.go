// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func TestLabel(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleLabel, t, "label_plot.png")
}

func TestLabelPanic(t *testing.T) {
	for _, tc := range []struct {
		x, y float64
		txt  string
		opts []hplot.LabelOption
		err  string
	}{
		{
			x:    1.1,
			txt:  "invalid-x",
			opts: []hplot.LabelOption{hplot.WithLabelNormalized(true)},
			err:  "hplot: normalized label x-position is outside [0,1]: 1.1",
		},
		{
			y:    1.1,
			txt:  "invalid-y",
			opts: []hplot.LabelOption{hplot.WithLabelNormalized(true)},
			err:  "hplot: normalized label y-position is outside [0,1]: 1.1",
		},
		{
			x:    0.99,
			y:    0,
			txt:  "very long text in x",
			opts: []hplot.LabelOption{hplot.WithLabelNormalized(true)},
			err:  "hplot: label (0.99, 0) falls outside data canvas",
		},
		{
			x:    0,
			y:    0.99,
			txt:  "very tall text in y\n1\n2\n",
			opts: []hplot.LabelOption{hplot.WithLabelNormalized(true)},
			err:  "hplot: label (0, 0.99) falls outside data canvas",
		},
	} {
		t.Run(tc.txt, func(t *testing.T) {
			defer func() {
				e := recover()
				if e == nil {
					t.Fatalf("expected a panic %q", tc.err)
				}
				if got, want := e.(error).Error(), tc.err; got != want {
					t.Fatalf("invalid panic message\ngot= %q\nwant=%q",
						got, want,
					)
				}
			}()

			lbl := hplot.NewLabel(tc.x, tc.y, tc.txt, tc.opts...)

			p := hplot.New()
			p.X.Min = -10
			p.X.Max = +10
			p.Y.Min = -10
			p.Y.Max = +10
			p.Add(lbl)

			const (
				sz = 10 * vg.Centimeter
			)
			dc, err := draw.NewFormattedCanvas(sz, sz, "png")
			if err != nil {
				t.Fatalf("could not create draw canvas: %+v", err)
			}

			p.Draw(draw.New(dc))
		})
	}
}

func TestLabelWithLog(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {

		// Creating a new plot
		p := hplot.New()
		p.Title.Text = "Plot labels"
		p.X.Min = 1
		p.X.Max = 110
		p.Y.Min = 1
		p.Y.Max = 110

		p.X.Scale = plot.LogScale{}
		p.Y.Scale = plot.LogScale{}

		p.Add(hplot.NewLabel(
			0.5, 0.5,
			"(0.5,0.5)\nMy Label",
			hplot.WithLabelNormalized(true),
		))

		p.Add(hplot.NewLabel(
			0.95, 0.95,
			"(0.95,0.95)\nAuto-adjust",
			hplot.WithLabelNormalized(true),
			hplot.WithLabelAutoAdjust(true),
		))

		p.Add(plotter.NewGlyphBoxes())
		p.Add(hplot.NewGrid())

		// Save the plot to a PNG file.
		err := p.Save(15*vg.Centimeter, -1, "testdata/label_log_plot.png")
		if err != nil {
			t.Fatalf("error saving plot: %v\n", err)
		}

	}, t, "label_log_plot.png")
}
