// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"os"
	"testing"
)

func TestInspect(t *testing.T) {
	for _, fname := range []string{
		"../../testdata/run-header_golden.slcio",
		"../../testdata/event_golden.slcio",
	} {

		t.Run(fname, func(t *testing.T) {
			buf := new(bytes.Buffer)
			inspect(buf, fname, -1, false)

			got := buf.Bytes()
			want, err := os.ReadFile(fname + ".txt")
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want) {
				t.Fatalf("error.\ngot = %q\nwant= %q\n", string(got), string(want))
			}
		})
	}
}
