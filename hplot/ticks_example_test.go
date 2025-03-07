// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hplot_test

import (
	"image/color"
	"log"
	"time"

	"git.sr.ht/~sbinet/epok"
	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ExampleTicks() {
	tp := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 8})
	tp.Align = true

	for i := range tp.Tiles.Rows {
		for j := range tp.Tiles.Cols {
			p := tp.Plot(j, i)
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
				var (
					n    = 10
					xfmt = ""
				)
				switch i {
				case 4:
					n = 5
				case 5:
					n = 5
					xfmt = "%g"
				}
				p.X.Tick.Marker = hplot.Ticks{N: n, Format: xfmt}
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

func ExampleTicks_yearly() {
	cnv := epok.UTCUnixTimeConverter{}

	p := hplot.New()
	p.Title.Text = "Time series (yearly)"
	p.Y.Label.Text = "Goroutines"

	p.Y.Min = 0
	p.Y.Max = 4
	p.X.AutoRescale = true
	p.X.Tick.Marker = epok.Ticks{
		Ruler: epok.Rules{
			Major: epok.Rule{
				Freq:  epok.Yearly,
				Range: epok.Range{Step: 5},
			},
		},
		Format:    "2006-01-02\n15:04:05",
		Converter: cnv,
	}

	xysFrom := func(vs ...float64) plotter.XYs {
		o := make(plotter.XYs, len(vs))
		for i := range o {
			o[i].X = vs[i]
			o[i].Y = float64(i + 1)
		}
		return o
	}
	data := xysFrom(
		cnv.FromTime(parse("2010-02-03 01:02:03")),
		cnv.FromTime(parse("2011-03-04 11:22:33")),
		cnv.FromTime(parse("2015-02-03 04:05:06")),
		cnv.FromTime(parse("2020-02-03 07:08:09")),
	)

	line, pnts, err := hplot.NewLinePoints(data)
	if err != nil {
		log.Fatalf("could not create plotter: %+v", err)
	}

	line.Color = color.RGBA{B: 255, A: 255}
	pnts.Shape = draw.CircleGlyph{}
	pnts.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, pnts, hplot.NewGrid())

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, "testdata/timeseries_yearly.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

func ExampleTicks_monthly() {
	cnv := epok.UTCUnixTimeConverter{}

	p := hplot.New()
	p.Title.Text = "Time series (monthly)"
	p.Y.Label.Text = "Goroutines"

	p.Y.Min = 0
	p.Y.Max = 4
	p.X.AutoRescale = true
	p.X.Tick.Marker = epok.Ticks{
		Ruler: epok.Rules{
			Major: epok.Rule{
				Freq:  epok.Monthly,
				Range: epok.RangeFrom(1, 13, 2),
			},
		},
		Format:    "2006\nJan-02\n15:04:05",
		Converter: cnv,
	}

	xysFrom := func(vs ...float64) plotter.XYs {
		o := make(plotter.XYs, len(vs))
		for i := range o {
			o[i].X = vs[i]
			o[i].Y = float64(i + 1)
		}
		return o
	}
	data := xysFrom(
		cnv.FromTime(parse("2010-01-02 01:02:03")),
		cnv.FromTime(parse("2010-02-01 01:02:03")),
		cnv.FromTime(parse("2010-02-04 11:22:33")),
		cnv.FromTime(parse("2010-03-04 01:02:03")),
		cnv.FromTime(parse("2010-04-05 01:02:03")),
		cnv.FromTime(parse("2010-04-05 01:02:03")),
		cnv.FromTime(parse("2010-05-01 00:02:03")),
		cnv.FromTime(parse("2010-05-04 04:04:04")),
		cnv.FromTime(parse("2010-05-08 11:12:13")),
		cnv.FromTime(parse("2010-06-15 01:02:03")),
		cnv.FromTime(parse("2010-07-04 04:04:43")),
		cnv.FromTime(parse("2010-07-14 14:17:09")),
		cnv.FromTime(parse("2010-08-04 21:22:23")),
		cnv.FromTime(parse("2010-08-15 11:12:13")),
		cnv.FromTime(parse("2010-09-01 21:52:53")),
		cnv.FromTime(parse("2010-10-25 01:19:23")),
		cnv.FromTime(parse("2010-11-30 11:32:53")),
		cnv.FromTime(parse("2010-12-24 23:59:59")),
		cnv.FromTime(parse("2010-12-31 23:59:59")),
		cnv.FromTime(parse("2011-01-12 01:02:03")),
	)

	line, pnts, err := hplot.NewLinePoints(data)
	if err != nil {
		log.Fatalf("could not create plotter: %+v", err)
	}

	line.Color = color.RGBA{B: 255, A: 255}
	pnts.Shape = draw.CircleGlyph{}
	pnts.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, pnts, hplot.NewGrid())

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, "testdata/timeseries_monthly.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

func ExampleTicks_daily() {
	cnv := epok.UTCUnixTimeConverter{}

	p := hplot.New()
	p.Title.Text = "Time series (daily)"
	p.Y.Label.Text = "Goroutines"

	p.Y.Min = 0
	p.Y.Max = 4
	p.X.AutoRescale = true
	p.X.Tick.Marker = epok.Ticks{
		Ruler: epok.Rules{
			Major: epok.Rule{
				Freq:  epok.Daily,
				Range: epok.RangeFrom(1, 29, 14),
			},
			Minor: epok.Rule{
				Freq:  epok.Daily,
				Range: epok.RangeFrom(1, 32, 1),
			},
		},
		Format:    "2006\nJan-02\n15:04:05",
		Converter: cnv,
	}

	xysFrom := func(vs ...float64) plotter.XYs {
		o := make(plotter.XYs, len(vs))
		for i := range o {
			o[i].X = vs[i]
			o[i].Y = float64(i + 1)
		}
		return o
	}
	data := xysFrom(
		cnv.FromTime(parse("2020-01-01 01:02:03")),
		cnv.FromTime(parse("2020-01-02 02:03:04")),
		cnv.FromTime(parse("2020-01-12 03:04:05")),
		cnv.FromTime(parse("2020-01-22 04:05:06")),
		cnv.FromTime(parse("2020-02-03 05:06:07")),
		cnv.FromTime(parse("2020-02-13 06:07:08")),
		cnv.FromTime(parse("2020-02-23 07:08:09")),
		cnv.FromTime(parse("2020-03-01 00:00:00")),
	)

	line, pnts, err := hplot.NewLinePoints(data)
	if err != nil {
		log.Fatalf("could not create plotter: %+v", err)
	}

	line.Color = color.RGBA{B: 255, A: 255}
	pnts.Shape = draw.CircleGlyph{}
	pnts.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, pnts, hplot.NewGrid())

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, "testdata/timeseries_daily.png")
	if err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}
}

func parse(vs ...string) time.Time {
	format := "2006-01-02 15:04:05"
	if len(vs) > 1 {
		format = vs[1]
	}
	t, err := time.Parse(format, vs[0])
	if err != nil {
		panic(err)
	}
	return t
}
