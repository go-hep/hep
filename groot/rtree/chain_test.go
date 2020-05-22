// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree_test

import (
	"fmt"
	"io"
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

func TestChainScanStruct(t *testing.T) {
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
		evt.SliF64 = make([]float64, evt.N)
		evt.StdStr = fmt.Sprintf("std-%03d", i)
		evt.StlVecF64 = make([]float64, int(evt.N))
		evt.StlVecStr = make([]string, int(evt.N))
		for ii := 0; ii < int(evt.N); ii++ {
			evt.SliF64[ii] = float64(i)
			evt.StlVecF64[ii] = float64(i)
			evt.StlVecStr[ii] = fmt.Sprintf("vec-%03d", i)
		}
		evt.End = fmt.Sprintf("end-%03d", i)

		return data
	}

	sc, err := rtree.NewTreeScanner(chain, &Data{})
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	for sc.Next() {
		var d1 Data
		err := sc.Scan(&d1)
		if err != nil {
			t.Fatal(err)
		}
		i := sc.Entry()
		if !reflect.DeepEqual(d1, want(i)) {
			t.Fatalf("entry[%d]:\ngot= %#v\nwant=%#v\n", i, d1, want(i))
		}

		var d2 Data
		err = sc.Scan(&d2)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(d2, want(i)) {
			t.Fatalf("entry[%d]:\ngot= %#v\nwant=%#v\n", i, d2, want(i))
		}
		total.got++
	}

	if err := sc.Err(); err != nil && err != io.EOF {
		t.Fatal(err)
	}

	if total.got != total.want {
		t.Fatalf("entries scanned differ: got=%d want=%d", total.got, total.want)
	}
}

var SumF64 float64 // global variable to prevent unwanted compiler optimization

func BenchmarkChainTreeScannerStruct(b *testing.B) {
	files := []string{
		"../testdata/chain.1.root",
		"../testdata/chain.2.root",
	}

	trees := make([]rtree.Tree, len(files))
	for i, fname := range files {
		f, err := riofs.Open(fname)
		if err != nil {
			b.Fatalf("could not open ROOT file %q: %v", fname, err)
		}
		defer f.Close()

		obj, err := f.Get("tree")
		if err != nil {
			b.Fatal(err)
		}

		trees[i] = obj.(rtree.Tree)
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

	sc, err := rtree.NewTreeScanner(chain, &Data{})
	if err != nil {
		b.Fatal(err)
	}
	defer sc.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := sc.SeekEntry(0)
		if err != nil {
			b.Fatal(err)
		}
		for sc.Next() {
			var data Data
			err := sc.Scan(&data)
			if err != nil {
				b.Fatal(err)
			}
			SumF64 += data.Event.F64
		}
	}
	if err := sc.Err(); err != nil && err != io.EOF {
		b.Fatal(err)
	}
}

func TestChainSeekEntryTreeScannerPtr(t *testing.T) {
	files := []string{
		"../testdata/chain.flat.1.root",
		"../testdata/chain.flat.2.root",
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
	}
	chain := rtree.Chain(trees...)

	type Data struct {
		F64    float64     `groot:"F64"`
		ArrF64 [10]float64 `groot:"ArrF64"`
		N      int32       `groot:"N"`
		SliF64 []float64   `groot:"SliF64"`
	}

	sc, err := rtree.NewTreeScanner(chain, &Data{})
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	var entries = []int64{0, 1, 4, 2, 5, 6, 7, 2, 7, 3, 0, 9}
	for _, entry := range entries {
		err := sc.SeekEntry(entry)
		if err != nil {
			t.Fatalf("Could not seek to entry %d: %v", entry, err)
		}

		if !sc.Next() {
			t.Fatalf("Could not read entry %d", entry)
		}

		var data Data

		err = sc.Scan(&data)
		if err != nil {
			t.Fatal(err)
		}

		i := sc.Entry()

		if i != entry {
			t.Fatalf("did not seek to entry %d. got=%d, want=%d", entry, i, entry)
		}
		if data.F64 != float64(i) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.F64, float64(i))
		}
		var arr [10]float64
		for ii := range arr {
			arr[ii] = float64(i)
		}
		if data.ArrF64 != arr {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.ArrF64, arr)
		}
		sli := arr[:int(data.N)]
		if !reflect.DeepEqual(sli, data.SliF64) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.SliF64, sli)
		}
	}
}

