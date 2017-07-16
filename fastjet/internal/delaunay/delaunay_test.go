// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package delaunay

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

const tol = 1e-3

func TestSimple(t *testing.T) {
	// NewPoint(x, y, id)
	a := NewPoint(0, 0, 0)
	b := NewPoint(0, 2, 1)
	c := NewPoint(1, 0, 2)
	d := NewPoint(4, 4, 3)
	ep := NewPoint(3, 2, -1)
	points := []*Point{
		a,
		b,
		c,
		ep,
		d,
	}
	del := HierarchicalDelaunay(points, 4, 0, 4, 0)
	del.RemovePoint(ep)
	tri := del.Triangles()
	exp := []*Triangle{
		NewTriangle(a, b, c),
		NewTriangle(b, c, d),
	}
	got, want := len(tri), len(exp)
	if got != want {
		t.Fatalf("got=%d delaunay triangles, want=%d", got, want)
	}
	for i := range tri {
		ok := false
		for j := range exp {
			if tri[i].Equals(exp[j]) {
				ok = true
				// remove triangles that have been matched from slice
				// in case there are duplicate triangles. So that it
				// wouldn't pass the test when it shouldn't
				exp = append(exp[:j], exp[j+1:]...)
				break
			}
		}
		if !ok {
			t.Fatalf("Triangle T%s not as expected", tri[i])
		}
	}
	nn := make([]*Point, 0, len(points)-1)
	nd := make([]float64, 0, len(points)-1)
	for i, p := range points {
		if i == 3 { // skip the removed point
			continue
		}
		n, d := p.NearestNeighbor()
		nn = append(nn, n)
		nd = append(nd, d)
	}
	expN := []*Point{c, a, a, b}
	expD := []float64{1.0, 2.0, 1.0, 4.4721}
	got, want = len(nn), len(expN)
	if got != want {
		t.Fatalf("got=%d nearest neighbors, want=%d", got, want)
	}
	for i := range nn {
		if !nn[i].Equals(expN[i]) {
			t.Fatalf("got=N%s nearest neighbor, want=N%s", nn[i], expN[i])
		}
		if math.Abs(nd[i]-expD[i]) > tol {
			t.Fatalf("got=%f distance, want=%f for point P%s with neighbour N%s", nn[i].dist, expD[i], points[i], nn[i])
		}
	}
}

func TestWalkSimple(t *testing.T) {
	// NewPoint(x, y, id)
	a := NewPoint(0, 0, 0)
	b := NewPoint(0, 2, 1)
	c := NewPoint(1, 0, 2)
	d := NewPoint(4, 4, 3)
	ep := NewPoint(3, 2, -1)
	points := []*Point{
		a,
		b,
		c,
		ep,
		d,
	}
	del := WalkDelaunay(points)
	tri := del.Triangles()
	del.RemovePoint(ep)
	tri = del.Triangles()
	exp := []*Triangle{
		NewTriangle(a, b, c),
		NewTriangle(b, c, d),
	}
	got, want := len(tri), len(exp)
	if got != want {
		t.Fatalf("got=%d delaunay triangles, want=%d", got, want)
	}
	for i := range tri {
		ok := false
		for j := range exp {
			if tri[i].Equals(exp[j]) {
				ok = true
				// remove triangles that have been matched from slice
				// in case there are duplicate triangles. So that it
				// wouldn't pass the test when it shouldn't
				exp = append(exp[:j], exp[j+1:]...)
				break
			}
		}
		if !ok {
			t.Fatalf("Triangle T%s not as expected", tri[i])
		}
	}
	nn := make([]*Point, 0, len(points)-1)
	nd := make([]float64, 0, len(points)-1)
	for i, p := range points {
		if i == 3 { // skip the removed point
			continue
		}
		n, d := p.NearestNeighbor()
		nn = append(nn, n)
		nd = append(nd, d)
	}
	expN := []*Point{c, a, a, b}
	expD := []float64{1.0, 2.0, 1.0, 4.4721}
	got, want = len(nn), len(expN)
	if got != want {
		t.Fatalf("got=%d nearest neighbors, want=%d", got, want)
	}
	for i := range nn {
		if !nn[i].Equals(expN[i]) {
			t.Fatalf("got=N%s nearest neighbor, want=N%s", nn[i], expN[i])
		}
		if math.Abs(nd[i]-expD[i]) > tol {
			t.Fatalf("got=%f distance, want=%f for point P%s with neighbour N%s", nn[i].dist, expD[i], points[i], nn[i])
		}
	}
}

