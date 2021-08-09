// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestROOTls(t *testing.T) {
	tmp, err := os.MkdirTemp("", "root-ls-")
	if err != nil {
		t.Fatalf("could not create tmp dir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, tc := range []struct {
		name string
		rc   int
	}{
		{
			name: "../../testdata/dirs-6.14.00.root",
		},
		{
			name: "../../testdata/graphs.root",
		},
		{
			name: "../../testdata/small-flat-tree.root",
		},
		{
			name: filepath.Join(tmp, "not-there.root"),
			rc:   1,
		},
		{
			name: "-h",
			rc:   0,
		},
		{
			name: "-=3",
			rc:   1,
		},
		{
			name: "-cpu-profile=" + filepath.Join(tmp, "cpu.prof"),
			rc:   1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out := new(bytes.Buffer)
			rc := run(out, out, []string{
				"-sinfos", "-t", tc.name,
			})

			if rc != tc.rc {
				t.Fatalf(
					"invalid exit-code for root-ls: got=%d, want=%d\n%s",
					rc, tc.rc, out.String(),
				)
			}
			if rc != 0 || tc.name == "-h" {
				return
			}

			ref := filepath.Join("testdata", filepath.Base(tc.name)+".txt")
			want, err := os.ReadFile(ref)
			if err != nil {
				t.Fatalf("could not open reference file: %v", err)
			}
			if got, want := out.String(), string(want); got != want {
				t.Fatalf("error:\ngot = %v\nwant= %v\n", got, want)
			}
		})
	}
}
