// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"testing"
)

func testSession_Ping(t *testing.T, addr string) {
	client, err := newSession(context.Background(), addr, "gopher", "", nil)
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	err = client.Ping(context.Background())
	if err != nil {
		t.Fatalf("invalid ping call: %v", err)
	}
}

func TestSession_Ping(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testSession_Ping(t, addr)
		})
	}
}

func BenchmarkSession_Ping(b *testing.B) {
	for _, addr := range testClientAddrs {
		b.Run(addr, func(b *testing.B) {
			benchmarkSession_Ping(b, addr)
		})
	}
}

func benchmarkSession_Ping(b *testing.B, addr string) {
	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		b.Fatalf("could not create client: %v", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := client.sessions[client.initialSessionID].Ping(context.Background())
		if err != nil {
			b.Errorf("could not ping: %v", err)
		}
	}
	if err := client.Close(); err != nil {
		b.Fatalf("could not close client: %v", err)
	}
}
