// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

// Point2D is a position in a 2-dim space
type Point2D struct {
	X    float64 // x-position
	Y    float64 // y-position
	ErrX Range   // error on x-position
	ErrY Range   // error on y-position
}

// XMin returns the X value minus negative X-error
func (p Point2D) XMin() float64 {
	return p.X - p.ErrX.Min
}

// XMax returns the X value plus positive X-error
func (p Point2D) XMax() float64 {
	return p.X + p.ErrX.Max
}

// YMin returns the Y value minus negative Y-error
func (p Point2D) YMin() float64 {
	return p.Y - p.ErrY.Min
}

// YMax returns the Y value plus positive Y-error
func (p Point2D) YMax() float64 {
	return p.Y + p.ErrY.Max
}

// ScaleX rescales the X value by a factor f.
func (p *Point2D) ScaleX(f float64) {
	p.X *= f
	p.ErrX.Min *= f
	p.ErrX.Max *= f
}

// ScaleY rescales the Y value by a factor f.
func (p *Point2D) ScaleY(f float64) {
	p.Y *= f
	p.ErrY.Min *= f
	p.ErrY.Max *= f
}

// ScaleXY rescales the X and Y values by a factor f.
func (p *Point2D) ScaleXY(f float64) {
	p.ScaleX(f)
	p.ScaleY(f)
}

// points2D implements sort.Interface
type points2D []Point2D

func (p points2D) Len() int { return len(p) }
func (p points2D) Less(i, j int) bool {
	pi := p[i]
	pj := p[j]
	if pi.X != pj.X {
		return pi.X < pj.X
	}
	if pi.ErrX.Min != pj.ErrX.Min {
		return pi.ErrX.Min < pj.ErrX.Min
	}
	if pi.ErrX.Max != pj.ErrX.Max {
		return pi.ErrX.Max < pj.ErrX.Max
	}
	if pi.Y != pj.Y {
		return pi.Y < pj.Y
	}
	if pi.ErrY.Min != pj.ErrY.Min {
		return pi.ErrY.Min < pj.ErrY.Min
	}
	if pi.ErrY.Max != pj.ErrY.Max {
		return pi.ErrY.Max < pj.ErrY.Max
	}
	return false
}
func (p points2D) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

// point1D is a position in a 1-dim space
type point1D struct {
	x  float64    // x-position
	ex [2]float64 // error on x-position
}
