// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rio2yoda converts rio files containing hbook values (H1D, H2D, P1D, ...)
// into YODA files.
//
// Example:
//
//  $> rio2yoda file1.rio file2.rio > out.yoda
//  $> rio2yoda -o out.yoda file1.rio file2.rio
//  $> rio2yoda -o out.yoda.gz file1.rio file2.rio
package main

import (
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/rio"
)

func main() {

	log.SetFlags(0)
	log.SetPrefix("rio2yoda: ")

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: rio2yoda [options] <file.rio> [<file2.rio> [...]]
	
ex:
 $> rio2yoda file1.rio > out.yoda
 $> rio2yoda -o out.yoda file1.rio file2.rio
 $> rio2yoda -o out.yoda.gz file1.rio file2.rio
 `)
	}

	oname := flag.String("o", "", "path to YODA output file")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Printf("missing input file name")
		flag.Usage()
		flag.PrintDefaults()
	}

	var out io.WriteCloser = os.Stdout
	if *oname != "" {
		f, err := os.Create(*oname)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		out = f
		if filepath.Ext(*oname) == ".gz" {
			wz := gzip.NewWriter(f)
			defer func() {
				err = wz.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()
			out = wz
		}
	} else {
		defer out.Close()
	}

	for _, fname := range flag.Args() {
		convert(out, fname)
	}
}

func convert(out io.Writer, fname string) {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r, err := rio.Open(f)
	if err != nil {
		log.Fatalf("error opening rio stream %q: %v\n", fname, err)
	}
	defer r.Close()

	for _, key := range r.Keys() {
		rt := typeFrom(key.Blocks[0].Type)
		if rt == nil {
			continue
		}
		v := reflect.New(rt.Elem())
		err = r.Get(key.Name, v.Interface())
		if err != nil {
			log.Fatalf(
				"error reading block %q from file %q: %v\n",
				key.Name, fname, err,
			)
		}

		yoda, err := v.Interface().(yodaMarshaler).MarshalYODA()
		if err != nil {
			log.Fatalf(
				"error marshaling block %q from file %q to YODA: %v\n",
				key.Name, fname, err,
			)
		}

		_, err = out.Write(yoda)
		if err != nil {
			log.Fatalf(
				"error streaming out YODA format for block %q from file %q: %v\n",
				key.Name, fname, err,
			)
		}
	}
}

func nameFromType(rt reflect.Type) string {
	if rt == nil {
		return "interface"
	}
	// Default to printed representation for unnamed types
	name := rt.String()

	// But for named types (or pointers to them), qualify with import path.
	// Dereference one pointer looking for a named type.
	star := ""
	if rt.Name() == "" {
		pt := rt
		if pt.Kind() == reflect.Ptr {
			star = "*"
			rt = pt.Elem()
		}
	}

	if rt.Name() != "" {
		switch rt.PkgPath() {
		case "":
			name = star + rt.Name()
		default:
			name = star + rt.PkgPath() + "." + rt.Name()
		}
	}

	return name
}

func typeFrom(name string) reflect.Type {
	for _, t := range hbookTypes {
		if name == nameFromType(t) {
			return t
		}
	}
	return nil
}

var hbookTypes = []reflect.Type{
	reflect.TypeOf((*hbook.H1D)(nil)),
	reflect.TypeOf((*hbook.H2D)(nil)),
	reflect.TypeOf((*hbook.P1D)(nil)),
	reflect.TypeOf((*hbook.S2D)(nil)),
}

type yodaMarshaler interface {
	MarshalYODA() ([]byte, error)
}
