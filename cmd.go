// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/peterh/liner"
)

type Func func(args []string) error

type Cmdr interface {
	Name() string
	Run(args []string) error
	Help(w io.Writer)
	Complete(line string) []string
}

type Cmd struct {
	rl   *liner.State
	cmds map[string]Cmdr
	fmgr fileMgr
	hmgr histMgr
}

func newCmd() *Cmd {
	c := Cmd{
		rl:   liner.NewLiner(),
		fmgr: newFileMgr(),
		hmgr: newHistMgr(),
	}
	c.cmds = map[string]Cmdr{
		"/help":        &cmdHelp{&c},
		"/file/open":   &cmdFileOpen{&c},
		"/file/close":  &cmdFileClose{&c},
		"/file/create": &cmdFileCreate{&c},
		"/file/ls":     &cmdFileList{&c},

		"/hist/open": &cmdH1DOpen{&c},
		"/hist/plot": &cmdH1DPlot{&c},
	}

	c.rl.SetTabCompletionStyle(liner.TabPrints)
	c.rl.SetCompleter(func(line string) []string {
		var o []string
		for k := range c.cmds {
			if strings.HasPrefix(k, line) {
				o = append(o, k+" ")
			}
		}
		if len(o) > 0 {
			return o
		}

		for k, cmd := range c.cmds {
			if strings.HasPrefix(line, k) {
				o = append(o, cmd.Complete(line)...)
			}
		}
		return o
	})

	f, err := os.Open(".pawgo.history")
	if err == nil {
		defer f.Close()
		c.rl.ReadHistory(f)
	}

	return &c
}

func (c *Cmd) Close() error {
	var err error

	err = c.fmgr.Close()

	f, err := os.Create(".pawgo.history")
	if err == nil {
		defer f.Close()
		c.rl.WriteHistory(f)
	}

	e := c.rl.Close()
	if e != nil {
		if err != nil {
			err = e
		}
	}

	return err
}

func (c *Cmd) Run() error {
	for {
		o, err := c.rl.Prompt("paw> ")
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}
		//fmt.Printf("<@ %q\n", o)
		if o == "" {
			continue
		}
		err = c.exec(o)
		if err != nil {
			fmt.Printf("**error** %v\n", err)
			if err == io.EOF {
				return err
			}
		}
		c.rl.AppendHistory(o)
	}
	panic("unreachable")
}

func (c *Cmd) exec(line string) error {
	args := strings.Split(line, " ")
	cmd, ok := c.cmds[args[0]]
	if !ok {
		return fmt.Errorf("unknown command %q", args[0])
	}
	err := cmd.Run(args[1:])
	return err
}
