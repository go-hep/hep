// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hep

import (
	"bytes"
	"os/exec"
	"testing"
)

func TestGofmt(t *testing.T) {
	cmd := exec.Command("go", "tool", "golang.org/x/tools/cmd/goimports", "-d", ".")
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	cmd.Stderr = buf

	err := cmd.Run()
	if err != nil {
		t.Fatalf("error running goimports:\n%s\n%v", buf.String(), err)
	}

	if len(buf.Bytes()) != 0 {
		t.Fatalf("some files were not gofmt'ed:\n%s\n", buf.String())
	}
}
