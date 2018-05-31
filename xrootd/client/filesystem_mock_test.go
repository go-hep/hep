// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"net"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/protocol"
	"go-hep.org/x/hep/xrootd/protocol/dirlist"
)

func TestFilesystem_Dirlist_Mock(t *testing.T) {
	path := "/tmp/test"
	response := ".\n0 0 0 0\ntestfile\n0 20 0 10\ntestfile2\n0 21 2 12\x00"

	var want = []protocol.EntryStat{
		{
			Name:        "testfile",
			Size:        20,
			Mtime:       10,
			HasStatInfo: true,
		},
		{
			Name:        "testfile2",
			Size:        21,
			Mtime:       12,
			HasStatInfo: true,
			IsDir:       true,
		},
	}

	var wantRequest = dirlist.Request{
		Options: dirlist.WithStatInfo,
		Path:    path,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		defer cancel()

		data, err := readRequest(conn)
		if err != nil {
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest dirlist.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != dirlist.RequestID {
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", dirlist.RequestID, gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		response := dirlist.Response{[]byte(response)}

		responseHeader := protocol.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: int32(len(response.Data)),
		}

		responseData, err := marshalResponse(responseHeader, response)
		if err != nil {
			t.Fatalf("could not marshal response: %v", err)
		}

		if err := writeResponse(conn, responseData); err != nil {
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := Filesystem{client}
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

func TestFilesystem_Dirlist_Mock_WithoutStatInfo(t *testing.T) {
	path := "/tmp/test"
	response := "testfile\ntestfile2\x00"

	var want = []protocol.EntryStat{
		{
			Name:        "testfile",
			HasStatInfo: false,
		},
		{
			Name:        "testfile2",
			HasStatInfo: false,
		},
	}

	var wantRequest = dirlist.Request{
		Options: dirlist.WithStatInfo,
		Path:    path,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		defer cancel()

		data, err := readRequest(conn)
		if err != nil {
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest dirlist.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != dirlist.RequestID {
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", dirlist.RequestID, gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		response := dirlist.Response{[]byte(response)}

		responseHeader := protocol.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: int32(len(response.Data)),
		}

		responseData, err := marshalResponse(responseHeader, response)
		if err != nil {
			t.Fatalf("could not marshal response: %v", err)
		}

		if err := writeResponse(conn, responseData); err != nil {
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := Filesystem{client}
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
