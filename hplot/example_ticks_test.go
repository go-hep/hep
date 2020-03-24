// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"log"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ExampleTicks() {
	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 8})
	tp.Align = true

	for i := 0; i < tp.Tiles.Rows; i++ {
		for j := 0; j < tp.Tiles.Cols; j++ {
			p := tp.Plot(i, j)
			switch i {
			case 0:
				p.X.Min = 0
				p.X.Max = 1
				switch j {
				case 0:
					p.Title.Text = "hplot.Ticks"
				default:
					p.Title.Text = "plot.Ticks"
				}
			case 1:
				p.X.Min = 0
				p.X.Max = 10
			case 2:
				p.X.Min = 0
				p.X.Max = 100
			case 3:
				p.X.Min = 0
				p.X.Max = 1000
			case 4:
				p.X.Min = 0
				p.X.Max = 10000
			case 5:
				p.X.Min = 0
				p.X.Max = 10000
			case 6:
				p.X.Min = 0
				p.X.Max = 1.2
			case 7:
				p.X.Min = 0
				p.X.Max = 120
			}
			if j == 0 {
				n := 20
				switch i {
				case 4:
					n = 10
				case 5:
					n = 5
				}
				p.X.Tick.Marker = hplot.Ticks{N: n}
			}
			p.Add(hplot.NewGrid())
		}
	}

	const sz = 20 * vg.Centimeter
	err := tp.Save(sz, sz, "testdata/ticks.png")
	if err != nil {
		log.Fatalf("error: %+v\n", err)
	}
}
