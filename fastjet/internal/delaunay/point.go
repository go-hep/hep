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

// NewPoint returns Point for the given x,y coordinates and id
func NewPoint(x float64, y float64) *Point {
	return &Point{
		X:         x,
		Y:         y,
		isOutside: false,
		dist:      math.Inf(1),
	}
}

// inTriangle checks whether the point is in the triangle and whether it is on an edge
// using barycentric coordinates
func (p *Point) inTriangle(t *Triangle) (inside, edge bool) {
	barcen1 := ((t.B.Y-t.C.Y)*(p.X-t.C.X) + (t.C.X-t.B.X)*(p.Y-t.C.Y)) / det(t)
	barcen2 := ((t.C.Y-t.A.Y)*(p.X-t.C.X) + (t.A.X-t.C.X)*(p.Y-t.C.Y)) / det(t)
	barcen3 := 1 - barcen1 - barcen2

	// use the math/big package to handle the geometric predicates
	b1 := big.NewFloat(barcen1)
	b2 := big.NewFloat(barcen2)
	b3 := big.NewFloat(barcen3)
	// inside triangle if all barycentric coordinates are between 0 and 1
	if b1.Cmp(zero) > 0 && b1.Cmp(one) < 0 && b2.Cmp(zero) > 0 && b2.Cmp(one) < 0 && b3.Cmp(zero) > 0 && b3.Cmp(one) < 0 {
		return true, false
	}
	// either outside triangle or on edge
	// it is outside if one of the coordinates is less than 0
	// it is on the edge if one or more of the coordinates is one and the rest greater 0
	in := b1.Cmp(zero) > 0 && b2.Cmp(zero) > 0 && b3.Cmp(zero) > 0
	return in, in
}

func det(t *Triangle) float64 {
	return (t.B.Y-t.C.Y)*(t.A.X-t.C.X) + (t.C.X-t.B.X)*(t.A.Y-t.C.Y)
}

// orientation returns the orientation between the two points forming a line and p
// if b is counterclockwise of a then orientation returns
// >0 if p is inside the edge, <0 if p is on the other side of the edge
// and 0 if p is on the edge
func (p *Point) orientation(a, b *Point) float64 {
	return (b.X-a.X)*(p.Y-a.Y) - (p.X-a.X)*(b.Y-a.Y)
}

// NearestNeighbor returns the nearest Neighbor and the distance to that neighbor
func (p *Point) NearestNeighbor() (*Point, float64) {
	return p.nearest, math.Sqrt(p.dist)
}

// Equals checks whether two points are the same
func (p *Point) Equals(v *Point) bool {
	return p == v || (p.X == v.X && p.Y == v.Y)
}

// EqualsApprox compares whether b's coordinates are within tolerance of p
func (p *Point) EqualsApprox(v *Point, tol float64) bool {
	return p == v || (floats.EqualWithinAbs(p.X, v.X, tol) && floats.EqualWithinAbs(p.Y, v.Y, tol))
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f ID:%d)", p.X, p.Y, p.id)
}

// distance returns the squared distance between two points
func (p *Point) distance(v *Point) float64 {
	dx := p.X - v.X
	dy := p.Y - v.Y
	return dx*dx + dy*dy
}

// findNearest looks at all adjacent points and finds the nearest one
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

// removeOuter finds the point around a point to be removed in clockwise order.
// it is only used by the walk method
// last is the index of the last point found in the points slice
func (p *Point) removeOuter(outer []*Point, last int) []*Point {
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
