// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes

import (
	"bytes"
	"encoding/binary"
	"io"
	"math"
	"reflect"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/root"
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

// grow grows the buffer's capacity, if necessary, to guarantee space foranother n bytes.
// After grow(n), at least n bytes can be written to the buffer without
// another allocation.
// If n is negative, grow will panic.
func (w *wbuff) grow(n int) {
	if n < 0 {
		panic(errors.Errorf("rbytes: negative count"))
	}
	if n == 0 {
		return
	}
	plen := len(w.p)
	pcap := cap(w.p)
	if plen+n < pcap {
		w.p = w.p[:plen+n]
		return
	}
	w.p = append(w.p, make([]byte, pcap+n)...)
}

// WBuffer is a write-only ROOT buffer for streaming.
type WBuffer struct {
	w      wbuff
	err    error
	offset uint32
	refs   map[interface{}]int64
	sictx  StreamerInfoContext
}

func NewWBuffer(data []byte, refs map[interface{}]int64, offset uint32, ctx StreamerInfoContext) *WBuffer {
	if refs == nil {
		refs = make(map[interface{}]int64)
	}

	return &WBuffer{
		w:      wbuff{p: data, c: 0},
		refs:   refs,
		offset: offset,
		sictx:  ctx,
	}
}

func (w *WBuffer) Grow(n int)     { w.w.grow(n) }
func (w *WBuffer) buffer() []byte { return w.w.p[:w.w.c] }
func (w *WBuffer) Bytes() []byte  { return w.w.p[:w.w.c] }

func (w *WBuffer) Err() error       { return w.err }
func (w *WBuffer) SetErr(err error) { w.err = err }

func (w *WBuffer) Pos() int64 {
	return int64(w.w.c) + int64(w.offset)
}

func (w *WBuffer) SetPos(pos int64) { w.setPos(pos) }
func (w *WBuffer) setPos(pos int64) {
	pos -= int64(w.offset)
	w.w.c = int(pos)
}

func (w *WBuffer) Write(p []byte) (int, error) {
	if w.err != nil {
		return 0, w.err
	}
	w.w.grow(len(p))
	n, err := w.w.Write(p)
	w.err = err
	return n, w.err
}

func (w *WBuffer) SetByteCount(beg int64, class string) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	cur := w.Pos()
	bcnt := cur - beg - 4
	w.setPos(beg)
	w.WriteU32(uint32(bcnt | kByteCountMask))
	w.setPos(cur)

	return int(bcnt + 4), w.err
}

func (w *WBuffer) WriteVersion(vers int16) {
	if w.err != nil {
		return
	}
	w.w.grow(6)
	w.writeU32(0) // byte-count placeholder
	w.writeU16(uint16(vers))
}

func (w *WBuffer) WriteObjectAny(obj root.Object) error {
	if w.err != nil {
		return w.err
	}

	if reflect.ValueOf(obj).IsNil() {
		w.WriteU32(0) // NULL pointer
		return w.err
	}

	pos := w.Pos()
	w.WriteU32(0) // placeholder for bytecount.

	bcnt, err := w.WriteClass(pos, obj)
	if err != nil {
		w.err = err
		return w.err
	}
	end := w.Pos()
	w.setPos(pos)
	w.writeU32(bcnt)
	w.setPos(end)

	return w.err
}

func (w *WBuffer) WriteClass(beg int64, obj root.Object) (uint32, error) {
	if w.err != nil {
		return 0, w.err
	}

	start := w.Pos()
	if ref64, dup := w.refs[obj]; dup {
		// we've already seen this value.
		w.WriteU32(uint32(ref64) | kClassMask)
		bcnt := w.Pos() - start
		return uint32(bcnt | kByteCountMask), w.err
	}

	class := obj.Class()
	ref64, ok := w.refs[class]
	if !ok {
		// first time we see this type
		w.WriteU32(uint32(kNewClassTag))
		w.WriteCString(class)
		w.refs[class] = (start + kMapOffset) | kClassMask

		// add to refs before writing value, to handle self reference
		w.refs[obj] = beg + kMapOffset

		mobj := obj.(Marshaler)
		if _, err := mobj.MarshalROOT(w); err != nil {
			w.err = err
			return 0, w.err
		}

		bcnt := w.Pos() - start
		return uint32(bcnt | kByteCountMask), w.err
	}

	// first time we see this value
	w.WriteU32(uint32(ref64) | kClassMask)
	if _, err := obj.(Marshaler).MarshalROOT(w); err != nil {
		w.err = err
		return 0, w.err
	}

	w.refs[obj] = beg + kMapOffset
	bcnt := w.Pos() - start
	return uint32(bcnt | kByteCountMask), w.err
}

func (w *WBuffer) write(v []byte) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v))
	_, w.err = w.w.Write(v)
}

func (w *WBuffer) WriteI8(v int8) {
	if w.err != nil {
		return
	}
	w.w.grow(1)
	w.writeI8(v)
}

func (w *WBuffer) writeI8(v int8) {
	w.w.p[w.w.c] = byte(v)
	w.w.c++
}

func (w *WBuffer) WriteI16(v int16) {
	if w.err != nil {
		return
	}
	const sz = 2
	w.w.grow(sz)
	w.writeI16(v)
}

func (w *WBuffer) writeI16(v int16) {
	const sz = 2
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint16(w.w.p[beg:end], uint16(v))
	w.w.c += sz
}

