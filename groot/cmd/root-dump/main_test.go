// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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
				t.Fatalf("error:\n%s\n", diff(t, got, want))
			}
		})
	}
}

func diff(t *testing.T, chk, ref string) string {
	t.Helper()

	if !hasDiffCmd {
		return fmt.Sprintf("=== got ===\n%s\n=== want ===\n%s\n", chk, ref)
	}

	tmpdir, err := os.MkdirTemp("", "groot-diff-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %+v", err)
	}
	defer os.RemoveAll(tmpdir)

	got := filepath.Join(tmpdir, "got.txt")
	err = os.WriteFile(got, []byte(chk), 0644)
	if err != nil {
		t.Fatalf("could not create %s file: %+v", got, err)
	}

	want := filepath.Join(tmpdir, "want.txt")
	err = os.WriteFile(want, []byte(ref), 0644)
	if err != nil {
		t.Fatalf("could not create %s file: %+v", want, err)
	}

	out := new(bytes.Buffer)
	cmd := exec.Command("diff", "-urN", want, got)
	cmd.Stdout = out
	cmd.Stderr = out
	err = cmd.Run()
	return out.String() + "\nerror: " + err.Error()
}

var hasDiffCmd = false

func init() {
	_, err := exec.LookPath("diff")
	if err == nil {
		hasDiffCmd = true
	}
}
