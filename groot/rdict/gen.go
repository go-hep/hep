// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"bytes"
	"fmt"
	"go/format"
	"go/importer"
	"go/types"
	"log"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rmeta"
)

var (
	binMa *types.Interface // encoding.BinaryMarshaler
	binUn *types.Interface // encoding.BinaryUnmarshaler

	rootVers *types.Interface // rbytes.RVersioner

	gosizes types.Sizes
)

// Generator holds the state of the ROOT streaemer generation.
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
		buf: new(bytes.Buffer),
		pkg: pkg,
		imps: map[string]int{
			"go-hep.org/x/hep/groot/rbase":  1,
			"go-hep.org/x/hep/groot/rbytes": 1,
			"go-hep.org/x/hep/groot/rdict":  1,
			"go-hep.org/x/hep/groot/rmeta":  1,
		},
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
		log.Printf("typ: %q: %+v\n", typeName, typ)
	}

	if !types.Implements(tn.Type(), rootVers) && !types.Implements(types.NewPointer(tn.Type()), rootVers) {
		log.Fatalf("type %q does not implement %q.", tn.Pkg().Path()+"."+tn.Name(), "go-hep.org/x/hep/groot/rbytes.RVersioner")
	}

	g.genStreamer(typ, typeName)
	// g.genMarshal(typ, typeName)
	// g.genUnmarshal(typ, typeName)
}

func (g *Generator) genMarshal(t types.Type, typeName string) {
	g.printf(`// MarshalROOT implements rbytes.Marshaler
func (o *%[1]s) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	ws, err := w.WStreamer(o)
	if err != nil {
		return 0, err
	}
	return ws.WStream(w)
}
`,
		typeName,
	)
}

func (g *Generator) genUnmarshal(t types.Type, typeName string) {
	g.printf(`// UnmarshalROOT implements rbytes.Unmarshaler
func (o *%[1]s) UnmarshalROOT(r *rbytes.RBuffer) error {
	rs, err := r.RStreamer(o)
	if err != nil {
		return err
	}
	return rs.RStream(r)
}
`,
		typeName,
	)
}

func (g *Generator) genStreamer(t types.Type, typeName string) {
	g.printf(`func init() {
	// Streamer for %[1]s.
	rdict.Streamers.Add(rdict.NewStreamerInfo(%[2]q, []rbytes.StreamerElement{
`,
		typeName,
		g.pkg.Path()+"."+typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		n := ft.Name()
		if tag := typ.Tag(i); tag != "" {
			nn := reflect.StructTag(tag).Get("groot")
			if nn != "" {
				n = nn
			}
		}
		g.genStreamerType(ft.Type(), n)
	}

	g.printf("}))\n}\n\n")
}

func (g *Generator) genStreamerType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {
		case types.Bool:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Uint8:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Uint16:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Uint32, types.Uint:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Uint64:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Int8:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Int16:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Int32, types.Int:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Int64:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Float32:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Float64:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 1))
		case types.Complex64:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.Complex128:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.String:
			g.printf("%s,\n", g.se(ut, n, "", 1))

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Array:
		// FIXME(sbinet): collect+visit element type.
		g.printf(
			"&rdict.StreamerBasicType{StreamerElement: %s},\n",
			g.se(ut.Elem(), n, "+ rmeta.OffsetL", ut.Len()),
		)
	case *types.Slice:
		// FIXME(sbinet): collect+visit element type.
		g.printf("rdict.NewStreamerSTL(%q, rmeta.STLvector, %d),\n", n, gotype2RMeta(ut.Elem()))

	case *types.Struct:
		g.printf(
			"&rdict.StreamerObjectAny{StreamerElement:rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Any,\nSize: %[4]d,\nEName:%[3]q,\n}.New()},\n",
			n, "",
			t.String(), gosizes.Sizeof(ut),
		)

	default:
		log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
	}
}

