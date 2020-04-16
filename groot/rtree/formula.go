// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"reflect"
	"sort"
	"strings"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
	"go-hep.org/x/hep/groot/root"
)

// Formula is a mathematical formula bound to variables (branches) of
// a given ROOT tree.
//
// Formulae are attached to a rtree.Reader.
type Formula struct {
	r    *Reader
	expr string
	prog string
	eval *interp.Interpreter

	fct  interface{}
	rfun reflect.Value
}

func newFormula(r *Reader, expr string, imports []string) (Formula, error) {
	var (
		eval = interp.New(interp.Options{})
		pkg  = "groot_rtree"
		uses = interp.Exports{
			pkg: make(map[string]reflect.Value),
		}
		prog = new(strings.Builder)
	)

	idents, err := identsFromExpr(expr)
	if err != nil {
		return Formula{}, fmt.Errorf("rtree: could not parse expression: %w", err)
	}

	var (
		loaded = make(map[string]*ReadVar, len(r.rvars))
		needed = make(map[string]*ReadVar, len(idents))
		rvars  = NewReadVars(r.t)
		all    = make(map[string]*ReadVar, len(rvars))
	)

	for i := range r.rvars {
		rvar := &r.rvars[i]
		loaded[rvar.Name] = rvar
		all[rvar.Name] = rvar
	}
	for i := range rvars {
		rvar := &rvars[i]
		if _, ok := all[rvar.Name]; ok {
			continue
		}
		all[rvar.Name] = rvar
	}
	for k := range idents {
		rvar, ok := all[k]
		if !ok {
			continue
		}
		if _, ok := loaded[k]; !ok {
			r.rvars = append(r.rvars, *rvar)
			rvar = &r.rvars[len(r.rvars)-1]
			loaded[k] = rvar
		}
		needed[k] = rvar
	}

	for _, name := range imports {
		if _, ok := stdlib.Symbols[name]; !ok {
			return Formula{}, fmt.Errorf("rtree: no known stdlib import for %q", name)
		}
		fmt.Fprintf(prog, "import %q\n", name)
	}

	names := make([]string, 0, len(needed))
	for k := range needed {
		names = append(names, k)
	}
	sort.Strings(names)

	ret, err := formulaAnalyze(needed, imports, expr)
	if err != nil {
		return Formula{}, fmt.Errorf("rtree: could not analyze formula type: %w", err)
	}

	// FIXME(sbinet): instead of returning the result of evaluating
	// the user expression by value, we could perhaps pass a pointer
	// as argument, storing the result in there.
	// we'd then zap everything to unsafe.Pointer so signatures match.
	// (that's to support, e.g., root.Double32, which isn't a float64)
	//
	// Alternatively, we could define the "return" value inside the fake
	// "groot_rtree" package, say, groot_rtree.Out.
	// we'd need to generate a reflect.Func wrapping it.

	fmt.Fprintf(prog, "import %q\n", pkg)
	fmt.Fprintf(prog, "func _groot_rtree_func_eval() %v {\n", ret)

	for _, n := range names {
		rvar := needed[n]
		name := "Var_" + rvar.Name
		uses[pkg][name] = reflect.ValueOf(rvar.Value)
		fmt.Fprintf(prog, "\t%s := *%s.%s // %T\n", rvar.Name, pkg, name, rvar.Value)
	}

	eval.Use(stdlib.Symbols)
	eval.Use(uses)

	fmt.Fprintf(prog, "\treturn %s(%s)\n}", ret, expr)

	_, err = eval.Eval(prog.String())
	if err != nil {
		return Formula{}, fmt.Errorf("rtree: could not define formula eval-func: %w", err)
	}

	f, err := eval.Eval("_groot_rtree_func_eval")
	if err != nil {
		return Formula{}, fmt.Errorf("rtree: could not retrieve formula eval-func: %w", err)
	}

	form := Formula{
		r:    r,
		expr: expr,
		prog: prog.String(),
		eval: eval,
		fct:  f.Interface(),
		rfun: f,
	}

	return form, nil
}

func (form *Formula) Func() interface{} { return form.fct }

func (form *Formula) Eval() interface{} {
	return form.rfun.Call(nil)[0].Interface()
}

