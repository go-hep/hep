// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package locate contains the types related to the locate request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 51) for more details.
package locate // import "go-hep.org/x/hep/xrootd/xrdproto/locate"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// locate options.
const (
	AddPeers   = 1 << 0  // AddPeers adds eligible peers to the location output
	Refresh    = 1 << 7  // Refresh updates cached information on the file’s location
	PreferName = 1 << 8  // PreferName indicates a hostname response is preferred
	NoWait     = 1 << 13 // NoWait provides informations as soon as possible
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3027

// Request holds the locate request parameters.
type Request struct {
	Options uint16 // Options to apply when Path is opened
	_       [14]byte
	Path    string // Path of the file to locate
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (*Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(w *xrdenc.WBuffer) error {
	w.WriteU16(o.Options)
	w.Next(14)
	w.WriteStr(o.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(r *xrdenc.RBuffer) error {
	o.Options = r.ReadU16()
	r.Skip(14)
	o.Path = r.ReadStr()
	return nil
}

// Response is the response issued by the server to a locate request.
type Response struct {
	Data []byte
}

// RespID implements xrdproto.Response.RespID.
func (*Response) RespID() uint16 { return RequestID }

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(w *xrdenc.WBuffer) error {
	w.WriteBytes(o.Data)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Response) UnmarshalXrd(r *xrdenc.RBuffer) error {
	o.Data = make([]byte, r.Len())
	r.ReadBytes(o.Data)
	return nil
}
