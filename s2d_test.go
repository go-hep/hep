// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"reflect"
	"testing"

	"github.com/go-hep/hbook"
)

func TestS2D(t *testing.T) {
	s := hbook.NewS2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})
	if s == nil {
		t.Fatal("nil pointer to S2D")
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

func ExampleS2D() {
	s := hbook.NewS2D(hbook.Point2D{X: 1, Y: 1}, hbook.Point2D{X: 2, Y: 1.5}, hbook.Point2D{X: -1, Y: +2})
	if s == nil {
		log.Fatal("nil pointer to S2D")
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

func ExampleS2D_newS2DFrom() {
	s := hbook.NewS2DFrom([]float64{1, 2, -1}, []float64{1, 1.5, 2})
	if s == nil {
		log.Fatal("nil pointer to S2D")
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

func ExampleS2D_newS2DFromH1D() {
	h := hbook.NewH1D(20, -4, +4)
	h.Fill(1, 2)
	h.Fill(2, 3)
	h.Fill(3, 1)
	h.Fill(1, 1)
	h.Fill(-2, 1)
	h.Fill(-3, 1)

	s := hbook.NewS2DFromH1D(h)
	s.Sort()
	for _, pt := range s.Points() {
		fmt.Printf("point=(%+3.2f +/- (%+3.2f,%+3.2f), %+3.2f +/- (%+3.2f, %+3.2f))\n", pt.X, pt.ErrX.Min, pt.ErrX.Max, pt.Y, pt.ErrY.Min, pt.ErrY.Max)
	}

	// Output:
	// point=(-3.80 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-3.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-3.00 +/- (+0.20,+0.20), +2.50 +/- (+2.50, +2.50))
	// point=(-2.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-2.20 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-1.80 +/- (+0.20,+0.20), +2.50 +/- (+2.50, +2.50))
	// point=(-1.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-1.00 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-0.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(-0.20 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+0.20 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+0.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+1.00 +/- (+0.20,+0.20), +7.50 +/- (+5.59, +5.59))
	// point=(+1.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+1.80 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+2.20 +/- (+0.20,+0.20), +7.50 +/- (+7.50, +7.50))
	// point=(+2.60 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+3.00 +/- (+0.20,+0.20), +2.50 +/- (+2.50, +2.50))
	// point=(+3.40 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
	// point=(+3.80 +/- (+0.20,+0.20), +0.00 +/- (+0.00, +0.00))
}

func TestS2DWriteYODA(t *testing.T) {
	h := hbook.NewH1D(20, -4, +4)
	h.Fill(1, 2)
	h.Fill(2, 3)
	h.Fill(3, 1)
	h.Fill(1, 1)
	h.Fill(-2, 1)
	h.Fill(-3, 1)

	s := hbook.NewS2DFromH1D(h)

	chk, err := s.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	ref, err := ioutil.ReadFile("testdata/s2d_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("s2d file differ:\n=== got ===\n%s\n=== want ===\n%s\n",
			string(chk),
			string(ref),
		)
	}
}

func TestS2DReadYODA(t *testing.T) {
	ref, err := ioutil.ReadFile("testdata/s2d_golden.yoda")
	if err != nil {
		t.Fatal(err)
	}

	var s hbook.S2D
	err = s.UnmarshalYODA(ref)
	if err != nil {
		t.Fatal(err)
	}

	chk, err := s.MarshalYODA()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(chk, ref) {
		t.Fatalf("s2d file differ:\n=== got ===\n%s\n=== want ===\n%s\n",
			string(chk),
			string(ref),
		)
	}
}

func TestS2DSerialization(t *testing.T) {
	sref := hbook.NewS2D()
	for i := 0; i < 10; i++ {
		v := float64(i)
		sref.Fill(hbook.Point2D{X: v, Y: v, ErrX: hbook.Range{Min: v, Max: 2 * v}, ErrY: hbook.Range{Min: v, Max: 3 * v}})
	}
	sref.Annotation()["title"] = "scatter2d title"
	sref.Annotation()["name"] = "s2d-name"

	{
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err := enc.Encode(sref)
		if err != nil {
			t.Fatalf("could not serialize scatter2d: %v\n", err)
		}

		var snew hbook.S2D
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&snew)
		if err != nil {
			t.Fatalf("could not deserialize scatter2d: %v\n", err)
		}

		if !reflect.DeepEqual(sref, &snew) {
			t.Fatalf("ref=%v\nnew=%v\n", sref, &snew)
		}
	}
}
