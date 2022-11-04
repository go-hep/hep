// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command brio-gen generates (un)marshaler code for types.
package main

import (
	"bytes"
	"os"
	"testing"

	// make sure this has been compiled for TestGenerate
	_ "go-hep.org/x/hep/brio/cmd/brio-gen/internal/briotest"
	"go-hep.org/x/hep/internal/diff"
)

func TestGenerate(t *testing.T) {
	for _, tc := range []struct {
		name  string
		types []string
		want  string
	}{
		{
			name:  "image",
			types: []string{"Point"},
			want:  "testdata/image_brio.go",
		},
		{
			name:  "go-hep.org/x/hep/brio/cmd/brio-gen/internal/briotest",
			types: []string{"Hist", "Bin"},
			want:  "testdata/briotest_brio.go",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := generate(buf, tc.name, tc.types)
			if err != nil {
				t.Fatal(err)
			}

			got := buf.Bytes()
			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatal(err)
			}
			outfile := tc.want + "_got"
			if !bytes.Equal(got, want) {
				t.Fatalf("generated code error:\n%s\n", diff.Format(string(got), string(want)))
			}
			// Remove output if test passes
			// Note: output file are referenced in .gitignore
			_ = os.Remove(outfile)
		})
	}
}
