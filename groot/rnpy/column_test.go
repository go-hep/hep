// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rnpy

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func TestColumn(t *testing.T) {
	for _, tc := range []struct {
		fname string
		tname string
		aname rtree.ReadVar
		want  any
		err   error
	}{
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			aname: rtree.ReadVar{Name: "one"},
			want:  []int32{1, 2, 3, 4},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			aname: rtree.ReadVar{Name: "two"},
			want:  []float32{1.1, 2.2, 3.3, 4.4},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			aname: rtree.ReadVar{Name: "three"},
			want:  []string{"uno", "dos", "tres", "quatro"},
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			aname: rtree.ReadVar{Name: "ArrF64"},
			want: func() any {
				o := make([][10]float64, 10)
				for i := range o {
					for j := range o[i] {
						o[i][j] = float64(i)
					}
				}
				return o
			}(),
		},
		{
			fname: "../testdata/root_numpy_struct.root",
			tname: "test",
			aname: rtree.ReadVar{Name: "branch1", Leaf: "intleaf"},
			want:  []int32{10},
		},
		{
			fname: "../testdata/root_numpy_struct.root",
			tname: "test",
			aname: rtree.ReadVar{Name: "branch1", Leaf: "floatleaf"},
			want:  []float32{15.5},
		},
		{
			fname: "../testdata/root_numpy_struct.root",
			tname: "test",
			aname: rtree.ReadVar{Name: "branch2", Leaf: "intleaf"},
			want:  []int32{20},
		},
		{
			fname: "../testdata/root_numpy_struct.root",
			tname: "test",
			aname: rtree.ReadVar{Name: "branch2", Leaf: "floatleaf"},
			want:  []float32{781.2},
		},
		{
			fname: "../testdata/simple.root",
			tname: "tree",
			aname: rtree.ReadVar{Name: "not_there"},
			err:   fmt.Errorf(`rnpy: no rvar named "not_there"`),
		},
		{
			fname: "../testdata/root_numpy_struct.root",
			tname: "test",
			aname: rtree.ReadVar{Name: "branch1", Leaf: "not_there"},
			err:   fmt.Errorf(`rnpy: no rvar named "branch1.not_there"`),
		},
		{
			fname: "../testdata/root_numpy_struct.root",
			tname: "test",
			aname: rtree.ReadVar{Name: "branch2", Leaf: "not_there"},
			err:   fmt.Errorf(`rnpy: no rvar named "branch2.not_there"`),
		},
		{
			fname: "../testdata/leaves.root",
			tname: "tree",
			aname: rtree.ReadVar{Name: "SliF64"},
			err:   fmt.Errorf(`rnpy: invalid branch or leaf type []float64`),
		},
	} {
		t.Run("", func(t *testing.T) {
			f, err := riofs.Open(tc.fname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			obj, err := riofs.Dir(f).Get(tc.tname)
			if err != nil {
				t.Fatal(err)
			}
			tree := obj.(rtree.Tree)

			var sli any

			col, err := NewColumn(tree, tc.aname)
			if err == nil {
				sli, err = col.Slice()
			}

			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %s\nwant=%s", got, want)
				}
				return
			case err != nil && tc.err == nil:
				t.Fatalf("could not create column: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("expected an error: %+v", tc.err)
			case err == nil && tc.err == nil:
				// ok.
			}

			if got, want := sli, tc.want; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid slice:\ngot= %+v\nwant=%+v", got, want)
			}
		})
	}
}

func TestColumnReadFail(t *testing.T) {
	f, err := riofs.Open("../testdata/leaves.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	obj, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}
	tree := obj.(rtree.Tree)

	col, err := NewColumn(tree, rtree.ReadVar{Name: "F64"})
	if err != nil {
		t.Fatal(err)
	}

	want := fmt.Errorf(`rnpy: could not create ROOT reader for "boo": rtree: could not create reader: rtree: tree "tree" has no branch named "boo"`)

	col.rvar.Name = "boo"
	_, err = col.Slice()
	if got, want := err.Error(), want.Error(); got != want {
		t.Fatalf("invalid error:\ngot= %+v\nwant=%+v", got, want)
	}
}
