// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dirlist contains the structures describing request and response
// for dirlist request used to obtain the contents of a directory.
package dirlist // import "go-hep.org/x/hep/xrootd/xrdproto/dirlist"

import (
	"bytes"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3004

// Response is a response for the dirlist request,
// which contains a slice of entries containing the entry name and the entry stat info.
type Response struct {
	Entries []xrdfs.EntryStat
}

// RespID implements xrdproto.Response.RespID.
func (resp *Response) RespID() uint16 { return RequestID }

// MarshalXrd implements xrdproto.Marshaler.
func (o Response) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	// TODO: implement
	panic(errors.New("xrootd: MarshalXrd is not implemented"))
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler
// When stat information is supported by the server, the format is
//  ".\n"
//  "0 0 0 0\n"
//  "dirname\n"
//  "id size flags modtime\n"
//  ...
//  0
// Otherwise, the format is the following:
//  "dirname\n"
//  ...
//  0
// See xrootd protocol specification, page 45 for further details.
func (o *Response) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	if rBuffer.Len() == 0 {
		return nil
	}

	data := bytes.TrimRight(rBuffer.Bytes(), "\x00")
	lines := bytes.Split(data, []byte{'\n'})

	// FIXME(sbinet): drop the extra call to bytes.Equal when
	//  https://github.com/xrootd/xrootd/issues/739
	// is fixed or clarified.
	if !(bytes.HasPrefix(data, []byte(".\n0 0 0 0\n")) || bytes.Equal(data, []byte(".\n0 0 0 0"))) {
		// That means that the server doesn't support returning stat information.
		o.Entries = make([]xrdfs.EntryStat, len(lines))
		for i, v := range lines {
			o.Entries[i] = xrdfs.EntryStat{EntryName: string(v)}
		}
		return nil
	}

	if len(lines)%2 != 0 {
		return errors.Errorf("xrootd: wrong response size for the dirlist request: want even number of lines, got %d", len(lines))
	}

	lines = lines[2:]
	o.Entries = make([]xrdfs.EntryStat, len(lines)/2)

	for i := 0; i < len(lines); i += 2 {
		var rbuf = xrdenc.NewRBuffer(lines[i+1])
		err := o.Entries[i/2].UnmarshalXrd(rbuf)
		if err != nil {
			return err
		}
		o.Entries[i/2].EntryName = string(lines[i])
	}

	return nil
}

// Request holds the dirlist request parameters.
type Request struct {
	_       [15]byte
	Options RequestOptions
	Path    string
}

// RequestOptions specifies what should be returned as part of response.
type RequestOptions byte

const (
	None         RequestOptions = 0 // None specifies that no addition information except entry names is required.
	WithStatInfo RequestOptions = 2 // WithStatInfo specifies that stat information for every entry is required.
)

// NewRequest forms a Request according to provided path.
func NewRequest(path string) *Request {
	return &Request{Options: WithStatInfo, Path: path}
}

// ReqID implements xrdproto.Request.ReqID.
func (req *Request) ReqID() uint16 { return RequestID }

// ShouldSign implements xrdproto.Request.ShouldSign.
func (req *Request) ShouldSign() bool { return false }

// MarshalXrd implements xrdproto.Marshaler.
func (o Request) MarshalXrd(wBuffer *xrdenc.WBuffer) error {
	wBuffer.Next(15)
	wBuffer.WriteU8(byte(o.Options))
	wBuffer.WriteStr(o.Path)
	return nil
}

// UnmarshalXrd implements xrdproto.Unmarshaler.
func (o *Request) UnmarshalXrd(rBuffer *xrdenc.RBuffer) error {
	rBuffer.Skip(15)
	o.Options = RequestOptions(rBuffer.ReadU8())
	o.Path = rBuffer.ReadStr()
	return nil
}
