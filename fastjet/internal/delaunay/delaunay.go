// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math/rand"

	"go-hep.org/x/hep/fastjet/internal/predicates"
)

// Delaunay holds necessary information for the delaunay triangulation.
type Delaunay struct {
	// triangles is a slice of all triangles that have been created. It is used to get the final
	// list of triangles in the delaunay triangulation.
	triangles triangles
	// root is a triangle that contains all points. It is used as the starting point in the hierarchy
	// to locate a point. The variable root's nil-ness also indicates which method to use to locate a point.
	root *Triangle
	r    *rand.Rand
}

// HierarchicalDelaunay creates a Delaunay Triangulation using the delaunay hierarchy.
//
// The worst case time complexity is O(n*log(n)).
//
// The three root points (-2^30,-2^30), (2^30,-2^30) and (0,2^30) can't be in the circumcircle of any three non-collinear points in the
// triangulation. Additionally all points have to be inside the Triangle formed by these three points.
// If any of these conditions doesn't apply use the WalkDelaunay function instead.
// 2^30 = 1,073,741,824
//
// To locate a point this algorithm uses a Directed Acyclic Graph with a single root.
// All triangles in the current triangulation are leaf triangles.
// To find the triangle which contains the point, the algorithm follows the graph until
// a leaf is reached.
func HierarchicalDelaunay() *Delaunay {
	a := NewPoint(-1<<30, -1<<30)
	b := NewPoint(1<<30, -1<<30)
	c := NewPoint(0, 1<<30)
	root := NewTriangle(a, b, c)
	return &Delaunay{
		root: root,
	}
}

func WalkDelaunay(points []*Point, r *rand.Rand) *Delaunay {
	panic(fmt.Errorf("delaunay: WalkDelaunay not implemented"))
}

// Triangles returns the triangles that form the delaunay triangulation.
func (d *Delaunay) Triangles() []*Triangle {
	// rt are triangles that contain the root points
	rt := make(triangles, len(d.root.A.adjacentTriangles)+len(d.root.B.adjacentTriangles)+len(d.root.C.adjacentTriangles))
	n := copy(rt, d.root.A.adjacentTriangles)
	n += copy(rt[n:], d.root.B.adjacentTriangles)
	copy(rt[n:], d.root.C.adjacentTriangles)
	// make a copy so that the original slice is not modified
	triangles := make(triangles, len(d.triangles))
	copy(triangles, d.triangles)
	// keep only the triangles that are in the current triangulation
	triangles = triangles.keepCurrentTriangles()
	// remove all triangles that contain the root points
	return triangles.remove(rt...)
}

// Insert inserts the point into the triangulation. It returns the points
// whose nearest neighbor changed due to the insertion. The slice may contain
// duplicates.
func (d *Delaunay) Insert(p *Point) (updatedNearestNeighbor []*Point) {
	if len(p.adjacentTriangles) == 0 {
		p.adjacentTriangles = make(triangles, 0)
	}
	p.adjacentTriangles = p.adjacentTriangles[:0]
	t, l := d.locatePointHierarchy(p, d.root)
	switch l {
	case inside:
		return d.insertInside(p, t)
	case onEdge:
		return d.insertOnEdge(p, t)
	}
	panic(fmt.Errorf("delaunay: no triangle containing point %v", p))
}

// Remove removes the point from the triangulation. It returns the points
// whose nearest neighbor changed due to the removal. The slice may contain
// duplicates.
func (d *Delaunay) Remove(p *Point) (updatedNearestNeighbor []*Point) {
	if len(p.adjacentTriangles) < 3 {
		panic(fmt.Errorf("delaunay: can't remove point %v, not enough adjacent triangles", p))
	}
	points := p.surroundingPoints()
	ts := make([]*Triangle, len(p.adjacentTriangles))
	copy(ts, p.adjacentTriangles)
	for _, t := range ts {
		updtemp := t.remove()
		updatedNearestNeighbor = append(updatedNearestNeighbor, updtemp...)
	}
	updtemp := d.retriangulateAndSew(points, ts)
	return append(updatedNearestNeighbor, updtemp...)
}

// locatePointHierarchy locates the point using the delaunay hierarchy.
//
// It returns the triangle that contains the point and the location
// indicates whether it is on an edge or not. The worst case time complexity
// for locating a point is O(log(n)).
func (d *Delaunay) locatePointHierarchy(p *Point, t *Triangle) (*Triangle, location) {
	l := p.inTriangle(t)
	if l == outside {
		return nil, l
	}
	if len(t.children) == 0 {
		// leaf triangle
		return t, l
	}
	// go down the children to eventually reach a leaf triangle
	for _, child := range t.children {
		t, l = d.locatePointHierarchy(p, child)
		if l != outside {
			return t, l
		}
	}
	panic(fmt.Errorf("delaunay: error locating Point %v in Triangle %v", p, t))
}

// insertInside inserts a point inside a triangle. It returns the points whose nearest
// neighbor changed during the process.
func (d *Delaunay) insertInside(p *Point, t *Triangle) []*Point {
	// form three new triangles
	t1 := NewTriangle(t.A, t.B, p)
	t2 := NewTriangle(t.B, t.C, p)
	t3 := NewTriangle(t.C, t.A, p)
	updated := t1.add()
	updtemp := t2.add()
	updated = append(updated, updtemp...)
	updtemp = t3.add()
	updated = append(updated, updtemp...)
	updtemp = t.remove()
	updated = append(updated, updtemp...)
	t.children = append(t.children, t1, t2, t3)
	d.triangles = append(d.triangles, t1, t2, t3)
	// change the edges so it is a valid delaunay triangulation
	updtemp = d.swapDelaunay(t1, p)
	updated = append(updated, updtemp...)
	updtemp = d.swapDelaunay(t2, p)
	updated = append(updated, updtemp...)
	updtemp = d.swapDelaunay(t3, p)
	updated = append(updated, updtemp...)
	return updated
}

// insertOnEdge inserts a point on an Edge between two triangles. It returns the points
// whose nearest neighbor changed.
func (d *Delaunay) insertOnEdge(p *Point, t *Triangle) []*Point {
	// To increase performance find the points in t, where p1 has the least adjacent triangles and
	// p2 the second least.
	var p1, p2, p3 *Point
	if len(t.A.adjacentTriangles) < len(t.B.adjacentTriangles) {
		if len(t.C.adjacentTriangles) < len(t.A.adjacentTriangles) {
			p1 = t.C
			p2 = t.A
			p3 = t.B
		} else {
			p1 = t.A
			if len(t.C.adjacentTriangles) < len(t.B.adjacentTriangles) {
				p2 = t.C
				p3 = t.B
			} else {
				p2 = t.B
				p3 = t.C
			}
		}
	} else {
		if len(t.C.adjacentTriangles) < len(t.B.adjacentTriangles) {
			p1 = t.C
			p2 = t.B
			p3 = t.A
		} else {
			p1 = t.B
			if len(t.A.adjacentTriangles) < len(t.C.adjacentTriangles) {
				p2 = t.A
				p3 = t.C
			} else {
				p2 = t.C
				p3 = t.A
			}
		}
	}
	// pA1 and pA2 will be the points adjacent to the edge. pO1 and pO2 will be the points opposite to the edge
	var pA1, pA2, pO1, pO2 *Point
	// find second triangle adjacent to edge
	var t2 *Triangle
	// exactly two points in each adjacent triangle to the edge have to be adjacent to that edge.
	found := false
	for _, ta := range p1.adjacentTriangles {
		if !ta.Equals(t) && p.inTriangle(ta) == onEdge {
			found = true
			t2 = ta
			break
		}
	}
	if found {
		pA1 = p1
		found = false
		for _, ta := range p2.adjacentTriangles {
			if ta.Equals(t2) {
				found = true
				break
			}
		}
		if found {
			pA2 = p2
			pO1 = p3
		} else {
			pA2 = p3
			pO1 = p2
		}
	} else {
		pO1 = p1
		for _, ta := range p2.adjacentTriangles {
			if !ta.Equals(t) && p.inTriangle(ta) == onEdge {
				found = true
				t2 = ta
				break
			}
		}
		if !found {
			panic(fmt.Errorf("delaunay: can't find second triangle with edge to %v and %v on edge", t, p))
		}
		pA1 = p2
		pA2 = p3
	}
	switch {
	case !t2.A.Equals(pA1) && !t2.A.Equals(pA2):
		pO2 = t2.A
	case !t2.B.Equals(pA1) && !t2.B.Equals(pA2):
		pO2 = t2.B
	case !t2.C.Equals(pA1) && !t2.C.Equals(pA2):
		pO2 = t2.C
	default:
		panic(fmt.Errorf("delaunay: no point in %v that is not adjacent to the edge of %v", t2, p))
	}
	// form four new triangles
	nt1 := NewTriangle(pA1, p, pO1)
	nt2 := NewTriangle(pA1, p, pO2)
	nt3 := NewTriangle(pA2, p, pO1)
	nt4 := NewTriangle(pA2, p, pO2)
	updated := nt1.add()
	updtemp := nt2.add()
	updated = append(updated, updtemp...)
	updtemp = nt3.add()
	updated = append(updated, updtemp...)
	updtemp = nt4.add()
	updated = append(updated, updtemp...)
	updtemp = t.remove()
	updated = append(updated, updtemp...)
	updtemp = t2.remove()
	updated = append(updated, updtemp...)
	t.children = append(t.children, nt1, nt3)
	t2.children = append(t2.children, nt2, nt4)
	d.triangles = append(d.triangles, nt1, nt2, nt3, nt4)
	// change the edges so it is a valid delaunay triangulation
	updtemp = d.swapDelaunay(nt1, p)
	updated = append(updated, updtemp...)
	updtemp = d.swapDelaunay(nt2, p)
	updated = append(updated, updtemp...)
	updtemp = d.swapDelaunay(nt3, p)
	updated = append(updated, updtemp...)
	updtemp = d.swapDelaunay(nt4, p)
	updated = append(updated, updtemp...)
	return updated
}

