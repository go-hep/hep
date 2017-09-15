// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"strings"
)

// cmdHistOpen opens a histogram
type cmdHistOpen struct {
	ctx *Cmd
}

func (cmd *cmdHistOpen) Name() string {
	return "/hist/open"
}

func (cmd *cmdHistOpen) Run(args []string) error {
	var err error
	if len(args) < 2 {
		return fmt.Errorf("%s: need histo-id and histo-name (got=%v)", cmd.Name(), args)
	}

	hid := args[0]

	// e.g: /file/id/1/my-histo
	hname := args[1]

	err = cmd.ctx.hmgr.open(cmd.ctx.fmgr, hid, hname)
	return err
}

func (cmd *cmdHistOpen) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- open a histogram\n", cmd.Name())
}

func (cmd *cmdHistOpen) Complete(line string) []string {
	var o []string
	args := strings.Split(line, " ")
	switch len(args) {
	case 0, 1:
		return o
	case 2:
		return o
	case 3:
		if args[2] == "" {
			args[2] = "/file/id/"
		}
		for id := range cmd.ctx.fmgr.rfds {
			switch {
			case strings.HasPrefix("/file/id/"+id+"/", args[2]):
				r := cmd.ctx.fmgr.rfds[id]
				v := "/file/id/" + id + "/"
				for _, k := range r.rio.Keys() {
					if strings.HasPrefix(v+k.Name, args[2]) {
						o = append(o, strings.Join(args[:2], " ")+" "+v+k.Name)
					}
				}
			case strings.HasPrefix("/file/id/"+id, args[2]):
				o = append(o, strings.Join(args[:2], " ")+" /file/id/"+id)
			}
		}
	}
	return o
}

// cmdHistPlot plots a histogram
type cmdHistPlot struct {
	ctx *Cmd
}

func (cmd *cmdHistPlot) Name() string {
	return "/hist/plot"
}

func (cmd *cmdHistPlot) Run(args []string) error {
	var err error
	if len(args) < 1 {
		return fmt.Errorf("%s: need a histo-id to plot", cmd.Name())
	}

	hid := args[0]
	err = cmd.ctx.hmgr.plot(cmd.ctx.fmgr, cmd.ctx.wmgr, hid)
	return err
}

func (cmd *cmdHistPlot) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- plot a histogram\n", cmd.Name())
}

func (cmd *cmdHistPlot) Complete(line string) []string {
	var o []string
	args := strings.Split(line, " ")
	switch len(args) {
	case 0, 1:
		return o
	case 2:
		if strings.HasPrefix(args[1], "/") {
			for id, r := range cmd.ctx.fmgr.rfds {
				for _, k := range r.rio.Keys() {
					name := "/file/id/" + id + "/" + k.Name
					if strings.HasPrefix(name, args[1]) {
						o = append(o, args[0]+" "+name)
					}
				}
			}
			return o
		}
		for k := range cmd.ctx.hmgr.hmap {
			if strings.HasPrefix(k, args[1]) {
				o = append(o, args[0]+" "+k)
			}
		}
		return o
	}
	return o
}
