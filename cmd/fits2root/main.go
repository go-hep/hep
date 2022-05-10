// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// fits2root converts the content of a FITS table to a ROOT file and tree.
//
// Usage: fits2root [OPTIONS] -f input.fits
//
// Example:
//
//	$> fits2root -f ./input.fits -t MyHDU
//
// Options:
//
//	-f string
//	  	path to input FITS file name
//	-o string
//	  	path to output ROOT file name (default "output.root")
//	-t string
//	  	name of the FITS table to convert
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/astrogo/fitsio"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	log.SetPrefix("fits2root: ")
	log.SetFlags(0)

	fname := flag.String("f", "", "path to input FITS file name")
	oname := flag.String("o", "output.root", "path to output ROOT file name")
	tname := flag.String("t", "", "name of the FITS table to convert")

	flag.Usage = func() {
		fmt.Printf(`fits2root converts the content of a FITS table to a ROOT file and tree.

Usage: fits2root [OPTIONS] -f input.fits

Example:

 $> fits2root -f ./input.fits -t MyHDU

Options:
`)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *fname == "" {
		flag.Usage()
		log.Fatalf("missing path to input FITS file argument")
	}

	err := process(*oname, *tname, *fname)
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func process(oname, tname, fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return fmt.Errorf("could not open input file %q: %w", fname, err)
	}
	defer f.Close()

	fits, err := fitsio.Open(f)
	if err != nil {
		return fmt.Errorf("could not open FITS file %q: %w", fname, err)
	}
	defer fits.Close()

	hdu := fits.Get(tname)
	if hdu == nil {
		return fmt.Errorf("could not retrieve HDU %q from FITS file %q", tname, fname)
	}
	defer hdu.Close()

	tbl, ok := hdu.(*fitsio.Table)
	if !ok {
		return fmt.Errorf("HDU %q from FITS file %q is not a Table", tname, fname)
	}

	o, err := groot.Create(oname)
	if err != nil {
		return fmt.Errorf("could not create output ROOT file %q: %w", oname, err)
	}
	defer o.Close()

	var (
		wvars = make([]rtree.WriteVar, tbl.NumCols())
		wargs = make([]interface{}, tbl.NumCols())
	)
	for i, col := range tbl.Cols() {
		wvars[i] = wvarFrom(col)
		wargs[i] = wvars[i].Value
	}

	if tname == "" {
		tname = "FITS"
	}

	tree, err := rtree.NewWriter(o, tname, wvars)
	if err != nil {
		return fmt.Errorf("could not create output ROOT tree %q: %w", tname, err)
	}

	rows, err := tbl.Read(0, tbl.NumRows())
	if err != nil {
		return fmt.Errorf("could not read FITS table range [0, %d): %w", tbl.NumRows(), err)
	}
	defer rows.Close()

	var irow int64
	for rows.Next() {
		err = rows.Scan(wargs...)
		if err != nil {
			return fmt.Errorf("could not read row %d: %w", irow, err)
		}

		_, err = tree.Write()
		if err != nil {
			return fmt.Errorf("could not write row %d: %w", irow, err)
		}

		irow++
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

func wvarFrom(col fitsio.Column) rtree.WriteVar {
	rt := col.Type()
	switch rt.Kind() {
	case reflect.Bool:
		return rtree.WriteVar{
			Name:  col.Name,
			Value: new(bool),
		}

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rtree.WriteVar{
			Name:  col.Name,
			Value: reflect.New(rt).Interface(),
		}

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rtree.WriteVar{
			Name:  col.Name,
			Value: reflect.New(rt).Interface(),
		}

	case reflect.Float32, reflect.Float64:
		return rtree.WriteVar{
			Name:  col.Name,
			Value: reflect.New(rt).Interface(),
		}

	case reflect.String:
		return rtree.WriteVar{
			Name:  col.Name,
			Value: new(string),
		}

	case reflect.Array:
		return rtree.WriteVar{
			Name:  col.Name,
			Value: reflect.New(rt).Interface(),
		}

	default:
		panic(fmt.Errorf("fits2root: invalid column type %+v (kind=%v)", col, rt.Kind()))
	}
}
