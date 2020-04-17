// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func run(bin string, args ...string) error {
	buf := new(bytes.Buffer)
	cmd := exec.Command(bin, args...)
	cmd.Stdout = buf
	cmd.Stderr = buf
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf(
			"error running %q:\n%s\nerr: %w",
			strings.Join(cmd.Args, " "),
			string(buf.Bytes()),
			err,
		)
	}

	return nil
}

func TestIssue120(t *testing.T) {
	cmd := "pawgo"
	if runtime.GOOS == "windows" {
		cmd = "pawgo.exe"
	}
	err := run(cmd, "./testdata/issue-120.paw")
	if err != nil {
		t.Fatalf("%+v", err)
	}
}
