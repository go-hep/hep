// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package open contains the structures describing request and response for open request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 63) for details.
package open // import "go-hep.org/x/hep/xrootd/xrdproto/open"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3010

// Response is a response for the open request,
// which contains the file handle, the compression page size,
// the compression type and the stat information.
type Response struct {
	FileHandle  xrdfs.FileHandle
	Compression *xrdfs.FileCompression
	Stat        *xrdfs.EntryStat
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.FileHandle[:])
	if o.Compression == nil {
		return nil
	}
	if err := o.Compression.MarshalXrd(wBuffer); err != nil {
		return err
	}

	if o.Stat == nil {
		return nil
	}
	if err := o.Stat.MarshalXrd(wBuffer); err != nil {
		return err
	}
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.FileHandle[:])
	if rBuffer.Len() == 0 {
		return nil
	}
	o.Compression = &xrdfs.FileCompression{}
	if err := o.Compression.UnmarshalXrd(rBuffer); err != nil {
		return err
	}
	if rBuffer.Len() == 0 {
		return nil
	}
	o.Stat = &xrdfs.EntryStat{}
	if err := o.Stat.UnmarshalXrd(rBuffer); err != nil {
		return err
	}
	return nil
}

// RespID implements xrdproto.Response.RespID.
func (resp *Response) RespID() uint16 { return RequestID }

// Request holds open request parameters.
type Request struct {
	Mode    xrdfs.OpenMode
	Options xrdfs.OpenOptions
	_       [12]byte
	Path    string
}

// Opaque implements xrdproto.FilepathRequest.Opaque.
func (req *Request) Opaque() string {
	return xrdproto.Opaque(req.Path)
}

// SetOpaque implements xrdproto.FilepathRequest.SetOpaque.
func (req *Request) SetOpaque(opaque string) {
	xrdproto.SetOpaque(&req.Path, opaque)
}

// NewRequest forms a Request according to provided path, mode, and options.
func NewRequest(path string, mode xrdfs.OpenMode, options xrdfs.OpenOptions) *Request {
	return &Request{Mode: mode, Options: options, Path: path}
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteU16(uint16(o.Mode))
	wBuffer.WriteU16(uint16(o.Options))
	wBuffer.Next(12)
	wBuffer.WriteStr(o.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.Mode = xrdfs.OpenMode(rBuffer.ReadU16())
	o.Options = xrdfs.OpenOptions(rBuffer.ReadU16())
	rBuffer.Skip(12)
	o.Path = rBuffer.ReadStr()
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool {
	// According to specification, the open request needs to be signed
	// if any of the following options has been specified.
	return req.Options&xrdfs.OpenOptionsDelete != 0 ||
		req.Options&xrdfs.OpenOptionsNew != 0 ||
		req.Options&xrdfs.OpenOptionsOpenUpdate != 0 ||
		req.Options&xrdfs.OpenOptionsMkPath != 0 ||
		req.Options&xrdfs.OpenOptionsOpenAppend != 0
}

var (
	_ xrdproto.FilepathRequest = (*Request)(nil)
)
