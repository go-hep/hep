// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtree"
)

func TestReader(t *testing.T) {
	f, err := groot.Open("../testdata/simple.root")
	if err != nil {
		t.Fatalf("could not open ROOT file: %+v", err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatalf("could not retrieve ROOT tree: %+v", err)
	}
	tree := o.(rtree.Tree)

	for _, tc := range []struct {
		name  string
		rvars []rtree.ReadVar
		ropts []rtree.ReadOption
		beg   int64
		end   int64
		fun   func(rtree.RCtx) error
		enew  error
		eloop error
	}{
		{
			name: "ok",
			beg:  0, end: -1,
			fun: func(rtree.RCtx) error { return nil },
		},
		{
			name: "empty-range",
			beg:  4, end: -1,
			fun: func(rtree.RCtx) error { return nil },
		},
		{
			name:  "invalid-rvar",
			rvars: []rtree.ReadVar{{Name: "not-there", Value: new(int16)}},
			beg:   0, end: -1,
			fun:  func(rtree.RCtx) error { return nil },
			enew: fmt.Errorf(`rtree: could not create scanner: rtree: Tree "tree" has no branch named "not-there"`),
		},
		{
			name:  "invalid-ropt",
			ropts: []rtree.ReadOption{func(r *rtree.Reader) error { return io.EOF }},
			beg:   0, end: -1,
			fun:  func(rtree.RCtx) error { return nil },
			enew: fmt.Errorf(`rtree: could not set reader option 1: EOF`),
		},
		{
			name: "negative-start",
			beg:  -1, end: -1,
			fun:  func(rtree.RCtx) error { return nil },
			enew: fmt.Errorf("rtree: invalid event reader range [-1, 4) (start=-1 < 0)"),
		},
		{
			name: "start-greater-than-end",
			beg:  2, end: 1,
			fun:  func(rtree.RCtx) error { return nil },
			enew: fmt.Errorf("rtree: invalid event reader range [2, 1) (start=2 > end=1)"),
		},
		{
			name: "start-greater-than-nentries",
			beg:  5, end: 10,
			fun:  func(rtree.RCtx) error { return nil },
			enew: fmt.Errorf("rtree: invalid event reader range [5, 10) (start=5 > tree-entries=4)"),
		},
		{
			name: "end-greater-than-nentries",
			beg:  0, end: 5,
			fun:  func(rtree.RCtx) error { return nil },
			enew: fmt.Errorf("rtree: invalid event reader range [0, 5) (end=5 > tree-entries=4)"),
		},
		{
			name: "process-error",
			beg:  0, end: 4,
			fun: func(ctx rtree.RCtx) error {
				if ctx.Entry == 2 {
					return io.EOF
				}
				return nil
			},
			eloop: fmt.Errorf("rtree: could not process entry 2: EOF"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var (
				v1 int32
				v2 float32
				v3 string

				rvars = []rtree.ReadVar{
					{Name: "one", Value: &v1},
					{Name: "two", Value: &v2},
					{Name: "three", Value: &v3},
				}
			)

			if tc.rvars != nil {
				rvars = tc.rvars
			}

			ropts := []rtree.ReadOption{rtree.WithRange(tc.beg, tc.end)}
			if tc.ropts != nil {
				ropts = append(ropts, tc.ropts...)
			}

			r, err := rtree.NewReader(tree, rvars, ropts...)
			switch {
			case err != nil && tc.enew != nil:
				if got, want := err.Error(), tc.enew.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %v\nwant=%v", got, want)
				}
				return
			case err != nil && tc.enew == nil:
				t.Fatalf("unexpected error: %v", err)
			case err == nil && tc.enew != nil:
				t.Fatalf("expected an error: got=%v, want=%v", err, tc.enew)
			case err == nil && tc.enew == nil:
				// ok.
			}
			defer r.Close()

			err = r.Read(tc.fun)

			switch {
			case err != nil && tc.eloop != nil:
				if got, want := err.Error(), tc.eloop.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %v\nwant=%v", got, want)
				}
			case err != nil && tc.eloop == nil:
				t.Fatalf("unexpected error: %v", err)
			case err == nil && tc.eloop != nil:
				t.Fatalf("expected an error: got=%v, want=%v", err, tc.eloop)
			case err == nil && tc.eloop == nil:
				// ok.
			}

			err = r.Close()
			if err != nil {
				t.Fatalf("could not close tree reader: %+v", err)
			}

			// check r.Close is idem-potent.
			err = r.Close()
			if err != nil {
				t.Fatalf("tree reader close not idem-potent: %+v", err)
			}
		})
	}
}

