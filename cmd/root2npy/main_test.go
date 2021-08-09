// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sbinet/npyio"
)

func TestProcess(t *testing.T) {
	loadRef := func(fname string) string {
		t.Helper()
		raw, err := os.ReadFile(fname)
		if err != nil {
			t.Fatalf("could not load reference file %q: %+v", fname, err)
		}
		return string(raw)
	}

	tmp, err := os.MkdirTemp("", "root2npy-")
	if err != nil {
		t.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, tc := range []struct {
		name string
		tree string
		want string
	}{
		{
			name: "../../groot/testdata/simple.root",
			tree: "tree",
			want: loadRef("testdata/simple.root.txt"),
		},
		{
			name: "../../groot/testdata/leaves.root",
			tree: "tree",
			want: loadRef("testdata/leaves.root.txt"),
		},
		{
			name: "../../groot/testdata/ndim.root",
			tree: "tree",
			want: loadRef("testdata/ndim.root.txt"),
		},
		{
			name: "../../groot/testdata/small-flat-tree.root",
			tree: "tree",
			want: loadRef("testdata/small-flat-tree.root.txt"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oname := filepath.Join(tmp, filepath.Base(tc.name)+".npz")
			err := process(oname, tc.name, tc.tree)
			if err != nil {
				t.Fatalf("could not run root2npy: %+v", err)
			}

			f, err := os.Open(oname)
			if err != nil {
				t.Fatalf("could not open %q: %+v", oname, err)
			}
			defer f.Close()

			type namer interface{ Name() string }

			r := struct {
				io.ReaderAt
				io.Seeker
				namer
			}{
				ReaderAt: f,
				Seeker:   f,
				namer:    nilNamer{},
			}

			var got strings.Builder
			err = npyio.Dump(&got, r)
			if err != nil {
				t.Fatalf("could not read output file")
			}

			if got, want := got.String(), tc.want; got != want {
				t.Fatalf("invalid npy:\ngot:\n%s\nwant:\n%s\n", got, want)
			}
		})
	}
}

type nilNamer struct{}

func (nilNamer) Name() string { return "output.npz" }
