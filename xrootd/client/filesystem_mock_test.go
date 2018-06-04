// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"net"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
)

func TestFileSystem_Dirlist_Mock(t *testing.T) {
	path := "/tmp/test"
	response := ".\n0 0 0 0\ntestfile\n0 20 0 10\ntestfile2\n0 21 2 12\x00"

	var want = []xrdfs.EntryStat{
		{
			EntryName:   "testfile",
			EntrySize:   20,
			Mtime:       10,
			HasStatInfo: true,
		},
		{
			EntryName:   "testfile2",
			EntrySize:   21,
			Mtime:       12,
			HasStatInfo: true,
			Flags:       xrdfs.StatIsDir,
		},
	}

	var wantRequest = dirlist.Request{
		Options: dirlist.WithStatInfo,
		Path:    path,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest dirlist.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != dirlist.RequestID {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", dirlist.RequestID, gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: int32(len(response)),
		}

		responseData, err := marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response: %v", err)
		}
		responseData = append(responseData, []byte(response)...)

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := client.FS()
		got, err := fs.Dirlist(context.Background(), path)
		if err != nil {
			t.Fatalf("invalid dirlist call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("dirlist info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_Dirlist_Mock_WithoutStatInfo(t *testing.T) {
	path := "/tmp/test"
	response := "testfile\ntestfile2\x00"

	var want = []xrdfs.EntryStat{
		{
			EntryName:   "testfile",
			HasStatInfo: false,
		},
		{
			EntryName:   "testfile2",
			HasStatInfo: false,
		},
	}

	var wantRequest = dirlist.Request{
		Options: dirlist.WithStatInfo,
		Path:    path,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest dirlist.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != dirlist.RequestID {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", dirlist.RequestID, gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: int32(len(response)),
		}

		responseData, err := marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response: %v", err)
		}
		responseData = append(responseData, []byte(response)...)

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := fileSystem{client}
		got, err := fs.Dirlist(context.Background(), path)
		if err != nil {
			t.Fatalf("invalid dirlist call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("dirlist info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
