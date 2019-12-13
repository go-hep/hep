// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command npy2root converts the content of a NumPy data file to a ROOT file and tree.
//
// Usage: npy2root [OPTIONS] input.npy
//
// The NumPy data file format is described here:
//
//  https://docs.scipy.org/doc/numpy/neps/npy-format.html
//
// Example:
//
//  $> npyio-ls input.npy
//  ================================================================================
//  file: input.npy
//  npy-header: Header{Major:1, Minor:0, Descr:{Type:<f8, Fortran:false, Shape:[2 3]}}
//  data = [0 1 2 3 4 5]
//
//  $> npy2root -o output.root -t mytree ./input.npy
//  $> root-ls -t ./output.root
//  === [./output.root] ===
//  version: 61804
//    TTree   mytree       mytree  (entries=2)
//      numpy "numpy[3]/D" TBranch
//
//  $> root-dump ./output.root
//  >>> file[./output.root]
//  key[000]: mytree;1 "mytree" (TTree)
//  [000][numpy]: [0 1 2]
//  [001][numpy]: [3 4 5]
//
// Options:
//   -o string
//     	path to output ROOT file (default "output.root")
//   -t string
//     	name of the output ROOT tree (default "tree")
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/apache/arrow/go/arrow/arrio"
	"github.com/sbinet/npyio"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rarrow"
	"go-hep.org/x/hep/groot/rtree"
	"golang.org/x/xerrors"
)

func main() {
	log.SetPrefix("npy2root: ")
	log.SetFlags(0)

	oname := flag.String("o", "output.root", "path to output ROOT file")
	tname := flag.String("t", "tree", "name of the output ROOT tree")

	flag.Usage = func() {
		fmt.Printf(`npy2root converts the content of a NumPy data file to a ROOT file and tree.

Usage: npy2root [OPTIONS] input.npy

The NumPy data file format is described here:

 https://docs.scipy.org/doc/numpy/neps/npy-format.html

Example:

 $> npyio-ls input.npy
 ================================================================================
 file: input.npy
 npy-header: Header{Major:1, Minor:0, Descr:{Type:<f8, Fortran:false, Shape:[2 3]}}
 data = [0 1 2 3 4 5]

 $> npy2root -o output.root -t mytree ./input.npy
 $> root-ls -t ./output.root
 === [./output.root] ===
 version: 61804
   TTree   mytree       mytree  (entries=2)
     numpy "numpy[3]/D" TBranch

 $> root-dump ./output.root
 >>> file[./output.root]
 key[000]: mytree;1 "mytree" (TTree)
 [000][numpy]: [0 1 2]
 [001][numpy]: [3 4 5]

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing input NumPy data file")
	}

	fname := flag.Arg(0)
	err := process(*oname, *tname, fname)
	if err != nil {
		log.Fatalf("could not convert %q: %+v", fname, err)
	}
}

func process(oname, tname, fname string) error {
	src, err := os.Open(fname)
	if err != nil {
		return xerrors.Errorf("could not open numpy file %q: %w", fname, err)
	}
	defer src.Close()

	npy, err := npyio.NewReader(src)
	if err != nil {
		return xerrors.Errorf("could not create numpy file reader %q: %w", fname, err)
	}

	rec := NewRecord(npy)

	dst, err := groot.Create(oname)
	if err != nil {
		return xerrors.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer dst.Close()

	t, err := rarrow.NewFlatTreeWriter(dst, tname, rec.Schema(), rtree.WithTitle(tname))
	if err != nil {
		return xerrors.Errorf("could not create output ROOT tree %q: %w", tname, err)
	}

	_, err = arrio.Copy(t, NewRecordReader(rec))
	if err != nil {
		return xerrors.Errorf("could not copy numpy array to ROOT tree %q: %w", tname, err)
	}

	err = t.Close()
	if err != nil {
		return xerrors.Errorf("could not close output ROOT tree %q: %w", tname, err)
	}

	err = dst.Close()
	if err != nil {
		return xerrors.Errorf("could not close output ROOT file %q: %w", oname, err)
	}

	return nil
}
