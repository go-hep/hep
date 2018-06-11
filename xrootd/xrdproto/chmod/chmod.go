// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package chmod contains the structures describing chmod request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 106) for details.
package chmod // import "go-hep.org/x/hep/xrootd/xrdproto/chmod"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3002

// Request holds chmod request parameters, such as the directory path and the mode to be applied.
type Request struct {
	_    [14]byte
	Mode xrdfs.OpenMode
	Path string
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.Next(14)
	wBuffer.WriteU16(uint16(o.Mode))
	wBuffer.WriteStr(o.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.Skip(14)
	o.Mode = xrdfs.OpenMode(rBuffer.ReadU16())
	o.Path = rBuffer.ReadStr()
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }
