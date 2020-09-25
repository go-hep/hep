// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root2fits converts the content of a ROOT tree to a FITS (binary) table.
//
// Usage: root2fits [OPTIONS] -f input.root
//
// Example:
//
//  $> root2fits -f ./input.root -t tree
//
// Options:
//   -f string
//     	path to input ROOT file name
//   -o string
//     	path to output FITS file name (default "output.fits")
//   -t string
//     	name of the ROOT tree to convert
//
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/astrogo/fitsio"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	log.SetPrefix("root2fits: ")
	log.SetFlags(0)

	fname := flag.String("f", "", "path to input ROOT file name")
	oname := flag.String("o", "output.fits", "path to output FITS file name")
	tname := flag.String("t", "", "name of the ROOT tree to convert")

	flag.Usage = func() {
		fmt.Printf(`root2fits converts the content of a ROOT tree to a FITS (binary) table.

Usage: root2fits [OPTIONS] -f input.root

Example:

 $> root2fits -f ./input.root -t tree

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *fname == "" {
		flag.Usage()
		log.Fatalf("missing path to input ROOT file argument")
	}

	if *tname == "" {
		flag.Usage()
		log.Fatalf("missing ROOT tree name to convert")
	}

	err := process(*oname, *tname, *fname)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func process(oname, tname, fname string) error {
	f, err := groot.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open input ROOT file %q: %w", fname, err)
	}
	defer f.Close()

	obj, err := riofs.Dir(f).Get(tname)
	if err != nil {
		return fmt.Errorf("could not retrieve ROOT tree %q from file %q: %w", tname, fname, err)
	}

	tree, ok := obj.(rtree.Tree)
	if !ok {
		return fmt.Errorf("ROOT object %q from file %q is not a tree", tname, fname)
	}

	o, err := os.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output file %q: %w", oname, err)
	}
	defer o.Close()

	fits, err := fitsio.Create(o)
	if err != nil {
		return fmt.Errorf("could not create output FITS file %q: %w", oname, err)
	}
	defer fits.Close()

	phdu, err := fitsio.NewPrimaryHDU(nil)
	if err != nil {
		return fmt.Errorf("could not create primary HDU: %w", err)
	}
	err = fits.Write(phdu)
	if err != nil {
		return fmt.Errorf("could not write primary HDU: %w", err)
	}

	tbl, err := tableFrom(tree)
	if err != nil {
		return fmt.Errorf("could not create output FITS table: %w", err)
	}
	defer tbl.Close()

	rvars := rtree.NewReadVars(tree)
	wvars := make([]interface{}, len(rvars))
	for i, rvar := range rvars {
		wvars[i] = rvar.Value
	}

	r, err := rtree.NewReader(tree, rvars)
	if err != nil {
		return fmt.Errorf("could not create ROOT reader: %w", err)
	}
	defer r.Close()

	err = r.Read(func(ctx rtree.RCtx) error {
		err = tbl.Write(wvars...)
		if err != nil {
			return fmt.Errorf("could not write entry %d to FITS table: %w", ctx.Entry, err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("could not read input ROOT file: %w", err)
	}

	err = fits.Write(tbl)
	if err != nil {
		return fmt.Errorf("could not write output FITS table %q: %w", tname, err)
	}

	err = tbl.Close()
	if err != nil {
		return fmt.Errorf("could not close output FITS table %q: %w", tname, err)
	}

	err = fits.Close()
	if err != nil {
		return fmt.Errorf("could not close output FITS file %q: %w", oname, err)
	}

	err = o.Close()
	if err != nil {
		return fmt.Errorf("could not close output file %q: %w", oname, err)
	}
	return nil
}

func tableFrom(tree rtree.Tree) (*fitsio.Table, error) {
	rvars := rtree.NewReadVars(tree)
	cols := make([]fitsio.Column, len(rvars))
	for i, rvar := range rvars {
		cols[i] = colFrom(rvar, tree)
	}

	return fitsio.NewTable(tree.Name(), cols, fitsio.BINARY_TBL)
}

func colFrom(rvar rtree.ReadVar, tree rtree.Tree) fitsio.Column {
	var (
		rt     = reflect.TypeOf(rvar.Value).Elem()
		format = formatFrom(rt)
	)

	switch rt.Kind() {
	case reflect.String:
		var (
			br   = tree.Branch(rvar.Name)
			leaf = br.Leaf(rvar.Leaf).(*rtree.LeafC)
			max  = leaf.Maximum()
		)
		format = fmt.Sprintf("%d%s", max, format)
	}

	return fitsio.Column{
		Name:   rvar.Name,
		Format: format,
	}
}

func formatFrom(rt reflect.Type) string {
	var format = ""
	switch rt.Kind() {
	case reflect.Bool:
		format = "L"
	case reflect.Int8, reflect.Uint8:
		format = "B"
	case reflect.Int16, reflect.Uint16:
		format = "I"
	case reflect.Int32, reflect.Uint32:
		format = "J"
	case reflect.Int64, reflect.Uint64:
		format = "K"
	case reflect.Float32:
		format = "E"
	case reflect.Float64:
		format = "D"
	case reflect.String:
		format = "A"
	case reflect.Array:
		elmt := formatFrom(rt.Elem())
		format = fmt.Sprintf("%d%s", rt.Len(), elmt)
	default:
		panic(fmt.Errorf("invalid branch type %T", reflect.New(rt).Elem().Interface()))
	}

	return format
}
