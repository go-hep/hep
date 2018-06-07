// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"net"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/rm"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
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

func testFileSystem_Open_Mock(t *testing.T, wantFileHandle xrdfs.FileHandle, wantFileCompression *xrdfs.FileCompression, wantFileInfo *xrdfs.EntryStat) {
	path := "/tmp/test"

	var wantRequest = open.Request{
		Path:    path,
		Mode:    xrdfs.OpenModeOtherRead,
		Options: xrdfs.OpenOptionsOpenRead,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest open.Request
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

		response := open.Response{
			FileHandle:  wantFileHandle,
			Compression: wantFileCompression,
			Stat:        wantFileInfo,
		}

		var wBuffer xrdenc.WBuffer
		err = response.MarshalXrd(&wBuffer)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response: %v", err)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: int32(len(wBuffer.Bytes())),
		}

		responseData, err := marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response header: %v", err)
		}
		responseData = append(responseData, wBuffer.Bytes()...)

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := client.FS()
		gotFile, err := fs.Open(context.Background(), path, xrdfs.OpenModeOtherRead, xrdfs.OpenOptionsOpenRead)
		if err != nil {
			t.Fatalf("invalid open call: %v", err)
		}
		// FIXME: consider calling defer gotFile.Close(context.Background()).

		if !reflect.DeepEqual(gotFile.Handle(), wantFileHandle) {
			t.Errorf("Filesystem.Open()\ngotFile.Handle() = %v\nwantFileHandle = %v", gotFile.Handle(), wantFileHandle)
		}

		if !reflect.DeepEqual(gotFile.Compression(), wantFileCompression) {
			t.Errorf("Filesystem.Open()\ngotFile.Compression() = %v\nwantFileCompression = %v", gotFile.Compression(), wantFileCompression)
		}
		if !reflect.DeepEqual(gotFile.Info(), wantFileInfo) {
			t.Errorf("Filesystem.Open()\ngotFile.Info() = %v\nwantFileInfo = %v", gotFile.Info(), wantFileInfo)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_Open_Mock(t *testing.T) {
	testCases := []struct {
		name        string
		handle      xrdfs.FileHandle
		compression *xrdfs.FileCompression
		stat        *xrdfs.EntryStat
	}{
		{"WithoutCompressionAndStat", xrdfs.FileHandle{0, 0, 0, 0}, nil, nil},
		{"WithEmptyCompression", xrdfs.FileHandle{0, 0, 0, 0}, &xrdfs.FileCompression{}, nil},
		{"WithCompression", xrdfs.FileHandle{0, 0, 0, 0}, &xrdfs.FileCompression{10, [4]byte{'t', 'e', 's', 't'}}, nil},
		{"WithStat", xrdfs.FileHandle{0, 0, 0, 0}, &xrdfs.FileCompression{}, &xrdfs.EntryStat{HasStatInfo: true, EntrySize: 10}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testFileSystem_Open_Mock(t, tc.handle, tc.compression, tc.stat)
		})
	}
}

func TestFileSystem_RemoveFile_Mock(t *testing.T) {
	var (
		path        = "/tmp/test"
		wantRequest = rm.Request{Path: path}
	)

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest rm.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != gotRequest.ReqID() {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", gotRequest.ReqID(), gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		responseHeader := xrdproto.ResponseHeader{StreamID: gotHeader.StreamID}

		responseData, err := marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response: %v", err)
		}

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := fileSystem{client}
		err := fs.RemoveFile(context.Background(), path)
		if err != nil {
			t.Fatalf("invalid rm call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_Truncate_Mock(t *testing.T) {
	var (
		path              = "/tmp/test"
		wantSize    int64 = 10
		wantRequest       = truncate.Request{Path: path, Size: wantSize}
	)

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest truncate.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != gotRequest.ReqID() {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", gotRequest.ReqID(), gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		responseHeader := xrdproto.ResponseHeader{StreamID: gotHeader.StreamID}

		responseData, err := marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response: %v", err)
		}

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := fileSystem{client}
		err := fs.Truncate(context.Background(), path, wantSize)
		if err != nil {
			t.Fatalf("invalid truncate call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
