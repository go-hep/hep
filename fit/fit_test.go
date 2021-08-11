// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"os"
	"path"
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func checkPlot(ex func(), t *testing.T, filenames ...string) {
	checkPlotApprox(ex, t, 0, filenames...)
}

func checkPlotApprox(ex func(), t *testing.T, delta float64, filenames ...string) {
	t.Helper()
	cmpimg.CheckPlotApprox(ex, t, delta, filenames...)
	if t.Failed() {
		return
	}
	for _, fname := range filenames {
		_ = os.Remove(path.Join("testdata", fname))
	}
}
