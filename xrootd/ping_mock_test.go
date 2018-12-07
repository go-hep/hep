// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"net"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/ping"
)

func TestSession_Ping_Mock(t *testing.T) {
	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest ping.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		err := client.sessions[client.initialSessionID].Ping(context.Background())
		if err != nil {
			t.Fatalf("invalid ping call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
