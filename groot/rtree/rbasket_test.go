// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
)

func TestReadAheadBasket(t *testing.T) {
	for _, tc := range []struct {
		fname  string
		tree   string
		branch string
		conc   int
		want   []rspan
	}{
		{
			fname:  "../testdata/simple.root",
			tree:   "tree",
			branch: "two",
			conc:   -1,
			want:   []rspan{{beg: 0, end: 4, pos: 304, sz: 86}},
		},
		{
			fname:  "../testdata/simple.root",
			tree:   "tree",
			branch: "two",
			conc:   0,
			want:   []rspan{{beg: 0, end: 4, pos: 304, sz: 86}},
		},
		{
			fname:  "../testdata/simple.root",
			tree:   "tree",
			branch: "two",
			conc:   1,
			want:   []rspan{{beg: 0, end: 4, pos: 304, sz: 86}},
		},
		{
			fname:  "../testdata/simple.root",
			tree:   "tree",
			branch: "two",
			conc:   10,
			want:   []rspan{{beg: 0, end: 4, pos: 304, sz: 86}},
		},
		{
			fname:  "../testdata/small-flat-tree.root",
			tree:   "tree",
			branch: "Float64",
			conc:   10,
			want:   []rspan{{beg: 0, end: 100, pos: 1551, sz: 297}},
		},
		{
			fname:  "../testdata/small-flat-tree.root",
			tree:   "tree",
			branch: "SliceFloat64",
			conc:   -1,
			want:   []rspan{{beg: 0, end: 100, pos: 8112, sz: 690}},
		},
		{
			fname:  "../testdata/chain.flat.1.root",
			tree:   "tree",
			branch: "SliF64",
			conc:   10,
			want:   []rspan{{beg: 0, end: 5, pos: 3770, sz: 125}},
		},
	} {
		t.Run(tc.fname+"-"+tc.branch, func(t *testing.T) {
			f, err := riofs.Open(tc.fname)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			o, err := f.Get(tc.tree)
			if err != nil {
				t.Fatal(err)
			}

			var (
				tree = o.(Tree)
				b    = tree.Branch(tc.branch)
				beg  = int64(0)
				end  = tree.Entries()
			)

			ra := newBkReader(b, tc.conc, beg, end)
			defer ra.close()

			var got []rspan
			for i := range ra.spans {
				rbk, err := ra.read()
				if err != nil {
					t.Fatalf("could not read basket %d: %+v", i, err)
				}
				got = append(got, rbk.span)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("invalid spans:\ngot= %#v\nwant=%#v", got, tc.want)
			}
		})
	}
}
