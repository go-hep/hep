// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package riofs

import "bytes"

// memFile is a simple in-memory read-only ROOT file
type memFile struct {
	r *bytes.Reader
}

func (r *memFile) Close() error                                 { return nil }
func (r *memFile) Read(p []byte) (int, error)                   { return r.r.Read(p) }
func (r *memFile) ReadAt(p []byte, off int64) (int, error)      { return r.r.ReadAt(p, off) }
func (r *memFile) Seek(offset int64, whence int) (int64, error) { return r.r.Seek(offset, whence) }

var (
	_ Reader = (*memFile)(nil)
)
