// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build xrootd_test_with_server

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"testing"
)

var testClientAddrs = []string{"0.0.0.0:9001"}

func TestNewClient(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testNewClient(t, addr)
		})
	}
}

func testNewClient(t *testing.T, addr string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	client, err := NewClient(ctx, addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	if err := client.Close(); err != nil {
		t.Fatalf("could not close client: %v", err)
	}
}

func BenchmarkNewClient(b *testing.B) {
	for _, addr := range testClientAddrs {
		b.Run(addr, func(b *testing.B) {
			benchmarkNewClient(b, addr)
		})
	}
}

func benchmarkNewClient(b *testing.B, addr string) {
	for i := 0; i < b.N; i++ {
		client, err := NewClient(context.Background(), addr, "gopher")
		if err != nil {
			b.Fatalf("could not create client: %v", err)
		}
		if err := client.Close(); err != nil {
			b.Fatalf("could not close client: %v", err)
		}
	}
}
