// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/rtypes"
)

type genGoType struct {
	pkg string
	buf *bytes.Buffer
	ctx rbytes.StreamerInfoContext

	verbose bool

	// set of imported packages.
	// usually: "go-hep.org/x/hep/groot/rbase", ".../rcont
	imps map[string]int
}

// NewGenGoType generates code for Go types from a ROOT StreamerInfo.
func NewGenGoType(pkg string, sictx rbytes.StreamerInfoContext, verbose bool) (*genGoType, error) {
	return &genGoType{
		pkg:     pkg,
		buf:     new(bytes.Buffer),
		ctx:     sictx,
		verbose: verbose,
		imps: map[string]int{
			"reflect":                       1,
			"go-hep.org/x/hep/groot/rbytes": 1,
			"go-hep.org/x/hep/groot/root":   1,
			"go-hep.org/x/hep/groot/rtypes": 1,
		},
	}, nil
}

// Generate implements rdict.Generator
func (g *genGoType) Generate(name string) error {
	if g.verbose {
		log.Printf("generating type for %q...", name)
	}
	si, err := g.ctx.StreamerInfo(name, -1)
	if err != nil {
		return errors.Wrapf(err, "rdict: could not find streamer for %q", name)
	}

	return g.genType(si)
}

func (g *genGoType) genType(si rbytes.StreamerInfo) error {
	name := si.Name()
	if title := si.Title(); title != "" {
		g.printf("// %s has been automatically generated.\n// %s\n", name, title)
	}
	g.printf("type %s struct{\n", name)
	for i, se := range si.Elements() {
		g.genField(si, i, se)
	}
	g.printf("}\n\n")
	g.printf("func (*%s) Class() string { return %q }\n", name, name)
	g.printf("func (*%s) RVersion() int16 { return %d }\n", name, si.ClassVersion())
	g.printf("\n")
	g.genMarshal(si)
	g.genUnmarshal(si)

	g.printf(`func init() {
	f := func() reflect.Value {
		var o %s
		return reflect.ValueOf(&o)
	}
	rtypes.Factory.Add(%q, f)
}

`,
		name, name,
	)

	g.genStreamerInfo(si)

	ifaces := []string{"root.Object", "rbytes.RVersioner", "rbytes.Marshaler", "rbytes.Unmarshaler"}
	g.printf("var (\n")
	for _, n := range ifaces {
		g.printf("\t_ %s = (*%s)(nil)\n", n, name)
	}
	g.printf(")\n\n")

	return nil
}

func (g *genGoType) genField(si rbytes.StreamerInfo, i int, se rbytes.StreamerElement) {
	const (
		docFmt = "\t%s\t%s\t%s\t%s\n"
	)
	doc := g.doc(se.Title())
	if doc != "" {
		doc = "// " + doc
	}
	switch se := se.(type) {
	case *StreamerBase:
		g.printf(docFmt, fmt.Sprintf("base%d", i), g.typename(se), g.stag(i, se), "// base class")

	case *StreamerBasicPointer:
		g.printf(docFmt, se.Name(), g.typename(se), g.stag(i, se), doc)

	case *StreamerBasicType:
		tname := g.typename(se)
		switch se.ArrayLen() {
		case 0:
		default:
			tname = fmt.Sprintf("[%d]%s", se.ArrayLen(), tname)
		}
		g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)

	case *StreamerLoop:
		tname := g.typename(se)
		g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)

	case *StreamerObject:
		tname := g.typename(se)
		switch se.ArrayLen() {
		case 0:
		default:
			tname = fmt.Sprintf("[%d]%s", se.ArrayLen(), tname)
		}
		g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)

	case *StreamerObjectAny:
		tname := g.typename(se)
		switch se.ArrayLen() {
		case 0:
		default:
			tname = fmt.Sprintf("[%d]%s", se.ArrayLen(), tname)
		}
		g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)

	case *StreamerObjectAnyPointer:
		tname := g.typename(se)
		g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)

	case *StreamerObjectPointer:
		tname := g.typename(se)
		g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)

	case *StreamerString, *StreamerSTLstring:
		g.printf(docFmt, se.Name(), "string", g.stag(i, se), doc)

	case *StreamerSTL:
		switch se.STLVectorType() {
		case rmeta.STLvector:
			tname := g.typename(se)
			g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)
		default:
			panic(errors.Errorf("STL-type not implemented %#v", se))
		}
	default:
		g.printf("\t%s\t%s // %T -- %s\n", se.Name(), g.typename(se), se, doc)
	}
}

