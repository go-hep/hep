// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package verifyw contains the structures describing verifyw request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 124) for details.
package verifyw // import "go-hep.org/x/hep/xrootd/xrdproto/verifyw"

import (
	"encoding/binary"
	"hash/crc32"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3026

// Type identifies the checksum algorithm used.
type Type uint8

const (
	NoCRC Type = iota // NoCRC identifies that no crc is used.
	CRC32             // CRC#@ identifies that 32-bit crc is used.
)

// Request holds verifyw request parameters.
type Request struct {
	Handle       xrdfs.FileHandle
	Offset       int64
	PathID       uint8
	Verification Type
	_            [2]uint8
	Data         []uint8
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.Handle[:])
	wBuffer.WriteI64(o.Offset)
	wBuffer.WriteU8(o.PathID)
	wBuffer.WriteU8(uint8(o.Verification))
	wBuffer.Next(2)
	wBuffer.WriteLen(len(o.Data))
	wBuffer.WriteBytes(o.Data)

	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.Handle[:])
	o.Offset = rBuffer.ReadI64()
	o.PathID = rBuffer.ReadU8()
	o.Verification = Type(rBuffer.ReadU8())
	rBuffer.Skip(2)
	o.Data = make([]uint8, rBuffer.ReadLen())
	rBuffer.ReadBytes(o.Data)
	return nil
}

// NewRequestCRC32 forms a Request with crc32 verification according to provided parameters.
func NewRequestCRC32(handle xrdfs.FileHandle, offset int64, data []uint8) *Request {
	req := &Request{Handle: handle, Offset: offset, Verification: CRC32}
	crc := crc32.ChecksumIEEE(data)
	crcData := make([]uint8, 4, 4+len(data))
	binary.BigEndian.PutUint32(crcData, crc)
	req.Data = append(crcData, data...)
	return req
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }
