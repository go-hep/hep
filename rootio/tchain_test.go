// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio_test

import (
	"fmt"
	"io"
	"reflect"
	"testing"

	"go-hep.org/x/hep/rootio"
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
			fnames:  []string{"testdata/chain.1.root"},
			entries: 10,
			name:    "tree",
			title:   "my tree title",
		},
		{
			// twice the same tree
			fnames:  []string{"testdata/chain.1.root", "testdata/chain.1.root"},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
		},
		{
			// two different trees (with the same schema)
			fnames:  []string{"testdata/chain.1.root", "testdata/chain.2.root"},
			entries: 20,
			name:    "tree",
			title:   "my tree title",
		},
		// TODO(sbinet): add a test with 2 trees with different schemas)
	} {
		t.Run("", func(t *testing.T) {
			files := make([]*rootio.File, len(tc.fnames))
			trees := make([]rootio.Tree, len(tc.fnames))
			for i, fname := range tc.fnames {
				f, err := rootio.Open(fname)
				if err != nil {
					t.Fatalf("could not open ROOT file %q: %v", fname, err)
				}
				defer f.Close()
				files[i] = f

				obj, err := f.Get(tc.name)
				if err != nil {
					t.Fatal(err)
				}

				trees[i] = obj.(rootio.Tree)
			}

			chain := rootio.Chain(trees...)

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

func TestChainScan(t *testing.T) {
	files := []string{
		"testdata/chain.1.root",
		"testdata/chain.2.root", // FIXME(sbinet): implement for >1 tree
	}

	trees := make([]rootio.Tree, len(files))
	for i, fname := range files {
		f, err := rootio.Open(fname)
		if err != nil {
			t.Fatalf("could not open ROOT file %q: %v", fname, err)
		}
		defer f.Close()

		obj, err := f.Get("tree")
		if err != nil {
			t.Fatal(err)
		}

		trees[i] = obj.(rootio.Tree)
	}

	chain := rootio.Chain(trees...)

	type Data struct {
		Event struct {
			Beg       string      `rootio:"Beg"`
			F64       float64     `rootio:"F64"`
			ArrF64    [10]float64 `rootio:"ArrayF64"`
			N         int32       `rootio:"N"`
			SliF64    []float64   `rootio:"SliceF64"`
			StdStr    string      `rootio:"StdStr"`
			StlVecF64 []float64   `rootio:"StlVecF64"`
			StlVecStr []string    `rootio:"StlVecStr"`
			End       string      `rootio:"End"`
		} `rootio:"evt"`
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

	sc, err := rootio.NewTreeScanner(chain, &Data{})
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
	}
	if err := sc.Err(); err != nil && err != io.EOF {
		t.Fatal(err)
	}
}

var SumF64 float64 // global variable to prevent unwanted compiler optimization

func BenchmarkChainTreeScannerStruct(b *testing.B) {
	files := []string{
		"testdata/chain.1.root",
		"testdata/chain.2.root",
		"testdata/chain.1.root",
		"testdata/chain.2.root",
	}

	trees := make([]rootio.Tree, len(files))
	for i, fname := range files {
		f, err := rootio.Open(fname)
		if err != nil {
			b.Fatalf("could not open ROOT file %q: %v", fname, err)
		}
		defer f.Close()

		obj, err := f.Get("tree")
		if err != nil {
			b.Fatal(err)
		}

		trees[i] = obj.(rootio.Tree)
	}
	chain := rootio.Chain(trees...)

	type Data struct {
		Event struct {
			Beg       string      `rootio:"Beg"`
			F64       float64     `rootio:"F64"`
			ArrF64    [10]float64 `rootio:"ArrayF64"`
			N         int32       `rootio:"N"`
			SliF64    []float64   `rootio:"SliceF64"`
			StdStr    string      `rootio:"StdStr"`
			StlVecF64 []float64   `rootio:"StlVecF64"`
			StlVecStr []string    `rootio:"StlVecStr"`
			End       string      `rootio:"End"`
		} `rootio:"evt"`
	}

	sc, err := rootio.NewTreeScanner(chain, &Data{})
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
