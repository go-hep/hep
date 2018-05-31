// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"go/importer"
	"go/types"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var (
		typeNames = flag.String("t", "", "comma-separated list of type names")
		pkgPath   = flag.String("p", "", "package import path")
	)

	flag.Parse()

	log.SetPrefix("gen-xrd: ")
	log.SetFlags(0)

	if *typeNames == "" {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")
	g, err := NewGenerator(*pkgPath)
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

	_, err = io.Copy(os.Stdout, bytes.NewReader(buf))
	if err != nil {
		log.Fatalf("error generating (un)marshaler code: %v\n", err)
	}
}

// Generator holds the state of the generation.
type Generator struct {
	buf *bytes.Buffer
	pkg *types.Package

	Verbose bool // enable verbose mode
}

// NewGenerator returns a new code generator for package p,
// where p is the package's import path.
func NewGenerator(p string) (*Generator, error) {
	pkg, err := importer.Default().Import(p)
	if err != nil {
		return nil, err
	}

	return &Generator{
		buf: new(bytes.Buffer),
		pkg: pkg,
	}, nil
}

func (g *Generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *Generator) Generate(typeName string) {
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
	if g.Verbose {
		log.Printf("typ: %+v\n", typ)
	}

	g.genMarshalXrd(typ, typeName)
	g.genUnmarshalXrd(typ, typeName)
}

func (g *Generator) genMarshalXrd(t types.Type, typeName string) {
	g.printf(`// MarshalXrd implements xrootd/protocol.Marshaler
func (o *%[1]s) MarshalXrd() (data []byte, err error) {
	var buf [8]byte
`,
		typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		g.genMarshalType(ft.Type(), "o."+ft.Name())
	}

	g.printf("return data, err\n}\n\n")
}

func (g *Generator) genMarshalType(t types.Type, n string) {

	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {

		case types.Bool:
			g.printf("if %s { data = append(data, uint8(1))\n", n)
			g.printf("}else { data = append(data, uint8(0)) }\n")

		case types.Uint:
			g.printf("binary.LittleEndian.PutUint64(buf[:8], uint64(%s))\n", n)
			g.printf("data = append(data, buf[:8]...)\n")

		case types.Uint8:
			g.printf("data = append(data, byte(%s))\n", n)

		case types.Uint16:
			g.printf(
				"binary.LittleEndian.PutUint16(buf[:2], %s)\n",
				n,
			)
			g.printf("data = append(data, buf[:2]...)\n")

		case types.Uint32:
			g.printf(
				"binary.LittleEndian.PutUint32(buf[:4], %s)\n",
				n,
			)
			g.printf("data = append(data, buf[:4]...)\n")

		case types.Uint64:
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:8], %s)\n",
				n,
			)
			g.printf("data = append(data, buf[:8]...)\n")

		case types.Int:
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:8], uint64(%s))\n",
				n,
			)
			g.printf("data = append(data, buf[:8]...)\n")

		case types.Int8:
			g.printf("data = append(data, byte(%s))\n", n)

		case types.Int16:
			g.printf(
				"binary.LittleEndian.PutUint16(buf[:2], uint16(%s))\n",
				n,
			)
			g.printf("data = append(data, buf[:2]...)\n")

		case types.Int32:
			g.printf(
				"binary.LittleEndian.PutUint32(buf[:4], uint32(%s))\n",
				n,
			)
			g.printf("data = append(data, buf[:4]...)\n")

		case types.Int64:
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:8], uint64(%s))\n",
				n,
			)
			g.printf("data = append(data, buf[:8]...)\n")

		case types.Float32:
			g.printf(
				"binary.LittleEndian.PutUint32(buf[:4], math.Float32bits(%s))\n",
				n,
			)
			g.printf("data = append(data, buf[:4]...)\n")

		case types.Float64:
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(%s))\n",
				n,
			)
			g.printf("data = append(data, buf[:8]...)\n")

		case types.Complex64:
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:4], math.Float32bits(real(%s)))\n",
				n,
			)
			g.printf("data = append(data, buf[:4]...)\n")
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:4], math.Float32bits(imag(%s)))\n",
				n,
			)
			g.printf("data = append(data, buf[:4]...)\n")

		case types.Complex128:
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(real(%s)))\n",
				n,
			)
			g.printf("data = append(data, buf[:8]...)\n")
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(imag(%s)))\n",
				n,
			)
			g.printf("data = append(data, buf[:8]...)\n")

		case types.String:
			g.printf(
				"binary.LittleEndian.PutUint64(buf[:8], uint64(len(%s)))\n",
				n,
			)
			g.printf("data = append(data, buf[:8]...)\n")
			g.printf("data = append(data, []byte(%s)...)\n", n)

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Struct:
		g.printf("{\nsub, err := %s.MarshalBinary()\n", n)
		g.printf("if err != nil {\nreturn nil, err\n}\n")
		g.printf("binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))\n")
		g.printf("data = append(data, buf[:8]...)\n")
		g.printf("data = append(data, sub...)\n")
		g.printf("}\n")

	case *types.Array:
		if isByteType(ut.Elem()) {
			g.printf("data = append(data, %s[:]...)\n", n)
		} else {
			g.printf("for i := range %s {\n", n)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.printf("o := %s[i]\n", n)
			} else {
				g.printf("o := &%s[i]\n", n)
			}
			g.genMarshalType(ut.Elem(), "o")
			g.printf("}\n")
		}

	case *types.Slice:
		g.printf(
			"binary.LittleEndian.PutUint64(buf[:8], uint64(len(%s)))\n",
			n,
		)
		g.printf("data = append(data, buf[:8]...)\n")
		if isByteType(ut.Elem()) {
			g.printf("data = append(data, %s...)\n", n)
		} else {
			g.printf("for i := range %s {\n", n)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.printf("o := %s[i]\n", n)
			} else {
				g.printf("o := &%s[i]\n", n)
			}
			g.genMarshalType(ut.Elem(), "o")
			g.printf("}\n")
		}

	case *types.Pointer:
		g.printf("{\n")
		g.printf("v := *%s\n", n)
		g.genUnmarshalXrd(ut.Elem(), "v")
		g.printf("}\n")

	case *types.Interface:
		log.Fatalf("marshal interface not supported (type=%v)\n", t)

	default:
		log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
	}
}

