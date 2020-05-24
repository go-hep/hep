// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
	"reflect"
)

// funcAr00U32 implements rfunc.Formula
type funcAr00U32 struct {
	fct func() uint32
}

func newFuncAr00U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr00U32{
		fct: fct.(func() uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr00U32{}.fct)] = newFuncAr00U32
}

func (f *funcAr00U32) RVars() []string { return nil }

func (f *funcAr00U32) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

func (f *funcAr00U32) Func() interface{} {
	return func() uint32 {
		return f.fct()
	}
}

var (
	_ Formula = (*funcAr00U32)(nil)
)

// funcAr01U32 implements rfunc.Formula
type funcAr01U32 struct {
	rvars []string
	arg0  *uint32
	fct   func(arg00 uint32) uint32
}

func newFuncAr01U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr01U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr01U32{}.fct)] = newFuncAr01U32
}

func (f *funcAr01U32) RVars() []string { return f.rvars }

func (f *funcAr01U32) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

func (f *funcAr01U32) Func() interface{} {
	return func() uint32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*funcAr01U32)(nil)
)

// funcAr02U32 implements rfunc.Formula
type funcAr02U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	fct   func(arg00 uint32, arg01 uint32) uint32
}

func newFuncAr02U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr02U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr02U32{}.fct)] = newFuncAr02U32
}

func (f *funcAr02U32) RVars() []string { return f.rvars }

func (f *funcAr02U32) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

func (f *funcAr02U32) Func() interface{} {
	return func() uint32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*funcAr02U32)(nil)
)

// funcAr03U32 implements rfunc.Formula
type funcAr03U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32) uint32
}

func newFuncAr03U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr03U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr03U32{}.fct)] = newFuncAr03U32
}

func (f *funcAr03U32) RVars() []string { return f.rvars }

func (f *funcAr03U32) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

func (f *funcAr03U32) Func() interface{} {
	return func() uint32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*funcAr03U32)(nil)
)

// funcAr04U32 implements rfunc.Formula
type funcAr04U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	arg3  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32) uint32
}

func newFuncAr04U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr04U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr04U32{}.fct)] = newFuncAr04U32
}

func (f *funcAr04U32) RVars() []string { return f.rvars }

func (f *funcAr04U32) Bind(args []interface{}) error {
	if got, want := len(args), 4; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	return nil
}

func (f *funcAr04U32) Func() interface{} {
	return func() uint32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
		)
	}
}

var (
	_ Formula = (*funcAr04U32)(nil)
)

// funcAr05U32 implements rfunc.Formula
type funcAr05U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	arg3  *uint32
	arg4  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32) uint32
}

func newFuncAr05U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr05U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr05U32{}.fct)] = newFuncAr05U32
}

func (f *funcAr05U32) RVars() []string { return f.rvars }

func (f *funcAr05U32) Bind(args []interface{}) error {
	if got, want := len(args), 5; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	return nil
}

func (f *funcAr05U32) Func() interface{} {
	return func() uint32 {
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
	_ Formula = (*funcAr05U32)(nil)
)

// funcAr06U32 implements rfunc.Formula
type funcAr06U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	arg3  *uint32
	arg4  *uint32
	arg5  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32) uint32
}

func newFuncAr06U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr06U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr06U32{}.fct)] = newFuncAr06U32
}

func (f *funcAr06U32) RVars() []string { return f.rvars }

func (f *funcAr06U32) Bind(args []interface{}) error {
	if got, want := len(args), 6; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	return nil
}

func (f *funcAr06U32) Func() interface{} {
	return func() uint32 {
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
	_ Formula = (*funcAr06U32)(nil)
)

// funcAr07U32 implements rfunc.Formula
type funcAr07U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	arg3  *uint32
	arg4  *uint32
	arg5  *uint32
	arg6  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32) uint32
}

func newFuncAr07U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr07U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr07U32{}.fct)] = newFuncAr07U32
}

func (f *funcAr07U32) RVars() []string { return f.rvars }

func (f *funcAr07U32) Bind(args []interface{}) error {
	if got, want := len(args), 7; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	return nil
}

func (f *funcAr07U32) Func() interface{} {
	return func() uint32 {
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
	_ Formula = (*funcAr07U32)(nil)
)

// funcAr08U32 implements rfunc.Formula
type funcAr08U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	arg3  *uint32
	arg4  *uint32
	arg5  *uint32
	arg6  *uint32
	arg7  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32, arg07 uint32) uint32
}

func newFuncAr08U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr08U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32, arg07 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr08U32{}.fct)] = newFuncAr08U32
}

func (f *funcAr08U32) RVars() []string { return f.rvars }

func (f *funcAr08U32) Bind(args []interface{}) error {
	if got, want := len(args), 8; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	return nil
}

func (f *funcAr08U32) Func() interface{} {
	return func() uint32 {
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
	_ Formula = (*funcAr08U32)(nil)
)

// funcAr09U32 implements rfunc.Formula
type funcAr09U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	arg3  *uint32
	arg4  *uint32
	arg5  *uint32
	arg6  *uint32
	arg7  *uint32
	arg8  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32, arg07 uint32, arg08 uint32) uint32
}

func newFuncAr09U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr09U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32, arg07 uint32, arg08 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr09U32{}.fct)] = newFuncAr09U32
}

func (f *funcAr09U32) RVars() []string { return f.rvars }

func (f *funcAr09U32) Bind(args []interface{}) error {
	if got, want := len(args), 9; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	return nil
}

func (f *funcAr09U32) Func() interface{} {
	return func() uint32 {
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
	_ Formula = (*funcAr09U32)(nil)
)

// funcAr10U32 implements rfunc.Formula
type funcAr10U32 struct {
	rvars []string
	arg0  *uint32
	arg1  *uint32
	arg2  *uint32
	arg3  *uint32
	arg4  *uint32
	arg5  *uint32
	arg6  *uint32
	arg7  *uint32
	arg8  *uint32
	arg9  *uint32
	fct   func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32, arg07 uint32, arg08 uint32, arg09 uint32) uint32
}

func newFuncAr10U32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr10U32{
		rvars: rvars,
		fct:   fct.(func(arg00 uint32, arg01 uint32, arg02 uint32, arg03 uint32, arg04 uint32, arg05 uint32, arg06 uint32, arg07 uint32, arg08 uint32, arg09 uint32) uint32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr10U32{}.fct)] = newFuncAr10U32
}

func (f *funcAr10U32) RVars() []string { return f.rvars }

func (f *funcAr10U32) Bind(args []interface{}) error {
	if got, want := len(args), 10; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	{
		ptr, ok := args[9].(*uint32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 9 (name=%s) mismatch: got=%T, want=*uint32",
				f.rvars[9], args[9],
			)
		}
		f.arg9 = ptr
	}
	return nil
}

func (f *funcAr10U32) Func() interface{} {
	return func() uint32 {
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
	_ Formula = (*funcAr10U32)(nil)
)
