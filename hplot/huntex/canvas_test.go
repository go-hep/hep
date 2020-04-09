// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package huntex_test

import (
	"image/color"
	"math"
	"testing"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

func TestCanvas(t *testing.T) {
	cmpimg.CheckPlot(func() {
		p := hplot.New()
		p.Title.Text = `Gaussian with $\mu=1$ and $\sigma=0$`
		p.X.Label.Text = `$\alpha$`
		p.Y.Label.Text = `$\Delta$`

		fct := hplot.NewFunction(math.Cos)
		fct.LineStyle.Color = color.RGBA{R: 255, A: 255}
		p.Legend.Add(`$\beta$`, fct)

		var err error

		err = p.Save(10*vg.Centimeter, -1, "testdata/plot.png")
		if err != nil {
			t.Fatalf("could not save plot: %+v", err)
		}
		err = p.Save(10*vg.Centimeter, -1, "testdata/plot.tex")
		if err != nil {
			t.Fatalf("could not save plot: %+v", err)
		}
	}, t, "plot.png", "plot.tex")

}
