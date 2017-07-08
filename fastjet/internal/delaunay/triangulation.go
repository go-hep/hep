// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
)

// Delaunay holds necessary information for the
// delaunay triangulation
type Delaunay struct {
	triangles              triangles // all triangles created
	root                   *Triangle // triangle that contains all points. Used to find a triangle
	maxX, minX, maxY, minY float64
}

// NewDelaunay creates a delaunay triangulation with the given points
// all points have to be inside the user defined bounds
func NewDelaunay(points []*Point, maxX float64, minX float64, maxY float64, minY float64) *Delaunay {
	// root Triangle is a triangle that contains all points
	dx := maxX - minX
	dy := maxY - minY
	a := NewPoint(minX-8*dx, maxY+10, -1)
	b := NewPoint(maxX+8*dx, maxY+10, -1)
	c := NewPoint(minX+dx/2, minY-8*dy, -1)
	root := NewTriangle(a, b, c)
	d := &Delaunay{
		root: root,
		maxX: maxX,
		minX: minX,
		maxY: maxY,
		minY: minY,
	}
	d.triangles = make([]*Triangle, 0)
	for _, p := range points {
		d.InsertPoint(p)
	}
	return d
}

// Triangle returns all delaunay Triangles
func (d *Delaunay) Triangles() []*Triangle {
	// remove triangles that contain the root points
	rt := make(triangles, len(d.root.A.adjacentTriangles), len(d.root.A.adjacentTriangles)+len(d.root.B.adjacentTriangles)+len(d.root.C.adjacentTriangles))
	copy(rt, d.root.A.adjacentTriangles)
	rt = append(rt, d.root.B.adjacentTriangles...)
	rt = append(rt, d.root.C.adjacentTriangles...)
	return d.triangles.finalize(rt...)
}

func (d *Delaunay) InsertPoint(p *Point) {
	p.adjacentTriangles = make(triangles, 0)
	t, onE := findTriangle(d.root, p)
	if t == nil {
		// should only happen when user gives wrong max and min values
		panic(fmt.Errorf("delaunay: no triangle which contains P%s", p))
	}
	if onE {
		d.insertPonE(t, p)
	} else {
		d.insertP(t, p)
	}
}

func (d *Delaunay) RemovePoint(p *Point) {
	if len(p.adjacentTriangles) < 3 {
		panic(fmt.Errorf("delaunay: can't remove point P%s not enough adjacent triangles", p))
	}
	// remove triangles adjacent to p
	for _, t := range p.adjacentTriangles {
		switch {
		case p.Equals(t.A):
			t.B.adjacentTriangles = t.B.adjacentTriangles.remove(t)
			t.C.adjacentTriangles = t.C.adjacentTriangles.remove(t)
		case p.Equals(t.B):
			t.A.adjacentTriangles = t.A.adjacentTriangles.remove(t)
			t.C.adjacentTriangles = t.C.adjacentTriangles.remove(t)
		case p.Equals(t.C):
			t.A.adjacentTriangles = t.A.adjacentTriangles.remove(t)
			t.B.adjacentTriangles = t.B.adjacentTriangles.remove(t)
		default:
			panic(fmt.Errorf("delaunay: internal error with adjacent triangles for P%s and T%s", p, t))
		}
	}
	// find points on polygon around the point in counterclockwise order
	points := make([]*Point, len(p.adjacentTriangles))
	t := p.adjacentTriangles[0]
	j := 1
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
		panic(fmt.Errorf("delaunay: internal error with adjacent triangles for P%s and T%s", p, t))
	}
	for i := 0; j < len(points)-1; {
		if i >= len(p.adjacentTriangles) {
			panic(fmt.Errorf("delaunay: internal error with adjacent triangles for P%s. Can't find counterclockwise neighbor of P%s", p, points[j]))
		}
		// k is the index of the previous triangle
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
	d.removeP(points, p.adjacentTriangles)
}

