// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gen

import (
	"bytes"
	"fmt"
	"go/format"
	"go/importer"
	"go/types"
	"log"
	"os"
	"strings"
)

// Generator holds the state of the generation.
type Generator struct {
	buf *bytes.Buffer
	pkg *types.Package

	// set of imported packages.
	// usually: "encoding/binary", "math"
	imps map[string]int

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
		buf:  new(bytes.Buffer),
		pkg:  pkg,
		imps: map[string]int{"encoding/binary": 1},
	}, nil
}

func (g *Generator) Printf(format string, args ...interface{}) {
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

	g.genMarshal(typ, typeName)
	g.genUnmarshal(typ, typeName)
}

func (g *Generator) genMarshal(t types.Type, typeName string) {
	g.Printf(`// MarshalBinary implements encoding.BinaryMarshaler
func (o *%[1]s) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
`,
		typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		g.genMarshalType(ft.Type(), "o."+ft.Name())
	}

	g.Printf("return data, err\n}\n\n")
}

func (g *Generator) genMarshalType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {

		case types.Bool:
			g.Printf("if %s { data = append(data, uint8(1))\n", n)
			g.Printf("}else { data = append(data, uint8(0)) }\n")

		case types.Uint:
			g.Printf("binary.LittleEndian.PutUint64(buf[:8], uint64(%s))\n", n)
			g.Printf("data = append(data, buf[:8])\n")

		case types.Uint8:
			g.Printf("data = append(data, byte(%s))\n", n)

		case types.Uint16:
			g.Printf(
				"binary.LittleEndian.PutUint16(buf[:2], %s)\n",
				n,
			)
			g.Printf("data = append(data, buf[:2])\n")

		case types.Uint32:
			g.Printf(
				"binary.LittleEndian.PutUint32(buf[:4], %s)\n",
				n,
			)
			g.Printf("data = append(data, buf[:4])\n")

		case types.Uint64:
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:8], %s)\n",
				n,
			)
			g.Printf("data = append(data, buf[:8])\n")

		case types.Int:
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:8], int64(%s))\n",
				n,
			)
			g.Printf("data = append(data, buf[:8])\n")

		case types.Int8:
			g.Printf("data = append(data, byte(%s))\n", n)

		case types.Int16:
			g.Printf(
				"binary.LittleEndian.PutUint16(buf[:2], uint16(%s))\n",
				n,
			)
			g.Printf("data = append(data, buf[:2])\n")

		case types.Int32:
			g.Printf(
				"binary.LittleEndian.PutUint32(buf[:4], uint32(%s))\n",
				n,
			)
			g.Printf("data = append(data, buf[:4])\n")

		case types.Int64:
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:8], uint64(%s))\n",
				n,
			)
			g.Printf("data = append(data, buf[:8])\n")

		case types.Float32:
			g.imps["math"] = 1
			g.Printf(
				"binary.LittleEndian.PutUint32(buf[:4], math.Float32bits(%s))\n",
				n,
			)
			g.Printf("data = append(data, buf[:4])\n")

		case types.Float64:
			g.imps["math"] = 1
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(%s))\n",
				n,
			)
			g.Printf("data = append(data, buf[:8])\n")

		case types.Complex64:
			g.imps["math"] = 1
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:4], math.Float32bits(real(%s)))\n",
				n,
			)
			g.Printf("data = append(data, buf[:4])\n")
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:4], math.Float32bits(imag(%s)))\n",
				n,
			)
			g.Printf("data = append(data, buf[:4])\n")

		case types.Complex128:
			g.imps["math"] = 1
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(real(%s)))\n",
				n,
			)
			g.Printf("data = append(data, buf[:8])\n")
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(imag(%s)))\n",
				n,
			)
			g.Printf("data = append(data, buf[:8])\n")

		case types.String:
			g.Printf(
				"binary.LittleEndian.PutUint64(buf[:8], uint64(len(%s)))\n",
				n,
			)
			g.Printf("data = append(data, buf[:8])\n")
			g.Printf("data = append(data, []byte(%s)...)\n")

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Struct:
		g.Printf("{\nsub, err := %s.MarshalBinary()\n", n)
		g.Printf("if err != nil {\nreturn nil, err\n}\n")
		g.Printf("binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))\n")
		g.Printf("data = append(data, sub...)\n")
		g.Printf("}\n")

	case *types.Array:
		if isByteType(ut.Elem()) {
			g.Printf("data = append(data, %s[:]...)\n", n)
		} else {
			g.Printf("for i := range %s {\n", n)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.Printf("o := %s[i]\n", n)
			} else {
				g.Printf("o := &%s[i]\n", n)
			}
			g.genMarshalType(ut.Elem(), "o")
			g.Printf("}\n")
		}

	case *types.Slice:
		g.Printf(
			"binary.LittleEndian.PutUint64(buf[:8], uint64(len(%s)))\n",
			n,
		)
		g.Printf("data = append(data, buf[:8])\n")
		if isByteType(ut.Elem()) {
			g.Printf("data = append(data, %s...)\n", n)
		} else {
			g.Printf("for i := range %s {\n", n)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.Printf("o := %s[i]\n", n)
			} else {
				g.Printf("o := &%s[i]\n", n)
			}
			g.genMarshalType(ut.Elem(), "o")
			g.Printf("}\n")
		}

	case *types.Pointer:
		g.Printf("{\n")
		g.Printf("v := *%s\n", n)
		g.genUnmarshal(ut.Elem(), "v")
		g.Printf("}\n")

	default:
		log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
	}
}