func TestMedium(t *testing.T) {
	// NewPoint(x, y, id)
	p1 := NewPoint(-1.5, 3.2, 1)
	p2 := NewPoint(1.8, 3.3, 2)
	p3 := NewPoint(-3.7, 1.5, 3)
	p4 := NewPoint(-1.5, 1.3, 4)
	p5 := NewPoint(0.8, 1.2, 5)
	p6 := NewPoint(3.3, 1.5, 6)
	p7 := NewPoint(-4, -1, 7)
	p8 := NewPoint(-2.3, -0.7, 8)
	p9 := NewPoint(0, -0.5, 9)
	p10 := NewPoint(2, -1.5, 10)
	p11 := NewPoint(3.7, -0.8, 11)
	p12 := NewPoint(-3.5, -2.9, 12)
	p13 := NewPoint(-0.9, -3.9, 13)
	p14 := NewPoint(2, -3.5, 14)
	p15 := NewPoint(3.5, -2.25, 15)
	pE1 := NewPoint(0, 0, -1)
	pE2 := NewPoint(-2.3, -0.6, -1)
	pE3 := NewPoint(2, 1.2, -1)
	pE4 := NewPoint(-2.8, -0.5, -1)
	ps := []*Point{p1, p2, p3, p4, p5, p6, pE3, pE4,
		p9, p10, p11, p12, p13, p14}
	d := HierarchicalDelaunay(ps, 4, -4, 4, -4)
	d.RemovePoint(pE4)
	d.InsertPoint(pE1)
	d.RemovePoint(pE3)
	d.InsertPoint(p15)
	d.InsertPoint(pE2)
	d.RemovePoint(pE1)
	d.InsertPoint(p7)
	d.InsertPoint(p8)
	d.RemovePoint(pE2)
	ts := d.Triangles()
	exp := []*Triangle{
		NewTriangle(p1, p3, p4),
		NewTriangle(p1, p4, p5),
		NewTriangle(p1, p5, p2),
		NewTriangle(p2, p5, p6),
		NewTriangle(p3, p4, p8),
		NewTriangle(p3, p8, p7),
		NewTriangle(p4, p8, p9),
		NewTriangle(p4, p9, p5),
		NewTriangle(p5, p9, p10),
		NewTriangle(p5, p10, p6),
		NewTriangle(p6, p10, p11),
		NewTriangle(p7, p8, p12),
		NewTriangle(p8, p12, p13),
		NewTriangle(p8, p13, p9),
		NewTriangle(p9, p13, p10),
		NewTriangle(p10, p13, p14),
		NewTriangle(p10, p14, p15),
		NewTriangle(p10, p15, p11),
	}
	got, want := len(ts), len(exp)
	if got != want {
		t.Fatalf("got=%d delaunay triangles, want=%d", got, want)
	}
	for i := range ts {
		ok := false
		for j := range exp {
			if ts[i].Equals(exp[j]) {
				ok = true
				// remove triangles that have been matched from slice
				// in case there are duplicate triangles. So that it
				// wouldn't pass the test when it shouldn't
				exp = append(exp[:j], exp[j+1:]...)
				break
			}
		}
		if !ok {
			t.Fatalf("Triangle T%s not as expected", ts[i])
		}
	}
	expv := []*Point{NewPoint(-2.523, 2.25, -1),
		NewPoint(-0.307, 2.25, -1), NewPoint(-0.373, 0.714, -1),
		NewPoint(-1.204, 0.022, -1), NewPoint(-2.672, 0.609, -1)}
	exparea := 4.3322215
	v := NewVoronoi(d)
	area, points := v.VoronoiCell(p4)
	got, want = len(points), len(expv)
	if got != want {
		t.Fatalf("got=%d voronoi points, want=%d", got, want)
	}
	for i := range points {
		ok := false
		for j := range expv {
			if points[i].EqualsApprox(expv[j], tol) {
				ok = true
				// remove points that have been matched from slice
				// in case there are duplicate points. So that it
				// wouldn't pass the test when it shouldn't
				expv = append(expv[:j], expv[j+1:]...)
				break
			}
		}
		if !ok {
			t.Fatalf("Point in the Voronoi Diagram P%s not as expected", points[i])
		}
	}

	if math.Abs(area-exparea) > tol {
		t.Fatalf("got=%f voronoi area, want=%f", got, want)
	}
	nn, dist := p11.NearestNeighbor()
	if !nn.Equals(p15) {
		t.Fatalf("got=N%s nearest neighbor, want=N%s", nn, p15)
	}
	expdist := 1.463
	if math.Abs(dist-expdist) > tol {
		t.Fatalf("got=%d distance, want=%d", dist, expdist)
	}
}

