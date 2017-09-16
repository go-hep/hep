// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bio

import (
	"bufio"
	"io"
)

type Reader struct {
	raw io.ReadSeeker
	r   *bufio.Reader
	pos int64
}

func NewReader(r io.ReadSeeker) *Reader {
	pos, _ := r.Seek(0, 0)
	rr := Reader{
		raw: r,
		r:   bufio.NewReaderSize(r, 32*1024*1024),
		pos: pos,
	}
	return &rr
}

func (r *Reader) Close() error { return nil }

func (r *Reader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	r.pos += int64(n)
	return n, err
}

func (r *Reader) ReadAt(p []byte, off int64) (int, error) {
	if off >= r.pos {
		d, err := r.r.Discard(int(off - r.pos))
		if err != nil {
			return d, err
		}
		//		log.Printf("reuse buffered reader (off=%d, pos=%d)", off, r.pos)
		return r.Read(p)
	}
	//	log.Printf("reset buffered reader (off=%d, pos=%d)", off, r.pos)
	var err error
	r.pos, err = r.raw.Seek(off, io.SeekStart)
	if err != nil {
		return 0, err
	}
	r.r.Reset(r.raw)
	return r.Read(p)
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	var err error
	r.pos, err = r.raw.Seek(offset, whence)
	// FIXME(sbinet): only reset when needed
	r.r.Reset(r.raw)
	return r.pos, err
}
