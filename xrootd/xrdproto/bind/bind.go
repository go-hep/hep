// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package bind contains the structures describing bind request and response.
package bind // import "go-hep.org/x/hep/xrootd/xrdproto/bind"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3024

// Request holds the bind request parameters.
type Request struct {
	SessionID [16]byte // SessionID is the session identifier returned by login request.
	_         int32
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.SessionID[:])
	wBuffer.Next(4)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.SessionID[:])
	rBuffer.Skip(4)
	return nil
}

// Response is a response for the bind request, which contains the path id.
type Response struct {
	PathID xrdproto.PathID
}

// RespID implements xrdproto.Response.RespID.
func (resp *Response) RespID() uint16 { return RequestID }

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteU8(uint8(o.PathID))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.PathID = xrdproto.PathID(rBuffer.ReadU8())
	return nil
}
