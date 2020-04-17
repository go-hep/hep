// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package query contains the types related to the query request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 79) for more details.
package query // import "go-hep.org/x/hep/xrootd/xrdproto/query"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// Query parameters.
const (
	Stats          = 1  // Query server statistics
	Prepare        = 2  // Query prepare status
	Checksum       = 3  // Query file checksum
	XAttr          = 4  // Query file extended attributes
	Space          = 5  // Query server logical space statistics
	CancelChecksum = 6  // Query file checksum cancellation
	Config         = 7  // Query server configuration
	Visa           = 8  // Query file visa attributes
	Opaque1        = 16 // Query implementation-dependent information
	Opaque2        = 32 // Query implementation-dependent information
	Opaque3        = 64 // Query implementation-dependent information
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3001

// Request holds the query request parameters.
type Request struct {
	Query  uint16
	_      [2]byte
	Handle xrdfs.FileHandle
	_      [8]byte
	Args   []byte
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (*Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(w *xrdenc.WBuffer) error {
	w.WriteU16(o.Query)
	w.Next(2)
	w.WriteBytes(o.Handle[:])
	w.Next(8)
	w.WriteI32(int32(len(o.Args)))
	w.WriteBytes(o.Args)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(r *xrdenc.RBuffer) error {
	o.Query = r.ReadU16()
	r.Skip(2)
	r.ReadBytes(o.Handle[:])
	r.Skip(8)
	n := r.ReadI32()
	if n > 0 {
		o.Args = make([]byte, n)
		r.ReadBytes(o.Args)
	}
	return nil
}

// Response is the response issued by the server to a query request.
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
