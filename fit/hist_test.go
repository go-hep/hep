// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"testing"
)

func TestH1D(t *testing.T) {
	checkPlot(ExampleH1D_gaussian, t, "h1d-gauss-plot.png")
}
