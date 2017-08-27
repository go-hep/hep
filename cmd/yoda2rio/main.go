// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// yoda2rio converts YODA files containing hbook-like values (H1D, H2D, P1D, ...)
// into rio files.
//
// Example:
//
//  $> yoda2rio rivet.yoda >| rivet.rio
//  $> yoda2rio rivet.yoda.gz >| rivet.rio
package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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
 $> yoda2rio rivet.yoda    >| rivet.rio
 $> yoda2rio rivet.yoda.gz >| rivet.rio
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
	var r io.ReadCloser
	r, err := os.Open(fname)
	if err != nil {
		log.Fatalf("error opening file [%s]: %v\n", fname, err)
	}
	defer r.Close()

	if filepath.Ext(fname) == ".gz" {
		rz, err := gzip.NewReader(r)
		if err != nil {
			log.Fatal(err)
		}
		defer rz.Close()
		r = rz
	}

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
