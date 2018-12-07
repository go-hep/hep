// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrootd // import "go-hep.org/x/hep/xrootd"

import (
	"context"
	"net"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
	"go-hep.org/x/hep/xrootd/xrdproto/chmod"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/mkdir"
	"go-hep.org/x/hep/xrootd/xrdproto/mv"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/rm"
	"go-hep.org/x/hep/xrootd/xrdproto/rmdir"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	"go-hep.org/x/hep/xrootd/xrdproto/statx"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
)

func TestFileSystem_Dirlist_Mock(t *testing.T) {
	t.Parallel()

	path := "/tmp/test"
	want := []xrdfs.EntryStat{
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
	wantRequest := dirlist.Request{
		Options: dirlist.WithStatInfo,
		Path:    path,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
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

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, dirlist.Response{Entries: want, WithStatInfo: true})
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
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
	t.Parallel()

	path := "/tmp/test"

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
		data, err := xrdproto.ReadRequest(conn)
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

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, dirlist.Response{Entries: want, WithStatInfo: false})
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
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
		data, err := xrdproto.ReadRequest(conn)
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

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		response := open.Response{
			FileHandle:  wantFileHandle,
			Compression: wantFileCompression,
			Stat:        wantFileInfo,
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, response)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
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
			t.Errorf("FileSystem.Open()\ngotFile.Handle() = %v\nwantFileHandle = %v", gotFile.Handle(), wantFileHandle)
		}

		if !reflect.DeepEqual(gotFile.Compression(), wantFileCompression) {
			t.Errorf("FileSystem.Open()\ngotFile.Compression() = %v\nwantFileCompression = %v", gotFile.Compression(), wantFileCompression)
		}
		if !reflect.DeepEqual(gotFile.Info(), wantFileInfo) {
			t.Errorf("FileSystem.Open()\ngotFile.Info() = %v\nwantFileInfo = %v", gotFile.Info(), wantFileInfo)
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
			t.Parallel()

			testFileSystem_Open_Mock(t, tc.handle, tc.compression, tc.stat)
		})
	}
}

func TestFileSystem_RemoveFile_Mock(t *testing.T) {
	t.Parallel()

	var (
		path        = "/tmp/test"
		wantRequest = rm.Request{Path: path}
	)

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
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

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
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
	t.Parallel()

	var (
		path              = "/tmp/test"
		wantSize    int64 = 10
		wantRequest       = truncate.Request{Path: path, Size: wantSize}
	)

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
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

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
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

func TestFileSystem_Stat_Mock(t *testing.T) {
	t.Parallel()

	path := "/tmp/test"

	var want = xrdfs.EntryStat{
		EntrySize:   20,
		Mtime:       10,
		HasStatInfo: true,
	}

	var wantRequest = stat.Request{Path: path}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
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

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, stat.DefaultResponse{EntryStat: want})
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := client.FS()
		got, err := fs.Stat(context.Background(), path)
		if err != nil {
			t.Fatalf("invalid stat call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("stat info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_VirtualStat_Mock(t *testing.T) {
	t.Parallel()

	path := "/tmp/test"

	var want = xrdfs.VirtualFSStat{
		NumberRW:           1,
		FreeRW:             100,
		UtilizationRW:      10,
		NumberStaging:      2,
		FreeStaging:        200,
		UtilizationStaging: 20,
	}

	var wantRequest = stat.Request{Path: path, Options: stat.OptionsVFS}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
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

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, stat.VirtualFSResponse{VirtualFSStat: want})
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := client.FS()
		got, err := fs.VirtualStat(context.Background(), path)
		if err != nil {
			t.Fatalf("invalid stat call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("virtual stat info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_Mkdir_Mock(t *testing.T) {
	t.Parallel()

	path := "/tmp/test"
	wantRequest := mkdir.Request{Path: path, Mode: xrdfs.OpenModeOwnerRead | xrdfs.OpenModeOwnerWrite}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest mkdir.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := client.FS()
		err := fs.Mkdir(context.Background(), path, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite)
		if err != nil {
			t.Fatalf("invalid mkdir call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_MkdirAll_Mock(t *testing.T) {
	t.Parallel()

	path := "/tmp/test"
	wantRequest := mkdir.Request{
		Path:    path,
		Mode:    xrdfs.OpenModeOwnerRead | xrdfs.OpenModeOwnerWrite,
		Options: mkdir.OptionsMakePath,
	}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest mkdir.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := client.FS()
		err := fs.MkdirAll(context.Background(), path, xrdfs.OpenModeOwnerRead|xrdfs.OpenModeOwnerWrite)
		if err != nil {
			t.Fatalf("invalid mkdir call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_RemoveDir_Mock(t *testing.T) {
	t.Parallel()

	var (
		path        = "/tmp/test"
		wantRequest = rmdir.Request{Path: path}
	)

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest rmdir.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := fileSystem{client}
		err := fs.RemoveDir(context.Background(), path)
		if err != nil {
			t.Fatalf("invalid rmdir call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_Rename_Mock(t *testing.T) {
	t.Parallel()

	var (
		oldpath     = "/tmp/test1"
		newpath     = "/tmp/test2"
		wantRequest = mv.Request{OldPath: oldpath, NewPath: newpath}
	)

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest mv.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := fileSystem{client}
		err := fs.Rename(context.Background(), oldpath, newpath)
		if err != nil {
			t.Fatalf("invalid mv call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_Chmod_Mock(t *testing.T) {
	t.Parallel()

	var (
		path        = "/tmp/test2"
		perm        = xrdfs.OpenModeOwnerRead | xrdfs.OpenModeOwnerWrite
		wantRequest = chmod.Request{Path: path, Mode: perm}
	)

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest chmod.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, nil)
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := fileSystem{client}
		err := fs.Chmod(context.Background(), path, perm)
		if err != nil {
			t.Fatalf("invalid chmod call: %v", err)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}

func TestFileSystem_Statx_Mock(t *testing.T) {
	paths := []string{"/tmp/test", "/tmp/test2"}
	want := []xrdfs.StatFlags{xrdfs.StatIsDir, xrdfs.StatIsOffline}
	wantRequest := statx.Request{Paths: "/tmp/test\n/tmp/test2"}

	serverFunc := func(cancel func(), conn net.Conn) {
		data, err := xrdproto.ReadRequest(conn)
		if err != nil {
			cancel()
			t.Fatalf("could not read request: %v", err)
		}

		var gotRequest statx.Request
		gotHeader, err := unmarshalRequest(data, &gotRequest)
		if err != nil {
			cancel()
			t.Fatalf("could not unmarshal request: %v", err)
		}

		if !reflect.DeepEqual(gotRequest, wantRequest) {
			cancel()
			t.Fatalf("request info does not match:\ngot = %v\nwant = %v", gotRequest, wantRequest)
		}

		err = xrdproto.WriteResponse(conn, gotHeader.StreamID, xrdproto.Ok, statx.Response{StatFlags: want})
		if err != nil {
			cancel()
			t.Fatalf("could not write response: %v", err)
		}
	}

	clientFunc := func(cancel func(), client *Client) {
		fs := client.FS()
		got, err := fs.Statx(context.Background(), paths)
		if err != nil {
			t.Fatalf("invalid statx call: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("statx info does not match:\ngot = %v\nwant = %v", got, want)
		}
	}

	testClientWithMockServer(serverFunc, clientFunc)
}
