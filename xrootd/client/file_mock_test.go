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
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

func TestFile_Close_Mock(t *testing.T) {
	handle := xrdfs.FileHandle{1, 2, 3, 4}
	wantRequest := xrdclose.Request{Handle: handle}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest xrdclose.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != wantRequest.ReqID() {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", wantRequest.ReqID(), gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: 0,
		}

		responseData, err := marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response header: %v", err)
		}

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		file := file{fs: client.FS().(*fileSystem), handle: handle}

		err := file.Close(context.Background())
		if err != nil {
			t.Fatalf("invalid close call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
