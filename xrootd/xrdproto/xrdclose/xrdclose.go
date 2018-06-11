// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xrdclose contains the structures describing request and response for close request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 41) for details.
package xrdclose // import "go-hep.org/x/hep/xrootd/xrdproto/xrdclose"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3003

// Request holds close request parameters, such as
// the file handle and the size, in bytes, that the file
// is to have. The close operation fails and the file
// is erased if it is not of the indicated size. A size of 0
// suppresses the check.
type Request struct {
	Handle xrdfs.FileHandle
	Size   int64
	_      [4]byte
	_      int32
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.Handle[:])
	wBuffer.WriteI64(o.Size)
	wBuffer.Next(8)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.Handle[:])
	o.Size = rBuffer.ReadI64()
	rBuffer.Skip(8)
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }
