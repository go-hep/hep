// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
	"reflect"
)

// funcAr00I32 implements rfunc.Formula
type funcAr00I32 struct {
	fct func() int32
}

func newFuncAr00I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr00I32{
		fct: fct.(func() int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr00I32{}.fct)] = newFuncAr00I32
}

func (f *funcAr00I32) RVars() []string { return nil }

func (f *funcAr00I32) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

func (f *funcAr00I32) Func() interface{} {
	return func() int32 {
		return f.fct()
	}
}

var (
	_ Formula = (*funcAr00I32)(nil)
)

// funcAr01I32 implements rfunc.Formula
type funcAr01I32 struct {
	rvars []string
	arg0  *int32
	fct   func(arg00 int32) int32
}

func newFuncAr01I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr01I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr01I32{}.fct)] = newFuncAr01I32
}

func (f *funcAr01I32) RVars() []string { return f.rvars }

func (f *funcAr01I32) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

func (f *funcAr01I32) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*funcAr01I32)(nil)
)

// funcAr02I32 implements rfunc.Formula
type funcAr02I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	fct   func(arg00 int32, arg01 int32) int32
}

func newFuncAr02I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr02I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr02I32{}.fct)] = newFuncAr02I32
}

func (f *funcAr02I32) RVars() []string { return f.rvars }

func (f *funcAr02I32) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

func (f *funcAr02I32) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*funcAr02I32)(nil)
)

// funcAr03I32 implements rfunc.Formula
type funcAr03I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32) int32
}

func newFuncAr03I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr03I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr03I32{}.fct)] = newFuncAr03I32
}

func (f *funcAr03I32) RVars() []string { return f.rvars }

func (f *funcAr03I32) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

func (f *funcAr03I32) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*funcAr03I32)(nil)
)

// funcAr04I32 implements rfunc.Formula
type funcAr04I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	arg3  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32, arg03 int32) int32
}

func newFuncAr04I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr04I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32, arg03 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr04I32{}.fct)] = newFuncAr04I32
}

func (f *funcAr04I32) RVars() []string { return f.rvars }

func (f *funcAr04I32) Bind(args []interface{}) error {
	if got, want := len(args), 4; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	return nil
}

func (f *funcAr04I32) Func() interface{} {
	return func() int32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
		)
	}
}

var (
	_ Formula = (*funcAr04I32)(nil)
)

// funcAr05I32 implements rfunc.Formula
type funcAr05I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	arg3  *int32
	arg4  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32) int32
}

func newFuncAr05I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr05I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr05I32{}.fct)] = newFuncAr05I32
}

func (f *funcAr05I32) RVars() []string { return f.rvars }

func (f *funcAr05I32) Bind(args []interface{}) error {
	if got, want := len(args), 5; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	return nil
}

func (f *funcAr05I32) Func() interface{} {
	return func() int32 {
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
	_ Formula = (*funcAr05I32)(nil)
)

// funcAr06I32 implements rfunc.Formula
type funcAr06I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	arg3  *int32
	arg4  *int32
	arg5  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32) int32
}

func newFuncAr06I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr06I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr06I32{}.fct)] = newFuncAr06I32
}

func (f *funcAr06I32) RVars() []string { return f.rvars }

func (f *funcAr06I32) Bind(args []interface{}) error {
	if got, want := len(args), 6; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	return nil
}

func (f *funcAr06I32) Func() interface{} {
	return func() int32 {
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
	_ Formula = (*funcAr06I32)(nil)
)

// funcAr07I32 implements rfunc.Formula
type funcAr07I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	arg3  *int32
	arg4  *int32
	arg5  *int32
	arg6  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32) int32
}

func newFuncAr07I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr07I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr07I32{}.fct)] = newFuncAr07I32
}

func (f *funcAr07I32) RVars() []string { return f.rvars }

func (f *funcAr07I32) Bind(args []interface{}) error {
	if got, want := len(args), 7; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	return nil
}

func (f *funcAr07I32) Func() interface{} {
	return func() int32 {
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
	_ Formula = (*funcAr07I32)(nil)
)

// funcAr08I32 implements rfunc.Formula
type funcAr08I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	arg3  *int32
	arg4  *int32
	arg5  *int32
	arg6  *int32
	arg7  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32) int32
}

func newFuncAr08I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr08I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr08I32{}.fct)] = newFuncAr08I32
}

func (f *funcAr08I32) RVars() []string { return f.rvars }

func (f *funcAr08I32) Bind(args []interface{}) error {
	if got, want := len(args), 8; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	return nil
}

func (f *funcAr08I32) Func() interface{} {
	return func() int32 {
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
	_ Formula = (*funcAr08I32)(nil)
)

// funcAr09I32 implements rfunc.Formula
type funcAr09I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	arg3  *int32
	arg4  *int32
	arg5  *int32
	arg6  *int32
	arg7  *int32
	arg8  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32, arg08 int32) int32
}

func newFuncAr09I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr09I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32, arg08 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr09I32{}.fct)] = newFuncAr09I32
}

func (f *funcAr09I32) RVars() []string { return f.rvars }

func (f *funcAr09I32) Bind(args []interface{}) error {
	if got, want := len(args), 9; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	return nil
}

func (f *funcAr09I32) Func() interface{} {
	return func() int32 {
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
	_ Formula = (*funcAr09I32)(nil)
)

// funcAr10I32 implements rfunc.Formula
type funcAr10I32 struct {
	rvars []string
	arg0  *int32
	arg1  *int32
	arg2  *int32
	arg3  *int32
	arg4  *int32
	arg5  *int32
	arg6  *int32
	arg7  *int32
	arg8  *int32
	arg9  *int32
	fct   func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32, arg08 int32, arg09 int32) int32
}

func newFuncAr10I32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr10I32{
		rvars: rvars,
		fct:   fct.(func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32, arg08 int32, arg09 int32) int32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr10I32{}.fct)] = newFuncAr10I32
}

func (f *funcAr10I32) RVars() []string { return f.rvars }

func (f *funcAr10I32) Bind(args []interface{}) error {
	if got, want := len(args), 10; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	{
		ptr, ok := args[9].(*int32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 9 (name=%s) mismatch: got=%T, want=*int32",
				f.rvars[9], args[9],
			)
		}
		f.arg9 = ptr
	}
	return nil
}

func (f *funcAr10I32) Func() interface{} {
	return func() int32 {
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
	_ Formula = (*funcAr10I32)(nil)
)
