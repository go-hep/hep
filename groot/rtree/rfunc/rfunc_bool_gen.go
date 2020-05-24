// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
	"reflect"
)

// funcAr00Bool implements rfunc.Formula
type funcAr00Bool struct {
	fct func() bool
}

func newFuncAr00Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr00Bool{
		fct: fct.(func() bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr00Bool{}.fct)] = newFuncAr00Bool
}

func (f *funcAr00Bool) RVars() []string { return nil }

func (f *funcAr00Bool) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

func (f *funcAr00Bool) Func() interface{} {
	return func() bool {
		return f.fct()
	}
}

var (
	_ Formula = (*funcAr00Bool)(nil)
)

// funcAr01Bool implements rfunc.Formula
type funcAr01Bool struct {
	rvars []string
	arg0  *float64
	fct   func(arg00 float64) bool
}

func newFuncAr01Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr01Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr01Bool{}.fct)] = newFuncAr01Bool
}

func (f *funcAr01Bool) RVars() []string { return f.rvars }

func (f *funcAr01Bool) Bind(args []interface{}) error {
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

func (f *funcAr01Bool) Func() interface{} {
	return func() bool {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*funcAr01Bool)(nil)
)

// funcAr02Bool implements rfunc.Formula
type funcAr02Bool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	fct   func(arg00 float64, arg01 float64) bool
}

func newFuncAr02Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr02Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr02Bool{}.fct)] = newFuncAr02Bool
}

func (f *funcAr02Bool) RVars() []string { return f.rvars }

func (f *funcAr02Bool) Bind(args []interface{}) error {
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

func (f *funcAr02Bool) Func() interface{} {
	return func() bool {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*funcAr02Bool)(nil)
)

// funcAr03Bool implements rfunc.Formula
type funcAr03Bool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64) bool
}

func newFuncAr03Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr03Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr03Bool{}.fct)] = newFuncAr03Bool
}

func (f *funcAr03Bool) RVars() []string { return f.rvars }

func (f *funcAr03Bool) Bind(args []interface{}) error {
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

func (f *funcAr03Bool) Func() interface{} {
	return func() bool {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*funcAr03Bool)(nil)
)

// funcAr04Bool implements rfunc.Formula
type funcAr04Bool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64) bool
}

func newFuncAr04Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr04Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr04Bool{}.fct)] = newFuncAr04Bool
}

func (f *funcAr04Bool) RVars() []string { return f.rvars }

func (f *funcAr04Bool) Bind(args []interface{}) error {
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

func (f *funcAr04Bool) Func() interface{} {
	return func() bool {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
		)
	}
}

var (
	_ Formula = (*funcAr04Bool)(nil)
)

// funcAr05Bool implements rfunc.Formula
type funcAr05Bool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64) bool
}

func newFuncAr05Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr05Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr05Bool{}.fct)] = newFuncAr05Bool
}

func (f *funcAr05Bool) RVars() []string { return f.rvars }

func (f *funcAr05Bool) Bind(args []interface{}) error {
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

func (f *funcAr05Bool) Func() interface{} {
	return func() bool {
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
	_ Formula = (*funcAr05Bool)(nil)
)

// funcAr06Bool implements rfunc.Formula
type funcAr06Bool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64) bool
}

func newFuncAr06Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr06Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr06Bool{}.fct)] = newFuncAr06Bool
}

func (f *funcAr06Bool) RVars() []string { return f.rvars }

func (f *funcAr06Bool) Bind(args []interface{}) error {
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

func (f *funcAr06Bool) Func() interface{} {
	return func() bool {
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
	_ Formula = (*funcAr06Bool)(nil)
)

// funcAr07Bool implements rfunc.Formula
type funcAr07Bool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	arg6  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64) bool
}

func newFuncAr07Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr07Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr07Bool{}.fct)] = newFuncAr07Bool
}

func (f *funcAr07Bool) RVars() []string { return f.rvars }

func (f *funcAr07Bool) Bind(args []interface{}) error {
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

func (f *funcAr07Bool) Func() interface{} {
	return func() bool {
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
	_ Formula = (*funcAr07Bool)(nil)
)

// funcAr08Bool implements rfunc.Formula
type funcAr08Bool struct {
	rvars []string
	arg0  *float64
	arg1  *float64
	arg2  *float64
	arg3  *float64
	arg4  *float64
	arg5  *float64
	arg6  *float64
	arg7  *float64
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64) bool
}

func newFuncAr08Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr08Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr08Bool{}.fct)] = newFuncAr08Bool
}

func (f *funcAr08Bool) RVars() []string { return f.rvars }

func (f *funcAr08Bool) Bind(args []interface{}) error {
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

func (f *funcAr08Bool) Func() interface{} {
	return func() bool {
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
	_ Formula = (*funcAr08Bool)(nil)
)

// funcAr09Bool implements rfunc.Formula
type funcAr09Bool struct {
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
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64) bool
}

func newFuncAr09Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr09Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr09Bool{}.fct)] = newFuncAr09Bool
}

func (f *funcAr09Bool) RVars() []string { return f.rvars }

func (f *funcAr09Bool) Bind(args []interface{}) error {
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

func (f *funcAr09Bool) Func() interface{} {
	return func() bool {
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
	_ Formula = (*funcAr09Bool)(nil)
)

// funcAr10Bool implements rfunc.Formula
type funcAr10Bool struct {
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
	fct   func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64, arg09 float64) bool
}

func newFuncAr10Bool(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr10Bool{
		rvars: rvars,
		fct:   fct.(func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64, arg09 float64) bool),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr10Bool{}.fct)] = newFuncAr10Bool
}

func (f *funcAr10Bool) RVars() []string { return f.rvars }

func (f *funcAr10Bool) Bind(args []interface{}) error {
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

func (f *funcAr10Bool) Func() interface{} {
	return func() bool {
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
	_ Formula = (*funcAr10Bool)(nil)
)
