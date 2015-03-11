// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

func main() {
	rc := xmain()
	os.Exit(rc)
}

func xmain() int {

	fname := flag.String("f", "", "paw script to execute")
	flag.Parse()

	fmt.Printf(`
:::::::::::::::::::::::::::::
:::   Welcome to PAW-Go   :::
:::::::::::::::::::::::::::::

Type /? for help.
^D to quit.

`)

	icmd := newCmd()
	defer icmd.Close()

	switch *fname {
	case "":
		err := icmd.Run()
		if err != nil {
			panic(err)
		}
	default:
		f, err := os.Open(*fname)
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
