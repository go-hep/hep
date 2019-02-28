// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"io"

	"github.com/pkg/errors"
)

// A Compressor takes data written to it and writes the compressed form of that
// data to an underlying writer.
type Compressor interface {
	io.WriteCloser

	// Reset clears the state of the Writer such that it is equivalent to its
	// initial state, but instead writing to w.
	Reset(w io.Writer) error

	// Flush flushes the Writer to its underlying io.Writer.
	Flush() error

	// Kind returns the (de)compressor kind
	//Kind() CompressorKind
}

// A Decompressor reads data from the underlying io.Reader and decompresses it.
type Decompressor interface {
	io.ReadCloser

	// Reset clears the state of the Reader such that it is equivalent to its
	// initial state, but instead reading from r.
	Reset(r io.Reader) error

	// Kind returns the (de)compressor kind
	//Kind() CompressorKind
}

// Xpressor provides compressor and decompressor functions
type Xpressor interface {
	Inflate(r io.Reader) (Decompressor, error)
	Deflate(w io.Writer, opts Options) (Compressor, error)
}

// CompressorKind names the various compressors
type CompressorKind uint16

// builtin compressor types
const (
	CompressDefault CompressorKind = iota
	CompressNone
	CompressFlate
	CompressZlib
	CompressGzip
	CompressLZA
	CompressLZO
	CompressSnappy

	CompressUser CompressorKind = 0xffff // keep last
)

func (ck CompressorKind) String() string {
	switch ck {
	case CompressUser:
		return "user"
	case CompressDefault:
		return "default"
	case CompressNone:
		return "none"
	case CompressFlate:
		return "flate"
	case CompressZlib:
		return "zlib"
	case CompressGzip:
		return "gzip"
	case CompressLZA:
		return "lza"
	case CompressLZO:
		return "lzo"
	case CompressSnappy:
		return "snappy"
	}
	return "N/A"
}

// NewCompressor creates a Compressor writing to w, with compression level according to opts.
func (ck CompressorKind) NewCompressor(w io.Writer, opts Options) (Compressor, error) {
	x, ok := xcomprs[ck]
	if !ok {
		return nil, errors.Errorf("rio: no compressor registered with %q (%v)", ck.String(), int(ck))
	}
	return x.Deflate(w, opts)
}

// NewDecompressor creates a Decompressor reading from r.
func (ck CompressorKind) NewDecompressor(r io.Reader) (Decompressor, error) {
	x, ok := xcomprs[ck]
	if !ok {
		return nil, errors.Errorf("rio: no decompressor registered with %q (%v)", ck.String(), int(ck))
	}
	return x.Inflate(r)
}

type xpressor struct {
	inflater func(r io.Reader) (Decompressor, error)
	deflater func(w io.Writer, opts Options) (Compressor, error)
}

func (x xpressor) Inflate(r io.Reader) (Decompressor, error) {
	return x.inflater(r)
}

func (x xpressor) Deflate(w io.Writer, o Options) (Compressor, error) {
	return x.deflater(w, o)
}

type noneCompressor struct{}

func (noneCompressor) Inflate(r io.Reader) (Decompressor, error) {
	return &nopReadCloser{r}, nil
}

func (noneCompressor) Deflate(w io.Writer, o Options) (Compressor, error) {
	return &nopWriteCloser{w}, nil
}

type nopWriteCloser struct {
	io.Writer
}

func (nw *nopWriteCloser) Close() error {
	return nil
}

func (nw *nopWriteCloser) Flush() error {
	return nil
}

func (nw *nopWriteCloser) Reset(w io.Writer) error {
	nw.Writer = w
	return nil
}

type nopReadCloser struct {
	io.Reader
}

func (nr *nopReadCloser) Close() error {
	return nil
}

func (nr *nopReadCloser) Reset(r io.Reader) error {
	nr.Reader = r
	return nil
}

// flate ---

type flateCompressor struct {
	*flate.Writer
}

func newFlateCompressor(w io.Writer, o Options) (Compressor, error) {
	cw, err := flate.NewWriter(w, o.CompressorLevel())
	return &flateCompressor{cw}, err
}

func (cw *flateCompressor) Reset(w io.Writer) error {
	cw.Writer.Reset(w)
	return nil
}

type flateDecompressor struct {
	io.ReadCloser
}

func newFlateDecompressor(r io.Reader) (Decompressor, error) {
	xr := flate.NewReader(r)
	return &flateDecompressor{xr}, nil
}

func (xr *flateDecompressor) Reset(r io.Reader) error {
	return xr.ReadCloser.(flate.Resetter).Reset(r, nil)
}

// zlib ---

type zlibCompressor struct {
	*zlib.Writer
}

func newZlibCompressor(w io.Writer, o Options) (Compressor, error) {
	cw, err := zlib.NewWriterLevel(w, o.CompressorLevel())
	return &zlibCompressor{cw}, err
}

func (cw *zlibCompressor) Reset(w io.Writer) error {
	cw.Writer.Reset(w)
	return nil
}

type zlibDecompressor struct {
	io.ReadCloser
}

func newZlibDecompressor(r io.Reader) (Decompressor, error) {
	xr, err := zlib.NewReader(r)
	return &zlibDecompressor{xr}, err
}

func (xr *zlibDecompressor) Reset(r io.Reader) error {
	return xr.ReadCloser.(zlib.Resetter).Reset(r, nil)
}

// gzip ---

type gzipCompressor struct {
	*gzip.Writer
}

func newGzipCompressor(w io.Writer, o Options) (Compressor, error) {
	cw, err := gzip.NewWriterLevel(w, o.CompressorLevel())
	return &gzipCompressor{cw}, err
}

func (cw *gzipCompressor) Reset(w io.Writer) error {
	cw.Writer.Reset(w)
	return nil
}

type gzipDecompressor struct {
	*gzip.Reader
}

func newGzipDecompressor(r io.Reader) (Decompressor, error) {
	xr, err := gzip.NewReader(r)
	return &gzipDecompressor{xr}, err
}

// RegisterCompressor registers a compressor/decompressor.
// It silently replaces the compressor/decompressor if one was already registered.
func RegisterCompressor(kind CompressorKind, x Xpressor) {
	xcomprs[kind] = x
}

var xcomprs map[CompressorKind]Xpressor

func init() {

	xcomprs = map[CompressorKind]Xpressor{

		CompressNone: noneCompressor{},

		CompressFlate: xpressor{
			inflater: newFlateDecompressor,
			deflater: newFlateCompressor,
		},

		CompressZlib: xpressor{
			inflater: newZlibDecompressor,
			deflater: newZlibCompressor,
		},

		CompressGzip: xpressor{
			inflater: newGzipDecompressor,
			deflater: newGzipCompressor,
		},
	}

	xcomprs[CompressDefault] = xcomprs[CompressZlib]
}
