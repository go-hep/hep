// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
	"reflect"
)

// funcAr00U64 implements rfunc.Formula
type funcAr00U64 struct {
	fct func() uint64
}

func newFuncAr00U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr00U64{
		fct: fct.(func() uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr00U64{}.fct)] = newFuncAr00U64
}

func (f *funcAr00U64) RVars() []string { return nil }

func (f *funcAr00U64) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

func (f *funcAr00U64) Func() interface{} {
	return func() uint64 {
		return f.fct()
	}
}

var (
	_ Formula = (*funcAr00U64)(nil)
)

// funcAr01U64 implements rfunc.Formula
type funcAr01U64 struct {
	rvars []string
	arg0  *uint64
	fct   func(arg00 uint64) uint64
}

func newFuncAr01U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr01U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr01U64{}.fct)] = newFuncAr01U64
}

func (f *funcAr01U64) RVars() []string { return f.rvars }

func (f *funcAr01U64) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

func (f *funcAr01U64) Func() interface{} {
	return func() uint64 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*funcAr01U64)(nil)
)

// funcAr02U64 implements rfunc.Formula
type funcAr02U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	fct   func(arg00 uint64, arg01 uint64) uint64
}

func newFuncAr02U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr02U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr02U64{}.fct)] = newFuncAr02U64
}

func (f *funcAr02U64) RVars() []string { return f.rvars }

func (f *funcAr02U64) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

func (f *funcAr02U64) Func() interface{} {
	return func() uint64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*funcAr02U64)(nil)
)

// funcAr03U64 implements rfunc.Formula
type funcAr03U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64) uint64
}

func newFuncAr03U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr03U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr03U64{}.fct)] = newFuncAr03U64
}

func (f *funcAr03U64) RVars() []string { return f.rvars }

func (f *funcAr03U64) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

func (f *funcAr03U64) Func() interface{} {
	return func() uint64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*funcAr03U64)(nil)
)

// funcAr04U64 implements rfunc.Formula
type funcAr04U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	arg3  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64) uint64
}

func newFuncAr04U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr04U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr04U64{}.fct)] = newFuncAr04U64
}

func (f *funcAr04U64) RVars() []string { return f.rvars }

func (f *funcAr04U64) Bind(args []interface{}) error {
	if got, want := len(args), 4; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	return nil
}

func (f *funcAr04U64) Func() interface{} {
	return func() uint64 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
		)
	}
}

var (
	_ Formula = (*funcAr04U64)(nil)
)

// funcAr05U64 implements rfunc.Formula
type funcAr05U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	arg3  *uint64
	arg4  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64) uint64
}

func newFuncAr05U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr05U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr05U64{}.fct)] = newFuncAr05U64
}

func (f *funcAr05U64) RVars() []string { return f.rvars }

func (f *funcAr05U64) Bind(args []interface{}) error {
	if got, want := len(args), 5; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	return nil
}

func (f *funcAr05U64) Func() interface{} {
	return func() uint64 {
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
	_ Formula = (*funcAr05U64)(nil)
)

// funcAr06U64 implements rfunc.Formula
type funcAr06U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	arg3  *uint64
	arg4  *uint64
	arg5  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64) uint64
}

func newFuncAr06U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr06U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr06U64{}.fct)] = newFuncAr06U64
}

func (f *funcAr06U64) RVars() []string { return f.rvars }

func (f *funcAr06U64) Bind(args []interface{}) error {
	if got, want := len(args), 6; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	return nil
}

func (f *funcAr06U64) Func() interface{} {
	return func() uint64 {
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
	_ Formula = (*funcAr06U64)(nil)
)

// funcAr07U64 implements rfunc.Formula
type funcAr07U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	arg3  *uint64
	arg4  *uint64
	arg5  *uint64
	arg6  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64) uint64
}

func newFuncAr07U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr07U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr07U64{}.fct)] = newFuncAr07U64
}

func (f *funcAr07U64) RVars() []string { return f.rvars }

func (f *funcAr07U64) Bind(args []interface{}) error {
	if got, want := len(args), 7; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	return nil
}

func (f *funcAr07U64) Func() interface{} {
	return func() uint64 {
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
	_ Formula = (*funcAr07U64)(nil)
)

// funcAr08U64 implements rfunc.Formula
type funcAr08U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	arg3  *uint64
	arg4  *uint64
	arg5  *uint64
	arg6  *uint64
	arg7  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64) uint64
}

func newFuncAr08U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr08U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr08U64{}.fct)] = newFuncAr08U64
}

func (f *funcAr08U64) RVars() []string { return f.rvars }

func (f *funcAr08U64) Bind(args []interface{}) error {
	if got, want := len(args), 8; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	return nil
}

func (f *funcAr08U64) Func() interface{} {
	return func() uint64 {
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
	_ Formula = (*funcAr08U64)(nil)
)

// funcAr09U64 implements rfunc.Formula
type funcAr09U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	arg3  *uint64
	arg4  *uint64
	arg5  *uint64
	arg6  *uint64
	arg7  *uint64
	arg8  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64, arg08 uint64) uint64
}

func newFuncAr09U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr09U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64, arg08 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr09U64{}.fct)] = newFuncAr09U64
}

func (f *funcAr09U64) RVars() []string { return f.rvars }

func (f *funcAr09U64) Bind(args []interface{}) error {
	if got, want := len(args), 9; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	return nil
}

func (f *funcAr09U64) Func() interface{} {
	return func() uint64 {
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
	_ Formula = (*funcAr09U64)(nil)
)

// funcAr10U64 implements rfunc.Formula
type funcAr10U64 struct {
	rvars []string
	arg0  *uint64
	arg1  *uint64
	arg2  *uint64
	arg3  *uint64
	arg4  *uint64
	arg5  *uint64
	arg6  *uint64
	arg7  *uint64
	arg8  *uint64
	arg9  *uint64
	fct   func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64, arg08 uint64, arg09 uint64) uint64
}

func newFuncAr10U64(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr10U64{
		rvars: rvars,
		fct:   fct.(func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64, arg08 uint64, arg09 uint64) uint64),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr10U64{}.fct)] = newFuncAr10U64
}

func (f *funcAr10U64) RVars() []string { return f.rvars }

func (f *funcAr10U64) Bind(args []interface{}) error {
	if got, want := len(args), 10; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	{
		ptr, ok := args[9].(*uint64)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 9 (name=%s) mismatch: got=%T, want=*uint64",
				f.rvars[9], args[9],
			)
		}
		f.arg9 = ptr
	}
	return nil
}

func (f *funcAr10U64) Func() interface{} {
	return func() uint64 {
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
	_ Formula = (*funcAr10U64)(nil)
)
