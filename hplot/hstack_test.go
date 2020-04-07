// Copyright ©2020 The go-hep Authors. All rights reserved.
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

func TestHStack(t *testing.T) {
	cmpimg.CheckPlot(ExampleHStack, t, "hstack.png")
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
