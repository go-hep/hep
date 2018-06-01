// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package handshake contains the structures describing request and response
// for handshake request (see XRootD specification).
package handshake // import "go-hep.org/x/hep/xrootd/protocol/handshake"

import (
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/protocol"
)

// Response is a response for the handshake request,
// which contains protocol version and server type.
type Response struct {
	ProtocolVersion int32
	ServerType      protocol.ServerType
}

// Request holds the handshake request parameters.
type Request struct {
	Reserved1 int32
	Reserved2 int32
	Reserved3 int32
	Reserved4 int32
	Reserved5 int32
}

// NewRequest forms a Request that comply with the XRootD protocol v3.1.0.
func NewRequest() Request {
	return Request{0, 0, 0, 4, 2012}
}

func (req *Request) MarshalXrd() ([]byte, error) {
	var enc xrdenc.Encoder
	enc.WriteI32(req.Reserved1)
	enc.WriteI32(req.Reserved2)
	enc.WriteI32(req.Reserved3)
	enc.WriteI32(req.Reserved4)
	enc.WriteI32(req.Reserved5)
	return enc.Bytes(), nil
}
