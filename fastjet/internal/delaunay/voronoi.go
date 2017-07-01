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
	maxX float64
	minX float64
	maxY float64
	minY float64
}

// FIXME can't do any delaunay operation after calling NewVoronoi
func NewVoronoi(d *Delaunay) *Voronoi {
	border := make([]*Triangle, 0, len(d.root.A.adjacentTriangles)+len(d.root.B.adjacentTriangles)+len(d.root.C.adjacentTriangles))
	border = append(border, d.root.A.adjacentTriangles[:]...)
	border = append(border, d.root.B.adjacentTriangles[:]...)
	border = append(border, d.root.C.adjacentTriangles[:]...)
	for i, t := range d.triangles {
		if len(d.triangles[i].children) == 0 {
			for j, tri := range border {
				if tri.Equals(t) {
					border = append(border[:j], border[j+1:]...)
					t.A.isOutside = true
					if !t.A.Equals(d.root.A) && !t.A.Equals(d.root.B) && !t.A.Equals(d.root.C) {
						t.A.adjacentTriangles = remove(t.A.adjacentTriangles, t)
					}
					t.B.isOutside = true
					if !t.B.Equals(d.root.A) && !t.B.Equals(d.root.B) && !t.B.Equals(d.root.C) {
						t.B.adjacentTriangles = remove(t.B.adjacentTriangles, t)
					}
					t.C.isOutside = true
					if !t.C.Equals(d.root.A) && !t.C.Equals(d.root.B) && !t.C.Equals(d.root.C) {
						t.C.adjacentTriangles = remove(t.C.adjacentTriangles, t)
					}
					break
				}
			}
		}
	}
	v := &Voronoi{maxX: d.maxX, minX: d.minX, maxY: d.maxY, minY: d.minY}
	return v
}

// finds the voronoi area of a point
// returns the area and the points associated with it
func (v *Voronoi) VoronoiCell(p *Point) (float64, []*Point) {
	// find all points that form the voronoi cell for that point
	// in clockwise order
	centers := make([]*Point, len(p.adjacentTriangles), len(p.adjacentTriangles)+3)
	var t, first *Triangle
	if p.isOutside {
		// find first triangle on the outside
		for _, t1 := range p.adjacentTriangles {
			var pt *Point
			if t1.A.Equals(p) {
				pt = t1.C
			} else if t1.B.Equals(p) {
				pt = t1.A
			} else if t1.C.Equals(p) {
				pt = t1.B
			} else {
				// should never happen
				panic(fmt.Errorf("voronoi: internal error with adjacent triangles for point P%s", p))
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
			// should only happen when user makes a mistake
			panic(fmt.Errorf("voronoi: internal error for point P%s", p))
		}
	} else {
		t = p.adjacentTriangles[0]
	}
	x, y := t.centerOfCircumCircle()
	c := NewPoint(x, y, -1)
	centers[0] = c
	// if p has only one adjacent triangle last is first
	last := first
	for i := 1; i < len(p.adjacentTriangles); i++ {
		t = findClockTri(p, t)
		x, y = t.centerOfCircumCircle()
		c = NewPoint(x, y, -1)
		centers[i] = c
		if i == len(p.adjacentTriangles)-1 {
			last = t
		}
	}
	if p.isOutside {
		// FIXME problem when circumcenter is outside border
		// extend polygon to border
		var pt, pt2 *Point
		if p.Equals(last.A) {
			pt = last.B
			pt2 = last.C
		} else if p.Equals(last.B) {
			pt = last.C
			pt2 = last.A
		} else if p.Equals(last.C) {
			pt = last.A
			pt2 = last.B
		} else {
			// should never happen
			panic(fmt.Errorf("voronoi: internal error with adjacent triangles for point P%s", p))
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
			// FIXME
			panic(fmt.Errorf("delaunay: not implemented for the constellation in triangle T%s", last))
		}
		var side int
		const top, right, bottom, left = 0, 1, 2, 3
		if x2-x1 > 0 { // right
			cy := m*v.maxX - m*x1 + y1
			if y2-y1 > 0 { // top
				if cy < v.maxY {
					side = right
					centers = append(centers, NewPoint(v.maxX, cy, -1))
				} else {
					side = top
					cx := v.maxY/m + x1 - y1/m
					centers = append(centers, NewPoint(cx, v.maxY, -1))
				}
			} else if y2-y1 < 0 { // bottom
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
		} else if x2-x1 < 0 { // left
			cy := m*v.minX - m*x1 + y1
			if y2-y1 > 0 { // top
				if cy <= v.maxY {
					side = left
					centers = append(centers, NewPoint(v.minX, cy, -1))
				} else {
					side = top
					cx := v.maxY/m + x1 - y1/m
					centers = append(centers, NewPoint(cx, v.maxY, -1))
				}
			} else if y2-y1 < 0 { // bottom
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
		// find second border point
		if p.Equals(first.A) {
			pt = first.C
			pt2 = first.B
		} else if p.Equals(first.B) {
			pt = first.A
			pt2 = first.C
		} else if p.Equals(first.C) {
			pt = first.B
			pt2 = first.A
		} else {
			// should never happen
			panic(fmt.Errorf("voronoi: internal error with adjacent triangles for point P%s", p))
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
			// FIXME
			panic(fmt.Errorf("delaunay: not implemented for the constellation in triangle T%s", first))
		}
		if x2-x1 > 0 { // right
			cy := m*v.maxX - m*x1 + y1
			if y2-y1 > 0 { // top
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
			} else if y2-y1 < 0 { // bottom
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
		} else if x2-x1 < 0 { // left
			cy := m*v.minX - m*x1 + y1
			if y2-y1 > 0 { // top
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
			} else if y2-y1 < 0 { // bottom
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
	} else {
		area := polyArea(centers)
		return area, centers
	}
}

// determines if a point is above or below the line formed by two other points
func aboveLine(lp1, lp2, p *Point) bool {
	m := (lp1.Y - lp2.Y) / (lp1.X - lp2.X)
	c := -m*lp1.X + lp1.Y
	b := p.Y >= m*p.X+c
	return b
}

// finds the area of an irregular polygon
// needs the points in clockwise order
func polyArea(points []*Point) float64 {
	area := 0.0
	j := len(points) - 1
	for i := 0; i < len(points); i++ {
		area += (points[j].X + points[i].X) * (points[j].Y - points[i].Y)
		j = i
	}
	return area / 2
}

// finds the next triangle in clockwise order
func findClockTri(p *Point, t *Triangle) *Triangle {
	// points in a triangle are ordered counter clockwise
	var p2 *Point
	if p.Equals(t.A) {
		p2 = t.B
	} else if p.Equals(t.B) {
		p2 = t.C
	} else if p.Equals(t.C) {
		p2 = t.A
	} else {
		// should only happen when user makes a mistake
		panic(fmt.Errorf("delaunay: can't find point P%s in Triangle T%s", p, t))
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
