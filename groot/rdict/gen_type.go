// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"log"
	"path/filepath"
	"sort"
	"strings"

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

	rdict string // whether to prepend 'rdict.'
}

// GenCxxStreamerInfo generates the textual representation of the provided streamer info.
func GenCxxStreamerInfo(w io.Writer, si rbytes.StreamerInfo, verbose bool) error {
	g, err := NewGenGoType("go-hep.org/x/hep/groot/rdict", nil, verbose)
	if err != nil {
		return fmt.Errorf("rdict: could not create streamer info generator: %w", err)
	}
	g.rdict = ""

	g.printf("%sNewCxxStreamerInfo(%q, %d, 0x%x, []rbytes.StreamerElement{\n", g.rdict, si.Name(), si.ClassVersion(), si.CheckSum())
	for i, se := range si.Elements() {
		g.genStreamerType(si, i, se)
	}
	g.printf("})")

	src, err := format.Source(g.buf.Bytes())
	if err != nil {
		return fmt.Errorf("rdict: could not format streamer code for %q: %w", si.Name(), err)
	}

	_, err = w.Write(src)
	if err != nil {
		return fmt.Errorf("rdict: could not write streamer info generated data: %w", err)
	}

	return err
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
		rdict: "rdict.",
	}, nil
}

// Generate implements rdict.Generator
func (g *genGoType) Generate(name string) error {
	if g.verbose {
		log.Printf("generating type for %q...", name)
	}
	si, err := g.ctx.StreamerInfo(name, -1)
	if err != nil {
		return fmt.Errorf("rdict: could not find streamer for %q: %w", name, err)
	}

	return g.genType(si)
}

func (g *genGoType) genType(si rbytes.StreamerInfo) error {
	name := si.Name()
	if title := si.Title(); title != "" {
		g.printf("// %s has been automatically generated.\n", name)
		g.printf("// %s\n", title)
	}
	goname := name
	goname = strings.Replace(goname, "::", "__", -1) // handle namespaces
	goname = strings.Replace(goname, "<", "_", -1)   // handle C++ templates
	goname = strings.Replace(goname, ">", "_", -1)   // handle C++ templates
	goname = strings.Replace(goname, ",", "_", -1)   // handle C++ templates
	g.printf("type %s struct{\n", goname)
	for i, se := range si.Elements() {
		g.genField(si, i, se)
	}
	g.printf("}\n\n")
	g.printf("func (*%s) Class() string {\nreturn %q\n}\n\n", goname, name)
	g.printf("func (*%s) RVersion() int16 {\nreturn %d\n}\n", goname, si.ClassVersion())
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
		goname, name,
	)

	g.genStreamerInfo(si)

	ifaces := []string{"root.Object", "rbytes.RVersioner", "rbytes.Marshaler", "rbytes.Unmarshaler"}
	g.printf("var (\n")
	for _, n := range ifaces {
		g.printf("\t_ %s = (*%s)(nil)\n", n, goname)
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
		switch se.STLType() {
		case rmeta.STLvector, rmeta.STLmap:
			tname := g.typename(se)
			g.printf(docFmt, se.Name(), tname, g.stag(i, se), doc)
		default:
			panic(fmt.Errorf("STL-type not implemented %#v", se))
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
			panic(fmt.Errorf("gen-type: unknown C++ builtin %q", tname))
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
			panic(fmt.Errorf("gen-type: unknown C++ builtin %q", tname))
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
		switch se.STLType() {
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
			case rmeta.Int64, rmeta.Long64:
				return "[]int64"
			case rmeta.Uint8:
				return "[]uint8"
			case rmeta.Uint16:
				return "[]uint16"
			case rmeta.Uint32:
				return "[]uint32"
			case rmeta.Uint64, rmeta.ULong64:
				return "[]uint64"
			case rmeta.Float32:
				return "[]float32"
			case rmeta.Float64:
				return "[]float64"
			case rmeta.Float16:
				return "[]root.Float16"
			case rmeta.Double32:
				return "[]root.Double32"
			case rmeta.Object:
				etn := se.ElemTypeName()
				switch etn[0] {
				case "string":
					return "[]string"
				default:
					etname := g.cxx2go(etn[0], qualNone)
					return "[]" + etname
				}
			default:
				panic(fmt.Errorf("invalid stl-vector element type: %v -- %#v", se.ContainedType(), se))
			}
		case rmeta.STLmap:
			types := rmeta.CxxTemplateArgsOf(se.TypeName())
			if len(types) != 2 {
				panic(fmt.Errorf(
					"invalid stl-map: got %d template arguments instead of 2 for type %q",
					len(types), se.TypeName(),
				))
			}
			k := g.cxx2go(types[0], qualNone)
			v := g.cxx2go(types[1], qualNone)
			return "map[" + k + "]" + v
		default:
			panic(fmt.Errorf("STL-type not implemented %#v", se))
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
	name = f(name)
	name = strings.Replace(name, "::", "__", -1) // handle namespaces
	return prefix + name
}

func (g *genGoType) genMarshal(si rbytes.StreamerInfo) {
	g.printf(`// MarshalROOT implements rbytes.Marshaler
func (o *%[1]s) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(o.RVersion())

`,
		g.cxx2go(si.Name(), qualNone),
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
			case rmeta.OffsetP + rmeta.Int64, rmeta.OffsetP + rmeta.Long64:
				wfunc = "WriteFastArrayI64"
			case rmeta.OffsetP + rmeta.Uint8:
				wfunc = "WriteFastArrayU8"
			case rmeta.OffsetP + rmeta.Uint16:
				wfunc = "WriteFastArrayU16"
			case rmeta.OffsetP + rmeta.Uint32:
				wfunc = "WriteFastArrayU32"
			case rmeta.OffsetP + rmeta.Uint64, rmeta.OffsetP + rmeta.ULong64:
				wfunc = "WriteFastArrayU64"
			case rmeta.OffsetP + rmeta.Float32:
				wfunc = "WriteFastArrayF32"
			case rmeta.OffsetP + rmeta.Float64:
				wfunc = "WriteFastArrayF64"
			default:
				panic(fmt.Errorf("invalid element type: %v", se.Type()))
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
					panic(fmt.Errorf("invalid counter size %d for %s.%s", se.Size(), si.Name(), se.Name()))
				}

			case rmeta.Bits:
				g.printf("w.WriteI32(int32(o.%s))\n", se.Name())

			case rmeta.Int8:
				g.printf("w.WriteI8(o.%s)\n", se.Name())
			case rmeta.Int16:
				g.printf("w.WriteI16(o.%s)\n", se.Name())
			case rmeta.Int32:
				g.printf("w.WriteI32(o.%s)\n", se.Name())
			case rmeta.Int64, rmeta.Long64:
				g.printf("w.WriteI64(o.%s)\n", se.Name())

			case rmeta.Uint8:
				g.printf("w.WriteU8(o.%s)\n", se.Name())
			case rmeta.Uint16:
				g.printf("w.WriteU16(o.%s)\n", se.Name())
			case rmeta.Uint32:
				g.printf("w.WriteU32(o.%s)\n", se.Name())
			case rmeta.Uint64, rmeta.ULong64:
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
				panic(fmt.Errorf("invalid basic type %v (%d) for %s.%s", se.Type(), se.Type(), si.Name(), se.Name()))
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
			case rmeta.OffsetL + rmeta.Int64, rmeta.OffsetL + rmeta.Long64:
				wfunc = "WriteFastArrayI64"
			case rmeta.OffsetL + rmeta.Uint8:
				wfunc = "WriteFastArrayU8"
			case rmeta.OffsetL + rmeta.Uint16:
				wfunc = "WriteFastArrayU16"
			case rmeta.OffsetL + rmeta.Uint32:
				wfunc = "WriteFastArrayU32"
			case rmeta.OffsetL + rmeta.Uint64, rmeta.OffsetL + rmeta.ULong64:
				wfunc = "WriteFastArrayU64"
			case rmeta.OffsetL + rmeta.Float32:
				wfunc = "WriteFastArrayF32"
			case rmeta.OffsetL + rmeta.Float64:
				wfunc = "WriteFastArrayF64"
			default:
				panic(fmt.Errorf("invalid array element type: %v", se.Type()))
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
		switch se.STLType() {
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
			case rmeta.Int64, rmeta.Long64:
				wfunc = "WriteFastArrayI64"
			case rmeta.Uint8:
				wfunc = "WriteFastArrayU8"
			case rmeta.Uint16:
				wfunc = "WriteFastArrayU16"
			case rmeta.Uint32:
				wfunc = "WriteFastArrayU32"
			case rmeta.Uint64, rmeta.ULong64:
				wfunc = "WriteFastArrayU64"
			case rmeta.Float32:
				wfunc = "WriteFastArrayF32"
			case rmeta.Float64:
				wfunc = "WriteFastArrayF64"
			case rmeta.Object:
				etn := se.ElemTypeName()
				switch etn[0] {
				case "string":
					wfunc = "WriteFastArrayString"
				default:
					panic(fmt.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
				}
			default:
				panic(fmt.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
			}
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
			panic(fmt.Errorf("STL-type not implemented %#v", se))
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
		g.cxx2go(si.Name(), qualNone),
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
			rsize := ""
			switch se.Type() {
			case rmeta.OffsetP + rmeta.Bool:
				rfunc = "ReadArrayBool"
				rsize = "ResizeBool"
			case rmeta.OffsetP + rmeta.Int8:
				rfunc = "ReadArrayI8"
				rsize = "ResizeI8"
			case rmeta.OffsetP + rmeta.Int16:
				rfunc = "ReadArrayI16"
				rsize = "ResizeI16"
			case rmeta.OffsetP + rmeta.Int32:
				rfunc = "ReadArrayI32"
				rsize = "ResizeI32"
			case rmeta.OffsetP + rmeta.Int64, rmeta.OffsetP + rmeta.Long64:
				rfunc = "ReadArrayI64"
				rsize = "ResizeI64"
			case rmeta.OffsetP + rmeta.Uint8:
				rfunc = "ReadArrayU8"
				rsize = "ResizeU8"
			case rmeta.OffsetP + rmeta.Uint16:
				rfunc = "ReadArrayU16"
				rsize = "ResizeU16"
			case rmeta.OffsetP + rmeta.Uint32:
				rfunc = "ReadArrayU32"
				rsize = "ResizeU32"
			case rmeta.OffsetP + rmeta.Uint64, rmeta.OffsetP + rmeta.ULong64:
				rfunc = "ReadArrayU64"
				rsize = "ResizeU64"
			case rmeta.OffsetP + rmeta.Float32:
				rfunc = "ReadArrayF32"
				rsize = "ResizeF32"
			case rmeta.OffsetP + rmeta.Float64:
				rfunc = "ReadArrayF64"
				rsize = "ResizeF64"
			default:
				panic(fmt.Errorf("invalid element type: %v", se.Type()))
			}
			g.printf("o.%s = rbytes.%s(nil, int(o.%s))\n", se.Name(), rsize, n)
			g.printf("r.%s(o.%s)\n", rfunc, se.Name())
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
					panic(fmt.Errorf("invalid counter size %d for %s.%s", se.Size(), si.Name(), se.Name()))
				}

			case rmeta.Bits:
				g.printf("o.%s = r.ReadI32()\n", se.Name())

			case rmeta.Int8:
				g.printf("o.%s = r.ReadI8()\n", se.Name())
			case rmeta.Int16:
				g.printf("o.%s = r.ReadI16()\n", se.Name())
			case rmeta.Int32:
				g.printf("o.%s = r.ReadI32()\n", se.Name())
			case rmeta.Int64, rmeta.Long64:
				g.printf("o.%s = r.ReadI64()\n", se.Name())

			case rmeta.Uint8:
				g.printf("o.%s = r.ReadU8()\n", se.Name())
			case rmeta.Uint16:
				g.printf("o.%s = r.ReadU16()\n", se.Name())
			case rmeta.Uint32:
				g.printf("o.%s = r.ReadU32()\n", se.Name())
			case rmeta.Uint64, rmeta.ULong64:
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
				panic(fmt.Errorf("invalid basic type %v (%d) for %s.%s", se.Type(), se.Type(), si.Name(), se.Name()))
			}
		default:
			rfunc := ""
			switch se.Type() {
			case rmeta.OffsetL + rmeta.Bool:
				rfunc = "ReadArrayBool"
			case rmeta.OffsetL + rmeta.Int8:
				rfunc = "ReadArrayI8"
			case rmeta.OffsetL + rmeta.Int16:
				rfunc = "ReadArrayI16"
			case rmeta.OffsetL + rmeta.Int32:
				rfunc = "ReadArrayI32"
			case rmeta.OffsetL + rmeta.Int64, rmeta.OffsetL + rmeta.Long64:
				rfunc = "ReadArrayI64"
			case rmeta.OffsetL + rmeta.Uint8:
				rfunc = "ReadArrayU8"
			case rmeta.OffsetL + rmeta.Uint16:
				rfunc = "ReadArrayU16"
			case rmeta.OffsetL + rmeta.Uint32:
				rfunc = "ReadArrayU32"
			case rmeta.OffsetL + rmeta.Uint64, rmeta.OffsetL + rmeta.ULong64:
				rfunc = "ReadArrayU64"
			case rmeta.OffsetL + rmeta.Float32:
				rfunc = "ReadArrayF32"
			case rmeta.OffsetL + rmeta.Float64:
				rfunc = "ReadArrayF64"
			default:
				panic(fmt.Errorf("invalid array element type: %v", se.Type()))
			}
			g.printf("r.%s(o.%s[:])\n", rfunc, se.Name())
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
		switch se.STLType() {
		case rmeta.STLvector:
			rfunc := ""
			rsize := ""
			switch se.ContainedType() {
			case rmeta.Bool:
				rfunc = "ReadArrayBool"
				rsize = "ResizeBool"
			case rmeta.Int8:
				rfunc = "ReadArrayI8"
				rsize = "ResizeI8"
			case rmeta.Int16:
				rfunc = "ReadArrayI16"
				rsize = "ResizeI16"
			case rmeta.Int32:
				rfunc = "ReadArrayI32"
				rsize = "ResizeI32"
			case rmeta.Int64, rmeta.Long64:
				rfunc = "ReadArrayI64"
				rsize = "ResizeI64"
			case rmeta.Uint8:
				rfunc = "ReadArrayU8"
				rsize = "ResizeU8"
			case rmeta.Uint16:
				rfunc = "ReadArrayU16"
				rsize = "ResizeU16"
			case rmeta.Uint32:
				rfunc = "ReadArrayU32"
				rsize = "ResizeU32"
			case rmeta.Uint64, rmeta.ULong64:
				rfunc = "ReadArrayU64"
				rsize = "ResizeU64"
			case rmeta.Float32:
				rfunc = "ReadArrayF32"
				rsize = "ResizeF32"
			case rmeta.Float64:
				rfunc = "ReadArrayF64"
				rsize = "ResizeF64"
			case rmeta.Object:
				etn := se.ElemTypeName()
				switch etn[0] {
				case "string":
					rfunc = "ReadArrayString"
					rsize = "ResizeStr"
				default:
					panic(fmt.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
				}
			default:
				panic(fmt.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
			}
			g.imps["fmt"] = 1
			g.imps["go-hep.org/x/hep/groot/rvers"] = 1
			g.printf("{\n")
			g.printf("vers, pos, bcnt := r.ReadVersion(%q)\n", se.TypeName())
			g.printf("if vers != rvers.StreamerInfo {\n")
			g.printf("r.SetErr(fmt.Errorf(\"rbytes: invalid version for \\\"%s\\\". got=%%v, want=%%v\", vers, rvers.StreamerInfo))\n", se.TypeName())
			g.printf("return r.Err()\n")
			g.printf("}\n")
			g.printf("o.%s = rbytes.%s(nil, int(r.ReadI32()))\n", se.Name(), rsize)
			g.printf("r.%s(o.%s)\n", rfunc, se.Name())
			g.printf("r.CheckByteCount(pos, bcnt, start, %q)\n", se.TypeName())
			g.printf("}\n")

		default:
			panic(fmt.Errorf("STL-type not implemented %#v", se))
		}

	default:
		g.printf("// %s -- %T\n", se.Name(), se)
	}
}

func (g *genGoType) genStreamerInfo(si rbytes.StreamerInfo) {
	g.printf(`func init() {
		// Streamer for %[1]s.
		%[4]sStreamerInfos.Add(%[4]sNewCxxStreamerInfo(%[1]q, %[2]d, 0x%[3]x, []rbytes.StreamerElement{
`,
		si.Name(), si.ClassVersion(), si.CheckSum(),
		g.rdict,
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
	if g.rdict != "" {
		g.imps["go-hep.org/x/hep/groot/rdict"] = 1
	}
	g.imps["go-hep.org/x/hep/groot/rmeta"] = 1

	switch se := se.(type) {
	case *StreamerBase:
		g.printf(`%[13]sNewStreamerBase(%[13]sElement{
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
			g.rdict,
		)

	case *StreamerBasicType:
		g.printf(`&%[13]sStreamerBasicType{StreamerElement: %[13]sElement{
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
			g.rdict,
		)

	case *StreamerBasicPointer:
		g.printf(`%[16]sNewStreamerBasicPointer(%[16]sElement{
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
			g.rdict,
		)

	case *StreamerLoop:
		g.printf(`%[15]sNewStreamerLoop(%[15]sElement{
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
			g.rdict,
		)

	case *StreamerObject:
		g.printf(`&%[13]sStreamerObject{StreamerElement: %[13]sElement{
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
			g.rdict,
		)

	case *StreamerObjectPointer:
		g.printf(`&%[13]sStreamerObjectPointer{StreamerElement: %[13]sElement{
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
			g.rdict,
		)

	case *StreamerObjectAny:
		g.printf(`&%[13]sStreamerObjectAny{StreamerElement: %[13]sElement{
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
			g.rdict,
		)

	case *StreamerObjectAnyPointer:
		g.printf(`&%[13]sStreamerObjectAnyPointer{StreamerElement: %[13]sElement{
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
			g.rdict,
		)

	case *StreamerString:
		g.printf(`&%[13]sStreamerString{StreamerElement: %[13]sElement{
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
			g.rdict,
		)

	case *StreamerSTL:
		g.printf(`%[15]sNewCxxStreamerSTL(%[15]sElement{
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
			g.rdict,
		)

	case *StreamerSTLstring:
		g.printf(`&%[15]sStreamerSTLstring{*%[15]sNewCxxStreamerSTL(%[15]sElement{
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
			g.rdict,
		)

	case *StreamerArtificial:
		g.printf(`&%[13]sStreamerArtificial{StreamerElement: %[13]sElement{
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
			g.rdict,
		)

	default:
		g.printf("// %s -- %T\n", se.Name(), se)
		panic(fmt.Errorf("rdict: unknown streamer type %T (%#v)", se, se))
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
