// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestPointEquals(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		want           bool
	}{
		{5, 5, 5, 5, true},
		{0.5, 4, 1, 4, false},
		{-3, 20, 20, -3, false},
		{2.5, 4, 2.5, 4, true},
	}
	for _, test := range tests {
		p1 := NewPoint(test.x1, test.y1)
		p2 := NewPoint(test.x2, test.y2)
		got := p1.Equals(p2)
		if got != test.want {
			t.Errorf("%v == %v, got = %v, want = %v", p1, p2, got, test.want)
		}
	}
}

func TestDistance(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2 float64
		want           float64
	}{
		{1, 1, 3, 3, 8},
		{0, 0, 0, 5, 25},
	}
	for _, test := range tests {
		p1 := NewPoint(test.x1, test.y1)
		p2 := NewPoint(test.x2, test.y2)
		got := p1.distance(p2)
		if got != test.want {
			t.Errorf("Distance between %v and %v, got = %v, want = %v", p1, p2, got, test.want)
		}
	}
}

func TestFindNearest(t *testing.T) {
	const tol float64 = 0.001
	tests := []struct {
		p                 *Point
		adjacentTriangles []*Triangle
		wantDist          float64
		wantNeighbor      *Point
	}{
		{
			NewPoint(0, 0),
			[]*Triangle{
				NewTriangle(NewPoint(0, 0), NewPoint(2, 0), NewPoint(0, 5)),
				NewTriangle(NewPoint(3, 3), NewPoint(0, 0), NewPoint(0, 5)),
			},
			2,
			NewPoint(2, 0),
		},
	}
	for _, test := range tests {
		test.p.adjacentTriangles = test.adjacentTriangles
		test.p.findNearest()
		gotNeighbor, gotDistance := test.p.NearestNeighbor()
		if !floats.EqualWithinAbs(gotDistance, test.wantDist, tol) || !test.wantNeighbor.Equals(gotNeighbor) {
			t.Errorf("Nearest Neighbor for P%v with adjacent triangles %v, \n gotDist = %v, wantDist = %v, gotNeighbor = %v, wantNeighbor = %v",
				test.p, test.adjacentTriangles, gotDistance, test.wantDist, gotNeighbor, test.wantNeighbor)
		}
	}
}

func TestFindSecondNearest(t *testing.T) {
	tests := []struct {
		p                 *Point
		adjacentTriangles []*Triangle
		wantDist          float64
		wantNeighbor      *Point
	}{
		{
			NewPoint(0, 0),
			[]*Triangle{
				NewTriangle(NewPoint(0, 0), NewPoint(2, 0), NewPoint(0, 5)),
				NewTriangle(NewPoint(3, 3), NewPoint(0, 0), NewPoint(0, 5)),
			},
			math.Sqrt(18),
			NewPoint(3, 3),
		},
		{
			// with a duplicate point for (0,0)
			NewPoint(0, 0),
			[]*Triangle{
				NewTriangle(NewPoint(0, 0), NewPoint(2, 0), NewPoint(0, 5)),
				NewTriangle(NewPoint(3, 3), NewPoint(0, 0), NewPoint(0, 5)),
			},
			2,
			NewPoint(2, 0),
		},
		{
			// with a duplicate point for (2,0)
			NewPoint(0, 0),
			[]*Triangle{
				NewTriangle(NewPoint(0, 0), NewPoint(2, 0), NewPoint(0, 5)),
				NewTriangle(NewPoint(3, 3), NewPoint(0, 0), NewPoint(0, 5)),
			},
			2,
			NewPoint(2, 0),
		},
		{
			// with a duplicate point for (2,0)
			NewPoint(0, 0),
			[]*Triangle{
				NewTriangle(NewPoint(0, 5), NewPoint(0, 0), NewPoint(2, 0)),
				NewTriangle(NewPoint(3, 3), NewPoint(0, 0), NewPoint(0, 5)),
			},
			2,
			NewPoint(2, 0),
		},
	}
	// set up the duplicate for (0,0) in test[1]
	tests[1].p.dist2 = 0
	tests[1].p.nearest = NewPoint(0, 0)
	// set up a duplicate for (2,0) in test[2]
	tests[2].adjacentTriangles[0].B.dist2 = 0
	tests[2].adjacentTriangles[0].B.nearest = NewPoint(2, 0)
	// set up a duplicate for (2,0) in test[3]
	tests[3].adjacentTriangles[0].C.dist2 = 0
	tests[3].adjacentTriangles[0].C.nearest = NewPoint(2, 0)
	for _, test := range tests {
		test.p.adjacentTriangles = test.adjacentTriangles
		gotNeighbor, gotDistance := test.p.SecondNearestNeighbor()
		if !floats.EqualWithinAbs(gotDistance, test.wantDist, tol) || !test.wantNeighbor.Equals(gotNeighbor) {
			t.Errorf("Second Nearest Neighbor for P%v with adjacent triangles %v, \n gotDist = %v, wantDist = %v, gotNeighbor = %v, wantNeighbor = %v",
				test.p, test.adjacentTriangles, gotDistance, test.wantDist, gotNeighbor, test.wantNeighbor)
		}
	}
}

