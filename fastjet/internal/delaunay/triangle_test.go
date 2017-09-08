// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"testing"

	"gonum.org/v1/gonum/floats"
)

func TestTriangleEquals(t *testing.T) {
	tests := []struct {
		a, b *Triangle
		want bool
	}{
		{NewTriangle(NewPoint(3, 2), NewPoint(0, 0), NewPoint(5, 0)), NewTriangle(NewPoint(3, 2), NewPoint(0, 0), NewPoint(5, 0)), true},
		{NewTriangle(NewPoint(3, 2), NewPoint(0, 0), NewPoint(5, 0)), NewTriangle(NewPoint(5, 0), NewPoint(3, 2), NewPoint(0, 0)), true},
		{NewTriangle(NewPoint(3, 2), NewPoint(0, 0), NewPoint(5, 0)), NewTriangle(NewPoint(3, 1), NewPoint(0, 0), NewPoint(5, 0)), false},
		{NewTriangle(NewPoint(3, 2), NewPoint(0, 0), NewPoint(5, 0)), NewTriangle(NewPoint(5, 0), NewPoint(3, 1), NewPoint(0, 0)), false},
	}
	for i, test := range tests {
		got := test.a.Equals(test.b)
		if got != test.want {
			t.Errorf("Test case %v: %v == %v,\n got = %v, want = %v", i, test.a, test.b, got, test.want)
		}
	}
}

func TestTrianglesRemove(t *testing.T) {
	tests := []struct {
		triangles triangles
		toRemove  []*Triangle
	}{
		{
			triangles{NewTriangle(NewPoint(0, 0), NewPoint(0, 1), NewPoint(1, 0)), NewTriangle(NewPoint(0, 0), NewPoint(5, 1), NewPoint(1, 0)), NewTriangle(NewPoint(0, 5), NewPoint(0, 1), NewPoint(1, 0))},
			[]*Triangle{NewTriangle(NewPoint(0, 0), NewPoint(0, 1), NewPoint(1, 0)), NewTriangle(NewPoint(0, 5), NewPoint(0, 1), NewPoint(1, 0))},
		},
	}
	for i, test := range tests {
		toRemove := make([]*Triangle, len(test.toRemove))
		copy(toRemove, test.toRemove)
		triangles := test.triangles.remove(toRemove...)
		for _, tri := range triangles {
			for _, removed := range test.toRemove {
				if tri.Equals(removed) {
					t.Errorf("Test case %v: Removed triangle %v still in triangles %v", i, removed, triangles)
				}
			}
		}
	}
}

func TestTriangleAdd(t *testing.T) {
	points := []*Point{
		NewPoint(0, 0),
		NewPoint(2, 3),
		NewPoint(0, 2),
		NewPoint(5, 5),
		NewPoint(0.5, -6),
	}
	t1 := NewTriangle(points[0], points[1], points[2])
	t1.add()
	t2 := NewTriangle(points[1], points[2], points[3])
	t2.add()
	t3 := NewTriangle(points[0], points[2], points[4])
	t3.add()
	want := []*Point{
		points[2], points[2], points[0], points[1], points[0],
	}
	var got []*Point
	for _, p := range points {
		got = append(got, p.nearest)
	}
	for i := range got {
		if !got[i].Equals(want[i]) {
			t.Errorf("After adding triangle nearest neighbor of points[%d]=%v, got = %v, want = %v", i, points[i], got[i], want[i])
		}
	}
}

func TestTriangleRemove(t *testing.T) {
	points := []*Point{
		NewPoint(0, 0),
		NewPoint(2, 3),
		NewPoint(0, 2),
		NewPoint(1, 10),
	}
	t1 := NewTriangle(points[0], points[1], points[2])
	t1.add()
	t2 := NewTriangle(points[0], points[1], points[3])
	t2.add()
	t1.remove()
	want := []*Point{
		points[1], points[0], nil, points[1],
	}
	var got []*Point
	for _, p := range points {
		got = append(got, p.nearest)
	}
	for i := range got {
		if !got[i].Equals(want[i]) {
			t.Errorf("After removing triangle nearest neighbor of points[%d]=%v, got = %v, want = %v", i, points[i], got[i], want[i])
		}
	}
}

func TestTriangleCircumcenter(t *testing.T) {
	tests := []struct {
		t     *Triangle
		wantX float64
		wantY float64
	}{
		{NewTriangle(NewPoint(0, 0), NewPoint(2, 0), NewPoint(0, 2)), 1, 1},
		{NewTriangle(NewPoint(-1, 4), NewPoint(3, 8), NewPoint(4, 12)), -5.1667, 12.1667},
		{NewTriangle(NewPoint(6, 9), NewPoint(8, 3), NewPoint(5, 15)), 44.5, 18.5},
		{NewTriangle(NewPoint(2, 8), NewPoint(3, 7), NewPoint(3, 8)), 2.5, 7.5},
	}
	for _, test := range tests {
		gotX, gotY := test.t.circumcenter()
		if !floats.EqualWithinAbs(gotX, test.wantX, tol) || !floats.EqualWithinAbs(gotY, test.wantY, tol) {
			t.Errorf("Circumcenter of %v, got = (%f,%f), want = (%f,%f)", test.t, gotX, gotY, test.wantX, test.wantY)
		}
	}
}