func (w *WBuffer) WriteI32(v int32) {
	if w.err != nil {
		return
	}
	const sz = 4
	w.w.grow(sz)
	w.writeI32(v)
}

func (w *WBuffer) writeI32(v int32) {
	const sz = 4
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint32(w.w.p[beg:end], uint32(v))
	w.w.c += sz
}

func (w *WBuffer) WriteI64(v int64) {
	if w.err != nil {
		return
	}
	const sz = 8
	w.w.grow(sz)
	w.writeI64(v)
}

func (w *WBuffer) writeI64(v int64) {
	const sz = 8
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint64(w.w.p[beg:end], uint64(v))
	w.w.c += sz
}

func (w *WBuffer) WriteU8(v uint8) {
	if w.err != nil {
		return
	}
	w.w.grow(1)
	w.writeU8(v)
}

func (w *WBuffer) writeU8(v uint8) {
	w.w.p[w.w.c] = v
	w.w.c++
}

func (w *WBuffer) WriteU16(v uint16) {
	if w.err != nil {
		return
	}
	const sz = 2
	w.w.grow(sz)
	w.writeU16(v)
}

func (w *WBuffer) writeU16(v uint16) {
	const sz = 2
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint16(w.w.p[beg:end], uint16(v))
	w.w.c += sz
}

func (w *WBuffer) WriteU32(v uint32) {
	if w.err != nil {
		return
	}
	const sz = 4
	w.w.grow(sz)
	w.writeU32(v)
}

func (w *WBuffer) writeU32(v uint32) {
	const sz = 4
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint32(w.w.p[beg:end], v)
	w.w.c += sz
}

func (w *WBuffer) WriteU64(v uint64) {
	if w.err != nil {
		return
	}
	const sz = 8
	w.w.grow(sz)
	w.writeU64(v)
}

func (w *WBuffer) writeU64(v uint64) {
	const sz = 8
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint64(w.w.p[beg:end], v)
	w.w.c += sz
}

func (w *WBuffer) WriteF32(v float32) {
	if w.err != nil {
		return
	}
	const sz = 4
	w.w.grow(sz)
	w.writeF32(v)
}

func (w *WBuffer) writeF32(v float32) {
	const sz = 4
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint32(w.w.p[beg:end], math.Float32bits(v))
	w.w.c += sz
}

func (w *WBuffer) WriteF64(v float64) {
	if w.err != nil {
		return
	}
	const sz = 8
	w.w.grow(sz)
	w.writeF64(v)
}

func (w *WBuffer) writeF64(v float64) {
	const sz = 8
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint64(w.w.p[beg:end], math.Float64bits(v))
	w.w.c += sz
}

func (w *WBuffer) WriteBool(v bool) {
	var o uint8
	if v {
		o = 1
	}
	w.w.grow(1)
	w.writeU8(o)
}

func (w *WBuffer) WriteString(v string) {
	if w.err != nil {
		return
	}
	l := len(v)
	if l < 255 {
		w.w.grow(1 + l)
		w.writeU8(uint8(l))
		w.write([]byte(v))
		return
	}
	w.w.grow(1 + 4 + l)
	w.writeU8(255)
	w.writeU32(uint32(l))
	w.write([]byte(v))
}

func (w *WBuffer) WriteCString(v string) {
	if w.err != nil {
		return
	}
	b := []byte(v)
	i := bytes.Index(b, []byte{0})
	switch {
	case i < 0:
		b = append(b, 0)
		w.write(b)
	default:
		b = b[:i+1]
		w.write(b)
	}
}

func (w *WBuffer) WriteStaticArrayI32(v []int32) {
	if w.err != nil {
		return
	}
	w.w.grow(4 + 4*len(v))
	w.writeI32(int32(len(v)))
	for _, v := range v {
		w.writeI32(v)
	}
}

func (w *WBuffer) WriteFastArrayBool(v []bool) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v))
	for _, v := range v {
		var b byte = 0
		if v {
			b = 1
		}
		w.writeU8(b)
	}
}

func (w *WBuffer) WriteFastArrayI8(v []int8) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v))
	for _, v := range v {
		w.writeI8(v)
	}
}

func (w *WBuffer) WriteFastArrayI16(v []int16) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 2)
	for _, v := range v {
		w.writeI16(v)
	}
}

func (w *WBuffer) WriteFastArrayI32(v []int32) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 4)
	for _, v := range v {
		w.writeI32(v)
	}
}

func (w *WBuffer) WriteFastArrayI64(v []int64) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 8)
	for _, v := range v {
		w.writeI64(v)
	}
}

func (w *WBuffer) WriteFastArrayU8(v []uint8) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v))
	for _, v := range v {
		w.writeU8(v)
	}
}

func (w *WBuffer) WriteFastArrayU16(v []uint16) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 2)
	for _, v := range v {
		w.writeU16(v)
	}
}

func (w *WBuffer) WriteFastArrayU32(v []uint32) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 4)
	for _, v := range v {
		w.writeU32(v)
	}
}

func (w *WBuffer) WriteFastArrayU64(v []uint64) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 8)
	for _, v := range v {
		w.writeU64(v)
	}
}

func (w *WBuffer) WriteFastArrayF32(v []float32) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 4)
	for _, v := range v {
		w.writeF32(v)
	}
}

func (w *WBuffer) WriteFastArrayF64(v []float64) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v) * 8)
	for _, v := range v {
		w.writeF64(v)
	}
}
