// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
)

type encoder struct {
	w    *rbytes.WBuffer
	si   *StreamerInfo
	kind rbytes.StreamKind
	wops []wstreamer
}

func newEncoder(w *rbytes.WBuffer, si *StreamerInfo, kind rbytes.StreamKind, ops []wstreamer) (*encoder, error) {
	return &encoder{w, si, kind, ops}, nil
}

func (enc *encoder) EncodeROOT(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rdict: invalid kind (got=%T, want=pointer)", ptr)
	}

	var (
		typename = enc.si.Name()
		typevers = int16(enc.si.ClassVersion())
		pos      = enc.w.WriteVersion(typevers)
		err      error
	)

	for i, op := range enc.wops {
		_, err = op.wstream(enc.w, ptr)
		if err != nil {
			return fmt.Errorf("rdict: could not write element %d from %q: %w", i, typename, err)
		}
	}

	_, err = enc.w.SetByteCount(pos, typename)
	if err != nil {
		return fmt.Errorf("rdict: could not set byte count for %q: %w", typename, err)
	}
	return nil
}

var (
	_ rbytes.Encoder = (*encoder)(nil)
)

type wstreamerInfo struct {
	recv interface{}
	wops []wstreamer
	kind rbytes.StreamKind
	si   *StreamerInfo
}

func newWStreamerInfo(si *StreamerInfo, kind rbytes.StreamKind, wops []wstreamer) (*wstreamerInfo, error) {
	return &wstreamerInfo{
		recv: nil,
		wops: wops,
		kind: kind,
		si:   si,
	}, nil
}

func (ww *wstreamerInfo) Bind(recv interface{}) error {
	rv := reflect.ValueOf(recv)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rdict: invalid kind (got=%T, want=pointer)", recv)
	}
	ww.recv = recv
	if len(ww.wops) == 1 && ww.wops[0].cfg.descr.elem.Name() == "This" {
		// binding directly to 'recv'. assume no offset is to be applied
		ww.wops[0].cfg.offset = -1
	}
	return nil
}

func (ww *wstreamerInfo) WStreamROOT(w *rbytes.WBuffer) error {
	for i, op := range ww.wops {
		_, err := op.wstream(w, ww.recv)
		if err != nil {
			typename := ww.si.Name()
			return fmt.Errorf("rdict: could not write element %d from %q: %w", i, typename, err)
		}
	}
	return nil
}

var (
	_ rbytes.WStreamer = (*wstreamerInfo)(nil)
	_ rbytes.Binder    = (*wstreamerInfo)(nil)
)
