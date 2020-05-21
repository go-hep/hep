// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"golang.org/x/exp/shiny/screen"

	"github.com/google/shlex"
	"github.com/peterh/liner"
)

type Cmdr interface {
	Name() string
	Run(args []string) error
	Help(w io.Writer)
	Complete(line string) []string
}

type Cmd struct {
	msg  *log.Logger
	rl   *liner.State
	cmds map[string]Cmdr
	wmgr *winMgr
	fmgr *fileMgr
	hmgr *histMgr
}

func newCmd(stdout io.Writer, scr screen.Screen) *Cmd {
	msg := log.New(stdout, "paw: ", 0)
	c := Cmd{
		msg:  msg,
		rl:   liner.NewLiner(),
		wmgr: newWinMgr(scr, msg),
		fmgr: newFileMgr(msg),
		hmgr: newHistMgr(msg),
	}
	c.cmds = map[string]Cmdr{
		"/?": &cmdHelp{&c},
		"/!": &cmdShell{&c},

		"/file/open":   &cmdFileOpen{&c},
		"/file/close":  &cmdFileClose{&c},
		"/file/create": &cmdFileCreate{&c},
		"/file/ls":     &cmdFileList{&c},

		"/hist/open": &cmdHistOpen{&c},
		"/hist/plot": &cmdHistPlot{&c},

		"/quit": &cmdQuit{&c},
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
		_, _ = c.rl.ReadHistory(f)
	}

	return &c
}

func (c *Cmd) Close() error {
	var err error

	err = c.fmgr.Close()
	if err != nil {
		return err
	}

	f, err := os.Create(".pawgo.history")
	if err == nil {
		defer f.Close()
		_, _ = c.rl.WriteHistory(f)
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
				_, _ = c.msg.Writer().Write([]byte("\n"))
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
			if err == io.EOF {
				return err
			}
			c.msg.Printf("error: %v\n", err)
		}
		c.rl.AppendHistory(o)
	}
}

func (c *Cmd) RunScript(r io.Reader) error {
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		err := scan.Err()
		if err != nil {
			break
		}
		line := scan.Text()
		if line == "" || line[0] == '#' {
			continue
		}
		fmt.Fprintf(c.msg.Writer(), "# %s\n", line)
		err = c.exec(line)
		if err == io.EOF {
			return err
		}
		if err != nil {
			c.msg.Printf("error executing %q: %v\n", line, err)
			return err
		}
	}

	err := scan.Err()
	if err == io.EOF {
		err = nil
	}
	return err
}

func (c *Cmd) exec(line string) error {
	args, err := shlex.Split(line)
	if err != nil {
		return fmt.Errorf("paw: splitting line failed: %w", err)
	}
	cmd, ok := c.cmds[args[0]]
	if !ok {
		return fmt.Errorf("unknown command %q", args[0])
	}
	err = cmd.Run(args[1:])
	return err
}
