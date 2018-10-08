// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package admin_test

import (
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/admin"
)

func TestRequest(t *testing.T) {
	for _, want := range []admin.Request{
		{Req: ""},
		{Req: "hello"},
		{Req: "sudo ls"},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got admin.Request
			)

			if want.ReqID() != admin.RequestID {
				t.Fatalf("invalid request ID: got=%d want=%d", want.ReqID(), admin.RequestID)
			}

			if want.ShouldSign() {
				t.Fatalf("invalid")
			}

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatalf("could not marshal request: %v", err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatalf("could not unmarshal request: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed:\ngot = %#v\nwant= %#v\n", got, want)
			}
		})
	}
}

func TestResponse(t *testing.T) {
	for _, want := range []admin.Response{
		{Data: []byte("")},
		{Data: []byte("1234")},
		{Data: []byte("localhost:1024")},
		{Data: []byte("0.0.0.0:1024")},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got admin.Response
			)

			if want.RespID() != admin.RequestID {
				t.Fatalf("invalid response ID: got=%d want=%d", want.RespID(), admin.RequestID)
			}

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatalf("could not marshal response: %v", err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatalf("could not unmarshal response: %v", err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip failed:\ngot = %#v\nwant= %#v\n", got, want)
			}
		})
	}
}

var (
	_ xrdproto.Request  = (*admin.Request)(nil)
	_ xrdproto.Response = (*admin.Response)(nil)
)
