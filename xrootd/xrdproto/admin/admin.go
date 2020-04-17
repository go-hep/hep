// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package admin contains the types related to the admin request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 37) for more details.
package admin // import "go-hep.org/x/hep/xrootd/xrdproto/admin"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3020

// Request holds the admin request parameters.
type Request struct {
	_   [16]byte
	Req string
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (*Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(w *xrdenc.WBuffer) error {
	w.Next(16)
	w.WriteStr(o.Req)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(r *xrdenc.RBuffer) error {
	r.Skip(16)
	o.Req = r.ReadStr()
	return nil
}

// Response is the response issued by the server to an admin request.
type Response struct {
	Data []byte
}

// RespID implements xrdproto.Response.RespID.
func (*Response) RespID() uint16 { return RequestID }

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(w *xrdenc.WBuffer) error {
	w.WriteBytes(o.Data)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Response) UnmarshalXrd(r *xrdenc.RBuffer) error {
	o.Data = make([]byte, r.Len())
	r.ReadBytes(o.Data)
	return nil
}
