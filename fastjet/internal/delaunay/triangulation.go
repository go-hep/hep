// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"fmt"
	"math/big"
	"math/rand"
)

var (
	zero = big.NewFloat(0)
)

// Delaunay holds necessary information for the delaunay triangulation
type Delaunay struct {
	// triangles is a slice of all triangles that have been created. It is used to get the final
	// list of triangles in the delaunay triangulation
	triangles triangles
	// root is a triangle that contains all points. It is used as the starting point in the hierarchy
	// to find the triangle that contains a point
	root                   *Triangle
	maxX, minX, maxY, minY float64
	// hierarchy indicates which method to use to find the triangle that contains the point
	useHierarchical bool
	r               *rand.Rand
}

// NewDelaunay creates a delaunay triangulation with the given points
// all points have to be inside the user defined bounds.
// It uses a hierarchy to find the triangle which contains a point.
// It has a worst time complexity of O(nln(n)).
func NewDelaunay(pts []*Point, maxX, minX, maxY, minY float64) *Delaunay {
	// root Triangle is a triangle that contains all points
	dx := maxX - minX
	dy := maxY - minY
	a := NewPoint(minX-8*dx, maxY+10)
	b := NewPoint(maxX+8*dx, maxY+10)
	c := NewPoint(minX+dx/2, minY-8*dy)
	root := NewTriangle(a, b, c)
	d := &Delaunay{
		root:            root,
		maxX:            maxX,
		minX:            minX,
		maxY:            maxY,
		minY:            minY,
		useHierarchical: true,
	}
	d.triangles = make([]*Triangle, 0)
	for _, p := range pts {
		d.Insert(p)
	}
	return d
}

// NewUnboundedDelaunay creates a delaunay triangulation with the given points.
// It uses the remembering stochastic walk method to find the triangle in which p is inserted.
// It has a worst time complexity of O(n^5/3).
func NewUnboundedDelaunay(pts []*Point, r *rand.Rand) *Delaunay {
	// it needs at least one triangle to start with
	if len(pts) < 3 {
		panic(fmt.Errorf("delaunay: not enough points"))
	}
	j := 2
	for ; ; j++ {
		if j >= len(pts) {
			panic(fmt.Errorf("delaunay: all points are in a line"))
		}
		if pts[j].orientation(pts[0], pts[1]).Cmp(zero) != 0 {
			break
		}
	}
	d := &Delaunay{
		useHierarchical: false,
		r:               r,
	}
	d.triangles = make([]*Triangle, 1)
	// create first triangle
	d.triangles[0] = NewTriangle(pts[0], pts[1], pts[j])
	d.triangles[0].isInTriangulation = true
	pts[0].adjacentTriangles = pts[0].adjacentTriangles.append(d.triangles[0])
	pts[1].adjacentTriangles = pts[1].adjacentTriangles.append(d.triangles[0])
	pts[j].adjacentTriangles = pts[j].adjacentTriangles.append(d.triangles[0])
	for i := 2; i < len(pts); i++ {
		if i == j {
			// point was inserted with first triangle
			continue
		}
		d.Insert(pts[i])
	}
	return d
}

// Triangles returns all delaunay Triangles
func (d *Delaunay) Triangles() []*Triangle {
	if !d.useHierarchical {
		return d.triangles.finalize()
	}
	// remove triangles that contain the root points
	rt := make(triangles, len(d.root.A.adjacentTriangles)+len(d.root.B.adjacentTriangles)+len(d.root.C.adjacentTriangles))
	n := copy(rt, d.root.A.adjacentTriangles)
	n += copy(rt[n:], d.root.B.adjacentTriangles)
	copy(rt[n:], d.root.C.adjacentTriangles)
	return d.triangles.finalize(rt...)
}

