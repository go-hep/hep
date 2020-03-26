// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"fmt"
	"io"
	"testing"

	"go-hep.org/x/hep/groot"
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
