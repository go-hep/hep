// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !ci
// +build !ci

package xrootd

import (
	"context"
	"testing"
)

func testNewSubSession(t *testing.T, addr string) {
	session, err := newSession(context.Background(), addr, "gopher", "", nil)
	if err != nil {
		t.Fatalf("could not create initialSession: %v", err)
	}
	defer session.Close()

	subSession, err := newSubSession(context.Background(), session)
	if err != nil {
		t.Fatalf("could not create subSession: %v", err)
	}

	if subSession.pathID == 0 {
		t.Fatalf("incorrect subSession.pathID value of 0 was received")
	}

	session.subs[subSession.pathID] = subSession
}

func TestNewSubSession(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testNewSubSession(t, addr)
		})
	}
}
