// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd

import (
	"context"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto/protocol"
)

func testSession_Protocol(t *testing.T, addr string) {
	want := protocol.Response{BinaryProtocolVersion: 784, Flags: protocol.IsServer}

	session, err := newSession(context.Background(), addr, "gopher", "", nil)
	if err != nil {
		t.Fatalf("could not create initialSession: %v", err)
	}
	defer session.Close()

	got, err := session.Protocol(context.Background())
	if err != nil {
		t.Fatalf("invalid protocol call: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("session.Protocol()\ngot = %v\nwant = %v", got, want)
	}
}

func TestSession_Protocol(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testSession_Protocol(t, addr)
		})
	}
}
