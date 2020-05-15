// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/root"
)

type FormulaFunc struct {
	rvars []reflect.Value
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

	rvars, missing := formulaAutoLoad(r, branches)
	if len(rvars) != len(branches) {
		return nil, fmt.Errorf("rtree: could not find all needed ReadVars (missing: %v)", missing)
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
		rvars: make([]reflect.Value, len(rvars)),
		args:  make([]reflect.Value, len(rvars)),
		ufct:  rv,
	}

	for i := range form.rvars {
		form.args[i] = reflect.New(rv.Type().In(i)).Elem()
		form.rvars[i] = reflect.ValueOf(rvars[i].Value)
	}

	var rfct reflect.Value
	switch reflect.New(rv.Type().Out(0)).Elem().Interface().(type) {
	case bool:
		ufct := func() bool {
			form.eval()
			return form.out[0].Interface().(bool)
		}
		rfct = reflect.ValueOf(ufct)
	case uint8:
		ufct := func() uint8 {
			form.eval()
			return form.out[0].Interface().(uint8)
		}
		rfct = reflect.ValueOf(ufct)
	case uint16:
		ufct := func() uint16 {
			form.eval()
			return form.out[0].Interface().(uint16)
		}
		rfct = reflect.ValueOf(ufct)
	case uint32:
		ufct := func() uint32 {
			form.eval()
			return form.out[0].Interface().(uint32)
		}
		rfct = reflect.ValueOf(ufct)
	case uint64:
		ufct := func() uint64 {
			form.eval()
			return form.out[0].Interface().(uint64)
		}
		rfct = reflect.ValueOf(ufct)
	case int8:
		ufct := func() int8 {
			form.eval()
			return form.out[0].Interface().(int8)
		}
		rfct = reflect.ValueOf(ufct)
	case int16:
		ufct := func() int16 {
			form.eval()
			return form.out[0].Interface().(int16)
		}
		rfct = reflect.ValueOf(ufct)
	case int32:
		ufct := func() int32 {
			form.eval()
			return form.out[0].Interface().(int32)
		}
		rfct = reflect.ValueOf(ufct)
	case int64:
		ufct := func() int64 {
			form.eval()
			return form.out[0].Interface().(int64)
		}
		rfct = reflect.ValueOf(ufct)
	case string:
		ufct := func() string {
			form.eval()
			return form.out[0].Interface().(string)
		}
		rfct = reflect.ValueOf(ufct)
	case root.Float16:
		ufct := func() root.Float16 {
			form.eval()
			return form.out[0].Interface().(root.Float16)
		}
		rfct = reflect.ValueOf(ufct)
	case float32:
		ufct := func() float32 {
			form.eval()
			return form.out[0].Interface().(float32)
		}
		rfct = reflect.ValueOf(ufct)
	case root.Double32:
		ufct := func() root.Double32 {
			form.eval()
			return form.out[0].Interface().(root.Double32)
		}
		rfct = reflect.ValueOf(ufct)
	case float64:
		ufct := func() float64 {
			form.eval()
			return form.out[0].Float()
		}
		rfct = reflect.ValueOf(ufct)
	default:
		rfct = reflect.MakeFunc(
			reflect.FuncOf(nil, []reflect.Type{rv.Type().Out(0)}, false),
			func(in []reflect.Value) []reflect.Value {
				form.eval()
				return form.out
			},
		)
	}
	form.rfct = rfct

	return form, nil
}

func (form *FormulaFunc) eval() {
	for i, rvar := range form.rvars {
		form.args[i].Set(rvar.Elem())
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

func formulaAutoLoad(r *Reader, idents []string) ([]*ReadVar, []string) {
	var (
		loaded  = make(map[string]*ReadVar, len(r.r.rvars))
		needed  = make([]*ReadVar, 0, len(idents))
		rvars   = NewReadVars(r.r.tree)
		all     = make(map[string]*ReadVar, len(rvars))
		missing []string
	)

	for i := range r.r.rvars {
		rvar := &r.r.rvars[i]
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
			missing = append(missing, name)
			continue
		}
		if _, ok := loaded[name]; !ok {
			r.r.rvars = append(r.r.rvars, *rvar)
			rvar = &r.r.rvars[len(r.r.rvars)-1]
			loaded[name] = rvar
		}
		needed = append(needed, rvar)
	}

	return needed, missing
}
