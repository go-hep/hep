// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bind // import "go-hep.org/x/hep/xrootd/requests/bind"

const RequestID uint16 = 3024

type Response struct {
	PathID byte
}

type Request struct {
	SessionID [16]byte
	Reserved  int32
}

func NewRequest(sessionID [16]byte) Request {
	return Request{sessionID, 0}
}
