// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package errstack // import "go-hep.org/x/hep/fwk/utils/errstack"

import (
	"fmt"
	"runtime"
)

const (
	// The size of the stack trace buffer for reading.  Stack traces will not be
	// longer than this.
	StackTraceSize = 4096
)

// Error wraps an error and attaches the stack trace at the moment of emission.
type Error struct {
	Err   error
	Stack []byte
}

// Error returns the error string
func (err *Error) Error() string {
	return fmt.Sprintf("%v\nstack: %v\n", err.Err, string(err.Stack))
}

// New returns a new Error.
// If err is already a errstack.Error, it's a no-op.
func New(err error) error {
	if err == nil {
		return nil
	}

	if err, ok := err.(*Error); ok {
		return err
	}

	var buf [StackTraceSize]byte
	n := runtime.Stack(buf[:], false)
	stack := make([]byte, n)
	copy(stack, buf[:n])

	return &Error{
		Err:   err,
		Stack: stack,
	}
}

// Newf creates a new Error.
func Newf(format string, args ...interface{}) error {

	var buf [StackTraceSize]byte
	n := runtime.Stack(buf[:], false)
	stack := make([]byte, n)
	copy(stack, buf[:n])

	return &Error{
		Err:   fmt.Errorf(format, args...),
		Stack: stack,
	}
}
