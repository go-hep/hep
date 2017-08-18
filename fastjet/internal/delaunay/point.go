// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math"

	"go-hep.org/x/hep/fastjet/internal/predicates"
)

// Point in the X-Y Plane.
//
// It holds dynamic information about the
// adjacent triangles, the nearest neighbor and the distance to that neighbor.
type Point struct {
	x, y              float64   // x and y are the coordinates of the point.
	adjacentTriangles triangles // adjacentTriangles is a list of triangles containing the point.
	nearest           *Point
	dist2             float64 // dist2 is the squared distance to the nearest neighbor.
	// id is used when points are removed. Copies of the points around the point to
	// be removed are made. The ID is set incremental in counterclockwise order. It identifies
	// the original. It is also used to determine whether a Triangle is inside or outside the
	// polygon formed by all those points.
	id int
}

// NewPoint returns Point for the given x,y coordinates
func NewPoint(x, y float64) *Point {
	return &Point{
		x:     x,
		y:     y,
		dist2: math.Inf(1),
	}
}

// NearestNeighbor returns the nearest neighbor and the distance to that neighbor.
func (p *Point) NearestNeighbor() (*Point, float64) {
	return p.nearest, math.Sqrt(p.dist2)
}

// Coordinates returns the x,y coordinates of a Point.
func (p *Point) Coordinates() (x, y float64) {
	return p.x, p.y
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

// Equals checks whether two points are the same.
func (p *Point) Equals(v *Point) bool {
	return p == v || (p.x == v.x && p.y == v.y)
}

// distance returns the squared distance between two points.
func (p *Point) distance(v *Point) float64 {
	dx := p.x - v.x
	dy := p.y - v.y
	return dx*dx + dy*dy
}

// findNearest looks at all adjacent points of p and finds the nearest one.
// p's nearest neighbor will be updated.
func (p *Point) findNearest() {
	var newNearest *Point
	min := math.Inf(1)
	for _, t := range p.adjacentTriangles {
		var dist float64
		var np *Point
		// find the point in t that is closest to p, but that is not p.
		switch {
		case p.Equals(t.A):
			distB := p.distance(t.B)
			distC := p.distance(t.C)
			if distB <= distC {
				dist = distB
				np = t.B
			} else {
				dist = distC
				np = t.C
			}
		case p.Equals(t.B):
			distA := p.distance(t.A)
			distC := p.distance(t.C)
			if distA <= distC {
				dist = distA
				np = t.A
			} else {
				dist = distC
				np = t.C
			}
		case p.Equals(t.C):
			distA := p.distance(t.A)
			distB := p.distance(t.B)
			if distA <= distB {
				dist = distA
				np = t.A
			} else {
				dist = distB
				np = t.B
			}
		default:
			panic(fmt.Errorf("delaunay: point P%s not found in T%s", p, t))
		}
		// check whether the distance found is smaller than the previous smallest distance.
		if dist < min {
			min = dist
			newNearest = np
		}
	}
	// update p's nearest Neighbor
	p.dist2 = min
	p.nearest = newNearest
}

// surroundingPoints returns the points that surround p in counterclockwise order.
func (p *Point) surroundingPoints() []*Point {
	points := make([]*Point, len(p.adjacentTriangles))
	t := p.adjacentTriangles[0]
	// j is the index of the previous point
	j := 1
	// k is the index of the previous triangle
	k := 0
	switch {
	case p.Equals(t.A):
		points[0] = t.B
		points[1] = t.C
	case p.Equals(t.B):
		points[0] = t.C
		points[1] = t.A
	case p.Equals(t.C):
		points[0] = t.A
		points[1] = t.B
	default:
		panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", p, t))
	}
	for i := 0; j < len(points)-1; {
		if i >= len(p.adjacentTriangles) {
			panic(fmt.Errorf("delaunay: internal error with adjacent triangles for %v. Can't find counterclockwise neighbor of %v", p, points[j]))
		}
		// it needs to find the triangle next to k and not k again
		if p.adjacentTriangles[i].Equals(p.adjacentTriangles[k]) {
			i++
			continue
		}
		t = p.adjacentTriangles[i]
		switch {
		case points[j].Equals(t.A):
			j++
			points[j] = t.B
			k = i
			// start the loop over
			i = 0
			continue
		case points[j].Equals(t.B):
			j++
			points[j] = t.C
			k = i
			// start the loop over
			i = 0
			continue
		case points[j].Equals(t.C):
			j++
			points[j] = t.A
			k = i
			// start the loop over
			i = 0
			continue
		}
		i++
	}
	return points
}

// inTriangle checks whether the point is in the triangle and whether it is on an edge.
func (p *Point) inTriangle(t *Triangle) location {
	o1 := predicates.Orientation(t.A.x, t.A.y, t.B.x, t.B.y, p.x, p.y)
	o2 := predicates.Orientation(t.B.x, t.B.y, t.C.x, t.C.y, p.x, p.y)
	o3 := predicates.Orientation(t.C.x, t.C.y, t.A.x, t.A.y, p.x, p.y)
	if o1 == predicates.CCW && o2 == predicates.CCW && o3 == predicates.CCW {
		return inside
	}
	if o1 == predicates.CW || o2 == predicates.CW || o3 == predicates.CW {
		return outside
	}
	return onEdge
}

// location is the position of a point relative to a triangle
type location int

const (
	inside location = iota
	onEdge
	outside
)

func (l location) String() string {
	switch l {
	case inside:
		return "Inside Triangle"
	case onEdge:
		return "On Edge of Triangle"
	case outside:
		return "Outside Triangle"
	default:
		panic(fmt.Errorf("delaunay: unknown location %d", int(l)))
	}
}
