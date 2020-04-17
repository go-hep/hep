// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"net/http"
)

type hwriter struct {
	body *bytes.Buffer
	code int
	hdr  http.Header
}

func newResponseWriter() *hwriter {
	return &hwriter{
		hdr:  make(http.Header),
		body: new(bytes.Buffer),
	}
}

func (w *hwriter) Header() http.Header         { return w.hdr }
func (w *hwriter) Write(p []byte) (int, error) { return w.body.Write(p) }
func (w *hwriter) WriteHeader(code int)        { w.code = code }
func (w *hwriter) reset() {
	w.body.Reset()
	for k := range w.hdr {
		delete(w.hdr, k)
	}
}

var (
	_ http.ResponseWriter = (*hwriter)(nil)
)
