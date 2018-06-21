// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"context"

	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdproto/read"
	"go-hep.org/x/hep/xrootd/xrdproto/stat"
	"go-hep.org/x/hep/xrootd/xrdproto/sync"
	"go-hep.org/x/hep/xrootd/xrdproto/truncate"
	"go-hep.org/x/hep/xrootd/xrdproto/verifyw"
	"go-hep.org/x/hep/xrootd/xrdproto/write"
	"go-hep.org/x/hep/xrootd/xrdproto/xrdclose"
)

// File implements access to a content and meta information of file over XRootD.
type file struct {
	fs          *fileSystem
	handle      xrdfs.FileHandle
	compression *xrdfs.FileCompression
	info        *xrdfs.EntryStat
	sessionID   string
}

// Compression returns the compression info.
func (f file) Compression() *xrdfs.FileCompression {
	return f.compression
}

// Info returns the cached stat info.
// Note that it may return nil if info was not yet fetched and info may be not up-to-date.
func (f file) Info() *xrdfs.EntryStat {
	return f.info
}

// Handle returns the file handle.
func (f file) Handle() xrdfs.FileHandle {
	return f.handle
}

// Close closes the file.
func (f file) Close(ctx context.Context) error {
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, nil, &xrdclose.Request{Handle: f.handle})
	if err != nil {
		return err
	}
	f.sessionID = newSessionID
	return nil
}

// CloseVerify closes the file and checks whether the file has the provided size.
// A zero size suppresses the verification.
func (f file) CloseVerify(ctx context.Context, size int64) error {
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, nil, &xrdclose.Request{Handle: f.handle, Size: size})
	if err != nil {
		return err
	}
	f.sessionID = newSessionID
	return nil
}

// Sync commits all pending writes to an open file.
func (f file) Sync(ctx context.Context) error {
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, nil, &sync.Request{Handle: f.handle})
	if err != nil {
		return err
	}
	f.sessionID = newSessionID
	return nil
}

// ReadAtContext reads len(p) bytes into p starting at offset off.
func (f file) ReadAtContext(ctx context.Context, p []byte, off int64) (n int, err error) {
	resp := read.Response{Data: p}
	req := &read.Request{Handle: f.handle, Offset: off, Length: int32(len(p))}
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, &resp, req)
	if err != nil {
		return 0, err
	}
	f.sessionID = newSessionID
	return len(resp.Data), nil
}

// ReadAt reads len(p) bytes into p starting at offset off.
func (f file) ReadAt(p []byte, off int64) (n int, err error) {
	return f.ReadAtContext(context.Background(), p, off)
}

// WriteAtContext writes len(p) bytes from p to the file at offset off.
func (f file) WriteAtContext(ctx context.Context, p []byte, off int64) error {
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, nil, &write.Request{Handle: f.handle, Offset: off, Data: p})
	if err != nil {
		return err
	}
	f.sessionID = newSessionID
	return nil
}

// WriteAt writes len(p) bytes from p to the file at offset off.
func (f file) WriteAt(p []byte, off int64) (n int, err error) {
	err = f.WriteAtContext(context.Background(), p, off)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

// Truncate changes the size of the named file.
func (f file) Truncate(ctx context.Context, size int64) error {
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, nil, &truncate.Request{Handle: f.handle, Size: size})
	if err != nil {
		return err
	}
	f.sessionID = newSessionID
	return nil
}

// StatVirtualFS fetches the virtual fs stat info from the XRootD server.
// TODO: note that calling stat with vfs and handle may be invalid.
// See https://github.com/xrootd/xrootd/issues/728 for the details.
func (f file) StatVirtualFS(ctx context.Context) (xrdfs.VirtualFSStat, error) {
	var resp stat.VirtualFSResponse
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, &resp, &stat.Request{FileHandle: f.handle, Options: stat.OptionsVFS})
	if err != nil {
		return xrdfs.VirtualFSStat{}, err
	}
	f.sessionID = newSessionID
	return resp.VirtualFSStat, nil
}

// Stat fetches the stat info of this file from the XRootD server.
// Note that Stat re-fetches value returned by the Info, so after the call to Stat
// calls to Info may return different value than before.
func (f *file) Stat(ctx context.Context) (xrdfs.EntryStat, error) {
	var resp stat.DefaultResponse
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, &resp, &stat.Request{FileHandle: f.handle})
	if err != nil {
		return xrdfs.EntryStat{}, err
	}
	f.info = &resp.EntryStat
	f.sessionID = newSessionID
	return resp.EntryStat, nil
}

// VerifyWriteAt writes len(p) bytes from p to the file at offset off using crc32 verification.
//
// TODO: note that verifyw is not supported by the XRootD server.
// See https://github.com/xrootd/xrootd/issues/738 for the details.
func (f file) VerifyWriteAt(ctx context.Context, p []byte, off int64) error {
	newSessionID, err := f.fs.c.sendSession(ctx, f.sessionID, nil, verifyw.NewRequestCRC32(f.handle, off, p))
	if err != nil {
		return err
	}
	f.sessionID = newSessionID
	return nil
}

var (
	_ xrdfs.File = (*file)(nil)
)
