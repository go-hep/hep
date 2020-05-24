// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"reflect"
	"testing"
)

func TestFuncAr00F64(t *testing.T) {

	var rvars []string

	fct := func() float64 {
		return 42
	}

	form, err := newFuncAr00F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr00F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	var ptrs []interface{}

	{
		bad := make([]interface{}, 1)
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr01F64(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 float64) float64 {
		return 42
	}

	form, err := newFuncAr01F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr01F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 1)
	ptrs[0] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr02F64(t *testing.T) {

	rvars := make([]string, 2)
	rvars[0] = "name-0"
	rvars[1] = "name-1"

	fct := func(arg00 float64, arg01 float64) float64 {
		return 42
	}

	form, err := newFuncAr02F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr02F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 2)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr03F64(t *testing.T) {

	rvars := make([]string, 3)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"

	fct := func(arg00 float64, arg01 float64, arg02 float64) float64 {
		return 42
	}

	form, err := newFuncAr03F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr03F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 3)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr04F64(t *testing.T) {

	rvars := make([]string, 4)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"

	fct := func(arg00 float64, arg01 float64, arg02 float64, arg03 float64) float64 {
		return 42
	}

	form, err := newFuncAr04F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr04F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 4)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)
	ptrs[3] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr05F64(t *testing.T) {

	rvars := make([]string, 5)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"

	fct := func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64) float64 {
		return 42
	}

	form, err := newFuncAr05F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr05F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 5)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)
	ptrs[3] = new(float64)
	ptrs[4] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr06F64(t *testing.T) {

	rvars := make([]string, 6)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"

	fct := func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64) float64 {
		return 42
	}

	form, err := newFuncAr06F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr06F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 6)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)
	ptrs[3] = new(float64)
	ptrs[4] = new(float64)
	ptrs[5] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr07F64(t *testing.T) {

	rvars := make([]string, 7)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"

	fct := func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64) float64 {
		return 42
	}

	form, err := newFuncAr07F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr07F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 7)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)
	ptrs[3] = new(float64)
	ptrs[4] = new(float64)
	ptrs[5] = new(float64)
	ptrs[6] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr08F64(t *testing.T) {

	rvars := make([]string, 8)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"
	rvars[7] = "name-7"

	fct := func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64) float64 {
		return 42
	}

	form, err := newFuncAr08F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr08F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 8)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)
	ptrs[3] = new(float64)
	ptrs[4] = new(float64)
	ptrs[5] = new(float64)
	ptrs[6] = new(float64)
	ptrs[7] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr09F64(t *testing.T) {

	rvars := make([]string, 9)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"
	rvars[7] = "name-7"
	rvars[8] = "name-8"

	fct := func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64) float64 {
		return 42
	}

	form, err := newFuncAr09F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr09F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 9)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)
	ptrs[3] = new(float64)
	ptrs[4] = new(float64)
	ptrs[5] = new(float64)
	ptrs[6] = new(float64)
	ptrs[7] = new(float64)
	ptrs[8] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr10F64(t *testing.T) {

	rvars := make([]string, 10)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"
	rvars[7] = "name-7"
	rvars[8] = "name-8"
	rvars[9] = "name-9"

	fct := func(arg00 float64, arg01 float64, arg02 float64, arg03 float64, arg04 float64, arg05 float64, arg06 float64, arg07 float64, arg08 float64, arg09 float64) float64 {
		return 42
	}

	form, err := newFuncAr10F64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr10F64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 10)
	ptrs[0] = new(float64)
	ptrs[1] = new(float64)
	ptrs[2] = new(float64)
	ptrs[3] = new(float64)
	ptrs[4] = new(float64)
	ptrs[5] = new(float64)
	ptrs[6] = new(float64)
	ptrs[7] = new(float64)
	ptrs[8] = new(float64)
	ptrs[9] = new(float64)

	{
		bad := make([]interface{}, len(ptrs))
		copy(bad, ptrs)
		for i := len(ptrs) - 1; i >= 0; i-- {
			bad[i] = interface{}(nil)
			err := form.Bind(bad)
			if err == nil {
				t.Fatalf("expected an error for empty iface")
			}
		}
		bad = append(bad, interface{}(nil))
		err := form.Bind(bad)
		if err == nil {
			t.Fatalf("expected an error for invalid args length")
		}
	}

	err = form.Bind(ptrs)
	if err != nil {
		t.Fatalf("could not bind formula: %+v", err)
	}

	got := form.Func().(func() float64)()
	if got, want := got, float64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