func (g *genGoType) stag(i int, se rbytes.StreamerElement) string {
	switch se := se.(type) {
	case *StreamerBase:
		return fmt.Sprintf("`groot:\"BASE-%s\"`", se.Name())
	}
	meta, _ := g.rcomment(se.Title())
	if meta != "" {
		return fmt.Sprintf("`groot:\"%s,meta=%s\"`", se.Name(), meta)
	}
	return fmt.Sprintf("`groot:%q`", se.Name())
}

func (g *genGoType) doc(title string) string {
	_, doc := g.rcomment(title)
	return doc
}

func (g *genGoType) rcomment(s string) (meta, comment string) {
	comment = s
	for strings.HasPrefix(comment, "[") {
		beg := strings.Index(comment, "[")
		end := strings.Index(comment, "]")
		meta += comment[beg : end+1]
		comment = comment[end+1:]
	}
	if strings.HasPrefix(comment, " ") {
		comment = strings.TrimSpace(comment)
	}
	return meta, comment
}

func (g *genGoType) typename(se rbytes.StreamerElement) string {
	tname := se.TypeName()
	switch se := se.(type) {
	case *StreamerBase:
		return g.cxx2go(se.Name(), 0)

	case *StreamerBasicPointer:
		tname = tname[:len(tname)-1] // drop last '*'
		t, ok := rmeta.CxxBuiltins[tname]
		if !ok {
			panic(errors.Errorf("gen-type: unknown C++ builtin %q", tname))
		}
		switch {
		case strings.HasPrefix(se.Title(), "["):
			return g.cxx2go(t.Name(), qualSlice)
		default:
			return g.cxx2go(t.Name(), qualStar)
		}

	case *StreamerBasicType:
		switch se.Type() {
		case rmeta.Float16:
			g.imps["go-hep.org/x/hep/groot/root"] = 1
			return "root.Float16"
		case rmeta.Double32:
			g.imps["go-hep.org/x/hep/groot/root"] = 1
			return "root.Double32"
		}
		t, ok := rmeta.CxxBuiltins[tname]
		if !ok {
			panic(errors.Errorf("gen-type: unknown C++ builtin %q", tname))
		}
		return t.Name()

	case *StreamerLoop:
		if strings.HasSuffix(tname, "*") {
			tname = tname[:len(tname)-1]
		}
		return "[]" + g.cxx2go(tname, qualNone)

	case *StreamerObject:
		return g.cxx2go(tname, qualNone)

	case *StreamerObjectAny:
		return g.cxx2go(tname, qualNone)

	case *StreamerObjectAnyPointer:
		tname = tname[:len(tname)-1] // drop last '*'
		return g.cxx2go(tname, qualStar)

	case *StreamerObjectPointer:
		tname = tname[:len(tname)-1] // drop last '*'
		return g.cxx2go(tname, qualStar)

	case *StreamerString, *StreamerSTLstring:
		return "string"

	case *StreamerSTL:
		switch se.STLVectorType() {
		case rmeta.STLvector:
			switch se.ContainedType() {
			case rmeta.Bool:
				return "[]bool"
			case rmeta.Int8:
				return "[]int8"
			case rmeta.Int16:
				return "[]int16"
			case rmeta.Int32:
				return "[]int32"
			case rmeta.Int64:
				return "[]int64"
			case rmeta.Uint8:
				return "[]uint8"
			case rmeta.Uint16:
				return "[]uint16"
			case rmeta.Uint32:
				return "[]uint32"
			case rmeta.Uint64:
				return "[]uint64"
			case rmeta.Float32, rmeta.Float16:
				return "[]float32"
			case rmeta.Float64, rmeta.Double32:
				return "[]float64"
			case rmeta.Object:
				switch se.ElemTypeName() {
				case "string":
					return "[]string"
				default:
					etname := g.cxx2go(se.ElemTypeName(), qualNone)
					return "[]" + etname
				}
			default:
				panic(errors.Errorf("invalid stl-vector element type: %v -- %#v", se.ContainedType(), se))
			}
		default:
			panic(errors.Errorf("STL-type not implemented %#v", se))
		}
	}
	return tname
}

type qualKind uint8

func (q qualKind) String() string {
	switch q {
	case qualNone:
		return ""
	case qualStar:
		return "*"
	case qualSlice:
		return "[]"
	}
	panic(fmt.Errorf("invalid qual-kind value %d", int(q)))
}

