// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package read contains the structures describing request and response for read request.
// See xrootd protocol specification (http://xrootd.org/doc/dev45/XRdv310.pdf, p. 99) for details.
package read // import "go-hep.org/x/hep/xrootd/xrdproto/read"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"golang.org/x/xerrors"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3013

// Response is a response for the read request, which contains the read data.
type Response struct {
	Data []uint8
}

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.Data)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	src := rBuffer.Len()
	dst := len(o.Data)
	if src > dst {
		o.Data = make([]byte, src)
	}
	n := copy(o.Data, rBuffer.Bytes())
	o.Data = o.Data[:n]
	return nil
}

// RespID implements xrdproto.Response.RespID.
func (resp *Response) RespID() uint16 { return RequestID }

// Request holds read request parameters.
type Request struct {
	Handle       xrdfs.FileHandle
	Offset       int64
	Length       int32
	OptionalArgs *OptionalArgs
}

// Request holds optional read request parameters.
type OptionalArgs struct {
	// PathID is the path id returned by bind request.
	// The response data is sent to this path, if possible.
	PathID xrdproto.PathID
	_      [7]uint8
	// ReadAhead is the slice of the pre-read requests.
	ReadAheads []ReadAhead
}

// MarshalXrd implements xrdproto.Marshaler.
func (o OptionalArgs) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	alen := len(o.ReadAheads)*16 + 8
	wBuffer.WriteLen(alen)
	wBuffer.WriteU8(uint8(o.PathID))
	wBuffer.Next(7)
	for _, x := range o.ReadAheads {
		err := x.MarshalXrd(wBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *OptionalArgs) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	alen := rBuffer.ReadLen()
	o.PathID = xrdproto.PathID(rBuffer.ReadU8())
	rBuffer.Skip(7)
	if alen < 8 || (alen-8)%16 != 0 {
		return xerrors.Errorf("xrootd: invalid alen is specified: should be greater or equal to 8"+
			"and (alen - 8) should be dividable by 16, got: %v", alen)
	}
	if alen <= 8 {
		return nil
	}
	o.ReadAheads = make([]ReadAhead, (alen-8)/16)
	for i := 0; i < len(o.ReadAheads); i++ {
		err := o.ReadAheads[i].UnmarshalXrd(rBuffer)
		if err != nil {
			return err
		}
	}
	return nil
}

// ReadAhead is the pre-read request. It is considered only a hint
// and can be used to schedule the pre-reading of data that will be asked
// in the very near future.
type ReadAhead struct {
	Handle xrdfs.FileHandle
	Length int32
	Offset int64
}

// MarshalXrd implements xrdproto.Marshaler.
func (o ReadAhead) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.Handle[:])
	wBuffer.WriteI32(o.Length)
	wBuffer.WriteI64(o.Offset)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *ReadAhead) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.Handle[:])
	o.Length = rBuffer.ReadI32()
	o.Offset = rBuffer.ReadI64()
	return nil
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.WriteBytes(o.Handle[:])
	wBuffer.WriteI64(o.Offset)
	wBuffer.WriteI32(o.Length)
	if o.OptionalArgs == nil {
		wBuffer.WriteLen(0)
		return nil
	}
	return o.OptionalArgs.MarshalXrd(wBuffer)
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.ReadBytes(o.Handle[:])
	o.Offset = rBuffer.ReadI64()
	o.Length = rBuffer.ReadI32()
	if rBuffer.Len() > 4 {
		o.OptionalArgs = &OptionalArgs{}
		return o.OptionalArgs.UnmarshalXrd(rBuffer)
	}
	alen := rBuffer.ReadLen()
	if alen == 0 {
		return nil
	}
	return xerrors.Errorf("xrootd: no data is passed after alen of %d", alen)
}

// PathID implements xrdproto.DataRequest.PathID.
func (o *Request) PathID() xrdproto.PathID {
	if o.OptionalArgs == nil {
		return 0
	}
	return o.OptionalArgs.PathID
}

// PathID implements xrdproto.DataRequest.SetPathID.
func (o *Request) SetPathID(pathID xrdproto.PathID) {
	if o.OptionalArgs == nil {
		o.OptionalArgs = &OptionalArgs{PathID: pathID}
		return
	}
	o.OptionalArgs.PathID = pathID
}

// PathID implements xrdproto.DataRequest.Direction.
func (o *Request) Direction() xrdproto.DataRequestDirection {
	return xrdproto.DataRequestRead
}

// PathID implements xrdproto.DataRequest.PathData.
func (o *Request) PathData() []byte {
	return nil
}

var (
	_ xrdproto.DataRequest = (*Request)(nil)
)
