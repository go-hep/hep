// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"errors"
)

var (
	ErrStreamNoRecMarker   = errors.New("sio: no record marker found")
	ErrRecordNoBlockMarker = errors.New("sio: no block marker found")
	ErrBlockConnected      = errors.New("sio: block already connected")

	// ErrBlockShortRead means that the deserialization of a SIO block
	// read too few bytes with regard to what was written out in the
	// Block header length.
	ErrBlockShortRead = errors.New("sio: read too few bytes")

	errPointerIDOverflow = errors.New("sio: pointer id overflow")
)
