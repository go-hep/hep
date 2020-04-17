// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package stat contains the structures describing request and response for stat request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 113) for details.
package stat // import "go-hep.org/x/hep/xrootd/xrdproto/stat"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3017

// DefaultResponse is a response for the stat request which contains stat information such as:
// the OS-dependent identifier, the size of the data, the entry attributes and the modification time.
type DefaultResponse struct {
	xrdfs.EntryStat
}

// VirtualFSResponse is a response for the stat request
// which contains virtual file system stat information such as:
// nrw - the number of nodes that can provide read/write access,
// frw - the size of the largest contiguous area of r/w free space,
// urw - the percent utilization of the partition represented by frw,
// nstg - the number of nodes that can provide staging access,
// fstg - the size of the largest contiguous area of staging free space,
// ustg - the percent utilization of the partition represebted by fstg,
type VirtualFSResponse struct {
	xrdfs.VirtualFSStat
}

// RespID implements xrdproto.Response.RespID.
func (resp *VirtualFSResponse) RespID() uint16 { return RequestID }

// RespID implements xrdproto.Response.RespID.
func (resp *DefaultResponse) RespID() uint16 { return RequestID }

// Options are stat processing options.
type Options uint8

const (
	OptionsVFS Options = 1 // OptionsVFS indicates that virtual file system information is requested.
)

// Request holds open request parameters.
type Request struct {
	Options    Options
	_          [11]uint8
	FileHandle xrdfs.FileHandle
	Path       string
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteU8(uint8(o.Options))
	wBuffer.Next(11)
	wBuffer.WriteBytes(o.FileHandle[:])
	wBuffer.WriteStr(o.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	o.Options = Options(rBuffer.ReadU8())
	rBuffer.Skip(11)
	rBuffer.ReadBytes(o.FileHandle[:])
	o.Path = rBuffer.ReadStr()
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
