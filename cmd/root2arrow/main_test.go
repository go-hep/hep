// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/cmd/root2arrow"

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func init() {
	_, err := exec.LookPath("arrow-cat")
	if err == nil {
		return
	}

	o, err := exec.Command("go", "install", "git.sr.ht/~sbinet/go-arrow/ipc/cmd/arrow-cat").CombinedOutput()
	if err != nil {
		panic(fmt.Errorf("could not install arrow-cat command:\n%v\nerr: %w", string(o), err))
	}

}

func TestFile(t *testing.T) {
	for _, tc := range []struct {
		file   string
		tree   string
		stream bool
		want   string
	}{
		{
			file: "../../groot/testdata/simple.root",
			tree: "tree",
			want: "testdata/simple.root.file",
		},
		{
			file:   "../../groot/testdata/simple.root",
			tree:   "tree",
			stream: true,
			want:   "testdata/simple.root.stream",
		},
		{
			file: "../../groot/testdata/leaves.root",
			tree: "tree",
			want: "testdata/leaves.root.file",
		},
		{
			file:   "../../groot/testdata/leaves.root",
			tree:   "tree",
			stream: true,
			want:   "testdata/leaves.root.stream",
		},
		{
			file: "../../groot/testdata/embedded-std-vector.root",
			tree: "modules",
			want: "testdata/embedded-std-vector.root.file",
		},
		{
			file:   "../../groot/testdata/embedded-std-vector.root",
			tree:   "modules",
			stream: true,
			want:   "testdata/embedded-std-vector.root.stream",
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			f, err := os.CreateTemp("", "root2arrow-")
			if err != nil {
				t.Fatal(err)
			}
			f.Close()
			defer os.Remove(f.Name())

			err = process(f.Name(), tc.file, tc.tree, tc.stream)
			if err != nil {
				t.Fatal(err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			got, err := arrowCat(f.Name())
			if err != nil {
				t.Fatal(err)
			}

			if got, want := string(got), string(want); got != want {
				diff := cmp.Diff(want, got)
				t.Fatalf(
					"arrow file/stream differ: -- (-ref +got)\n%s",
					diff,
				)
			}
		})
	}

}

func arrowCat(fname string) ([]byte, error) {
	return exec.Command("arrow-cat", fname).CombinedOutput()
}
