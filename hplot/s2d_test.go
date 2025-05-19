// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"fmt"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
)

func TestS2D(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleS2D, t, "s2d.png")
}

func TestScatter2DWithErrorBars(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleS2D_withErrorBars, t, "s2d_errbars.png")
}

func TestScatter2DWithBand(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleS2D_withBand, t, "s2d_band.png")
}

func TestScatter2DWithStepsKind(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleS2D_withStepsKind, t, "s2d_steps.png")
}

func TestScatter2DWithPreMidPostSteps(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleS2D_withPreMidPostSteps, t, "s2d_premidpost_steps.png")
}

func TestScatter2DWithStepsKindWithBand(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleS2D_withStepsKind_withBand, t, "s2d_steps_band.png")
}

func TestScatter2DSteps(t *testing.T) {
	for _, tc := range []struct {
		name string
		pts  []hbook.Point2D
		opts []hplot.Options
		want error
	}{
		{
			name: "histeps_no_xerr",
			pts: []hbook.Point2D{
				{X: 1, Y: 1, ErrY: hbook.Range{Min: 2, Max: 3}},
				{X: 2, Y: 2, ErrY: hbook.Range{Min: 5, Max: 2}},
				{X: 3, Y: 3, ErrY: hbook.Range{Min: 2, Max: 2}},
				{X: 4, Y: 4, ErrY: hbook.Range{Min: 1.2, Max: 2}},
			},
			opts: []hplot.Options{hplot.WithStepsKind(hplot.HiSteps)},
			want: fmt.Errorf("s2d with HiSteps needs XErr informations for all points"),
		},
		{
			name: "histeps_missing_some_xerr",
			pts: []hbook.Point2D{
				{X: 1, Y: 1, ErrY: hbook.Range{Min: 2, Max: 3}, ErrX: hbook.Range{Min: 1, Max: 2}},
				{X: 2, Y: 2, ErrY: hbook.Range{Min: 5, Max: 2}},
				{X: 3, Y: 3, ErrY: hbook.Range{Min: 2, Max: 2}, ErrX: hbook.Range{Min: 1, Max: 2}},
				{X: 4, Y: 4, ErrY: hbook.Range{Min: 1.2, Max: 2}},
			},
			opts: []hplot.Options{hplot.WithStepsKind(hplot.HiSteps)},
			want: fmt.Errorf("s2d with HiSteps needs XErr informations for all points"),
		},
		{
			name: "presteps_with_band", // TODO(sbinet)
			pts: []hbook.Point2D{
				{X: 1, Y: 1},
				{X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 4, Y: 4},
			},
			opts: []hplot.Options{hplot.WithStepsKind(hplot.PreSteps), hplot.WithBand(true)},
			want: fmt.Errorf("presteps not implemented"),
		},
		{
			name: "midsteps_with_band", // TODO(sbinet)
			pts: []hbook.Point2D{
				{X: 1, Y: 1},
				{X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 4, Y: 4},
			},
			opts: []hplot.Options{hplot.WithStepsKind(hplot.MidSteps), hplot.WithBand(true)},
			want: fmt.Errorf("midsteps not implemented"),
		},
		{
			name: "poststeps_with_band", // TODO(sbinet)
			pts: []hbook.Point2D{
				{X: 1, Y: 1},
				{X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 4, Y: 4},
			},
			opts: []hplot.Options{hplot.WithStepsKind(hplot.PostSteps), hplot.WithBand(true)},
			want: fmt.Errorf("poststeps not implemented"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				err := recover()
				switch {
				case err == nil && tc.want != nil:
					t.Fatalf("expected a panic")
				case err == nil && tc.want == nil:
					// ok.
				case err != nil && tc.want == nil:
					panic(err) // bubble up
				case err != nil && tc.want != nil:
					var got string
					switch err := err.(type) {
					case error:
						got = err.Error()
					case string:
						got = err
					default:
						panic(fmt.Errorf("invalid recover type %T", err))
					}
					if got, want := got, tc.want.Error(); got != want {
						t.Fatalf("invalid error:\ngot= %q\nwant=%q", got, want)
					}
				}
			}()

			s2d := hbook.NewS2D(tc.pts...)

			_ = hplot.NewS2D(s2d, tc.opts...)
		})
	}
}
