// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build xrootd_test_with_server

package client

import (
	"context"
	"reflect"
	"testing"

	"go-hep.org/x/hep/xrootd/protocol"
)

func testFilesystem_Dirlist(t *testing.T, addr string) {
	var want = []protocol.EntryStat{
		{
			Name:         "admin",
			HasStatInfo:  true,
			Id:           81604443394,
			IsExecutable: true,
			IsOther:      true,
			IsWritable:   true,
			IsReadable:   true,
			Mtime:        1519143383,
		},
	}

	client, err := NewClient(context.Background(), addr, "gopher")
	if err != nil {
		t.Fatalf("could not create client: %v", err)
	}

	fs := Filesystem{client}

	got, err := fs.Dirlist(context.Background(), "/tmp/.xrootd")
	if err != nil {
		t.Fatalf("invalid protocol call: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Filesystem.Dirlist()\ngot = %v\nwant = %v", got, want)
	}

	client.Close()
}

func TestFilesystem_Dirlist(t *testing.T) {
	for _, addr := range testClientAddrs {
		t.Run(addr, func(t *testing.T) {
			testFilesystem_Dirlist(t, addr)
		})
	}
}
