// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
	"reflect"
)

// funcAr00I64 implements rfunc.Formula
type funcAr00I64 struct {
	fct func() int64
}

func newFuncAr00I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr00I64{
		fct: fct.(func() int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr00I64{}.fct)] = newFuncAr00I64
}

func (f *funcAr00I64) RVars() []string { return nil }

func (f *funcAr00I64) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

func (f *funcAr00I64) Func() interface{} {
	return func() int64 {
		return f.fct()
	}
}

var (
	_ Formula = (*funcAr00I64)(nil)
)

// funcAr01I64 implements rfunc.Formula
type funcAr01I64 struct {
	rvars []string
	arg0  *int64
	fct   func(arg00 int64) int64
}

func newFuncAr01I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr01I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr01I64{}.fct)] = newFuncAr01I64
}

func (f *funcAr01I64) RVars() []string { return f.rvars }

func (f *funcAr01I64) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

func (f *funcAr01I64) Func() interface{} {
	return func() int64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*funcAr01I64)(nil)
)

// funcAr02I64 implements rfunc.Formula
type funcAr02I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	fct   func(arg00 int64, arg01 int64) int64
}

func newFuncAr02I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr02I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr02I64{}.fct)] = newFuncAr02I64
}

func (f *funcAr02I64) RVars() []string { return f.rvars }

func (f *funcAr02I64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

func (f *funcAr02I64) Func() interface{} {
	return func() int64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*funcAr02I64)(nil)
)

// funcAr03I64 implements rfunc.Formula
type funcAr03I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64) int64
}

func newFuncAr03I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr03I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr03I64{}.fct)] = newFuncAr03I64
}

func (f *funcAr03I64) RVars() []string { return f.rvars }

func (f *funcAr03I64) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

func (f *funcAr03I64) Func() interface{} {
	return func() int64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*funcAr03I64)(nil)
)

// funcAr04I64 implements rfunc.Formula
type funcAr04I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	arg3  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64, arg03 int64) int64
}

func newFuncAr04I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr04I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64, arg03 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr04I64{}.fct)] = newFuncAr04I64
}

func (f *funcAr04I64) RVars() []string { return f.rvars }

func (f *funcAr04I64) Bind(args []interface{}) error {
	if got, want := len(args), 4; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	return nil
}

func (f *funcAr04I64) Func() interface{} {
	return func() int64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
		)
	}
}

var (
	_ Formula = (*funcAr04I64)(nil)
)

// funcAr05I64 implements rfunc.Formula
type funcAr05I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	arg3  *int64
	arg4  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64) int64
}

func newFuncAr05I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr05I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr05I64{}.fct)] = newFuncAr05I64
}

func (f *funcAr05I64) RVars() []string { return f.rvars }

func (f *funcAr05I64) Bind(args []interface{}) error {
	if got, want := len(args), 5; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	return nil
}

func (f *funcAr05I64) Func() interface{} {
	return func() int64 {
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
	_ Formula = (*funcAr05I64)(nil)
)

// funcAr06I64 implements rfunc.Formula
type funcAr06I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	arg3  *int64
	arg4  *int64
	arg5  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64) int64
}

func newFuncAr06I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr06I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr06I64{}.fct)] = newFuncAr06I64
}

func (f *funcAr06I64) RVars() []string { return f.rvars }

func (f *funcAr06I64) Bind(args []interface{}) error {
	if got, want := len(args), 6; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	return nil
}

func (f *funcAr06I64) Func() interface{} {
	return func() int64 {
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
	_ Formula = (*funcAr06I64)(nil)
)

// funcAr07I64 implements rfunc.Formula
type funcAr07I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	arg3  *int64
	arg4  *int64
	arg5  *int64
	arg6  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64) int64
}

func newFuncAr07I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr07I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr07I64{}.fct)] = newFuncAr07I64
}

func (f *funcAr07I64) RVars() []string { return f.rvars }

func (f *funcAr07I64) Bind(args []interface{}) error {
	if got, want := len(args), 7; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	return nil
}

func (f *funcAr07I64) Func() interface{} {
	return func() int64 {
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
	_ Formula = (*funcAr07I64)(nil)
)

// funcAr08I64 implements rfunc.Formula
type funcAr08I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	arg3  *int64
	arg4  *int64
	arg5  *int64
	arg6  *int64
	arg7  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64, arg07 int64) int64
}

func newFuncAr08I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr08I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64, arg07 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr08I64{}.fct)] = newFuncAr08I64
}

func (f *funcAr08I64) RVars() []string { return f.rvars }

func (f *funcAr08I64) Bind(args []interface{}) error {
	if got, want := len(args), 8; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	return nil
}

func (f *funcAr08I64) Func() interface{} {
	return func() int64 {
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
	_ Formula = (*funcAr08I64)(nil)
)

// funcAr09I64 implements rfunc.Formula
type funcAr09I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	arg3  *int64
	arg4  *int64
	arg5  *int64
	arg6  *int64
	arg7  *int64
	arg8  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64, arg07 int64, arg08 int64) int64
}

func newFuncAr09I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr09I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64, arg07 int64, arg08 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr09I64{}.fct)] = newFuncAr09I64
}

func (f *funcAr09I64) RVars() []string { return f.rvars }

func (f *funcAr09I64) Bind(args []interface{}) error {
	if got, want := len(args), 9; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	return nil
}

func (f *funcAr09I64) Func() interface{} {
	return func() int64 {
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
	_ Formula = (*funcAr09I64)(nil)
)

// funcAr10I64 implements rfunc.Formula
type funcAr10I64 struct {
	rvars []string
	arg0  *int64
	arg1  *int64
	arg2  *int64
	arg3  *int64
	arg4  *int64
	arg5  *int64
	arg6  *int64
	arg7  *int64
	arg8  *int64
	arg9  *int64
	fct   func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64, arg07 int64, arg08 int64, arg09 int64) int64
}

func newFuncAr10I64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr10I64{
		rvars: rvars,
		fct:   fct.(func(arg00 int64, arg01 int64, arg02 int64, arg03 int64, arg04 int64, arg05 int64, arg06 int64, arg07 int64, arg08 int64, arg09 int64) int64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr10I64{}.fct)] = newFuncAr10I64
}

func (f *funcAr10I64) RVars() []string { return f.rvars }

func (f *funcAr10I64) Bind(args []interface{}) error {
	if got, want := len(args), 10; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	{
		ptr, ok := args[9].(*int64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 9 (name=%s) mismatch: got=%T, want=*int64",
				f.rvars[9], args[9],
			)
		}
		f.arg9 = ptr
	}
	return nil
}

func (f *funcAr10I64) Func() interface{} {
	return func() int64 {
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
	_ Formula = (*funcAr10I64)(nil)
)
