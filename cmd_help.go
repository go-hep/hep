// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"text/tabwriter"
)

// cmdHelp prints the help
type cmdHelp struct {
	ctx *Cmd
}

func (cmd *cmdHelp) Name() string {
	return "/help"
}

func (cmd *cmdHelp) Run(args []string) error {
	var err error
	switch len(args) {
	case 1:
		var cmds []string
		for k := range cmd.ctx.cmds {
			cmds = append(cmds, k)
		}
		sort.Strings(cmds)
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		for _, k := range cmds {
			c := cmd.ctx.cmds[k]
			c.Help(w)
		}
		w.Flush()
	}
	return err
}

func (cmd *cmdHelp) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- print help\n", cmd.Name())
}

func (cmd *cmdHelp) Complete(line string) []string {
	var o []string
	return o
}
