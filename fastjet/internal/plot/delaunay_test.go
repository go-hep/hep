// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plot

import (
	"bytes"
	"image"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"go-hep.org/x/hep/fastjet/internal/delaunay"
)

func TestDelaunayPlot(t *testing.T) {
	fname := "testdata/delaunay_example.png"
	ref := "testdata/delaunay_golden.png"
	err := DelaunayPlotExample(fname)
	if err != nil {
		t.Fatal(err)
	}
	read := func(name string) image.Image {
		raw, err := ioutil.ReadFile(name)
		if err != nil {
			t.Fatalf("error reading %s: %v", name, err)
		}

		img, _, err := image.Decode(bytes.NewReader(raw))
		if err != nil {
			t.Fatalf("error decoding %s: %v", name, err)
		}
		return img
	}
	got := read(fname)
	want := read(ref)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("error: %s differs from reference file", fname)
	}
	os.Remove(fname)
}

func DelaunayPlotExample(path string) error {
	d := delaunay.HierarchicalDelaunay()
	points := []*delaunay.Point{
		delaunay.NewPoint(-1.5, 3.2),
		delaunay.NewPoint(1.8, 3.3),
		delaunay.NewPoint(-3.7, 1.5),
		delaunay.NewPoint(-1.5, 1.3),
		delaunay.NewPoint(0.8, 1.2),
		delaunay.NewPoint(3.3, 1.5),
		delaunay.NewPoint(-4, -1),
		delaunay.NewPoint(-2.3, -0.7),
		delaunay.NewPoint(0, -0.5),
		delaunay.NewPoint(2, -1.5),
		delaunay.NewPoint(3.7, -0.8),
		delaunay.NewPoint(-3.5, -2.9),
		delaunay.NewPoint(-0.9, -3.9),
		delaunay.NewPoint(2, -3.5),
		delaunay.NewPoint(3.5, -2.25),
	}
	for _, p := range points {
		d.Insert(p)
	}
	return Delaunay(path, d)
}
