// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
	"reflect"
)

// funcAr00F64 implements rfunc.Formula
type funcAr00F64 struct {
	fct func() float64
}

func newFuncAr00F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr00F64{
		fct: fct.(func() float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr00F64{}.fct)] = newFuncAr00F64
}

func (f *funcAr00F64) RVars() []string { return nil }

func (f *funcAr00F64) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

func (f *funcAr00F64) Func() interface{} {
	return func() float64 {
		return f.fct()
	}
}

var (
	_ Formula = (*funcAr00F64)(nil)
)

// funcAr01F64 implements rfunc.Formula
type funcAr01F64 struct {
	rvars []string
	arg0  *float64
	fct   func(arg00 float64) float64
}

func newFuncAr01F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr01F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr01F64{}.fct)] = newFuncAr01F64
}

func (f *funcAr01F64) RVars() []string { return f.rvars }

func (f *funcAr01F64) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

func (f *funcAr01F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*funcAr01F64)(nil)
)

// funcAr02F64 implements rfunc.Formula
type funcAr02F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	fct   func(arg00 float64, arg01 float64) float64
}

func newFuncAr02F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr02F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr02F64{}.fct)] = newFuncAr02F64
}

func (f *funcAr02F64) RVars() []string { return f.rvars }

func (f *funcAr02F64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

func (f *funcAr02F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*funcAr02F64)(nil)
)

// funcAr03F64 implements rfunc.Formula
type funcAr03F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64) float64
}

func newFuncAr03F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr03F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr03F64{}.fct)] = newFuncAr03F64
}

func (f *funcAr03F64) RVars() []string { return f.rvars }

func (f *funcAr03F64) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

func (f *funcAr03F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*funcAr03F64)(nil)
)

// funcAr04F64 implements rfunc.Formula
type funcAr04F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64) float64
}

func newFuncAr04F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr04F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr04F64{}.fct)] = newFuncAr04F64
}

func (f *funcAr04F64) RVars() []string { return f.rvars }

func (f *funcAr04F64) Bind(args []interface{}) error {
	if got, want := len(args), 4; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	return nil
}

func (f *funcAr04F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
		)
	}
}

var (
	_ Formula = (*funcAr04F64)(nil)
)

// funcAr05F64 implements rfunc.Formula
type funcAr05F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64) float64
}

func newFuncAr05F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr05F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr05F64{}.fct)] = newFuncAr05F64
}

func (f *funcAr05F64) RVars() []string { return f.rvars }

func (f *funcAr05F64) Bind(args []interface{}) error {
	if got, want := len(args), 5; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	return nil
}

func (f *funcAr05F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
			*f.arg4,
		)
	}
}

var (
	_ Formula = (*funcAr05F64)(nil)
)

// funcAr06F64 implements rfunc.Formula
type funcAr06F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64) float64
}

func newFuncAr06F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr06F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr06F64{}.fct)] = newFuncAr06F64
}

func (f *funcAr06F64) RVars() []string { return f.rvars }

func (f *funcAr06F64) Bind(args []interface{}) error {
	if got, want := len(args), 6; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	return nil
}

func (f *funcAr06F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
			*f.arg4,
			*f.arg5,
		)
	}
}

var (
	_ Formula = (*funcAr06F64)(nil)
)

// funcAr07F64 implements rfunc.Formula
type funcAr07F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	arg6  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64) float64
}

func newFuncAr07F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr07F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr07F64{}.fct)] = newFuncAr07F64
}

func (f *funcAr07F64) RVars() []string { return f.rvars }

func (f *funcAr07F64) Bind(args []interface{}) error {
	if got, want := len(args), 7; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	return nil
}

func (f *funcAr07F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
			*f.arg4,
			*f.arg5,
			*f.arg6,
		)
	}
}

var (
	_ Formula = (*funcAr07F64)(nil)
)

// funcAr08F64 implements rfunc.Formula
type funcAr08F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	arg6  *float64
	arg7  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64) float64
}

func newFuncAr08F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr08F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr08F64{}.fct)] = newFuncAr08F64
}

func (f *funcAr08F64) RVars() []string { return f.rvars }

func (f *funcAr08F64) Bind(args []interface{}) error {
	if got, want := len(args), 8; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	return nil
}

func (f *funcAr08F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
			*f.arg4,
			*f.arg5,
			*f.arg6,
			*f.arg7,
		)
	}
}

var (
	_ Formula = (*funcAr08F64)(nil)
)

// funcAr09F64 implements rfunc.Formula
type funcAr09F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	arg6  *float64
	arg7  *float64
	arg8  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64) float64
}

func newFuncAr09F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr09F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr09F64{}.fct)] = newFuncAr09F64
}

func (f *funcAr09F64) RVars() []string { return f.rvars }

func (f *funcAr09F64) Bind(args []interface{}) error {
	if got, want := len(args), 9; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	return nil
}

func (f *funcAr09F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
			*f.arg4,
			*f.arg5,
			*f.arg6,
			*f.arg7,
			*f.arg8,
		)
	}
}

var (
	_ Formula = (*funcAr09F64)(nil)
)

// funcAr10F64 implements rfunc.Formula
type funcAr10F64 struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	arg6  *float64
	arg7  *float64
	arg8  *float64
	arg9  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64, arg09 float64) float64
}

func newFuncAr10F64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr10F64{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64, arg09 float64) float64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr10F64{}.fct)] = newFuncAr10F64
}

func (f *funcAr10F64) RVars() []string { return f.rvars }

func (f *funcAr10F64) Bind(args []interface{}) error {
	if got, want := len(args), 10; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	{
		ptr, ok := args[9].(*float64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 9 (name=%s) mismatch: got=%T, want=*float64",
				f.rvars[9], args[9],
			)
		}
		f.arg9 = ptr
	}
	return nil
}

func (f *funcAr10F64) Func() interface{} {
	return func() float64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
			*f.arg4,
			*f.arg5,
			*f.arg6,
			*f.arg7,
			*f.arg8,
			*f.arg9,
		)
	}
}

var (
	_ Formula = (*funcAr10F64)(nil)
)
