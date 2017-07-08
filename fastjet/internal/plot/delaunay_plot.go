// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"image/color"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
	"go-hep.org/x/hep/fastjet/internal/delaunay"
)

// Plot plots the delaunay triangulation
func Plot(path string, d *delaunay.Delaunay) error {
	ts := d.Triangles()
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Delaunay Triangulation"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	for _, t := range ts {
		pts := plotter.XYs{{X: t.A.X, Y: t.A.Y}, {X: t.B.X, Y: t.B.Y}, {X: t.C.X, Y: t.C.Y}}
		poly, err := plotter.NewPolygon(pts)
		if err != nil {
			return err
		}
		p.Add(poly)
	}
	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, path)
	return err
}

// PlotVoronoiAndDelaunay plots the delaunay triangulation for all points and the voronoi diagram for the given points
// delaunay is black and voronoi is red
func PlotVoronoiAndDelaunay(path string, d *delaunay.Delaunay, points []*delaunay.Point) error {
	ts := d.Triangles()
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Delaunay And Voronoi"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	for _, t := range ts {
		pts := plotter.XYs{{X: t.A.X, Y: t.A.Y}, {X: t.B.X, Y: t.B.Y}, {X: t.C.X, Y: t.C.Y}}
		poly, err := plotter.NewPolygon(pts)
		if err != nil {
			return err
		}
		p.Add(poly)
	}
	v := delaunay.NewVoronoi(d)
	for _, pt := range points {
		_, centers := v.VoronoiCell(pt)
		pts := make(plotter.XYs, len(centers))
		for i := 0; i < len(pts); i++ {
			pts[i].X = centers[i].X
			pts[i].Y = centers[i].Y
		}
		poly, err := plotter.NewPolygon(pts)
		if err != nil {
			panic(err)
		}
		poly.LineStyle = draw.LineStyle{
			Color:    color.RGBA{R: 196, B: 128, A: 255},
			Width:    vg.Points(1),
			Dashes:   []vg.Length{},
			DashOffs: 0,
		}
		p.Add(poly)
	}
	err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, path)
	return err
}
