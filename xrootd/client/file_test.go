// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/xrdfs"
)

func testFile_Close(t *testing.T, addr string) {
	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	file, err := fs.Open(context.Background(), "/tmp/dir1/file1.txt", xrdfs.OpenModeOtherRead, xrdfs.OpenOptionsNone)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}

	err = file.Close(context.Background())
	if err != nil {
		t.Fatalf("invalid close call: %v", err)
	}
}

func TestFile_Close(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_Close(t, addr)
		})
	}
}

func testFile_CloseVerify(t *testing.T, addr string) {
	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	file, err := fs.Open(context.Background(), "/tmp/test.txt", xrdfs.OpenModeOwnerWrite, xrdfs.OpenOptionsOpenUpdate|xrdfs.OpenOptionsOpenAppend)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}

	// TODO: Remove these 2 lines when XRootD server will follow protocol specification and fail such requests.
	// See https://github.com/xrootd/xrootd/issues/727.
	defer file.Close(context.Background())
	t.Skip("Skipping test because the XRootD C++ server doesn't fail request when the wrong size is passed.")

	err = file.CloseVerify(context.Background(), 14)
	if err == nil {
		t.Fatal("close call should fail when the wrong size is passed")
	}
}

func TestFile_CloseVerify(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_CloseVerify(t, addr)
		})
	}
}

func testFile_ReadAt(t *testing.T, addr string) {
	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}
	defer client.Close()

	fs := client.FS()

	file, err := fs.Open(context.Background(), "/tmp/file1.txt", xrdfs.OpenModeOtherRead, xrdfs.OpenOptionsNone)
	if err != nil {
		t.Fatalf("invalid open call: %v", err)
	}
	defer file.Close(context.Background())

	want := []byte("Hello XRootD.\n")
	got := make([]uint8, 20)
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

func TestFile_ReadAt(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFile_ReadAt(t, addr)
		})
	}
}
