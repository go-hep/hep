// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

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
