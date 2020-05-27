// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package genroot // import "go-hep.org/x/hep/groot/internal/genroot"

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"strings"
	"text/template"

	"golang.org/x/tools/go/packages"
)

// RFunc describes which function should be used as a template
// to implement the rtree/rfunc.Formula interface.
type RFunc struct {
	Pkg  string // Name of package hosting the formula to be generated.
	Path string // Import path of the package holding the function.
	Name string // Formula name.
	Def  string // Function name or signature.
}

// GenRFunc generates the rtree/rfunc.Formula implementation for fct.
func GenRFunc(w io.Writer, fct RFunc) error {
	gen, err := NewRFuncGenerator(w, fct)
	if err != nil {
		return fmt.Errorf("genroot: could not create rfunc generator: %w", err)
	}

	err = gen.Generate()
	if err != nil {
		return fmt.Errorf("genroot: could not generate rfunc formula implementation: %w", err)
	}
	return nil
}

type rfuncGen struct {
	w    io.Writer
	f    *types.Signature
	pkg  string // "rfunc." or ""
	name string
}

func NewRFuncGenerator(w io.Writer, fct RFunc) (*rfuncGen, error) {
	var (
		f   *types.Signature
		err error
	)
	switch fct.Path {
	case "":
		f, err = parseExpr(fct.Def)
		if err != nil {
			return nil, fmt.Errorf("genroot: could not parse function signature: %w", err)
		}
	default:
		cfg := &packages.Config{
			Mode: packages.NeedName |
				packages.NeedFiles |
				packages.NeedCompiledGoFiles |
				packages.NeedSyntax |
				packages.NeedTypes |
				packages.NeedTypesInfo,
		}
		pkgs, err := packages.Load(cfg, fct.Path)
		if err != nil {
			return nil, fmt.Errorf("genroot: could not load package of %q %s: %w", fct.Path, fct.Name, err)
		}
		var pkg *packages.Package
		for _, p := range pkgs {
			if p.PkgPath == fct.Path {
				pkg = p
				break
			}
		}
		if pkg == nil || len(pkg.Errors) > 0 {
			return nil, fmt.Errorf("genroot: could not find package %q", fct.Path)
		}

		var (
			scope = pkg.Types.Scope()
		)
		obj := scope.Lookup(fct.Def)
		if obj == nil {
			return nil, fmt.Errorf("genroot: could not find %s in package %q", fct.Def, fct.Path)
		}
		ft, ok := obj.(*types.Func)
		if !ok {
			return nil, fmt.Errorf("genroot: object %s in package %q is not a func (%T)", fct.Def, fct.Path, obj)
		}
		f = ft.Type().Underlying().(*types.Signature)
	}

	name := fct.Name
	if name == "" {
		switch fct.Path {
		case "":
			name = "FuncFormula"
		default:
			name = fct.Def + "Formula"
		}
	}

	gen := &rfuncGen{w: w, f: f, name: name}
	switch fct.Pkg {
	case "go-hep.org/x/hep/groot/rtree/rfunc":
		// no-op.
	default:
		gen.pkg = "rfunc."
	}

	return gen, nil
}

func (gen *rfuncGen) Generate() error {
	fct := rfuncTypeFrom(gen.name, gen.f)
	tmpl := template.Must(template.New("rfunc").Funcs(
		template.FuncMap{
			"Pkg": func() string {
				return gen.pkg
			},
		},
	).Parse(rfuncCodeTmpl))
	err := tmpl.Execute(gen.w, fct)
	if err != nil {
		return fmt.Errorf("genroot: could not execute template for %q: %w",
			fct.Name, err,
		)
	}
	return nil
}

func (gen *rfuncGen) GenerateTest(w io.Writer) error {
	fct := rfuncTypeFrom(gen.name, gen.f)
	tmpl := template.Must(template.New("rfunc").Funcs(
		template.FuncMap{
			"Pkg":  func() string { return gen.pkg },
			"Out0": func() string { return fct.Out[0] },
		},
	).Parse(rfuncTestTmpl))
	err := tmpl.Execute(w, fct)
	if err != nil {
		return fmt.Errorf("genroot: could not execute template for %q: %w",
			fct.Name, err,
		)
	}
	return nil
}

