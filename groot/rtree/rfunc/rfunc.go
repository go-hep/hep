// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rfunc provides types and funcs to implement user-provided formulae
// evaluated on data exposed by ROOT trees.
package rfunc // import "go-hep.org/x/hep/groot/rtree/rfunc"

//go:generate go run ./gen-rfuncs.go

import (
	"fmt"
	"reflect"
)

// Formula is the interface that describes the protocol between a user
// provided function (that evaluates a value based on some data in a ROOT
// tree) and the rtree.Reader (that presents data from a ROOT tree.)
type Formula interface {
	// RVars returns the names of the leaves that this formula needs.
	// The returned slice must contain the names in the same order than the
	// user formula function's arguments.
	RVars() []string

	// Bind provides the arguments to the user function.
	// ptrs is a slice of pointers to the rtree.ReadVars, in the same order
	// than requested by RVars.
	Bind(ptrs []interface{}) error

	// Func returns the user function closing on the bound pointer-to-arguments
	// and returning the expected evaluated value.
	Func() interface{}
}

// NewGenericFormula returns a new formula from the provided list of needed
// tree variables and the provided user function.
// NewGenericFormula uses reflect to bind read-vars and the generic function.
func NewGenericFormula(rvars []string, fct interface{}) (Formula, error) {
	return newGenericFormula(rvars, fct)
}

type genericFormula struct {
	names []string
	fct   interface{}

	ptrs []reflect.Value
	args []reflect.Value
	out  []reflect.Value

	rfct reflect.Value // formula-created function to eval read-vars
	ufct reflect.Value // user-provided function
}

func newGenericFormula(names []string, fct interface{}) (*genericFormula, error) {
	rv := reflect.ValueOf(fct)
	if rv.Kind() != reflect.Func {
		return nil, fmt.Errorf("rfunc: formula expects a func")
	}

	if len(names) != rv.Type().NumIn() {
		return nil, fmt.Errorf("rfunc: num-branches/func-arity mismatch")
	}

	if rv.Type().NumOut() != 1 {
		// FIXME(sbinet): allow any kind of function?
		return nil, fmt.Errorf("rfunc: invalid number of return values")
	}

	args := make([]reflect.Value, len(names))
	ptrs := make([]reflect.Value, len(names))
	for i := range args {
		args[i] = reflect.New(rv.Type().In(i)).Elem()
	}

	gen := &genericFormula{
		names: names,
		fct:   fct,

		ptrs: ptrs,
		args: args,
		ufct: rv,
	}

	gen.rfct = reflect.MakeFunc(
		reflect.FuncOf(nil, []reflect.Type{rv.Type().Out(0)}, false),
		func(in []reflect.Value) []reflect.Value {
			gen.eval()
			return gen.out
		},
	)

	return gen, nil
}

func (f *genericFormula) RVars() []string { return f.names }
func (f *genericFormula) Bind(args []interface{}) error {
	if got, want := len(args), len(f.ptrs); got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}

	for i := range args {
		var (
			got  = reflect.TypeOf(args[i]).Elem()
			want = f.args[i].Type()
		)
		if got != want {
			return fmt.Errorf(
				"rfunc: argument type %d (name=%s) mismatch: got=%T, want=%T",
				i, f.names[i],
				reflect.New(got).Elem().Interface(),
				reflect.New(want).Elem().Interface(),
			)
		}
		f.ptrs[i] = reflect.ValueOf(args[i])
		f.args[i] = reflect.New(f.ptrs[i].Type().Elem()).Elem()
	}

	return nil
}

func (f *genericFormula) eval() {
	for i := range f.ptrs {
		f.args[i].Set(f.ptrs[i].Elem())
	}
	f.out = f.ufct.Call(f.args)
}

func (f *genericFormula) Func() interface{} {
	return f.rfct.Interface()
}

var (
	_ Formula = (*genericFormula)(nil)
)
