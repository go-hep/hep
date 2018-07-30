// Copyright 2018 The go-hep Authors. All rights reserved.
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

// Auther is the interface that must be implemented by a security provider.
type Auther interface {
	Provider() string                          // Provider returns the name of the security provider.
	Request(params []string) (*Request, error) // Request forms an authorization Request according to passed parameters.
}
