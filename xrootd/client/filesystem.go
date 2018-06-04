// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
)

// FS returns a xrdfs.FileSystem which uses this client to make requests.
func (cli *Client) FS() xrdfs.FileSystem {
	return &fileSystem{cli}
}

// fileSystem contains filesystem-related methods of the XRootD protocol.
type fileSystem struct {
	c *Client
}

// Dirlist returns the contents of a directory together with the stat information.
func (fs *fileSystem) Dirlist(ctx context.Context, path string) ([]xrdfs.EntryStat, error) {
	var resp dirlist.Response
	err := fs.c.Send(ctx, &resp, dirlist.NewRequest(path))
	if err != nil {
		return nil, err
	}
	return resp.Entries, err
}

var (
	_ xrdfs.FileSystem = (*fileSystem)(nil)
)
