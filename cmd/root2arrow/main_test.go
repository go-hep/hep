// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/cmd/root2arrow"

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

func init() {
	_, err := exec.LookPath("arrow-cat")
	if err == nil {
		return
	}

	o, err := exec.Command("go", "install", "github.com/apache/arrow/go/arrow/ipc/cmd/arrow-cat").CombinedOutput()
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
	} {
		t.Run(tc.file, func(t *testing.T) {
			f, err := ioutil.TempFile("", "root2arrow-")
			if err != nil {
				t.Fatal(err)
			}
			f.Close()
			defer os.Remove(f.Name())

			err = process(f.Name(), tc.file, tc.tree, tc.stream)
			if err != nil {
				t.Fatal(err)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}

			got, err := arrowCat(f.Name())
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(got, want) {
				t.Fatalf("arrow file/stream differ")
			}
		})
	}

}

func arrowCat(fname string) ([]byte, error) {
	return exec.Command("arrow-cat", fname).CombinedOutput()
}
