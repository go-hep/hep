// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rm contains the structures describing rm request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 105) for details.
package rm // import "go-hep.org/x/hep/xrootd/xrdproto/rm"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3014

// Request holds rm request parameters, such as the file path.
type Request struct {
	_    [16]byte
	Path string
}

// MarshalXrd implements xrdproto.Marshaler.
func (req Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.Next(16)
	wBuffer.WriteStr(req.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (req *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.Skip(16)
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
	xrdproto.SetOpaque(&req.Path, opaque)
}
