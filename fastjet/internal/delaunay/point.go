// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math"
	"math/big"

	"github.com/gonum/floats"
)

// Point in the X-Y Plane, that holds dynamic information about the
// adjacent triangles, the nearest neighbor and the distance to that neighbor
type Point struct {
	// X and Y are the coordinates of the point.
	X, Y float64
	// adjacentTriangles is a list of triangles containing the point.
	adjacentTriangles triangles
	// isOutside indicates whether the point is outside the triangulation.
	isOutside bool
	// id is used when points are removed. Copies of the points around the point to
	// be removed are made. The ID is set incremental in counterclockwise order. It identifies
	// the original. It is also used to determine if a Triangle is inside or outside the
	// polygon formed by all those points.
	id      int
	nearest *Point
	// dist is the squared distance to the nearest neighbor.
	dist float64
}

// NewPoint returns Point for the given x,y coordinates
func NewPoint(x float64, y float64) *Point {
	return &Point{
		X:         x,
		Y:         y,
		isOutside: false,
		dist:      math.Inf(1),
	}
}

// inTriangle checks whether the point is in the triangle and whether it is on an edge
func (p *Point) inTriangle(t *Triangle) (inside, edge bool) {
	e1 := p.orientation(t.A, t.B)
	e2 := p.orientation(t.B, t.C)
	e3 := p.orientation(t.C, t.A)
	if e1.Cmp(zero) > 0 && e2.Cmp(zero) > 0 && e3.Cmp(zero) > 0 {
		// point is inside the triangle
		return true, false
	}
	// if all orientations are greater than 0 and one is 0 then the point is on an edge
	edge = e1.Cmp(zero) >= 0 && e2.Cmp(zero) >= 0 && e3.Cmp(zero) >= 0
	return edge, edge
}

// orientation returns the orientation between the two points forming a line and p.
// If b is counterclockwise of a then orientation returns
// greater than 0 if p is inside the edge, less than 0 if p is on the other side of the edge
// and 0 if p is on the edge.
// (b.X-a.X)*(p.Y-a.Y) - (p.X-a.X)*(b.Y-a.Y)
func (p *Point) orientation(a, b *Point) *big.Float {
	// use the math/big package to handle the geometric predicates
	var sum1, sum2, sum3, sum4, prod1, prod2, nprod2, result big.Float
	bx := big.NewFloat(b.X)
	by := big.NewFloat(b.Y)
	px := big.NewFloat(p.X)
	py := big.NewFloat(p.Y)
	nax := big.NewFloat(-a.X)
	nay := big.NewFloat(-a.Y)
	sum1.Add(bx, nax)
	sum2.Add(py, nay)
	sum3.Add(px, nax)
	sum4.Add(by, nay)
	prod1.Mul(&sum1, &sum2)
	prod2.Mul(&sum3, &sum4)
	nprod2.Neg(&prod2)
	return &result.Add(&prod1, &nprod2)
}

// NearestNeighbor returns the nearest Neighbor and the distance to that neighbor
func (p *Point) NearestNeighbor() (*Point, float64) {
	return p.nearest, math.Sqrt(p.dist)
}

// Equals checks whether two points are the same
func (p *Point) Equals(v *Point) bool {
	return p == v || (p.X == v.X && p.Y == v.Y)
}

// EqualsApprox compares whether v's coordinates are within a given tolerance of p
func (p *Point) EqualsApprox(v *Point, tolerance float64) bool {
	return p == v || (floats.EqualWithinAbs(p.X, v.X, tolerance) && floats.EqualWithinAbs(p.Y, v.Y, tolerance))
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}

// distance returns the squared distance between two points
func (p *Point) distance(v *Point) float64 {
	dx := p.X - v.X
	dy := p.Y - v.Y
	return dx*dx + dy*dy
}

