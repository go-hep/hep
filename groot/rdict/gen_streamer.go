// Copyright 2019 The go-hep Authors. All rights reserved.
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

	"go-hep.org/x/hep/groot/rmeta"
)

// genStreamer holds the state of the ROOT streamer generation.
type genStreamer struct {
	buf *bytes.Buffer
	pkg *types.Package

	// set of imported packages.
	// usually: "encoding/binary", "math"
	imps map[string]int

	verbose bool // enable verbose mode

	binMa *types.Interface // encoding.BinaryMarshaler
	binUn *types.Interface // encoding.BinaryUnmarshaler
	rvers *types.Interface // rbytes.RVersioner

	gosizes types.Sizes
}

// NewGenStreamer returns a new code generator for package p,
// where p is the package's import path.
func NewGenStreamer(p string, verbose bool) (Generator, error) {
	pkg, err := importer.Default().Import(p)
	if err != nil {
		return nil, err
	}

	g := &genStreamer{
		buf: new(bytes.Buffer),
		pkg: pkg,
		imps: map[string]int{
			"go-hep.org/x/hep/groot/rbytes": 1,
			"go-hep.org/x/hep/groot/rdict":  1,
			"go-hep.org/x/hep/groot/rmeta":  1,
		},
		verbose: verbose,
	}

	err = g.init()
	if err != nil {
		return nil, err
	}

	return g, nil
}

func (g *genStreamer) init() error {
	pkg, err := importer.Default().Import("encoding")
	if err != nil {
		return fmt.Errorf("rdict: could not find package \"encoding\": %w", err)
	}

	o := pkg.Scope().Lookup("BinaryMarshaler")
	if o == nil {
		return fmt.Errorf("rdict: could not find interface encoding.BinaryMarshaler")
	}
	g.binMa = o.(*types.TypeName).Type().Underlying().(*types.Interface)

	o = pkg.Scope().Lookup("BinaryUnmarshaler")
	if o == nil {
		return fmt.Errorf("rdict: could not find interface encoding.BinaryUnmarshaler")
	}
	g.binUn = o.(*types.TypeName).Type().Underlying().(*types.Interface)

	pkg, err = importer.Default().Import("go-hep.org/x/hep/groot/rbytes")
	if err != nil {
		return fmt.Errorf("rdict: could not find package %q: %w", "go-hep.org/x/hep/groot/rbytes", err)
	}

	o = pkg.Scope().Lookup("RVersioner")
	if o == nil {
		return fmt.Errorf("rdict: could not find interface rbytes.RVersioner")
	}
	g.rvers = o.(*types.TypeName).Type().Underlying().(*types.Interface)

	sz := int64(reflect.TypeOf(int(0)).Size())
	g.gosizes = &types.StdSizes{WordSize: sz, MaxAlign: sz}
	return nil
}

func (g *genStreamer) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

// Generate implements rdict.Generator
func (g *genStreamer) Generate(typeName string) error {
	scope := g.pkg.Scope()
	obj := scope.Lookup(typeName)
	if obj == nil {
		return fmt.Errorf("no such type %q in package %q\n", typeName, g.pkg.Path()+"/"+g.pkg.Name())
	}

	tn, ok := obj.(*types.TypeName)
	if !ok {
		return fmt.Errorf("%q is not a type (%v)\n", typeName, obj)
	}

	typ, ok := tn.Type().Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("%q is not a named struct (%v)\n", typeName, tn)
	}
	if g.verbose {
		log.Printf("typ: %q: %+v\n", typeName, typ)
	}

	if !types.Implements(tn.Type(), g.rvers) && !types.Implements(types.NewPointer(tn.Type()), g.rvers) {
		return fmt.Errorf("type %q does not implement %q.", tn.Pkg().Path()+"."+tn.Name(), "go-hep.org/x/hep/groot/rbytes.RVersioner")
	}

	g.genStreamer(typ, typeName)
	g.genMarshal(typ, typeName)
	// g.genUnmarshal(typ, typeName)

	return nil
}

func (g *genStreamer) genMarshal(t types.Type, typeName string) {
	g.printf(`// MarshalROOT implements rbytes.Marshaler
func (o *%[1]s) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

`,
		typeName,
	)

	typ := t.Underlying().(*types.Struct)
	for i := 0; i < typ.NumFields(); i++ {
		ft := typ.Field(i)
		n := ft.Name() // no `groot:"foo"` redirection.
		g.genMarshalType(ft.Type(), n)
	}

	g.printf("\n\treturn w.SetByteCount(pos, o.Class())\n}\n\n")
}