// forms a new triangulation inside the points, which form a polygon
func (d *Delaunay) removeP(points []*Point, parents []*Triangle) {
	// for performance improvement handle points with few adjacent points differently
	if len(points) == 3 {
		// polygon is already a delaunay triangle
		nt := NewTriangle(points[0], points[1], points[2])
		points[0].adjacentTriangles = points[0].adjacentTriangles.appendT(nt)
		points[1].adjacentTriangles = points[1].adjacentTriangles.appendT(nt)
		points[2].adjacentTriangles = points[2].adjacentTriangles.appendT(nt)
		d.triangles = append(d.triangles, nt)
		for i := range parents {
			parents[i].children = append(parents[i].children, nt)
		}
		/*	FIXME make the low degree optimization work
			} else if len(points) == 4 {
					// only two possible edges, so one incircle test can determine the valid edge
					nt1 := NewTriangle(points[0], points[1], points[2])
					var nt2 *Triangle
					if nt1.inCCircle(points[3]) {
						nt1 = NewTriangle(points[3], points[0], points[1])
						nt2 = NewTriangle(points[1], points[2], points[3])
						points[0].adjT = appendT(points[0].adjT, nt1)
						points[1].adjT = appendT(points[1].adjT, nt1, nt2)
						points[2].adjT = appendT(points[2].adjT, nt2)
						points[3].adjT = appendT(points[3].adjT, nt1, nt2)
					} else {
						nt2 = NewTriangle(points[3], points[0], points[2])
						points[0].adjT = appendT(points[0].adjT, nt1, nt2)
						points[1].adjT = appendT(points[1].adjT, nt1)
						points[2].adjT = appendT(points[2].adjT, nt1, nt2)
						points[3].adjT = appendT(points[3].adjT, nt2)
					}
					d.triangles = append(d.triangles, nt1, nt2)
					for i := range parents {
						parents[i].children = append(parents[i].children, nt1, nt2)
					}
				} else if len(points) == 5 {
					// use a decision tree to determine the correct triangles
					nt1 := NewTriangle(points[0], points[1], points[2])
					var nt2 *Triangle
					var nt3 *Triangle
					if nt1.inCCircle(points[3]) { // 3 in 012
						nt1 = NewTriangle(points[0], points[1], points[3])
						if nt1.inCCircle(points[4]) { // 4 in 013
							nt1 = NewTriangle(points[1], points[2], points[3])
							if nt1.inCCircle(points[4]) { // 4 in 123
								nt1 = NewTriangle(points[4], points[2], points[3])
								nt2 = NewTriangle(points[4], points[1], points[2])
								nt3 = NewTriangle(points[4], points[0], points[1])
								points[0].adjT = appendT(points[0].adjT, nt3)
								points[1].adjT = appendT(points[1].adjT, nt2, nt3)
								points[2].adjT = appendT(points[2].adjT, nt1, nt2)
								points[3].adjT = appendT(points[3].adjT, nt1)
								points[4].adjT = appendT(points[4].adjT, nt1, nt2, nt3)
							} else {
								// nt1 123
								nt2 = NewTriangle(points[4], points[1], points[3])
								nt3 = NewTriangle(points[4], points[0], points[1])
								points[0].adjT = appendT(points[0].adjT, nt3)
								points[1].adjT = appendT(points[1].adjT, nt1, nt2, nt3)
								points[2].adjT = appendT(points[2].adjT, nt1)
								points[3].adjT = appendT(points[3].adjT, nt1, nt2)
								points[4].adjT = appendT(points[4].adjT, nt2, nt3)
							}
						} else {
							// nt1 013
							nt2 = NewTriangle(points[3], points[1], points[2])
							nt3 = NewTriangle(points[3], points[4], points[0])
							points[0].adjT = appendT(points[0].adjT, nt1, nt3)
							points[1].adjT = appendT(points[1].adjT, nt1, nt2)
							points[2].adjT = appendT(points[2].adjT, nt2)
							points[3].adjT = appendT(points[3].adjT, nt1, nt2, nt3)
							points[4].adjT = appendT(points[4].adjT, nt3)
						}
					} else {
						nt2 = NewTriangle(points[0], points[2], points[3])
						if nt2.inCCircle(points[4]) { // 4 in 023
							if nt1.inCCircle(points[4]) { // 4 in 012
								nt1 = NewTriangle(points[4], points[0], points[1])
								nt2 = NewTriangle(points[4], points[1], points[2])
								nt3 = NewTriangle(points[4], points[2], points[3])
								points[0].adjT = appendT(points[0].adjT, nt1)
								points[1].adjT = appendT(points[1].adjT, nt1, nt2)
								points[2].adjT = appendT(points[2].adjT, nt2, nt3)
								points[3].adjT = appendT(points[3].adjT, nt3)
								points[4].adjT = appendT(points[4].adjT, nt1, nt2, nt3)
							} else {
								// nt1 012
								nt2 = NewTriangle(points[0], points[2], points[4])
								nt3 = NewTriangle(points[4], points[2], points[3])
								points[0].adjT = appendT(points[0].adjT, nt1, nt2)
								points[1].adjT = appendT(points[1].adjT, nt1)
								points[2].adjT = appendT(points[2].adjT, nt1, nt2, nt3)
								points[3].adjT = appendT(points[3].adjT, nt3)
								points[4].adjT = appendT(points[4].adjT, nt2, nt3)
							}
						} else {
							// nt1 012
							// nt2 023
							nt3 = NewTriangle(points[0], points[3], points[4])
							points[0].adjT = appendT(points[0].adjT, nt1, nt2, nt3)
							points[1].adjT = appendT(points[1].adjT, nt1)
							points[2].adjT = appendT(points[2].adjT, nt1, nt2)
							points[3].adjT = appendT(points[3].adjT, nt2, nt3)
							points[4].adjT = appendT(points[4].adjT, nt3)
						}
					}
					d.triangles = append(d.triangles, nt1, nt2, nt3)
					for i := range parents {
						parents[i].children = append(parents[i].children, nt1, nt2, nt3)
					}
				} else if len(points) == 6 {
					// use a decision tree to determine the correct triangles
					nt1 := NewTriangle(points[2], points[3], points[0])
					var nt2 *Triangle
					var nt3 *Triangle
					var nt4 *Triangle
					if nt1.inCCircle(points[1]) { // 1 in 230
						nt1 = NewTriangle(points[2], points[3], points[5])
						if nt1.inCCircle(points[4]) { // 4 in 235
							nt1 = NewTriangle(points[2], points[3], points[4])
							if nt1.inCCircle(points[1]) { // 1 in 234
								nt1 = NewTriangle(points[1], points[2], points[3])
								nt2 = NewTriangle(points[0], points[1], points[3])
								if nt2.inCCircle(points[4]) { // 4 in 013
									nt2 = NewTriangle(points[1], points[3], points[4])
									nt3 = NewTriangle(points[0], points[1], points[4])
									if nt3.inCCircle(points[5]) { // 5 in 014
										nt3 = NewTriangle(points[1], points[4], points[5])
										nt4 = NewTriangle(points[1], points[5], points[0])
									} else {
										nt4 = NewTriangle(points[0], points[4], points[5])
									}
								} else {
									nt3 = NewTriangle(points[0], points[3], points[4])
									nt4 = NewTriangle(points[0], points[4], points[5])
								}
							} else {
								nt2 = NewTriangle(points[1], points[2], points[4])
								if nt2.inCCircle(points[5]) { // 5 in 124
									nt2 = NewTriangle(points[5], points[2], points[4])
									nt3 = NewTriangle(points[5], points[1], points[2])
									nt4 = NewTriangle(points[5], points[0], points[1])
								} else {
									nt3 = NewTriangle(points[0], points[1], points[4])
									if nt3.inCCircle(points[5]) { // 5 in 014
										nt3 = NewTriangle(points[5], points[1], points[4])
										nt4 = NewTriangle(points[5], points[0], points[1])
									} else {
										nt4 = NewTriangle(points[5], points[0], points[4])
									}
								}
							}
						} else {
							nt1 = NewTriangle(points[2], points[3], points[5])
							if nt1.inCCircle(points[1]) { // 1 in 235
								nt1 = NewTriangle(points[1], points[2], points[3])
								nt2 = NewTriangle(points[3], points[4], points[5])
								if nt2.inCCircle(points[1]) { // 1 in 345
									nt2 = NewTriangle(points[0], points[1], points[3])
									if nt2.inCCircle(points[4]) { // 4 in 013
										nt2 = NewTriangle(points[1], points[3], points[4])
										nt3 = NewTriangle(points[0], points[1], points[4])
										if nt3.inCCircle(points[5]) { // 5 in 014
											nt3 = NewTriangle(points[1], points[4], points[5])
											nt4 = NewTriangle(points[1], points[5], points[0])
										} else {
											nt4 = NewTriangle(points[0], points[4], points[5])
										}
									} else {
										nt3 = NewTriangle(points[0], points[3], points[4])
										nt4 = NewTriangle(points[0], points[4], points[5])
									}
								} else {
									nt3 = NewTriangle(points[0], points[1], points[3])
									if nt3.inCCircle(points[5]) { // 5 in 013
										nt3 = NewTriangle(points[0], points[1], points[5])
										nt4 = NewTriangle(points[5], points[1], points[3])
									} else {
										nt4 = NewTriangle(points[5], points[0], points[3])
										if nt4.inCCircle(points[4]) { // 4 in 503
											nt2 = NewTriangle(points[3], points[4], points[0])
											nt4 = NewTriangle(points[0], points[4], points[5])
										} else {

										}
									}
								}
							} else {
								nt2 = NewTriangle(points[5], points[3], points[4])
								nt3 = NewTriangle(points[5], points[1], points[2])
								nt4 = NewTriangle(points[5], points[0], points[1])
							}
						}
					} else {
						nt1 = NewTriangle(points[2], points[3], points[5])
						if nt1.inCCircle(points[4]) { // 4 in 235
							nt1 = NewTriangle(points[2], points[3], points[0])
							if nt1.inCCircle(points[4]) { // 4 in 230
								nt1 = NewTriangle(points[4], points[2], points[3])
								nt2 = NewTriangle(points[0], points[1], points[2])
								if nt2.inCCircle(points[4]) { // 4 in 012
									nt2 = NewTriangle(points[1], points[2], points[5])
									if nt2.inCCircle(points[4]) { // 4 in 125
										nt2 = NewTriangle(points[4], points[1], points[2])
										nt3 = NewTriangle(points[0], points[1], points[5])
										if nt3.inCCircle(points[4]) { // 4 in 015
											nt3 = NewTriangle(points[0], points[1], points[4])
											nt4 = NewTriangle(points[0], points[4], points[5])
										} else {
											nt4 = NewTriangle(points[5], points[1], points[4])
										}
									} else {
										nt3 = NewTriangle(points[5], points[2], points[1])
										nt4 = NewTriangle(points[0], points[1], points[5])
									}
								} else {
									nt2 = NewTriangle(points[5], points[0], points[2])
									if nt2.inCCircle(points[4]) { // 4 in 502
										nt2 = NewTriangle(points[0], points[1], points[2])
										nt3 = NewTriangle(points[0], points[4], points[5])
										nt4 = NewTriangle(points[0], points[2], points[4])
									} else {
										nt3 = NewTriangle(points[5], points[2], points[4])
										nt4 = NewTriangle(points[0], points[1], points[2])

										if nt4.inCCircle(points[5]) { // 5 in 012
											nt2 = NewTriangle(points[5], points[1], points[2])
											nt4 = NewTriangle(points[0], points[1], points[5])
										} else {

										}
									}
								}
							} else {
								nt2 = NewTriangle(points[0], points[1], points[2])
								nt3 = NewTriangle(points[0], points[3], points[4])
								nt4 = NewTriangle(points[0], points[4], points[5])
							}
						} else {
							nt1 = NewTriangle(points[2], points[3], points[0])
							if nt1.inCCircle(points[5]) { // 5 in 230
								nt1 = NewTriangle(points[5], points[3], points[4])
								nt2 = NewTriangle(points[5], points[2], points[3])
								nt3 = NewTriangle(points[0], points[1], points[2])
								if nt3.inCCircle(points[5]) { // 5 in 012
									nt3 = NewTriangle(points[5], points[0], points[1])
									nt4 = NewTriangle(points[5], points[1], points[2])
								} else {
									nt4 = NewTriangle(points[5], points[0], points[2])
								}
							} else {
								nt2 = NewTriangle(points[0], points[1], points[2])
								nt3 = NewTriangle(points[5], points[0], points[3])
								if nt3.inCCircle(points[4]) { // 4 in 503
									nt3 = NewTriangle(points[0], points[3], points[4])
									nt4 = NewTriangle(points[0], points[4], points[5])
								} else {
									nt4 = NewTriangle(points[5], points[3], points[4])
								}
							}
						}
					}
					nt1.A.adjT = appendT(nt1.A.adjT, nt1)
					nt1.B.adjT = appendT(nt1.B.adjT, nt1)
					nt1.C.adjT = appendT(nt1.C.adjT, nt1)
					nt2.A.adjT = appendT(nt2.A.adjT, nt2)
					nt2.B.adjT = appendT(nt2.B.adjT, nt2)
					nt2.C.adjT = appendT(nt2.C.adjT, nt2)
					nt3.A.adjT = appendT(nt3.A.adjT, nt3)
					nt3.B.adjT = appendT(nt3.B.adjT, nt3)
					nt3.C.adjT = appendT(nt3.C.adjT, nt3)
					nt4.A.adjT = appendT(nt4.A.adjT, nt4)
					nt4.B.adjT = appendT(nt4.B.adjT, nt4)
					nt4.C.adjT = appendT(nt4.C.adjT, nt4)
					d.triangles = append(d.triangles, nt1, nt2, nt3, nt4)
					for i := range parents {
						parents[i].children = append(parents[i].children, nt1, nt2, nt3, nt4)
					}
		*/
	} else {
		// make copies of points on polygon and run a delaunay triangulation with them
		// indices of copies are in counter clockwise order, so that with the help of
		// areCounterclockwise it can be determined if a point is inside or outside the polygon
		copies := make([]*Point, len(points))
		for i, p := range points {
			copies[i] = NewPoint(p.X, p.Y, i)
		}
		// change limits in case root points are part of the polygon
		dx := d.maxX - d.minX
		dy := d.maxY - d.minY
		dn := NewDelaunay(copies, d.maxX+6*dx, d.minX-6*dx, d.maxY+10, d.minY-6*dy)
		ts := dn.Triangles()
		triangles := make([]*Triangle, 0, len(ts))
		for _, t := range ts {
			a := t.A.ID
			b := t.B.ID
			c := t.C.ID
			// only keep triangles that are inside the polygon
			// points are inside the triangle if the order of the indices inside the triangle
			// is counter clockwise
			if areCounterclockwise(a, b, c) {
				tr := NewTriangle(points[a], points[b], points[c])
				points[a].adjacentTriangles = points[a].adjacentTriangles.appendT(tr)
				points[b].adjacentTriangles = points[b].adjacentTriangles.appendT(tr)
				points[c].adjacentTriangles = points[c].adjacentTriangles.appendT(tr)
				triangles = append(triangles, tr)
			}
		}
		d.triangles = append(d.triangles, triangles[:]...)
		for i := range parents {
			parents[i].children = append(parents[i].children, triangles[:]...)
		}
	}
}