func TestChainSeekEntryTreeScannerVars(t *testing.T) {
	files := []string{
		"../testdata/chain.flat.1.root",
		"../testdata/chain.flat.2.root",
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
	}
	chain := rtree.Chain(trees...)

	type Data struct {
		F64    float64     `groot:"F64"`
		ArrF64 [10]float64 `groot:"ArrF64"`
		N      int32       `groot:"N"`
		SliF64 []float64   `groot:"SliF64"`
	}

	var data Data
	rvars := []rtree.ReadVar{
		{Name: "F64", Value: &data.F64},
		{Name: "ArrF64", Value: &data.ArrF64},
		{Name: "N", Value: &data.N},
		{Name: "SliF64", Value: &data.SliF64},
	}

	sc, err := rtree.NewTreeScannerVars(chain, rvars...)
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	var entries = []int64{0, 1, 4, 2, 5, 6, 7, 2, 7, 3, 0, 9}
	for _, entry := range entries {
		err := sc.SeekEntry(entry)
		if err != nil {
			t.Fatalf("Could not seek to entry %d: %v", entry, err)
		}

		if !sc.Next() {
			t.Fatalf("Could not read entry %d", entry)
		}

		var data Data

		err = sc.Scan(&data.F64, &data.ArrF64, &data.N, &data.SliF64)
		if err != nil {
			t.Fatal(err)
		}

		i := sc.Entry()

		if i != entry {
			t.Fatalf("did not seek to entry %d. got=%d, want=%d", entry, i, entry)
		}
		if data.F64 != float64(i) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.F64, float64(i))
		}
		var arr [10]float64
		for ii := range arr {
			arr[ii] = float64(i)
		}
		if data.ArrF64 != arr {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.ArrF64, arr)
		}
		sli := arr[:int(data.N)]
		if !reflect.DeepEqual(sli, data.SliF64) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.SliF64, sli)
		}
	}
}

func TestChainSeekEntryScannerPtr(t *testing.T) {
	files := []string{
		"../testdata/chain.flat.1.root",
		"../testdata/chain.flat.2.root",
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
	}
	chain := rtree.Chain(trees...)

	type Data struct {
		F64    float64     `groot:"F64"`
		ArrF64 [10]float64 `groot:"ArrF64"`
		N      int32       `groot:"N"`
		SliF64 []float64   `groot:"SliF64"`
	}

	var data Data
	sc, err := rtree.NewScanner(chain, &data)
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	var entries = []int64{0, 1, 4, 2, 5, 6, 7, 2, 7, 3, 0, 9}
	for _, entry := range entries {
		err := sc.SeekEntry(entry)
		if err != nil {
			t.Fatalf("Could not seek to entry %d: %v", entry, err)
		}

		if !sc.Next() {
			t.Fatalf("Could not read entry %d", entry)
		}

		err = sc.Scan()
		if err != nil {
			t.Fatal(err)
		}

		i := sc.Entry()

		if i != entry {
			t.Fatalf("did not seek to entry %d. got=%d, want=%d", entry, i, entry)
		}
		if data.F64 != float64(i) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.F64, float64(i))
		}
		var arr [10]float64
		for ii := range arr {
			arr[ii] = float64(i)
		}
		if data.ArrF64 != arr {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.ArrF64, arr)
		}
		sli := arr[:int(data.N)]
		if !reflect.DeepEqual(sli, data.SliF64) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.SliF64, sli)
		}
	}
}

func TestChainSeekEntryScannerVars(t *testing.T) {
	files := []string{
		"../testdata/chain.flat.1.root",
		"../testdata/chain.flat.2.root",
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
	}
	chain := rtree.Chain(trees...)

	type Data struct {
		F64    float64     `groot:"F64"`
		ArrF64 [10]float64 `groot:"ArrF64"`
		N      int32       `groot:"N"`
		SliF64 []float64   `groot:"SliF64"`
	}

	var data Data
	rvars := []rtree.ReadVar{
		{Name: "F64", Value: &data.F64},
		{Name: "ArrF64", Value: &data.ArrF64},
		{Name: "N", Value: &data.N},
		{Name: "SliF64", Value: &data.SliF64},
	}

	sc, err := rtree.NewScannerVars(chain, rvars...)
	if err != nil {
		t.Fatal(err)
	}
	defer sc.Close()

	var entries = []int64{0, 1, 4, 2, 5, 6, 7, 2, 7, 3, 0, 9}
	for _, entry := range entries {
		err := sc.SeekEntry(entry)
		if err != nil {
			t.Fatalf("Could not seek to entry %d: %v", entry, err)
		}

		if !sc.Next() {
			t.Fatalf("Could not read entry %d", entry)
		}

		err = sc.Scan()
		if err != nil {
			t.Fatal(err)
		}

		i := sc.Entry()

		if i != entry {
			t.Fatalf("did not seek to entry %d. got=%d, want=%d", entry, i, entry)
		}
		if data.F64 != float64(i) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.F64, float64(i))
		}
		var arr [10]float64
		for ii := range arr {
			arr[ii] = float64(i)
		}
		if data.ArrF64 != arr {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.ArrF64, arr)
		}
		sli := arr[:int(data.N)]
		if !reflect.DeepEqual(sli, data.SliF64) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.SliF64, sli)
		}
	}
}