func parseExpr(x string) (*types.Signature, error) {
	expr, err := parser.ParseExpr(x)
	if err != nil {
		return nil, fmt.Errorf("genroot: could not parse %q: %w", x, err)
	}
	switch expr := expr.(type) {
	case *ast.FuncType:
		var (
			pos token.Pos
			pkg *types.Package
			par *types.Tuple
			res *types.Tuple
			sig *types.Signature
		)
		typeFor := func(typ ast.Expr) types.Type {
			switch typ := typ.(type) {
			case *ast.Ident:
				t, ok := astTypesToGoTypes[typ.Name]
				if !ok {
					panic(fmt.Errorf("unknown ast.Ident type name %q", typ.Name))
				}
				return t
			default:
				panic(fmt.Errorf("unhandled ast.Expr: %#v (%T)", typ, typ))
			}
		}
		mk := func(lst *ast.FieldList) *types.Tuple {
			vs := make([]*types.Var, lst.NumFields())
			ns := make([]string, 0, len(vs))
			ts := make([]ast.Expr, 0, len(vs))
			for i, vs := range lst.List {
				switch len(vs.Names) {
				case 0:
					ns = append(ns, fmt.Sprintf("arg%02d", i))
					ts = append(ts, vs.Type)
				default:
					for _, n := range vs.Names {
						ts = append(ts, vs.Type)
						ns = append(ns, n.Name)
					}
				}
			}
			for i, v := range ns {
				vs[i] = types.NewVar(pos, pkg, v, typeFor(ts[i]))
			}
			return types.NewTuple(vs...)
		}
		par = mk(expr.Params)
		res = mk(expr.Results)
		sig = types.NewSignature(nil, par, res, false)
		return sig, nil
	default:
		panic(fmt.Errorf("error: expr=%T", expr))
	}
}

var (
	astTypesToGoTypes = map[string]types.Type{
		"bool":    types.Typ[types.Bool],
		"byte":    types.Typ[types.Byte],
		"uint8":   types.Typ[types.Uint8],
		"uint16":  types.Typ[types.Uint16],
		"uint32":  types.Typ[types.Uint32],
		"uint64":  types.Typ[types.Uint64],
		"int8":    types.Typ[types.Int8],
		"int16":   types.Typ[types.Int16],
		"int32":   types.Typ[types.Int32],
		"int64":   types.Typ[types.Int64],
		"uint":    types.Typ[types.Uint],
		"int":     types.Typ[types.Int],
		"float32": types.Typ[types.Float32],
		"float64": types.Typ[types.Float64],
		"string":  types.Typ[types.String],
	}
)

type rfuncType struct {
	Name string
	In   []string
	Out  []string
}

func rfuncTypeFrom(name string, sig *types.Signature) rfuncType {
	var (
		ps  = sig.Params()
		rs  = sig.Results()
		fct = rfuncType{
			Name: name,
			In:   make([]string, ps.Len()),
			Out:  make([]string, rs.Len()),
		}
	)

	for i := range fct.In {
		fct.In[i] = ps.At(i).Type().String()
	}

	for i := range fct.Out {
		fct.Out[i] = rs.At(i).Type().String()
	}

	return fct
}

func (f rfuncType) NumIn() int   { return len(f.In) }
func (f rfuncType) NumOut() int  { return len(f.Out) }
func (f rfuncType) Type() string { return f.Name }

func (f rfuncType) Func() string {
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

func (f rfuncType) Return() string {
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

func (f rfuncType) TestFunc() string {
	switch f.Out[0] {
	case "string":
		return `"42"`
	case "bool":
		return "true"
	default:
		return "42"
	}
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

// New{{.Type}} return a new formula, from the provided function.
func New{{.Type}}(rvars []string, fct {{.Func}}) *{{.Type}} {
	return &{{.Type}}{
{{- if gt .NumIn 0}}
		rvars: rvars,
{{- end}}
		fct: fct,
	}
}

{{if gt .NumIn 0}}
// RVars implements rfunc.Formula
func (f *{{.Type}}) RVars() []string { return f.rvars }
{{else}}
// RVars implements rfunc.Formula
func (f *{{.Type}}) RVars() []string { return nil }
{{end}}

// Bind implements rfunc.Formula
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

// Func implements rfunc.Formula
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
	_ {{Pkg}}Formula = (*{{.Type}})(nil)
)
`

const rfuncTestTmpl = `func Test{{.Type}}(t *testing.T) {
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

	form := New{{.Type}}(rvars, fct)

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

	err := form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func () {{.Return}})()
	if got, want := got, {{Out0}}({{.TestFunc}}); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
`
