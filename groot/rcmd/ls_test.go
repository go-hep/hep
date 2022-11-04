// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rcmd_test

import (
	"os"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/rcmd"
	"go-hep.org/x/hep/internal/diff"
)

func TestList(t *testing.T) {
	loadRef := func(fname string) string {
		t.Helper()
		raw, err := os.ReadFile(fname)
		if err != nil {
			t.Fatalf("could not load reference file %q: %+v", fname, err)
		}
		return string(raw)
	}

	opts := []rcmd.ListOption{
		rcmd.ListStreamers(true),
		rcmd.ListTrees(true),
	}

	for _, tc := range []struct {
		name string
		opts []rcmd.ListOption
		want string
	}{
		{
			name: "../testdata/simple.root",
			want: `=== [../testdata/simple.root] ===
version: 60600
TTree   tree    fake data (cycle=1)
`,
		},
		{
			name: "../testdata/simple.root",
			opts: []rcmd.ListOption{
				rcmd.ListTrees(true),
			},
			want: `=== [../testdata/simple.root] ===
version: 60600
  TTree   tree      fake data (entries=4)
    one   "one/I"   TBranch
    two   "two/F"   TBranch
    three "three/C" TBranch
`,
		},
		{
			name: "../testdata/simple.root",
			opts: opts,
			want: loadRef("./testdata/simple.root-ls.txt"),
		},
		{
			name: "../testdata/graphs.root",
			opts: opts,
			want: loadRef("./testdata/graphs.root-ls.txt"),
		},
		{
			name: "../testdata/small-flat-tree.root",
			opts: opts,
			want: loadRef("./testdata/small-flat-tree.root-ls.txt"),
		},
		{
			name: "../testdata/small-evnt-tree-fullsplit.root",
			opts: opts,
			want: loadRef("./testdata/small-evnt-tree-fullsplit.root-ls.txt"),
		},
		{
			name: "../testdata/small-evnt-tree-nosplit.root",
			opts: opts,
			want: loadRef("./testdata/small-evnt-tree-nosplit.root-ls.txt"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := new(strings.Builder)
			err := rcmd.List(out, tc.name, tc.opts...)
			if err != nil {
				t.Fatalf("could not run root-ls: %+v", err)
			}

			if got, want := out.String(), tc.want; got != want {
				t.Fatalf("invalid root-ls output:\n%s", diff.Format(got, want))
			}
		})
	}
}
