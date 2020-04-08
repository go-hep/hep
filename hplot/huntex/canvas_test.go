// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package huntex_test

import (
	"testing"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

func TestCanvas(t *testing.T) {
	cmpimg.CheckPlot(func() {
		p := hplot.New()
		p.Title.Text = `Gaussian with \mu=1 and \sigma=0`
		p.X.Label.Text = `\alpha`
		p.Y.Label.Text = `\Delta`

		err := p.Save(10*vg.Centimeter, -1, "testdata/plot.png")
		if err != nil {
			t.Fatalf("could not save plot: %+v", err)
		}
	}, t, "plot.png")

}
