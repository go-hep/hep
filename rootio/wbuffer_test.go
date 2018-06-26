// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"strings"
	"testing"
)

func TestWBuffer_WriteBool(t *testing.T) {
	data := make([]byte, 20)
	wbuf := NewWBuffer(data, nil, 0, nil)
	want := true
	wbuf.WriteBool(want)
	rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
	got := rbuf.ReadBool()
	if got != want {
		t.Fatalf("Invalid value. got:%v, want:%v", got, want)
	}
}

func TestWBuffer_WriteString_Long(t *testing.T) {
	data := make([]byte, 520)
	wbuf := NewWBuffer(data, nil, 0, nil)
	want := strings.Repeat("=", 512)
	wbuf.WriteString(want)
	rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
	got := rbuf.ReadString()
	if got != want {
		t.Fatalf("Invalid value.\ngot: %q\nwant:%q", got, want)
	}
}

func TestWBuffer_WriteString_Short(t *testing.T) {
	data := make([]byte, 520)
	wbuf := NewWBuffer(data, nil, 0, nil)
	want := strings.Repeat("=", 254)
	wbuf.WriteString(want)
	rbuf := NewRBuffer(wbuf.w.p, nil, 0, nil)
	got := rbuf.ReadString()
	if got != want {
		t.Fatalf("Invalid value.\ngot: %q\nwant:%q", got, want)
	}
}
