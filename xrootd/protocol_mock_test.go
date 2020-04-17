// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd

import (
	"context"
	"net"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
)

func TestSession_Protocol_WithSecurityInfo(t *testing.T) {
	var protocolVersion int32 = 0x310

	var want = protocol.Response{
		BinaryProtocolVersion: protocolVersion,
		HasSecurityInfo:       true,
		SecurityLevel:         xrdproto.Pedantic,
		SecurityOverrides:     []xrdproto.SecurityOverride{{RequestIndex: 1, RequestLevel: xrdproto.SignNeeded}},
		Flags:                 protocol.IsServer | protocol.IsManager | protocol.IsMeta | protocol.IsProxy | protocol.IsSupervisor,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
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

		if gotRequest.ClientProtocolVersion != protocolVersion {
			cancel()
			t.Fatalf("invalid client protocol version was specified:\nwant = %d\ngot = %d\n", protocolVersion, gotRequest.ClientProtocolVersion)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, want)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		client.sessions[client.initialSessionID].protocolVersion = protocolVersion
		got, err := client.sessions[client.initialSessionID].Protocol(context.Background())
		if err != nil {
			t.Fatalf("invalid protocol call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("protocol info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
