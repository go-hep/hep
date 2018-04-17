// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio_test

import (
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