// findNearest looks at all adjacent points of p and finds the nearest one
func (p *Point) findNearest() {
	var np *Point
	min := math.Inf(1)
	for _, t := range p.adjacentTriangles {
		switch {
		case p.Equals(t.A):
			dist := p.distance(t.B)
			if dist < min {
				min = dist
				np = t.B
			}
			dist = p.distance(t.C)
			if dist < min {
				min = dist
				np = t.C
			}
		case p.Equals(t.B):
			dist := p.distance(t.A)
			if dist < min {
				min = dist
				np = t.A
			}
			dist = p.distance(t.C)
			if dist < min {
				min = dist
				np = t.C
			}
		case p.Equals(t.C):
			dist := p.distance(t.A)
			if dist < min {
				min = dist
				np = t.A
			}
			dist = p.distance(t.B)
			if dist < min {
				min = dist
				np = t.B
			}
		default:
			panic(fmt.Errorf("delaunay: point P%s not found in T%s", p, t))
		}
	}
	p.dist = min
	p.nearest = np
}

// findBorderConnectors finds all points to the left and right of t that can connect to p
// without crossing a line. Therefore those point must be on the border and the edge to the
// previous point must face p.
func (p *Point) findBorderConnectors(l, r *Point, t *Triangle) (outerL, outerR []*Point) {
	// r is counterclockwise of l
	previousL, previousR := r, l
	outerL = append(outerL, l)
	outerR = append(outerR, r)
	var s *Point
	switch {
	case r.Equals(t.A):
		s = t.B
	case r.Equals(t.B):
		s = t.C
	case r.Equals(t.C):
		s = t.A
	default:
		panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", r, t))
	}
	// check if the third point can connect to p
	if p.orientation(r, s).Cmp(zero) < 0 {
		border := true
	loop1:
		for _, t1 := range s.adjacentTriangles {
			for _, t2 := range r.adjacentTriangles {
				if t1.Equals(t2) && !t1.Equals(t) {
					border = false
					break loop1
				}
			}
		}
		if border {
			outerR = append(outerR, s)
			previousR = r
		}
	} else if p.orientation(s, l).Cmp(zero) < 0 {
		border := true
	loop2:
		for _, t1 := range s.adjacentTriangles {
			for _, t2 := range l.adjacentTriangles {
				if t1.Equals(t2) && !t1.Equals(t) {
					border = false
					break loop2
				}
			}
		}
		if border {
			outerL = append(outerL, s)
			previousL = l
		}
	}
	found1 := true
	// l's clockwise neighbor needs to be checked
	for i := len(outerL) - 1; found1; i++ {
	outerloop1:
		for _, t1 := range outerL[i].adjacentTriangles {
			found1 = true
			if t1.A.Equals(previousL) || t1.B.Equals(previousL) || t1.C.Equals(previousL) {
				found1 = false
				continue
			}
			// potential new point
			var np *Point
			var p3 *Point
			switch {
			case outerL[i].Equals(t1.A):
				np = t1.C
				p3 = t1.B
			case outerL[i].Equals(t1.B):
				np = t1.A
				p3 = t1.C
			case outerL[i].Equals(t1.C):
				np = t1.B
				p3 = t1.A
			default:
				panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", outerL[i], t1))
			}
			// find point that is on the outside
			// point is on the outside if it only has one triangle in common with the point found last
			for _, t2 := range np.adjacentTriangles {
				if !t1.Equals(t2) && (t2.A.Equals(outerL[i]) || t2.B.Equals(outerL[i]) || t2.C.Equals(outerL[i])) {
					found1 = false
					continue outerloop1
				}
			}
			// check if it would cross a line when connecting to p
			if p.orientation(np, outerL[i]).Cmp(zero) >= 0 {
				found1 = false
				break
			}
			previousL = outerL[i]
			outerL = append(outerL, np)
			// check if third point in that triangle can connect to p
			if p.orientation(p3, np).Cmp(zero) < 0 {
				border := true
			outloop1:
				for _, tr1 := range p3.adjacentTriangles {
					for _, tr2 := range np.adjacentTriangles {
						if tr1.Equals(tr2) && !tr1.Equals(t1) {
							border = false
							break outloop1
						}
					}
				}
				if border {
					i++
					previousL = outerL[i]
					outerL = append(outerL, p3)
				}
			}
			break
		}
	}
	found2 := true
	for i := len(outerR) - 1; found2; i++ {
		// r's counterclockwise neighbor needs to be checked
	outerloop2:
		for _, t1 := range outerR[i].adjacentTriangles {
			found2 = true
			if t1.A.Equals(previousR) || t1.B.Equals(previousR) || t1.C.Equals(previousR) {
				found2 = false
				continue
			}
			// potential new point
			var np *Point
			var p3 *Point
			switch {
			case outerR[i].Equals(t1.A):
				np = t1.B
				p3 = t1.C
			case outerR[i].Equals(t1.B):
				np = t1.C
				p3 = t1.A
			case outerR[i].Equals(t1.C):
				np = t1.A
				p3 = t1.B
			default:
				panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", outerR[i], t1))
			}
			// find point that is on the outside
			// point is on the outside if it only has one triangle in common with the last point
			for _, t2 := range np.adjacentTriangles {
				if !t1.Equals(t2) && (t2.A.Equals(outerR[i]) || t2.B.Equals(outerR[i]) || t2.C.Equals(outerR[i])) {
					found2 = false
					continue outerloop2
				}
			}
			// check if it would cross a line when connecting to p
			if p.orientation(outerR[i], np).Cmp(zero) >= 0 {
				found2 = false
				break
			}
			outerR = append(outerR, np)
			previousR = outerR[i]
			// check if third point in that triangle can connect to p
			if p.orientation(np, p3).Cmp(zero) < 0 {
				border := true
			outloop2:
				for _, tr1 := range p3.adjacentTriangles {
					for _, tr2 := range np.adjacentTriangles {
						if tr1.Equals(tr2) && !tr1.Equals(t1) {
							border = false
							break outloop2
						}
					}
				}
				if border {
					i++
					previousR = outerR[i]
					outerR = append(outerR, p3)
				}
			}
			break
		}
	}
	return outerL, outerR
}

