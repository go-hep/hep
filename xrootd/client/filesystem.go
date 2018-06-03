// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto/dirlist"
	"go-hep.org/x/hep/xrootd/xrdproto/open"
	"go-hep.org/x/hep/xrootd/xrdproto/rm"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
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

// Open returns the file handle for a file together with the compression and the stat info.
func (fs *fileSystem) Open(ctx context.Context, path string, mode xrdfs.OpenMode, options xrdfs.OpenOptions) (xrdfs.File, error) {
	var resp open.Response
	err := fs.c.Send(ctx, &resp, open.NewRequest(path, mode, options))
	if err != nil {
		return nil, err
	}
	return &file{fs, resp.FileHandle, resp.Compression, resp.Stat}, nil
}

// RemoveFile removes a file.
func (fs *fileSystem) RemoveFile(ctx context.Context, path string) error {
	_, err := fs.c.call(ctx, &rm.Request{Path: path})
	return err
}

// Truncate changes the size of the named file.
func (fs *fileSystem) Truncate(ctx context.Context, path string, size int64) error {
	_, err := fs.c.call(ctx, &truncate.Request{Path: path, Size: size})
	return err
}

// Stat returns the entry stat info for the given path.
func (fs *fileSystem) Stat(ctx context.Context, path string) (xrdfs.EntryStat, error) {
	var resp stat.DefaultResponse
	err := fs.c.Send(ctx, &resp, &stat.Request{Path: path})
	if err != nil {
		return xrdfs.EntryStat{}, err
	}
	return resp.EntryStat, nil
}

// VirtualStat returns the virtual filesystem stat info for the given path.
// Note that path needs not to be an existing filesystem object, it is used as a path prefix in order to
// filter out servers and partitions that could not be used to hold objects whose path starts
// with the specified path prefix.
func (fs *fileSystem) VirtualStat(ctx context.Context, path string) (xrdfs.VirtualFSStat, error) {
	var resp stat.VirtualFSResponse
	err := fs.c.Send(ctx, &resp, &stat.Request{Path: path, Options: stat.OptionsVFS})
	if err != nil {
		return xrdfs.VirtualFSStat{}, err
	}
	return resp.VirtualFSStat, nil
}

var (
	_ xrdfs.FileSystem = (*fileSystem)(nil)
)