func (g *Generator) genUnmarshalXrd(t types.Type, typeName string) {
	g.printf(`// UnmarshalXrd implements xrootd/protocol.Unmarshaler
func (o *%[1]s) UnmarshalXrd(data []byte) (err error) {
`,
		typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		g.genUnmarshalType(ft.Type(), "o."+ft.Name())
	}

	g.printf("return err\n}\n\n")
}

func (g *Generator) genUnmarshalType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {

		case types.Bool:
			g.printf("if data[i] == 1 { %s = true }\n", n)
			g.printf("}else { %s = false }\n", n)
			g.printf("data = data[1:]\n")

		case types.Uint:
			g.printf("%s = uint(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.printf("data = data[8:]\n")

		case types.Uint8:
			g.printf("%s = data[0]\n", n)
			g.printf("data = data[1:]\n")

		case types.Uint16:
			g.printf("%s = binary.LittleEndian.Uint16(data[:2])\n", n)
			g.printf("data = data[2:]\n")

		case types.Uint32:
			g.printf("%s = binary.LittleEndian.Uint32(data[:4])\n", n)
			g.printf("data = data[4:]\n")

		case types.Uint64:
			g.printf("%s = binary.LittleEndian.Uint64(data[:8])\n", n)
			g.printf("data = data[8:]\n")

		case types.Int:
			g.printf("%s = int(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.printf("data = data[8:]\n")

		case types.Int8:
			g.printf("%s = int8(data[0])\n", n)
			g.printf("data = data[1:]\n")

		case types.Int16:
			g.printf("%s = int16(binary.LittleEndian.Uint16(data[:2]))\n", n)
			g.printf("data = data[2:]\n")

		case types.Int32:
			g.printf("%s = int32(binary.LittleEndian.Uint32(data[:4]))\n", n)
			g.printf("data = data[4:]\n")

		case types.Int64:
			g.printf("%s = int64(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.printf("data = data[8:]\n")

		case types.Float32:
			g.printf("%s = math.Float32frombits(binary.LittleEndian.Uint32(data[:4]))\n", n)
			g.printf("data = data[4:]\n")

		case types.Float64:
			g.printf("%s = math.Float64frombits(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.printf("data = data[8:]\n")

		case types.Complex64:
			g.printf("%s = complex(math.Float32frombits(binary.LittleEndian.Uint32(data[:4])), math.Float32frombits(binary.LittleEndian.Uint32(data[4:8])))\n", n)
			g.printf("data = data[8:]\n")

		case types.Complex128:
			g.printf("%s = complex(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])), math.Float64frombits(binary.LittleEndian.Uint64(data[8:16])))\n", n)
			g.printf("data = data[16:]\n")

		case types.String:
			g.printf("{\n")
			g.printf("n := int(binary.LittleEndian.Uint64(data[:8]))\n")
			g.printf("data = data[8:])\n")
			g.printf("%s = string(data[:n])\n", n)
			g.printf("data = data[n:]\n")
			g.printf("}\n")

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Struct:
		g.printf("{\n")
		g.printf("n := int(binary.LittleEndian.Uint64(data[:8]))\n")
		g.printf("data = data[8:]\n")
		g.printf("err = %s.UnmarshalBinary(data[:n])\n", n)
		g.printf("if err != nil {\nreturn err\n}\n")
		g.printf("data = data[n:]\n")
		g.printf("}\n")

	case *types.Array:
		if isByteType(ut.Elem()) {
			g.printf("copy(%s[:], data[:n])\n", n)
			g.printf("data = data[n:]\n")
		} else {
			g.printf("for i := range %s {\n", n)
			nn := n + "[i]"
			if pt, ok := ut.Elem().(*types.Pointer); ok {
				g.printf("var oi %s\n", qualTypeName(pt.Elem(), g.pkg))
				nn = "oi"
			}
			if _, ok := ut.Elem().Underlying().(*types.Struct); ok {
				g.printf("oi := &%s[i]\n", n)
				nn = "oi"
			}
			g.genUnmarshalType(ut.Elem(), nn)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.printf("%s[i] = oi\n", n)
			}
			g.printf("}\n")
		}

	case *types.Slice:
		g.printf("{\n")
		g.printf("n := int(binary.LittleEndian.Uint64(data[:8]))\n")
		g.printf("%[1]s = make([]%[2]s, n)\n", n, qualTypeName(ut.Elem(), g.pkg))
		g.printf("data = data[8:]\n")
		if isByteType(ut.Elem()) {
			g.printf("%[1]s = append(%[1]s, data[:n]...)\n", n)
			g.printf("data = data[n:]\n")
		} else {
			g.printf("for i := range %s {\n", n)
			nn := n + "[i]"
			if pt, ok := ut.Elem().(*types.Pointer); ok {
				g.printf("var oi %s\n", qualTypeName(pt.Elem(), g.pkg))
				nn = "oi"
			}
			if _, ok := ut.Elem().Underlying().(*types.Struct); ok {
				g.printf("oi := &%s[i]\n", n)
				nn = "oi"
			}
			g.genUnmarshalType(ut.Elem(), nn)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.printf("%s[i] = oi\n", n)
			}
			g.printf("}\n")
		}
		g.printf("}\n")

	case *types.Pointer:
		g.printf("{\n")
		elt := ut.Elem()
		g.printf("var v %s\n", qualTypeName(elt, g.pkg))
		g.genUnmarshalXrd(elt, "v")
		g.printf("%s = &v\n\n", n)
		g.printf("}\n")

	case *types.Interface:
		log.Fatalf("marshal interface not supported (type=%v)\n", t)

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

func qualTypeName(t types.Type, pkg *types.Package) string {
	n := types.TypeString(t, types.RelativeTo(pkg))
	i := strings.LastIndex(n, "/")
	if i < 0 {
		return n
	}
	return string(n[i+1:])
}

func (g *Generator) Format() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(g.buf.Bytes())

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("=== error ===\n%s\n", buf.Bytes())
	}
	return src, err
}
