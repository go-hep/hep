// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rio-ls-records displays the list of records stored in a given rio file.
package main

import (
	"flag"
	"fmt"
	"os"

	"go-hep.org/x/hep/rio"
)

func main() {
	var fname string

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	fmt.Printf("::: inspecting file [%s]...\n", fname)
	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	r, err := rio.NewReader(f)
	if err != nil {
		fmt.Printf("*** error creating rio.Reader: %v\n", err)
		os.Exit(2)
	}

	scan := rio.NewScanner(r)
	for scan.Scan() {
		// scans through the whole stream
		err = scan.Err()
		if err != nil {
			break
		}
		rec := scan.Record()
		fmt.Printf(" -> %v\n", rec.Name())
	}
	err = scan.Err()
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("::: inspecting file [%s]... [done]\n", fname)
}
