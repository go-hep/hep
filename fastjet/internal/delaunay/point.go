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
	adjacentTriangles []*Triangle
	isOutside         bool // Indicates whether the point is outside the triangulation.
	ID                int  // user defined identifier
	nearest           *Point
	dist              float64 // squared distance to the nearest neighbor
}

// return Point for the given x,y coordinates and id
func NewPoint(x float64, y float64, id int) *Point {
	return &Point{
		X:         x,
		Y:         y,
		isOutside: false,
		ID:        id,
		dist:      math.MaxFloat64,
	}
}

// returns the nearest Neighbor and the distance to that neighbor
func (p *Point) NearestNeighbor() (*Point, float64) {
	return p.nearest, math.Sqrt(p.dist)
}

func (p *Point) Equals(v *Point) bool {
	return p == v || (p.X == v.X && p.Y == v.Y)
}

// compare whether b's coordinates are within tolerance of p
func (p *Point) EqualsApprox(v *Point, tol float64) bool {
	return p == v || (floats.EqualWithinAbs(p.X, v.X, tol) && floats.EqualWithinAbs(p.Y, v.Y, tol))
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f ID:%d)", p.X, p.Y, p.ID)
}

func (p *Point) distance(v *Point) float64 {
	dx := p.X - v.X
	dy := p.Y - v.Y
	return dx*dx + dy*dy
}

// looks at all adjacent points and finds the nearest one
func (p *Point) findNearest() {
	var np *Point
	min := math.MaxFloat64
	for _, t := range p.adjacentTriangles {
		if p.Equals(t.A) {
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
		} else if p.Equals(t.B) {
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
		} else {
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
		}
	}
	p.dist = min
	p.nearest = np
}
