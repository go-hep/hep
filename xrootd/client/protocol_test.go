// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
)

func testClient_Protocol(t *testing.T, addr string) {
	want := protocol.Response{BinaryProtocolVersion: 784, Flags: protocol.IsServer}

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	got, err := client.Protocol(context.Background())
	if err != nil {
		t.Fatalf("invalid protocol call: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Client.Protocol()\ngot = %v\nwant = %v", got, want)
	}
}

func TestClient_Protocol(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testClient_Protocol(t, addr)
		})
	}
}
