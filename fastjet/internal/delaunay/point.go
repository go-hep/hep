// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math"

	"github.com/gonum/floats"
)

// Point in the X-Y Plane, that holds dynamic information about the
// adjacent triangles, the nearest neighbor and the distance to that neighbor
type Point struct {
	X, Y              float64 // X and Y are the coordinates of the point.
	adjacentTriangles triangles
	isOutside         bool // Indicates whether the point is outside the triangulation.
	ID                int  // user defined identifier. There are no restrictions on how the user uses it.
	// . The field is used when points are removed
	nearest *Point
	dist    float64 // squared distance to the nearest neighbor
}

// NewPoint returns Point for the given x,y coordinates and id
func NewPoint(x float64, y float64, id int) *Point {
	return &Point{
		X:         x,
		Y:         y,
		isOutside: false,
		ID:        id,
		dist:      math.Inf(1),
	}
}

// inTriangle checks whether the point is in the triangle and whether it is on an edge
// using barycentric coordinates
func (p *Point) inTriangle(t *Triangle) (inside, edge bool) {
	barcen1 := ((t.B.Y-t.C.Y)*(p.X-t.C.X) + (t.C.X-t.B.X)*(p.Y-t.C.Y)) / det(t)
	barcen2 := ((t.C.Y-t.A.Y)*(p.X-t.C.X) + (t.A.X-t.C.X)*(p.Y-t.C.Y)) / det(t)
	barcen3 := 1 - barcen1 - barcen2

	// inside triangle
	if barcen1 > 0 && barcen1 < 1 && barcen2 > 0 && barcen2 < 1 && barcen3 > 0 && barcen3 < 1 {
		return true, false
	}
	// either outside triangle or on edge
	in := barcen1 < 0 || barcen2 < 0 || barcen3 < 0
	return !in, !in
}

func det(t *Triangle) float64 {
	return (t.B.Y-t.C.Y)*(t.A.X-t.C.X) + (t.C.X-t.B.X)*(t.A.Y-t.C.Y)
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
	return fmt.Sprintf("(%f, %f ID:%d)", p.X, p.Y, p.ID)
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
