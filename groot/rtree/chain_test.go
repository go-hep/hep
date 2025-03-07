// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/internal/rtests"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/rtree"
)

func TestChain(t *testing.T) {
	for _, tc := range []struct {
		fnames  []string
		entries int64
		name    string
		title   string
		brs     []string
		brOK    string
		brNOT   string
		lvs     []string
		lvOK    string
		lvNOT   string
	}{
		{
			fnames:  nil,
			entries: 0,
			name:    "",
			title:   "",
		},
		{
			fnames:  []string{"../testdata/chain.1.root"},
			entries: 10,
			name:    "tree",
			title:   "my tree title",
			brs:     []string{"evt"},
			brOK:    "evt",
			brNOT:   "foo",
			lvs:     []string{"evt"},
			lvOK:    "evt",
			lvNOT:   "foo",
		},
		{
			fnames:  []string{rtests.XrdRemote("testdata/chain.1.root")},
			entries: 10,
			name:    "tree",
			title:   "my tree title",
			brs:     []string{"evt"},
			brOK:    "evt",
			brNOT:   "foo",
			lvs:     []string{"evt"},
			lvOK:    "evt",
			lvNOT:   "foo",
		},
		{
			// twice the same tree
			fnames:  []string{"../testdata/chain.1.root", "../testdata/chain.1.root"},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
			brs:     []string{"evt"},
			brOK:    "evt",
			brNOT:   "foo",
			lvs:     []string{"evt"},
			lvOK:    "evt",
			lvNOT:   "foo",
		},
		{
			// twice the same tree
			fnames: []string{
				rtests.XrdRemote("testdata/chain.1.root"),
				rtests.XrdRemote("testdata/chain.1.root"),
			},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
			brs:     []string{"evt"},
			brOK:    "evt",
			brNOT:   "foo",
			lvs:     []string{"evt"},
			lvOK:    "evt",
			lvNOT:   "foo",
		},
		{
			// two different trees (with the same schema)
			fnames:  []string{"../testdata/chain.1.root", "../testdata/chain.2.root"},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
			brs:     []string{"evt"},
			brOK:    "evt",
			brNOT:   "foo",
			lvs:     []string{"evt"},
			lvOK:    "evt",
			lvNOT:   "foo",
		},
		{
			// two different trees (with the same schema)
			fnames: []string{
				rtests.XrdRemote("testdata/chain.1.root"),
				rtests.XrdRemote("testdata/chain.2.root"),
			},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
			brs:     []string{"evt"},
			brOK:    "evt",
			brNOT:   "foo",
			lvs:     []string{"evt"},
			lvOK:    "evt",
			lvNOT:   "foo",
		},
		{
			// two different (flat) trees (with the same schema)
			fnames: []string{
				"../testdata/chain.flat.1.root",
				"../testdata/chain.flat.2.root",
			},
			entries: 10,
			name:    "tree",
			title:   "my tree title",
			brs: []string{
				"B",
				"Str",
				"I8", "I16", "I32", "I64",
				"U8", "U16", "U32", "U64",
				"F32", "F64",
				"ArrBs",
				"ArrI8", "ArrI16", "ArrI32", "ArrI64",
				"ArrU8", "ArrU16", "ArrU32", "ArrU64",
				"ArrF32", "ArrF64",
				"N",
				"SliBs",
				"SliI8", "SliI16", "SliI32", "SliI64",
				"SliU8", "SliU16", "SliU32", "SliU64",
				"SliF32", "SliF64",
			},
			brOK:  "N",
			brNOT: "foo",
			lvs: []string{
				"B",
				"Str",
				"I8", "I16", "I32", "I64",
				"U8", "U16", "U32", "U64",
				"F32", "F64",
				"ArrBs",
				"ArrI8", "ArrI16", "ArrI32", "ArrI64",
				"ArrU8", "ArrU16", "ArrU32", "ArrU64",
				"ArrF32", "ArrF64",
				"N",
				"SliBs",
				"SliI8", "SliI16", "SliI32", "SliI64",
				"SliU8", "SliU16", "SliU32", "SliU64",
				"SliF32", "SliF64",
			},
			lvOK:  "N",
			lvNOT: "foo",
		},
		// TODO(sbinet): add a test with 2 trees with different schemas)
	} {
		t.Run("", func(t *testing.T) {
			files := make([]*riofs.File, len(tc.fnames))
			trees := make([]rtree.Tree, len(tc.fnames))
			for i, fname := range tc.fnames {
				f, err := riofs.Open(fname)
				if err != nil {
					t.Fatalf("could not open ROOT file %q: %v", fname, err)
				}
				defer f.Close()
				files[i] = f

				obj, err := f.Get(tc.name)
				if err != nil {
					t.Fatal(err)
				}

				trees[i] = obj.(rtree.Tree)
			}

			chain := rtree.Chain(trees...)

			if got, want := chain.Class(), "TChain"; got != want {
				t.Fatalf("class name differ\ngot = %q, want= %q", got, want)
			}
			if got, want := chain.Name(), tc.name; got != want {
				t.Fatalf("names differ\ngot = %q, want= %q", got, want)
			}
			if got, want := chain.Title(), tc.title; got != want {
				t.Fatalf("titles differ\ngot = %q, want= %q", got, want)
			}
			if got, want := chain.Entries(), tc.entries; got != want {
				t.Fatalf("titles differ\ngot = %v, want= %v", got, want)
			}
			{
				brs := chain.Branches()
				n := min(len(brs), len(tc.brs))

				for i := range n {
					if got, want := brs[i].Name(), tc.brs[i]; got != want {
						t.Fatalf("invalid branch name[%d]: got=%q, want=%q", i, got, want)
					}
				}

				if got, want := len(brs), len(tc.brs); got != want {
					t.Fatalf("invalid number of branches: got=%d, want=%d", got, want)
				}

				if tc.brOK != "" {
					br := chain.Branch(tc.brOK)
					if br == nil {
						t.Fatalf("could not retrieve branch %q", tc.brOK)
					}
					if got, want := br.Name(), tc.brOK; got != want {
						t.Fatalf("invalid name for branch-ok: got=%q, want=%q", got, want)
					}
				}

				br := chain.Branch(tc.brNOT)
				if br != nil {
					t.Fatalf("unexpected branch for branch-not (%s): got=%#v", tc.brNOT, br)
				}
			}
			{
				lvs := chain.Leaves()
				n := min(len(lvs), len(tc.lvs))

				for i := range n {
					if got, want := lvs[i].Name(), tc.lvs[i]; got != want {
						t.Fatalf("invalid leaf name[%d]: got=%q, want=%q", i, got, want)
					}
				}

				if got, want := len(lvs), len(tc.lvs); got != want {
					t.Fatalf("invalid number of leaves: got=%d, want=%d", got, want)
				}

				if tc.lvOK != "" {
					lv := chain.Leaf(tc.lvOK)
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
				}
				lv := chain.Leaf(tc.lvNOT)
				if lv != nil {
					t.Fatalf("unexpected leaf for leaf-not (%s): got=%#v", tc.lvNOT, lv)
				}
			}
		})
	}
}

