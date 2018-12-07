// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package query_test

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"
	"testing"

	"go-hep.org/x/hep/xrootd"
	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/query"
)

func TestRequest(t *testing.T) {
	for _, want := range []query.Request{
		{
			Query: query.Config,
			Args:  []byte("bind_max chksum"),
		},
		{
			Query:  query.Visa,
			Handle: [4]byte{1, 2, 3, 4},
		},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got query.Request
			)

			if want.ReqID() != query.RequestID {
				t.Fatalf("invalid request ID: got=%d want=%d", want.ReqID(), query.RequestID)
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
	for _, want := range []query.Response{
		{
			Data: []byte("1234"),
		},
		{
			Data: []byte("oss.cgroup=public&oss.space=499337216&oss.free=444280832&oss.maxf=444280832&oss.used=55056384&oss.quota=-1\x00"),
		},
	} {
		t.Run("", func(t *testing.T) {
			var (
				err error
				w   = new(xrdenc.WBuffer)
				got query.Response
			)

			if want.RespID() != query.RequestID {
				t.Fatalf("invalid response ID: got=%d want=%d", want.RespID(), query.RequestID)
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

	cl, err := xrootd.NewClient(bkg, "ccxrootdgotest.in2p3.fr:9001", "gopher")
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}
	defer cl.Close()

	fs := cl.FS()
	f, err := fs.Open(bkg, "/tmp/dir1/file1.txt", xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		log.Fatalf("open error: %v", err)
	}
	defer f.Close(bkg)

	var (
		resp query.Response
		req  = query.Request{
			Query:  query.Space,
			Handle: f.Handle(),
			Args:   []byte("/tmp/dir1/file1.txt"),
		}
	)

	id, err := cl.Send(bkg, &resp, &req)
	if err != nil {
		log.Fatalf("space request error: %v", err)
	}
	fmt.Printf("sess: %s\n", id)
	// fmt.Printf("space: %q\n", resp.Data)

	cfg := []string{
		"bind_max",
		"chksum",
		"cid", "cms", "pio_max",
		"readv_ior_max",
		"readv_iov_max",
		"role",
		"sitename",
		"tpc",
		"version",
		"wan_port",
		"wan_window",
		"window",
	}

	req = query.Request{
		Query: query.Config,
		Args:  []byte(strings.Join(cfg, " ")),
	}

	id, err = cl.Send(bkg, &resp, &req)
	if err != nil {
		log.Fatalf("config request error: %v", err)
	}
	for i, v := range strings.Split(strings.TrimRight(string(resp.Data), "\n"), "\n") {
		if v == cfg[i] {
			fmt.Printf("config: %s=N/A\n", v)
			continue
		}
		fmt.Printf("config: %s=%q\n", cfg[i], v)
	}

	// Output:
	// sess: ccxrootdgotest.in2p3.fr:9001
	// config: bind_max="15"
	// config: chksum=N/A
	// config: cid=N/A
	// config: cms="none|"
	// config: pio_max="5"
	// config: readv_ior_max="2097136"
	// config: readv_iov_max="1024"
	// config: role="server"
	// config: sitename=N/A
	// config: tpc=N/A
	// config: version="v4.8.4"
	// config: wan_port=N/A
	// config: wan_window=N/A
	// config: window="87380"
}

var (
	_ xrdproto.Request  = (*query.Request)(nil)
	_ xrdproto.Response = (*query.Response)(nil)
)
