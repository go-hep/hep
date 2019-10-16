// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestH1D(t *testing.T) {
	cmpimg.CheckPlot(ExampleH1D, t, "h1d_plot.png")
}

func TestH1DtoPDF(t *testing.T) {
	cmpimg.CheckPlot(ExampleH1D_toPDF, t, "h1d_plot.pdf")
}

func TestH1DLogScale(t *testing.T) {
	cmpimg.CheckPlot(ExampleH1D_logScaleY, t, "h1d_logy.png")
}