// since the points in triangle are ordered counterclockwise
// and the indices around the polygon are ordered counterclockwise
// checking if the indices of A,B,C are counter clockwise
func areCounterclockwise(a, b, c int) bool {
	if b < c {
		return a < b || c < a
	}
	return a < b && c < a
}

// findTriangle goes down the hierarchy to find the triangle in which the point is located.
// it returns true if the point is on an edge
func findTriangle(t *Triangle, p *Point) (*Triangle, bool) {
	inside, edge := p.inTriangle(t)
	if !inside {
		return nil, false
	}
	if len(t.children) == 0 {
		return t, edge
	}
	for _, tc := range t.children {
		tt, oe := findTriangle(tc, p)
		if tt != nil {
			return tt, oe
		}
	}
	return nil, false
}

// insertP inserts a point inside a triangle
func (d *Delaunay) insertP(t *Triangle, p *Point) {
	// form three new triangles
	t1 := NewTriangle(t.A, t.B, p)
	t2 := NewTriangle(t.B, t.C, p)
	t3 := NewTriangle(t.A, p, t.C)
	// adjust the adjacent triangles for all points involved
	p.adjacentTriangles = p.adjacentTriangles.appendT(t1, t2, t3)
	t.A.adjacentTriangles = t.A.adjacentTriangles.remove(t)
	t.B.adjacentTriangles = t.B.adjacentTriangles.remove(t)
	t.C.adjacentTriangles = t.C.adjacentTriangles.remove(t)
	t.A.adjacentTriangles = t.A.adjacentTriangles.appendT(t1, t3)
	t.B.adjacentTriangles = t.B.adjacentTriangles.appendT(t1, t2)
	t.C.adjacentTriangles = t.C.adjacentTriangles.appendT(t2, t3)
	t.children = append(t.children, t1, t2, t3)
	d.triangles = append(d.triangles, t1, t2, t3)
	// change the edges so it is a valid delaunay triangulation
	d.validateEdge(t1, p)
	d.validateEdge(t2, p)
	d.validateEdge(t3, p)
}

