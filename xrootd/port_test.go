// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd

import "testing"

func TestParseAddr(t *testing.T) {
	for _, tc := range []struct {
		addr string
		want string
	}{
		{addr: "localhost:1094", want: "localhost:1094"},
		{addr: ":1094", want: ":1094"},
		{addr: "0.0.0.0:1094", want: "0.0.0.0:1094"},
		{addr: "0.0.0.0", want: "0.0.0.0:1094"},
		{addr: "192.168.0.1", want: "192.168.0.1:1094"},
		{addr: "ccxrootdgotest.in2p3.fr", want: "ccxrootdgotest.in2p3.fr:1094"},
		{addr: "ccxrootdgotest.in2p3.fr:8080", want: "ccxrootdgotest.in2p3.fr:8080"},
	} {
		t.Run(tc.addr, func(t *testing.T) {
			got := parseAddr(tc.addr)
			if got != tc.want {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}
