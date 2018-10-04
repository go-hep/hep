// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package prepare_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/client"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/prepare"
)

func TestRequest(t *testing.T) {
	for _, want := range []prepare.Request{
		{
			Options:  prepare.Stage,
			Priority: 2,
			Port:     8080,
			Paths:    []string{"/foo", "/foo/bar", `C:\\Users\Me`},
		},
		{
			Options:  prepare.Stage | prepare.NoErrors,
			Priority: 0,
			Paths:    []string{"/foo"},
		},
		{
			Options:  prepare.Stage | prepare.NoErrors,
			Priority: 0,
			Paths:    []string{},
		},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got prepare.Request
			)

			if want.ReqID() != prepare.RequestID {
				t.Fatalf("invalid request ID: got=%d want=%d", want.ReqID(), prepare.RequestID)
			}

			if want.ShouldSign() {
				t.Fatalf("invalid")
			}

			err = want.MarshalXrd(w)
			if err != nil {
				t.Fatal(err)
			}

			r := xrdenc.NewRBuffer(w.Bytes())
			err = got.UnmarshalXrd(r)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(got, want) {
				t.Fatalf(
					"round-trip failed:\ngot = %#v\nwant= %#v\n",
					got, want,
				)
			}
		})
	}
}

func TestResponse(t *testing.T) {
	for _, want := range []prepare.Response{
		{Data: []byte{}},
		{Data: []byte("1234")},
		{Data: []byte("")},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got prepare.Response
			)

			if want.RespID() != prepare.RequestID {
				t.Fatalf("invalid response ID: got=%d want=%d", want.RespID(), prepare.RequestID)
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

func Example() {
	bkg := context.Background()

	cl, err := client.NewClient(bkg, "ccxrootdgotest.in2p3.fr:9001", "gopher")
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}
	defer cl.Close()

	var (
		resp prepare.Response
		req  = prepare.Request{
			Options: prepare.Stage | prepare.Notify,
			Paths:   []string{"/tmp/dir1/file1.txt"},
		}
	)

	// staging request

	id, err := cl.Send(bkg, &resp, &req)
	if err != nil {
		log.Fatalf("stage request error: %v", err)
	}
	fmt.Printf("sess:   %s\n", id)
	fmt.Printf("stage:  %q\n", resp.Data[:12]) // Locator ID

	// cancel staging request

	locid := append([]byte(nil), resp.Data...)
	req = prepare.Request{
		Options: prepare.Cancel,
		Paths:   []string{string(locid)},
	}

	id, err = cl.Send(bkg, &resp, &req)
	if err != nil {
		log.Fatalf("cancel request error: %v", err)
	}
	fmt.Printf("cancel: %q\n", resp.Data)

	// Output:
	// sess:   ccxrootdgotest.in2p3.fr:9001
	// stage:  "23297f000001"
	// cancel: ""
}

var (
	_ xrdproto.Request  = (*prepare.Request)(nil)
	_ xrdproto.Response = (*prepare.Response)(nil)
)
