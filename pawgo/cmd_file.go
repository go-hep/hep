// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/google/shlex"
)

// cmdFileOpen opens a file for read access
type cmdFileOpen struct {
	ctx *Cmd
}

func (cmd *cmdFileOpen) Name() string {
	return "/file/open"
}

func (cmd *cmdFileOpen) Run(args []string) error {
	var err error
	id := args[0]
	fname := args[1]
	err = cmd.ctx.fmgr.open(id, fname)
	return err
}

func (cmd *cmdFileOpen) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- open file for read access\n", cmd.Name())
}

func (cmd *cmdFileOpen) Complete(line string) []string {
	var o []string
	args, err := shlex.Split(line)
	if err != nil {
		cmd.ctx.msg.Printf("error splitting line: %v\n", err)
		return o
	}
	switch len(args) {
	case 0:
		return o
	case 1:
		// fmt.Printf(">>> %q\n", args[0])
	case 2:
		// fmt.Printf("### %q %q\n", args[0], args[1])
	case 3:
		// fmt.Printf("+++ %q %q %q\n", args[0], args[1], args[2])
		if args[2] != "" {
			matches, err := filepath.Glob(args[2] + "*")
			//fmt.Printf(">>> matches: %v\nerr=%v\n", matches, err)
			if err != nil {
				return o
			}
			for _, m := range matches {
				mm := strings.Trim(m, "\t\n\r ")
				if mm != "" {
					args[2] = mm
					o = append(o, strings.Join(args, " "))
				}
			}
		}
	}

	return o
}

// cmdFileCreate creates a file for write access
type cmdFileCreate struct {
	ctx *Cmd
}

func (cmd *cmdFileCreate) Name() string {
	return "/file/create"
}

func (cmd *cmdFileCreate) Run(args []string) error {
	var err error
	id := args[0]
	fname := args[1]
	err = cmd.ctx.fmgr.create(id, fname)
	return err
}

func (cmd *cmdFileCreate) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- create file for write access\n", cmd.Name())
}

func (cmd *cmdFileCreate) Complete(line string) []string {
	var o []string
	return o
}

// cmdFileClose closes a file
type cmdFileClose struct {
	ctx *Cmd
}

func (cmd *cmdFileClose) Name() string {
	return "/file/close"
}

func (cmd *cmdFileClose) Run(args []string) error {
	var err error
	id := args[0]
	err = cmd.ctx.fmgr.close(id)
	return err
}

func (cmd *cmdFileClose) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- close a file\n", cmd.Name())
}

func (cmd *cmdFileClose) Complete(line string) []string {
	var o []string
	return o
}

// cmdFileList closes a file
type cmdFileList struct {
	ctx *Cmd
}

func (cmd *cmdFileList) Name() string {
	return "/file/list"
}

func (cmd *cmdFileList) Run(args []string) error {
	switch len(args) {
	case 0:
		return cmd.ctx.fmgr.ls("")
	case 1:
		return cmd.ctx.fmgr.ls(args[0])
	default:
		return fmt.Errorf("%s: need at most 1 file id", cmd.Name())
	}
}

func (cmd *cmdFileList) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- list a file's content\n", cmd.Name())
}

func (cmd *cmdFileList) Complete(line string) []string {
	var o []string
	return o
}
