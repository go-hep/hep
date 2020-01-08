// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/cmd/arrow2root"

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/rcmd"
)

func TestConvert(t *testing.T) {
	tmp, err := ioutil.TempDir("", "arrow2root-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, tc := range []struct {
		name   string
		panics string
	}{
		{
			name: "testdata/primitives.file.data",
		},
		{
			name: "testdata/arrays.file.data",
		},
		{
			name: "testdata/strings.file.data",
		},
		{
			name: "testdata/fixed_size_binaries.file.data",
		},
		{
			name: "testdata/lists.file.data",
		},
		{
			name:   "testdata/structs.file.data",
			panics: "invalid ARROW data-type: *arrow.StructType", // FIXME(sbinet): needs non-flat-tree writer support
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic (%s)", tc.panics)
					}
					if got, want := err.(error).Error(), tc.panics; got != want {
						t.Fatalf("invalid panic message:\ngot= %v\nwant=%v", got, want)
					}
				}()
			}
			oname := filepath.Join(tmp, filepath.Base(tc.name)+".root")
			tname := "tree"
			err := process(oname, tname, tc.name)
			if err != nil {
				t.Fatalf("could not convert %q: %+v", tc.name, err)
			}

			var (
				out  = new(strings.Builder)
				deep = true
			)
			err = rcmd.Dump(out, oname, deep, nil)
			if err != nil {
				t.Fatalf("could not dump ROOT file %q: %+v", oname, err)
			}

			want, err := ioutil.ReadFile(tc.name + ".txt")
			if err != nil {
				t.Fatalf("could not load reference file %q: %+v", tc.name, err)
			}

			if got, want := out.String(), string(want); got != want {
				t.Fatalf("invalid root-dump output:\ngot:\n%s\nwant:\n%s\n", got, want)
			}
		})
	}
}
