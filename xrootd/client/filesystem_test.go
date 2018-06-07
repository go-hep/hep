// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client

import (
	"context"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdfs"
)

func testFileSystem_Dirlist(t *testing.T, addr string) {
	var want = []xrdfs.EntryStat{
		{
			EntryName:   "file1.txt",
			HasStatInfo: true,
			ID:          60129606914,
			EntrySize:   0,
			Mtime:       1528218208,
			Flags:       xrdfs.StatIsReadable | xrdfs.StatIsWritable,
		},
	}

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	got, err := fs.Dirlist(context.Background(), "/tmp/dir1")
	if err != nil {
		t.Fatalf("invalid protocol call: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("FileSystem.Dirlist()\ngot = %v\nwant = %v", got, want)
	}
}

func TestFileSystem_Dirlist(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFileSystem_Dirlist(t, addr)
		})
	}
}

func testFileSystem_Open(t *testing.T, addr string, options xrdfs.OpenOptions, wantFileHandle xrdfs.FileHandle, wantFileCompression *xrdfs.FileCompression, wantFileInfo *xrdfs.EntryStat) {
	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	gotFile, err := fs.Open(context.Background(), "/tmp/dir1/file1.txt", xrdfs.OpenModeOtherRead, options)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer gotFile.Close(context.Background())

	if !reflect.DeepEqual(gotFile.Handle(), wantFileHandle) {
		t.Errorf("Filesystem.Open()\ngotFile.Handle() = %v\nwantFileHandle = %v", gotFile.Handle(), wantFileHandle)
	}

	if !reflect.DeepEqual(gotFile.Compression(), wantFileCompression) {
		// TODO: Remove this workaround when fix for https://github.com/xrootd/xrootd/issues/721 will be released.
		skippedDefaultCompressionValue := reflect.DeepEqual(wantFileCompression, &xrdfs.FileCompression{}) && gotFile.Compression() == nil
		if !skippedDefaultCompressionValue {
			t.Errorf("Filesystem.Open()\ngotFile.Compression() = %v\nwantFileCompression = %v", gotFile.Compression(), wantFileCompression)
		}
	}

	if !reflect.DeepEqual(gotFile.Info(), wantFileInfo) {
		t.Errorf("Filesystem.Open()\ngotFile.Info() = %v\nwantFileInfo = %v", gotFile.Info(), wantFileInfo)
	}
}

func TestFileSystem_Open(t *testing.T) {
	emptyCompression := xrdfs.FileCompression{}
	entryStat := &xrdfs.EntryStat{
		HasStatInfo: true,
		ID:          60129606914,
		Mtime:       1528218208,
		Flags:       xrdfs.StatIsWritable | xrdfs.StatIsReadable,
	}

	testCases := []struct {
		name        string
		options     xrdfs.OpenOptions
		handle      xrdfs.FileHandle
		compression *xrdfs.FileCompression
		info        *xrdfs.EntryStat
	}{
		{"WithoutCompressionAndStat", xrdfs.OpenOptionsOpenRead, xrdfs.FileHandle{0, 0, 0, 0}, nil, nil},
		{"WithCompression", xrdfs.OpenOptionsOpenRead | xrdfs.OpenOptionsCompress, xrdfs.FileHandle{0, 0, 0, 0}, &emptyCompression, nil},
		{"WithStat", xrdfs.OpenOptionsOpenRead | xrdfs.OpenOptionsReturnStatus, xrdfs.FileHandle{0, 0, 0, 0}, &emptyCompression, entryStat},
	}

	for _, addr := range testClientAddrs {
		for _, tc := range testCases {
			t.Run(addr+"/"+tc.name, func(t *testing.T) {
				testFileSystem_Open(t, addr, tc.options, tc.handle, tc.compression, tc.info)
			})
		}
	}
}

func testFileSystem_RemoveFile(t *testing.T, addr string) {
	fileName := "rm_test.txt"
	fileParent := "/tmp"
	filePath := fileParent + "/" + fileName

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsDelete)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}

	file.Close(context.Background())

	err = fs.RemoveFile(context.Background(), filePath)
	if err != nil {
		t.Fatalf("invalid rm call: %v", err)
	}

	got, err := fs.Dirlist(context.Background(), fileParent)
	if err != nil {
		t.Fatalf("invalid dirlist call: %v", err)
	}

	found := false
	for _, entry := range got {
		if entry.Name() == fileName {
			found = true
		}
	}

	if found {
		t.Errorf("file '%s' is still present after fs.RemoveFile()", filePath)
	}
}

func TestFileSystem_RemoveFile(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFileSystem_RemoveFile(t, addr)
		})
	}
}

func testFileSystem_Truncate(t *testing.T, addr string) {
	filePath := "/tmp/test_truncate_fs.txt"
	write := []uint8{1, 2, 3, 4, 5, 6, 7, 8}
	want := write[:4]

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()
	file, err := fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsOpenUpdate)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	_, err = file.WriteAt(write, 0)
	if err != nil {
		t.Fatalf("invalid write call: %v", err)
	}

	err = file.Sync(context.Background())
	if err != nil {
		t.Fatalf("invalid sync call: %v", err)
	}

	err = file.Close(context.Background())
	if err != nil {
		t.Fatalf("invalid close call: %v", err)
	}

	err = fs.Truncate(context.Background(), filePath, int64(len(want)))
	if err != nil {
		t.Fatalf("invalid truncate call: %v", err)
	}

	file, err = fs.Open(context.Background(), filePath, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	got := make([]uint8, len(want)+10)
	n, err := file.ReadAt(got, 0)
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

func TestFileSystem_Truncate(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFileSystem_Truncate(t, addr)
		})
	}
}
