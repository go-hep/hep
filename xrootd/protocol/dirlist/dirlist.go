// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package dirlist contains the structures describing request and response
// for dirlist request used to obtain the contents of a directory.
package dirlist // import "go-hep.org/x/hep/xrootd/protocol/dirlist"

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

// Request holds the dirlist request parameters.
type Request struct {
	// FIXME: Rename Reserved field to _ when automatically generated (un)marshalling will be available.
	Reserved   [15]byte
	Options    RequestOptions
	PathLength int32
	Path       []byte
}

// RequestOptions specifies what should be returned as part of response.
type RequestOptions byte

const (
	None         RequestOptions = 0 // None specifies that no addition information except entry names is required.
	WithStatInfo RequestOptions = 2 // WithStatInfo specifies that stat information for every entry is required.
)

// NewRequest forms a Request according to provided path.
func NewRequest(path string) Request {
	var pathBytes = make([]byte, len(path))
	copy(pathBytes, path)

	return Request{Options: WithStatInfo, PathLength: int32(len(path)), Path: pathBytes}
}
