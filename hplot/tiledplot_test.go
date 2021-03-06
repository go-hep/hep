// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestTiledPlot(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleTiledPlot, t, "tiled_plot_histogram.png")
}

func TestTiledPlotAlign(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleTiledPlot_align, t, "tiled_plot_aligned_histogram.png")
}