const (
	qualNone  qualKind = 0
	qualStar  qualKind = 1
	qualSlice qualKind = 2
)

func (g *genGoType) cxx2go(name string, qual qualKind) string {
	prefix := qual.String()
	if rtypes.Factory.HasKey(name) {
		t := rtypes.Factory.Get(name)().Type().Elem()
		pkg := t.PkgPath()
		g.imps[pkg]++
		switch t.Name() {
		case "tgraph":
			if qual == qualStar {
				prefix = ""
			}
			return prefix + "rhist.Graph"
		case "tgrapherrs", "tgraphasymerrs":
			if qual == qualStar {
				prefix = ""
			}
			return prefix + "rhist.GraphErrors"
		}

		return prefix + filepath.Base(pkg) + "." + t.Name()
	}
	var f func(name string) string
	f = func(name string) string {
		switch {
		case strings.HasSuffix(name, "*"):
			return "*" + f(name[:len(name)-1])
		}
		return name
	}
	return prefix + f(name)
}

func (g *genGoType) genMarshal(si rbytes.StreamerInfo) {
	g.printf(`// MarshalROOT implements rbytes.Marshaler
func (o *%[1]s) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

`,
		si.Name(),
	)

	for i, se := range si.Elements() {
		g.genMarshalField(si, i, se)
	}

	g.printf("\n\treturn w.SetByteCount(pos, o.Class())\n}\n\n")
}

