// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package xrdfuse // import "go-hep.org/x/hep/xrootd/xrdfuse"

import (
	"context"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// File represents a file on the remote XRootD server.
type File struct {
	nodefs.File
	xrdfile xrdfs.File
	fs      *FS
}

// Read implements nodefs.File.Read
func (f *File) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	n, err := f.xrdfile.ReadAt(dest, off)
	if err != nil {
		f.fs.handler(errors.WithMessage(err, "xrdfuse: error calling ReadAt"))
		return nil, fuse.EIO
	}

	return fuse.ReadResultData(dest[:n]), fuse.OK
}

// Write implements nodefs.File.Write
func (f *File) Write(data []byte, off int64) (uint32, fuse.Status) {
	n, err := f.xrdfile.WriteAt(data, off)
	if err != nil {
		f.fs.handler(errors.WithMessage(err, "xrdfuse: error calling WriteAt"))
		return 0, fuse.EIO
	}

	return uint32(n), fuse.OK
}

// Truncate implements nodefs.File.Truncate
func (f *File) Truncate(size uint64) fuse.Status {
	err := f.xrdfile.Truncate(context.Background(), int64(size))
	if err != nil {
		f.fs.handler(errors.WithMessage(err, "xrdfuse: error calling Truncate"))
		return fuse.EIO
	}

	return fuse.OK
}

// Fsync implements nodefs.File.Fsync
func (f *File) Fsync(flags int) (code fuse.Status) {
	return f.Flush()
}

// Flush implements nodefs.File.Flush
func (f *File) Flush() (code fuse.Status) {
	err := f.xrdfile.Sync(context.Background())
	if err != nil {
		f.fs.handler(errors.WithMessage(err, "xrdfuse: error calling Sync"))
		return fuse.EIO
	}

	return fuse.OK
}

// GetAttr implements nodefs.File.GetAttr
func (f *File) GetAttr(out *fuse.Attr) fuse.Status {
	stat, err := f.xrdfile.Stat(context.Background())
	if err != nil {
		f.fs.handler(errors.WithMessage(err, "xrdfuse: error calling Stat"))
		return fuse.EIO
	}

	out.Size = uint64(stat.Size())
	out.Mtime = uint64(stat.Mtime)
	out.Mode = entryStatToMode(stat)

	return fuse.OK
}

// Release implements nodefs.File.Release
func (f *File) Release() {
	err := f.xrdfile.Close(context.Background())
	if err != nil {
		f.fs.handler(errors.WithMessage(err, "xrdfuse: error calling Close"))
	}
}

var (
	_ nodefs.File = (*File)(nil)
)
