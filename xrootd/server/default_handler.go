// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server // import "go-hep.org/x/hep/xrootd/server"

import (
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/handshake"
	"go-hep.org/x/hep/xrootd/xrdproto/login"
	"go-hep.org/x/hep/xrootd/xrdproto/mv"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	"go-hep.org/x/hep/xrootd/xrdproto/sync"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
	"go-hep.org/x/hep/xrootd/xrdproto/write"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

// defaultHandler implements Handler with some general functionality added.
// Any unimplemented request returns InvalidRequest error.
type defaultHandler struct {
}

// Default returns the defaultHandler implementing Handler with some general functionality added.
// Any unimplemented request returns InvalidRequest error.
func Default() Handler {
	return &defaultHandler{}
}

// Login implements Handler.Login.
func (h *defaultHandler) Login(sessionID [16]byte, request *login.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	return &login.Response{SessionID: sessionID}, xrdproto.Ok
}

// Protocol implements Handler.Protocol.
func (h *defaultHandler) Protocol(sessionID [16]byte, request *protocol.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := &protocol.Response{BinaryProtocolVersion: 0x310, Flags: protocol.IsServer}
	return resp, xrdproto.Ok
}

// Dirlist implements Handler.Dirlist.
func (h *defaultHandler) Dirlist(sessionID [16]byte, request *dirlist.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Dirlist request is not implemented"}
	return resp, xrdproto.Error
}

// Handshake implements Handler.Handshake.
func (*defaultHandler) Handshake() (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := handshake.Response{ProtocolVersion: 0x310, ServerType: xrdproto.DataServer}
	return &resp, xrdproto.Ok
}

// CloseSession implements Handler.CloseSession.
func (h *defaultHandler) CloseSession(sessionID [16]byte) error { return nil }

// Open implements Handler.Open.
func (h *defaultHandler) Open(sessionID [16]byte, request *open.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Open request is not implemented"}
	return resp, xrdproto.Error
}

// Close implements Handler.Close.
func (h *defaultHandler) Close(sessionID [16]byte, request *xrdclose.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Close request is not implemented"}
	return resp, xrdproto.Error
}

// Read implements Handler.Read.
func (h *defaultHandler) Read(sessionID [16]byte, request *read.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Read request is not implemented"}
	return resp, xrdproto.Error
}

// Write implements Handler.Write.
func (h *defaultHandler) Write(sessionID [16]byte, request *write.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Write request is not implemented"}
	return resp, xrdproto.Error
}

// Stat implements Handler.Stat.
func (h *defaultHandler) Stat(sessionID [16]byte, request *stat.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Stat request is not implemented"}
	return resp, xrdproto.Error
}

// Sync implements Handler.Sync.
func (h *defaultHandler) Sync(sessionID [16]byte, request *sync.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Sync request is not implemented"}
	return resp, xrdproto.Error
}

// Truncate implements Handler.Truncate.
func (h *defaultHandler) Truncate(sessionID [16]byte, request *truncate.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Truncate request is not implemented"}
	return resp, xrdproto.Error
}

// Rename implements Handler.Rename.
func (h *defaultHandler) Rename(sessionID [16]byte, request *mv.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus) {
	resp := xrdproto.ServerError{Code: xrdproto.InvalidRequest, Message: "Rename request is not implemented"}
	return resp, xrdproto.Error
}
