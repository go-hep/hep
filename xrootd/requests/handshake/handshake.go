// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package handshake // import "go-hep.org/x/hep/xrootd/requests/handshake"

//go:generate stringer -type=ServerType

type ServerType int32

const (
	LoadBalancingServer ServerType = iota
	DataServer          ServerType = iota
)

type Response struct {
	ProtocolVersion int32
	ServerType      ServerType
}

type Request struct {
	Reserved1 int32
	Reserved2 int32
	Reserved3 int32
	Reserved4 int32
	Reserved5 int32
}

func NewRequest() Request {
	return Request{0, 0, 0, 4, 2012}
}
