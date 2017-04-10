// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"io"
)

// Reader is the interface that wraps the basic io.Reader interface
// and adds SIO pointer tagging capabilities.
type Reader interface {
	io.Reader
}

// Writer is the interface that wraps the basic io.Writer interface
// and adds SIO pointer tagging capabilities.
type Writer interface {
	io.Writer
}

// Marshaler is the interface implemented by an object that can marshal
// itself into a binary, sio-compatible, form.
type Marshaler interface {
	MarshalSio(w Writer) error
}

// Unmarshaler is the interface implemented by an object that can
// unmarshal a binary, sio-compatible, representation of itself.
type Unmarshaler interface {
	UnmarshalSio(r Reader) error
}

// Code is the interface implemented by an object that can
// unmarshal and marshal itself from and to a binary, sio-compatible, form.
type Codec interface {
	Marshaler
	Unmarshaler
}
