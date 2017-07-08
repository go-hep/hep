// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math"
)

// Voronoi holds the border information
type Voronoi struct {
	maxX, minX, maxY, minY float64
}

// FIXME can't do any delaunay operation after calling NewVoronoi
func NewVoronoi(d *Delaunay) *Voronoi {
	border := make(triangles, len(d.root.A.adjacentTriangles)+len(d.root.B.adjacentTriangles)+len(d.root.C.adjacentTriangles))
	n := copy(border, d.root.A.adjacentTriangles)
	n += copy(border[n:], d.root.B.adjacentTriangles)
	copy(border[n:], d.root.C.adjacentTriangles)
	for _, t := range d.triangles {
		if len(t.children) != 0 {
			continue
		}
		for j, tri := range border {
			if tri.Equals(t) {
				border = append(border[:j], border[j+1:]...)
				t.A.isOutside = true
				if !t.A.Equals(d.root.A) && !t.A.Equals(d.root.B) && !t.A.Equals(d.root.C) {
					t.A.adjacentTriangles = t.A.adjacentTriangles.remove(t)
				}
				t.B.isOutside = true
				if !t.B.Equals(d.root.A) && !t.B.Equals(d.root.B) && !t.B.Equals(d.root.C) {
					t.B.adjacentTriangles = t.B.adjacentTriangles.remove(t)
				}
				t.C.isOutside = true
				if !t.C.Equals(d.root.A) && !t.C.Equals(d.root.B) && !t.C.Equals(d.root.C) {
					t.C.adjacentTriangles = t.C.adjacentTriangles.remove(t)
				}
				break
			}
		}
	}
	return &Voronoi{maxX: d.maxX, minX: d.minX, maxY: d.maxY, minY: d.minY}
}

