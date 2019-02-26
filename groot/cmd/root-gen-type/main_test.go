// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	regen = flag.Bool("regen", false, "regenerate reference files")
)

func TestGenerate(t *testing.T) {
	dir, err := ioutil.TempDir("", "groot-gen-type-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for _, tc := range []struct {
		fname     string
		want      string
		types     []string
		verbose   bool
		streamers bool
	}{
		{
			fname:     "../../testdata/small-evnt-tree-fullsplit.root",
			want:      "testdata/small-evnt-tree-fullsplit.txt",
			types:     []string{"Event", "P3"},
			streamers: true,
		},
	} {
		t.Run(tc.fname, func(t *testing.T) {
			oname := filepath.Base(tc.fname) + ".go"
			o, err := os.Create(filepath.Join(dir, oname))
			if err != nil {
				t.Fatal(err)
			}
			defer o.Close()

			err = generate(o, "main", tc.types, tc.fname, tc.verbose, tc.streamers)
			if err != nil {
				t.Fatalf("could not generate types: %v", err)
			}

			err = o.Close()
			if err != nil {
				t.Fatal(err)
			}

			got, err := ioutil.ReadFile(o.Name())
			if err != nil {
				t.Fatalf("could not read generated file: %v", err)
			}

			if *regen {
				ioutil.WriteFile(tc.want, got, 0644)
			}

			want, err := ioutil.ReadFile(tc.want)
			if err != nil {
				t.Fatalf("could not read reference file: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("error:\n%v", diff(t, string(got), string(want)))
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
