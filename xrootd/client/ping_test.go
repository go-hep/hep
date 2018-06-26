// Copyright 2018 The go-hep Authors.  All rights reserved.
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