// swapDelaunay finds the triangle adjacent to t and opposite to p.
// Then it checks whether p is in the circumcircle. If p is in the circumcircle
// that means that the triangle is not a valid delaunay triangle.
// Therefore the edge in between the two triangles is swapped, creating
// two new triangles that need to be checked.
// It returns the points whose nearest neighbors changed during the
// process.
func (d *Delaunay) swapDelaunay(t *Triangle, p *Point) []*Point {
	// find points in the triangle that are not p
	var p2, p3 *Point
	switch {
	case p.Equals(t.A):
		p2 = t.B
		p3 = t.C
	case p.Equals(t.B):
		p2 = t.C
		p3 = t.A
	case p.Equals(t.C):
		p2 = t.A
		p3 = t.B
	default:
		panic(fmt.Errorf("delaunay: can't find point %v in Triangle %v", p, t))
	}
	// find triangle opposite to p
	var ta *Triangle
loop:
	for _, t1 := range p2.adjacentTriangles {
		for _, t2 := range p3.adjacentTriangles {
			if !t1.Equals(t) && t1.Equals(t2) {
				ta = t1
				break loop
			}
		}

	}
	var updated []*Point
	if ta == nil {
		return updated
	}
	pos := predicates.Incircle(ta.A.x, ta.A.y, ta.B.x, ta.B.y, ta.C.x, ta.C.y, p.x, p.y)
	// swap edges if p is inside the circumcircle of ta
	if pos == predicates.Inside {
		var nt1, nt2 *Triangle
		nt1, nt2, updated = d.swapEdge(t, ta)
		updtemp := d.swapDelaunay(nt1, p)
		updated = append(updated, updtemp...)
		updtemp = d.swapDelaunay(nt2, p)
		updated = append(updated, updtemp...)
	}
	return updated
}