func TestChainOf(t *testing.T) {
	for _, tc := range []struct {
		fnames  []string
		entries int64
		name    string
		title   string
	}{
		{
			fnames:  nil,
			entries: 0,
			name:    "",
			title:   "",
		},
		{
			fnames:  []string{"../testdata/chain.1.root"},
			entries: 10,
			name:    "tree",
			title:   "my tree title",
		},
		{
			fnames:  []string{rtests.XrdRemote("testdata/chain.1.root")},
			entries: 10,
			name:    "tree",
			title:   "my tree title",
		},
		{
			// twice the same tree
			fnames:  []string{"../testdata/chain.1.root", "../testdata/chain.1.root"},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
		},
		{
			// twice the same tree
			fnames: []string{
				rtests.XrdRemote("testdata/chain.1.root"),
				rtests.XrdRemote("testdata/chain.1.root"),
			},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
		},
		{
			// two different trees (with the same schema)
			fnames:  []string{"../testdata/chain.1.root", "../testdata/chain.2.root"},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
		},
		{
			// two different trees (with the same schema)
			fnames: []string{
				rtests.XrdRemote("testdata/chain.1.root"),
				rtests.XrdRemote("testdata/chain.2.root"),
			},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
		},
		// TODO(sbinet): add a test with 2 trees with different schemas)
	} {
		t.Run("", func(t *testing.T) {
			chain, closer, err := rtree.ChainOf(tc.name, tc.fnames...)
			if err != nil {
				t.Fatalf("could not create chain: %v", err)
			}
			defer func() {
				_ = closer()
			}()

			if got, want := chain.Name(), tc.name; got != want {
				t.Fatalf("names differ\ngot = %q, want= %q", got, want)
			}
			if got, want := chain.Title(), tc.title; got != want {
				t.Fatalf("titles differ\ngot = %q, want= %q", got, want)
			}
			if got, want := chain.Entries(), tc.entries; got != want {
				t.Fatalf("titles differ\ngot = %v, want= %v", got, want)
			}
		})
	}
}

