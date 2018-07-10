// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"encoding/binary"
	"hash/crc32"
	"net"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/internal/xrdenc"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	"go-hep.org/x/hep/xrootd/xrdproto/sync"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
	"go-hep.org/x/hep/xrootd/xrdproto/verifyw"
	"go-hep.org/x/hep/xrootd/xrdproto/write"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

func TestFile_Close_Mock(t *testing.T) {
	t.Parallel()

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
		file := file{fs: client.FS().(*fileSystem), handle: handle, sessionID: client.initialSessionID}

		err := file.Close(context.Background())
		if err != nil {
			t.Fatalf("invalid close call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFile_Sync_Mock(t *testing.T) {
	t.Parallel()

	handle := xrdfs.FileHandle{1, 2, 3, 4}
	wantRequest := sync.Request{Handle: handle}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest sync.Request
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
		file := file{fs: client.FS().(*fileSystem), handle: handle, sessionID: client.initialSessionID}

		err := file.Sync(context.Background())
		if err != nil {
			t.Fatalf("invalid sync call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFile_ReadAt_Mock(t *testing.T) {
	t.Parallel()

	handle := xrdfs.FileHandle{1, 2, 3, 4}
	want := []byte("Hello XRootD.\n")
	askLength := int32(len(want) + 4)

	wantRequest := read.Request{Handle: handle, Offset: 1, Length: askLength, OptionalArgs: &read.OptionalArgs{PathID: 0}}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest read.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if gotHeader.RequestID != wantRequest.ReqID() {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", wantRequest.ReqID(), gotHeader.RequestID)
		}

		if !reflect.DeepEqual(*gotRequest.OptionalArgs, *wantRequest.OptionalArgs) {
			cancel()
			t.Fatalf("optional args do not match:\ngot = %v\nwant = %v", *gotRequest.OptionalArgs, *wantRequest.OptionalArgs)
		}

		gotRequest.OptionalArgs = nil
		wantRequest.OptionalArgs = nil

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: 5,
			Status:     xrdproto.OkSoFar,
		}

		responseData, err := marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response header: %v", err)
		}

		responseData = append(responseData, want[:5]...)

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}

		responseHeader = xrdproto.ResponseHeader{
			StreamID:   gotHeader.StreamID,
			DataLength: int32(len(want) - 5),
			Status:     xrdproto.Ok,
		}

		responseData, err = marshalResponse(responseHeader)
		if err != nil {
			cancel()
			t.Fatalf("could not marshal response header: %v", err)
		}

		responseData = append(responseData, want[5:]...)

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		file := file{fs: client.FS().(*fileSystem), handle: handle, sessionID: client.initialSessionID}
		got := make([]uint8, askLength)

		n, err := file.ReadAt(got, 1)
		if err != nil {
			t.Fatalf("invalid read call: %v", err)
		}
		if n != len(want) {
			t.Fatalf("read count does not match:\ngot = %v\nwant = %v", n, len(want))
		}

		if !reflect.DeepEqual(got[:n], want) {
			t.Fatalf("read data does not match:\ngot = %v\nwant = %v", got[:n], want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFile_WriteAt_Mock(t *testing.T) {
	t.Parallel()

	handle := xrdfs.FileHandle{1, 2, 3, 4}
	want := []byte("Hello XRootD.\n")

	wantRequest := write.Request{Handle: handle, Offset: 1, Data: want}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest write.Request
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

		responseHeader := xrdproto.ResponseHeader{StreamID: gotHeader.StreamID}

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
		file := file{fs: client.FS().(*fileSystem), handle: handle, sessionID: client.initialSessionID}

		n, err := file.WriteAt(want, 1)
		if err != nil {
			t.Fatalf("invalid write call: %v", err)
		}
		if n != len(want) {
			t.Fatalf("write count does not match:\ngot = %v\nwant = %v", n, len(want))
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFile_Truncate_Mock(t *testing.T) {
	t.Parallel()

	var (
		handle         = xrdfs.FileHandle{1, 2, 3, 4}
		wantSize int64 = 10
	)

	wantRequest := truncate.Request{Handle: handle, Size: wantSize}

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

		if gotHeader.RequestID != wantRequest.ReqID() {
			cancel()
			t.Fatalf("invalid request id was specified:\nwant = %d\ngot = %d\n", wantRequest.ReqID(), gotHeader.RequestID)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		responseHeader := xrdproto.ResponseHeader{
			StreamID: gotHeader.StreamID,
			Status:   xrdproto.Ok,
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
		file := file{fs: client.FS().(*fileSystem), handle: handle, sessionID: client.initialSessionID}

		err := file.Truncate(context.Background(), wantSize)
		if err != nil {
			t.Fatalf("invalid truncate call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFile_Stat_Mock(t *testing.T) {
	t.Parallel()

	handle := xrdfs.FileHandle{0, 1, 2, 3}

	var want = &xrdfs.EntryStat{
		EntrySize:   20,
		Mtime:       10,
		HasStatInfo: true,
	}

	var wantRequest = stat.Request{FileHandle: handle}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest stat.Request
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

		response := stat.DefaultResponse{EntryStat: *want}
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
			t.Fatalf("could not marshal response: %v", err)
		}
		responseData = append(responseData, wBuffer.Bytes()...)

		if err := writeResponse(conn, responseData); err != nil {
			cancel()
			t.Fatalf("invalid write: %s", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		var fs = client.FS().(*fileSystem)
		file := file{fs: fs, handle: handle, sessionID: client.initialSessionID}
		got, err := file.Stat(context.Background())
		if err != nil {
			t.Fatalf("invalid stat call: %v", err)
		}
		if !reflect.DeepEqual(&got, want) {
			t.Fatalf("stat info does not match:\ngot = %v\nwant = %v", &got, want)
		}
		if !reflect.DeepEqual(file.Info(), want) {
			t.Fatalf("stat info does not match:\nfile.Info() = %v\nwant = %v", file.Info(), want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFile_VerifyWriteAt_Mock(t *testing.T) {
	t.Parallel()

	handle := xrdfs.FileHandle{1, 2, 3, 4}
	data := []byte("Hello XRootD.\n")
	crc := crc32.ChecksumIEEE(data)
	crcData := make([]uint8, 4, 4+len(data))
	binary.BigEndian.PutUint32(crcData, crc)
	crcData = append(crcData, data...)

	wantRequest := verifyw.Request{Handle: handle, Offset: 1, Data: crcData, Verification: verifyw.CRC32}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := readRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest verifyw.Request
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

		responseHeader := xrdproto.ResponseHeader{StreamID: gotHeader.StreamID}

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
		file := file{fs: client.FS().(*fileSystem), handle: handle, sessionID: client.initialSessionID}

		err := file.VerifyWriteAt(context.Background(), data, 1)
		if err != nil {
			t.Fatalf("invalid verifyw call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
