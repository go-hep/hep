// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"go-hep.org/x/hep/fastjet/internal/delaunay"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// Delaunay plots the given delaunay triangulation and saves it to the given path.
func Delaunay(path string, d *delaunay.Delaunay) error {
	ts := d.Triangles()
	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = "Delaunay Triangulation"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	for _, t := range ts {
		ax, ay := t.A.Coordinates()
		bx, by := t.B.Coordinates()
		cx, cy := t.C.Coordinates()
		pts := plotter.XYs{{X: ax, Y: ay}, {X: bx, Y: by}, {X: cx, Y: cy}}
		poly, err := plotter.NewPolygon(pts)
		if err != nil {
			return err
		}
		p.Add(poly)
	}
	return p.Save(10*vg.Centimeter, 10*vg.Centimeter, path)
}
