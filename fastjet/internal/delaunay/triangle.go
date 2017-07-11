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
	return (b.Y-a.Y)*(c.X-b.X)-(c.Y-b.Y)*(b.X-a.X) > 0
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
				panic(fmt.Errorf("delaunay: error caluclating the circumcenter of triangle " + t.String("T")))
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

func (t *Triangle) String(name string) string {
	return fmt.Sprintf("{A%s, B%s, C%s}", t.A, t.B, t.C)
}

type triangles []*Triangle

// appendT appends to a slice of triangles and updates the nearest neighbor
// it is used when the adjacent triangles of a point change
func (triangles triangles) appendT(elems ...*Triangle) []*Triangle {
	// check if nearest neighbor changes by going through each triangles points
	// and checking if the distance to that point is less. It is done both ways.
	for _, t := range elems {
		d := t.A.distance(t.B)
		if d < t.A.dist {
			t.A.dist = d
			t.A.nearest = t.B
		}
		if d < t.B.dist {
			t.B.dist = d
			t.B.nearest = t.A
		}
		d = t.B.distance(t.C)
		if d < t.B.dist {
			t.B.dist = d
			t.B.nearest = t.C
		}
		if d < t.C.dist {
			t.C.dist = d
			t.C.nearest = t.B
		}
		d = t.A.distance(t.C)
		if d < t.A.dist {
			t.A.dist = d
			t.A.nearest = t.C
		}
		if d < t.C.dist {
			t.C.dist = d
			t.C.nearest = t.A
		}
	}
	return append(triangles, elems...)
}

// remove removes given triangles from a slice of triangles
func (triangles triangles) remove(elems ...*Triangle) []*Triangle {
	// check if nearest neighbor of any point is removed and if so
	// put that point in the update slice
	var update []*Point
	for _, t := range elems {
		if t.A.nearest != nil && (t.A.nearest.Equals(t.B) || t.A.nearest.Equals(t.C)) {
			update = append(update, t.A)
		}
		if t.B.nearest != nil && (t.B.nearest.Equals(t.A) || t.B.nearest.Equals(t.C)) {
			update = append(update, t.B)
		}
		if t.C.nearest != nil && (t.C.nearest.Equals(t.A) || t.C.nearest.Equals(t.B)) {
			update = append(update, t.C)
		}
	}
	for i := len(triangles) - 1; i >= 0; i-- {
		for j, tri := range elems {
			if tri.Equals(triangles[i]) {
				triangles = append(triangles[:i], triangles[i+1:]...)
				elems = append(elems[:j], elems[j+1:]...)
				break
			}
		}
	}
	// find the new nearest neighbor for the points who's nearest neighbor was removed
	for _, p := range update {
		p.findNearest()
	}
	return triangles
}

// finalize returns the final delaunay triangles
// only keeps leaf elements from the hierarchy
// removes given triangles
func (triangles triangles) finalize(elems ...*Triangle) []*Triangle {
	ft := make([]*Triangle, 0, len(triangles))
	for i, t := range triangles {
		if len(triangles[i].children) == 0 {
			keep := true
			for j, tri := range elems {
				if tri.Equals(t) {
					keep = false
					elems = append(elems[:j], elems[j+1:]...)
					break
				}
			}
			if keep {
				ft = append(ft, t)
			}
		}
	}
	return ft
}