func TestChainReaderStruct(t *testing.T) {
	files := []string{
		"../testdata/chain.1.root",
		"../testdata/chain.2.root",
	}
	var total struct {
		got, want int64
	}
	trees := make([]rtree.Tree, len(files))
	for i, fname := range files {
		f, err := riofs.Open(fname)
		if err != nil {
			t.Fatalf("could not open ROOT file %q: %v", fname, err)
		}
		defer f.Close()

		obj, err := f.Get("tree")
		if err != nil {
			t.Fatal(err)
		}

		trees[i] = obj.(rtree.Tree)
		total.want += trees[i].Entries()
	}

	chain := rtree.Chain(trees...)

	type Data struct {
		Event struct {
			Beg       string      `groot:"Beg"`
			F64       float64     `groot:"F64"`
			ArrF64    [10]float64 `groot:"ArrayF64"`
			N         int32       `groot:"N"`
			SliF64    []float64   `groot:"SliceF64"`
			StdStr    string      `groot:"StdStr"`
			StlVecF64 []float64   `groot:"StlVecF64"`
			StlVecStr []string    `groot:"StlVecStr"`
			End       string      `groot:"End"`
		} `groot:"evt"`
	}

	want := func(i int64) (data Data) {
		evt := &data.Event
		evt.Beg = fmt.Sprintf("beg-%03d", i)
		evt.F64 = float64(i)
		for j := range evt.ArrF64 {
			evt.ArrF64[j] = float64(i)
		}
		evt.N = int32(i) % 10
		evt.StdStr = fmt.Sprintf("std-%03d", i)
		switch i {
		case 0:
			evt.SliF64 = nil
			evt.StlVecF64 = nil
			evt.StlVecStr = nil
		default:
			evt.SliF64 = make([]float64, evt.N)
			evt.StlVecF64 = make([]float64, int(evt.N))
			evt.StlVecStr = make([]string, int(evt.N))
		}
		for ii := range int(evt.N) {
			evt.SliF64[ii] = float64(i)
			evt.StlVecF64[ii] = float64(i)
			evt.StlVecStr[ii] = fmt.Sprintf("vec-%03d", i)
		}
		evt.End = fmt.Sprintf("end-%03d", i)

		return data
	}

	var data Data
	r, err := rtree.NewReader(chain, rtree.ReadVarsFromStruct(&data))
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error {
		i := ctx.Entry
		if !reflect.DeepEqual(data, want(i)) {
			return fmt.Errorf("entry[%d]:\ngot= %#v\nwant=%#v\n", i, data, want(i))
		}
		total.got++
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	if total.got != total.want {
		t.Fatalf("entries differ: got=%d want=%d", total.got, total.want)
	}
}
