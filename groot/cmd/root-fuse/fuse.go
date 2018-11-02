// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"

	"go-hep.org/x/hep/groot/riofs"
	"go-hep.org/x/hep/groot/root"
)

type FS struct {
	pathfs.FileSystem
	root  *riofs.File
	mtime uint64
	mode  uint32
}

func NewFS(f *riofs.File) *FS {
	st, err := f.Stat()
	if err != nil {
		return nil
	}
	return &FS{
		FileSystem: pathfs.NewDefaultFileSystem(),
		root:       f,
		mtime:      uint64(st.ModTime().Unix()),
		mode:       fuse.S_IFDIR | 0444 | 0111,
	}
}

func (fs *FS) OpenDir(name string, ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	v, err := fs.get(name)
	if err != nil {
		return nil, fuse.ENOENT
	}
	dir, ok := v.(riofs.Directory)
	if !ok {
		return nil, fuse.ENOTDIR
	}
	o := make([]fuse.DirEntry, len(dir.Keys()))
	for i, key := range dir.Keys() {
		o[i] = fuse.DirEntry{
			Name: key.Name(),
			Mode: 0, // FIXME(sbinet)
		}
	}
	return o, fuse.OK
}

func (fs *FS) Open(name string, flags uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	obj, err := fs.get(name)
	if err != nil {
		return nil, fuse.ENOENT
	}
	named := obj.(root.Named)
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "name:  %s\ntitle: %s\ntype:  %s\n",
		named.Name(),
		named.Title(),
		obj.Class(),
	)
	return &File{
		File: nodefs.NewDefaultFile(), fs: fs, obj: obj, data: buf.Bytes(),
	}, fuse.OK
}

// GetAttr implements paths.FileSystem.GetAttr
func (fs *FS) GetAttr(name string, ctx *fuse.Context) (*fuse.Attr, fuse.Status) {
	obj, err := fs.get(name)
	if err != nil {
		return nil, fuse.ENOENT
	}

	named := obj.(root.Named)
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "name:  %s\ntitle: %s\ntype:  %s\n",
		named.Name(),
		named.Title(),
		obj.Class(),
	)
	f := &File{
		File: nodefs.NewDefaultFile(), fs: fs, obj: obj, data: buf.Bytes(),
	}
	defer f.Release()

	var attr fuse.Attr
	status := f.GetAttr(&attr)
	return &attr, status
}

func (fs *FS) get(name string) (root.Object, error) {
	if name == "" {
		return fs.root, nil
	}
	dirs := strings.Split(name, "/")
	var (
		ctx riofs.Directory = fs.root
		obj root.Object     = fs.root
		err error
	)
	for i, n := range dirs {
		obj, err = ctx.Get(n)
		if err != nil {
			return nil, err
		}
		if i != len(dirs)-1 {
			ctx = obj.(riofs.Directory)
		}
	}
	return obj, nil
}

func (fs *FS) StatFs(name string) *fuse.StatfsOut {
	var out fuse.StatfsOut
	if name != "" {
		return &out
	}
	st, err := fs.root.Stat()
	if err != nil {
		return &out
	}

	out.Blocks = statBlocks(uint64(st.Size()))
	out.Bsize = blockSize
	out.NameLen = uint32(len(fs.root.Name()))
	out.Files = uint64(len(fs.root.Keys()))
	return &out
}

const blockSize = 512 // FIXME(sbinet): arbritrary size, at this point.

func statBlocks(size uint64) uint64 {
	r := size / blockSize
	if size%blockSize > 0 {
		r++
	}
	return r
}

type File struct {
	nodefs.File
	fs   *FS
	obj  root.Object
	data []byte
}

// Release implements nodefs.File.Release
func (f *File) Release() {
	f.data = nil
	f.obj = nil
	f.fs = nil
}

// Flush implements nodefs.File.Flush
func (f *File) Flush() fuse.Status {
	f.data = nil
	f.obj = nil
	f.fs = nil
	return fuse.OK
}

// GetAttr implements nodefs.File.GetAttr
func (f *File) GetAttr(out *fuse.Attr) fuse.Status {
	out.Size = uint64(len(f.data))
	out.Mtime = f.fs.mtime
	out.Mode = 0
	switch f.obj.(type) {
	case riofs.Directory:
		out.Mode |= fuse.S_IFDIR
	default:
		out.Mode |= fuse.S_IFREG
	}
	out.Mode |= 0444
	if f.obj == f.fs.root {
		out.Mode |= 0111
	}

	return fuse.OK
}

// Read implements nodefs.File.Read
func (f *File) Read(p []byte, off int64) (fuse.ReadResult, fuse.Status) {
	copy(p, f.data[off:])
	return fuse.ReadResultData(p), fuse.OK
}

var (
	_ pathfs.FileSystem = (*FS)(nil)
	_ nodefs.File       = (*File)(nil)
)
