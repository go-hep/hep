// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"
	stdpath "path"

	"go-hep.org/x/hep/xrootd/xrdfs"
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
	_, err := fs.c.Send(ctx, &resp, dirlist.NewRequest(path))
	if err != nil {
		return nil, err
	}
	return resp.Entries, err
}

// Open returns the file handle for a file together with the compression and the stat info.
func (fs *fileSystem) Open(ctx context.Context, path string, mode xrdfs.OpenMode, options xrdfs.OpenOptions) (xrdfs.File, error) {
	var resp open.Response
	server, err := fs.c.Send(ctx, &resp, open.NewRequest(path, mode, options))
	if err != nil {
		return nil, err
	}
	return &file{fs, resp.FileHandle, resp.Compression, resp.Stat, server}, nil
}

// RemoveFile removes a file.
func (fs *fileSystem) RemoveFile(ctx context.Context, path string) error {
	_, err := fs.c.Send(ctx, nil, &rm.Request{Path: path})
	return err
}

// Truncate changes the size of the named file.
func (fs *fileSystem) Truncate(ctx context.Context, path string, size int64) error {
	_, err := fs.c.Send(ctx, nil, &truncate.Request{Path: path, Size: size})
	return err
}

// Stat returns the entry stat info for the given path.
func (fs *fileSystem) Stat(ctx context.Context, path string) (xrdfs.EntryStat, error) {
	var resp stat.DefaultResponse
	_, err := fs.c.Send(ctx, &resp, &stat.Request{Path: path})
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
	_, err := fs.c.Send(ctx, &resp, &stat.Request{Path: path, Options: stat.OptionsVFS})
	if err != nil {
		return xrdfs.VirtualFSStat{}, err
	}
	return resp.VirtualFSStat, nil
}

// Mkdir creates a new directory with the specified name and permission bits.
func (fs *fileSystem) Mkdir(ctx context.Context, path string, perm xrdfs.OpenMode) error {
	_, err := fs.c.Send(ctx, nil, &mkdir.Request{Path: path, Mode: perm})
	return err
}

// MkdirAll creates a directory named path, along with any necessary parents,
// and returns nil, or else returns an error.
// The permission bits perm are used for all directories that MkdirAll creates.
func (fs *fileSystem) MkdirAll(ctx context.Context, path string, perm xrdfs.OpenMode) error {
	_, err := fs.c.Send(ctx, nil, &mkdir.Request{Path: path, Mode: perm, Options: mkdir.OptionsMakePath})
	return err
}

// RemoveDir removes a directory.
// The directory to be removed must be empty.
func (fs *fileSystem) RemoveDir(ctx context.Context, path string) error {
	_, err := fs.c.Send(ctx, nil, &rmdir.Request{Path: path})
	return err
}

// RemoveAll removes path and any children it contains.
// It removes everything it can but returns the first error it encounters.
// If the path does not exist, RemoveAll returns nil (no error.)
func (fs *fileSystem) RemoveAll(ctx context.Context, path string) error {
	st, err := fs.Stat(ctx, path)
	if err != nil {
		return err
	}
	switch {
	case st.IsDir():
		entries, err := fs.Dirlist(ctx, path)
		if err != nil {
			return err
		}
		for _, e := range entries {
			name := stdpath.Join(path, e.Name())
			err := fs.RemoveAll(ctx, name)
			if err != nil {
				return err
			}
		}
		return fs.RemoveDir(ctx, path)
	default:
		return fs.RemoveFile(ctx, path)
	}
}

// Rename renames (moves) oldpath to newpath.
func (fs *fileSystem) Rename(ctx context.Context, oldpath, newpath string) error {
	_, err := fs.c.Send(ctx, nil, &mv.Request{OldPath: oldpath, NewPath: newpath})
	return err
}

// Chmod changes the permissions of the named file to perm.
func (fs *fileSystem) Chmod(ctx context.Context, path string, perm xrdfs.OpenMode) error {
	_, err := fs.c.Send(ctx, nil, &chmod.Request{Path: path, Mode: perm})
	return err
}

// Statx obtains type information for one or more paths.
// Only a limited number of flags is meaningful such as StatIsExecutable, StatIsDir, StatIsOther, StatIsOffline.
func (fs *fileSystem) Statx(ctx context.Context, paths []string) ([]xrdfs.StatFlags, error) {
	var resp statx.Response
	_, err := fs.c.Send(ctx, &resp, statx.NewRequest(paths))
	if err != nil {
		return nil, err
	}
	return resp.StatFlags, nil
}

var (
	_ xrdfs.FileSystem = (*fileSystem)(nil)
)
