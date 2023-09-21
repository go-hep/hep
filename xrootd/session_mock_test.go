// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"errors"
	"net"
	"os"
	"testing"
	"time"

	"go-hep.org/x/hep/xrootd/internal/mux"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/ping"
	"go-hep.org/x/hep/xrootd/xrdproto/signing"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
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

		sleepTime := time.Since(responseTime)
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

func TestSession_ConnectionAbort(t *testing.T) {
	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest truncate.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
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
	serverFuncForSecondConnection := func(cancel func(), conn net.Conn) {
		_, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Errorf("could not read request: %v", err)
		}
		conn.Close()
	}

	clientFunc := func(cancel func(), client *Client) {
		p1, p2 := net.Pipe()
		go serverFuncForSecondConnection(cancel, p2)
		session := &cliSession{
			cancel:           cancel,
			ctx:              context.Background(),
			conn:             p1,
			mux:              mux.New(),
			requests:         make(map[xrdproto.StreamID]pendingRequest),
			client:           client,
			signRequirements: signing.Default(),
			sessionID:        client.initialSessionID + "2",
			isSub:            true,
		}
		defer session.Close()
		defer p1.Close()
		client.sessions[session.sessionID] = session
		go session.consume()

		f := file{sessionID: session.sessionID, fs: client.FS().(*fileSystem)}
		err := f.Truncate(context.Background(), 0)
		if err != nil {
			t.Fatalf("invalid truncate call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestSessionCloseNil(t *testing.T) {
	var sess *cliSession
	err := sess.Close()
	if !errors.Is(err, os.ErrInvalid) {
		t.Fatalf("invalid error: got=%v, want=%v", err, os.ErrInvalid)
	}
}