func (g *genGoType) genMarshalField(si rbytes.StreamerInfo, i int, se rbytes.StreamerElement) {
	switch se := se.(type) {
	case *StreamerBase:
		g.printf("o.%s.MarshalROOT(w)\n", fmt.Sprintf("base%d", i))

	case *StreamerBasicPointer:
		title := se.Title()
		switch {
		case strings.HasPrefix(title, "["):
			n := title[strings.Index(title, "[")+1 : strings.Index(title, "]")]
			g.printf("w.WriteI8(1) // is-array\n")
			wfunc := ""
			switch se.Type() {
			case rmeta.OffsetP + rmeta.Bool:
				wfunc = "WriteFastArrayBool"
			case rmeta.OffsetP + rmeta.Int8:
				wfunc = "WriteFastArrayI8"
			case rmeta.OffsetP + rmeta.Int16:
				wfunc = "WriteFastArrayI16"
			case rmeta.OffsetP + rmeta.Int32:
				wfunc = "WriteFastArrayI32"
			case rmeta.OffsetP + rmeta.Int64:
				wfunc = "WriteFastArrayI64"
			case rmeta.OffsetP + rmeta.Uint8:
				wfunc = "WriteFastArrayU8"
			case rmeta.OffsetP + rmeta.Uint16:
				wfunc = "WriteFastArrayU16"
			case rmeta.OffsetP + rmeta.Uint32:
				wfunc = "WriteFastArrayU32"
			case rmeta.OffsetP + rmeta.Uint64:
				wfunc = "WriteFastArrayU64"
			case rmeta.OffsetP + rmeta.Float32:
				wfunc = "WriteFastArrayF32"
			case rmeta.OffsetP + rmeta.Float64:
				wfunc = "WriteFastArrayF64"
			default:
				panic(errors.Errorf("invalid element type: %v", se.Type()))
			}
			g.printf("w.%s(o.%s[:o.%s])\n", wfunc, se.Name(), n)
		default:
			panic("not implemented")
		}

	case *StreamerBasicType:
		switch se.ArrayLen() {
		case 0:
			switch se.Type() {
			case rmeta.Bool:
				g.printf("w.WriteBool(o.%s)\n", se.Name())

			case rmeta.Counter:
				switch se.Size() {
				case 4:
					g.printf("w.WriteI32(int32(o.%s))\n", se.Name())
				case 8:
					g.printf("w.WriteI64(int64(o.%s))\n", se.Name())
				default:
					panic(errors.Errorf("invalid counter size %d for %s.%s", se.Size(), si.Name(), se.Name()))
				}

			case rmeta.Bits:
				g.printf("w.WriteI32(int32(o.%s))\n", se.Name())

			case rmeta.Int8:
				g.printf("w.WriteI8(o.%s)\n", se.Name())
			case rmeta.Int16:
				g.printf("w.WriteI16(o.%s)\n", se.Name())
			case rmeta.Int32:
				g.printf("w.WriteI32(o.%s)\n", se.Name())
			case rmeta.Int64:
				g.printf("w.WriteI64(o.%s)\n", se.Name())

			case rmeta.Uint8:
				g.printf("w.WriteU8(o.%s)\n", se.Name())
			case rmeta.Uint16:
				g.printf("w.WriteU16(o.%s)\n", se.Name())
			case rmeta.Uint32:
				g.printf("w.WriteU32(o.%s)\n", se.Name())
			case rmeta.Uint64:
				g.printf("w.WriteU64(o.%s)\n", se.Name())

			case rmeta.Float32:
				g.printf("w.WriteF32(o.%s)\n", se.Name())
			case rmeta.Float64:
				g.printf("w.WriteF64(o.%s)\n", se.Name())

			case rmeta.Float16:
				g.printf("w.WriteF32(float32(o.%s)) // FIXME(sbinet)\n", se.Name()) // FIXME(sbinet): handle compression
			case rmeta.Double32:
				g.printf("w.WriteF32(float32(o.%s)) // FIXME(sbinet)\n", se.Name()) // FIXME(sbinet): handle compression

			default:
				panic(errors.Errorf("invalid basic type %v (%d) for %s.%s", se.Type(), se.Type(), si.Name(), se.Name()))
			}
		default:
			n := int(se.ArrayLen())
			wfunc := ""
			switch se.Type() {
			case rmeta.OffsetL + rmeta.Bool:
				wfunc = "WriteFastArrayBool"
			case rmeta.OffsetL + rmeta.Int8:
				wfunc = "WriteFastArrayI8"
			case rmeta.OffsetL + rmeta.Int16:
				wfunc = "WriteFastArrayI16"
			case rmeta.OffsetL + rmeta.Int32:
				wfunc = "WriteFastArrayI32"
			case rmeta.OffsetL + rmeta.Int64:
				wfunc = "WriteFastArrayI64"
			case rmeta.OffsetL + rmeta.Uint8:
				wfunc = "WriteFastArrayU8"
			case rmeta.OffsetL + rmeta.Uint16:
				wfunc = "WriteFastArrayU16"
			case rmeta.OffsetL + rmeta.Uint32:
				wfunc = "WriteFastArrayU32"
			case rmeta.OffsetL + rmeta.Uint64:
				wfunc = "WriteFastArrayU64"
			case rmeta.OffsetL + rmeta.Float32:
				wfunc = "WriteFastArrayF32"
			case rmeta.OffsetL + rmeta.Float64:
				wfunc = "WriteFastArrayF64"
			default:
				panic(errors.Errorf("invalid array element type: %v", se.Type()))
			}
			g.printf("w.%s(o.%s[:%d])\n", wfunc, se.Name(), n)
		}

	case *StreamerLoop:
		// FIXME(sbinet): implement. handle mbr-wise
		g.printf("panic(\"o.%s: not implemented (TStreamerLoop)\")\n", se.Name())

	case *StreamerObject:
		// FIXME(sbinet): check semantics
		switch se.ArrayLen() {
		case 0:
			g.printf("o.%s.MarshalROOT(w) // obj\n", se.Name())
		default:
			g.printf("for i := range o.%s {\n", se.Name())
			g.printf("o.%s[i].MarshalROOT(w) // obj\n", se.Name())
			g.printf("}\n")
		}

	case *StreamerObjectAny:
		// FIXME(sbinet): check semantics
		switch se.ArrayLen() {
		case 0:
			g.printf("o.%s.MarshalROOT(w) // obj-any\n", se.Name())
		default:
			g.printf("for i := range o.%s {\n", se.Name())
			g.printf("o.%s[i].MarshalROOT(w) // obj-any\n", se.Name())
			g.printf("}\n")
		}

	case *StreamerObjectAnyPointer:
		// FIXME(sbinet): check semantics
		g.printf("w.WriteObjectAny(o.%s) // obj-any-ptr\n", se.Name())

	case *StreamerObjectPointer:
		// FIXME(sbinet): check semantics
		g.printf("w.WriteObjectAny(o.%s) // obj-ptr \n", se.Name())

	case *StreamerString:
		g.printf("w.WriteString(o.%s)\n", se.Name())

	case *StreamerSTLstring:
		g.printf("w.WriteSTLString(o.%s)\n", se.Name())

	case *StreamerSTL:
		switch se.STLVectorType() {
		case rmeta.STLvector:
			wfunc := ""
			switch se.ContainedType() {
			case rmeta.Bool:
				wfunc = "WriteFastArrayBool"
			case rmeta.Int8:
				wfunc = "WriteFastArrayI8"
			case rmeta.Int16:
				wfunc = "WriteFastArrayI16"
			case rmeta.Int32:
				wfunc = "WriteFastArrayI32"
			case rmeta.Int64:
				wfunc = "WriteFastArrayI64"
			case rmeta.Uint8:
				wfunc = "WriteFastArrayU8"
			case rmeta.Uint16:
				wfunc = "WriteFastArrayU16"
			case rmeta.Uint32:
				wfunc = "WriteFastArrayU32"
			case rmeta.Uint64:
				wfunc = "WriteFastArrayU64"
			case rmeta.Float32:
				wfunc = "WriteFastArrayF32"
			case rmeta.Float64:
				wfunc = "WriteFastArrayF64"
			case rmeta.Object:
				switch se.ElemTypeName() {
				case "string":
					wfunc = "WriteFastArrayString"
				default:
					panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
				}
			default:
				panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
			}
			g.imps["github.com/pkg/errors"] = 1
			g.imps["go-hep.org/x/hep/groot/rvers"] = 1
			g.printf("{\n")
			g.printf("pos := w.WriteVersion(rvers.StreamerInfo)\n")
			g.printf("w.WriteI32(int32(len(o.%s)))\n", se.Name())
			g.printf("w.%s(o.%s)\n", wfunc, se.Name())
			g.printf("if _, err := w.SetByteCount(pos, %q); err != nil {\n", se.TypeName())
			g.printf("w.SetErr(err)\n")
			g.printf("return 0, w.Err()\n")
			g.printf("}\n")
			g.printf("}\n")

		default:
			panic(errors.Errorf("STL-type not implemented %#v", se))
		}

	default:
		g.printf("// %s -- %T\n", se.Name(), se)
	}
}

