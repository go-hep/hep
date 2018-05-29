// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"net"
	"reflect"
	"testing"

	xrdproto "go-hep.org/x/hep/xrootd/protocol"
	"go-hep.org/x/hep/xrootd/protocol/protocol"
)

func TestClient_Protocol_WithSecurityInfo(t *testing.T) {
	var protocolVersion int32 = 0x310

	var want = protocol.Response{
		BinaryProtocolVersion: protocolVersion,
		HasSecurityInfo:       true,
		SecurityLevel:         xrdproto.Pedantic,
		SecurityOverrides:     []xrdproto.SecurityOverride{{1, xrdproto.SignNeeded}},
		Flags:                 protocol.IsServer | protocol.IsManager | protocol.IsMeta | protocol.IsProxy | protocol.IsSupervisor,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest protocol.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != protocol.RequestID {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", protocol.RequestID, gotHeader.RequestID)
		}

		if gotRequest.ClientProtocolVersion != protocolVersion {
			cancel()
			t.Fatalf("invalid client protocol version was specified:\nwant = %d\ngot = %d\n", protocolVersion, gotRequest.ClientProtocolVersion)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: 14 + xrdproto.SecurityOverrideLength,
		}

		response, err := marshalResponse(responseHeader, want)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response: %v", err)
		}

		if err := writeResponse(conn, response); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		client.protocolVersion = protocolVersion
		got, err := client.Protocol(context.Background())
		if err != nil {
			t.Fatalf("invalid protocol call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("protocol info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