func TestSurroundingPoints(t *testing.T) {
	p := NewPoint(0, 0)
	points := []*Point{NewPoint(1, 0), NewPoint(0, 1), NewPoint(-1, 0), NewPoint(0, -1)}
	triangles := []*Triangle{NewTriangle(p, points[0], points[1]), NewTriangle(p, points[1], points[2]), NewTriangle(p, points[2], points[3]), NewTriangle(p, points[3], points[0])}
	p.adjacentTriangles = triangles
	want := points
	got := p.surroundingPoints()
	if len(got) != len(want) {
		t.Errorf("SurroundingPoints for %v, got = %d, want = %d", p, len(got), len(want))
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("SurroundingPoints for %v, got[%d] = %v, want[%d] = %v", p, i, got, i, want)
		}
	}
}

func TestInTriangle(t *testing.T) {
	tests := []struct {
		x1, y1, x2, y2, x3, y3, x, y float64
		want                         location
	}{
		{2, 2, 5, 3, 6, 2, 5, 2.5, inside},
		{2, 2, 5, 3, 6, 2, 5, 2, onEdge},
		{2, 2, 5, 3, 6, 2, 4, 3, outside},
		{2, 2, 5, 3, 6, 2, 10, 2, outside},
		{5, 3, 2, 2, 6, 1.5, 5, 2.5, inside},
		{5, 3, 2, 2, 6, 1.5, 6, 2.5, outside},
	}
	for _, test := range tests {
		tri := NewTriangle(NewPoint(test.x1, test.y1), NewPoint(test.x2, test.y2), NewPoint(test.x3, test.y3))
		p := NewPoint(test.x, test.y)
		got := p.inTriangle(tri)
		if got != test.want {
			t.Fatalf("Point %v in Triangle %v, got = %v, want = %v", p, tri, got, test.want)
		}
	}
}

func TestPointsRemove(t *testing.T) {
	pts := points{NewPoint(0, 0), NewPoint(1, 1), NewPoint(2, 2), NewPoint(0, 0), NewPoint(3, 3)}
	want := points{pts[1], pts[2], pts[4]}
	got := pts.remove(NewPoint(0, 0))
	if len(got) != len(want) {
		t.Errorf("got %d points, want %d", len(got), len(want))
	}
	for i := range got {
		if !want[i].Equals(got[i]) {
			t.Errorf("got[%d] = %v, want[%d] = %v", i, got[i], i, want[i])
		}
	}
}

func TestPolyArea(t *testing.T) {
	tests := []struct {
		points points
		want   float64
	}{
		{points{NewPoint(0, 0), NewPoint(0, 1), NewPoint(1, 1), NewPoint(1, 0)}, 1},
		{points{NewPoint(0, 2), NewPoint(-1, 4), NewPoint(4, 5), NewPoint(5, -1), NewPoint(-2, 0)}, 27},
	}
	for i, test := range tests {
		got := test.points.polyArea()
		if !floats.EqualWithinAbs(got, test.want, tol) {
			t.Errorf("test %d: got = %f, want = %f", i, got, test.want)
		}
	}
}

func TestFindClockwiseTriangle(t *testing.T) {
	p := NewPoint(0, 0)
	p1 := NewPoint(2, 0)
	p2 := NewPoint(0, 3)
	p3 := NewPoint(-2, 0)
	p4 := NewPoint(0, -2)
	t1 := NewTriangle(p, p1, p2)
	t2 := NewTriangle(NewPoint(-4, 1), NewPoint(-4, 0), p3)
	t3 := NewTriangle(p, p2, p3)
	t4 := NewTriangle(p3, p4, p)
	t5 := NewTriangle(p4, p1, p)
	t6 := NewTriangle(p1, p2, NewPoint(5, 5))
	p.adjacentTriangles = triangles{t1, t3, t4, t5}
	p1.adjacentTriangles = triangles{t1, t5, t6}
	p2.adjacentTriangles = triangles{t1, t3, t6}
	p3.adjacentTriangles = triangles{t2, t3}
	p4.adjacentTriangles = triangles{t4, t5}
	want := []*Triangle{t1, t5, t4, t3}
	got := make([]*Triangle, len(p.adjacentTriangles))
	got[0] = p.adjacentTriangles[0]
	for i := 1; i < len(got); i++ {
		got[i] = p.findClockwiseTriangle(got[i-1])
	}
	if len(got) != len(want) {
		t.Fatalf("got = %d triangles, want = %d", len(got), len(want))
	}
	for i := range got {
		if !got[i].Equals(want[i]) {
			t.Errorf("got[%d] = %v, want[%d] = %v", i, got[i], i, want[i])
		}
	}
}
