// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package handshake contains the structures describing request and response
// for handshake request (see XRootD specification).
package handshake // import "go-hep.org/x/hep/xrootd/protocol/handshake"

// ServerType is the general server type kept for compatibility
// with 2.0 protocol version (see xrootd protocol specification v3.1.0, p. 5).
type ServerType int32

const (
	// LoadBalancingServer indicates whether this is a load-balancing server.
	LoadBalancingServer ServerType = iota
	// DataServer indicates whether this is a data server.
	DataServer ServerType = iota
)

// Response is a response for the handshake request,
// which contains protocol version and server type.
type Response struct {
	ProtocolVersion int32
	ServerType      ServerType
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
