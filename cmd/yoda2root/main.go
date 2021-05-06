// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// yoda2root converts YODA files containing hbook-like values (H1D, H2D, P1D, ...)
// into ROOT files.
//
// Example:
//
//  $> yoda2root rivet.yoda rivet.root
//  $> yoda2root rivet.yoda.gz rivet.root
package main // import "go-hep.org/x/hep/cmd/yoda2root"

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/hbook/yodacnv"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("yoda2root: ")
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: yoda2root [options] <file1.yoda> [<file2.yoda> [...]] file.root

ex:
 $> yoda2root rivet.yoda rivet.root
 $> yoda2root rivet.yoda.gz rivet.root
`)
	}

	flag.Parse()

	if flag.NArg() < 2 {
		log.Printf("missing input and/or output file name(s)")
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}

	oname := flag.Arg(flag.NArg() - 1)
	f, err := groot.Create(oname)
	if err != nil {
		log.Fatalf("could not create output ROOT file: %v", err)
	}
	defer f.Close()

	for _, fname := range flag.Args()[:flag.NArg()-1] {
		err = convert(f, fname)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("could not close output ROOT file: %v", err)
	}
}

func convert(w *riofs.File, fname string) error {
	var r io.ReadCloser
	r, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("error opening file [%s]: %w", fname, err)
	}
	defer r.Close()

	if filepath.Ext(fname) == ".gz" {
		rz, err := gzip.NewReader(r)
		if err != nil {
			return fmt.Errorf("could not open gzip file [%s]: %w", fname, err)
		}
		defer rz.Close()
		r = rz
	}

	vs, err := yodacnv.Read(r)
	if err != nil {
		return fmt.Errorf("error decoding YODA file [%s]: %w", fname, err)
	}

	for i, v := range vs {
		var (
			obj root.Object
			key string
		)
		switch v := v.(type) {
		case *hbook.H1D:
			key = "h1"
			obj = rootcnv.FromH1D(v)
		case *hbook.H2D:
			key = "h2"
			obj = rootcnv.FromH2D(v)
		case *hbook.S2D:
			key = "scatter"
			obj = rootcnv.FromS2D(v)
		default:
			log.Printf("%s: no YODA -> ROOT conversion for %T", fname, v)
			continue
		}
		switch v.Name() {
		case "":
			key = fmt.Sprintf("yoda-%s-%03d", key, i)
		default:
			key = v.Name()
		}
		err = riofs.Dir(w).Put(key, obj)
		if err != nil {
			return fmt.Errorf("error writing %q from YODA file [%s]: %w", v.Name(), fname, err)
		}
	}

	return nil
}
