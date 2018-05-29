// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package handshake contains the structures describing request and response
// for handshake request (see XRootD specification).
package handshake // import "go-hep.org/x/hep/xrootd/protocol/handshake"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/protocol"
)

// Response is a response for the handshake request,
// which contains protocol version and server type.
type Response struct {
	ProtocolVersion int32
	ServerType      protocol.ServerType
}

// MarshalXrd implements xrootd/protocol.Marshaler
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteI32(o.ProtocolVersion)
	wBuffer.WriteI32(int32(o.ServerType))
	return nil
}

// UnmarshalXrd implements xrootd/protocol.Unmarshaler
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.ProtocolVersion = rBuffer.ReadI32()
	o.ServerType = protocol.ServerType(rBuffer.ReadI32())
	return nil
}

// Request holds the handshake request parameters.
type Request struct {
	Reserved1 int32
	Reserved2 int32
	Reserved3 int32
	Reserved4 int32
	Reserved5 int32
}

// NewRequest forms a Request that comply with the XRootD protocol v3.1.0.
func NewRequest() Request {
	return Request{0, 0, 0, 4, 2012}
}

// MarshalXrd implements xrootd/protocol.Marshaler
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteI32(o.Reserved1)
	wBuffer.WriteI32(o.Reserved2)
	wBuffer.WriteI32(o.Reserved3)
	wBuffer.WriteI32(o.Reserved4)
	wBuffer.WriteI32(o.Reserved5)
	return nil
}

// UnmarshalXrd implements xrootd/protocol.Unmarshaler
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.Reserved1 = rBuffer.ReadI32()
	o.Reserved2 = rBuffer.ReadI32()
	o.Reserved3 = rBuffer.ReadI32()
	o.Reserved4 = rBuffer.ReadI32()
	o.Reserved5 = rBuffer.ReadI32()
	return nil
}
