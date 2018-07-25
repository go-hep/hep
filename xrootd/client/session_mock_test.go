// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"net"
	"testing"
	"time"

	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/ping"
)

func TestSession_WaitResponse(t *testing.T) {
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

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Wait, xrdproto.WaitResponse{Duration: time.Second})
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}

		responseTime := time.Now()

		data, err = xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		sleepTime := time.Now().Sub(responseTime)
		if sleepTime < time.Second/2 {
			t.Errorf("client should wait around 1 second before re-issuing request, waited %v", sleepTime)
		}

		gotHeader, err = unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, xrdproto.WaitResponse{Duration: time.Second})
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
