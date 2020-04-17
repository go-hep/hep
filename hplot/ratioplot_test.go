// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestRatioPlot(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleRatioPlot, t, "diff_plot.png")
}
