// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mkdir contains the structures describing mkdir request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 105) for details.
package mkdir // import "go-hep.org/x/hep/xrootd/xrdproto/mkdir"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3008

// Options are the options to apply when path is created.
type Options uint8

// OptionsMakePath indicates whether directory path
// should be created if it does not already exist.
// When a directory path is created, the directory permission
// specified in Mode is propagated along the newly created path.
const OptionsMakePath Options = 1

// Request holds mkdir request parameters, such as the file path.
type Request struct {
	Options Options
	_       [13]uint8
	Mode    xrdfs.OpenMode
	Path    string
}

// MarshalXrd implements xrdproto.Marshaler.
func (req Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteU8(uint8(req.Options))
	wBuffer.Next(13)
	wBuffer.WriteU16(uint16(req.Mode))
	wBuffer.WriteStr(req.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (req *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	req.Options = Options(rBuffer.ReadU8())
	rBuffer.Skip(13)
	req.Mode = xrdfs.OpenMode(rBuffer.ReadU16())
	req.Path = rBuffer.ReadStr()
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }
