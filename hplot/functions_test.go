// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestFunction(t *testing.T) {
	cmpimg.CheckPlot(ExampleFunction, t, "functions.png")
}

func TestFunctionLogY(t *testing.T) {
	cmpimg.CheckPlot(ExampleFunction_logY, t, "functions_logy.png")
}
