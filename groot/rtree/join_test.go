// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/riofs"
)

func TestJoin(t *testing.T) {
	loadTree := func(fname, tname string) (Tree, func() error) {
		f, err := riofs.Open(fname)
		if err != nil {
			t.Fatalf("could not open %q: %+v", fname, err)
		}
		o, err := riofs.Dir(f).Get(tname)
		if err != nil {
			t.Fatalf("could not get tree %q from %q: %+v", tname, fname, err)
		}

		return o.(Tree), f.Close
	}
	chk := func(f func() error) {
		err := f()
		if err != nil {
			t.Fatal(err)
		}
	}

	j1, close1 := loadTree("../testdata/join1.root", "j1")
	defer chk(close1)

	j2, close2 := loadTree("../testdata/join2.root", "j2")
	defer chk(close2)

	j3, close3 := loadTree("../testdata/join3.root", "j3")
	defer chk(close3)

	j41, close41 := loadTree("../testdata/join4.root", "j41")
	defer chk(close41)

	j42, close42 := loadTree("../testdata/join4.root", "j42")
	defer chk(close42)

	for _, tc := range []struct {
		test  string
		trees []Tree
		name  string
		title string
		nevts int64
		rvars []ReadVar
		brs   []string
		brOK  string
		brNOT string
		lvs   []string
		lvOK  string
		lvNOT string
		err   error
	}{
		{
			test:  "empty",
			trees: nil,
			err:   fmt.Errorf("rtree: no trees to join"),
		},
		{
			test:  "entries-differ-j1-j41",
			trees: []Tree{j1, j41},
			err:   fmt.Errorf("rtree: invalid number of entries in tree j41 (got=11, want=10)"),
		},
		{
			test:  "entries-differ-j41-j1",
			trees: []Tree{j41, j1},
			err:   fmt.Errorf("rtree: invalid number of entries in tree j1 (got=10, want=11)"),
		},
		{
			test:  "branch-collision-1-same-type",
			trees: []Tree{j2, j42},
			err:   fmt.Errorf("rtree: trees j2 and j42 both have a branch named b22"),
		},
		{
			test:  "branch-collision-2-same-type",
			trees: []Tree{j42, j2},
			err:   fmt.Errorf("rtree: trees j42 and j2 both have a branch named b22"),
		},
		{
			test:  "branch-collision-1-diff-type",
			trees: []Tree{j1, j42},
			err:   fmt.Errorf("rtree: trees j1 and j42 both have a branch named b11"),
		},
		{
			test:  "branch-collision-2-diff-type",
			trees: []Tree{j42, j1},
			err:   fmt.Errorf("rtree: trees j42 and j1 both have a branch named b11"),
		},
		{
			test:  "join-j1",
			trees: []Tree{j1},
			nevts: 10,
			name:  "join_j1",
			title: "j1-tree",
			rvars: []ReadVar{
				{Name: "b10", Leaf: "b10", Value: new(float64)},
				{Name: "b11", Leaf: "b11", Value: new(int64)},
				{Name: "b12", Leaf: "b12", Value: new(string)},
			},
			brs: []string{
				"b10", "b11", "b12",
			},
			brOK:  "b10",
			brNOT: "b40",
			lvs: []string{
				"b10", "b11", "b12",
			},
			lvOK:  "b10",
			lvNOT: "b40",
		},
		{
			test:  "join-j2",
			trees: []Tree{j2},
			nevts: 10,
			name:  "join_j2",
			title: "j2-tree",
			rvars: []ReadVar{
				{Name: "b20", Leaf: "b20", Value: new(float64)},
				{Name: "b21", Leaf: "b21", Value: new(int64)},
				{Name: "b22", Leaf: "b22", Value: new(string)},
			},
			brs: []string{
				"b20", "b21", "b22",
			},
			brOK:  "b20",
			brNOT: "b40",
			lvs: []string{
				"b20", "b21", "b22",
			},
			lvOK:  "b20",
			lvNOT: "b40",
		},
		{
			test:  "join-j3",
			trees: []Tree{j3},
			nevts: 10,
			name:  "join_j3",
			title: "j3-tree",
			rvars: []ReadVar{
				{Name: "b30", Leaf: "b30", Value: new(float64)},
				{Name: "b31", Leaf: "b31", Value: new(int64)},
				{Name: "b32", Leaf: "b32", Value: new(string)},
			},
			brs: []string{
				"b30", "b31", "b32",
			},
			brOK:  "b30",
			brNOT: "b40",
			lvs: []string{
				"b30", "b31", "b32",
			},
			lvOK:  "b30",
			lvNOT: "b40",
		},
		{
			test:  "join-j1-j2-j3",
			trees: []Tree{j1, j2, j3},
			nevts: 10,
			name:  "join_j1_j2_j3",
			title: "j1-tree, j2-tree, j3-tree",
			rvars: []ReadVar{
				{Name: "b10", Leaf: "b10", Value: new(float64)},
				{Name: "b11", Leaf: "b11", Value: new(int64)},
				{Name: "b12", Leaf: "b12", Value: new(string)},
				{Name: "b20", Leaf: "b20", Value: new(float64)},
				{Name: "b21", Leaf: "b21", Value: new(int64)},
				{Name: "b22", Leaf: "b22", Value: new(string)},
				{Name: "b30", Leaf: "b30", Value: new(float64)},
				{Name: "b31", Leaf: "b31", Value: new(int64)},
				{Name: "b32", Leaf: "b32", Value: new(string)},
			},
			brs: []string{
				"b10", "b11", "b12",
				"b20", "b21", "b22",
				"b30", "b31", "b32",
			},
			brOK:  "b10",
			brNOT: "b40",
			lvs: []string{
				"b10", "b11", "b12",
				"b20", "b21", "b22",
				"b30", "b31", "b32",
			},
			lvOK:  "b10",
			lvNOT: "b40",
		},
		{
			test:  "join-j1-j2",
			trees: []Tree{j1, j2},
			nevts: 10,
			name:  "join_j1_j2",
			title: "j1-tree, j2-tree",
			rvars: []ReadVar{
				{Name: "b10", Leaf: "b10", Value: new(float64)},
				{Name: "b11", Leaf: "b11", Value: new(int64)},
				{Name: "b12", Leaf: "b12", Value: new(string)},
				{Name: "b20", Leaf: "b20", Value: new(float64)},
				{Name: "b21", Leaf: "b21", Value: new(int64)},
				{Name: "b22", Leaf: "b22", Value: new(string)},
			},
			brs: []string{
				"b10", "b11", "b12",
				"b20", "b21", "b22",
			},
			brOK:  "b10",
			brNOT: "b30",
			lvs: []string{
				"b10", "b11", "b12",
				"b20", "b21", "b22",
			},
			lvOK:  "b10",
			lvNOT: "b30",
		},
		{
			test:  "join-j2-j1",
			trees: []Tree{j2, j1},
			nevts: 10,
			name:  "join_j2_j1",
			title: "j2-tree, j1-tree",
			rvars: []ReadVar{
				{Name: "b20", Leaf: "b20", Value: new(float64)},
				{Name: "b21", Leaf: "b21", Value: new(int64)},
				{Name: "b22", Leaf: "b22", Value: new(string)},
				{Name: "b10", Leaf: "b10", Value: new(float64)},
				{Name: "b11", Leaf: "b11", Value: new(int64)},
				{Name: "b12", Leaf: "b12", Value: new(string)},
			},
			brs: []string{
				"b20", "b21", "b22",
				"b10", "b11", "b12",
			},
			brOK:  "b10",
			brNOT: "b30",
			lvs: []string{
				"b20", "b21", "b22",
				"b10", "b11", "b12",
			},
			lvOK:  "b10",
			lvNOT: "b30",
		},
		{
			test:  "join-j1-j3",
			trees: []Tree{j1, j3},
			nevts: 10,
			name:  "join_j1_j3",
			title: "j1-tree, j3-tree",
			rvars: []ReadVar{
				{Name: "b10", Leaf: "b10", Value: new(float64)},
				{Name: "b11", Leaf: "b11", Value: new(int64)},
				{Name: "b12", Leaf: "b12", Value: new(string)},
				{Name: "b30", Leaf: "b30", Value: new(float64)},
				{Name: "b31", Leaf: "b31", Value: new(int64)},
				{Name: "b32", Leaf: "b32", Value: new(string)},
			},
			brs: []string{
				"b10", "b11", "b12",
				"b30", "b31", "b32",
			},
			brOK:  "b10",
			brNOT: "b40",
			lvs: []string{
				"b10", "b11", "b12",
				"b30", "b31", "b32",
			},
			lvOK:  "b10",
			lvNOT: "b40",
		},
		{
			test:  "join-j2-j3",
			trees: []Tree{j2, j3},
			nevts: 10,
			name:  "join_j2_j3",
			title: "j2-tree, j3-tree",
			rvars: []ReadVar{
				{Name: "b20", Leaf: "b20", Value: new(float64)},
				{Name: "b21", Leaf: "b21", Value: new(int64)},
				{Name: "b22", Leaf: "b22", Value: new(string)},
				{Name: "b30", Leaf: "b30", Value: new(float64)},
				{Name: "b31", Leaf: "b31", Value: new(int64)},
				{Name: "b32", Leaf: "b32", Value: new(string)},
			},
			brs: []string{
				"b20", "b21", "b22",
				"b30", "b31", "b32",
			},
			brOK:  "b20",
			brNOT: "b10",
			lvs: []string{
				"b20", "b21", "b22",
				"b30", "b31", "b32",
			},
			lvOK:  "b20",
			lvNOT: "b10",
		},
	} {
		t.Run(tc.test, func(t *testing.T) {
			tree, err := Join(tc.trees...)
			switch {
			case err != nil && tc.err != nil:
				if got, want := err.Error(), tc.err.Error(); got != want {
					t.Fatalf("invalid error:\ngot= %s\nwant=%s", got, want)
				}
				return
			case err != nil && tc.err == nil:
				t.Fatalf("could not join trees: %+v", err)
			case err == nil && tc.err != nil:
				t.Fatalf("invalid error:\ngot= %v\nwant: %s", err, tc.err.Error())
			case err == nil && tc.err == nil:
				// ok.
			}

			if got, want := tree.Class(), "TJoin"; got != want {
				t.Fatalf("invalid class:\ngot= %q\nwant=%q", got, want)
			}

			if got, want := tree.Name(), tc.name; got != want {
				t.Fatalf("invalid name:\ngot= %q\nwant=%q", got, want)
			}

			if got, want := tree.Title(), tc.title; got != want {
				t.Fatalf("invalid title:\ngot= %q\nwant=%q", got, want)
			}

			if got, want := tree.Entries(), tc.nevts; got != want {
				t.Fatalf("invalid entries: got=%d, want=%d", got, want)
			}
			{
				rvars := NewReadVars(tree)
				n := len(tc.rvars)
				if len(rvars) < n {
					n = len(rvars)
				}

				for i := 0; i < n; i++ {
					got := rvars[i]
					want := tc.rvars[i]
					if got.Name != want.Name {
						t.Fatalf("invalid rvar-name[%d]: got=%q, want=%q", i, got.Name, want.Name)
					}
					if got.Leaf != want.Leaf {
						t.Fatalf("invalid rvar-leaf[%d]: got=%q, want=%q", i, got.Leaf, want.Leaf)
					}
					if got, want := reflect.TypeOf(got.Value), reflect.TypeOf(want.Value); got != want {
						t.Fatalf("invalid rvar (name=%q) type[%d]: got=%v, want=%v", rvars[i].Name, i, got, want)
					}
				}

				if got, want := len(rvars), len(tc.rvars); got != want {
					t.Fatalf("invalid number of rvars: got=%d, want=%d", got, want)
				}
			}
			{
				brs := tree.Branches()
				n := len(tc.brs)
				if len(brs) < n {
					n = len(brs)
				}

				for i := 0; i < n; i++ {
					if got, want := brs[i].Name(), tc.brs[i]; got != want {
						t.Fatalf("invalid branch name[%d]: got=%q, want=%q", i, got, want)
					}
				}

				if got, want := len(brs), len(tc.brs); got != want {
					t.Fatalf("invalid number of branches: got=%d, want=%d", got, want)
				}

				br := tree.Branch(tc.brOK)
				if br == nil {
					t.Fatalf("could not retrieve branch %q", tc.brOK)
				}
				if got, want := br.Name(), tc.brOK; got != want {
					t.Fatalf("invalid name for branch-ok: got=%q, want=%q", got, want)
				}

				br = tree.Branch(tc.brNOT)
				if br != nil {
					t.Fatalf("unexpected branch for branch-not (%s): got=%#v", tc.brNOT, br)
				}
			}
			{
				lvs := tree.Leaves()
				n := len(tc.lvs)
				if len(lvs) < n {
					n = len(lvs)
				}

				for i := 0; i < n; i++ {
					if got, want := lvs[i].Name(), tc.lvs[i]; got != want {
						t.Fatalf("invalid leaf name[%d]: got=%q, want=%q", i, got, want)
					}
				}

				if got, want := len(lvs), len(tc.lvs); got != want {
					t.Fatalf("invalid number of leaves: got=%d, want=%d", got, want)
				}

				lv := tree.Leaf(tc.lvOK)
				if lv == nil {
					t.Fatalf("could not retrieve leaf %q", tc.lvOK)
				}
				if got, want := lv.Name(), tc.lvOK; got != want {
					t.Fatalf("invalid name for leaf-ok: got=%q, want=%q", got, want)
				}
				br := lv.Branch()
				if br == nil || br.Name() != tc.lvOK {
					t.Fatalf("invalid leaf-branch: ptr-ok=%v", br != nil)
				}

				lv = tree.Leaf(tc.lvNOT)
				if lv != nil {
					t.Fatalf("unexpected leaf for leaf-not (%s): got=%#v", tc.lvNOT, lv)
				}
			}
		})
	}
}
