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

func newFormula(r *Reader, expr string, imports []string) (*Formula, error) {
	var (
		ir   = fast.New()
		code = new(strings.Builder)
	)

	idents, err := identsFromExpr(expr)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not parse expression: %w", err)
	}

	needed := formulaAutoLoad(r, idents)

	for _, name := range imports {
		_ = ir.ImportPackage("", name)
		//if err != nil {
		//	return Formula{}, fmt.Errorf("rtree: could not import %q into formula interpreter", name)
		//}
	}

	ret, err := formulaAnalyze(needed, imports, expr)
	if err != nil {
		return nil, fmt.Errorf("rtree: could not analyze formula type: %w", err)
	}

	fmt.Fprintf(code, "func _groot_eval() {\n")

	for _, rvar := range needed {
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

	form := &Formula{
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

func identsFromExpr(s string) ([]string, error) {
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

	idents := make([]string, 0, len(set))
	for k := range set {
		idents = append(idents, k)
	}
	sort.Strings(idents)

	return idents, nil
}

var (
	_ formula = (*Formula)(nil)
)

type FormulaFunc struct {
	rvars []*ReadVar
	args  []reflect.Value
	out   []reflect.Value
	rfct  reflect.Value // formula-created function to eval read-vars
	ufct  reflect.Value // user-provided function
}

func newFormulaFunc(r *Reader, branches []string, fct interface{}) (*FormulaFunc, error) {
	rv := reflect.ValueOf(fct)
	if rv.Kind() != reflect.Func {
		return nil, fmt.Errorf("rtree: FormulaFunc expects a func")
	}

	if len(branches) != rv.Type().NumIn() {
		return nil, fmt.Errorf("rtree: num-branches/func-arity mismatch")
	}

	if rv.Type().NumOut() != 1 {
		// FIXME(sbinet): allow any kind of function?
		return nil, fmt.Errorf("rtree: invalid number of return values")
	}

	rvars := formulaAutoLoad(r, branches)
	if len(rvars) != len(branches) {
		return nil, fmt.Errorf("rtree: could not find all needed ReadVars")
	}

	for i, rvar := range rvars {
		btyp := reflect.TypeOf(rvar.Value).Elem()
		atyp := rv.Type().In(i)
		if btyp != atyp {
			return nil, fmt.Errorf(
				"rtree: argument type %d mismatch: func=%T, read-var[%s]=%T",
				i,
				reflect.New(atyp).Elem().Interface(),
				rvar.Name,
				reflect.New(btyp).Elem().Interface(),
			)
		}
	}

	form := &FormulaFunc{
		rvars: rvars,
		args:  make([]reflect.Value, len(rvars)),
		ufct:  rv,
	}

	for i := range form.rvars {
		form.args[i] = reflect.New(rv.Type().In(i)).Elem()
	}

	rfct := reflect.MakeFunc(
		reflect.FuncOf(nil, []reflect.Type{rv.Type().Out(0)}, false),
		func(in []reflect.Value) []reflect.Value {
			form.eval()
			return form.out
		},
	)
	form.rfct = rfct

	return form, nil
}

func (form *FormulaFunc) eval() {
	for i, rvar := range form.rvars {
		form.args[i].Set(reflect.ValueOf(rvar.Value).Elem())
	}
	form.out = form.ufct.Call(form.args)
}

func (form *FormulaFunc) Eval() interface{} {
	form.eval()
	return form.out[0].Interface()
}

func (form *FormulaFunc) Func() interface{} {
	return form.rfct.Interface()
}

var (
	_ formula = (*FormulaFunc)(nil)
)

func formulaAnalyze(rvars []*ReadVar, imports []string, expr string) (reflect.Type, error) {
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

func formulaAutoLoad(r *Reader, idents []string) []*ReadVar {
	var (
		loaded = make(map[string]*ReadVar, len(r.rvars))
		needed = make([]*ReadVar, 0, len(idents))
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
	for _, name := range idents {
		rvar, ok := all[name]
		if !ok {
			continue
		}
		if _, ok := loaded[name]; !ok {
			r.rvars = append(r.rvars, *rvar)
			rvar = &r.rvars[len(r.rvars)-1]
			loaded[name] = rvar
		}
		needed = append(needed, rvar)
	}

	return needed
}
