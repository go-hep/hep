// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sync contains the structures describing sync request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 119) for details.
package sync // import "go-hep.org/x/hep/xrootd/xrdproto/sync"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3016

// Request holds sync request parameters, such as the file handle.
type Request struct {
	Handle xrdfs.FileHandle
	_      [12]byte
	_      int32
}

// MarshalXrd implements xrdproto.Marshaler
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.Handle[:])
	wBuffer.Next(16)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.Handle[:])
	rBuffer.Skip(16)
	return nil
}

// ReqID implements xrdproto.Request.ReqID
func (req *Request) ReqID() uint16 { return RequestID }
