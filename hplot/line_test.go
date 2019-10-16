// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestVLine(t *testing.T) {
	cmpimg.CheckPlot(ExampleVLine, t, "vline.png")
}

func TestHLine(t *testing.T) {
	cmpimg.CheckPlot(ExampleHLine, t, "hline.png")
}
