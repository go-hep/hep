// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// gen-code generates code for simple ROOT classes hierarchies.
package main

import (
	"fmt"
	"log"
	"os"
	"text/template"
)

func main() {
	genLeaves()
	genArrays()
}

func genLeaves() {
	f, err := os.Create("leaf_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, srcHeader)

	for _, typ := range []struct {
		Name string
		Type string
		Func string
	}{
		{
			Name: "LeafC",
			Type: "int32",
			Func: "r.ReadI32()",
		},
		{
			Name: "LeafI",
			Type: "int32",
			Func: "r.ReadI32()",
		},
		{
			Name: "LeafL",
			Type: "int64",
			Func: "r.ReadI64()",
		},
		{
			Name: "LeafF",
			Type: "float32",
			Func: "r.ReadF32()",
		},
		{
			Name: "LeafD",
			Type: "float64",
			Func: "r.ReadF64()",
		},
	} {
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
}

func genArrays() {
	f, err := os.Create("array_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintf(f, srcHeader)

	for _, typ := range []struct {
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
