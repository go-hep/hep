// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import "go-hep.org/x/hep/xrootd/xrdfs"

// File implements access to a content and meta information of file over XRootD.
type file struct {
	fs          *fileSystem
	handle      xrdfs.FileHandle
	compression *xrdfs.FileCompression
	info        *xrdfs.EntryStat
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

var (
	_ xrdfs.File = (*file)(nil)
)