func (d *Delaunay) Insert(p *Point) {
	p.adjacentTriangles = make(triangles, 0)
	var t *Triangle
	var isOnEdge bool
	if d.useHierarchical {
		t, isOnEdge = findTriangle(d.root, p)
	} else {
		var l, r *Point
		var start *Triangle
		// find a triangle to start the walk
		// use the last changed triangle to start from
		for i := len(d.triangles) - 1; ; i-- {
			if d.triangles[i].isInTriangulation {
				start = d.triangles[i]
				break
			}
		}
		t, isOnEdge, l, r = d.walkTriangle(start, p)
		if l != nil || r != nil {
			// point is on the outside of the current triangulation
			d.addPoint(p, l, r, t)
			return
		}
	}
	if t == nil {
		panic(fmt.Errorf("delaunay: no triangle which contains P%v. Min and Max values must be wrong.", p))
	}
	if isOnEdge {
		d.insertAtEdge(p, t)
	} else {
		d.insertPoint(p, t)
	}
}

// addPoint adds point on the outside. It checks the neighbors of l and r and creates
// triangles between p,l,r and the neighbors if they can reach p without crossing any lines
func (d *Delaunay) addPoint(p, l, r *Point, t *Triangle) {
	// need to find points next to l and r on the border and check if there orientation is less than 0
	outerL, outerR := p.findBorderConnectors(l, r, t)
	// create triangles with all the points that can connect to p
	nts := make(triangles, 1)
	nts[0] = NewTriangle(l, r, p)
	l.adjacentTriangles = l.adjacentTriangles.append(nts[0])
	r.adjacentTriangles = r.adjacentTriangles.append(nts[0])
	for i := 0; i < len(outerL)-1; i++ {
		nt := NewTriangle(outerL[i+1], p, outerL[i])
		nts = append(nts, nt)
		outerL[i].adjacentTriangles = outerL[i].adjacentTriangles.append(nt)
		outerL[i+1].adjacentTriangles = outerL[i+1].adjacentTriangles.append(nt)
	}
	for i := 0; i < len(outerR)-1; i++ {
		nt := NewTriangle(outerR[i+1], outerR[i], p)
		nts = append(nts, nt)
		outerR[i].adjacentTriangles = outerR[i].adjacentTriangles.append(nt)
		outerR[i+1].adjacentTriangles = outerR[i+1].adjacentTriangles.append(nt)
	}
	p.adjacentTriangles = p.adjacentTriangles.append(nts...)
	d.triangles = append(d.triangles, nts...)
	for _, t := range nts {
		t.isInTriangulation = true
	}
	// validate the edges
	for _, t := range nts {
		if t.isInTriangulation {
			d.swapDelaunay(t, p)
		}
	}
}

