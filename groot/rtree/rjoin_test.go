// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"testing"

	"go-hep.org/x/hep/groot/riofs"
)

func TestRJoin(t *testing.T) {
	get := func(fname, tname string) (Tree, func() error) {
		t.Helper()

		f, err := riofs.Open(fname)
		if err != nil {
			t.Fatal(err)
		}
		tree, err := f.Get(tname)
		if err != nil {
			_ = f.Close()
			t.Fatal(err)
		}
		return tree.(Tree), f.Close
	}
	chk := func(f func() error) {
		err := f()
		if err != nil {
			t.Fatal(err)
		}
	}

	t1, close1 := get("../testdata/join1.root", "j1")
	defer chk(close1)

	t2, close2 := get("../testdata/join2.root", "j2")
	defer chk(close2)

	t3, close3 := get("../testdata/join3.root", "j3")
	defer chk(close3)

	for _, tc := range []struct {
		name  string
		trees []Tree
		rvars []ReadVar
		err   error
		beg   int64
		end   int64
		fct   func(RCtx) error
	}{
		{
			name:  "all",
			trees: []Tree{t1, t2, t3},
			beg:   0,
			end:   10,
			fct:   func(RCtx) error { return nil },
		},
		{
			name:  "sub-range",
			trees: []Tree{t1, t2, t3},
			beg:   3,
			end:   4,
			fct:   func(RCtx) error { return nil },
		},
		{
			name:  "empty-range",
			trees: []Tree{t1, t2, t3},
			beg:   3,
			end:   3,
			fct:   func(RCtx) error { return nil },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			join, err := Join(tc.trees...)
			if err != nil {
				t.Fatalf("could not create joined-tree: %+v", err)
			}

			rvars := tc.rvars
			if rvars == nil {
				rvars = NewReadVars(join)
			}

			r, err := NewReader(join, rvars, WithRange(tc.beg, tc.end))
			if err != nil {
				t.Fatalf("could not create reader: %+v", err)
			}
			defer r.Close()

			err = r.Read(tc.fct)
			if err != nil {
				t.Fatalf("could not run reader: %+v", err)
			}
		})
	}
}