// insertPonE inserts a point on an edge
func (d *Delaunay) insertPonE(t1 *Triangle, p *Point) {
	// find second triangle adjacent to edge
	var t2 *Triangle
	for _, t2 = range d.triangles {
		_, edge := p.inTriangle(t2)
		if edge && !t2.Equals(t1) {
			break
		}
	}
	// find points opposite and adjacent to the edge
	var p1, p2, pO1, pO2 *Point
	switch {
	case !t1.A.Equals(t2.A) && !t1.A.Equals(t2.B) && !t1.A.Equals(t2.C):
		p1 = t1.A
		pO1 = t1.B
		pO2 = t1.C
	case !t1.B.Equals(t2.A) && !t1.B.Equals(t2.B) && !t1.B.Equals(t2.C):
		p1 = t1.B
		pO1 = t1.A
		pO2 = t1.C
	case !t1.C.Equals(t2.A) && !t1.C.Equals(t2.B) && !t1.C.Equals(t2.C):
		p1 = t1.C
		pO1 = t1.B
		pO2 = t1.A
	default:
		panic(fmt.Errorf("delaunay: triangle T1%s doesn't have points not in T2%s", t1, t2))
	}
	switch {
	case !t2.A.Equals(t1.A) && !t2.A.Equals(t1.B) && !t2.A.Equals(t1.C):
		p2 = t2.A
	case !t2.B.Equals(t1.A) && !t2.B.Equals(t1.B) && !t2.B.Equals(t1.C):
		p2 = t2.B
	case !t2.C.Equals(t1.A) && !t2.C.Equals(t1.B) && !t2.C.Equals(t1.C):
		p2 = t2.C
	default:
		panic(fmt.Errorf("delaunay: triangle T2%s doesn't have points not in T1%s", t2, t1))
	}
	// form four new triangles
	nt1 := NewTriangle(p1, p, pO1)
	nt2 := NewTriangle(p, p2, pO1)
	nt3 := NewTriangle(p1, p, pO2)
	nt4 := NewTriangle(p, p2, pO2)
	// adjust the adjacent triangles for all points involved
	p.adjacentTriangles = p.adjacentTriangles.appendT(nt1, nt2, nt3, nt4)
	p1.adjacentTriangles = p1.adjacentTriangles.remove(t1)
	p2.adjacentTriangles = p2.adjacentTriangles.remove(t2)
	pO1.adjacentTriangles = pO1.adjacentTriangles.remove(t1, t2)
	pO2.adjacentTriangles = pO2.adjacentTriangles.remove(t1, t2)
	p1.adjacentTriangles = p1.adjacentTriangles.appendT(nt1, nt3)
	p2.adjacentTriangles = p2.adjacentTriangles.appendT(nt2, nt4)
	pO1.adjacentTriangles = pO1.adjacentTriangles.appendT(nt1, nt2)
	pO2.adjacentTriangles = pO2.adjacentTriangles.appendT(nt3, nt4)
	t1.children = append(t1.children, nt1, nt3)
	t2.children = append(t2.children, nt2, nt4)
	d.triangles = append(d.triangles, nt1, nt2, nt3, nt4)
	// change the edges so it is a valid delaunay triangulation
	d.validateEdge(nt1, p)
	d.validateEdge(nt2, p)
	d.validateEdge(nt3, p)
	d.validateEdge(nt4, p)
}