func (d *Delaunay) Remove(p *Point) {
	// in hierarchical delaunay points on the outside are never removed. The root points stay.
	// in walk delaunay points on the outside are removed and therefore can have less than 3
	// adjacent triangles
	if len(p.adjacentTriangles) < 3 && d.useHierarchical {
		panic(fmt.Errorf("delaunay: can't remove point P%v not enough adjacent triangles", p))
	}
	// remove triangles adjacent to p
	for _, t := range p.adjacentTriangles {
		t.isInTriangulation = false
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
			panic(fmt.Errorf("delaunay: point %v not in adjacent triangle %v", p, t))
		}
	}
	if len(p.adjacentTriangles) == 1 {
		// can't form a new triangle with only one adjacent triangle
		return
	}
	if len(p.adjacentTriangles) == 0 {
		panic(fmt.Errorf("delaunay: no adjacent triangles of %v. Need at least one triangle to remove", p))
	}
	// find points on polygon around the point in counterclockwise order
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
			if d.useHierarchical {
				panic(fmt.Errorf("delaunay: internal error with adjacent triangles for P%v. Can't find counterclockwise neighbor of P%v", p, points[j]))
			}
			// the bound was reached, now the rest of the points are found by going clockwise from the starting point
			// an outer point has one more adjacent points than triangles
			outerpoints := make([]*Point, len(points)+1)
			copy(outerpoints, points)
			pts := p.findRemainingSurrounding(outerpoints, j)
			d.removePoints(pts, nil)
			return
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
	if !d.useHierarchical {
		// check if point is on the outside by checking if there is a triangle that contains
		// the last and the first point found
		found := false
		if len(p.adjacentTriangles) >= 3 {
			for _, t := range p.adjacentTriangles {
				if (t.A == points[0] || t.B == points[0] || t.C == points[0]) &&
					(t.A == points[len(points)-1] || t.B == points[len(points)-1] ||
						t.C == points[len(points)-1]) {
					found = true
					break
				}
			}
		}
		// it is on the outside since the first and the last point don't have a triangle in common
		if !found {
			// an outer point has one more adjacent points than triangles
			outerpoints := make([]*Point, len(points)+1)
			copy(outerpoints, points)
			// check if you can find the remaining point by going counterclockwise from the last point
			for i := 0; ; {
				if i >= len(p.adjacentTriangles) {
					// the bound was reached, now the rest of the points are found by going clockwise from the starting point
					pts := p.findRemainingSurrounding(outerpoints, j)
					d.removePoints(pts, nil)
					return
				}
				// it needs to find the triangle next to k and not k again
				if p.adjacentTriangles[i].Equals(p.adjacentTriangles[k]) {
					i++
					continue
				}
				t = p.adjacentTriangles[i]
				switch {
				case outerpoints[j].Equals(t.A):
					j++
					outerpoints[j] = t.B
					break
				case outerpoints[j].Equals(t.B):
					j++
					outerpoints[j] = t.C
					break
				case outerpoints[j].Equals(t.C):
					j++
					outerpoints[j] = t.A
					break
				}
				i++
			}
			pts := p.findRemainingSurrounding(outerpoints, len(outerpoints)-1)
			d.removePoints(pts, nil)
			return
		}
	}
	d.removePoints(points, p.adjacentTriangles)
}

// removePoints forms a new triangulation inside the points, which form a polygon
func (d *Delaunay) removePoints(points []*Point, parents []*Triangle) {
	// for performance improvement handle points with few adjacent points differently
	/*FIXME make the low degree optimization work
	new changes have been added since attempting to implement this
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

	} else {*/
	// make copies of points on polygon and run a delaunay triangulation with them
	// indices of copies are in counter clockwise order, so that with the help of
	// areCounterclockwise it can be determined if a point is inside or outside the polygon.
	// A,B,C are ordered counterclockwise, so if the numbers in A,B,C are counterclockwise it is
	// inside the polygon.
	copies := make([]*Point, len(points))
	for i, p := range points {
		copies[i] = NewPoint(p.X, p.Y)
		copies[i].id = i
	}
	var dn *Delaunay
	if d.useHierarchical {
		// change limits to create a root triangle that's far outside of the original root triangle
		dx := d.maxX - d.minX
		dy := d.maxY - d.minY
		dn = NewDelaunay(copies, d.maxX+6*dx, d.minX-6*dx, d.maxY+10, d.minY-6*dy)
	} else {
		dn = NewUnboundedDelaunay(copies, d.r)
	}
	ts := dn.Triangles()
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
			tr.isInTriangulation = true
			points[a].adjacentTriangles = points[a].adjacentTriangles.append(tr)
			points[b].adjacentTriangles = points[b].adjacentTriangles.append(tr)
			points[c].adjacentTriangles = points[c].adjacentTriangles.append(tr)
			triangles = append(triangles, tr)
		}
	}
	d.triangles = append(d.triangles, triangles...)
	if d.useHierarchical {
		for i := range parents {
			parents[i].children = append(parents[i].children, triangles...)
		}
	}
	//	}
}

// areCounterclockwise return whether three points are in counterclockwise order.
// Since the points in triangle are ordered counterclockwise and the indices around
// the polygon are ordered counterclockwise checking if the indices of A,B,C
// are counter clockwise
func areCounterclockwise(a, b, c int) bool {
	if b < c {
		return a < b || c < a
	}
	return a < b && c < a
}

