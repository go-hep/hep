// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command hepmc2root converts a HepMC2 ASCII file into a ROOT file and (flat) tree.
//
// Usage: hepmc2root [OPTIONS] hepmc.ascii
//
// Example:
//
// $> hepmc2root ./hepmc.ascii
// $> hepmc2root -o out.root -t mytree ./hepmc.ascii
//
// Options:
//
//	-o string
//	  	path to output ROOT file name (default "out.root")
//	-t string
//	  	name of the output tree (default "tree")
package main // import "go-hep.org/x/hep/cmd/hepmc2root"

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hepmc"
	"go-hep.org/x/hep/hepmc/rootcnv"
)

func main() {
	log.SetPrefix("hepmc2root: ")
	log.SetFlags(0)

	oname := flag.String("o", "out.root", "path to output ROOT file name")
	tname := flag.String("t", "tree", "name of the output tree")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `hepmc2root converts a HepMC2 ASCII file into a ROOT file and (flat) tree.

Usage: hepmc2root [OPTIONS] hepmc.ascii

Example:

$> hepmc2root ./hepmc.ascii
$> hepmc2root -o out.root -t mytree ./hepmc.ascii

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing input HepMC filename argument")
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
		return fmt.Errorf("could not open HepMC file %q: %w", fname, err)
	}
	defer f.Close()

	o, err := groot.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer o.Close()

	tree, err := rootcnv.NewFlatTreeWriter(o, tname, rtree.WithTitle(tname))
	if err != nil {
		return fmt.Errorf("could not create output ROOT tree %q: %w", tname, err)
	}

	_, err = hepmc.Copy(tree, hepmc.NewASCIIReader(f))
	if err != nil {
		return fmt.Errorf("could not write HepMC events to ROOT: %w", err)
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
