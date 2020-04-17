// Copyright ©2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
)

// cmdShell runs a shell command.
type cmdShell struct {
	ctx *Cmd
}

func (cmd *cmdShell) Name() string {
	return "/!"
}

func (cmd *cmdShell) Run(args []string) error {
	sh := exec.Command(args[0], args[1:]...)
	sh.Stdin = os.Stdin
	sh.Stdout = os.Stdout
	sh.Stderr = os.Stderr
	return sh.Run()
}

func (cmd *cmdShell) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- run a shell command\n", cmd.Name())
}

func (cmd *cmdShell) Complete(line string) []string {
	var o []string
	args, err := shlex.Split(line)
	if err != nil {
		cmd.ctx.msg.Printf("error splitting line: %v\n", err)
		return o
	}
	if len(args) < 2 {
		return o
	}
	i := len(args) - 1
	if args[i] != "" {
		matches, err := filepath.Glob(args[i] + "*")
		//fmt.Printf(">>> matches: %v\nerr=%v\n", matches, err)
		if err != nil {
			return o
		}
		for _, m := range matches {
			mm := strings.Trim(m, "\t\n\r ")
			if mm != "" {
				args[i] = mm
				o = append(o, strings.Join(args, " "))
			}
		}
	}

	return o
}
