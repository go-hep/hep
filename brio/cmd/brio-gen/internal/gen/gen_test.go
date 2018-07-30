// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen_test

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"testing"

	"go-hep.org/x/hep/brio/cmd/brio-gen/internal/gen"
)

func TestGenerator(t *testing.T) {
	const pkg = "go-hep.org/x/hep/brio/cmd/brio-gen/internal/gen/_test/pkg"
	err := exec.Command("go", "get", pkg).Run()
	if err != nil {
		t.Fatalf("could not build test package: %v", err)
	}

	g, err := gen.NewGenerator(pkg)
	if err != nil {
		t.Fatal(err)
	}

	g.Generate("T1")
	g.Generate("T2")
	g.Generate("T3")

	got, err := g.Format()
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/brio_gen_golden.go.txt")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("files differ.\ngot = %q\nwant= %q\n", string(got), string(want))
	}
}
