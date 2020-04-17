// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdenc // import "go-hep.org/x/hep/xrootd/internal/xrdenc"

import "encoding/binary"

// WBuffer encodes values to a buffer according to the XRootD protocol.
type WBuffer struct {
	buf []byte
}

func (w *WBuffer) Bytes() []byte { return w.buf }

func (w *WBuffer) WriteU8(v uint8) {
	w.buf = append(w.buf, v)
}

func (w *WBuffer) WriteU16(v uint16) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], v)
	w.buf = append(w.buf, buf[:]...)
}

func (w *WBuffer) WriteI32(v int32) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v))
	w.buf = append(w.buf, buf[:]...)
}

func (w *WBuffer) WriteI64(v int64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(v))
	w.buf = append(w.buf, buf[:]...)
}

func (w *WBuffer) WriteBool(v bool) {
	if v {
		w.buf = append(w.buf, 1)
		return
	}
	w.buf = append(w.buf, 0)
}

func (w *WBuffer) WriteLen(n int) {
	w.WriteI32(int32(n))
}

func (w *WBuffer) WriteBytes(vs []byte) {
	w.buf = append(w.buf, vs...)
}

func (w *WBuffer) WriteStr(str string) {
	w.WriteLen(len(str))
	w.WriteBytes([]byte(str))
}

func (w *WBuffer) Next(n int) {
	w.buf = append(w.buf, make([]byte, n)...)
}

// RBuffer decodes values from a buffer according to the XRootD protocol.
type RBuffer struct {
	buf []byte
	pos int
}

func NewRBuffer(data []byte) *RBuffer {
	return &RBuffer{buf: data}
}

func (r *RBuffer) ReadU8() uint8 {
	o := r.buf[r.pos]
	r.pos++
	return o
}

func (r *RBuffer) Len() int {
	return len(r.buf) - r.pos
}

func (r *RBuffer) ReadU16() uint16 {
	beg := r.pos
	end := r.pos + 2
	r.pos += 2
	o := binary.BigEndian.Uint16(r.buf[beg:end])
	return o
}

func (r *RBuffer) ReadI32() int32 {
	beg := r.pos
	end := r.pos + 4
	r.pos += 4
	o := binary.BigEndian.Uint32(r.buf[beg:end])
	return int32(o)
}

func (r *RBuffer) ReadI64() int64 {
	beg := r.pos
	end := r.pos + 8
	r.pos += 8
	o := binary.BigEndian.Uint64(r.buf[beg:end])
	return int64(o)
}

func (r *RBuffer) ReadBool() bool {
	r.pos += 1
	return r.buf[r.pos] != 0
}

func (r *RBuffer) ReadLen() int {
	return int(r.ReadI32())
}

func (r *RBuffer) ReadBytes(data []byte) {
	n := len(data)
	beg := r.pos
	end := r.pos + n
	copy(data, r.buf[beg:end])
	r.pos += n
	return
}

func (r *RBuffer) ReadStr() string {
	n := r.ReadLen()
	beg := r.pos
	end := r.pos + n
	r.pos += n
	return string(r.buf[beg:end])
}

func (r *RBuffer) Skip(n int) {
	r.pos += n
}

func (r *RBuffer) Bytes() []byte {
	return r.buf[r.pos:]
}

func (r *RBuffer) Pos() int {
	return r.pos
}
