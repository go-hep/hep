// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
)

func main() {
	rc := 0
	driver.Main(func(scr screen.Screen) {
		rc = xmain(scr)
	})
	os.Exit(rc)
}

func xmain(scr screen.Screen) int {

	interactive := flag.Bool(
		"i", false,
		"enable interactive mode: drop into PAW-Go prompt after processing script files",
	)

	flag.Parse()

	fmt.Printf(`
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D or /quit to quit.

`)

	icmd := newCmd(scr)
	defer icmd.Close()

	if flag.NArg() > 0 {
		for _, fname := range flag.Args() {
			f, err := os.Open(fname)
			if err != nil {
				icmd.msg.Printf("error: %v\n", err)
				return 1
			}
			defer f.Close()

			err = icmd.RunScript(f)
			if err == io.EOF {
				return 0
			}
			if err != nil {
				icmd.msg.Printf("error running script [%s]: %v\n", f.Name(), err)
				return 1
			}
		}
		if !*interactive {
			return 0
		}
	}

	err := icmd.Run()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		icmd.msg.Printf("error running interpreter: %v\n", err)
		return 1
	}

	return 0
}
