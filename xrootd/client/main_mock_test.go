// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"net"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/internal/mux"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/signing"
)

var testClientAddrs []string

func testClientWithMockServer(serverFunc func(cancel func(), conn net.Conn), clientFunc func(cancel func(), client *Client)) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, conn := net.Pipe()
	defer server.Close()
	defer conn.Close()

	client := &Client{cancel: cancel, sessions: make(map[string]*session)}
	session := &session{cancel: cancel, ctx: ctx, conn: conn, mux: mux.New(), requests: make(map[xrdproto.StreamID]pendingRequest), client: client, signRequirements: signing.Default()}
	client.initialSessionID = "test.org:1234"
	client.sessions[client.initialSessionID] = session
	defer client.Close()

	go serverFunc(func() { client.Close() }, server)
	go session.consume()

	clientFunc(cancel, client)
}

func unmarshalRequest(data []byte, request xrdproto.Request) (xrdproto.RequestHeader, error) {
	var header xrdproto.RequestHeader
	rBuffer := xrdenc.NewRBuffer(data)

	if err := header.UnmarshalXrd(rBuffer); err != nil {
		return xrdproto.RequestHeader{}, err
	}
	if header.RequestID != request.ReqID() {
		return xrdproto.RequestHeader{}, errors.Errorf("xrootd: unexpected request id was specified:\nwant = %d\ngot = %d\n", request.ReqID(), header.RequestID)
	}
	if err := request.UnmarshalXrd(rBuffer); err != nil {
		return xrdproto.RequestHeader{}, err
	}

	return header, nil
}