func (g *genGoType) genUnmarshal(si rbytes.StreamerInfo) {
	g.printf(`// UnmarshalROOT implements rbytes.Unmarshaler
func (o *%[1]s) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}
	
	start := r.Pos()
	/*vers*/ _, pos, bcnt := r.ReadVersion(o.Class())

`,
		si.Name(),
	)

	for i, se := range si.Elements() {
		g.genUnmarshalField(si, i, se)
	}

	g.printf("\nr.CheckByteCount(pos, bcnt, start, o.Class())\nreturn r.Err()\n}\n\n")
}

func (g *genGoType) genUnmarshalField(si rbytes.StreamerInfo, i int, se rbytes.StreamerElement) {
	switch se := se.(type) {
	case *StreamerBase:
		g.printf("o.%s.UnmarshalROOT(r)\n", fmt.Sprintf("base%d", i))

	case *StreamerBasicPointer:
		title := se.Title()
		switch {
		case strings.HasPrefix(title, "["):
			n := title[strings.Index(title, "[")+1 : strings.Index(title, "]")]
			g.printf("_ = r.ReadI8() // is-array\n")
			rfunc := ""
			switch se.Type() {
			case rmeta.OffsetP + rmeta.Bool:
				rfunc = "ReadFastArrayBool"
			case rmeta.OffsetP + rmeta.Int8:
				rfunc = "ReadFastArrayI8"
			case rmeta.OffsetP + rmeta.Int16:
				rfunc = "ReadFastArrayI16"
			case rmeta.OffsetP + rmeta.Int32:
				rfunc = "ReadFastArrayI32"
			case rmeta.OffsetP + rmeta.Int64:
				rfunc = "ReadFastArrayI64"
			case rmeta.OffsetP + rmeta.Uint8:
				rfunc = "ReadFastArrayU8"
			case rmeta.OffsetP + rmeta.Uint16:
				rfunc = "ReadFastArrayU16"
			case rmeta.OffsetP + rmeta.Uint32:
				rfunc = "ReadFastArrayU32"
			case rmeta.OffsetP + rmeta.Uint64:
				rfunc = "ReadFastArrayU64"
			case rmeta.OffsetP + rmeta.Float32:
				rfunc = "ReadFastArrayF32"
			case rmeta.OffsetP + rmeta.Float64:
				rfunc = "ReadFastArrayF64"
			default:
				panic(errors.Errorf("invalid element type: %v", se.Type()))
			}
			g.printf("o.%s = r.%s(int(o.%s))\n", se.Name(), rfunc, n)
		default:
			panic("not implemented")
		}

	case *StreamerBasicType:
		switch se.ArrayLen() {
		case 0:
			switch se.Type() {
			case rmeta.Bool:
				g.printf("o.%s = r.ReadBool()\n", se.Name())

			case rmeta.Counter:
				switch se.Size() {
				case 4:
					g.printf("o.%s = r.ReadI32()\n", se.Name())
				case 8:
					g.printf("o.%s = r.ReadI64()\n", se.Name())
				default:
					panic(errors.Errorf("invalid counter size %d for %s.%s", se.Size(), si.Name(), se.Name()))
				}

			case rmeta.Bits:
				g.printf("o.%s = r.ReadI32()\n", se.Name())

			case rmeta.Int8:
				g.printf("o.%s = r.ReadI8()\n", se.Name())
			case rmeta.Int16:
				g.printf("o.%s = r.ReadI16()\n", se.Name())
			case rmeta.Int32:
				g.printf("o.%s = r.ReadI32()\n", se.Name())
			case rmeta.Int64:
				g.printf("o.%s = r.ReadI64()\n", se.Name())

			case rmeta.Uint8:
				g.printf("o.%s = r.ReadU8()\n", se.Name())
			case rmeta.Uint16:
				g.printf("o.%s = r.ReadU16()\n", se.Name())
			case rmeta.Uint32:
				g.printf("o.%s = r.ReadU32()\n", se.Name())
			case rmeta.Uint64:
				g.printf("o.%s = r.ReadU64()\n", se.Name())

			case rmeta.Float32:
				g.printf("o.%s = r.ReadF32()\n", se.Name())
			case rmeta.Float64:
				g.printf("o.%s = r.ReadF64()\n", se.Name())

			case rmeta.Float16:
				g.printf("o.%s = root.Float16(r.ReadF32()) // FIXME(sbinet)\n", se.Name()) // FIXME(sbinet): handle compression,factor
			case rmeta.Double32:
				g.printf("o.%s = root.Double32(r.ReadF32()) // FIXME(sbinet)\n", se.Name()) // FIXME(sbinet): handle compression,factor

			default:
				panic(errors.Errorf("invalid basic type %v (%d) for %s.%s", se.Type(), se.Type(), si.Name(), se.Name()))
			}
		default:
			rfunc := ""
			switch se.Type() {
			case rmeta.OffsetL + rmeta.Bool:
				rfunc = "ReadFastArrayBool"
			case rmeta.OffsetL + rmeta.Int8:
				rfunc = "ReadFastArrayI8"
			case rmeta.OffsetL + rmeta.Int16:
				rfunc = "ReadFastArrayI16"
			case rmeta.OffsetL + rmeta.Int32:
				rfunc = "ReadFastArrayI32"
			case rmeta.OffsetL + rmeta.Int64:
				rfunc = "ReadFastArrayI64"
			case rmeta.OffsetL + rmeta.Uint8:
				rfunc = "ReadFastArrayU8"
			case rmeta.OffsetL + rmeta.Uint16:
				rfunc = "ReadFastArrayU16"
			case rmeta.OffsetL + rmeta.Uint32:
				rfunc = "ReadFastArrayU32"
			case rmeta.OffsetL + rmeta.Uint64:
				rfunc = "ReadFastArrayU64"
			case rmeta.OffsetL + rmeta.Float32:
				rfunc = "ReadFastArrayF32"
			case rmeta.OffsetL + rmeta.Float64:
				rfunc = "ReadFastArrayF64"
			default:
				panic(errors.Errorf("invalid array element type: %v", se.Type()))
			}
			g.printf("copy(o.%s[:], r.%s(len(o.%s)))\n", se.Name(), rfunc, se.Name())
		}

	case *StreamerLoop:
		// FIXME(sbinet): implement. handle mbr-wise
		g.printf("panic(\"o.%s: not implemented (TStreamerLoop)\")\n", se.Name())

	case *StreamerObject:
		// FIXME(sbinet): check semantics
		switch se.ArrayLen() {
		case 0:
			g.printf("o.%s.UnmarshalROOT(r) // obj\n", se.Name())
		default:
			g.printf("for i := range o.%s {\n", se.Name())
			g.printf("o.%s[i].UnmarshalROOT(r) // obj\n", se.Name())
			g.printf("}\n")
		}

	case *StreamerObjectAny:
		// FIXME(sbinet): check semantics
		switch se.ArrayLen() {
		case 0:
			g.printf("o.%s.UnmarshalROOT(r) // obj-any\n", se.Name())
		default:
			g.printf("for i := range o.%s {\n", se.Name())
			g.printf("o.%s[i].UnmarshalROOT(r) // obj-any\n", se.Name())
			g.printf("}\n")
		}

	case *StreamerObjectAnyPointer:
		// FIXME(sbinet): check semantics
		g.printf("{\n")
		g.printf("o.%s = nil\n", se.Name())
		g.printf("if oo := r.ReadObjectAny(); oo != nil {  // obj-any-ptr\n")
		g.printf("o.%s = oo.(%s)\n", se.Name(), g.typename(se))
		g.printf("}\n}\n")

	case *StreamerObjectPointer:
		// FIXME(sbinet): check semantics
		g.printf("{\n")
		g.printf("o.%s = nil\n", se.Name())
		g.printf("if oo := r.ReadObjectAny(); oo != nil {  // obj-ptr\n")
		g.printf("o.%s = oo.(%s)\n", se.Name(), g.typename(se))
		g.printf("}\n}\n")

	case *StreamerString:
		g.printf("o.%s = r.ReadString()\n", se.Name())

	case *StreamerSTLstring:
		g.printf("o.%s = r.ReadSTLString()\n", se.Name())

	case *StreamerSTL:
		switch se.STLVectorType() {
		case rmeta.STLvector:
			rfunc := ""
			switch se.ContainedType() {
			case rmeta.Bool:
				rfunc = "ReadFastArrayBool"
			case rmeta.Int8:
				rfunc = "ReadFastArrayI8"
			case rmeta.Int16:
				rfunc = "ReadFastArrayI16"
			case rmeta.Int32:
				rfunc = "ReadFastArrayI32"
			case rmeta.Int64:
				rfunc = "ReadFastArrayI64"
			case rmeta.Uint8:
				rfunc = "ReadFastArrayU8"
			case rmeta.Uint16:
				rfunc = "ReadFastArrayU16"
			case rmeta.Uint32:
				rfunc = "ReadFastArrayU32"
			case rmeta.Uint64:
				rfunc = "ReadFastArrayU64"
			case rmeta.Float32:
				rfunc = "ReadFastArrayF32"
			case rmeta.Float64:
				rfunc = "ReadFastArrayF64"
			case rmeta.Object:
				switch se.ElemTypeName() {
				case "string":
					rfunc = "ReadFastArrayString"
				default:
					panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
				}
			default:
				panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
			}
			g.imps["github.com/pkg/errors"] = 1
			g.imps["go-hep.org/x/hep/groot/rvers"] = 1
			g.printf("{\n")
			g.printf("vers, pos, bcnt := r.ReadVersion(%q)\n", se.TypeName())
			g.printf("if vers != rvers.StreamerInfo {\n")
			g.printf("r.SetErr(errors.Errorf(\"rbytes: invalid version for \\\"%s\\\". got=%%v, want=%%v\", vers, rvers.StreamerInfo))\n", se.TypeName())
			g.printf("return r.Err()\n")
			g.printf("}\n")
			g.printf("o.%s = r.%s(int(r.ReadI32()))\n", se.Name(), rfunc)
			g.printf("r.CheckByteCount(pos, bcnt, start, %q)\n", se.TypeName())
			g.printf("}\n")
		default:
			panic(errors.Errorf("STL-type not implemented %#v", se))
		}

	default:
		g.printf("// %s -- %T\n", se.Name(), se)
	}
}

