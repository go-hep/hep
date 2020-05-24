// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"

	"go-hep.org/x/hep/groot/internal/genroot"
)

func main() {
	genRFuncs()
}

func genRFuncs() {

	for _, typ := range []struct {
		N int
		F Func
	}{
		{
			N: 10,
			F: Func{
				Name:    "Bool",
				InType:  "float64",
				OutType: "bool",
			},
		},
		{
			N: 10,
			F: Func{
				Name:    "U32",
				InType:  "uint32",
				OutType: "uint32",
			},
		},
		{
			N: 10,
			F: Func{
				Name:    "U64",
				InType:  "uint64",
				OutType: "uint64",
			},
		},
		{
			N: 10,
			F: Func{
				Name:    "I32",
				InType:  "int32",
				OutType: "int32",
			},
		},
		{
			N: 10,
			F: Func{
				Name:    "I64",
				InType:  "int64",
				OutType: "int64",
			},
		},
		{
			N: 10,
			F: Func{
				Name:    "F32",
				InType:  "float32",
				OutType: "float32",
			},
		},
		{
			N: 10,
			F: Func{
				Name:    "F64",
				InType:  "float64",
				OutType: "float64",
			},
		},
	} {
		genRFunc(typ.N, typ.F)
	}
}

func genRFunc(arity int, typ Func) {
	genRFuncCode(arity, typ)
	genRFuncTest(arity, typ)
}

func genRFuncCode(arity int, typ Func) {
	f, err := os.Create("./rfunc_" + typ.OutType + "_gen.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rfunc", f,
		"fmt",
		"reflect",
	)

	for i := 0; i <= arity; i++ {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		fct := typ.gen(i)
		tmpl := template.Must(template.New(fct.Name).Parse(rfuncCodeTmpl))
		err = tmpl.Execute(f, fct)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", fct.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)
}

func genRFuncTest(arity int, typ Func) {
	f, err := os.Create("./rfunc_" + typ.OutType + "_gen_test.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	genroot.GenImports("rfunc", f,
		"reflect",
		"testing",
	)

	for i := 0; i <= arity; i++ {
		if i > 0 {
			fmt.Fprintf(f, "\n")
		}
		fct := typ.gen(i)
		tmpl := template.Must(template.New(fct.Name).Parse(rfuncTestTmpl))
		err = tmpl.Execute(f, fct)
		if err != nil {
			log.Fatalf("error executing template for %q: %v\n", fct.Name, err)
		}
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
	genroot.GoFmt(f)
}

type Func struct {
	Name    string
	InType  string
	OutType string
	In      []string
	Out     []string
}

func (f Func) gen(arity int) Func {
	ins := make([]string, arity)
	for i := range ins {
		ins[i] = f.InType
	}
	out := []string{f.OutType}
	return Func{
		Name:    f.Name,
		InType:  f.InType,
		OutType: f.OutType,
		In:      ins,
		Out:     out,
	}
}

func (f Func) NumIn() int  { return len(f.In) }
func (f Func) NumOut() int { return len(f.Out) }
func (f Func) Type() string {
	return fmt.Sprintf("funcAr%02d%s", f.NumIn(), f.Name)
}

func (f Func) Func() string {
	sig := new(strings.Builder)
	sig.WriteString("func(")
	for i, typ := range f.In {
		if i > 0 {
			sig.WriteString(", ")
		}
		fmt.Fprintf(sig, "arg%02d %s", i, typ)
	}
	sig.WriteString(")")

	sig.WriteString(f.Return())

	return sig.String()
}

func (f Func) Return() string {
	sig := new(strings.Builder)
	switch len(f.Out) {
	case 0:
		// no-op
	case 1:
		sig.WriteString(" ")
	default:
		sig.WriteString(" (")
	}
	for i, typ := range f.Out {
		if i > 0 {
			sig.WriteString(", ")
		}
		sig.WriteString(typ)
	}
	switch len(f.Out) {
	case 0, 1:
		// no-op
	default:
		sig.WriteString(")")
	}

	return sig.String()
}

func (f Func) TestFunc() string {
	switch f.OutType {
	case "string":
		return `"42"`
	case "bool":
		return "true"
	default:
		return "42"
	}
}

func (f Func) ExportType() string {
	return strings.Title(f.Type())
}

const rfuncCodeTmpl = `// {{.Type}} implements rfunc.Formula
type {{.Type}} struct {
{{- if gt .NumIn 0}}
	rvars []string
{{- end}}
{{- range $i, $typ := .In}}
	arg{{$i}} *{{$typ}}
{{- end}}
	fct {{.Func}}
}

func new{{.ExportType}}(rvars []string, fct interface{}) (Formula, error) {
	return &{{.Type}}{
{{- if gt .NumIn 0}}
		rvars: rvars,
{{- end}}
		fct: fct.({{.Func}}),
	}, nil
}

func init() {
	funcs[reflect.TypeOf({{.Type}}{}.fct)] = new{{.ExportType}}
}

{{if gt .NumIn 0}}
func (f *{{.Type}}) RVars() []string { return f.rvars }
{{else}}
func (f *{{.Type}}) RVars() []string { return nil }
{{end}}

func (f *{{.Type}}) Bind(args []interface{}) error {
	if got, want := len(args), {{.NumIn}}; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
{{- range $i, $typ := .In}}
	{
		ptr, ok := args[{{$i}}].(*{{$typ}})
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type {{$i}} (name=%s) mismatch: got=%T, want=*{{$typ}}",
				f.rvars[{{$i}}], args[{{$i}}],
			)
		}
		f.arg{{$i}} = ptr
	}
{{- end}}
	return nil
}

func (f *{{.Type}}) Func() interface{} {
	return func() {{.Return}} {
		return f.fct(
{{- range $i, $typ := .In}}
			*f.arg{{$i}},
{{- end}}
		)
	}
}

var (
	_ Formula = (*{{.Type}})(nil)
)
`

const rfuncTestTmpl = `func Test{{.ExportType}}(t *testing.T) {
{{if gt .NumIn 0}}
	rvars := make([]string, {{.NumIn}})
{{- else}}
	var rvars []string
{{- end}}
{{- range $i, $typ := .In}}
	rvars[{{$i}}] = "name-{{$i}}"
{{- end}}

	fct := {{.Func}} {
		return {{.TestFunc}}
	}

	form, err := new{{.ExportType}}(rvars, fct)
	if err != nil {
		t.Fatalf("could not create {{.Type}} formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

{{if gt .NumIn 0}}
	ptrs := make([]interface{}, {{.NumIn}})
{{- range $i, $typ := .In}}
	ptrs[{{$i}}] = new({{$typ}})
{{- end}}
{{else}}
	var ptrs []interface{}
{{- end}}

{{if gt .NumIn 0}}
	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs)-1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}
{{- else}}
	{
		bad := make([]interface{}, 1)
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}
{{- end}}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func () {{.Return}})()
	if got, want := got, {{.OutType}}({{.TestFunc}}); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
`
