// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// pawgo is a simple interactive shell to quickly plot hbook histograms from
// rio files.
//
// Example:
//
//  $> pawgo
//  paw> /file/open f testdata/issue-120.rio
//  paw> /file/ls f
//  /file/id/f name=testdata/issue-120.rio
//   	- MonoH_Truth/jets	(type="*go-hep.org/x/hep/hbook.H1D")
//
//  paw> /hist/open h /file/id/f/MonoH_Truth/jets
//  paw> /hist/plot h
//  == h1d: name="MonoH_Truth/jets"
//  entries=20000
//  mean=  +2.554
//  RMS=   +2.891
//  paw> /?
//  /! 		-- run a shell command
//  /? 		-- print help
//  /file/close 	-- close a file
//  /file/create 	-- create file for write access
//  /file/list 	-- list a file's content
//  /file/open 	-- open file for read access
//  /hist/open 	-- open a histogram
//  /hist/plot 	-- plot a histogram
//  /quit 		-- quit PAW-Go
package main // import "go-hep.org/x/hep/pawgo"

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
	defer fmt.Printf("bye.\n")

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