// walkTriangle finds the triangle which contains p by using a remembering stochastic walk
func (d *Delaunay) walkTriangle(start *Triangle, p *Point) (t *Triangle, onEdge bool, l *Point, r *Point) {
	found := false
	var previous *Triangle
	for !found {
		found = true
		// k is a random int {0,1,2}
		// it is used to pick a random edge
		// the randomness prevents loops in walks
		var k int
		if d.r == nil {
			k = rand.Intn(3)
		} else {
			k = d.r.Intn(3)
		}
		for i := k; i <= k+2; i++ {
			c := i % 3
			switch c {
			case 0:
				l = start.A
				r = start.B
			case 1:
				l = start.B
				r = start.C
			case 2:
				l = start.C
				r = start.A
			}
			// remembering improvement
			// skip edge if it goes to previous
			inc := 0
			if previous != nil {
				if previous.A == l || previous.A == r {
					inc++
				}
				if previous.B == l || previous.B == r {
					inc++
				}
				if previous.C == l || previous.C == r {
					inc++
				}
			}
			if inc < 2 {
				orient := p.orientation(l, r)
				if orient.Cmp(zero) < 0 {
					// p is on the other side of the line formed by the two points
					// therefore cross the edge
					previous = start
					start = nil
					for _, t1 := range l.adjacentTriangles {
						for _, t2 := range r.adjacentTriangles {
							if t1.Equals(t2) && !t1.Equals(previous) {
								start = t1
								break
							}
						}
					}
					if start == nil {
						// if t is nil the point is outside the current triangulation
						return previous, false, l, r
					}
					found = false
					break
				} else if orient.Cmp(zero) == 0 {
					ab := p.orientation(start.A, start.B)
					bc := p.orientation(start.B, start.C)
					ca := p.orientation(start.C, start.A)
					// p is on the edge if it is on the line formed by the points and if it is in between the 2 other edges
					// in that triangle
					if ab.Cmp(zero) >= 0 && bc.Cmp(zero) >= 0 && ca.Cmp(zero) >= 0 {
						return start, true, nil, nil
					}
				}
			}
		}
	}
	return start, false, nil, nil
}

// findTriangle goes down the hierarchy to find the triangle in which the point is located.
// It returns the triangle and whether a point is on an edge.
func findTriangle(t *Triangle, p *Point) (*Triangle, bool) {
	// get information about the point in respect to its position to the triangle.
	inside, edge := p.inTriangle(t)
	if !inside {
		return nil, false
	}
	// leaf triangle
	if len(t.children) == 0 {
		return t, edge
	}
	for _, tc := range t.children {
		tt, oe := findTriangle(tc, p)
		// if tt is nil then look at other children. Otherwise return the tt.
		if tt != nil {
			return tt, oe
		}
	}
	return nil, false
}

// insertPoint inserts a point inside a triangle
func (d *Delaunay) insertPoint(new *Point, t *Triangle) {
	// form three new triangles
	t1 := NewTriangle(t.A, t.B, new)
	t1.isInTriangulation = true
	t2 := NewTriangle(t.B, t.C, new)
	t2.isInTriangulation = true
	t3 := NewTriangle(t.A, new, t.C)
	t3.isInTriangulation = true
	// adjust the adjacent triangles for all points involved
	new.adjacentTriangles = new.adjacentTriangles.append(t1, t2, t3)
	t.isInTriangulation = false
	t.A.adjacentTriangles = t.A.adjacentTriangles.remove(t)
	t.B.adjacentTriangles = t.B.adjacentTriangles.remove(t)
	t.C.adjacentTriangles = t.C.adjacentTriangles.remove(t)
	t.A.adjacentTriangles = t.A.adjacentTriangles.append(t1, t3)
	t.B.adjacentTriangles = t.B.adjacentTriangles.append(t1, t2)
	t.C.adjacentTriangles = t.C.adjacentTriangles.append(t2, t3)
	if d.useHierarchical {
		t.children = append(t.children, t1, t2, t3)
	}
	d.triangles = append(d.triangles, t1, t2, t3)
	// change the edges so it is a valid delaunay triangulation
	d.swapDelaunay(t1, new)
	d.swapDelaunay(t2, new)
	d.swapDelaunay(t3, new)
}