func TestWalkMedium(t *testing.T) {
	// NewPoint(x, y, id)
	p1 := NewPoint(-1.5, 3.2, 1)
	p2 := NewPoint(1.8, 3.3, 2)
	p3 := NewPoint(-3.7, 1.5, 3)
	p4 := NewPoint(-1.5, 1.3, 4)
	p5 := NewPoint(0.8, 1.2, 5)
	p6 := NewPoint(3.3, 1.5, 6)
	p7 := NewPoint(-4, -1, 7)
	p8 := NewPoint(-2.3, -0.7, 8)
	p9 := NewPoint(0, -0.5, 9)
	p10 := NewPoint(2, -1.5, 10)
	p11 := NewPoint(3.7, -0.8, 11)
	p12 := NewPoint(-3.5, -2.9, 12)
	p13 := NewPoint(-0.9, -3.9, 13)
	p14 := NewPoint(2, -3.5, 14)
	p15 := NewPoint(3.5, -2.25, 15)
	pE1 := NewPoint(0, 0, -1)
	pE2 := NewPoint(-2.3, -0.6, -2)
	pE3 := NewPoint(2, 1.2, -3)
	pE4 := NewPoint(-2.8, -0.5, -4)
	ps := []*Point{p1, p2, p3, p4, p5, p6, pE3, pE4,
		p9, p10, p11, p12, p13, p14}
	d := WalkDelaunay(ps)
	d.RemovePoint(pE4)
	d.InsertPoint(pE1)
	d.RemovePoint(pE3)
	d.InsertPoint(p15)
	d.InsertPoint(pE2)
	d.RemovePoint(pE1)
	d.InsertPoint(p7)
	d.InsertPoint(p8)
	d.RemovePoint(pE2)
	ts := d.Triangles()
	exp := []*Triangle{
		NewTriangle(p1, p3, p4),
		NewTriangle(p1, p4, p5),
		NewTriangle(p1, p5, p2),
		NewTriangle(p2, p5, p6),
		NewTriangle(p3, p4, p8),
		NewTriangle(p3, p8, p7),
		NewTriangle(p4, p8, p9),
		NewTriangle(p4, p9, p5),
		NewTriangle(p5, p9, p10),
		NewTriangle(p5, p10, p6),
		NewTriangle(p6, p10, p11),
		NewTriangle(p7, p8, p12),
		NewTriangle(p8, p12, p13),
		NewTriangle(p8, p13, p9),
		NewTriangle(p9, p13, p10),
		NewTriangle(p10, p13, p14),
		NewTriangle(p10, p14, p15),
		NewTriangle(p10, p15, p11),
	}
	got, want := len(ts), len(exp)
	if got != want {
		t.Fatalf("got=%d delaunay triangles, want=%d", got, want)
	}
	for i := range ts {
		ok := false
		for j := range exp {
			if ts[i].Equals(exp[j]) {
				ok = true
				// remove triangles that have been matched from slice
				// in case there are duplicate triangles. So that it
				// wouldn't pass the test when it shouldn't
				exp = append(exp[:j], exp[j+1:]...)
				break
			}
		}
		if !ok {
			t.Fatalf("Triangle T%s not as expected", ts[i])
		}
	}
	nn, dist := p11.NearestNeighbor()
	if !nn.Equals(p15) {
		t.Fatalf("got=N%s nearest neighbor, want=N%s", nn, p15)
	}
	expdist := 1.463
	if math.Abs(dist-expdist) > tol {
		t.Fatalf("got=%d distance, want=%d", dist, expdist)
	}
}

