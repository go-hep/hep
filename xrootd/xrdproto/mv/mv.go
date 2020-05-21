// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mv contains the structures describing mv request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 106) for details.
package mv // import "go-hep.org/x/hep/xrootd/xrdproto/mv"

import (
	"errors"
	"fmt"
	"strings"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3009

// Request holds mv request parameters.
type Request struct {
	_       [14]byte
	OldPath string
	NewPath string
}

// MarshalXrd implements xrdproto.Marshaler.
func (req Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.Next(14)
	wBuffer.WriteU16(uint16(len(req.OldPath)))
	wBuffer.WriteLen(len(req.OldPath) + len(req.NewPath) + 1)
	wBuffer.WriteBytes([]byte(req.OldPath))
	wBuffer.WriteBytes([]byte{' '})
	wBuffer.WriteBytes([]byte(req.NewPath))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (req *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.Skip(14)
	fromLen := int(rBuffer.ReadU16())
	paths := rBuffer.ReadStr()
	if fromLen >= len(paths) {
		return fmt.Errorf("xrootd: wrong mv request. Want fromLen < %d, got %d", len(paths)-1, fromLen)
	}
	if fromLen == 0 {
		fromLen = strings.Index(paths, " ")
		if fromLen == -1 {
			return errors.New("xrootd: wrong mv request. Want paths to be separated by ' ', none found")
		}
	}
	req.OldPath = string(paths[:fromLen])
	req.NewPath = string(paths[fromLen+1:])
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }

// Opaque implements xrdproto.FilepathRequest.Opaque.
func (req *Request) Opaque() string {
	return xrdproto.Opaque(req.NewPath)
}

// SetOpaque implements xrdproto.FilepathRequest.SetOpaque.
func (req *Request) SetOpaque(opaque string) {
	xrdproto.SetOpaque(&req.NewPath, opaque)
}
