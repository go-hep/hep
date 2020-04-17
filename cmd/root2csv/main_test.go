// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/cmd/root2csv"

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestROOT2CSV(t *testing.T) {
	for _, tc := range []struct {
		file string
		tree string
		want string
		skip bool
	}{
		{
			file: "../../groot/testdata/simple.root",
			tree: "tree",
			want: "testdata/simple.root.csv",
		},
		{
			file: "../../groot/testdata/leaves.root",
			tree: "tree",
			want: "testdata/leaves.root.csv",
		},
		{
			file: "../../groot/testdata/small-evnt-tree-fullsplit.root",
			tree: "tree",
			want: "testdata/small-evnt-tree-fullsplit.root.csv",
		},
		{
			file: "../../groot/testdata/small-evnt-tree-nosplit.root",
			tree: "tree",
			want: "testdata/small-evnt-tree-nosplit.root.csv",
			skip: true, // FIXME(sbinet)
		},
	} {
		t.Run(tc.file, func(t *testing.T) {
			if tc.skip {
				t.Skipf("not ready (FIXME)")
			}

			f, err := ioutil.TempFile("", "root2csv-")
			if err != nil {
				t.Fatal(err)
			}
			f.Close()
			defer os.Remove(f.Name())

			err = process(f.Name(), tc.file, tc.tree)
			if err != nil {
				t.Fatal(err)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			got, err := ioutil.ReadFile(f.Name())
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got, want) {
				t.Fatalf("CSV files differ")
			}
		})
	}

}
