// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package xrdio provides a File type that implements various interfaces from the io package.
package xrdio // import "go-hep.org/x/hep/xrootd/xrdio"

import (
	"context"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/pkg/errors"
	xrdclient "go-hep.org/x/hep/xrootd/client"
	"go-hep.org/x/hep/xrootd/xrdfs"
)

// File wraps a xrdfs.File and implements the following interfaces:
//  - io.Closer
//  - io.Reader
//  - io.Writer
//  - io.ReaderAt
//  - io.WriterAt
//  - io.Seeker
type File struct {
	cli *xrdclient.Client
	fs  xrdfs.FileSystem
	f   xrdfs.File

	name string
	pos  int64
	size int64
}

// Open opens the name file, where name is the absolute location of that file
// (xrootd server address and path to the file on that server.)
//
// Example:
//
//  f, err := xrdio.Open("root://server.example.com:1094//some/path/to/file")
func Open(name string) (*File, error) {
	urn, err := url.Parse(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	host := urn.Hostname()
	port := urn.Port()

	path := urn.Path
	if strings.HasPrefix(path, "//") {
		path = path[1:]
	}

	user := ""
	if urn.User != nil {
		user = urn.User.Username()
	}

	addr := host
	if port != "" {
		addr += ":" + port
	}

	xrd, err := xrdclient.NewClient(context.Background(), addr, user)
	if err != nil {
		return nil, errors.Errorf("xrdio: could not connect to xrootd server %q: %v", host, err)
	}

	fs := xrd.FS()
	f, err := fs.Open(context.Background(), path, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		xrd.Close()
		return nil, errors.Errorf("xrdio: could not open %q: %v", name, err)
	}

	xf := &File{cli: xrd, fs: fs, f: f, name: path}
	fi, err := xf.Stat()
	if err != nil {
		xrd.Close()
		return nil, errors.Errorf("xrdio: could not stat %q: %v", name, err)
	}
	xf.size = fi.Size()

	return xf, nil
}

// OpenFrom opens the file name via the given filesystem handle.
// name is the absolute path of the wanted file on the server.
//
// Example:
//
//  f, err := xrdio.OpenFrom(fs, "/some/path/to/file")
func OpenFrom(fs xrdfs.FileSystem, name string) (*File, error) {
	f, err := fs.Open(context.Background(), name, xrdfs.OpenModeOwnerRead, xrdfs.OpenOptionsOpenRead)
	if err != nil {
		return nil, errors.Errorf("xrdio: could not open %q: %v", name, err)
	}

	xf := &File{fs: fs, f: f, name: name}
	fi, err := xf.Stat()
	if err != nil {
		return nil, errors.Errorf("xrdio: could not stat %q: %v", name, err)
	}
	xf.size = fi.Size()

	return xf, nil
}

// Name returns the name of the file.
func (f *File) Name() string {
	return f.name
}

// Close implements io.Closer.
func (f *File) Close() error {
	var (
		err1 = f.f.Close(context.Background())
		err2 error
	)

	if f.cli != nil {
		err2 = f.cli.Close()
	}
	if err1 != nil {
		return errors.WithStack(err1)
	}
	if err2 != nil {
		return errors.WithStack(err2)
	}
	return nil
}

// Read implements io.Reader.
func (f *File) Read(data []byte) (int, error) {
	n, err := f.f.ReadAt(data, f.pos)
	f.pos += int64(n)
	if err != nil {
		return n, err
	}
	if f.pos == f.size {
		err = io.EOF
	}
	return n, err
}

// ReadAt implements io.ReaderAt.
func (f *File) ReadAt(data []byte, offset int64) (int, error) {
	return f.f.ReadAt(data, offset)
}

// Write implements io.Writer.
func (f *File) Write(data []byte) (int, error) {
	n, err := f.f.WriteAt(data, f.pos)
	f.pos += int64(n)
	return n, err
}

// WriteAt implements io.WriterAt.
func (f *File) WriteAt(data []byte, offset int64) (int, error) {
	return f.f.WriteAt(data, offset)
}

// Seek implements io.Seeker
func (f *File) Seek(offset int64, whence int) (int64, error) {
	var err error
	switch whence {
	case io.SeekStart:
		f.pos = offset
	case io.SeekEnd:
		st, err := f.Stat()
		if err != nil {
			return 0, errors.Errorf("xrdio: could not xrootd-stat %q: %v", f.Name(), err)
		}
		f.pos = st.Size() - offset
	case io.SeekCurrent:
		f.pos += offset
	}
	return f.pos, err
}

func (f *File) Stat() (os.FileInfo, error) {
	v, err := f.f.Stat(context.Background())
	return v, err
}

var (
	_ io.Closer   = (*File)(nil)
	_ io.Reader   = (*File)(nil)
	_ io.ReaderAt = (*File)(nil)
	_ io.Writer   = (*File)(nil)
	_ io.WriterAt = (*File)(nil)
	_ io.Seeker   = (*File)(nil)
)