// VoronoiCell finds the voronoi area of a point
// returns the area and the points associated with it
func (v *Voronoi) VoronoiCell(p *Point) (area float64, centers []*Point) {
	// find all points that form the voronoi cell for that point
	// in clockwise order
	centers = make([]*Point, len(p.adjacentTriangles), len(p.adjacentTriangles)+3)
	var t, first *Triangle
	if p.isOutside {
		// find first triangle on the outside
		for _, t1 := range p.adjacentTriangles {
			var pt *Point
			// assign pt to the Point clockwise of p
			switch {
			case p.Equals(t1.A):
				pt = t1.C
			case p.Equals(t1.B):
				pt = t1.A
			case p.Equals(t1.C):
				pt = t1.B
			default:
				panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", p, t))
			}
			border := true
			for _, t2 := range p.adjacentTriangles {
				if t1.Equals(t2) {
					continue
				}
				if pt.Equals(t2.A) || pt.Equals(t2.B) || pt.Equals(t2.C) {
					border = false
				}
			}
			if border {
				t = t1
				break
			}
		}
		first = t
		if t == nil {
			panic(fmt.Errorf("voronoi: invalid point P%v", p))
		}
	} else {
		t = p.adjacentTriangles[0]
	}
	x, y := t.centerOfCircumcircle()
	c := NewPoint(x, y, -1)
	centers[0] = c
	// if p has only one adjacent triangle last is first
	last := first
	for i := 1; i < len(p.adjacentTriangles); i++ {
		t = findClockTri(p, t)
		x, y = t.centerOfCircumcircle()
		c = NewPoint(x, y, -1)
		centers[i] = c
		if i == len(p.adjacentTriangles)-1 {
			last = t
		}
	}
	if !p.isOutside {
		area = polyArea(centers)
		return area, centers
	}
	// FIXME problem when circumcenter is outside border
	// extend polygon to border
	var pt, pt2 *Point
	switch {
	case p.Equals(last.A):
		pt = last.B
		pt2 = last.C
	case p.Equals(last.B):
		pt = last.C
		pt2 = last.A
	case p.Equals(last.C):
		pt = last.A
		pt2 = last.B
	default:
		panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", p, t))
	}
	var x1, y1, x2, y2 float64
	// check if circumcenter is outside the polygon
	if aboveLine(p, pt, centers[len(centers)-1]) == aboveLine(p, pt, pt2) {
		x1 = centers[len(centers)-1].X
		y1 = centers[len(centers)-1].Y
		x2 = (p.X + pt.X) / 2
		y2 = (p.Y + pt.Y) / 2
	} else {
		x1 = (p.X + pt.X) / 2
		y1 = (p.Y + pt.Y) / 2
		x2 = centers[len(centers)-1].X
		y2 = centers[len(centers)-1].Y
	}
	m := (y2 - y1) / (x2 - x1)
	if math.IsNaN(m) {
		// FIXME circumcenter is exactly on midpoint of border edge. Need to get line by using angles.
		panic(fmt.Errorf("delaunay: not implemented for the constellation in triangle T%v", last))
	}
	var side int
	const (
		top = iota + 1
		right
		bottom
		left
	)
	if x2 > x1 { // right
		cy := m*v.maxX - m*x1 + y1
		if y2 > y1 { // top
			if cy < v.maxY {
				side = right
				centers = append(centers, NewPoint(v.maxX, cy, -1))
			} else {
				side = top
				cx := v.maxY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.maxY, -1))
			}
		} else if y2 < y1 { // bottom
			if cy >= v.minY {
				side = right
				centers = append(centers, NewPoint(v.maxX, cy, -1))
			} else {
				side = bottom
				cx := v.minY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.minY, -1))
			}
		} else { // horizontal
			side = right
			centers = append(centers, NewPoint(v.maxX, cy, -1))
		}
	} else if x2 < x1 { // left
		cy := m*v.minX - m*x1 + y1
		if y2 > y1 { // top
			if cy <= v.maxY {
				side = left
				centers = append(centers, NewPoint(v.minX, cy, -1))
			} else {
				side = top
				cx := v.maxY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.maxY, -1))
			}
		} else if y2 < y1 { // bottom
			if cy > v.minY {
				side = left
				centers = append(centers, NewPoint(v.minX, cy, -1))
			} else {
				side = bottom
				cx := v.minY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.minY, -1))
			}
		} else { // horizontal
			side = left
			centers = append(centers, NewPoint(v.minX, cy, -1))
		}
	}
	if side == 0 {
		panic(fmt.Errorf("delaunay: internal error calculating voronoi cell for border point P%v", p))
	}
	// find second border point
	switch {
	case p.Equals(first.A):
		pt = first.C
		pt2 = first.B
	case p.Equals(first.B):
		pt = first.A
		pt2 = first.C
	case p.Equals(first.C):
		pt = first.B
		pt2 = first.A
	default:
		panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", p, t))
	}
	if aboveLine(p, pt, centers[0]) == aboveLine(p, pt, pt2) {
		x1 = centers[0].X
		y1 = centers[0].Y
		x2 = (p.X + pt.X) / 2
		y2 = (p.Y + pt.Y) / 2
	} else {
		x1 = (p.X + pt.X) / 2
		y1 = (p.Y + pt.Y) / 2
		x2 = centers[0].X
		y2 = centers[0].Y
	}
	m = (y2 - y1) / (x2 - x1)
	if math.IsNaN(m) {
		// FIXME Circumcenter is exactly on midpoint of border edge. Need to get line by using angles.
		panic(fmt.Errorf("delaunay: not implemented for the constellation in triangle T%v", first))
	}
	if x2 > x1 { // right
		cy := m*v.maxX - m*x1 + y1
		if y2 > y1 { // top
			if cy < v.maxY {
				if side == top {
					centers = append(centers, NewPoint(v.maxX, v.maxY, -1))
				}
				centers = append(centers, NewPoint(v.maxX, cy, -1))
			} else {
				if side == left {
					centers = append(centers, NewPoint(v.minX, v.maxY, -1))
				}
				cx := v.maxY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.maxY, -1))
			}
		} else if y2 < y1 { // bottom
			if cy >= v.minY {
				if side == top {
					centers = append(centers, NewPoint(v.maxX, v.maxY, -1))
				}
				centers = append(centers, NewPoint(v.maxX, cy, -1))
			} else {
				if side == right {
					centers = append(centers, NewPoint(v.maxX, v.minY, -1))
				}
				cx := v.minY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.minY, -1))
			}
		} else { // horizontal
			if side == top {
				centers = append(centers, NewPoint(v.maxX, v.maxY, -1))
			}
			centers = append(centers, NewPoint(v.maxX, cy, -1))
		}
	} else if x2 < x1 { // left
		cy := m*v.minX - m*x1 + y1
		if y2 > y1 { // top
			if cy <= v.maxY {
				if side == bottom {
					centers = append(centers, NewPoint(v.minX, v.minY, -1))
				}
				centers = append(centers, NewPoint(v.minX, cy, -1))
			} else {
				if side == left {
					centers = append(centers, NewPoint(v.minX, v.maxY, -1))
				}
				cx := v.maxY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.maxY, -1))
			}
		} else if y2 < y1 { // bottom
			if cy > v.minY {
				if side == bottom {
					centers = append(centers, NewPoint(v.minX, v.minY, -1))
				}
				centers = append(centers, NewPoint(v.minX, cy, -1))
			} else {
				if side == right {
					centers = append(centers, NewPoint(v.maxX, v.minY, -1))
				}
				cx := v.minY/m + x1 - y1/m
				centers = append(centers, NewPoint(cx, v.minY, -1))
			}
		} else { // horizontal
			if side == bottom {
				centers = append(centers, NewPoint(v.minX, v.minY, -1))
			}
			centers = append(centers, NewPoint(v.minX, cy, -1))
		}
	}
	return math.Inf(1), centers
}

// aboveLine determines if a point is above or below the line formed by two other points
func aboveLine(lp1, lp2, p *Point) bool {
	m := (lp1.Y - lp2.Y) / (lp1.X - lp2.X)
	c := -m*lp1.X + lp1.Y
	b := p.Y >= m*p.X+c
	return b
}

// polyArea finds the area of an irregular polygon
// needs the points in clockwise order
func polyArea(points []*Point) float64 {
	var area float64
	j := len(points) - 1
	for i := 0; i < len(points); i++ {
		area += (points[j].X + points[i].X) * (points[j].Y - points[i].Y)
		j = i
	}
	return area * 0.5
}

// findClockTri finds the next triangle in clockwise order
func findClockTri(p *Point, t *Triangle) *Triangle {
	// points in a triangle are ordered counter clockwise
	var p2 *Point
	// find point counterclockwise of p
	switch {
	case p.Equals(t.A):
		p2 = t.B
	case p.Equals(t.B):
		p2 = t.C
	case p.Equals(t.C):
		p2 = t.A
	default:
		panic(fmt.Errorf("delaunay: can't find point P%v in Triangle T%v", p, t))
	}
	for _, t1 := range p.adjacentTriangles {
		for _, t2 := range p2.adjacentTriangles {
			if !t1.Equals(t) && t1.Equals(t2) {
				return t1
			}
		}
	}
	return nil
}