// insertAtBorderEdge inserts a point on an edge that part of the border.
// This method is only used by the walk method.
func (d *Delaunay) insertAtBorderEdge(new *Point, t *Triangle) {
	ab := new.orientation(t.A, t.B)
	bc := new.orientation(t.B, t.C)
	ca := new.orientation(t.C, t.A)
	var op, adj1, adj2 *Point
	switch {
	case ab.Cmp(zero) == 0:
		op = t.C
		adj1 = t.A
		adj2 = t.B
	case bc.Cmp(zero) == 0:
		op = t.A
		adj1 = t.B
		adj2 = t.C
	case ca.Cmp(zero) == 0:
		op = t.B
		adj1 = t.C
		adj2 = t.A
	default:
		panic(fmt.Errorf("delaunay: %v is not on edge of %t", new, t))
	}
	// form two new triangles
	nt1 := NewTriangle(op, adj1, new)
	nt2 := NewTriangle(op, adj2, new)
	t.isInTriangulation = false
	nt1.isInTriangulation = true
	nt2.isInTriangulation = true
	op.adjacentTriangles = op.adjacentTriangles.remove(t)
	adj1.adjacentTriangles = adj1.adjacentTriangles.remove(t)
	adj2.adjacentTriangles = adj2.adjacentTriangles.remove(t)
	op.adjacentTriangles = op.adjacentTriangles.append(nt1, nt2)
	adj1.adjacentTriangles = adj1.adjacentTriangles.append(nt1)
	adj2.adjacentTriangles = adj2.adjacentTriangles.append(nt2)
	new.adjacentTriangles = new.adjacentTriangles.append(nt1, nt2)
	d.triangles = append(d.triangles, nt1, nt2)
	d.swapDelaunay(nt1, new)
	d.swapDelaunay(nt2, new)
}

// insertAtEdge inserts a point on an edge between two triangles
func (d *Delaunay) insertAtEdge(new *Point, t *Triangle) {
	// find second triangle adjacent to edge
	var t2 *Triangle
	found := false
	for _, t2 = range d.triangles {
		_, edge := new.inTriangle(t2)
		if edge && t2.isInTriangulation && !t2.Equals(t) {
			found = true
			break
		}
	}
	if !found {
		// point is on a border edge
		d.insertAtBorderEdge(new, t)
	}
	// find points opposite and adjacent to the edge
	var p1, p2, pO1, pO2 *Point
	switch {
	case !t.A.Equals(t2.A) && !t.A.Equals(t2.B) && !t.A.Equals(t2.C):
		p1 = t.A
		pO1 = t.B
		pO2 = t.C
	case !t.B.Equals(t2.A) && !t.B.Equals(t2.B) && !t.B.Equals(t2.C):
		p1 = t.B
		pO1 = t.A
		pO2 = t.C
	case !t.C.Equals(t2.A) && !t.C.Equals(t2.B) && !t.C.Equals(t2.C):
		p1 = t.C
		pO1 = t.B
		pO2 = t.A
	default:
		panic(fmt.Errorf("delaunay: triangle T1%v doesn't have points not in T2%v", t, t2))
	}
	switch {
	case !t2.A.Equals(t.A) && !t2.A.Equals(t.B) && !t2.A.Equals(t.C):
		p2 = t2.A
	case !t2.B.Equals(t.A) && !t2.B.Equals(t.B) && !t2.B.Equals(t.C):
		p2 = t2.B
	case !t2.C.Equals(t.A) && !t2.C.Equals(t.B) && !t2.C.Equals(t.C):
		p2 = t2.C
	default:
		panic(fmt.Errorf("delaunay: triangle T2%v doesn't have points not in T1%v", t2, t))
	}
	// form four new triangles
	nt1 := NewTriangle(p1, new, pO1)
	nt1.isInTriangulation = true
	nt2 := NewTriangle(new, p2, pO1)
	nt2.isInTriangulation = true
	nt3 := NewTriangle(p1, new, pO2)
	nt3.isInTriangulation = true
	nt4 := NewTriangle(new, p2, pO2)
	nt4.isInTriangulation = true
	// adjust the adjacent triangles for all points involved
	new.adjacentTriangles = new.adjacentTriangles.append(nt1, nt2, nt3, nt4)
	t.isInTriangulation = false
	t2.isInTriangulation = false
	p1.adjacentTriangles = p1.adjacentTriangles.remove(t)
	p2.adjacentTriangles = p2.adjacentTriangles.remove(t2)
	pO1.adjacentTriangles = pO1.adjacentTriangles.remove(t, t2)
	pO2.adjacentTriangles = pO2.adjacentTriangles.remove(t, t2)
	p1.adjacentTriangles = p1.adjacentTriangles.append(nt1, nt3)
	p2.adjacentTriangles = p2.adjacentTriangles.append(nt2, nt4)
	pO1.adjacentTriangles = pO1.adjacentTriangles.append(nt1, nt2)
	pO2.adjacentTriangles = pO2.adjacentTriangles.append(nt3, nt4)
	if d.useHierarchical {
		t.children = append(t.children, nt1, nt3)
		t2.children = append(t2.children, nt2, nt4)
	}
	d.triangles = append(d.triangles, nt1, nt2, nt3, nt4)
	// change the edges so it is a valid delaunay triangulation
	d.swapDelaunay(nt1, new)
	d.swapDelaunay(nt2, new)
	d.swapDelaunay(nt3, new)
	d.swapDelaunay(nt4, new)
}

