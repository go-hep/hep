// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"
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

	for _, name := range imports {
		if _, ok := stdlib.Symbols[name]; !ok {
			return Formula{}, fmt.Errorf("rtree: no known stdlib import for %q", name)
		}
		fmt.Fprintf(prog, "import %q\n", name)
	}

	fmt.Fprintf(prog, "import %q\n", pkg)
	fmt.Fprintf(prog, "func _groot_rtree_func_eval() interface{} {\n")

	for _, rvar := range r.rvars {
		name := "Var_" + rvar.Name
		uses[pkg][name] = reflect.ValueOf(rvar.Value)
		// FIXME(sbinet): only load rvars that are actually used.
		fmt.Fprintf(prog, "\t%s := *%s.%s // %T\n", rvar.Name, pkg, name, rvar.Value)
	}

	eval.Use(stdlib.Symbols)
	eval.Use(uses)

	fmt.Fprintf(prog,
		"\t_groot_return := %s\n\treturn &_groot_return\n}",
		expr,
	)

	_, err := eval.Eval(prog.String())
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
