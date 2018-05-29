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

		flags := protocol.IsManager | protocol.IsServer | protocol.IsMeta | protocol.IsProxy | protocol.IsSupervisor

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: protocol.GeneralResponseLength + protocol.SecurityInfoLength + protocol.SecurityOverrideLength,
		}

		protocolResponse := protocol.GeneralResponse{protocolVersion, flags}

		protocolSecurityInfo := protocol.SecurityInfo{
			SecurityOptions:       protocol.None,
			SecurityLevel:         protocol.Pedantic,
			SecurityOverridesSize: 1,
		}

		securityOverride := protocol.SecurityOverride{1, protocol.SignNeeded}

		response, err := marshalResponse(responseHeader, protocolResponse, protocolSecurityInfo, securityOverride)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response: %v", err)
		}

		if err := writeResponse(conn, response); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	var want = ProtocolInfo{
		BinaryProtocolVersion: protocolVersion,
		ServerType:            xrdproto.DataServer,
		IsManager:             true,
		IsServer:              true,
		IsMeta:                true,
		IsProxy:               true,
		IsSupervisor:          true,
		SecurityLevel:         protocol.Pedantic,
		SecurityOverrides:     []protocol.SecurityOverride{{1, protocol.SignNeeded}},
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
