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

	// set of imported packages.
	// usually: "encoding/binary", "math"
	imps map[string]int
}

func newGenerator(path string) (*Generator, error) {
	pkg, err := importer.Default().Import(path)
	if err != nil {
		return nil, err
	}

	return &Generator{
		buf:  new(bytes.Buffer),
		pkg:  pkg,
		imps: map[string]int{"encoding/binary": 1},
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

	g.Printf(`// MarshalBinary implements encoding.BinaryMarshaler
func (o *%[1]s) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
`,
		typeName,
	)

	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		log.Printf("field[%d]: %v\n", i, ft)
		ut := ft.Type().Underlying()
		switch ut := ut.(type) {
		case *types.Basic:
			switch kind := ut.Kind(); kind {
			case types.Uint:
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:8], uint64(o.%s))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:8])\n")
			case types.Uint8:
				g.Printf("\tdata = append(data, byte(o.%s))\n", ft.Name())
			case types.Uint16:
				g.Printf(
					"\tbinary.LittleEndian.PutUint16(buf[:2], o.%s)\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:2])\n")
			case types.Uint32:
				g.Printf(
					"\tbinary.LittleEndian.PutUint32(buf[:4], o.%s)\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:4])\n")
			case types.Uint64:
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:8], o.%s)\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:8])\n")
			case types.Int:
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:8], int64(o.%s))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:8])\n")
			case types.Int8:
				g.Printf("\tdata = append(data, byte(o.%s))\n", ft.Name())
			case types.Int16:
				g.Printf(
					"\tbinary.LittleEndian.PutUint16(buf[:2], uint16(o.%s))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:2])\n")
			case types.Int32:
				g.Printf(
					"\tbinary.LittleEndian.PutUint32(buf[:4], uint32(o.%s))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:4])\n")
			case types.Int64:
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:8], uint64(o.%s))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:8])\n")
			case types.Float32:
				g.imps["math"]++
				g.Printf(
					"\tbinary.LittleEndian.PutUint32(buf[:4], math.Float32bits(o.%s))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:4])\n")
			case types.Float64:
				g.imps["math"]++
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.%s))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:8])\n")
			case types.Complex64:
				g.imps["math"]++
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:4], math.Float32bits(real(o.%s)))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:4])\n")
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:4], math.Float32bits(imag(o.%s)))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:4])\n")
			case types.Complex128:
				g.imps["math"]++
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:8], math.Float64bits(real(o.%s)))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:8])\n")
				g.Printf(
					"\tbinary.LittleEndian.PutUint64(buf[:8], math.Float64bits(imag(o.%s)))\n",
					ft.Name(),
				)
				g.Printf("\tdata = append(data, buf[:8])\n")
			}
		}
	}

	g.Printf("\treturn data, err\n}\n\n")
}