func TestNewReadVars(t *testing.T) {
	f, err := groot.Open("../testdata/leaves.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := o.(rtree.Tree)

	vars := rtree.NewReadVars(tree)
	want := []rtree.ReadVar{
		{Name: "B", Leaf: "B", Value: new(bool)},
		{Name: "Str", Leaf: "Str", Value: new(string)},
		{Name: "I8", Leaf: "I8", Value: new(int8)},
		{Name: "I16", Leaf: "I16", Value: new(int16)},
		{Name: "I32", Leaf: "I32", Value: new(int32)},
		{Name: "I64", Leaf: "I64", Value: new(int64)},
		{Name: "U8", Leaf: "U8", Value: new(uint8)},
		{Name: "U16", Leaf: "U16", Value: new(uint16)},
		{Name: "U32", Leaf: "U32", Value: new(uint32)},
		{Name: "U64", Leaf: "U64", Value: new(uint64)},
		{Name: "F32", Leaf: "F32", Value: new(float32)},
		{Name: "F64", Leaf: "F64", Value: new(float64)},
		{Name: "D16", Leaf: "D16", Value: new(root.Float16)},
		{Name: "D32", Leaf: "D32", Value: new(root.Double32)},
		// arrays
		{Name: "ArrBs", Leaf: "ArrBs", Value: new([10]bool)},
		{Name: "ArrI8", Leaf: "ArrI8", Value: new([10]int8)},
		{Name: "ArrI16", Leaf: "ArrI16", Value: new([10]int16)},
		{Name: "ArrI32", Leaf: "ArrI32", Value: new([10]int32)},
		{Name: "ArrI64", Leaf: "ArrI64", Value: new([10]int64)},
		{Name: "ArrU8", Leaf: "ArrU8", Value: new([10]uint8)},
		{Name: "ArrU16", Leaf: "ArrU16", Value: new([10]uint16)},
		{Name: "ArrU32", Leaf: "ArrU32", Value: new([10]uint32)},
		{Name: "ArrU64", Leaf: "ArrU64", Value: new([10]uint64)},
		{Name: "ArrF32", Leaf: "ArrF32", Value: new([10]float32)},
		{Name: "ArrF64", Leaf: "ArrF64", Value: new([10]float64)},
		{Name: "ArrD16", Leaf: "ArrD16", Value: new([10]root.Float16)},
		{Name: "ArrD32", Leaf: "ArrD32", Value: new([10]root.Double32)},
		// slices
		{Name: "N", Leaf: "N", Value: new(int32)},
		{Name: "SliBs", Leaf: "SliBs", Value: new([]bool)},
		{Name: "SliI8", Leaf: "SliI8", Value: new([]int8)},
		{Name: "SliI16", Leaf: "SliI16", Value: new([]int16)},
		{Name: "SliI32", Leaf: "SliI32", Value: new([]int32)},
		{Name: "SliI64", Leaf: "SliI64", Value: new([]int64)},
		{Name: "SliU8", Leaf: "SliU8", Value: new([]uint8)},
		{Name: "SliU16", Leaf: "SliU16", Value: new([]uint16)},
		{Name: "SliU32", Leaf: "SliU32", Value: new([]uint32)},
		{Name: "SliU64", Leaf: "SliU64", Value: new([]uint64)},
		{Name: "SliF32", Leaf: "SliF32", Value: new([]float32)},
		{Name: "SliF64", Leaf: "SliF64", Value: new([]float64)},
		{Name: "SliD16", Leaf: "SliD16", Value: new([]root.Float16)},
		{Name: "SliD32", Leaf: "SliD32", Value: new([]root.Double32)},
	}

	n := len(want)
	if len(vars) < n {
		n = len(vars)
	}

	for i := 0; i < n; i++ {
		got := vars[i]
		if got.Name != want[i].Name {
			t.Fatalf("invalid read-var name[%d]: got=%q, want=%q", i, got.Name, want[i].Name)
		}
		if got.Leaf != want[i].Leaf {
			t.Fatalf("invalid read-var (name=%q) leaf-name[%d]: got=%q, want=%q", got.Name, i, got.Leaf, want[i].Leaf)
		}
		if got, want := reflect.TypeOf(got.Value), reflect.TypeOf(want[i].Value); got != want {
			t.Fatalf("invalid read-var (name=%q) type[%d]: got=%v, want=%v", vars[i].Name, i, got, want)
		}
	}

	if len(want) != len(vars) {
		t.Fatalf("invalid lengths. got=%d, want=%d", len(vars), len(want))
	}
}
