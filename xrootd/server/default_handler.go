// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server // import "go-hep.org/x/hep/xrootd/server"

import (
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/handshake"
	"go-hep.org/x/hep/xrootd/xrdproto/login"
	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
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
func (h *defaultHandler) CloseSession(sessionID [16]byte) {}
