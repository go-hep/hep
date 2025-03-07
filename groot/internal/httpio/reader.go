// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package httpio

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
)

// Reader presents an HTTP resource as an io.Reader and io.ReaderAt.
type Reader struct {
	cli    *http.Client
	req    *http.Request
	ctx    context.Context
	cancel context.CancelFunc

	pool sync.Pool

	r    *io.SectionReader
	len  int64
	etag string
}

// Open returns a Reader from the provided URL.
func Open(uri string, opts ...Option) (r *Reader, err error) {
	cfg := newConfig()
	for _, opt := range opts {
		err := opt(cfg)
		if err != nil {
			return nil, fmt.Errorf("httpio: could not open %q: %w", uri, err)
		}
	}

	r = &Reader{
		cli: cfg.cli,
	}
	r.ctx, r.cancel = context.WithCancel(cfg.ctx)

	req, err := http.NewRequestWithContext(r.ctx, http.MethodGet, uri, nil)
	if err != nil {
		r.cancel()
		return nil, fmt.Errorf("httpio: could not create HTTP request: %w", err)
	}
	if cfg.auth.usr != "" || cfg.auth.pwd != "" {
		req.SetBasicAuth(cfg.auth.usr, cfg.auth.pwd)
	}
	r.req = req.Clone(r.ctx)

	hdr, err := r.cli.Head(r.req.URL.String())
	if err != nil {
		r.cancel()
		return nil, fmt.Errorf("httpio: could not send HEAD request: %w", err)
	}
	defer hdr.Body.Close()
	_, _ = io.Copy(io.Discard, hdr.Body)

	if hdr.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("httpio: invalid HEAD response code=%v", hdr.StatusCode)
	}

	if hdr.Header.Get("accept-ranges") != "bytes" {
		return nil, fmt.Errorf("httpio: invalid HEAD response: %w", errAcceptRange)
	}

	r.len = hdr.ContentLength
	r.etag = hdr.Header.Get("Etag")
	r.r = io.NewSectionReader(r, 0, r.len)

	r.req.Header.Set("Range", "")
	r.pool = sync.Pool{
		New: func() any {
			return r.req.Clone(r.ctx)
		},
	}

	return r, nil
}

// Size returns the number of bytes available for reading via ReadAt.
func (r *Reader) Size() int64 {
	return r.len
}

// Name returns the name of the file as presented to Open.
func (r *Reader) Name() string {
	return r.req.URL.String()
}

// Close implements the io.Closer interface.
func (r *Reader) Close() error {
	r.cancel()
	r.cli = nil
	r.req = nil
	return nil
}

// Read implements the io.Reader interface.
func (r *Reader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

// Seek implements the io.Seeker interface.
func (r *Reader) Seek(offset int64, whence int) (int64, error) {
	return r.r.Seek(offset, whence)
}

// ReadAt implements the io.ReaderAt interface.
func (r *Reader) ReadAt(p []byte, off int64) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	rng := rng(off, off+int64(len(p))-1)
	req := r.getReq(rng)
	defer r.pool.Put(req)

	resp, err := r.cli.Do(req)
	if err != nil {
		return 0, fmt.Errorf("httpio: could not send GET request: %w", err)
	}
	defer resp.Body.Close()

	n, _ := io.ReadFull(resp.Body, p)

	if etag := resp.Header.Get("Etag"); etag != r.etag {
		return n, fmt.Errorf("httpio: resource changed")
	}

	switch resp.StatusCode {
	case http.StatusPartialContent:
		// ok.
	case http.StatusRequestedRangeNotSatisfiable:
		return 0, io.EOF
	default:
		return n, fmt.Errorf("httpio: invalid GET response: code=%v", resp.StatusCode)
	}

	if int64(len(p)) > r.len {
		return n, io.EOF
	}

	return n, nil
}

func (r *Reader) getReq(rng string) *http.Request {
	o := r.pool.Get().(*http.Request)
	o.Header = r.req.Header.Clone()
	o.Header["Range"][0] = rng
	return o
}

func rng(beg, end int64) string {
	return "bytes=" + strconv.Itoa(int(beg)) + "-" + strconv.Itoa(int(end))
}

var (
	_ io.Reader   = (*Reader)(nil)
	_ io.Seeker   = (*Reader)(nil)
	_ io.ReaderAt = (*Reader)(nil)
	_ io.Closer   = (*Reader)(nil)
)
