// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-hep/hbook"
)

func TestScatter2D(t *testing.T) {
	s := hbook.NewScatter2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})
	if s == nil {
		t.Fatal("nil pointer to Scatter2D")
	}

	if got, want := s.Len(), 3; got != want {
		t.Errorf("got len=%d. want=%d\n", got, want)
	}

	pt := hbook.Point2D{X: 10, Y: -10, ErrX: hbook.Range{Min: 5, Max: 5}, ErrY: hbook.Range{Min: 6, Max: 6}}
	s.Fill(pt)

	if got, want := s.Len(), 4; got != want {
		t.Errorf("got len=%d. want=%d\n", got, want)
	}

	if got, want := s.Point(3), pt; got != want {
		t.Errorf("invalid pt[%d]:\ngot= %+v\nwant=%+v\n", 3, got, want)
	}
}

func ExampleScatter2D() {
	s := hbook.NewScatter2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})
	if s == nil {
		log.Fatal("nil pointer to Scatter2D")
	}

	fmt.Printf("len=%d\n", s.Len())

	s.Fill(hbook.Point2D{X: 10, Y: -10, ErrX: hbook.Range{Min: 5, Max: 5}, ErrY: hbook.Range{Min: 6, Max: 6}})
	fmt.Printf("len=%d\n", s.Len())
	fmt.Printf("pt[%d]=%+v\n", 3, s.Point(3))

	// Output:
	// len=3
	// len=4
	// pt[3]={X:10 Y:-10 ErrX:{Min:5 Max:5} ErrY:{Min:6 Max:6}}
}

func ExampleScatter2D_newScatter2DFrom() {
	s := hbook.NewScatter2DFrom([]float64{1, 2, -1}, []float64{1, 1.5, 2})
	if s == nil {
		log.Fatal("nil pointer to Scatter2D")
	}

	fmt.Printf("len=%d\n", s.Len())

	s.Fill(hbook.Point2D{X: 10, Y: -10, ErrX: hbook.Range{Min: 5, Max: 5}, ErrY: hbook.Range{Min: 6, Max: 6}})
	fmt.Printf("len=%d\n", s.Len())
	fmt.Printf("pt[%d]=%+v\n", 3, s.Point(3))

	// Output:
	// len=3
	// len=4
	// pt[3]={X:10 Y:-10 ErrX:{Min:5 Max:5} ErrY:{Min:6 Max:6}}
}
