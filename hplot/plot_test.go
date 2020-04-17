// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"os"
	"testing"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/cmpimg"
	"gonum.org/v1/plot/vg"
)

func TestPlotWriterTo(t *testing.T) {
	checkPlot(cmpimg.CheckPlot)(func() {
		p := hplot.New()
		p.Title.Text = "title"
		p.X.Label.Text = "x"
		p.Y.Label.Text = "y"

		c, err := p.WriterTo(5*vg.Centimeter, 5*vg.Centimeter, "png")
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Create("testdata/plot_writerto.png")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		_, err = c.WriteTo(f)
		if err != nil {
			t.Fatal(err)
		}

		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}, t, "plot_writerto.png")
}
