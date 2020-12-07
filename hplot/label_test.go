// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestLabel(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(exampleLabel, t, "label_plot.png")
}
