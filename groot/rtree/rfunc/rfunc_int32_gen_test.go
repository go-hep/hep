// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"reflect"
	"testing"
)

func TestFuncAr00I32(t *testing.T) {

	var rvars []string

	fct := func() int32 {
		return 42
	}

	form, err := newFuncAr00I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr00I32 formula: %+v", err)
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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr01I32(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 int32) int32 {
		return 42
	}

	form, err := newFuncAr01I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr01I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 1)
	ptrs[0] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr02I32(t *testing.T) {

	rvars := make([]string, 2)
	rvars[0] = "name-0"
	rvars[1] = "name-1"

	fct := func(arg00 int32, arg01 int32) int32 {
		return 42
	}

	form, err := newFuncAr02I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr02I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 2)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr03I32(t *testing.T) {

	rvars := make([]string, 3)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"

	fct := func(arg00 int32, arg01 int32, arg02 int32) int32 {
		return 42
	}

	form, err := newFuncAr03I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr03I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 3)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr04I32(t *testing.T) {

	rvars := make([]string, 4)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"

	fct := func(arg00 int32, arg01 int32, arg02 int32, arg03 int32) int32 {
		return 42
	}

	form, err := newFuncAr04I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr04I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 4)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)
	ptrs[3] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr05I32(t *testing.T) {

	rvars := make([]string, 5)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"

	fct := func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32) int32 {
		return 42
	}

	form, err := newFuncAr05I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr05I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 5)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)
	ptrs[3] = new(int32)
	ptrs[4] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr06I32(t *testing.T) {

	rvars := make([]string, 6)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"

	fct := func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32) int32 {
		return 42
	}

	form, err := newFuncAr06I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr06I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 6)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)
	ptrs[3] = new(int32)
	ptrs[4] = new(int32)
	ptrs[5] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr07I32(t *testing.T) {

	rvars := make([]string, 7)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"

	fct := func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32) int32 {
		return 42
	}

	form, err := newFuncAr07I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr07I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 7)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)
	ptrs[3] = new(int32)
	ptrs[4] = new(int32)
	ptrs[5] = new(int32)
	ptrs[6] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr08I32(t *testing.T) {

	rvars := make([]string, 8)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"
	rvars[7] = "name-7"

	fct := func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32) int32 {
		return 42
	}

	form, err := newFuncAr08I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr08I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 8)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)
	ptrs[3] = new(int32)
	ptrs[4] = new(int32)
	ptrs[5] = new(int32)
	ptrs[6] = new(int32)
	ptrs[7] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr09I32(t *testing.T) {

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

	fct := func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32, arg08 int32) int32 {
		return 42
	}

	form, err := newFuncAr09I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr09I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 9)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)
	ptrs[3] = new(int32)
	ptrs[4] = new(int32)
	ptrs[5] = new(int32)
	ptrs[6] = new(int32)
	ptrs[7] = new(int32)
	ptrs[8] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr10I32(t *testing.T) {

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

	fct := func(arg00 int32, arg01 int32, arg02 int32, arg03 int32, arg04 int32, arg05 int32, arg06 int32, arg07 int32, arg08 int32, arg09 int32) int32 {
		return 42
	}

	form, err := newFuncAr10I32(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr10I32 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 10)
	ptrs[0] = new(int32)
	ptrs[1] = new(int32)
	ptrs[2] = new(int32)
	ptrs[3] = new(int32)
	ptrs[4] = new(int32)
	ptrs[5] = new(int32)
	ptrs[6] = new(int32)
	ptrs[7] = new(int32)
	ptrs[8] = new(int32)
	ptrs[9] = new(int32)

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

	got := form.Func().(func() int32)()
	if got, want := got, int32(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
