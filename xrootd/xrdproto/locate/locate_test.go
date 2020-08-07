// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package locate_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/locate"
)

func Example() {
	bkg := context.Background()

	cl, err := xrootd.NewClient(bkg, "ccxrootdgotest.in2p3.fr:9001", "gopher")
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}
	defer cl.Close()

	var (
		resp locate.Response
		req  = locate.Request{
			Options: locate.PreferName | locate.Refresh | locate.AddPeers,
			Path:    "/tmp/dir1/file1.txt",
		}
	)

	id, err := cl.Send(bkg, &resp, &req)
	if err != nil {
		log.Fatalf("locate request error: %v", err)
	}
	fmt.Printf("sess: %s\n", id)
	fmt.Printf("locate: %q\n", bytes.TrimRight(resp.Data, "\x00"))

	// Output:
	// sess: ccxrootdgotest.in2p3.fr:9001
	// locate: "Sr[::172.17.0.61]:9001"
}

func TestRequest(t *testing.T) {
	for _, want := range []locate.Request{
		{
			Options: locate.AddPeers,
			Path:    "/tmp/dir1/file1.txt",
		},
		{
			Options: locate.AddPeers | locate.Refresh,
			Path:    "/tmp/dir1/file1.txt?foo=bar",
		},
		{
			Options: locate.NoWait | locate.PreferName,
			Path:    "*file1.txt?foo=bar",
		},
		{
			Options: locate.NoWait | locate.PreferName,
			Path:    "*",
		},
	} {
		t.Run(want.Path, func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got locate.Request
			)

			if want.ReqID() != locate.RequestID {
				t.Fatalf("invalid request ID: got=%d want=%d", want.ReqID(), locate.RequestID)
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
	for _, want := range []locate.Response{
		{Data: []byte("")},
		{Data: []byte("1234")},
		{Data: []byte("localhost:1024")},
		{Data: []byte("0.0.0.0:1024")},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got locate.Response
			)

			if want.RespID() != locate.RequestID {
				t.Fatalf("invalid response ID: got=%d want=%d", want.RespID(), locate.RequestID)
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
	_ xrdproto.Request  = (*locate.Request)(nil)
	_ xrdproto.Response = (*locate.Response)(nil)
)
