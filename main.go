package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/peterh/liner"
)

type Func func(args []string) error

type Cmd struct {
	rl   *liner.State
	cmds map[string]Func
	fmgr fileMgr
	hmgr histMgr
}

func newCmd() *Cmd {
	c := Cmd{
		rl:   liner.NewLiner(),
		fmgr: newFileMgr(),
		hmgr: newHistMgr(),
	}
	c.cmds = map[string]Func{
		"/file/open":   c.cmdOpenFile,
		"/file/close":  c.cmdCloseFile,
		"/file/create": c.cmdCreateFile,
		"/file/ls":     c.cmdListFile,

		"/hist/open": c.cmdOpenH1D,
		"/hist/plot": c.cmdPlotH1D,
	}

	c.rl.SetTabCompletionStyle(liner.TabPrints)
	c.rl.SetCompleter(func(line string) []string {
		var o []string
		for k := range c.cmds {
			if strings.HasPrefix(k, line) {
				o = append(o, k)
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

func (c *Cmd) cmdOpenFile(args []string) error {
	//fmt.Printf("/file/open: %v\n", args)
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	fname := args[1]
	err = c.fmgr.open(id, fname)
	return err
}

func (c *Cmd) cmdCreateFile(args []string) error {
	//fmt.Printf("/file/create: %v\n", args)
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	fname := args[1]
	err = c.fmgr.create(id, fname)
	return err
}

func (c *Cmd) cmdCloseFile(args []string) error {
	//fmt.Printf("/file/close: %v\n", args)
	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	err = c.fmgr.close(id)
	return err
}

func (c *Cmd) cmdListFile(args []string) error {
	//fmt.Printf("/file/ls: %v\n", args)
	if len(args) != 1 {
		return fmt.Errorf("/file/ls: need a file id")
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	err = c.fmgr.ls(id)
	return err
}

func (c *Cmd) cmdOpenH1D(args []string) error {
	var err error
	if len(args) != 2 {
		return fmt.Errorf("/hist/open: need histo-id and histo-name")
	}

	hid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	// e.g: /file/id/1/my-histo
	hname := args[1]

	err = c.hmgr.openH1D(&c.fmgr, hid, hname)
	return err
}

func (c *Cmd) cmdPlotH1D(args []string) error {
	var err error
	if len(args) < 1 {
		return fmt.Errorf("/hist/plot: need a histo-id to plot")
	}

	hid, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}

	err = c.hmgr.plotH1D(hid)
	return err
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
	fct, ok := c.cmds[args[0]]
	if !ok {
		return fmt.Errorf("unknown command %q", args[0])
	}
	err := fct(args[1:])
	return err
}

func main() {

	fname := flag.String("f", "", "paw script to execute")
	flag.Parse()

	os.Exit(xmain(*fname))
}

func xmain(fname string) int {
	icmd := newCmd()
	defer icmd.Close()

	switch fname {
	case "":
		err := icmd.Run()
		if err != nil {
			panic(err)
		}
	default:
		f, err := os.Open(fname)
		if err != nil {
			panic(err)
		}
		scan := bufio.NewScanner(f)
		for scan.Scan() {
			line := scan.Text()
			fmt.Printf("paw> %s\n", line)
			err := icmd.exec(line)
			if err != nil {
				fmt.Printf("**error** %v\n", err)
				return 1
			}
		}
	}

	return 0
}
