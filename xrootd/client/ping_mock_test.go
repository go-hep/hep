// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"net"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/ping"
)

func TestSession_Ping_Mock(t *testing.T) {
	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
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

		if gotHeader.RequestID != gotRequest.ReqID() {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", gotRequest.ReqID(), gotHeader.RequestID)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: 0,
		}

		response, err := marshalResponse(responseHeader)
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
		err := client.sessions[client.initialSessionID].Ping(context.Background())
		if err != nil {
			t.Fatalf("invalid ping call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
