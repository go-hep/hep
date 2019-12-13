// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sbinet/npyio"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func TestConvert(t *testing.T) {
	for _, tc := range []struct {
		name string
		want interface{}
	}{
		// 4 scalars
		{
			name: "bool_4x1",
			want: [4]bool{true, true, false, true},
		},
		{
			name: "uint8_4x1",
			want: [4]uint8{0, 1, 2, 3},
		},
		{
			name: "uint16_4x1",
			want: [4]uint16{0, 1, 2, 3},
		},
		{
			name: "uint32_4x1",
			want: [4]uint32{0, 1, 2, 3},
		},
		{
			name: "uint64_4x1",
			want: [4]uint64{0, 1, 2, 3},
		},
		{
			name: "int8_4x1",
			want: [4]int8{0, 1, 2, 3},
		},
		{
			name: "int16_4x1",
			want: [4]int16{0, 1, 2, 3},
		},
		{
			name: "int32_4x1",
			want: [4]int32{0, 1, 2, 3},
		},
		{
			name: "int64_4x1",
			want: [4]int64{0, 1, 2, 3},
		},
		{
			name: "float32_4x1",
			want: [4]float32{0, 1, 2, 3},
		},
		{
			name: "float64_4x1",
			want: [4]float64{0, 1, 2, 3},
		},
		{
			name: "nans_4x1",
			want: [4]float64{math.Inf(-1), math.Inf(+1), math.NaN(), 0},
		},
		// 4 1d-arrays
		{
			name: "bool_4x2",
			want: [4][2]bool{{true, false}, {true, false}, {false, true}, {true, false}},
		},
		{
			name: "uint8_4x2",
			want: [4][2]uint8{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "uint16_4x2",
			want: [4][2]uint16{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "uint32_4x2",
			want: [4][2]uint32{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "uint64_4x2",
			want: [4][2]uint64{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int8_4x2",
			want: [4][2]int8{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int16_4x2",
			want: [4][2]int16{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int32_4x2",
			want: [4][2]int32{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "int64_4x2",
			want: [4][2]int64{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "float32_4x2",
			want: [4][2]float32{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
		{
			name: "float64_4x2",
			want: [4][2]float64{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tmp, err := ioutil.TempDir("", "npy2root-")
			if err != nil {
				t.Fatalf("%+v", err)
			}
			defer os.RemoveAll(tmp)

			fname := filepath.Join(tmp, "data.npy")
			src, err := os.Create(fname)
			if err != nil {
				t.Fatalf("could not create NumPy data file: %+v", err)
			}
			defer src.Close()

			err = npyio.Write(src, tc.want)
			if err != nil {
				t.Fatalf("could not save NumPy data file: %+v", err)
			}

			err = src.Close()
			if err != nil {
				t.Fatalf("could not close NumPy data file: %+v", err)
			}

			oname := filepath.Join(tmp, "out.root")
			err = process(oname, "tree", fname)
			if err != nil {
				t.Fatalf("could not create ROOT data file: %+v", err)
			}

			f, err := groot.Open(oname)
			if err != nil {
				t.Fatalf("could not open ROOT file: %+v", err)
			}
			defer f.Close()

			obj, err := f.Get("tree")
			if err != nil {
				t.Fatalf("could not get ROOT tree: %+v", err)
			}

			tree := obj.(rtree.Tree)
			rvars := rtree.NewScanVars(tree)
			scan, err := rtree.NewScannerVars(tree, rvars...)
			if err != nil {
				t.Fatalf("could not create tree scanner: %+v", err)
			}
			defer scan.Close()

			want := reflect.ValueOf(tc.want)
			n := 0
			for scan.Next() {
				err := scan.Scan()
				if err != nil {
					t.Fatalf("could not read entry %d: %+v", scan.Entry(), err)
				}

				i := int(scan.Entry())
				want := want.Index(i).Interface()
				got := reflect.ValueOf(rvars[0].Value).Elem().Interface()
				ok := false
				switch want := want.(type) {
				case float64:
					got := got.(float64)
					switch {
					case math.IsNaN(want):
						ok = math.IsNaN(got)
					case math.IsInf(want, +1):
						ok = math.IsInf(got, +1)
					case math.IsInf(want, -1):
						ok = math.IsInf(got, -1)
					default:
						ok = got == want
					}
				default:
					ok = reflect.DeepEqual(got, want)
				}

				if !ok {
					t.Fatalf("invalid value for entry %d:\ngot= %v\nwant=%v", scan.Entry(), got, want)
				}
				n++
			}

			if got, want := n, want.Len(); got != want {
				t.Fatalf("invalid number of events: got=%d, want=%d", got, want)
			}
		})
	}
}
