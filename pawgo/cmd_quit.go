// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
)

// cmdQuit exits the application.
type cmdQuit struct {
	ctx *Cmd
}

func (cmd *cmdQuit) Name() string {
	return "/quit"
}

func (cmd *cmdQuit) Run(args []string) error {
	return io.EOF
}

func (cmd *cmdQuit) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- quit PAW-Go\n", cmd.Name())
}

func (cmd *cmdQuit) Complete(line string) []string {
	var o []string
	return o
}
