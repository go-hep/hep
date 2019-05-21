// Copyright 2019 The go-hep Authors. All rights reserved.
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

func TestTable(t *testing.T) {
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
			tbl := NewTable(tree, WithAllocator(mem))
			defer tbl.Release()

			tbl.Retain()
			tbl.Release()

			tr := array.NewTableReader(tbl, -1)
			defer tr.Release()

			recs := 0
			out := new(strings.Builder)
			fmt.Fprintf(out, "file: %s\n", tc.file)
			for tr.Next() {
				rec := tr.Record()
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
