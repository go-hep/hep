// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"fmt"
	"reflect"
)

// funcAr00F32 implements rfunc.Formula
type funcAr00F32 struct {
	fct func() float32
}

func newFuncAr00F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr00F32{
		fct: fct.(func() float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr00F32{}.fct)] = newFuncAr00F32
}

func (f *funcAr00F32) RVars() []string { return nil }

func (f *funcAr00F32) Bind(args []interface{}) error {
	if got, want := len(args), 0; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	return nil
}

func (f *funcAr00F32) Func() interface{} {
	return func() float32 {
		return f.fct()
	}
}

var (
	_ Formula = (*funcAr00F32)(nil)
)

// funcAr01F32 implements rfunc.Formula
type funcAr01F32 struct {
	rvars []string
	arg0  *float32
	fct   func(arg00 float32) float32
}

func newFuncAr01F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr01F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr01F32{}.fct)] = newFuncAr01F32
}

func (f *funcAr01F32) RVars() []string { return f.rvars }

func (f *funcAr01F32) Bind(args []interface{}) error {
	if got, want := len(args), 1; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	return nil
}

func (f *funcAr01F32) Func() interface{} {
	return func() float32 {
		return f.fct(
			*f.arg0,
		)
	}
}

var (
	_ Formula = (*funcAr01F32)(nil)
)

// funcAr02F32 implements rfunc.Formula
type funcAr02F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	fct   func(arg00 float32, arg01 float32) float32
}

func newFuncAr02F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr02F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr02F32{}.fct)] = newFuncAr02F32
}

func (f *funcAr02F32) RVars() []string { return f.rvars }

func (f *funcAr02F32) Bind(args []interface{}) error {
	if got, want := len(args), 2; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	return nil
}

func (f *funcAr02F32) Func() interface{} {
	return func() float32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
		)
	}
}

var (
	_ Formula = (*funcAr02F32)(nil)
)

// funcAr03F32 implements rfunc.Formula
type funcAr03F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32) float32
}

func newFuncAr03F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr03F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr03F32{}.fct)] = newFuncAr03F32
}

func (f *funcAr03F32) RVars() []string { return f.rvars }

func (f *funcAr03F32) Bind(args []interface{}) error {
	if got, want := len(args), 3; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	return nil
}

func (f *funcAr03F32) Func() interface{} {
	return func() float32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
		)
	}
}

var (
	_ Formula = (*funcAr03F32)(nil)
)

// funcAr04F32 implements rfunc.Formula
type funcAr04F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	arg3  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32, arg03 float32) float32
}

func newFuncAr04F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr04F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32, arg03 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr04F32{}.fct)] = newFuncAr04F32
}

func (f *funcAr04F32) RVars() []string { return f.rvars }

func (f *funcAr04F32) Bind(args []interface{}) error {
	if got, want := len(args), 4; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	return nil
}

func (f *funcAr04F32) Func() interface{} {
	return func() float32 {
		return f.fct(
			*f.arg0,
			*f.arg1,
			*f.arg2,
			*f.arg3,
		)
	}
}

var (
	_ Formula = (*funcAr04F32)(nil)
)

// funcAr05F32 implements rfunc.Formula
type funcAr05F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	arg3  *float32
	arg4  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32) float32
}

func newFuncAr05F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr05F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr05F32{}.fct)] = newFuncAr05F32
}

func (f *funcAr05F32) RVars() []string { return f.rvars }

func (f *funcAr05F32) Bind(args []interface{}) error {
	if got, want := len(args), 5; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	return nil
}

func (f *funcAr05F32) Func() interface{} {
	return func() float32 {
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
	_ Formula = (*funcAr05F32)(nil)
)

// funcAr06F32 implements rfunc.Formula
type funcAr06F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	arg3  *float32
	arg4  *float32
	arg5  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32) float32
}

func newFuncAr06F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr06F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr06F32{}.fct)] = newFuncAr06F32
}

func (f *funcAr06F32) RVars() []string { return f.rvars }

func (f *funcAr06F32) Bind(args []interface{}) error {
	if got, want := len(args), 6; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	return nil
}

func (f *funcAr06F32) Func() interface{} {
	return func() float32 {
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
	_ Formula = (*funcAr06F32)(nil)
)

// funcAr07F32 implements rfunc.Formula
type funcAr07F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	arg3  *float32
	arg4  *float32
	arg5  *float32
	arg6  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32) float32
}

func newFuncAr07F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr07F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr07F32{}.fct)] = newFuncAr07F32
}

