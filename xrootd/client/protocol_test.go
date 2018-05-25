// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build xrootd_test_with_server

package client

import (
	"context"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/protocol"
)

func testClient_Protocol(t *testing.T, addr string) {
	var want = ProtocolInfo{
		BinaryProtocolVersion: 784,
		ServerType:            protocol.DataServer,
		IsServer:              true,
	}

	client, err := NewClient(context.Background(), addr)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	got, err := client.Protocol(context.Background())
	if err != nil {
		t.Fatalf("invalid protocol call: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Client.Protocol()\ngot = %v\nwant = %v", got, want)
	}

	client.Close()
}

func TestClient_Protocol(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testClient_Protocol(t, addr)
		})
	}
}
