// Copyright 2018 The go-hep Authors. All rights reserved.
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
			defer closer()

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
	svars := []rtree.ScanVar{
		{Name: "F64", Value: &data.F64},
		{Name: "ArrF64", Value: &data.ArrF64},
		{Name: "N", Value: &data.N},
		{Name: "SliF64", Value: &data.SliF64},
	}

	sc, err := rtree.NewTreeScannerVars(chain, svars...)
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
		if data.N == 0 {
			sli = nil
		}
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
	svars := []rtree.ScanVar{
		{Name: "F64", Value: &data.F64},
		{Name: "ArrF64", Value: &data.ArrF64},
		{Name: "N", Value: &data.N},
		{Name: "SliF64", Value: &data.SliF64},
	}

	sc, err := rtree.NewScannerVars(chain, svars...)
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
		if data.N == 0 {
			sli = nil
		}
		if !reflect.DeepEqual(sli, data.SliF64) {
			t.Fatalf("entry [%d] : got= %#v want=%#v\n", i, data.SliF64, sli)
		}
	}
}
