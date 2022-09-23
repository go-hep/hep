// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go-hep.org/x/hep/groot/rcmd"
)

func TestConvert(t *testing.T) {
	tmp, err := os.MkdirTemp("", "hepmc2root-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmp)

	for _, name := range []string{
		"testdata/small.hepmc",
	} {
		t.Run(name, func(t *testing.T) {
			oname := filepath.Join(tmp, filepath.Base(name)+".root")
			tname := "tree"

			err := process(oname, tname, name)
			if err != nil {
				t.Fatalf("could not convert %q: %+v", name, err)
			}

			var (
				out  = new(strings.Builder)
				deep = true
			)
			err = rcmd.Dump(out, oname, deep, nil)
			if err != nil {
				t.Fatalf("could not dump ROOT file %q: %+v", oname, err)
			}

			want, err := os.ReadFile(name + ".txt")
			if err != nil {
				t.Fatalf("could not load reference file %q: %+v", name, err)
			}

			if got, want := out.String(), string(want); got != want {
				t.Fatalf("invalid root-dump output:\ngot:\n%s\nwant:\n%s\n", got, want)
			}
		})
	}
}
