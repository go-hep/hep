// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math"

	"go-hep.org/x/hep/fastjet/internal/predicates"
	"golang.org/x/xerrors"
)

// Point in the X-Y Plane.
//
// It holds dynamic information about the
// adjacent triangles, the nearest neighbor and the distance to that neighbor.
//
// One should use the Equal method of Point to test whether 2 points are equal.
type Point struct {
	x, y              float64   // x and y are the coordinates of the point.
	adjacentTriangles triangles // adjacentTriangles is a list of triangles containing the point.
	nearest           *Point
	dist2             float64 // dist2 is the squared distance to the nearest neighbor.
	// id is a unique identifier, that is assigned incrementally to a point on insertion.
	// It is used when points are removed. Copies of the points around the point to be removed are made.
	// The ID is set incremental in counterclockwise order. It identifies the original. It is also used
	// to determine whether a Triangle is inside or outside the polygon formed by all those points.
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

// ID returns the ID of the point. It is a unique identifier that is incrementally assigned
// to a point on insertion.
func (p *Point) ID() int {
	return p.id
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.x, p.y)
}

// Equals checks whether two points are the same.
func (p *Point) Equals(v *Point) bool {
	return p == v || (p.x == v.x && p.y == v.y)
}

// SecondNearestNeighbor looks at all adjacent points of p and returns the second nearest one
// and the distance to that point.
func (p *Point) SecondNearestNeighbor() (*Point, float64) {
	var nearest, secondNearest *Point
	min, secondMin := math.Inf(1), math.Inf(1)
	if p.dist2 == 0 {
		// p has a duplicate
		nearest = p.nearest
		min = 0
	}
	for _, t := range p.adjacentTriangles {
		var p2, p3 *Point
		// find the point in t that is not p
		switch {
		case p.Equals(t.A):
			p2 = t.B
			p3 = t.C
		case p.Equals(t.B):
			p2 = t.A
			p3 = t.C
		case p.Equals(t.C):
			p2 = t.A
			p3 = t.B
		default:
			panic(xerrors.Errorf("delaunay: point %v not found in %v", p, t))
		}
		dist := p.distance(p2)
		switch {
		case dist < min:
			min, secondMin = dist, min
			nearest, secondNearest = p2, nearest
			if p2.dist2 == 0 {
				// p2 has a duplicate
				secondNearest = p2.nearest
				secondMin = dist
			}
		case dist < secondMin:
			secondMin = dist
			secondNearest = p2
		}
		dist = p.distance(p3)
		switch {
		case dist < min:
			min, secondMin = dist, min
			nearest, secondNearest = p3, nearest
			if p3.dist2 == 0 {
				// p3 has a duplicate
				secondNearest = p3.nearest
				secondMin = dist
			}
		case dist < secondMin:
			secondMin = dist
			secondNearest = p3
		}
	}
	return secondNearest, math.Sqrt(secondMin)
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
			panic(xerrors.Errorf("delaunay: point P%s not found in T%s", p, t))
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
		panic(xerrors.Errorf("delaunay: point %v not in adjacent triangle %v", p, t))
	}
	for i := 0; j < len(points)-1; {
		if i >= len(p.adjacentTriangles) {
			panic(xerrors.Errorf("delaunay: internal error with adjacent triangles for %v. Can't find counterclockwise neighbor of %v", p, points[j]))
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

// findClockwiseTriangle finds the next triangle in clockwise order.
func (p *Point) findClockwiseTriangle(t *Triangle) *Triangle {
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
		panic(xerrors.Errorf("delaunay: can't find Point %v in Triangle %v", p, t))
	}
	for _, t1 := range p.adjacentTriangles {
		for _, t2 := range p2.adjacentTriangles {
			if !t1.Equals(t) && t1.Equals(t2) {
				return t1
			}
		}
	}
	panic(xerrors.Errorf("delaunay: no clockwise neighbor of Triangle %v around Point %v", t, p))
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
		panic(xerrors.Errorf("delaunay: unknown location %d", int(l)))
	}
}

type points []*Point

// remove removes given points from a slice of points.
//
// remove will remove all occurrences of the points.
func (ps points) remove(pts ...*Point) points {
	out := make(points, 0, len(ps))
	for _, p := range ps {
		keep := true
		for _, pt := range pts {
			if p.Equals(pt) {
				keep = false
				break
			}
		}
		if keep {
			out = append(out, p)
		}

	}
	return out
}

// polyArea finds the area of an irregular polygon.
// The points need to be in clockwise order.
func (points points) polyArea() float64 {
	var area float64
	j := len(points) - 1
	for i := 0; i < len(points); i++ {
		area += (points[j].x + points[i].x) * (points[j].y - points[i].y)
		j = i
	}
	return area * 0.5
}