// swapDelaunay finds the triangle adjacent to t and opposite to p.
// Then it checks whether p is in the circumcircle. If p is in the circumcircle
// that means that the triangle is not a valid delaunay triangle.
// Therefore the edge in between the two triangles is flipped, creating
// two new triangles that need to be checked.
func (d *Delaunay) swapDelaunay(t *Triangle, p *Point) {
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
		panic(fmt.Errorf("delaunay: can't find point P%v in Triangle T%v", p, t))
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
	// flip edges if p is inside circumcircle of ta
	if ta != nil && ta.inCircumcircle(p) {
		nt1, nt2 := d.swapEdge(t, ta)
		d.swapDelaunay(nt1, p)
		d.swapDelaunay(nt2, p)
	}
}

// swapEdge flips edge between two triangles.
// The edge in the middle of the two triangles is removed and
// an edge between the two opposite points is added
func (d *Delaunay) swapEdge(t1, t2 *Triangle) (nt1, nt2 *Triangle) {
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
	nt1.isInTriangulation = true
	nt2 = NewTriangle(adj1, adj2, opp2)
	nt2.isInTriangulation = true
	t1.isInTriangulation = false
	t2.isInTriangulation = false
	// update adjacent lists
	adj1.adjacentTriangles = adj1.adjacentTriangles.remove(t1)
	adj2.adjacentTriangles = adj2.adjacentTriangles.remove(t2)
	opp1.adjacentTriangles = opp1.adjacentTriangles.remove(t1, t2)
	opp2.adjacentTriangles = opp2.adjacentTriangles.remove(t1, t2)
	adj1.adjacentTriangles = adj1.adjacentTriangles.append(nt1, nt2)
	adj2.adjacentTriangles = adj2.adjacentTriangles.append(nt1, nt2)
	opp1.adjacentTriangles = opp1.adjacentTriangles.append(nt1)
	opp2.adjacentTriangles = opp2.adjacentTriangles.append(nt2)
	if d.useHierarchical {
		t1.children = append(t1.children, nt1, nt2)
		t2.children = append(t2.children, nt1, nt2)
	}
	d.triangles = append(d.triangles, nt1, nt2)
	return nt1, nt2
}
