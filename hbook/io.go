// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import "io"

type rbuffer struct {
	p []byte // buffer of data to read from
	c int    // current position in buffer of data
}

func newRBuffer(p []byte) *rbuffer {
	return &rbuffer{p: p}
}

func (r *rbuffer) Read(p []byte) (int, error) {
	if r.c >= len(r.p) {
		return 0, io.EOF
	}
	n := copy(p, r.p[r.c:])
	r.c += n
	return n, nil
}

func (r *rbuffer) Bytes() []byte { return r.p[r.c:] }

func (r *rbuffer) next(n int) []byte {
	m := len(r.p[r.c:])
	if n > m {
		n = m
	}
	p := r.p[r.c : r.c+n]
	r.c += n
	return p
}