func (g *Generator) se(t types.Type, n, rtype string, len int64) string {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {
		case types.Bool:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Bool %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				1*len,
			)
		case types.Uint8:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint8 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				1*len,
			)
		case types.Uint16:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint16 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				2*len,
			)
		case types.Uint32, types.Uint:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint32 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				4*len,
			)
		case types.Uint64:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint64 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				8*len,
			)
		case types.Int8:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int8 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				1*len,
			)
		case types.Int16:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int16 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				2*len,
			)
		case types.Int32, types.Int:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int32 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				4*len,
			)
		case types.Int64:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int64 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				8*len,
			)
		case types.Float32:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Float32 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				4*len,
			)
		case types.Float64:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Float64 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				8*len,
			)
		case types.Complex64:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.Complex128:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.String:
			return fmt.Sprintf("&rdict.StreamerString{rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.TString %[4]s,\nSize: 24,\nEName:%[3]q,\n}.New()}",
				n, "",
				"TString",
				rtype,
			)
		}
	}

	return ""
}

func (g *Generator) Format() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf(`// DO NOT EDIT; automatically generated by %[1]s

package %[2]s

import (
`,
		"root-gen-streamer",
		g.pkg.Name(),
	))

	// FIXME(sbinet): separate stdlib from 3rd-party imports.

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

func init() {
	pkg, err := importer.Default().Import("encoding")
	if err != nil {
		log.Fatalf("error finding package \"encoding\": %v\n", err)
	}

	o := pkg.Scope().Lookup("BinaryMarshaler")
	if o == nil {
		log.Fatalf("could not find interface encoding.BinaryMarshaler\n")
	}
	binMa = o.(*types.TypeName).Type().Underlying().(*types.Interface)

	o = pkg.Scope().Lookup("BinaryUnmarshaler")
	if o == nil {
		log.Fatalf("could not find interface encoding.BinaryUnmarshaler\n")
	}
	binUn = o.(*types.TypeName).Type().Underlying().(*types.Interface)

	pkg, err = importer.Default().Import("go-hep.org/x/hep/groot/rbytes")
	if err != nil {
		log.Fatalf("could not find package %q: %v", "go-hep.org/x/hep/groot/rbytes", err)
	}

	o = pkg.Scope().Lookup("RVersioner")
	if o == nil {
		log.Fatalf("could not find interface rbytes.RVersioner")
	}
	rootVers = o.(*types.TypeName).Type().Underlying().(*types.Interface)

	sz := int64(reflect.TypeOf(int(0)).Size())
	gosizes = &types.StdSizes{WordSize: sz, MaxAlign: sz}
}

func gotype2RMeta(t types.Type) rmeta.Enum {
	switch ut := t.Underlying().(type) {
	case *types.Basic:
		switch ut.Kind() {
		case types.Bool:
			return rmeta.Bool
		case types.Uint8:
			return rmeta.Uint8
		case types.Uint16:
			return rmeta.Uint16
		case types.Uint32, types.Uint:
			return rmeta.Uint32
		case types.Uint64:
			return rmeta.Uint64
		case types.Int8:
			return rmeta.Int8
		case types.Int16:
			return rmeta.Int16
		case types.Int32, types.Int:
			return rmeta.Int32
		case types.Int64:
			return rmeta.Int64
		case types.Float32:
			return rmeta.Float32
		case types.Float64:
			return rmeta.Float64
		case types.String:
			return rmeta.TString
		}
	case *types.Struct:
		return rmeta.Any
	case *types.Slice:
		return rmeta.STL
	case *types.Array:
		return rmeta.OffsetL + gotype2RMeta(ut.Elem())
	}
	return -1
}

// GoName2Cxx translates a fully-qualified Go type name to a C++ one.
// e.g.:
//  - go-hep.org/x/hep/hbook.H1D -> go_hep_org::x::hep::hbook::H1D
func GoName2Cxx(name string) string {
	repl := strings.NewReplacer(
		"-", "_",
		"/", "::",
		".", "_",
	)
	i := strings.LastIndex(name, ".")
	if i > 0 {
		name = name[:i] + "::" + name[i+1:]
	}
	return repl.Replace(name)
}

// Typename returns a language dependant typename, usually encoded inside a
// StreamerInfo's title.
func Typename(name, title string) (string, bool) {
	if title == "" {
		return name, false
	}
	i := strings.Index(title, ";")
	if i <= 0 {
		return name, false
	}
	title = strings.TrimSpace(title[i+1:])
	return title, true
}
