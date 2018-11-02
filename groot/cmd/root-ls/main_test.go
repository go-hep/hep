// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestROOTls(t *testing.T) {
	for _, name := range []string{
		"../../testdata/dirs-6.14.00.root",
		"../../testdata/graphs.root",
		"../../testdata/small-flat-tree.root",
	} {
		t.Run(name, func(t *testing.T) {
			out := new(bytes.Buffer)
			cmd := rootls{stdout: out, streamers: true, trees: true}
			err := cmd.ls(name)
			if err != nil {
				t.Fatal(err)
			}
			ref := filepath.Join("testdata", filepath.Base(name)+".txt")
			want, err := ioutil.ReadFile(ref)
			if err != nil {
				t.Fatalf("could not open reference file: %v", err)
			}
			if got, want := out.String(), string(want); got != want {
				t.Fatalf("error:\ngot = %v\nwant= %v\n", got, want)
			}
		})
	}
}
