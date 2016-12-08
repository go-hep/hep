// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"math"
	"sort"
)

// Scatter2D is a collection of 2-dim data points with errors.
type Scatter2D struct {
	pts []Point2D
	ann Annotation
}

// NewScatter2D creates a new 2-dim scatter with pts as an optional
// initial set of data points.
//
// If n <= 0, the initial capacity is 0.
func NewScatter2D(pts ...Point2D) *Scatter2D {
	s := &Scatter2D{
		pts: make([]Point2D, len(pts)),
		ann: make(Annotation),
	}
	copy(s.pts, pts)
	return s
}

// NewScatter2DFrom creates a new 2-dim scatter with x,y data slices.
//
// It panics if the lengths of the 2 slices don't match.
func NewScatter2DFrom(x, y []float64) *Scatter2D {
	if len(x) != len(y) {
		panic("hbook: len differ")
	}

	s := &Scatter2D{
		pts: make([]Point2D, len(x)),
		ann: make(Annotation),
	}
	for i := range s.pts {
		pt := &s.pts[i]
		pt.X = x[i]
		pt.Y = y[i]
	}
	return s
}

// Annotation returns the annotations attached to the
// scatter. (e.g. name, title, ...)
func (s *Scatter2D) Annotation() Annotation {
	return s.ann
}

// Name returns the name of this scatter
func (s *Scatter2D) Name() string {
	v, ok := s.ann["name"]
	if !ok {
		return ""
	}
	n, ok := v.(string)
	if !ok {
		return ""
	}
	return n
}

// Rank returns the number of dimensions of this scatter.
func (*Scatter2D) Rank() int {
	return 2
}

// Entries returns the number of entries of this histogram.
func (s *Scatter2D) Entries() int64 {
	return int64(len(s.pts))
}

// Fill adds new points to the scatter.
func (s *Scatter2D) Fill(pts ...Point2D) {
	if len(pts) == 0 {
		return
	}

	i := len(s.pts)
	s.pts = append(s.pts, make([]Point2D, len(pts))...)
	copy(s.pts[i:], pts)
}

// Sort sorts the data points by x,y and x-err,y-err.
func (s *Scatter2D) Sort() {
	sort.Sort(points2D(s.pts))
}

// Points returns the points of the scatter.
//
// Users may not modify the returned slice.
// Users may not rely on the stability of the indices as the slice of points
// may be re-sorted at any point in time.
func (s *Scatter2D) Points() []Point2D {
	return s.pts
}

// Point returns the point at index i.
//
// Point panics if i is out of bounds.
func (s *Scatter2D) Point(i int) Point2D {
	return s.pts[i]
}

// ScaleX rescales the X values by a factor f.
func (s *Scatter2D) ScaleX(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleX(f)
	}
}

// ScaleY rescales the Y values by a factor f.
func (s *Scatter2D) ScaleY(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleY(f)
	}
}

// ScaleXY rescales the X and Y values by a factor f.
func (s *Scatter2D) ScaleXY(f float64) {
	for i := range s.pts {
		p := &s.pts[i]
		p.ScaleX(f)
		p.ScaleY(f)
	}
}

// Len returns the number of points in the scatter.
//
// Len implements the gonum/plot/plotter.XYer interface.
func (s *Scatter2D) Len() int {
	return len(s.pts)
}

// XY returns the x, y pair at index i.
//
// XY panics if i is out of bounds.
// XY implements the gonum/plot/plotter.XYer interface.
func (s *Scatter2D) XY(i int) (x, y float64) {
	pt := s.pts[i]
	x = pt.X
	y = pt.Y
	return
}

// XError returns the two error values for X data.
//
// XError implements the gonum/plot/plotter.XErrorer interface.
func (s *Scatter2D) XError(i int) (float64, float64) {
	pt := s.pts[i]
	return pt.ErrX.Min, pt.ErrX.Max
}

// YError returns the two error values for Y data.
//
// YError implements the gonum/plot/plotter.YErrorer interface.
func (s *Scatter2D) YError(i int) (float64, float64) {
	pt := s.pts[i]
	return pt.ErrY.Min, pt.ErrY.Max
}

// DataRange returns the minimum and maximum
// x and y values, implementing the gonum/plot.DataRanger
// interface.
func (s *Scatter2D) DataRange() (xmin, xmax, ymin, ymax float64) {
	xmin = math.Inf(+1)
	ymin = math.Inf(+1)
	xmax = math.Inf(-1)
	ymax = math.Inf(-1)
	for _, p := range s.pts {
		xmin = math.Min(p.XMin(), xmin)
		xmax = math.Max(p.XMax(), xmax)
		ymin = math.Min(p.YMin(), ymin)
		ymax = math.Max(p.YMax(), ymax)
	}
	return
}
