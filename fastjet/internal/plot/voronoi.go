// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"math"

	"go-hep.org/x/hep/fastjet/internal/delaunay"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Voronoi plots the Voronoi diagram of the points and saves it to the given path.
//
// It will compute the delaunay triangulation first.
func Voronoi(path string, points []*delaunay.Point) error {
	// First insert all points into a Delaunay triangulation.
	d := delaunay.HierarchicalDelaunay()
	for _, p := range points {
		d.Insert(p)
	}
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Voronoi Diagram"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	// pts are Scatters of all the points inserted.
	pts := make(plotter.XYs, len(points))
	// find the min and max coordinates of the inserted points.
	minX, minY, maxX, maxY := math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)
	for i, pt := range points {
		pts[i].X, pts[i].Y = pt.Coordinates()
		if pts[i].X < minX {
			minX = pts[i].X
		}
		if pts[i].Y < minY {
			minY = pts[i].Y
		}
		if pts[i].X > maxX {
			maxX = pts[i].X
		}
		if pts[i].Y > maxY {
			maxY = pts[i].Y
		}
		voronoi, _ := d.VoronoiCell(pt)
		vs := make(plotter.XYs, len(voronoi))
		for j := range vs {
			vs[j].X, vs[j].Y = voronoi[j].Coordinates()
		}
		poly, err := plotter.NewPolygon(vs)
		if err != nil {
			return err
		}
		p.Add(poly)
	}
	s, err := plotter.NewScatter(pts)
	if err != nil {
		return err
	}
	p.Add(s)
	// cut plot off to only show the points plus a little extra space.
	xpad := (maxX - minX) * .05
	ypad := (maxY - minY) * .05
	p.X.Max = maxX + xpad
	p.X.Min = minX - xpad
	p.Y.Max = maxY + ypad
	p.Y.Min = minY - ypad
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, path)
}