// swapEdge swaps edge between two triangles.
// The edge in the middle of the two triangles is removed and
// an edge between the two opposite points is added.
func (d *Delaunay) swapEdge(t1, t2 *Triangle) (nt1, nt2 *Triangle, updated []*Point) {
	// find points adjacent and opposite to edge
	var adj1, adj2, opp1, opp2 *Point
	switch {
	case !t1.A.Equals(t2.A) && !t1.A.Equals(t2.B) && !t1.A.Equals(t2.C):
		adj1 = t1.A
		opp1 = t1.B
		opp2 = t1.C
	case !t1.B.Equals(t2.A) && !t1.B.Equals(t2.B) && !t1.B.Equals(t2.C):
		adj1 = t1.B
		opp1 = t1.A
		opp2 = t1.C
	case !t1.C.Equals(t2.A) && !t1.C.Equals(t2.B) && !t1.C.Equals(t2.C):
		adj1 = t1.C
		opp1 = t1.B
		opp2 = t1.A
	default:
		panic(fmt.Errorf("delaunay: triangle T1%v is equal to T2%v", t1, t2))
	}
	switch {
	case !t2.A.Equals(t1.A) && !t2.A.Equals(t1.B) && !t2.A.Equals(t1.C):
		adj2 = t2.A
	case !t2.B.Equals(t1.A) && !t2.B.Equals(t1.B) && !t2.B.Equals(t1.C):
		adj2 = t2.B
	case !t2.C.Equals(t1.A) && !t2.C.Equals(t1.B) && !t2.C.Equals(t1.C):
		adj2 = t2.C
	default:
		panic(fmt.Errorf("delaunay: triangle T2%v is equal to T1%v", t2, t1))
	}
	// create two new triangles
	nt1 = NewTriangle(adj1, adj2, opp1)
	nt2 = NewTriangle(adj1, adj2, opp2)
	updated = nt1.add()
	updtemp := nt2.add()
	updated = append(updated, updtemp...)
	updtemp = t1.remove()
	updated = append(updated, updtemp...)
	updtemp = t2.remove()
	updated = append(updated, updtemp...)
	t1.children = append(t1.children, nt1, nt2)
	t2.children = append(t2.children, nt1, nt2)
	d.triangles = append(d.triangles, nt1, nt2)
	return nt1, nt2, updated
}

// retriangulateAndSew uses the re-triangulate and sew method to find the delaunay triangles
// inside the polygon formed by the CCW-ordered points. If k = len(points) then it has a
// worst-time complexity of O(k*log(k)).
func (d *Delaunay) retriangulateAndSew(points []*Point, parents []*Triangle) (updated []*Point) {
	nd := HierarchicalDelaunay()
	// change limits to create a root triangle that's far outside of the original root triangle
	nd.root.A.x = -1 << 35
	nd.root.A.y = -1 << 35
	nd.root.B.x = 1 << 35
	nd.root.B.y = -1 << 35
	nd.root.C.y = 1 << 35
	// make copies of points on polygon and run a delaunay triangulation with them
	// indices of copies are in counter clockwise order, so that with the help of
	// areCounterclockwise it can be determined if a point is inside or outside the polygon.
	// A,B,C are ordered counterclockwise, so if the numbers in A,B,C are counterclockwise it is
	// inside the polygon.
	copies := make([]*Point, len(points))
	for i, p := range points {
		copies[i] = NewPoint(p.x, p.y)
		copies[i].id = i
		nd.Insert(copies[i])
	}
	ts := nd.Triangles()
	triangles := make([]*Triangle, 0, len(ts))
	for _, t := range ts {
		a := t.A.id
		b := t.B.id
		c := t.C.id
		// only keep triangles that are inside the polygon
		// points are inside the polygon if the order of the indices inside the triangle
		// is counterclockwise
		if areCounterclockwise(a, b, c) {
			tr := NewTriangle(points[a], points[b], points[c])
			updtemp := tr.add()
			updated = append(updated, updtemp...)
			triangles = append(triangles, tr)
		}
	}
	d.triangles = append(d.triangles, triangles...)
	for i := range parents {
		parents[i].children = append(parents[i].children, triangles...)
	}
	return updated
}

// areCounterclockwise is a helper function for retriangulateAndSew. It returns
// whether three points are in counterclockwise order.
// Since the points in triangle are ordered counterclockwise and the indices around
// the polygon are ordered counterclockwise checking if the indices of A,B,C
// are counter clockwise is enough.
func areCounterclockwise(a, b, c int) bool {
	if b < c {
		return a < b || c < a
	}
	return a < b && c < a
}
