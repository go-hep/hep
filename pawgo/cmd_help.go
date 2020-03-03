// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
)

// cmdHelp prints the help
type cmdHelp struct {
	ctx *Cmd
}

func (cmd *cmdHelp) Name() string {
	return "/?"
}

func (cmd *cmdHelp) Run(args []string) error {
	var err error
	switch len(args) {
	case 0:
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
	case 1:
		c, ok := cmd.ctx.cmds[args[0]]
		if !ok {
			return fmt.Errorf("unknown command %q", args[0])
		}
		w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
		c.Help(w)
		w.Flush()
	}
	return err
}

func (cmd *cmdHelp) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- print help\n", cmd.Name())
}

func (cmd *cmdHelp) Complete(line string) []string {
	var o []string
	args := strings.Split(line, " ")
	switch len(args) {
	case 0, 1:
		return o
	case 2:
		if args[1] == "" {
			args[1] = "/"
		}
		for k := range cmd.ctx.cmds {
			if strings.HasPrefix(k, args[1]) {
				o = append(o, strings.Join(args[:1], " ")+" "+k)
			}
		}
	}
	return o
}
