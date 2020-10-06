// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rarrow // import "go-hep.org/x/hep/groot/rarrow"

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/memory"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func TestRecord(t *testing.T) {
	for _, tc := range []struct {
		file string
		tree string
		want string
	}{
		{
			file: "../testdata/simple.root",
			tree: "tree",
			want: "testdata/simple.root.txt",
		},
		{
			file: "../testdata/small-flat-tree.root",
			tree: "tree",
			want: "testdata/small-flat-tree.root.txt",
		},
		{
			file: "../testdata/small-evnt-tree-fullsplit.root",
			tree: "tree",
			want: "testdata/small-evnt-tree-fullsplit.root.txt",
		},
		{
			file: "../testdata/small-evnt-tree-nosplit.root",
			tree: "tree",
			want: "testdata/small-evnt-tree-nosplit.root.txt",
		},
		{
			// n-dim arrays
			// FIXME(sbinet): arrays of Float16_t and Double32_t are flatten.
			// This is because of:
			// https://sft.its.cern.ch/jira/browse/ROOT-10149
			file: "../testdata/ndim.root",
			tree: "tree",
			want: "testdata/ndim.root.txt",
		},
		{
			// slice of n-dim arrays
			// FIXME(sbinet): arrays of Float16_t and Double32_t are flatten.
			// This is because of:
			// https://sft.its.cern.ch/jira/browse/ROOT-10149
			file: "../testdata/ndim-slice.root",
			tree: "tree",
			want: "testdata/ndim-slice.root.txt",
		},
	} {
		t.Run(tc.file, func(t *testing.T) {
			f, err := groot.Open(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			o, err := riofs.Dir(f).Get(tc.tree)
			if err != nil {
				t.Fatal(err)
			}

			mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
			defer mem.AssertSize(t, 0)

			tree := o.(rtree.Tree)
			rec := NewRecord(tree, WithAllocator(mem))
			defer rec.Release()

			if got, want := rec.NumCols(), int64(len(tree.Branches())); got != want {
				t.Fatalf("invalid number of columns: got=%d, want=%d", got, want)
			}

			if got, want := rec.NumRows(), tree.Entries(); got != want {
				t.Fatalf("invalid number of rows: got=%d, want=%d", got, want)
			}

			for i, branch := range tree.Branches() {
				col := rec.Column(i)
				if got, want := int64(col.Len()), rec.NumRows(); got != want {
					t.Fatalf("invalid column size[%d]: got=%d, want=%d", i, got, want)
				}
				name := rec.ColumnName(i)
				if got, want := name, branch.Name(); got != want {
					t.Fatalf("invalid column name[%d]: got=%q, want=%q", i, got, want)
				}
			}

			rec.Retain()
			rec.Release()

			for _, tc := range []struct{ beg, end, want int64 }{
				{0, rec.NumRows(), rec.NumRows()},
				{0, rec.NumRows() - 1, rec.NumRows() - 1},
				{0, 0, 0},
				{0, 1, 1},
				{0, 2, 2},
				{1, 2, 1},
			} {
				t.Run(fmt.Sprintf("slice-%d-%d", tc.beg, tc.end), func(t *testing.T) {
					sub := rec.NewSlice(tc.beg, tc.end)
					defer sub.Release()

					if got, want := sub.NumCols(), rec.NumCols(); got != want {
						t.Fatalf("invalid number of sub-cols: got=%d, want=%d", got, want)
					}

					if got, want := sub.NumRows(), tc.want; got != want {
						t.Fatalf("invalid number of sub-rows: got=%d, want=%d", got, want)
					}
				})
			}

			rr, err := array.NewRecordReader(rec.Schema(), []array.Record{rec})
			if err != nil {
				t.Fatal(err)
			}
			defer rr.Release()

			recs := 0
			out := new(strings.Builder)
			fmt.Fprintf(out, "file: %s\n", tc.file)
			for rr.Next() {
				rec := rr.Record()
				for i, col := range rec.Columns() {
					fmt.Fprintf(out, "rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
				}
				recs++
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := out.String(), string(want); got != want {
				t.Fatalf("invalid table\ngot:\n%s\nwant:\n%s\n", got, want)
			}
		})
	}
}

func TestRecordReader(t *testing.T) {
	for _, tc := range []struct {
		file  string
		tree  string
		chunk int64
		want  string
	}{
		{
			file:  "../testdata/leaves.root",
			tree:  "tree",
			chunk: -1,
			want:  "testdata/leaves.root.txt",
		},
		{
			file:  "../testdata/simple.root",
			tree:  "tree",
			chunk: -1,
			want:  "testdata/simple.root.txt",
		},
		{
			file:  "../testdata/simple.root",
			tree:  "tree",
			chunk: 0,
			want:  "testdata/simple.root.chunk=1.txt",
		},
		{
			file:  "../testdata/simple.root",
			tree:  "tree",
			chunk: 1,
			want:  "testdata/simple.root.chunk=1.txt",
		},
		{
			file:  "../testdata/simple.root",
			tree:  "tree",
			chunk: 2,
			want:  "testdata/simple.root.chunk=2.txt",
		},
		{
			file:  "../testdata/simple.root",
			tree:  "tree",
			chunk: 3,
			want:  "testdata/simple.root.chunk=3.txt",
		},
		{
			file:  "../testdata/simple.root",
			tree:  "tree",
			chunk: 4,
			want:  "testdata/simple.root.txt",
		},
		{
			file:  "../testdata/small-flat-tree.root",
			tree:  "tree",
			chunk: -1,
			want:  "testdata/small-flat-tree.root.txt",
		},
		{
			file:  "../testdata/small-evnt-tree-fullsplit.root",
			tree:  "tree",
			chunk: -1,
			want:  "testdata/small-evnt-tree-fullsplit.root.txt",
		},
		{
			file:  "../testdata/small-evnt-tree-nosplit.root",
			tree:  "tree",
			chunk: -1,
			want:  "testdata/small-evnt-tree-nosplit.root.txt",
		},
	} {
		t.Run(fmt.Sprintf("%s-with-chunk=%d", tc.file, tc.chunk), func(t *testing.T) {
			f, err := groot.Open(tc.file)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()

			o, err := riofs.Dir(f).Get(tc.tree)
			if err != nil {
				t.Fatal(err)
			}

			mem := memory.NewCheckedAllocator(memory.NewGoAllocator())
			defer mem.AssertSize(t, 0)

			tree := o.(rtree.Tree)
			rr := NewRecordReader(tree, WithAllocator(mem), WithChunk(tc.chunk))
			defer rr.Release()

			rr.Retain()
			rr.Release()

			recs := 0
			out := new(strings.Builder)
			fmt.Fprintf(out, "file: %s\n", tc.file)
			for rr.Next() {
				rec := rr.Record()
				for i, col := range rec.Columns() {
					fmt.Fprintf(out, "rec[%d][%s]: %v\n", recs, rec.Schema().Field(i).Name, col)
				}
				recs++
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			if got, want := out.String(), string(want); got != want {
				t.Fatalf("invalid table\ngot:\n%s\nwant:\n%s\n", got, want)
			}
		})
	}
}
