// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package endsess contains the types related to the endsess request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 47) for more details.
package endsess // import "go-hep.org/x/hep/xrootd/xrdproto/endsess"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3023

// Request holds the endsess request parameters.
type Request struct {
	SessionID [16]byte
	_         int32
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (*Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(w *xrdenc.WBuffer) error {
	w.WriteBytes(o.SessionID[:])
	w.WriteI32(0)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(r *xrdenc.RBuffer) error {
	r.ReadBytes(o.SessionID[:])
	_ = r.ReadI32()
	return nil
}
