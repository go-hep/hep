// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"encoding/binary"
	"io"
)

type wbuff struct {
	p []byte // buffer of data to write on
	c int    // current position in buffer of data
}

func (w *wbuff) Write(p []byte) (int, error) {
	if w.c >= len(w.p) {
		return 0, io.EOF
	}
	n := copy(w.p[w.c:], p)
	w.c += n
	return n, nil
}

func (w *wbuff) WriteByte(p byte) error {
	if w.c >= len(w.p) {
		return io.EOF
	}
	w.p[w.c] = p
	w.c++
	return nil
}

// WBuffer is a write-only ROOT buffer for streaming.
type WBuffer struct {
	w      *wbuff
	err    error
	offset uint32
	refs   map[int64]interface{}
	sictx  StreamerInfoContext
}

func NewWBuffer(data []byte, refs map[int64]interface{}, offset uint32, ctx StreamerInfoContext) *WBuffer {
	if refs == nil {
		refs = make(map[int64]interface{})
	}

	return &WBuffer{
		w:      &wbuff{p: data, c: 0},
		refs:   refs,
		offset: offset,
		sictx:  ctx,
	}
}

func (w *WBuffer) write(v []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write(v)
}

func (w *WBuffer) WriteI8(v int8) {
	if w.err != nil {
		return
	}
	j := byte(v)
	w.err = w.w.WriteByte(j)
}

func (w *WBuffer) WriteI16(v int16) {
	if w.err != nil {
		return
	}
	beg := w.w.c
	end := w.w.c + 2
	binary.BigEndian.PutUint16(w.w.p[beg:end], uint16(v))
	w.w.c += 2

}

func (w *WBuffer) WriteI32(v int32) {
	if w.err != nil {
		return
	}
	w.err = binary.Write(w.w, binary.BigEndian, v)
}

func (w *WBuffer) WriteI64(v int64) {
	if w.err != nil {
		return
	}
	w.err = binary.Write(w.w, binary.BigEndian, v)
}

func (w *WBuffer) WriteU8(v uint8) {
	if w.err != nil {
		return
	}
	j := byte(v)
	w.err = w.w.WriteByte(j)
}

func (w *WBuffer) WriteU16(v uint16) {
	if w.err != nil {
		return
	}
	beg := w.w.c
	end := w.w.c + 2
	binary.BigEndian.PutUint16(w.w.p[beg:end], v)
	w.w.c += 2
}

func (w *WBuffer) WriteU32(v uint32) {
	if w.err != nil {
		return
	}
	w.err = binary.Write(w.w, binary.BigEndian, v)
}

func (w *WBuffer) WriteU64(v uint64) {
	if w.err != nil {
		return
	}
	w.err = binary.Write(w.w, binary.BigEndian, v)
}

func (w *WBuffer) WriteF32(v float32) {
	if w.err != nil {
		return
	}
	w.err = binary.Write(w.w, binary.BigEndian, v)
}

func (w *WBuffer) WriteF64(v float64) {
	if w.err != nil {
		return
	}
	w.err = binary.Write(w.w, binary.BigEndian, v)
}

func (w *WBuffer) WriteBool(v bool) {
	if w.err != nil {
		return
	}
	var o byte = 0
	if v {
		o = 1
	}
	w.err = w.w.WriteByte(o)
}

func (w *WBuffer) WriteString(v string) {
	if w.err != nil {
		return
	}
	l := len(v)
	if l < 255 {
		w.WriteU8(uint8(l))
		w.write([]byte(v))
	} else {
		w.WriteU8(255)
		w.WriteU32(uint32(l))
		w.write([]byte(v))
	}
}
