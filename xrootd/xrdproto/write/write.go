// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package write contains the structures describing write request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 124) for details.
package write // import "go-hep.org/x/hep/xrootd/xrdproto/write"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3019

// Request holds write request parameters.
type Request struct {
	Handle xrdfs.FileHandle
	Offset int64
	PathID uint8
	_      [3]uint8
	Data   []uint8
}

// MarshalXrd implements xrdproto.Marshaler.
func (req Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(req.Handle[:])
	wBuffer.WriteI64(req.Offset)
	wBuffer.WriteU8(req.PathID)
	wBuffer.Next(3)
	wBuffer.WriteLen(len(req.Data))
	wBuffer.WriteBytes(req.Data)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (req *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(req.Handle[:])
	req.Offset = rBuffer.ReadI64()
	req.PathID = rBuffer.ReadU8()
	rBuffer.Skip(3)
	req.Data = make([]uint8, rBuffer.ReadLen())
	rBuffer.ReadBytes(req.Data)
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }
