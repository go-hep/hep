// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"go-hep.org/x/hep/groot/internal/rdatatest" // make sure this is compiled
)

var evt rdatatest.Event

func TestGenerate(t *testing.T) {

	for _, tc := range []struct {
		pkg   string
		types []string
		want  string
	}{
		{
			pkg:   "go-hep.org/x/hep/groot/internal/rdatatest",
			types: []string{"Event", "HLV", "Particle"},
			want:  "testdata/rdatatest.txt",
		},
	} {
		t.Run(tc.pkg, func(t *testing.T) {
			buf := new(bytes.Buffer)
			err := generate(buf, tc.pkg, tc.types)
			if err != nil {
				t.Fatalf("could not generate streamer: %v", err)
			}
			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference streamer: %v", err)
			}

			if got, want := buf.String(), string(want); got != want {
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

	tmpdir, err := ioutil.TempDir("", "groot-diff-")
	if err != nil {
		t.Fatalf("could not create tmpdir: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	got := filepath.Join(tmpdir, "got.txt")
	err = ioutil.WriteFile(got, []byte(chk), 0644)
	if err != nil {
		t.Fatalf("could not create %s file: %v", got, err)
	}

	want := filepath.Join(tmpdir, "want.txt")
	err = ioutil.WriteFile(want, []byte(ref), 0644)
	if err != nil {
		t.Fatalf("could not create %s file: %v", want, err)
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
