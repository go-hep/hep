// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// yoda2rio converts YODA files containing hbook-like values (H1D, H2D, P1D, ...)
// into rio files.
//
// Example:
//
//  $> yoda2rio rivet.yoda >| rivet.rio
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/hbook/yodacnv"
	"go-hep.org/x/hep/rio"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("yoda2rio: ")
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: yoda2rio [options] <file1.yoda> [<file2.yoda> [...]]

ex:
 $ yoda2rio rivet.yoda >| rivet.rio
`)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		log.Printf("missing input file name")
		flag.Usage()
		flag.PrintDefaults()
	}

	o, err := rio.NewWriter(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()

	for _, fname := range flag.Args() {
		convert(o, fname)
	}
}

func convert(w *rio.Writer, fname string) {
	r, err := os.Open(fname)
	if err != nil {
		log.Fatalf("error opening file [%s]: %v\n", fname, err)
	}
	defer r.Close()

	vs, err := yodacnv.Read(r)
	if err != nil {
		log.Fatalf("error decoding YODA file [%s]: %v\n", fname, err)
	}

	for _, v := range vs {
		err = w.WriteValue(v.Name(), v)
		if err != nil {
			log.Fatalf("error writing %q from YODA file [%s]: %v\n", v.Name(), fname, err)
		}
	}
}
