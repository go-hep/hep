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
func (o %[1]s) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
`,
		typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		g.genMarshalType(ft.Type(), "o."+ft.Name())
	}

	g.printf("return nil\n}\n\n")
}

func (g *Generator) genMarshalType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {

		case types.Bool:
			g.printf("wBuffer.WriteBool(%s)\n", g.upcasted(t, n))

		case types.Uint8:
			if n == "o._" {
				g.printf("wBuffer.Next(1)\n")
			} else {
				g.printf("wBuffer.WriteU8(%s)\n", g.upcasted(t, n))
			}

		case types.Uint16:
			g.printf("wBuffer.WriteU16(%s)\n", g.upcasted(t, n))

		case types.Uint32:
			g.printf("wBuffer.WriteI32(int32(%s))\n", g.upcasted(t, n))

		case types.Uint64:
			g.printf("wBuffer.WriteI64(int64(%s))\n", g.upcasted(t, n))

		case types.Int8:
			g.printf("wBuffer.WriteU8(uint8(%s))\n", g.upcasted(t, n))

		case types.Int16:
			g.printf("wBuffer.WriteU16(uint16(%s))\n", g.upcasted(t, n))

		case types.Int32:
			if n == "o._" {
				g.printf("wBuffer.Next(4)\n")
			} else {
				g.printf("wBuffer.WriteI32(%s)\n", g.upcasted(t, n))
			}

		case types.Int64:
			g.printf("wBuffer.WriteI64(%s)\n", g.upcasted(t, n))

		case types.String:
			g.printf("wBuffer.WriteStr(%s)\n", g.upcasted(t, n))

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Struct:
		g.printf("if err := %s.MarshalXrd(wBuffer); err != nil {\nreturn err\n}\n", n)

	case *types.Array:
		if !isByteType(ut.Elem()) {
			log.Fatalf("marshal array of type %v not supported", ut)
		}
		if n == "o._" {
			g.printf("wBuffer.Next(%d)\n", ut.Len())
		} else {
			g.printf("wBuffer.WriteBytes(%s[:])\n", n)
		}

	case *types.Slice:
		if !isByteType(ut.Elem()) {
			g.printf("wBuffer.WriteLen(len(%s))\n", n)
			g.printf(`for _, x := range %s {
	err := x.MarshalXrd(wBuffer)
	if err != nil {
		return err
	}
}
`, n)
		} else {
			g.printf("wBuffer.WriteLen(len(%s))\n", n)
			g.printf("wBuffer.WriteBytes(%s)\n", n)
		}

	default:
		log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
	}
}

func (g *Generator) genUnmarshalXrd(t types.Type, typeName string) {
	g.printf(`// UnmarshalXrd implements xrootd/protocol.Unmarshaler
func (o *%[1]s) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
`,
		typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		g.genUnmarshalType(ft.Type(), "o."+ft.Name())
	}

	g.printf("return nil\n}\n\n")
}

func (g *Generator) downcasted(t types.Type, expression string) string {
	if named, ok := t.(*types.Named); ok {
		cast := qualTypeName(named, g.pkg)
		return cast + "(" + expression + ")"
	}
	return expression
}

func (g *Generator) upcasted(t types.Type, expression string) string {
	if named, ok := t.(*types.Named); ok {
		ut := named.Underlying()
		if basic, ok := ut.(*types.Basic); ok {
			cast := basic.Name()
			return cast + "(" + expression + ")"
		}
	}
	return expression
}

func (g *Generator) genUnmarshalType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {
		case types.Bool:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "rBuffer.ReadBool()"))

		case types.Uint:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "uint(rBuffer.ReadI64())"))

		case types.Uint8:
			if n == "o._" {
				g.printf("rBuffer.Skip(1)\n")
			} else {
				g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "rBuffer.ReadU8()"))
			}

		case types.Uint16:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "rBuffer.ReadU16()"))

		case types.Uint32:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "uint32(rBuffer.ReadI32())"))

		case types.Uint64:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "uint64(rBuffer.ReadI64())"))

		case types.Int8:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "int8(rBuffer.ReadU8())"))
		case types.Int16:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "int16(rBuffer.ReadU16())"))

		case types.Int32:
			if n == "o._" {
				g.printf("rBuffer.Skip(4)\n")
			} else {
				g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "rBuffer.ReadI32()"))
			}

		case types.Int64:
			g.printf("%[1]s = %[2]s\n", n, g.downcasted(t, "rBuffer.ReadI64()"))

		case types.String:
			g.printf("%s = %[2]s\n", n, g.downcasted(t, "rBuffer.ReadStr()"))

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Struct:
		g.printf("if err := %s.UnmarshalXrd(rBuffer); err != nil {\n return err\n}\n", n)

	case *types.Array:
		if !isByteType(ut.Elem()) {
			log.Fatalf("unmarshal array of type %v not supported", ut)
		}
		if n == "o._" {
			g.printf("rBuffer.Skip(%d)\n", ut.Len())
		} else {
			g.printf("rBuffer.ReadBytes(%s[:])\n", n)
		}

	case *types.Slice:
		if !isByteType(ut.Elem()) {
			g.printf("%[1]s = make([]%[2]s, rBuffer.ReadLen())\n", n, qualTypeName(ut.Elem(), g.pkg))
			g.printf(`for i:=0; i<len(%[1]s); i++ {
    err := %[1]s[i].UnmarshalXrd(rBuffer)
    if err != nil {
		return err
    }
}
`, n)
		} else {
			g.printf("%[1]s = make([]%[2]s, rBuffer.ReadLen())\n", n, qualTypeName(ut.Elem(), g.pkg))
			g.printf("rBuffer.ReadBytes(%s)\n", n)
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
