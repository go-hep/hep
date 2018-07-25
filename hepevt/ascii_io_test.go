// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepevt_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hepevt"
)

var (
	small = hepevt.Event{
		Nevhep: 1,
		Nhep:   8,
		Isthep: []int{3, 3, 3, 3, 3, 1, 1, 1},
		Idhep:  []int{2212, 2212, -2, 1, -24, -2, 1, 22},
		Jmohep: [][2]int{{0, 0}, {0, 0}, {2, 2}, {1, 1}, {3, 4}, {5, 5}, {5, 5}, {3, 4}},
		Jdahep: [][2]int{{4, 4}, {3, 3}, {5, 8}, {5, 8}, {6, 7}, {0, 0}, {0, 0}, {0, 0}},
		Phep: [][5]float64{
			{+0.00000000E+00, +0.00000000E+00, +7.00000000E+03, +7.0000000000E+03, +0.00000000E+00},
			{+0.00000000E+00, +0.00000000E+00, -7.00000000E+03, +7.0000000000E+03, +0.00000000E+00},
			{-3.04700000E+00, -1.90000000E+01, -5.46290000E+01, +5.7920000000E+01, +0.00000000E+00},
			{+7.50000000E-01, -1.56900000E+00, +3.21910000E+01, +3.2238000000E+01, +0.00000000E+00},
			{+1.51700000E+00, -2.06800000E+01, -2.06050000E+01, +8.5925000000E+01, +0.00000000E+00},
			{+3.96200000E+00, -4.94980000E+01, -2.66870000E+01, +5.6373000000E+01, +0.00000000E+00},
			{-2.44500000E+00, +2.88160000E+01, +6.08200000E+00, +2.9552000000E+01, +0.00000000E+00},
			{-3.81300000E+00, +1.13000000E-01, -1.83300000E+00, +4.2330000000E+00, +0.00000000E+00},
		},
		Vhep: [][4]float64{
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
			{0, 0, 0, 0},
		},
	}
)

func TestDecoder(t *testing.T) {
	f, err := os.Open("testdata/small.hepevt")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	var evt hepevt.Event

	dec := hepevt.NewDecoder(f)
	err = dec.Decode(&evt)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(evt, small) {
		t.Fatalf("decoded event differ from reference:\ngot = %#v\nwant= %#v",
			evt, small,
		)
	}
}

func TestEncoder(t *testing.T) {
	f, err := ioutil.TempFile("", "hepevt-")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())

	enc := hepevt.NewEncoder(f)
	if err := enc.Encode(&small); err != nil {
		t.Fatal(err)
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/small.hepevt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("files differ.\ngot:\n%s\nwant:\n%s\n", got, want)
	}
}
