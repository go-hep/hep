// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"go-hep.org/x/hep/internal/diff"
)

func TestGenerator(t *testing.T) {
	for _, tc := range []struct {
		name  string
		rules []string
		want  string
	}{
		{
			name: "testdata/datalayout.yaml",
			rules: []string{
				"ex2::->ex2_",
				"ex42::->ex42_",
			},
			want: "testdata/datalayout.go",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := new(bytes.Buffer)
			err := process(got, "podio", tc.name, strings.Join(tc.rules, ","))
			if err != nil {
				t.Fatalf("could not process %q: %+v", tc.name, err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference file %q: %+v", tc.want, err)
			}

			if got, want := got.Bytes(), want; !bytes.Equal(got, want) {
				t.Fatalf("invalid generated code:\n%s", diff.Format(string(got), string(want)))
			}
		})
	}
}
