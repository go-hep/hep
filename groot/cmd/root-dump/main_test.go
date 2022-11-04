// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"

	"go-hep.org/x/hep/internal/diff"
)

func TestROOTDump(t *testing.T) {
	const deep = true
	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "../../testdata/simple.root",
			want: "testdata/simple.txt",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			o := new(bytes.Buffer)
			err := dump(o, tc.name, deep)
			if err != nil {
				t.Fatalf("could not dump %q: %+v", tc.name, err)
			}

			want, err := os.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference file: %+v", err)
			}

			if got, want := o.String(), string(want); got != want {
				t.Fatalf("error:\n%s\n", diff.Format(got, want))
			}
		})
	}
}
