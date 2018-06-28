// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xrdfuse contains the implementation of the FUSE API
// accessing a remote filesystem served over the XRootD protocol.
package xrdfuse // import "go-hep.org/x/hep/xrootd/xrdfuse"

import (
	"context"
	"os"
	"path"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/client"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto"
)

// FS implements a pathfs.FileSystem that makes requests to the remote server over the XRootD protocol.
type FS struct {
	pathfs.FileSystem
	client  *client.Client
	xrdfs   xrdfs.FileSystem
	root    string
	handler ErrorHandler
}

// ErrorHandler is the function which handles occurred error (e.g. logs it).
type ErrorHandler func(error)

// NewFS returns a new path.FileSystem representing the filesystem on the remote XRootD server.
// client is a client connected to the remote XRootD server.
// root is the path to the remote directory to be used as a root directory.
// handler is the function which handles occurred error (e.g. logs it). If the handler is nil,
// then a default handler is used that does nothing.
func NewFS(client *client.Client, root string, handler ErrorHandler) *FS {
	if handler == nil {
		handler = func(error) {}
	}
	return &FS{
		FileSystem: pathfs.NewDefaultFileSystem(),
		client:     client,
		xrdfs:      client.FS(),
		root:       root,
		handler:    handler,
	}
}

// GetAttr implements pathfs.FileSystem.GetAttr
func (fs *FS) GetAttr(name string, ctx *fuse.Context) (*fuse.Attr, fuse.Status) {
	stat, err := fs.xrdfs.Stat(context.Background(), path.Join(fs.root, name))
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Stat"))
	}
	if status != fuse.OK {
		return nil, status
	}

	return &fuse.Attr{
		Size:  uint64(stat.Size()),
		Mtime: uint64(stat.Mtime),
		Mode:  entryStatToMode(stat),
	}, fuse.OK
}

// OpenDir implements pathfs.FileSystem.OpenDir
func (fs *FS) OpenDir(name string, ctx *fuse.Context) (c []fuse.DirEntry, code fuse.Status) {
	entries, err := fs.xrdfs.Dirlist(context.Background(), path.Join(fs.root, name))
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Dirlist"))
	}
	if status != fuse.OK {
		return nil, status
	}

	s := make([]fuse.DirEntry, 0, len(entries))
	for _, entry := range entries {
		s = append(s, fuse.DirEntry{Name: entry.Name(), Mode: entryStatToMode(entry)})
	}

	return s, fuse.OK
}

func convertFlagsToMode(flags uint32) xrdfs.OpenMode {
	switch {
	case flags&uint32(os.O_RDWR) != 0:
		return xrdfs.OpenModeOwnerRead | xrdfs.OpenModeOwnerWrite
	case flags&uint32(os.O_RDONLY) != 0:
		return xrdfs.OpenModeOwnerRead
	case flags&uint32(os.O_WRONLY) != 0:
		return xrdfs.OpenModeOwnerWrite
	}
	return 0
}

func convertFlagsToOptions(flags uint32) xrdfs.OpenOptions {
	var result xrdfs.OpenOptions
	if flags&uint32(os.O_APPEND) != 0 {
		result |= xrdfs.OpenOptionsOpenAppend
	}
	if flags&uint32(os.O_CREATE) != 0 {
		result |= xrdfs.OpenOptionsNew
	}
	if flags&uint32(os.O_TRUNC) != 0 {
		result |= xrdfs.OpenOptionsDelete
	}
	if flags&fuse.O_ANYWRITE != 0 {
		result |= xrdfs.OpenOptionsOpenUpdate
	} else {
		result |= xrdfs.OpenOptionsOpenRead
	}
	return result
}

func convertModeToXrdMode(mode uint32) xrdfs.OpenMode {
	// XRootD open mode follows Unix permissions.
	return xrdfs.OpenMode(mode)
}