func BenchmarkNewDelaunay50(b *testing.B) {
	benchmarkDelaunay(50, b)
}

func BenchmarkNewDelaunay100(b *testing.B) {
	benchmarkDelaunay(100, b)
}

func BenchmarkNewDelaunay150(b *testing.B) {
	benchmarkDelaunay(150, b)
}

func BenchmarkNewDelaunay200(b *testing.B) {
	benchmarkDelaunay(200, b)
}

func BenchmarkNewDelaunay250(b *testing.B) {
	benchmarkDelaunay(250, b)
}

func BenchmarkNewDelaunay300(b *testing.B) {
	benchmarkDelaunay(300, b)
}

func BenchmarkNewDelaunay350(b *testing.B) {
	benchmarkDelaunay(350, b)
}

func BenchmarkNewDelaunay400(b *testing.B) {
	benchmarkDelaunay(400, b)
}

func BenchmarkNewDelaunay450(b *testing.B) {
	benchmarkDelaunay(450, b)
}

func BenchmarkNewDelaunay500(b *testing.B) {
	benchmarkDelaunay(500, b)
}

func BenchmarkNewDelaunay550(b *testing.B) {
	benchmarkDelaunay(550, b)
}

func BenchmarkNewDelaunay600(b *testing.B) {
	benchmarkDelaunay(600, b)
}

func BenchmarkNewDelaunay650(b *testing.B) {
	benchmarkDelaunay(650, b)
}

func BenchmarkNewDelaunay700(b *testing.B) {
	benchmarkDelaunay(700, b)
}

func BenchmarkNewDelaunay750(b *testing.B) {
	benchmarkDelaunay(750, b)
}

func BenchmarkNewDelaunay800(b *testing.B) {
	benchmarkDelaunay(800, b)
}

func BenchmarkNewDelaunay850(b *testing.B) {
	benchmarkDelaunay(850, b)
}

func BenchmarkNewDelaunay900(b *testing.B) {
	benchmarkDelaunay(900, b)
}

func BenchmarkNewDelaunay950(b *testing.B) {
	benchmarkDelaunay(950, b)
}

func BenchmarkNewDelaunay1000(b *testing.B) {
	benchmarkDelaunay(1000, b)
}

func BenchmarkDelaunay_VoronoiArea(b *testing.B) {
	points := make([]*Point, 100)
	rand.Seed(int64(time.Now().Nanosecond()))
	for j := 0; j < 100; j++ {
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		points[j] = NewPoint(x, y, j)
	}
	d := HierarchicalDelaunay(points, 1000, 0, 1000, 0)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		v := NewVoronoi(d)
		v.VoronoiCell(points[rand.Intn(100)])
	}
}

func benchmarkDelaunay(i int, b *testing.B) {
	points := make([]*Point, i)
	rand.Seed(int64(time.Now().Nanosecond()))
	for j := 0; j < i; j++ {
		x := rand.Float64() * 1000
		y := rand.Float64() * 1000
		points[j] = NewPoint(x, y, j)
	}
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		d := HierarchicalDelaunay(points, 1000, 0, 1000, 0)
		for _, p := range points {
			d.RemovePoint(p)
		}
	}
}