// validateEdge turns triangle into a valid delaunay triangle
func (d *Delaunay) validateEdge(t *Triangle, p *Point) {
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
		panic(fmt.Errorf("fastjet: delaunay can't find point P%s in Triangle T%s", p, t))
	}
	// find triangle opposite to p
	var ta *Triangle
	for _, t1 := range p2.adjacentTriangles {
		for _, t2 := range p3.adjacentTriangles {
			if !t1.Equals(t) && t1.Equals(t2) {
				ta = t1
			}
		}

	}
	// flip edges if ta is inside circumcircle of t
	if ta != nil && ta.inCircumcircle(p) {
		nt1, nt2 := d.flip(t, ta)
		d.validateEdge(nt1, p)
		d.validateEdge(nt2, p)
	}
}

// flip flips edge between two triangles
func (d *Delaunay) flip(t1 *Triangle, t2 *Triangle) (*Triangle, *Triangle) {
	// find points adjacent and opposite to edge
	var p1, p2, pO1, pO2 *Point
	switch {
	case !t1.A.Equals(t2.A) && !t1.A.Equals(t2.B) && !t1.A.Equals(t2.C):
		p1 = t1.A
		pO1 = t1.B
		pO2 = t1.C
	case !t1.B.Equals(t2.A) && !t1.B.Equals(t2.B) && !t1.B.Equals(t2.C):
		p1 = t1.B
		pO1 = t1.A
		pO2 = t1.C
	case !t1.C.Equals(t2.A) && !t1.C.Equals(t2.B) && !t1.C.Equals(t2.C):
		p1 = t1.C
		pO1 = t1.B
		pO2 = t1.A
	default:
		panic(fmt.Errorf("delaunay: triangle T1%s doesn't have points not in T2%s", t1, t2))
	}
	switch {
	case !t2.A.Equals(t1.A) && !t2.A.Equals(t1.B) && !t2.A.Equals(t1.C):
		p2 = t2.A
	case !t2.B.Equals(t1.A) && !t2.B.Equals(t1.B) && !t2.B.Equals(t1.C):
		p2 = t2.B
	case !t2.C.Equals(t1.A) && !t2.C.Equals(t1.B) && !t2.C.Equals(t1.C):
		p2 = t2.C
	default:
		panic(fmt.Errorf("delaunay: triangle T2%s doesn't have points not in T1%s", t2, t1))
	}
	// create two new triangles
	nt1 := NewTriangle(p1, p2, pO1)
	nt2 := NewTriangle(p1, p2, pO2)
	// update adjacent lists
	p1.adjacentTriangles = p1.adjacentTriangles.remove(t1)
	p2.adjacentTriangles = p2.adjacentTriangles.remove(t2)
	pO1.adjacentTriangles = pO1.adjacentTriangles.remove(t1, t2)
	pO2.adjacentTriangles = pO2.adjacentTriangles.remove(t1, t2)
	p1.adjacentTriangles = p1.adjacentTriangles.appendT(nt1, nt2)
	p2.adjacentTriangles = p2.adjacentTriangles.appendT(nt1, nt2)
	pO1.adjacentTriangles = pO1.adjacentTriangles.appendT(nt1)
	pO2.adjacentTriangles = pO2.adjacentTriangles.appendT(nt2)
	t1.children = append(t1.children, nt1, nt2)
	t2.children = append(t2.children, nt1, nt2)
	d.triangles = append(d.triangles, nt1, nt2)
	return nt1, nt2
}