// Open implements pathfs.FileSystem.Open
func (fs *FS) Open(name string, flags uint32, ctx *fuse.Context) (file nodefs.File, code fuse.Status) {
	mode := convertFlagsToMode(flags)
	options := convertFlagsToOptions(flags)
	f, err := fs.xrdfs.Open(context.Background(), path.Join(fs.root, name), mode, options)
	if serverError, ok := err.(xrdproto.ServerError); ok {
		if serverError.Code == xrdproto.InvalidRequestCode {
			// It is possible the request is invalid because the file already exists.
			// O_CREAT flag can be passed to the fuse API despite the fact that file
			// is already created. Open should correctly handle this situation, as far
			// as I can see from the docs.
			f, err = fs.xrdfs.Open(context.Background(), path.Join(fs.root, name), mode, options^xrdfs.OpenOptionsNew)
		}
	}
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Open"))
	}
	if status != fuse.OK {
		return nil, status
	}
	return &File{File: nodefs.NewDefaultFile(), xrdfile: f, fs: fs}, fuse.OK
}

// Mknod implements pathfs.FileSystem.Mknod
func (fs *FS) Mknod(name string, mode uint32, dev uint32, ctx *fuse.Context) fuse.Status {
	xrdmode := convertModeToXrdMode(mode)
	f, err := fs.xrdfs.Open(context.Background(), path.Join(fs.root, name), xrdmode, xrdfs.OpenOptionsNew)
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Open"))
	}
	if status != fuse.OK {
		return status
	}

	err = f.Close(context.Background())
	status = errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Close"))
	}
	return status
}

// Rename implements pathfs.FileSystem.Rename
func (fs *FS) Rename(oldName string, newName string, ctx *fuse.Context) fuse.Status {
	err := fs.xrdfs.Rename(context.Background(), path.Join(fs.root, oldName), path.Join(fs.root, newName))
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Rename"))
	}
	return status
}

// Unlink implements pathfs.FileSystem.Unlink
func (fs *FS) Unlink(name string, ctx *fuse.Context) fuse.Status {
	err := fs.xrdfs.RemoveFile(context.Background(), path.Join(fs.root, name))
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling RemoveFile"))
	}
	return status
}

// Rmdir implements pathfs.FileSystem.Rmdir
func (fs *FS) Rmdir(name string, ctx *fuse.Context) fuse.Status {
	err := fs.xrdfs.RemoveDir(context.Background(), path.Join(fs.root, name))
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling RemoveDir"))
	}
	return status
}

// Mkdir implements pathfs.FileSystem.Mkdir
func (fs *FS) Mkdir(name string, mode uint32, ctx *fuse.Context) fuse.Status {
	xrdmode := convertModeToXrdMode(mode)
	err := fs.xrdfs.Mkdir(context.Background(), path.Join(fs.root, name), xrdmode)
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Mkdir"))
	}
	return status
}

// Chmod implements pathfs.FileSystem.Chmod
func (fs *FS) Chmod(name string, mode uint32, ctx *fuse.Context) fuse.Status {
	xrdmode := convertModeToXrdMode(mode)
	err := fs.xrdfs.Chmod(context.Background(), path.Join(fs.root, name), xrdmode)
	status := errorToStatus(err)
	if status == fuse.EIO {
		fs.handler(errors.WithMessage(err, "xrdfuse: error calling Chmod"))
	}
	return status
}

func entryStatToMode(stat xrdfs.EntryStat) uint32 {
	mode := uint32(0)
	if stat.IsDir() {
		mode |= fuse.S_IFDIR
	} else {
		mode |= fuse.S_IFREG
	}
	if stat.IsReadable() {
		mode |= 0444
	}
	if stat.IsWritable() {
		mode |= 0222
	}
	if stat.IsExecutable() {
		mode |= 0111
	}
	return mode
}

func errorToStatus(err error) fuse.Status {
	if err == nil {
		return fuse.OK
	}
	if serverError, ok := err.(xrdproto.ServerError); ok {
		switch serverError.Code {
		case xrdproto.NotFoundCode:
			// File does not exists.
			return fuse.ENOENT
		case xrdproto.NotAuthorized:
			// Permission denied.
			return fuse.EACCES
		default:
			return fuse.EIO
		}
	}
	return fuse.EIO
}

var (
	_ pathfs.FileSystem = (*FS)(nil)
)
