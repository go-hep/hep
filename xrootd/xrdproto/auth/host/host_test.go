// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package host_test

import (
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto/auth"
	"go-hep.org/x/hep/xrootd/xrdproto/auth/host"
)

func TestAuthHost(t *testing.T) {
	hauth := host.Auth{Hostname: "example.org"}
	if got, want := hauth.Provider(), "host"; got != want {
		t.Fatalf("invalid auth type: got=%q, want=%q", got, want)
	}

	req, err := hauth.Request(nil)
	if err != nil {
		t.Fatalf("got err=%v", err)
	}

	want := &auth.Request{Type: host.Type, Credentials: "host\000example.org\000"}
	if *want != *req {
		t.Fatalf("invalid request:\ngot= %#v\nwant=%#v", req, want)
	}
}
