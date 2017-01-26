// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root-ls dumps the content of a ROOT file
package main // import "github.com/go-hep/rootio/cmd/root-ls"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/go-hep/rootio"
)

var g_prof = flag.String("profile", "", "filename of cpuprofile")

func main() {
	flag.Parse()

	if *g_prof != "" {
		f, err := os.Create(*g_prof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if flag.NArg() <= 0 {
		fmt.Fprintf(os.Stderr, "**error** you need to give a ROOT file\n")
		flag.Usage()
		os.Exit(1)
	}

	for ii, fname := range flag.Args() {

		if ii > 0 {
			fmt.Printf("\n")
		}

		fmt.Printf("=== [%s] ===\n", fname)
		f, err := rootio.Open(fname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error** %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("version: %v\n", f.Version())

		for _, k := range f.Keys() {
			fmt.Printf("%-8s %-40s %s\n", k.Class(), k.Name(), k.Title())
		}
	}
}
