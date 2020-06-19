// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"testing"

	"gonum.org/v1/plot/cmpimg"
)

func TestBinnedErrBand(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleBinnedErrBand, t, "binnederrband.png")
}

func TestBinnedErrBandFromH1D(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(ExampleBinnedErrBand_fromH1D, t, "binnederrband_fromh1d.png")
}
