// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"net"
	"os"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/login"
)

func TestSession_Login_Mock(t *testing.T) {
	username := "gopher"
	token := "token"

	var usernameBytes [8]byte
	copy(usernameBytes[:], username)

	var want = login.Response{
		SessionID:           [16]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		SecurityInformation: []byte("&P=unix"),
	}

	var wantRequest = login.Request{
		Pid:          int32(os.Getpid()),
		Username:     usernameBytes,
		Capabilities: 4,
		Token:        []byte(token),
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest login.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant= %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, want)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		got, err := client.sessions[client.initialSessionID].Login(context.Background(), username, token)
		if err != nil {
			t.Fatalf("invalid login call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("login info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
