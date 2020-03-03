// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// arrow2root converts the content of an ARROW file to a ROOT TTree.
package main // import "go-hep.org/x/hep/cmd/arrow2root"

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/apache/arrow/go/arrow/arrio"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rarrow"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	log.SetPrefix("arrow2root: ")
	log.SetFlags(0)

	oname := flag.String("o", "output.root", "path to output ROOT file name")
	tname := flag.String("t", "tree", "name of the output tree")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing input ARROW filename argument")
	}
	fname := flag.Arg(0)

	err := process(*oname, *tname, fname)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func process(oname, tname, fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open ARROW file %q: %w", fname, err)
	}
	defer f.Close()

	mem := memory.NewGoAllocator()
	r, err := ipc.NewFileReader(f, ipc.WithAllocator(mem))
	if err != nil {
		return fmt.Errorf("could not create ARROW IPC reader from %q: %w", fname, err)
	}
	defer r.Close()

	o, err := groot.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer o.Close()

	tree, err := rarrow.NewFlatTreeWriter(o, tname, r.Schema(), rtree.WithTitle(tname))
	if err != nil {
		return fmt.Errorf("could not create output ROOT tree %q: %w", tname, err)
	}

	_, err = arrio.Copy(tree, r)
	if err != nil {
		return fmt.Errorf("could not convert ARROW file to ROOT tree: %w", err)
	}

	err = tree.Close()
	if err != nil {
		return fmt.Errorf("could not close ROOT tree writer: %w", err)
	}

	err = o.Close()
	if err != nil {
		return fmt.Errorf("could not close output ROOT file %q: %w", oname, err)
	}

	return nil
}
