// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command brio-gen generates (un)marshaler code for types.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/importer"
	"go/types"
	"io"
	"log"
	"os"
	"strings"
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
	g, err := newGenerator(*pkgPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range types {
		g.generate(t)
	}

	// FIXME(sbinet): debug
	io.Copy(os.Stdout, g.buf)
}

// Generator holds the state of the generation.
type Generator struct {
	buf *bytes.Buffer
	pkg *types.Package
}

func newGenerator(path string) (*Generator, error) {
	pkg, err := importer.Default().Import(path)
	if err != nil {
		return nil, err
	}

	return &Generator{
		buf: new(bytes.Buffer),
		pkg: pkg,
	}, nil
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *Generator) generate(typeName string) {
	scope := g.pkg.Scope()
	obj := scope.Lookup(typeName)
	if obj == nil {
		log.Fatalf("no such type %q in package %q\n", typeName, g.pkg.Path()+"/"+g.pkg.Name())
	}

	tn, ok := obj.(*types.TypeName)
	if !ok {
		log.Fatalf("%q is not a type (%v)\n", typeName, obj)
	}

	typ, ok := tn.Type().Underlying().(*types.Struct)
	if !ok {
		log.Fatalf("%q is not a named struct (%v)\n", typeName, tn)
	}
	log.Printf("typ: %+v\n", typ)

	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		log.Printf("field[%d]: %v\n", i, ft)
	}

}
