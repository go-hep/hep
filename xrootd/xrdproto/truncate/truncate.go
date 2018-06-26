// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package truncate contains the structures describing truncate request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 121) for details.
package truncate // import "go-hep.org/x/hep/xrootd/xrdproto/truncate"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3028

// Request holds truncate request parameters.
// Either the Handle or the Path should be specified to identify the file.
type Request struct {
	Handle xrdfs.FileHandle
	Size   int64
	_      [4]uint8
	Path   string
}

// MarshalXrd implements xrdproto.Marshaler.
func (req Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(req.Handle[:])
	wBuffer.WriteI64(req.Size)
	wBuffer.Next(4)
	wBuffer.WriteStr(req.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (req *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(req.Handle[:])
	req.Size = rBuffer.ReadI64()
	rBuffer.Skip(4)
	req.Path = rBuffer.ReadStr()
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }

// Opaque implements xrdproto.FilepathRequest.Opaque.
func (req *Request) Opaque() string {
	return xrdproto.Opaque(req.Path)
}

// SetOpaque implements xrdproto.FilepathRequest.SetOpaque.
func (req *Request) SetOpaque(opaque string) {
	// Opaque is only required if path was specified.
	if len(req.Path) == 0 {
		return
	}
	xrdproto.SetOpaque(&req.Path, opaque)
}
