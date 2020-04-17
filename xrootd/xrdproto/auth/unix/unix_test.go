// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unix_test

import (
	"testing"

	"go-hep.org/x/hep/xrootd/xrdproto/auth"
	"go-hep.org/x/hep/xrootd/xrdproto/auth/unix"
)

func TestUnix(t *testing.T) {
	v := unix.Auth{User: "gopher", Group: "golang"}
	if got, want := v.Provider(), "unix"; got != want {
		t.Fatalf("invalid provider: got=%q, want=%q", got, want)
	}

	if got, want := v.Provider(), string(unix.Type[:]); got != want {
		t.Fatalf("invalid provider type: got=%q, want=%q", got, want)
	}

	want := &auth.Request{Type: unix.Type, Credentials: "unix\000gopher golang\000"}

	got, err := v.Request(nil)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}

	if got.Credentials != want.Credentials {
		t.Fatalf("invalid credentials:\ngot= %q\nwant=%q\n", got.Credentials, want.Credentials)
	}

	if got.Type != want.Type {
		t.Fatalf("invalid type: got=%q, want=%q", got.Type, want.Type)
	}
}
