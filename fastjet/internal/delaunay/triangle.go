// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math"
	"math/big"
)

// Triangle is a set of three points that make up a triangle, with hierarchical information to find triangles.
type Triangle struct {
	// children are triangles that lead to the removal of this triangle.
	// When a point is inserted the triangle that contains the point is found by going down the hierarchical tree.
	// The tree's root is the root triangle and the children slice contains the children triangles.
	children triangles
	// A,B,C are the points that make up the triangle
	A, B, C *Point
	// isInTriangulation holds whether a triangle is part of the triangulation
	isInTriangulation bool
}

// NewTriangle returns a triangle formed out of the three given points
func NewTriangle(a, b, c *Point) *Triangle {
	// make it counter clockwise
	if isClockwise(a, b, c) {
		b, c = c, b
	}
	return &Triangle{
		A: a,
		B: b,
		C: c,
	}
}

// isClockwise checks whether three points are in clockwise order
func isClockwise(a, b, c *Point) bool {
	return c.orientation(a, b).Cmp(zero) < 0
}

// inCircumcircle returns whether the point is inside the circumcircle of the triangle
func (t *Triangle) inCircumcircle(p *Point) bool {
	x, y := t.centerOfCircumcircle()
	// use the math/big package to handle the geometric predicates
	// r is the squared radius of the circumcircle
	r := big.NewFloat((x-t.A.X)*(x-t.A.X) + (y-t.A.Y)*(y-t.A.Y))
	// d is the squared distance from the point to the circumcenter
	d := big.NewFloat((p.X-x)*(p.X-x) + (p.Y-y)*(p.Y-y))
	// point is in circumcircle if squared distance to center is less than squared radius
	return d.Cmp(r) < 0
}

// centerOfCircumcircle returns the center of the triangle's circumcircle.
func (t *Triangle) centerOfCircumcircle() (x, y float64) {
	m1 := (t.A.X - t.B.X) / (t.B.Y - t.A.Y)
	m2 := (t.A.X - t.C.X) / (t.C.Y - t.A.Y)
	b1 := (t.A.Y+t.B.Y)*0.5 - m1*(t.A.X+t.B.X)*0.5
	b2 := (t.A.Y+t.C.Y)*0.5 - m2*(t.A.X+t.C.X)*0.5
	// x and y are the coordinates of the center of the circumcircle
	x = (b2 - b1) / (m1 - m2)
	y = m1*x + b1
	if math.IsNaN(y) {
		m1 = (t.B.X - t.A.X) / (t.A.Y - t.B.Y)
		m2 = (t.B.X - t.C.X) / (t.C.Y - t.B.Y)
		b1 = (t.B.Y+t.A.Y)*0.5 - m1*(t.B.X+t.A.X)*0.5
		b2 = (t.B.Y+t.C.Y)*0.5 - m2*(t.B.X+t.C.X)*0.5
		x = (b2 - b1) / (m1 - m2)
		y = m1*x + b1
		if math.IsNaN(y) {
			m1 = (t.C.X - t.A.X) / (t.A.Y - t.C.Y)
			m2 = (t.C.X - t.B.X) / (t.B.Y - t.C.Y)
			b1 = (t.C.Y+t.A.Y)*0.5 - m1*(t.C.X+t.A.X)*0.5
			b2 = (t.C.Y+t.B.Y)*0.5 - m2*(t.C.X+t.B.X)*0.5
			x = (b2 - b1) / (m1 - m2)
			y = m1*x + b1
			if math.IsNaN(y) {
				panic(fmt.Errorf("delaunay: error caluclating the circumcenter of triangle %v", t))
			}
		}
	}
	return x, y
}

// Equals checks whether two triangles are the same
func (t *Triangle) Equals(s *Triangle) bool {
	return t == s ||
		(t.A.Equals(s.A) || t.A.Equals(s.B) || t.A.Equals(s.C)) &&
			(t.B.Equals(s.A) || t.B.Equals(s.B) || t.B.Equals(s.C)) &&
			(t.C.Equals(s.A) || t.C.Equals(s.B) || t.C.Equals(s.C))
}

func (t *Triangle) String() string {
	return fmt.Sprintf("{A%s, B%s, C%s}", t.A, t.B, t.C)
}

type triangles []*Triangle

// append appends to a slice of triangles and updates the nearest neighbor
// it is used when the adjacent triangles of a point change
func (t triangles) append(triangles ...*Triangle) []*Triangle {
	// check if nearest neighbor changes by going through each triangles points
	// and checking if the distance to that point is less. It is done both ways.
	for _, tri := range triangles {
		d := tri.A.distance(tri.B)
		if d < tri.A.dist {
			tri.A.dist = d
			tri.A.nearest = tri.B
		}
		if d < tri.B.dist {
			tri.B.dist = d
			tri.B.nearest = tri.A
		}
		d = tri.B.distance(tri.C)
		if d < tri.B.dist {
			tri.B.dist = d
			tri.B.nearest = tri.C
		}
		if d < tri.C.dist {
			tri.C.dist = d
			tri.C.nearest = tri.B
		}
		d = tri.A.distance(tri.C)
		if d < tri.A.dist {
			tri.A.dist = d
			tri.A.nearest = tri.C
		}
		if d < tri.C.dist {
			tri.C.dist = d
			tri.C.nearest = tri.A
		}
	}
	return append(t, triangles...)
}

// remove removes given triangles from a slice of triangles
func (t triangles) remove(triangles ...*Triangle) []*Triangle {
	// check if nearest neighbor of any point is removed and if so
	// put that point in the update slice
	var update []*Point
	for _, tri := range triangles {
		if tri.A.nearest != nil && (tri.A.nearest.Equals(tri.B) || tri.A.nearest.Equals(tri.C)) {
			update = append(update, tri.A)
		}
		if tri.B.nearest != nil && (tri.B.nearest.Equals(tri.A) || tri.B.nearest.Equals(tri.C)) {
			update = append(update, tri.B)
		}
		if tri.C.nearest != nil && (tri.C.nearest.Equals(tri.A) || tri.C.nearest.Equals(tri.B)) {
			update = append(update, tri.C)
		}
	}
	for i := len(t) - 1; i >= 0; i-- {
		for j, tri := range triangles {
			if tri.Equals(t[i]) {
				t = append(t[:i], t[i+1:]...)
				triangles = append(triangles[:j], triangles[j+1:]...)
				break
			}
		}
	}
	// find the new nearest neighbor for the points who's nearest neighbor was removed
	for _, p := range update {
		p.findNearest()
	}
	return t
}

// finalize returns the final delaunay triangles
// only keeps leaf elements from the hierarchy
// removes given triangles
func (t triangles) finalize(triangles ...*Triangle) []*Triangle {
	ft := make([]*Triangle, 0, len(t))
	for _, a := range t {
		if a.isInTriangulation {
			keep := true
			for j, b := range triangles {
				if b.Equals(a) {
					keep = false
					triangles = append(triangles[:j], triangles[j+1:]...)
					break
				}
			}
			if keep {
				ft = append(ft, a)
			}
		}
	}
	return ft
}
