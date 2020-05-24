// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Automatically generated. DO NOT EDIT.

package rfunc

import (
	"reflect"
	"testing"
)

func TestFuncAr00U64(t *testing.T) {

	var rvars []string

	fct := func() uint64 {
		return 42
	}

	form, err := newFuncAr00U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr00U64 formula: %+v", err)
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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr01U64(t *testing.T) {

	rvars := make([]string, 1)
	rvars[0] = "name-0"

	fct := func(arg00 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr01U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr01U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 1)
	ptrs[0] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr02U64(t *testing.T) {

	rvars := make([]string, 2)
	rvars[0] = "name-0"
	rvars[1] = "name-1"

	fct := func(arg00 uint64, arg01 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr02U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr02U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 2)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr03U64(t *testing.T) {

	rvars := make([]string, 3)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr03U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr03U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 3)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr04U64(t *testing.T) {

	rvars := make([]string, 4)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr04U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr04U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 4)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)
	ptrs[3] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr05U64(t *testing.T) {

	rvars := make([]string, 5)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr05U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr05U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 5)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)
	ptrs[3] = new(uint64)
	ptrs[4] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr06U64(t *testing.T) {

	rvars := make([]string, 6)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr06U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr06U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 6)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)
	ptrs[3] = new(uint64)
	ptrs[4] = new(uint64)
	ptrs[5] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr07U64(t *testing.T) {

	rvars := make([]string, 7)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr07U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr07U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 7)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)
	ptrs[3] = new(uint64)
	ptrs[4] = new(uint64)
	ptrs[5] = new(uint64)
	ptrs[6] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr08U64(t *testing.T) {

	rvars := make([]string, 8)
	rvars[0] = "name-0"
	rvars[1] = "name-1"
	rvars[2] = "name-2"
	rvars[3] = "name-3"
	rvars[4] = "name-4"
	rvars[5] = "name-5"
	rvars[6] = "name-6"
	rvars[7] = "name-7"

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr08U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr08U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 8)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)
	ptrs[3] = new(uint64)
	ptrs[4] = new(uint64)
	ptrs[5] = new(uint64)
	ptrs[6] = new(uint64)
	ptrs[7] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr09U64(t *testing.T) {

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

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64, arg08 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr09U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr09U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 9)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)
	ptrs[3] = new(uint64)
	ptrs[4] = new(uint64)
	ptrs[5] = new(uint64)
	ptrs[6] = new(uint64)
	ptrs[7] = new(uint64)
	ptrs[8] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}

func TestFuncAr10U64(t *testing.T) {

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

	fct := func(arg00 uint64, arg01 uint64, arg02 uint64, arg03 uint64, arg04 uint64, arg05 uint64, arg06 uint64, arg07 uint64, arg08 uint64, arg09 uint64) uint64 {
		return 42
	}

	form, err := newFuncAr10U64(rvars, fct)
	if err != nil {
		t.Fatalf("could not create funcAr10U64 formula: %+v", err)
	}

	if got, want := form.RVars(), rvars; !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid rvars: got=%#v, want=%#v", got, want)
	}

	ptrs := make([]interface{}, 10)
	ptrs[0] = new(uint64)
	ptrs[1] = new(uint64)
	ptrs[2] = new(uint64)
	ptrs[3] = new(uint64)
	ptrs[4] = new(uint64)
	ptrs[5] = new(uint64)
	ptrs[6] = new(uint64)
	ptrs[7] = new(uint64)
	ptrs[8] = new(uint64)
	ptrs[9] = new(uint64)

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

	got := form.Func().(func() uint64)()
	if got, want := got, uint64(42); got != want {
		t.Fatalf("invalid output:\ngot= %v (%T)\nwant=%v (%T)", got, got, want, want)
	}
}
