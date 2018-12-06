// Copyright 2016 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command brio-gen generates (un)marshaler code for types.
//
// For each type given to brio-gen, marshaling and unmarshaling code will be
// generated so the types implement binary.Binary(Un)Marshaler.
//
//  - values are encoded using binary.LittleEndian,
//  - int and uint are encoded as (resp.) int64 and uint64,
//  - booleans are encoded as a single uint8 (0==false, 1==true),
//  - strings are encoded as a pair(uint64, []byte),
//  - arrays are encoded as a sequence of Ts, the length is implicit as it is
//    part of the array type,
//  - slices are encoded as a pair(uint64, T...)
//  - pointers are encoded as *T (like encoding/gob)
//  - interfaces are not supported.
//
//
// Usage: brio-gen [options]
//
// ex:
//  $> brio-gen -p image -t Point -o image_brio.go
//  $> brio-gen -p go-hep.org/x/hep/hbook -t Dist0D,Dist1D,Dist2D -o foo_brio.go
//
// options:
//   -o string
//     	output file name (default "brio_gen.go")
//   -p string
//     	package import path
//   -t string
//     	comma-separated list of type names
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"go-hep.org/x/hep/brio/cmd/brio-gen/internal/gen"
)

var (
	typeNames = flag.String("t", "", "comma-separated list of type names")
	pkgPath   = flag.String("p", "", "package import path")
	output    = flag.String("o", "brio_gen.go", "output file name")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: brio-gen [options]

ex:
 $> brio-gen -p image -t Point -o image_brio.go
 $> brio-gen -p go-hep.org/x/hep/hbook -t Dist0D,Dist1D,Dist2D -o foo_brio.go

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	log.SetPrefix("brio: ")
	log.SetFlags(0)

	if *typeNames == "" {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")
	out, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	err = generate(out, *pkgPath, types)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func generate(w io.Writer, pkg string, types []string) error {
	g, err := gen.NewGenerator(pkg)
	if err != nil {
		return err
	}

	for _, t := range types {
		g.Generate(t)
	}

	buf, err := g.Format()
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}
