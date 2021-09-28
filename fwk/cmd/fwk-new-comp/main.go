// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	g_type = flag.String("c", "task", "type of component to generate (task|svc)")
	g_pkg  = flag.String("p", "", "name of the package holding the component")
)

type Component struct {
	Package string
	Name    string
	Type    string
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("fwk-new-comp: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %[1]s [options] <component-name>

ex:
 $ %[1]s -c=task -p=mypackage mytask
 $ %[1]s -c=task -p mypackage mytask >| mytask.go
 $ %[1]s -c=svc  -p mypackage mysvc  >| mysvc.go

options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}

	flag.Parse()
	sc := run()
	os.Exit(sc)
}

func run() int {
	if *g_type != "svc" && *g_type != "task" {
		log.Printf("**error** invalid component type [%s]\n", *g_type)
		flag.Usage()
		return 1
	}

	if *g_pkg == "" {
		// take directory name
		wd, err := os.Getwd()
		if err != nil {
			log.Printf("**error** could not get directory name: %v\n", err)
			return 1
		}
		*g_pkg = filepath.Base(wd)

		if *g_pkg == "" || *g_pkg == "." {
			log.Printf(
				"**error** invalid package name %q. please specify via the '-p' flag.",
				*g_pkg,
			)
			return 1
		}
	}

	args := flag.Args()
	if len(args) <= 0 {
		log.Printf("**error** you need to give a component name\n")
		flag.Usage()
		return 1
	}

	c := Component{
		Package: *g_pkg,
		Name:    args[0],
		Type:    *g_type,
	}

	var err error
	switch *g_type {
	case "svc":
		err = gen_svc(c)
	case "task":
		err = gen_task(c)
	default:
		log.Printf("**error** invalid component type [%s]\n", *g_type)
		flag.Usage()
		return 1
	}

	if err != nil {
		log.Printf("**error** generating %q: %v\n", c.Name, err)
		return 1
	}

	return 0
}
