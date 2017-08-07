// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"

	"go-hep.org/x/hep/fastjet/internal/predicates"
)

// Triangle is a set of three points that make up a triangle. It stores hierarchical information to find triangles.
type Triangle struct {
	// children are triangles that lead to the removal of this triangle.
	// When a point is inserted the triangle that contains the point is found by going down the hierarchical tree.
	// The tree's root is the root triangle and the children slice contains the children triangles.
	children []*Triangle
	// A,B,C are the CCW-oriented points that make up the triangle
	A, B, C *Point
	// isInTriangulation holds whether a triangle is part of the triangulation
	isInTriangulation bool
}

// NewTriangle returns a triangle formed out of the three given points.
func NewTriangle(a, b, c *Point) *Triangle {
	// order the points counter clockwise
	o := predicates.Orientation(a.x, a.y, b.x, b.y, c.x, c.y)
	switch o {
	case predicates.CW:
		a, b = b, a
	case predicates.Colinear:
		panic(fmt.Errorf("delaunay: Can't form triangle, because Points a%v, b%v and c%v are colinear.", a, b, c))
	}
	return &Triangle{
		A: a,
		B: b,
		C: c,
	}
}

// add is used when the triangle t is added to the delaunay triangulation.
// It updates the information of the points in t.
// It returns all points whose nearest neighbor was updated.
func (t *Triangle) add() []*Point {
	t.isInTriangulation = true
	// update the adjacent lists
	t.A.adjacentTriangles = append(t.A.adjacentTriangles, t)
	t.B.adjacentTriangles = append(t.B.adjacentTriangles, t)
	t.C.adjacentTriangles = append(t.C.adjacentTriangles, t)
	// update nearest neighbor if one of the points is closer than current nearest.
	// First find the nearest neighbor in t.
	// Then check whether local min distance is smaller than current min distance.
	updated := make([]*Point, 0, 3)
	distAB := t.A.distance(t.B)
	distBC := t.B.distance(t.C)
	distCA := t.C.distance(t.A)
	var localMin float64
	var localNearest *Point
	if distAB < distCA {
		localMin = distAB
		localNearest = t.B
	} else {
		localMin = distCA
		localNearest = t.C
	}
	if localMin < t.A.dist2 {
		t.A.dist2 = localMin
		t.A.nearest = localNearest
		updated = append(updated, t.A)
	}
	if distAB < distBC {
		localMin = distAB
		localNearest = t.A
	} else {
		localMin = distBC
		localNearest = t.C
	}
	if localMin < t.B.dist2 {
		t.B.dist2 = localMin
		t.B.nearest = localNearest
		updated = append(updated, t.B)
	}
	if distBC < distCA {
		localMin = distBC
		localNearest = t.B
	} else {
		localMin = distCA
		localNearest = t.A
	}
	if localMin < t.C.dist2 {
		t.C.dist2 = localMin
		t.C.nearest = localNearest
		updated = append(updated, t.C)
	}
	return updated
}

// remove is used when the triangle t is removed from the delaunay triangulation.
// It updates the information of the points in t.
// It returns all points whose nearest neighbor was updated.
func (t *Triangle) remove() []*Point {
	t.isInTriangulation = false
	t.A.adjacentTriangles = t.A.adjacentTriangles.remove(t)
	t.B.adjacentTriangles = t.B.adjacentTriangles.remove(t)
	t.C.adjacentTriangles = t.C.adjacentTriangles.remove(t)
	// update the nearest neighbor if the nearest neighbor is in t.
	updated := make([]*Point, 0, 3)
	if t.A.nearest.Equals(t.B) || t.A.nearest.Equals(t.C) {
		t.A.findNearest()
		updated = append(updated, t.A)
	}
	if t.B.nearest.Equals(t.A) || t.B.nearest.Equals(t.C) {
		t.B.findNearest()
		updated = append(updated, t.B)
	}
	if t.C.nearest.Equals(t.A) || t.C.nearest.Equals(t.B) {
		t.C.findNearest()
		updated = append(updated, t.C)
	}
	return updated
}

// Equals checks whether two triangles are the same.
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

// remove removes given triangles from a slice of triangles.
//
// remove may modify the content of the input slice of triangles to remove from the receiver.
// remove assumes there are no duplicate triangles in the receiver list of triangles.
// remove will fail to remove the duplicates if that assumption does not hold.
func (ts triangles) remove(triangles ...*Triangle) triangles {
	k := len(triangles) - 1
	for i := len(ts) - 1; i >= 0; i-- {
		if k < 0 {
			break
		}
		for j := k; j >= 0; j-- {
			if triangles[j].Equals(ts[i]) {
				n := len(ts) - 1
				ts[i], ts[n] = ts[n], nil
				ts = ts[:n]
				triangles[j] = triangles[k]
				k--
				break
			}
		}
	}
	return ts
}
