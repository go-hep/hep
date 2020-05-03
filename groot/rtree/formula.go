// Copyright ©2020 The go-hep Authors. All rights reserved.
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

	"github.com/cosmos72/gomacro/fast"
)

// Formula is a mathematical formula bound to variables (branches) of
// a given ROOT tree.
//
// Formulae are attached to a rtree.Reader.
type Formula struct {
	ir   *fast.Interp
	expr *fast.Expr

	fct  interface{}
	recv reflect.Value
}

func newFormula(r *Reader, expr string, imports []string) (Formula, error) {
	var (
		ir   = fast.New()
		code = new(strings.Builder)
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
		_ = ir.ImportPackage("", name)
		//if err != nil {
		//	return Formula{}, fmt.Errorf("rtree: could not import %q into formula interpreter", name)
		//}
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

	fmt.Fprintf(code, "func _groot_eval() {\n")

	for _, n := range names {
		rvar := needed[n]
		name := "_groot_var_" + rvar.Name
		rval := reflect.ValueOf(rvar.Value)
		rtyp := ir.TypeOf(rval.Interface())
		ir.DeclVar(name, rtyp, rval.Interface())
		fmt.Fprintf(code, "\t%s := *%s\n", rvar.Name, name)
	}
	recv := reflect.New(ret)
	rtyp := ir.TypeOf(recv.Interface())
	ir.DeclVar("_groot_recv", rtyp, recv.Interface())

	fmt.Fprintf(code, "\t*_groot_recv = %s\n}\n", expr)
	fmt.Fprintf(code, "_groot_eval()\n")

	prog := ir.Compile(code.String())

	results := make([]reflect.Value, 1)
	otypes := []reflect.Type{ret}
	sig := reflect.FuncOf(nil, otypes, false)

	form := Formula{
		ir:   ir,
		expr: prog,
		recv: recv,
	}

	f := reflect.MakeFunc(
		sig,
		func(args []reflect.Value) []reflect.Value {
			form.eval()
			results[0] = form.recv.Elem()
			return results
		},
	)

	form.fct = f.Interface()

	return form, nil
}

func (form *Formula) Func() interface{} {
	return form.fct
}

func (form *Formula) Eval() interface{} {
	form.eval()
	return form.recv.Elem().Interface()
}

func (form *Formula) eval() {
	_, _ = form.ir.RunExpr(form.expr)
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

func formulaAnalyze(rvars map[string]*ReadVar, imports []string, expr string) (reflect.Type, error) {
	ir := fast.New()
	for _, name := range imports {
		_ = ir.ImportPackage("", name)
	}
	code := new(strings.Builder)
	for _, rvar := range rvars {
		name := "_groot_var_" + rvar.Name
		rval := reflect.ValueOf(rvar.Value)
		rtyp := ir.TypeOf(rval.Interface())
		ir.DeclVar(name, rtyp, rval.Interface())
		fmt.Fprintf(code, "%s := *%s // %T\n", rvar.Name, name, rvar.Value)
	}
	fmt.Fprintf(code, "_groot_recv := %s\n_groot_recv\n", expr)

	var (
		prog *fast.Expr
		err  error
	)
	func() {
		defer func() {
			e := recover()
			if e == nil {
				return
			}
			err = e.(error)
		}()
		prog = ir.Compile(code.String())
	}()
	if err != nil {
		return nil, fmt.Errorf("rtree: could not analyze formula: %w", err)
	}

	return prog.DefaultType().ReflectType(), nil
}
