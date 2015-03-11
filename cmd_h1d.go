// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
)

// cmdH1DOpen opens a histogram
type cmdH1DOpen struct {
	ctx *Cmd
}

func (cmd *cmdH1DOpen) Name() string {
	return "/hist/open"
}

func (cmd *cmdH1DOpen) Run(args []string) error {
	var err error
	if len(args) != 2 {
		return fmt.Errorf("%s: need histo-id and histo-name", cmd.Name())
	}

	hid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	// e.g: /file/id/1/my-histo
	hname := args[1]

	err = cmd.ctx.hmgr.openH1D(&cmd.ctx.fmgr, hid, hname)
	return err
}

func (cmd *cmdH1DOpen) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- open a histogram\n", cmd.Name())
}

func (cmd *cmdH1DOpen) Complete(line string) []string {
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
			sid := strconv.Itoa(id)
			switch {
			case strings.HasPrefix("/file/id/"+sid+"/", args[2]):
				r := cmd.ctx.fmgr.rfds[id]
				v := "/file/id/" + sid + "/"
				for _, k := range r.rio.Keys() {
					if strings.HasPrefix(v+k.Name, args[2]) {
						o = append(o, strings.Join(args[:2], " ")+" "+v+k.Name)
					}
				}
			case strings.HasPrefix("/file/id/"+sid, args[2]):
				o = append(o, strings.Join(args[:2], " ")+" /file/id/"+sid)
			}
		}
	}
	return o
}

// cmdH1DPlot opens a histogram
type cmdH1DPlot struct {
	ctx *Cmd
}

func (cmd *cmdH1DPlot) Name() string {
	return "/hist/plot"
}

func (cmd *cmdH1DPlot) Run(args []string) error {
	var err error
	if len(args) < 1 {
		return fmt.Errorf("%s: need a histo-id to plot", cmd.Name())
	}

	hid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	err = cmd.ctx.hmgr.plotH1D(hid)
	return err
}

func (cmd *cmdH1DPlot) Help(w io.Writer) {
	fmt.Fprintf(w, "%s \t-- plot a histogram\n", cmd.Name())
}

func (cmd *cmdH1DPlot) Complete(line string) []string {
	var o []string
	return o
}
