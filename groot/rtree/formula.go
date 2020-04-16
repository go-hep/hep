// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"go/ast"
	"go/parser"
	"reflect"
	"sort"
	"strings"

	"github.com/containous/yaegi/interp"
	"github.com/containous/yaegi/stdlib"
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
	fct  func() interface{}
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

	fmt.Fprintf(prog, "import %q\n", pkg)
	fmt.Fprintf(prog, "func _groot_rtree_func_eval() interface{} {\n")

	names := make([]string, 0, len(needed))
	for k := range needed {
		names = append(names, k)
	}
	sort.Strings(names)

	for _, n := range names {
		rvar := needed[n]
		name := "Var_" + rvar.Name
		uses[pkg][name] = reflect.ValueOf(rvar.Value)
		fmt.Fprintf(prog, "\t%s := *%s.%s // %T\n", rvar.Name, pkg, name, rvar.Value)
	}

	eval.Use(stdlib.Symbols)
	eval.Use(uses)

	fmt.Fprintf(prog,
		"\t_groot_return := %s\n\treturn &_groot_return\n}",
		expr,
	)

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
		fct:  f.Interface().(func() interface{}),
	}

	return form, nil
}

func (form *Formula) Eval() interface{} {
	return reflect.ValueOf(form.fct()).Elem().Interface()
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