// findRemainingSurrounding finds the points around a point to be removed in clockwise order.
// It is only used by the walk method when a point is on the border and not all points surrounding the
// point have been found by going counterclockwise.
// last is the index of the last point found in the points slice
func (p *Point) findRemainingSurrounding(outer []*Point, last int) []*Point {
	if len(outer)-1 > last {
		// need to find remaining points
		// here it needs to find the points in clockwise order from the starting point,
		// because going counterclockwise stopped when the border was reached
		t := p.adjacentTriangles[0]
		// j is the index of the previous point
		j := 0
		// k is the index of the previous triangle
		k := 0
		for i := 0; j > last+1 || j == 0; {
			if i >= len(p.adjacentTriangles) {
				panic(fmt.Errorf("delaunay: internal error with adjacent triangles for P%v. Can't find clockwise neighbor of P%v", p, outer[j]))
			}
			// it needs to find the triangle next to k and not k again
			if p.adjacentTriangles[i].Equals(p.adjacentTriangles[k]) {
				i++
				continue
			}
			t = p.adjacentTriangles[i]
			switch {
			case outer[j].Equals(t.A):
				if j == 0 {
					j = len(outer)
				}
				j--
				outer[j] = t.C
				k = i
				// start the loop over
				i = 0
				continue
			case outer[j].Equals(t.B):
				if j == 0 {
					j = len(outer)
				}
				j--
				outer[j] = t.A
				k = i
				// start the loop over
				i = 0
				continue
			case outer[j].Equals(t.C):
				if j == 0 {
					j = len(outer)
				}
				j--
				outer[j] = t.B
				k = i
				// start the loop over
				i = 0
				continue
			}
			i++
		}
	}
	return outer
}
