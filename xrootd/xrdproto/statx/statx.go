// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package statx contains the structures describing request and response for statx request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 113) for details.
// Note that only a limited number of flags is meaningful such as StatIsExecutable, StatIsDir, StatIsOther, StatIsOffline.
package statx // import "go-hep.org/x/hep/xrootd/xrdproto/statx"

import (
	"strings"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3022

// Response is a response for the statx request which contains the information about every requested path.
// Note that only limited number of flags is meaningful such as StatIsExecutable, StatIsDir, StatIsOther, StatIsOffline.
type Response struct {
	StatFlags []xrdfs.StatFlags
}

// RespID implements xrdproto.Response.RespID.
func (resp *Response) RespID() uint16 { return RequestID }

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	for _, x := range o.StatFlags {
		wBuffer.WriteU8(uint8(x))
	}
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.StatFlags = make([]xrdfs.StatFlags, rBuffer.Len())
	for i := range o.StatFlags {
		o.StatFlags[i] = xrdfs.StatFlags(rBuffer.ReadU8())
	}
	return nil
}

// Request holds open request parameters.
type Request struct {
	_     [16]uint8
	Paths string // Paths is the new-line separated path list.
}

// NewRequest forms a Request according to provided paths.
func NewRequest(paths []string) *Request {
	return &Request{Paths: strings.Join(paths, "\n")}
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.Next(16)
	wBuffer.WriteStr(o.Paths)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.Skip(16)
	o.Paths = rBuffer.ReadStr()
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }
