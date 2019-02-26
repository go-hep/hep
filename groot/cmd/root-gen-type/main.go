// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command root-gen-type generates a Go type from the StreamerInfo contained
// in a ROOT file.
package main // import "go-hep.org/x/hep/groot/cmd/root-gen-type"

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rdict"
	"go-hep.org/x/hep/groot/rmeta"
	"go-hep.org/x/hep/groot/rtypes"
	_ "go-hep.org/x/hep/groot/ztypes"
)

func main() {
	log.SetPrefix("root-gen-type: ")
	log.SetFlags(0)

	var (
		typeNames = flag.String("t", ".*", "comma-separated list of (regexp) type names")
		pkgPath   = flag.String("p", "", "package import path")
		output    = flag.String("o", "", "output file name")
		streamers = flag.Bool("streamers", false, "enable generation of MarshalROOT/UnmarshalROOT methods")
		verbose   = flag.Bool("v", false, "enable verbose mode")
	)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: root-gen-type [options] input.root

ex:
 $> root-gen-type -p mypkg -t MyType -o streamers_gen.go ./input.root

options:
`,
		)
		flag.PrintDefaults()
	}

	flag.Parse()

	if *typeNames == "" {
		flag.Usage()
		os.Exit(2)
	}

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	types := strings.Split(*typeNames, ",")

	var (
		err error
		out io.WriteCloser
	)

	switch *output {
	case "":
		out = os.Stdout
	default:
		out, err = os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
	}

	err = generate(out, *pkgPath, types, flag.Arg(0), *verbose, *streamers)
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func generate(w io.Writer, pkg string, types []string, fname string, verbose, streamers bool) error {
	f, err := groot.Open(fname)
	if err != nil {
		return err
	}

	g, err := newGenerator(pkg, f)
	if err != nil {
		return err
	}
	g.Verbose = verbose
	g.Streamers = streamers

	filters := make([]*regexp.Regexp, len(types))
	for i, t := range types {
		filters[i] = regexp.MustCompile(t)
	}

	accept := func(name string) string {
		for _, filter := range filters {
			if filter.MatchString(name) {
				return name
			}
		}
		return ""
	}
	for _, si := range f.StreamerInfos() {
		if t := accept(si.Name()); t != "" {
			g.Generate(t)
		}
	}

	buf, err := g.Format()
	if err != nil {
		return err
	}

	_, err = w.Write(buf)
	if err != nil {
		return err
	}

	return nil
}

type Generator struct {
	pkg string
	buf *bytes.Buffer
	ctx rbytes.StreamerInfoContext

	Streamers bool // generate streamers
	Verbose   bool

	// set of imported packages.
	// usually: "go-hep.org/x/hep/groot/rbase", ".../rcont
	imps map[string]int
}

func newGenerator(pkg string, sictx rbytes.StreamerInfoContext) (*Generator, error) {
	return &Generator{
		pkg: pkg,
		buf: new(bytes.Buffer),
		ctx: sictx,
		imps: map[string]int{
			"reflect":                       1,
			"go-hep.org/x/hep/groot/rbytes": 1,
			"go-hep.org/x/hep/groot/root":   1,
			"go-hep.org/x/hep/groot/rtypes": 1,
		},
	}, nil
}

func (g *Generator) Generate(name string) {
	if g.Verbose {
		log.Printf("generating type for %q...", name)
	}
	si, err := g.ctx.StreamerInfo(name, -1)
	if err != nil {
		log.Panic(err)
	}

	g.genType(si)
}

func (g *Generator) genType(si rbytes.StreamerInfo) {
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
	ifaces := []string{"root.Object", "rbytes.RVersioner"}
	if g.Streamers {
		g.printf("\n")
		g.genMarshal(si)
		g.genUnmarshal(si)
		ifaces = append(ifaces, "rbytes.Marshaler", "rbytes.Unmarshaler")
	}

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

	g.printf("var (\n")
	for _, n := range ifaces {
		g.printf("\t_ %s = (*%s)(nil)\n", n, name)
	}
	g.printf(")\n\n")
}

func (g *Generator) genField(si rbytes.StreamerInfo, i int, se rbytes.StreamerElement) {
	const (
		docFmt = "\t%s\t%s\t%s\n"
	)
	doc := se.Title()
	if doc != "" {
		doc = "// " + doc
	}
	switch se := se.(type) {
	case *rdict.StreamerBase:
		g.printf(docFmt, fmt.Sprintf("base%d", i), g.typename(se), "// base class")

	case *rdict.StreamerBasicPointer:
		g.printf(docFmt, se.Name(), g.typename(se), doc)

	case *rdict.StreamerBasicType:
		tname := g.typename(se)
		switch se.ArrayLen() {
		case 0:
		default:
			tname = fmt.Sprintf("[%d]%s", se.ArrayLen(), tname)
		}
		g.printf(docFmt, se.Name(), tname, doc)

	case *rdict.StreamerObject:
		tname := g.typename(se)
		g.printf(docFmt, se.Name(), tname, doc)

	case *rdict.StreamerObjectAny:
		tname := g.typename(se)
		g.printf(docFmt, se.Name(), tname, doc)

	case *rdict.StreamerObjectAnyPointer:
		tname := g.typename(se)
		g.printf(docFmt, se.Name(), tname, doc)

	case *rdict.StreamerObjectPointer:
		tname := g.typename(se)
		g.printf(docFmt, se.Name(), tname, doc)

	case *rdict.StreamerString, *rdict.StreamerSTLstring:
		g.printf(docFmt, se.Name(), "string", doc)

	case *rdict.StreamerSTL:
		switch se.STLVectorType() {
		case rmeta.STLvector:
			tname := g.typename(se)
			g.printf(docFmt, se.Name(), tname, doc)
		default:
			panic(errors.Errorf("STL-type not implemented %#v", se))
		}
	default:
		g.printf("\t%s\t%s // %T -- %s\n", se.Name(), g.typename(se), se, doc)
	}
}

func (g *Generator) typename(se rbytes.StreamerElement) string {
	tname := se.TypeName()
	switch se := se.(type) {
	case *rdict.StreamerBase:
		return g.cxx2go(se.Name(), 0)

	case *rdict.StreamerBasicPointer:
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

	case *rdict.StreamerBasicType:
		t, ok := rmeta.CxxBuiltins[tname]
		if !ok {
			panic(errors.Errorf("gen-type: unknown C++ builtin %q", tname))
		}
		return t.Name()

	case *rdict.StreamerObject:
		return g.cxx2go(tname, qualNone)

	case *rdict.StreamerObjectAny:
		return g.cxx2go(tname, qualNone)

	case *rdict.StreamerObjectAnyPointer:
		tname = tname[:len(tname)-1] // drop last '*'
		return g.cxx2go(tname, qualStar)

	case *rdict.StreamerObjectPointer:
		tname = tname[:len(tname)-1] // drop last '*'
		return g.cxx2go(tname, qualStar)

	case *rdict.StreamerString, *rdict.StreamerSTLstring:
		return "string"

	case *rdict.StreamerSTL:
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
				switch se.TypeName() {
				case "vector<string>":
					return "[]string"
				default:
					panic(errors.Errorf("invalid stl-vector element type: %v -- %#v", se.ContainedType(), se))
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

func (g *Generator) cxx2go(name string, qual qualKind) string {
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
	return prefix + name
}

func (g *Generator) genMarshal(si rbytes.StreamerInfo) {
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

func (g *Generator) genMarshalField(si rbytes.StreamerInfo, i int, se rbytes.StreamerElement) {
	switch se := se.(type) {
	case *rdict.StreamerBase:
		g.printf("o.%s.MarshalROOT(w)\n", fmt.Sprintf("base%d", i))

	case *rdict.StreamerBasicPointer:
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

	case *rdict.StreamerBasicType:
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

	case *rdict.StreamerObject:
		// FIXME(sbinet): check semantics
		g.printf("o.%s.MarshalROOT(w) // obj\n", se.Name())

	case *rdict.StreamerObjectAny:
		// FIXME(sbinet): check semantics
		g.printf("o.%s.MarshalROOT(w) // obj-any\n", se.Name())

	case *rdict.StreamerObjectAnyPointer:
		// FIXME(sbinet): check semantics
		g.printf("w.WriteObjectAny(o.%s) // obj-any-ptr\n", se.Name())

	case *rdict.StreamerObjectPointer:
		// FIXME(sbinet): check semantics
		g.printf("w.WriteObjectAny(o.%s) // obj-ptr \n", se.Name())

	case *rdict.StreamerString, *rdict.StreamerSTLstring:
		g.printf("w.WriteString(o.%s)\n", se.Name())

	case *rdict.StreamerSTL:
		switch se.STLVectorType() {
		case rmeta.STLvector:
			g.printf("w.WriteI32(int32(len(o.%s)))\n", se.Name())
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
				switch se.TypeName() {
				case "vector<string>":
					wfunc = "WriteFastArrayString"
				default:
					panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
				}
			default:
				panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
			}
			g.printf("w.%s(o.%s)\n", wfunc, se.Name())
		default:
			panic(errors.Errorf("STL-type not implemented %#v", se))
		}

	default:
		g.printf("// %s -- %T\n", se.Name(), se)
	}
}

func (g *Generator) genUnmarshal(si rbytes.StreamerInfo) {
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

func (g *Generator) genUnmarshalField(si rbytes.StreamerInfo, i int, se rbytes.StreamerElement) {
	switch se := se.(type) {
	case *rdict.StreamerBase:
		g.printf("o.%s.UnmarshalROOT(r)\n", fmt.Sprintf("base%d", i))

	case *rdict.StreamerBasicPointer:
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

	case *rdict.StreamerBasicType:
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

	case *rdict.StreamerObject:
		// FIXME(sbinet): check semantics
		g.printf("o.%s.UnmarshalROOT(r) // obj\n", se.Name())

	case *rdict.StreamerObjectAny:
		// FIXME(sbinet): check semantics
		g.printf("o.%s.UnmarshalROOT(r) // obj-any\n", se.Name())

	case *rdict.StreamerObjectAnyPointer:
		// FIXME(sbinet): check semantics
		g.printf("{\n")
		g.printf("o.%s = nil\n", se.Name())
		g.printf("if oo := r.ReadObjectAny(); oo != nil {  // obj-any-ptr\n")
		g.printf("o.%s = oo.(%s)\n", se.Name(), g.typename(se))
		g.printf("}\n}\n")

	case *rdict.StreamerObjectPointer:
		// FIXME(sbinet): check semantics
		g.printf("{\n")
		g.printf("o.%s = nil\n", se.Name())
		g.printf("if oo := r.ReadObjectAny(); oo != nil {  // obj-ptr\n")
		g.printf("o.%s = oo.(%s)\n", se.Name(), g.typename(se))
		g.printf("}\n}\n")

	case *rdict.StreamerString, *rdict.StreamerSTLstring:
		g.printf("o.%s = r.ReadString()\n", se.Name())

	case *rdict.StreamerSTL:
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
				switch se.TypeName() {
				case "vector<string>":
					rfunc = "ReadFastArrayString"
				default:
					panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
				}
			default:
				panic(errors.Errorf("invalid stl-vector element type: %v", se.ContainedType()))
			}
			g.printf("o.%s = r.%s(int(r.ReadI32()))\n", se.Name(), rfunc)
		default:
			panic(errors.Errorf("STL-type not implemented %#v", se))
		}

	default:
		g.printf("// %s -- %T\n", se.Name(), se)
	}
}

func (g *Generator) printf(format string, args ...interface{}) {
	fmt.Fprintf(g.buf, format, args...)
}

func (g *Generator) Format() ([]byte, error) {
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

func init() {
	f := func() reflect.Value {
		o := &rcont.ArrayL64{}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TArrayL", f)
}