func identsFromExpr(s string) (map[string]struct{}, error) {
	set := make(map[string]struct{})
	expr, err := parser.ParseExpr(s)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not parse formula %q: %w", s, err)
	}

	ast.Inspect(expr, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			set[node.Name] = struct{}{}
		case *ast.SelectorExpr:
			// e.g.: math.Pi or math.Exp(x)
			return false
		}
		return true
	})

	return set, nil
}

func formulaAnalyze(rvars map[string]*ReadVar, imports []string, expr string) (types.Type, error) {
	prog := new(strings.Builder)
	fmt.Fprintf(prog, "package main\n")
	fmt.Fprintf(prog, "\nimport (\n")
	for _, name := range imports {
		fmt.Fprintf(prog, "\t%q\n", name)
	}
	fmt.Fprintf(prog, ")\n")

	fmt.Fprintf(prog, "\nfunc main() {}\n")
	fmt.Fprintf(prog, "\nfunc _() {\n")
	for _, rvar := range rvars {
		rv := reflect.ValueOf(rvar.Value).Elem().Interface()
		// hack for root.Double32/root.Float16
		switch rv.(type) {
		case root.Double32:
			rv = float64(0)
		case root.Float16:
			rv = float32(0)
		}
		fmt.Fprintf(prog, "\tvar %s %T\n", rvar.Name, rv)
	}
	fmt.Fprintf(prog, "\n\tvar _groot_out = %s\n", expr)
	fmt.Fprintf(prog, "\t_ = _groot_out\n}\n")

	var (
		fset  = token.NewFileSet()
		input = prog.String()
	)

	f, err := parser.ParseFile(fset, "groot_rtree_formula.go", input, 0)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not create type analysis code: %w", err)
	}

	conf := types.Config{Importer: importer.Default()}
	info := &types.Info{Types: make(map[ast.Expr]types.TypeAndValue)}
	if _, err := conf.Check("cmd/groot-rtree-formula-type", fset, []*ast.File{f}, info); err != nil {
		return nil, fmt.Errorf("rtree: could not type-check formula analysis code: %w", err)
	}

	var typ types.Type
	ast.Inspect(f, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.Ident:
			if node.Name != "_groot_out" {
				return true
			}
			tv, ok := info.Types[node]
			if !ok {
				return true
			}
			typ = tv.Type
			return false
		}
		return true
	})

	if typ == nil {
		return nil, fmt.Errorf("rtree: could not find type of expression")
	}

	return typ.Underlying(), nil
}

//func reflectTypeFromGoTypes(typ types.Type) reflect.Type {
//	ut := typ.Underlying()
//	switch typ := ut.(type) {
//	case *types.Basic:
//		return map[types.BasicKind]reflect.Type{
//			types.Bool:       reflect.TypeOf(true),
//			types.Int:        reflect.TypeOf(int(0)),
//			types.Int8:       reflect.TypeOf(int8(0)),
//			types.Int16:      reflect.TypeOf(int16(0)),
//			types.Int32:      reflect.TypeOf(int32(0)),
//			types.Int64:      reflect.TypeOf(int64(0)),
//			types.Uint:       reflect.TypeOf(uint(0)),
//			types.Uint8:      reflect.TypeOf(uint8(0)),
//			types.Uint16:     reflect.TypeOf(uint16(0)),
//			types.Uint32:     reflect.TypeOf(uint32(0)),
//			types.Uint64:     reflect.TypeOf(uint64(0)),
//			types.Float32:    reflect.TypeOf(float32(0)),
//			types.Float64:    reflect.TypeOf(float64(0)),
//			types.Complex64:  reflect.TypeOf(complex(float32(0), float32(0))),
//			types.Complex128: reflect.TypeOf(complex(float64(0), float64(0))),
//			types.String:     reflect.TypeOf(""),
//		}[typ.Kind()]
//	case *types.Slice:
//		et := reflectTypeFromGoTypes(typ.Elem())
//		return reflect.SliceOf(et)
//	case *types.Array:
//		et := reflectTypeFromGoTypes(typ.Elem())
//		return reflect.ArrayOf(int(typ.Len()), et)
//	}
//	panic("not implemented")
//}
