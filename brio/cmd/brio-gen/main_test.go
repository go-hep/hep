// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command brio-gen generates (un)marshaler code for types.
package main

import (
	"bytes"
	"os"
	"os/exec"
	"testing"

	// make sure this has been compiled for TestGenerate
	_ "go-hep.org/x/hep/brio/cmd/brio-gen/internal/briotest"
)

func TestGenerate(t *testing.T) {
	diff, err := exec.LookPath("diff")
	hasDiff := err == nil

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
				err = os.WriteFile(outfile, got, 0644)
				if err == nil && hasDiff {
					out := new(bytes.Buffer)
					cmd := exec.Command(diff, "-urN", outfile, tc.want)
					cmd.Stdout = out
					cmd.Stderr = out
					err = cmd.Run()
					t.Fatalf("generated code error: %v\n%v\n", err, out.String())
				}
				t.Fatalf("generated code error.\ngot:\n%s\nwant:\n%s\n", string(got), string(want))
			}
			// Remove output if test passes
			// Note: output file are referenced in .gitignore
			_ = os.Remove(outfile)
		})
	}
}