func (g *Generator) genUnmarshal(t types.Type, typeName string) {
	g.Printf(`// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *%[1]s) UnmarshalBinary(data []byte) (err error) {
`,
		typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		g.genUnmarshalType(ft.Type(), "o."+ft.Name())
	}

	g.Printf("return err\n}\n\n")
}

func (g *Generator) genUnmarshalType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {

		case types.Bool:
			g.Printf("if data[i] == 1 { %s = true }\n", n)
			g.Printf("}else { %s = false }\n", n)
			g.Printf("data = data[1:]\n")

		case types.Uint:
			g.Printf("%s = uint(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.Printf("data = data[8:]\n")

		case types.Uint8:
			g.Printf("%s = data[0]\n", n)
			g.Printf("data = data[1:]\n")

		case types.Uint16:
			g.Printf("%s = binary.LittleEndian.Uint16(data[:2])\n", n)
			g.Printf("data = data[2:]\n")

		case types.Uint32:
			g.Printf("%s = binary.LittleEndian.Uint32(data[:4])\n", n)
			g.Printf("data = data[4:]\n")

		case types.Uint64:
			g.Printf("%s = binary.LittleEndian.Uint64(data[:8])\n", n)
			g.Printf("data = data[8:]\n")

		case types.Int:
			g.Printf("%s = int(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.Printf("data = data[8:]\n")

		case types.Int8:
			g.Printf("%s = int8(data[0])\n", n)
			g.Printf("data = data[1:]\n")

		case types.Int16:
			g.Printf("%s = int16(binary.LittleEndian.Uint16(data[:2]))\n", n)
			g.Printf("data = data[2:]\n")

		case types.Int32:
			g.Printf("%s = int32(binary.LittleEndian.Uint32(data[:4]))\n", n)
			g.Printf("data = data[4:]\n")

		case types.Int64:
			g.Printf("%s = int64(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.Printf("data = data[8:]\n")

		case types.Float32:
			g.imps["math"] = 1
			g.Printf("%s = math.Float32frombits(binary.LittleEndian.Uint32(data[:4]))\n", n)
			g.Printf("data = data[4:]\n")

		case types.Float64:
			g.imps["math"] = 1
			g.Printf("%s = math.Float64frombits(binary.LittleEndian.Uint64(data[:8]))\n", n)
			g.Printf("data = data[8:]\n")

		case types.Complex64:
			g.imps["math"] = 1
			g.Printf("%s = complex(math.Float32frombits(binary.LittleEndian.Uint32(data[:4])), math.Float32frombits(binary.LittleEndian.Uint32(data[4:8])))\n", n)
			g.Printf("data = data[8:]\n")

		case types.Complex128:
			g.imps["math"] = 1
			g.Printf("%s = complex(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])), math.Float64frombits(binary.LittleEndian.Uint64(data[8:16])))\n", n)
			g.Printf("data = data[16:]\n")

		case types.String:
			g.Printf("{\n")
			g.Printf("n := int(binary.LittleEndian.Uint64(data[:8]))\n")
			g.Printf("data = data[8:])\n")
			g.Printf("%s = string(data[:n])\n", n)
			g.Printf("data = data[n:]\n")
			g.Printf("}\n")

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Struct:
		g.Printf("{\n")
		g.Printf("n := int(binary.LittleEndian.Uint64(data[:8]))\n")
		g.Printf("data = data[8:]\n")
		g.Printf("err = %s.BinaryUnmarshal(data[:n])\n", n)
		g.Printf("if err != nil {\nreturn err\n}\n")
		g.Printf("data = data[n:]\n")
		g.Printf("}\n")

	case *types.Array:
		if isByteType(ut.Elem()) {
			g.Printf("copy(%s[:], data[:n])\n", n)
			g.Printf("data = data[n:]\n")
		} else {
			g.Printf("for i := range %s {\n", n)
			nn := n + "[i]"
			if pt, ok := ut.Elem().(*types.Pointer); ok {
				g.Printf("var oi %s\n", qualTypeName(pt.Elem()))
				nn = "oi"
			}
			if _, ok := ut.Elem().Underlying().(*types.Struct); ok {
				g.Printf("oi := &%s[i]\n", n)
				nn = "oi"
			}
			g.genUnmarshalType(ut.Elem(), nn)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.Printf("%s[i] = oi\n", n)
			}
			g.Printf("}\n")
		}

	case *types.Slice:
		g.Printf("{\n")
		g.Printf("n := int(binary.LittleEndian.Uint64(data[:8]))\n")
		g.Printf("data = data[8:]\n")
		if isByteType(ut.Elem()) {
			g.Printf("%[1]s = append(%[1]s, data[:n]...)\n", n)
			g.Printf("data = data[n:]\n")
		} else {
			g.Printf("for i := range %s {\n", n)
			nn := n + "[i]"
			if pt, ok := ut.Elem().(*types.Pointer); ok {
				g.Printf("var oi %s\n", qualTypeName(pt.Elem()))
				nn = "oi"
			}
			if _, ok := ut.Elem().Underlying().(*types.Struct); ok {
				g.Printf("oi := &%s[i]\n", n)
				nn = "oi"
			}
			g.genUnmarshalType(ut.Elem(), nn)
			if _, ok := ut.Elem().(*types.Pointer); ok {
				g.Printf("%s[i] = oi\n", n)
			}
			g.Printf("}\n")
		}
		g.Printf("}\n")

	case *types.Pointer:
		g.Printf("{\n")
		elt := ut.Elem()
		g.Printf("var v %s\n", qualTypeName(elt))
		g.genUnmarshal(elt, "v")
		g.Printf("%s = &v\n\n", n)
		g.Printf("}\n")

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

func qualTypeName(t types.Type) string {
	n := types.TypeString(t, nil)
	i := strings.LastIndex(n, "/")
	if i < 0 {
		return n
	}
	return string(n[i+1:])
}

func (g *Generator) Format() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf(`// DO NOT EDIT; automatically generated by %[1]s

package %[2]s

import (
	"encoding/binary"
`,
		os.Args[0],
		g.pkg.Name(),
	))

	for k := range g.imps {
		fmt.Fprintf(buf, "%q\n", k)
	}
	fmt.Fprintf(buf, ")\n\n")

	buf.Write(g.buf.Bytes())

	src, err := format.Source(buf.Bytes())
	if err != nil {
		log.Printf("=== error ===\n%s\n", buf.Bytes())
	}
	return src, err
}
