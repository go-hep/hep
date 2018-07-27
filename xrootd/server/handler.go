// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server // import "go-hep.org/x/hep/xrootd/server"

import (
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/login"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

// Handler provides a high-level API for the XRootD server.
// The Handler receives a parsed request and returns a response together with the status
// that will be send via Server to the client.
type Handler interface {
	// Handshake handles the XRootD handshake: http://xrootd.org/doc/dev45/XRdv310.htm#_Toc464248784.
	Handshake() (xrdproto.Marshaler, xrdproto.ResponseStatus)

	// Login handles the XRootD login request: http://xrootd.org/doc/dev45/XRdv310.htm#_Toc464248819.
	Login(sessionID [16]byte, request *login.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus)

	// Protocol handles the XRootD protocol request: http://xrootd.org/doc/dev45/XRdv310.htm#_Toc464248827.
	Protocol(sessionID [16]byte, request *protocol.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus)

	// Dirlist handles the XRootD dirlist request: http://xrootd.org/doc/dev45/XRdv310.htm#_Toc464248815.
	Dirlist(sessionID [16]byte, request *dirlist.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus)

	// CloseSession handles the aborting of user session. This can be used to free some user-related data.
	CloseSession(sessionID [16]byte) error

	// Open handles the XRootD open request: http://xrootd.org/doc/dev45/XRdv310.htm#_Toc464248823.
	Open(sessionID [16]byte, request *open.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus)

	// Close handles the XRootD close request: http://xrootd.org/doc/dev45/XRdv310.htm#_Toc464248813.
	Close(sessionID [16]byte, request *xrdclose.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus)

	// Read handles the XRootD read request: http://xrootd.org/doc/dev45/XRdv310.htm#_Toc464248841.
	Read(sessionID [16]byte, request *read.Request) (xrdproto.Marshaler, xrdproto.ResponseStatus)
}