func (g *genGoType) genStreamerInfo(si rbytes.StreamerInfo) {
	g.printf(`func init() {
		// Streamer for %[1]s.
		rdict.Streamers.Add(rdict.NewCxxStreamerInfo(%[1]q, %[2]d, 0x%[3]x, []rbytes.StreamerElement{
`,
		si.Name(), si.ClassVersion(), si.CheckSum(),
	)

	for i, se := range si.Elements() {
		g.genStreamerType(si, i, se)
	}

	g.printf("}))\n}\n\n")
}

func (g *genGoType) genStreamerType(si rbytes.StreamerInfo, i int, se rbytes.StreamerElement) {
	maxidx := func(v [5]int32) string {
		return fmt.Sprintf("[5]int32{%d, %d, %d, %d, %d}", v[0], v[1], v[2], v[3], v[4])
	}

	g.imps["go-hep.org/x/hep/groot/rbase"] = 1
	g.imps["go-hep.org/x/hep/groot/rdict"] = 1
	g.imps["go-hep.org/x/hep/groot/rmeta"] = 1

	switch se := se.(type) {
	case *StreamerBase:
		g.printf(`rdict.NewStreamerBase(rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.Base,
			Size:   %[3]d,
			ArrLen: %[4]d,
			ArrDim: %[5]d,
			MaxIdx: %[6]s,
			Offset: %[7]d,
			EName:  %[8]q,
			XMin:   %[9]f,
			XMax:   %[10]f,
			Factor: %[11]f,
		}.New(), %[12]d),
		`,
			se.Name(), se.Title(),
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
			se.vbase,
		)

	case *StreamerBasicType:
		g.printf(`&rdict.StreamerBasicType{StreamerElement: rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New()},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
		)

	case *StreamerBasicPointer:
		g.printf(`rdict.NewStreamerBasicPointer(rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   %[3]d,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New(), %[13]d, %[14]q, %[15]q),
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
			se.cvers, se.cname, se.ccls,
		)

	case *StreamerLoop:
		g.printf(`rdict.NewStreamerLoop(rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.StreamLoop,
			Size:   %[3]d,
			ArrLen: %[4]d,
			ArrDim: %[5]d,
			MaxIdx: %[6]s,
			Offset: %[7]d,
			EName:  %[8]q,
			XMin:   %[9]f,
			XMax:   %[10]f,
			Factor: %[11]f,
		}.New(), %[12]d, %[13]q, %[14]q),
		`,
			se.Name(), se.Title(),
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
			se.cvers, se.cname, se.cclass,
		)

	case *StreamerObject:
		g.printf(`&rdict.StreamerObject{StreamerElement: rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New()},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
		)

	case *StreamerObjectPointer:
		g.printf(`&rdict.StreamerObjectPointer{StreamerElement: rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New()},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
		)

	case *StreamerObjectAny:
		g.printf(`&rdict.StreamerObjectAny{StreamerElement: rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New()},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
		)

	case *StreamerObjectAnyPointer:
		g.printf(`&rdict.StreamerObjectAnyPointer{StreamerElement: rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New()},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
		)

	case *StreamerString:
		g.printf(`&rdict.StreamerString{StreamerElement: rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New()},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
		)

	case *StreamerSTL:
		g.printf(`rdict.NewCxxStreamerSTL(rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New(), %[13]d, %[14]d),
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
			se.vtype, se.ctype,
		)

	case *StreamerSTLstring:
		g.printf(`&rdict.StreamerSTLstring{*rdict.NewCxxStreamerSTL(rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New(), %[13]d, %[14]d)},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
			se.vtype, se.ctype,
		)

	case *StreamerArtificial:
		g.printf(`&rdict.StreamerArtificial{StreamerElement: rdict.Element{
			Name:   *rbase.NewNamed(%[1]q, %[2]q),
			Type:   rmeta.%[3]v,
			Size:   %[4]d,
			ArrLen: %[5]d,
			ArrDim: %[6]d,
			MaxIdx: %[7]s,
			Offset: %[8]d,
			EName:  %[9]q,
			XMin:   %[10]f,
			XMax:   %[11]f,
			Factor: %[12]f,
		}.New()},
		`,
			se.Name(), se.Title(),
			se.etype,
			se.esize,
			se.arrlen,
			se.arrdim,
			maxidx(se.maxidx),
			se.offset,
			se.ename,
			se.xmin, se.xmax, se.factor,
		)

	default:
		g.printf("// %s -- %T\n", se.Name(), se)
		panic(errors.Errorf("rdict: unknown streamer type %T (%#v)", se, se))
	}
}

func (g *genGoType) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

// Generate implements rdict.Generator
func (g *genGoType) Format() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf(`// DO NOT EDIT; automatically generated by %[1]s

package %[2]s

import (
`,
		"root-gen-type",
		g.pkg,
	))

	var (
		stdlib []string
		imps   []string
	)

	for k := range g.imps {
		switch {
		case !strings.Contains(k, "."): // stdlib
			stdlib = append(stdlib, k)
		default:
			imps = append(imps, k)
		}
	}

	sort.Strings(stdlib)
	for _, pkg := range stdlib {
		fmt.Fprintf(buf, "%q\n", pkg)
	}
	fmt.Fprintf(buf, "\n")

	sort.Strings(imps)
	for _, pkg := range imps {
		fmt.Fprintf(buf, "%q\n", pkg)
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
	_ Generator = (*genGoType)(nil)
)