func (f *funcAr07F32) RVars() []string { return f.rvars }

func (f *funcAr07F32) Bind(args []interface{}) error {
	if got, want := len(args), 7; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	return nil
}

func (f *funcAr07F32) Func() interface{} {
	return func() float32 {
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
	_ Formula = (*funcAr07F32)(nil)
)

// funcAr08F32 implements rfunc.Formula
type funcAr08F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	arg3  *float32
	arg4  *float32
	arg5  *float32
	arg6  *float32
	arg7  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32, arg07 float32) float32
}

func newFuncAr08F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr08F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32, arg07 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr08F32{}.fct)] = newFuncAr08F32
}

func (f *funcAr08F32) RVars() []string { return f.rvars }

func (f *funcAr08F32) Bind(args []interface{}) error {
	if got, want := len(args), 8; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	return nil
}

func (f *funcAr08F32) Func() interface{} {
	return func() float32 {
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
	_ Formula = (*funcAr08F32)(nil)
)

// funcAr09F32 implements rfunc.Formula
type funcAr09F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	arg3  *float32
	arg4  *float32
	arg5  *float32
	arg6  *float32
	arg7  *float32
	arg8  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32, arg07 float32, arg08 float32) float32
}

func newFuncAr09F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr09F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32, arg07 float32, arg08 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr09F32{}.fct)] = newFuncAr09F32
}

func (f *funcAr09F32) RVars() []string { return f.rvars }

func (f *funcAr09F32) Bind(args []interface{}) error {
	if got, want := len(args), 9; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	return nil
}

func (f *funcAr09F32) Func() interface{} {
	return func() float32 {
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
	_ Formula = (*funcAr09F32)(nil)
)

// funcAr10F32 implements rfunc.Formula
type funcAr10F32 struct {
	rvars []string
	arg0  *float32
	arg1  *float32
	arg2  *float32
	arg3  *float32
	arg4  *float32
	arg5  *float32
	arg6  *float32
	arg7  *float32
	arg8  *float32
	arg9  *float32
	fct   func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32, arg07 float32, arg08 float32, arg09 float32) float32
}

func newFuncAr10F32(rvars []string, fct interface{}) (Formula, error) {
	return &funcAr10F32{
		rvars: rvars,
		fct:   fct.(func(arg00 float32, arg01 float32, arg02 float32, arg03 float32, arg04 float32, arg05 float32, arg06 float32, arg07 float32, arg08 float32, arg09 float32) float32),
	}, nil
}

func init() {
	funcs[reflect.TypeOf(funcAr10F32{}.fct)] = newFuncAr10F32
}

func (f *funcAr10F32) RVars() []string { return f.rvars }

func (f *funcAr10F32) Bind(args []interface{}) error {
	if got, want := len(args), 10; got != want {
		return fmt.Errorf(
			"rfunc: invalid number of bind arguments (got=%d, want=%d)",
			got, want,
		)
	}
	{
		ptr, ok := args[0].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 0 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[0], args[0],
			)
		}
		f.arg0 = ptr
	}
	{
		ptr, ok := args[1].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 1 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[1], args[1],
			)
		}
		f.arg1 = ptr
	}
	{
		ptr, ok := args[2].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 2 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[2], args[2],
			)
		}
		f.arg2 = ptr
	}
	{
		ptr, ok := args[3].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 3 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[3], args[3],
			)
		}
		f.arg3 = ptr
	}
	{
		ptr, ok := args[4].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 4 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[4], args[4],
			)
		}
		f.arg4 = ptr
	}
	{
		ptr, ok := args[5].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 5 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[5], args[5],
			)
		}
		f.arg5 = ptr
	}
	{
		ptr, ok := args[6].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 6 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[6], args[6],
			)
		}
		f.arg6 = ptr
	}
	{
		ptr, ok := args[7].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 7 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[7], args[7],
			)
		}
		f.arg7 = ptr
	}
	{
		ptr, ok := args[8].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 8 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[8], args[8],
			)
		}
		f.arg8 = ptr
	}
	{
		ptr, ok := args[9].(*float32)
		if !ok {
			return fmt.Errorf(
				"rfunc: argument type 9 (name=%s) mismatch: got=%T, want=*float32",
				f.rvars[9], args[9],
			)
		}
		f.arg9 = ptr
	}
	return nil
}

func (f *funcAr10F32) Func() interface{} {
	return func() float32 {
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
	_ Formula = (*funcAr10F32)(nil)
)
