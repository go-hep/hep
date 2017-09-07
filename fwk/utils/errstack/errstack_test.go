// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errstack_test

import (
	"fmt"
	"reflect"
	"testing"

	"go-hep.org/x/hep/fwk/utils/errstack"
)

func TestNewNil(t *testing.T) {
	err := errstack.New(nil)
	if err != nil {
		t.Fatalf("expected a nil error. got=%#v\n", err)
	}
}

func TestNew(t *testing.T) {
	exp := fmt.Errorf("my bad %d", 42)

	err := errstack.New(exp)
	if err == nil {
		t.Fatalf("expected a non-nil error. got=%#v\n", err)
	}

	if err, ok := err.(*errstack.Error); !ok {
		t.Fatalf("expected an *errstack.Error. got=%T\n", err)
	} else {
		if !reflect.DeepEqual(err.Err, exp) {
			t.Fatalf("expected err.Err=%v.\ngot=%v\n", exp, err.Err)
		}
	}
}

func TestNewNoHysteresis(t *testing.T) {
	exp := fmt.Errorf("my bad %d", 42)

	err := errstack.New(exp)
	if err == nil {
		t.Fatalf("expected a non-nil error. got=%#v\n", err)
	}

	errs := err.(*errstack.Error)
	n := len(errs.Stack)

	err = newerr(err)
	if err == nil {
		t.Fatalf("expected a non-nil error. got=%#v\n", err)
	}

	errs = err.(*errstack.Error)
	if n != len(errs.Stack) {
		t.Fatalf("hysteresis detected:\nold-stack=%d\nnew-stack=%d\n",
			n,
			len(errs.Stack),
		)
	}

	err = newerr(errs.Err)
	errs = err.(*errstack.Error)
	if n == len(errs.Stack) {
		t.Fatalf("hysteresis error:\nold-stack=%d\nnew-stack=%d\n%v\n",
			n,
			len(errs.Stack),
			errs,
		)
	}

}

// adds another stack-frame
func newerr(err error) error {
	return errstack.New(err)
}

func TestNewf(t *testing.T) {
	err := errstack.Newf("my bad %d", 42)
	if err == nil {
		t.Fatalf("expected an error. got=%#v\n", err)
	}

	if err, ok := err.(*errstack.Error); !ok {
		t.Fatalf("expected an *errstack.Error. got=%T\n", err)
	} else {
		exp := fmt.Errorf("my bad %d", 42)
		if !reflect.DeepEqual(err.Err, exp) {
			t.Fatalf("expected err.Err=%v.\ngot=%v\n", exp, err.Err)
		}
	}

}
