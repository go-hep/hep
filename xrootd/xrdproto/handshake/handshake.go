// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package handshake contains the structures describing request and response
// for handshake request (see XRootD specification).
package handshake // import "go-hep.org/x/hep/xrootd/xrdproto/handshake"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// Response is a response for the handshake request,
// which contains protocol version and server type.
type Response struct {
	ProtocolVersion int32
	ServerType      xrdproto.ServerType
}

// MarshalXrd implements xrdproto.Marshaler
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteI32(o.ProtocolVersion)
	wBuffer.WriteI32(int32(o.ServerType))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.ProtocolVersion = rBuffer.ReadI32()
	o.ServerType = xrdproto.ServerType(rBuffer.ReadI32())
	return nil
}

// RequestLength is the length of the Request in bytes.
const RequestLength = 20

// Request holds the handshake request parameters.
type Request [5]int32

// NewRequest forms a Request that complies with the XRootD protocol v3.1.0.
func NewRequest() Request {
	return Request{0, 0, 0, 4, 2012}
}

// MarshalXrd implements xrdproto.Marshaler
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteI32(o[0])
	wBuffer.WriteI32(o[1])
	wBuffer.WriteI32(o[2])
	wBuffer.WriteI32(o[3])
	wBuffer.WriteI32(o[4])
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o[0] = rBuffer.ReadI32()
	o[1] = rBuffer.ReadI32()
	o[2] = rBuffer.ReadI32()
	o[3] = rBuffer.ReadI32()
	o[4] = rBuffer.ReadI32()
	return nil
}
