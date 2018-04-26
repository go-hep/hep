// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ping // import "go-hep.org/x/hep/xrootd/requests/ping"

const RequestID uint16 = 3011

type Request struct {
	Reserved1 [16]byte
	Reserved2 int32
}

func NewRequest() Request {
	return Request{[16]byte{}, 0}
}
