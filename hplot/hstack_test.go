// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"fmt"
	"image/color"
	"log"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

func TestHStack(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleHStack, t, "hstack.png")
}

func TestHStackPanic(t *testing.T) {
	for _, tc := range []struct {
		fct    func() []*hplot.H1D
		panics error
	}{
		{
			fct: func() []*hplot.H1D {
				return nil
			},
			panics: fmt.Errorf("hplot: not enough histograms to make a stack"),
		},
		{
			fct: func() []*hplot.H1D {
				return make([]*hplot.H1D, 0)
			},
			panics: fmt.Errorf("hplot: not enough histograms to make a stack"),
		},
		{
			fct: func() []*hplot.H1D {
				return []*hplot.H1D{
					hplot.NewH1D(hbook.NewH1D(10, 0, 10)),
					hplot.NewH1D(hbook.NewH1D(11, 0, 10)),
				}
			},
			panics: fmt.Errorf("hplot: bins length mismatch"),
		},
		{
			fct: func() []*hplot.H1D {
				return []*hplot.H1D{
					hplot.NewH1D(hbook.NewH1D(10, 0, 10)),
					hplot.NewH1D(hbook.NewH1D(10, 0, 11)),
				}
			},
			panics: fmt.Errorf("hplot: bin range mismatch"),
		},
	} {
		t.Run("", func(t *testing.T) {
			defer func() {
				err := recover()
				if err == nil {
					t.Fatalf("expected a panic")
				}
				switch err := err.(type) {
				case string:
					if got, want := err, tc.panics.Error(); got != want {
						t.Fatalf(
							"invalid panic message:\ngot= %v\nwant=%v",
							got, want,
						)
					}
				case error:
					if got, want := err.Error(), tc.panics.Error(); got != want {
						t.Fatalf(
							"invalid panic message:\ngot= %v\nwant=%v",
							got, want,
						)
					}
				}
			}()
			_ = hplot.NewHStack(tc.fct())
		})
	}
}

func TestHStackCornerBins(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		h1 := hbook.NewH1D(10, 0, 10)
		h2 := hbook.NewH1D(10, 0, 10)
		h3 := hbook.NewH1D(10, 0, 10)

		h1.FillN(
			[]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		)
		h2.FillN(
			[]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]float64{2, 2, 2, 2, 2, 2, 2, 2, 2, 2},
		)
		h3.FillN(
			[]float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			[]float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1},
		)

		colors := []color.Color{
			color.NRGBA{122, 195, 106, 150},
			color.NRGBA{90, 155, 212, 150},
			color.NRGBA{250, 167, 91, 150},
		}

		hh1 := hplot.NewH1D(h1)
		hh1.FillColor = colors[0]
		hh1.LineStyle.Color = color.Black

		hh2 := hplot.NewH1D(h2)
		hh2.FillColor = colors[1]
		hh2.LineStyle.Width = 0

		hh3 := hplot.NewH1D(h3)
		hh3.FillColor = colors[2]
		hh3.LineStyle.Color = color.Black

		hs := []*hplot.H1D{hh1, hh2, hh3}

		p := hplot.New()
		p.X.Label.Text = "X"
		p.Y.Label.Text = "Y"
		p.Y.Min = -0.5
		p.Y.Max = 15.5
		hstack := hplot.NewHStack(hs, hplot.WithLogY(false))
		p.Add(hstack, hplot.NewGrid())
		p.Legend.Add("h1", hs[0])
		p.Legend.Add("h2", hs[1])
		p.Legend.Add("h3", hs[2])
		p.Legend.Top = true
		p.Legend.Left = true

		err := p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/hstack_corner_bins.png")
		if err != nil {
			log.Fatalf("error: %+v", err)
		}

	}, t, "hstack_corner_bins.png")
}