func (g *genStreamer) genUnmarshal(t types.Type, typeName string) {
	g.printf(`// UnmarshalROOT implements rbytes.Unmarshaler
func (o *%[1]s) UnmarshalROOT(r *rbytes.RBuffer) error {
	rs, err := rdict.RStreamer(r, o)
	if err != nil {
		return err
	}
	return rs.RStreamROOT(r)
}
`,
		typeName,
	)
}

func (g *genStreamer) genStreamer(t types.Type, typeName string) {
	g.printf(`func init() {
	// Streamer for %[1]s.
	rdict.StreamerInfos.Add(rdict.NewStreamerInfo(%[2]q, int(((*%[1]s)(nil)).RVersion()), []rbytes.StreamerElement{
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

func (g *genStreamer) genStreamerType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {
		case types.Bool:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Uint8:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Uint16:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Uint32, types.Uint:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Uint64:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Int8:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Int16:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Int32, types.Int:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Int64:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Float32:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Float64:
			g.printf("&rdict.StreamerBasicType{StreamerElement: %s},\n", g.se(ut, n, "", 0))
		case types.Complex64:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.Complex128:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.String:
			g.printf("%s,\n", g.se(ut, n, "", 0))

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Array:
		// FIXME(sbinet): collect+visit element type.
		switch eut := ut.Elem().Underlying().(type) {
		case *types.Basic:
			switch kind := eut.Kind(); kind {
			default:
				g.printf(
					"&rdict.StreamerBasicType{StreamerElement: %s},\n",
					g.se(ut.Elem(), n, "+ rmeta.OffsetL", ut.Len()),
				)
			case types.String:
				g.printf(
					"%s,\n",
					g.se(ut.Elem(), n, "+ rmeta.OffsetL", ut.Len()),
				)
			}
		default:
			g.printf(
				"%s\n",
				g.se(ut.Elem(), n, "+ rmeta.OffsetL", ut.Len()),
			)
		}
	case *types.Slice:
		// FIXME(sbinet): collect+visit element type.
		g.printf("rdict.NewStreamerSTL(%q, rmeta.STLvector, rmeta.%v),\n", n, gotype2RMeta(ut.Elem()))

	case *types.Struct:
		g.imps["go-hep.org/x/hep/groot/rbase"]++
		g.printf(
			"&rdict.StreamerObjectAny{StreamerElement:rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Any,\nSize: %[4]d,\nEName:rdict.GoName2Cxx(%[3]q),\n}.New()},\n",
			n, "",
			t.String(), g.gosizes.Sizeof(ut),
		)

	default:
		log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
	}
}

func (g *genStreamer) wt(t types.Type, n, meth, arr string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {
		case types.Bool:
			g.printf("w.Write%sBool(o.%s%s)\n", meth, n, arr)
		case types.Uint8:
			g.printf("w.Write%sU8(o.%s%s)\n", meth, n, arr)
		case types.Uint16:
			g.printf("w.Write%sU16(o.%s%s)\n", meth, n, arr)
		case types.Uint32:
			g.printf("w.Write%sU32(o.%s%s)\n", meth, n, arr)
		case types.Uint64:
			g.printf("w.Write%sU64(o.%s%s)\n", meth, n, arr)
		case types.Int8:
			g.printf("w.Write%sI8(o.%s%s)\n", meth, n, arr)
		case types.Int16:
			g.printf("w.Write%sI16(o.%s%s)\n", meth, n, arr)
		case types.Int32:
			g.printf("w.Write%sI32(o.%s%s)\n", meth, n, arr)
		case types.Int64:
			g.printf("w.Write%sI64(o.%s%s)\n", meth, n, arr)
		case types.Float32:
			g.printf("w.Write%sF32(o.%s%s)\n", meth, n, arr)
		case types.Float64:
			g.printf("w.Write%sF64(o.%s%s)\n", meth, n, arr)

		case types.Uint:
			g.printf("w.Write%sU64(uint64(o.%s%s))\n", meth, n, arr)
		case types.Int:
			g.printf("w.Write%sI64(int64(o.%s%s))\n", meth, n, arr)

		case types.Complex64:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)
		case types.Complex128:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.String:
			g.printf("w.Write%sString(o.%s%s)\n", meth, n, arr)

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Struct:
		g.printf("o.%s.MarshalROOT(w)\n", n)

	default:
		log.Fatalf("unhandled marshal type: %v (underlying %v)", t, ut)
	}
}

func (g *genStreamer) se(t types.Type, n, rtype string, arrlen int64) string {
	elmt := Element{
		Size: 1,
	}
	if arrlen > 0 {
		elmt.Size = int32(arrlen)
		elmt.ArrLen = int32(arrlen)
		elmt.ArrDim = 1
	}

	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		g.imps["go-hep.org/x/hep/groot/rbase"]++
		switch kind := ut.Kind(); kind {
		case types.Bool:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Bool %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				1*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Uint8:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint8 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				1*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Uint16:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint16 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				2*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Uint32, types.Uint:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint32 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				4*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Uint64:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Uint64 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				8*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Int8:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int8 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				1*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Int16:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int16 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				2*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Int32, types.Int:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int32 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				4*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Int64:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Int64 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				8*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Float32:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Float32 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				4*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Float64:
			return fmt.Sprintf("rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.Float64 %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()",
				n, "",
				rmeta.GoType2Cxx[ut.Name()],
				rtype,
				8*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		case types.Complex64:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.Complex128:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.String:
			return fmt.Sprintf("&rdict.StreamerString{rdict.Element{\nName: *rbase.NewNamed(%[1]q, %[2]q),\nType: rmeta.TString %[4]s,\nSize: %[5]d,\nEName:%[3]q,\nArrLen:%[6]d,\nArrDim:%[7]d,\n}.New()}",
				n, "",
				"TString",
				rtype,
				24*elmt.Size,
				elmt.ArrLen,
				elmt.ArrDim,
			)
		}
	case *types.Struct:
		// FIXME(sbinet): implement.
		// FIXME(sbinet): prevent recursion.
		old := g.buf
		g.buf = new(bytes.Buffer)
		g.genStreamerType(t, n)
		str := g.buf.String()
		g.buf = old
		return str
	}

	log.Printf("gen-streamer: unhandled type: %v (underlying %v)", t, ut)
	return ""
}

func (g *genStreamer) genMarshalType(t types.Type, n string) {
	ut := t.Underlying()
	switch ut := ut.(type) {
	case *types.Basic:
		switch kind := ut.Kind(); kind {
		case types.Bool:
			g.wt(ut, n, "", "")
		case types.Uint8:
			g.wt(ut, n, "", "")
		case types.Uint16:
			g.wt(ut, n, "", "")
		case types.Uint32:
			g.wt(ut, n, "", "")
		case types.Uint64:
			g.wt(ut, n, "", "")
		case types.Int8:
			g.wt(ut, n, "", "")
		case types.Int16:
			g.wt(ut, n, "", "")
		case types.Int32:
			g.wt(ut, n, "", "")
		case types.Int64:
			g.wt(ut, n, "", "")
		case types.Float32:
			g.wt(ut, n, "", "")
		case types.Float64:
			g.wt(ut, n, "", "")

		case types.Uint:
			g.wt(ut, n, "", "")
		case types.Int:
			g.wt(ut, n, "", "")

		case types.Complex64:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)
		case types.Complex128:
			log.Fatalf("unhandled type: %v (underlying %v)\n", t, ut) // FIXME(sbinet)

		case types.String:
			g.wt(ut, n, "", "")

		default:
			log.Fatalf("unhandled type: %v (underlying: %v)\n", t, ut)
		}

	case *types.Array:
		switch ut.Elem().Underlying().(type) {
		case *types.Basic:
			g.wt(ut.Elem(), n, "FastArray", "[:]")
		default:
			g.printf("for i := range o.%s {\n", n)
			g.wt(ut.Elem(), n+"[i]", "", "")
			g.printf("}\n")
		}

	case *types.Slice:
		g.wt(ut.Elem(), n, "FastArray", "")

	case *types.Struct:
		g.printf("o.%s.MarshalROOT(w)\n", n)

	default:
		log.Fatalf("gen-marshal-type: unhandled type: %v (underlying: %v)\n", t, ut)
	}
}

// Generate implements rdict.Generator
func (g *genStreamer) Format() ([]byte, error) {
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

var (
	_ Generator = (*genStreamer)(nil)
)
