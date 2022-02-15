// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rbytes"
)

type decoder struct {
	r    *rbytes.RBuffer
	si   *StreamerInfo
	kind rbytes.StreamKind
	rops []rstreamer
}

func newDecoder(r *rbytes.RBuffer, si *StreamerInfo, kind rbytes.StreamKind, ops []rstreamer) (*decoder, error) {
	return &decoder{r, si, kind, ops}, nil
}

func (dec *decoder) DecodeROOT(ptr interface{}) error {
	rv := reflect.ValueOf(ptr)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rdict: invalid kind (got=%T, want=pointer)", ptr)
	}

	var (
		typename = dec.si.Name()
		typevers = int16(dec.si.ClassVersion())
		hdr      = dec.r.ReadHeader(typename)
	)

	if hdr.Vers != typevers {
		dec.r.SetErr(fmt.Errorf("rdict: inconsistent ROOT version type=%q (got=%d, want=%d)",
			typename, hdr.Vers, typevers,
		))
		return dec.r.Err()
	}

	for i, op := range dec.rops {
		err := op.rstream(dec.r, ptr)
		if err != nil {
			return fmt.Errorf("rdict: could not read element %d from %q: %w", i, typename, err)
		}
	}

	dec.r.CheckHeader(hdr)
	if err := dec.r.Err(); err != nil {
		return fmt.Errorf("rdict: invalid bytecount for %q: %w", typename, err)
	}

	return nil
}

var (
	_ rbytes.Decoder = (*decoder)(nil)
)

type rstreamerInfo struct {
	recv interface{}
	rops []rstreamer
	kind rbytes.StreamKind
	si   *StreamerInfo
}

func newRStreamerInfo(si *StreamerInfo, kind rbytes.StreamKind, rops []rstreamer) (*rstreamerInfo, error) {
	return &rstreamerInfo{
		recv: nil,
		rops: rops,
		kind: kind,
		si:   si,
	}, nil
}

func (rr *rstreamerInfo) Bind(recv interface{}) error {
	rv := reflect.ValueOf(recv)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("rdict: invalid kind (got=%T, want=pointer)", recv)
	}
	rr.recv = recv
	if recv, ok := recv.(rbytes.Unmarshaler); ok && rr.kind == rbytes.ObjectWise {
		// FIXME(sbinet): handle mbr-/obj-wise
		rr.rops = []rstreamer{{
			op: func(r *rbytes.RBuffer, _ interface{}, _ *streamerConfig) error {
				return recv.UnmarshalROOT(r)
			},
			cfg: nil,
		}}
		return nil
	}
	if len(rr.rops) == 1 {
		se := rr.rops[0].cfg.descr.elem
		if se.Name() == "This" ||
			strings.HasPrefix(se.TypeName(), "vector<") {
			// binding directly to 'recv'. assume no offset is to be applied
			rr.rops[0].cfg.offset = -1
		}
	}
	return nil
}

func (rr *rstreamerInfo) RStreamROOT(r *rbytes.RBuffer) error {
	for i, op := range rr.rops {
		err := op.rstream(r, rr.recv)
		if err != nil {
			typename := rr.si.Name()
			return fmt.Errorf("rdict: could not read element %d from %q: %w", i, typename, err)
		}
	}
	return nil
}

var (
	_ rbytes.RStreamer = (*rstreamerInfo)(nil)
	_ rbytes.Binder    = (*rstreamerInfo)(nil)
)
