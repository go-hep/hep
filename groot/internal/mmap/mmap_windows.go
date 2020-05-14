// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package mmap provides a way to memory-map a file.
package mmap

import (
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

// Reader reads a memory-mapped file.
type Reader struct {
	data []byte
	c    int
}

// Close closes the reader.
func (r *Reader) Close() error {
	if r.data == nil {
		return nil
	}
	data := r.data
	r.data = nil
	runtime.SetFinalizer(r, nil)
	return syscall.UnmapViewOfFile(uintptr(unsafe.Pointer(&data[0])))
}

// Len returns the length of the underlying memory-mapped file.
func (r *Reader) Len() int {
	return len(r.data)
}

// At returns the byte at index i.
func (r *Reader) At(i int) byte {
	return r.data[i]
}

func (r *Reader) Read(p []byte) (int, error) {
	if r.c >= len(r.data) {
		return 0, io.EOF
	}
	n := copy(p, r.data[r.c:])
	r.c += n
	return n, nil
}

func (r *Reader) ReadByte() (byte, error) {
	if r.c >= len(r.data) {
		return 0, io.EOF
	}
	v := r.data[r.c]
	r.c++
	return v, nil
}

// ReadAt implements the io.ReaderAt interface.
func (r *Reader) ReadAt(p []byte, off int64) (int, error) {
	if r.data == nil {
		return 0, errors.New("mmap: closed")
	}
	if off < 0 || int64(len(r.data)) < off {
		return 0, fmt.Errorf("mmap: invalid ReadAt offset %d", off)
	}
	n := copy(p, r.data[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.c = int(offset)
	case io.SeekCurrent:
		r.c += int(offset)
	case io.SeekEnd:
		r.c = len(r.data) - int(offset)
	default:
		return 0, fmt.Errorf("mmap: invalid whence")
	}
	if r.c < 0 {
		return 0, fmt.Errorf("mmap: negative position")
	}
	return int64(r.c), nil
}

// Open memory-maps the named file for reading.
func Open(filename string) (*Reader, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	size := fi.Size()
	if size == 0 {
		return &Reader{}, nil
	}
	if size < 0 {
		return nil, fmt.Errorf("mmap: file %q has negative size", filename)
	}
	if size != int64(int(size)) {
		return nil, fmt.Errorf("mmap: file %q is too large", filename)
	}

	low, high := uint32(size), uint32(size>>32)
	fmap, err := syscall.CreateFileMapping(syscall.Handle(f.Fd()), nil, syscall.PAGE_READONLY, high, low, nil)
	if err != nil {
		return nil, err
	}
	defer syscall.CloseHandle(fmap)
	ptr, err := syscall.MapViewOfFile(fmap, syscall.FILE_MAP_READ, 0, 0, uintptr(size))
	if err != nil {
		return nil, err
	}
	data := (*[maxBytes]byte)(unsafe.Pointer(ptr))[:size]

	r := &Reader{data: data}
	runtime.SetFinalizer(r, (*Reader).Close)
	return r, nil
}

var (
	_ io.Reader     = (*Reader)(nil)
	_ io.ReaderAt   = (*Reader)(nil)
	_ io.Seeker     = (*Reader)(nil)
	_ io.Closer     = (*Reader)(nil)
	_ io.ByteReader = (*Reader)(nil)
)
