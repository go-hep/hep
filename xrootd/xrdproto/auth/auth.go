// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package auth contains the structures describing auth request.
package auth // import "go-hep.org/x/hep/xrootd/xrdproto/auth"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3000

// Request holds the auth request parameters.
type Request struct {
	_           [12]byte
	Type        [4]byte
	Credentials string
}

// UnixType indicates that unix authentication protocol is used.
var UnixType = [4]byte{'u', 'n', 'i', 'x'}

// NewUnixRequest forms a Request according to provided parameters using unix authentication.
func NewUnixRequest(username, groupname string) *Request {
	return &Request{Type: UnixType, Credentials: "unix\000" + username + " " + groupname + "\000"}
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.Next(12)
	wBuffer.WriteBytes(o.Type[:])
	wBuffer.WriteStr(o.Credentials)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.Skip(12)
	rBuffer.ReadBytes(o.Type[:])
	o.Credentials = rBuffer.ReadStr()
	return nil
}
