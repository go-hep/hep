// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook_test

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hbook"
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

	ref, err := ioutil.ReadFile("testdata/s2d_v1_golden.yoda")
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

func TestS2DReadYODAv1(t *testing.T) {
	ref, err := ioutil.ReadFile("testdata/s2d_v1_golden.yoda")
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
