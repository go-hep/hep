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
