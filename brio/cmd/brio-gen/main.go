// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command brio-gen generates (un)marshaler code for types.
package main

import (
	"flag"
	"io/ioutil"
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
	flag.Parse()

	log.SetPrefix("brio: ")
	log.SetFlags(0)

	if *typeNames == "" {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")
	g, err := gen.NewGenerator(*pkgPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range types {
		g.Generate(t)
	}

	buf, err := g.Format()
	if err != nil {
		log.Fatalf("gofmt: %v\n", err)
	}

	err = ioutil.WriteFile(*output, buf, 0644)
	if err != nil {
		log.Fatalf("error writing file [%s]: %v\n", *output, err)
	}
}
