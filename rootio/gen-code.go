// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// gen-code generates code for simple ROOT classes hierarchies.
package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func main() {
	genLeaves()
	genArrays()
	genH1()
	genH2()
}

func gofmt(f *os.File) {
	fname := f.Name()
	src, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	src, err = format.Source(src)
	if err != nil {
		log.Fatalf("error formating sources of %q: %v\n", fname, err)
	}

	err = ioutil.WriteFile(fname, src, 0644)
	if err != nil {
		log.Fatalf("error writing back %q: %v\n", fname, err)
	}
}

func genLeaves() {
	f, err := os.Create("leaf_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, srcHeader)

	for i, typ := range []struct {
		Name       string
		Type       string
		DoUnsigned bool
		Func       string
		FuncArray  string
		Count      bool
	}{
		{
			Name:      "LeafO",
			Type:      "bool",
			Func:      "r.ReadBool()",
			FuncArray: "r.ReadFastArrayBool",
		},
		{
			Name:       "LeafS",
			Type:       "int16",
			DoUnsigned: true,
			Func:       "r.ReadI16()",
			FuncArray:  "r.ReadFastArrayI16",
			Count:      true,
		},
		{
			Name:      "LeafC",
			Type:      "int32",
			Func:      "r.ReadI32()",
			FuncArray: "r.ReadFastArrayI32",
		},
		{
			Name:       "LeafI",
			Type:       "int32",
			DoUnsigned: true,
			Func:       "r.ReadI32()",
			FuncArray:  "r.ReadFastArrayI32",
			Count:      true,
		},
		{
			Name:       "LeafL",
			Type:       "int64",
			DoUnsigned: true,
			Func:       "r.ReadI64()",
			FuncArray:  "r.ReadFastArrayI64",
			Count:      true,
		},
		{
			Name:      "LeafF",
			Type:      "float32",
			Func:      "r.ReadF32()",
			FuncArray: "r.ReadFastArrayF32",
		},
		{
			Name:      "LeafD",
			Type:      "float64",
			Func:      "r.ReadF64()",
			FuncArray: "r.ReadFastArrayF64",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(leafTmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}

func genArrays() {
	f, err := os.Create("array_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, srcHeader)

	for i, typ := range []struct {
		Name string
		Type string
		Func string
	}{
		{
			Name: "ArrayI",
			Type: "int32",
			Func: "r.ReadFastArrayI32",
		},
		{
			Name: "ArrayL64",
			Type: "int64",
			Func: "r.ReadFastArrayI64",
		},
		{
			Name: "ArrayF",
			Type: "float32",
			Func: "r.ReadFastArrayF32",
		},
		{
			Name: "ArrayD",
			Type: "float64",
			Func: "r.ReadFastArrayF64",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(arrayTmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}

func genH1() {
	f, err := os.Create("h1_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, srcHeader)

	for i, typ := range []struct {
		Name string
		Type string
	}{
		{
			Name: "H1F",
			Type: "ArrayF",
		},
		{
			Name: "H1D",
			Type: "ArrayD",
		},
		{
			Name: "H1I",
			Type: "ArrayI",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(h1Tmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}

func genH2() {
	f, err := os.Create("h2_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, srcHeader)

	for i, typ := range []struct {
		Name string
		Type string
	}{
		{
			Name: "H2F",
			Type: "ArrayF",
		},
		{
			Name: "H2D",
			Type: "ArrayD",
		},
		{
			Name: "H2I",
			Type: "ArrayI",
		},
	} {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		tmpl := template.Must(template.New(typ.Name).Parse(h2Tmpl))
		err = tmpl.Execute(f, typ)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", typ.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	gofmt(f)
}

const srcHeader = `// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rootio

import (
	"reflect"
)

`

const leafTmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	tleaf
	val []{{.Type}}
	min	{{.Type}}
	max {{.Type}}
}

// Class returns the ROOT class name.
func (leaf *{{.Name}}) Class() string {
	return "T{{.Name}}"
}

// Minimum returns the minimum value of the leaf.
func (leaf *{{.Name}}) Minimum() {{.Type}} {
	return leaf.min
}

// Maximum returns the maximum value of the leaf.
func (leaf *{{.Name}}) Maximum() {{.Type}} {
	return leaf.max
}

// Value returns the leaf value at index i.
func (leaf *{{.Name}}) Value(i int) interface{} {
	return leaf.val[i]
}

// value returns the leaf value.
func (leaf *{{.Name}}) value() interface{} {
	return leaf.val
}

{{if .Count}}
// ivalue returns the first leaf value as int
func (leaf *{{.Name}}) ivalue() int {
	return int(leaf.val[0])
}

// imax returns the leaf maximum value as int
func (leaf *{{.Name}}) imax() int {
	return int(leaf.max)
}
{{end}}

func (leaf *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	start := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	myprintf("{{.Name}}: %v %v %v\n", vers, pos, bcnt)

	if err := leaf.tleaf.UnmarshalROOT(r); err != nil {
		r.err = err
		return r.err
	}

	leaf.min = {{.Func}}
	leaf.max = {{.Func}}

	r.CheckByteCount(pos, bcnt, start, "T{{.Name}}")
	return r.Err()
}

func (leaf *{{.Name}}) readBasket(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if leaf.count == nil && len(leaf.val) == 1 {
		leaf.val[0] = {{.Func}}
	} else {
		if leaf.count != nil {
			entry := leaf.Branch().getReadEntry()
			if leaf.count.Branch().getReadEntry() != entry {
				leaf.count.Branch().getEntry(entry)
			}
			n := leaf.count.ivalue()
			max := leaf.count.imax()
			if n > max {
				n = max
			}
			leaf.val = {{.FuncArray}}(leaf.tleaf.len * n)
		} else {
			leaf.val = {{.FuncArray}}(leaf.tleaf.len)
		}
	}
	return r.err
}

func (leaf *{{.Name}}) scan(r *RBuffer, ptr interface{}) error {
	if r.err != nil {
		return r.err
	}

	if rv := reflect.Indirect(reflect.ValueOf(ptr)); rv.Kind() == reflect.Array {
		return leaf.scan(r, rv.Slice(0, rv.Len()).Interface())
	}

	switch v := ptr.(type) {
	case *{{.Type}}:
		*v = leaf.val[0]
	case *[]{{.Type}}:
		if len(*v) < len(leaf.val) || *v == nil {
			*v = make([]{{.Type}}, len(leaf.val))
		}
		copy(*v, leaf.val)
		*v = (*v)[:leaf.count.ivalue()]
	case []{{.Type}}:
		if len(v) < len(leaf.val) {
			v = make([]{{.Type}}, len(leaf.val))
		}
		copy(v, leaf.val)
{{if .DoUnsigned}}
	case *u{{.Type}}:
		*v = u{{.Type}}(leaf.val[0])
	case *[]u{{.Type}}:
		if len(*v) < len(leaf.val) || *v == nil {
			*v = make([]u{{.Type}}, len(leaf.val))
		}
		for i, u := range leaf.val {
			(*v)[i] = u{{.Type}}(u)
		}
		*v = (*v)[:leaf.count.ivalue()]
	case []u{{.Type}}:
		if len(v) < len(leaf.val) {
			v = make([]u{{.Type}}, len(leaf.val))
		}
		for i := range v {
			v[i] = u{{.Type}}(leaf.val[i])
		}
{{end}}
	default:
		panic(errorf("invalid ptr type %T (leaf=%s|%T)", v, leaf.Name(), leaf))
	}

	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &{{.Name}}{}
		return reflect.ValueOf(o)
	}
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var _ Object = (*{{.Name}})(nil)
var _ Named = (*{{.Name}})(nil)
var _ Leaf = (*{{.Name}})(nil)
var _ ROOTUnmarshaler = (*{{.Name}})(nil)
`

const arrayTmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	Data []{{.Type}}
}

// Class returns the ROOT class name.
func (*{{.Name}}) Class() string {
	return "T{{.Name}}"
}

func (arr *{{.Name}}) Len() int {
	return len(arr.Data)
}

func (arr *{{.Name}}) At(i int) {{.Type}} {
	return arr.Data[i]
}

func (arr *{{.Name}}) Get(i int) interface{} {
	return arr.Data[i]
}

func (arr *{{.Name}}) Set(i int, v interface{}) {
	arr.Data[i] = v.({{.Type}})
}

func (arr *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	n := int(r.ReadI32())
	arr.Data = {{.Func}}(n)

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &{{.Name}}{}
		return reflect.ValueOf(o)
	}
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var _ Array = (*{{.Name}})(nil)
var _ ROOTUnmarshaler = (*{{.Name}})(nil)
`

const h1Tmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	th1
	arr {{.Type}}
}

// Class returns the ROOT class name.
func (*{{.Name}}) Class() string {
	return "T{{.Name}}"
}

func (h *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: T{{.Name}} version too old (%d<2)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th1,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "T{{.Name}}")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &{{.Name}}{}
		return reflect.ValueOf(o)
	}
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var _ Object = (*{{.Name}})(nil)
var _ Named = (*{{.Name}})(nil)
var _ ROOTUnmarshaler = (*{{.Name}})(nil)
`

const h2Tmpl = `// {{.Name}} implements ROOT T{{.Name}}
type {{.Name}} struct {
	th2
	arr {{.Type}}
}

// Class returns the ROOT class name.
func (*{{.Name}}) Class() string {
	return "T{{.Name}}"
}

func (h *{{.Name}}) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion()
	if vers < 2 {
		return errorf("rootio: T{{.Name}} version too old (%d<2)", vers)
	}

	for _, v := range []ROOTUnmarshaler{
		&h.th2,
		&h.arr,
	} {
		if err := v.UnmarshalROOT(r); err != nil {
			r.err = err
			return r.err
		}
	}

	r.CheckByteCount(pos, bcnt, beg, "T{{.Name}}")
	return r.err
}

func init() {
	f := func() reflect.Value {
		o := &{{.Name}}{}
		return reflect.ValueOf(o)
	}
	Factory.add("T{{.Name}}", f)
	Factory.add("*rootio.{{.Name}}", f)
}

var _ Object = (*{{.Name}})(nil)
var _ Named = (*{{.Name}})(nil)
var _ ROOTUnmarshaler = (*{{.Name}})(nil)
`
