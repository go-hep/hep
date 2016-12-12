// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command brio-gen generates (un)marshaler code for types.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/importer"
	"go/types"
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

	buf, err := format.Source(g.buf.Bytes())
	if err != nil {
		log.Fatalf("gofmt: %v\n", err)
	}
	// FIXME(sbinet): debug
	os.Stdout.Write(buf)
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
		g.genMarshalType(ft.Type(), "o."+ft.Name())
	}

	g.Printf("\treturn data, err\n}\n\n")
}

func (g *Generator) genMarshalType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {

		case types.Bool:
			g.Printf("\tif %s { data = append(data, uint8(1))\n", n)
			g.Printf("\t}else { data = append(data, uint8(0)) }\n")

		case types.Uint:
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], uint64(%s))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")

		case types.Uint8:
			g.Printf("\tdata = append(data, byte(%s))\n", n)

		case types.Uint16:
			g.Printf(
				"\tbinary.LittleEndian.PutUint16(buf[:2], %s)\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:2])\n")

		case types.Uint32:
			g.Printf(
				"\tbinary.LittleEndian.PutUint32(buf[:4], %s)\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:4])\n")

		case types.Uint64:
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], %s)\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")

		case types.Int:
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], int64(%s))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")

		case types.Int8:
			g.Printf("\tdata = append(data, byte(%s))\n", n)

		case types.Int16:
			g.Printf(
				"\tbinary.LittleEndian.PutUint16(buf[:2], uint16(%s))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:2])\n")

		case types.Int32:
			g.Printf(
				"\tbinary.LittleEndian.PutUint32(buf[:4], uint32(%s))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:4])\n")

		case types.Int64:
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], uint64(%s))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")

		case types.Float32:
			g.imps["math"]++
			g.Printf(
				"\tbinary.LittleEndian.PutUint32(buf[:4], math.Float32bits(%s))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:4])\n")

		case types.Float64:
			g.imps["math"]++
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], math.Float64bits(%s))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")

		case types.Complex64:
			g.imps["math"]++
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:4], math.Float32bits(real(%s)))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:4])\n")
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:4], math.Float32bits(imag(%s)))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:4])\n")

		case types.Complex128:
			g.imps["math"]++
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], math.Float64bits(real(%s)))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], math.Float64bits(imag(%s)))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")

		case types.String:
			g.Printf(
				"\tbinary.LittleEndian.PutUint64(buf[:8], uint64(len(%s)))\n",
				n,
			)
			g.Printf("\tdata = append(data, buf[:8])\n")
			g.Printf("\tdata = append(data, []byte(%s)...)\n")
		}

	case *types.Struct:
		nn := strings.Replace(n, ".", "_", -1)
		g.Printf("\tdata_%s, err := %s.MarshalBinary()\n", nn, n)
		g.Printf("\tif err != nil {\n\t\treturn nil, err\n\t}\n")
		g.Printf("\tdata = append(data, data_%s...)\n", nn)

	case *types.Array:
		g.Printf("\tfor i := range %s {\n", n)
		if _, ok := ut.Elem().(*types.Pointer); ok {
			g.Printf("\t\to := %s[i]\n", n)
		} else {
			g.Printf("\t\to := &%s[i]\n", n)
		}
		g.genMarshalType(ut.Elem(), "o")
		g.Printf("\t}\n")

	case *types.Slice:
		g.Printf(
			"\tbinary.LittleEndian.PutUint64(buf[:8], uint64(len(%s)))\n",
			n,
		)
		g.Printf("\tdata = append(data, buf[:8])\n")
		if isByteType(ut.Elem()) {
			g.Printf("\tdata = append(data, %s...)\n", n)
		} else {
			g.Printf("\tfor i := range %s {\n", n)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.Printf("\t\to := %s[i]\n", n)
			} else {
				g.Printf("\t\to := &%s[i]\n", n)
			}
			g.genMarshalType(ut.Elem(), "o")
			g.Printf("\t}\n")
		}
	default:
		log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
	}
}

func isByteType(t types.Type) bool {
	b, ok := t.Underlying().(*types.Basic)
	if !ok {
		return false
	}
	return b.Kind() == types.Byte
}
