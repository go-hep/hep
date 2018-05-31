// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dirlist contains the structures describing request and response
// for dirlist request used to obtain the contents of a directory.
package dirlist // import "go-hep.org/x/hep/xrootd/protocol/dirlist"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3004

// Response is a response for the dirlist request,
// which contains a byte array containing encoded response.
// The format (if stat information is supported by the server) is:
// ".\n"
// "0 0 0 0\n"
// "dirname\n"
// "id size flags modtime\n"
// ...
// 0
// In case that the server doesn't support returning the stat information, the format is:
// "dirname\n"
// ...
// 0
// See xrootd protocol specification, page 45 for further details.
type Response struct {
	Data []byte
}

func (resp Response) MarshalXrd() ([]byte, error) {
	return resp.Data, nil
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

func (req *Request) ReqID() uint16 { return RequestID }

func (req *Request) MarshalXrd() ([]byte, error) {
	var enc xrdenc.Encoder
	enc.WriteReserved(15)
	enc.WriteU8(byte(req.Options))
	enc.WriteStr(req.Path)
	return enc.Bytes(), nil
}

func (req *Request) UnmarshalXrd(data []byte) error {
	dec := xrdenc.NewDecoder(data)
	dec.Skip(15)
	req.Options = RequestOptions(dec.ReadU8())
	req.Path = dec.ReadStr()
	return nil
}
