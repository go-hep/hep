// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package prepare contains the types related to the prepare request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 69) for more details.
package prepare // import "go-hep.org/x/hep/xrootd/xrdproto/prepare"

import (
	"strings"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// Prepare request options.
const (
	Cancel   = 1  // Cancel will cancel a prepare request.
	Notify   = 2  // Notify will send a message when the file has been processed.
	NoErrors = 4  // NoErrors will not send a notification for preparation errors.
	Stage    = 8  // Stage will stage the file to disk if it is not online.
	Write    = 16 // Write will prepare the file with write access.
	Colocate = 32 // Colocate will co-locate the staged files, if at all possible.
	Refresh  = 64 // Refresh will refresh the file access time even when location is known.
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3021

// Request holds the prepare request parameters.
type Request struct {
	Options  byte   // Options is a set of flags that apply to each path.
	Priority byte   // Priority the request will have. 0: lowest priority, 3: highest.
	Port     uint16 // UDP port number to which a message is to be sent.
	_        [12]byte
	Paths    []string
}

// MarshalXrd implements xrdproto.Marshaler.
func (req Request) MarshalXrd(w *xrdenc.WBuffer) error {
	w.WriteU8(req.Options)
	w.WriteU8(req.Priority)
	w.WriteU16(req.Port)
	w.Next(12)

	var raw []byte
	switch len(req.Paths) {
	case 0:
		// no-op
	default:
		for i, p := range req.Paths {
			if i > 0 {
				raw = append(raw, '\n')
			}
			raw = append(raw, []byte(p)...)
		}
	}
	w.WriteI32(int32(len(raw)))
	w.WriteBytes(raw)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (req *Request) UnmarshalXrd(r *xrdenc.RBuffer) error {
	req.Options = r.ReadU8()
	req.Priority = r.ReadU8()
	req.Port = r.ReadU16()
	r.Skip(12)
	n := r.ReadI32()
	raw := make([]byte, n)
	r.ReadBytes(raw)
	switch n {
	case 0:
		req.Paths = []string{}
	default:
		req.Paths = strings.Split(string(raw), "\n")
	}
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (*Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (*Request) ShouldSign() bool { return false }

// Response is the response issued by the server to a prepare request.
type Response struct {
	Data []byte
}

// RespID implements xrdproto.Response.RespID.
func (resp *Response) RespID() uint16 { return RequestID }

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
