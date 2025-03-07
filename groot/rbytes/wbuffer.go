// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbytes

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"reflect"

	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rvers"
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
		panic(fmt.Errorf("rbytes: negative count"))
	}
	if n == 0 {
		return
	}
	var (
		plen = len(w.p)
		pcap = cap(w.p)
		nlen = plen + n
	)
	if nlen < pcap {
		w.p = w.p[:nlen]
		return
	}
	if pcap < nlen {
		const sz = 4 * 1024
		pcap = ((nlen / sz) + 1) * sz
	}
	w.p = append(w.p[:cap(w.p)], make([]byte, 2*pcap)...)
	w.p = w.p[:nlen:cap(w.p)]
}

// WBuffer is a write-only ROOT buffer for streaming.
type WBuffer struct {
	w      wbuff
	err    error
	offset uint32
	refs   map[any]int64
	sictx  StreamerInfoContext
}

func NewWBuffer(data []byte, refs map[any]int64, offset uint32, ctx StreamerInfoContext) *WBuffer {
	if refs == nil {
		refs = make(map[any]int64)
	}

	return &WBuffer{
		w:      wbuff{p: data, c: 0},
		refs:   refs,
		offset: offset,
		sictx:  ctx,
	}
}

// StreamerInfo returns the named StreamerInfo.
// If version is negative, the latest version should be returned.
func (w *WBuffer) StreamerInfo(name string, version int) (StreamerInfo, error) {
	if w.sictx == nil {
		return nil, fmt.Errorf("rbytes: no streamers")
	}
	return w.sictx.StreamerInfo(name, version)
}

func (w *WBuffer) Grow(n int)     { w.w.grow(n) }
func (w *WBuffer) buffer() []byte { return w.w.p[:w.w.c] }
func (w *WBuffer) Bytes() []byte  { return w.w.p[:w.w.c] }

func (w *WBuffer) Err() error       { return w.err }
func (w *WBuffer) SetErr(err error) { w.err = err }

func (w *WBuffer) Pos() int64 {
	return int64(w.w.c) + int64(w.offset)
}

func (w *WBuffer) Len() int64 {
	return int64(w.w.c)
}

func (w *WBuffer) SetPos(pos int64) { w.setPos(pos) }
func (w *WBuffer) setPos(pos int64) {
	pos -= int64(w.offset)
	w.w.c = int(pos)
}

func (w *WBuffer) DumpHex(n int) {
	buf := w.buffer()
	if len(buf) > n {
		buf = buf[:n]
	}
	fmt.Printf("--- hex --- (pos=%d len=%d end=%d)\n%s\n", w.Pos(), n, w.Len(), string(hex.Dump(buf)))
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

func (w *WBuffer) WriteHeader(typename string, vers int16) Header {
	hdr := Header{
		Name: typename,
		Pos:  w.Pos(),
		Vers: vers,
	}
	if w.err != nil {
		return hdr
	}

	w.w.grow(6)
	w.writeU32(0) // byte-count placeholder
	w.writeU16(uint16(vers))

	return hdr
}

func (w *WBuffer) SetHeader(hdr Header) (int, error) {
	if w.err != nil {
		return 0, w.err
	}

	cur := w.Pos()
	cnt := cur - hdr.Pos - 4
	w.setPos(hdr.Pos)
	w.WriteU32(uint32(cnt | kByteCountMask))
	w.setPos(cur)

	return int(cnt + 4), w.err
}

func (w *WBuffer) WriteObject(obj Marshaler) {
	if w.err != nil {
		return
	}

	if v := reflect.ValueOf(obj); (v == reflect.Value{}) || v.IsNil() {
		w.WriteU32(0) // NULL pointer
		return
	}

	_, w.err = obj.MarshalROOT(w)
}

func (w *WBuffer) WriteObjectAny(obj root.Object) {
	if w.err != nil {
		return
	}

	if v := reflect.ValueOf(obj); (v == reflect.Value{}) || v.IsNil() {
		w.WriteU32(0) // NULL pointer
		return
	}

	pos := w.Pos()
	w.WriteU32(0) // placeholder for bytecount.

	bcnt := w.WriteClass(pos, obj)
	if w.err != nil {
		return
	}
	end := w.Pos()
	w.setPos(pos)
	w.writeU32(bcnt)
	w.setPos(end)
}

func (w *WBuffer) WriteClass(beg int64, obj root.Object) uint32 {
	if w.err != nil {
		return 0
	}

	start := w.Pos()
	if ref64, dup := w.refs[obj]; dup {
		// we've already seen this value.
		w.WriteU32(uint32(ref64))
		bcnt := w.Pos() - start
		return uint32(bcnt | kByteCountMask)
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
			return 0
		}

		bcnt := w.Pos() - start
		return uint32(bcnt | kByteCountMask)
	}

	// first time we see this value
	w.WriteU32(uint32(ref64) | kClassMask)
	if _, err := obj.(Marshaler).MarshalROOT(w); err != nil {
		return 0
	}

	w.refs[obj] = beg + kMapOffset
	bcnt := w.Pos() - start
	return uint32(bcnt | kByteCountMask)
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

func (w *WBuffer) writeI32(v int32) {
	const sz = 4
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint32(w.w.p[beg:end], uint32(v))
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

func (w *WBuffer) writeU16(v uint16) {
	const sz = 2
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint16(w.w.p[beg:end], uint16(v))
	w.w.c += sz
}

func (w *WBuffer) writeU32(v uint32) {
	const sz = 4
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint32(w.w.p[beg:end], v)
	w.w.c += sz
}

func (w *WBuffer) writeF32(v float32) {
	const sz = 4
	beg := w.w.c
	end := w.w.c + sz
	binary.BigEndian.PutUint32(w.w.p[beg:end], math.Float32bits(v))
	w.w.c += sz
}

func (w *WBuffer) WriteF16(v root.Float16, elm StreamerElement) {
	if w.err != nil {
		return
	}
	switch {
	case elm != nil && elm.Factor() != 0:
		// range specified.
		// normalize the float64 to the range and convert to an integer,
		// according to the provided scaling factor.
		var (
			x   = float64(v)
			min = float64(elm.XMin())
			max = float64(elm.XMax())
		)
		if x < min {
			x = min
		}
		if x > max {
			x = max
		}
		u32 := uint32(0.5 + elm.Factor()*(x-min))
		w.w.grow(4)
		w.writeU32(u32)
	default:
		var nbits uint32
		if elm != nil {
			nbits = uint32(elm.XMin())
		}
		if nbits == 0 {
			nbits = 12
		}
		// no range, but nbits.
		// truncate mantissa to nbits.
		u32 := math.Float32bits(float32(v))
		exp := uint8(0x000000ff & ((u32 << 1) >> 24))
		man := uint16(((1 << (nbits + 1)) - 1) & (u32 >> (23 - nbits - 1)))
		man++
		man = man >> 1
		if man&1<<nbits != 0 {
			man = (1 << nbits) - 1
		}
		if v < 0 {
			man |= 1 << (nbits + 1)
		}
		w.w.grow(3)
		w.writeU8(exp)
		w.writeU16(man)
	}
}

func (w *WBuffer) WriteD32(v root.Double32, elm StreamerElement) {
	if w.err != nil {
		return
	}
	switch {
	case elm != nil && elm.Factor() != 0:
		// range specified.
		// normalize the float64 to the range and convert to an integer,
		// according to the provided scaling factor.
		var (
			min = root.Double32(elm.XMin())
			max = root.Double32(elm.XMax())
		)
		if v < min {
			v = min
		}
		if v > max {
			v = max
		}
		u32 := uint32(0.5 + elm.Factor()*float64(v-min))
		w.w.grow(4)
		w.writeU32(u32)
	default:
		var nbits uint32
		if elm != nil {
			nbits = uint32(elm.XMin())
		}
		switch {
		case nbits == 0:
			// no range, no bits: convert to float32
			w.w.grow(4)
			w.writeF32(float32(v))
		default:
			// no range, but nbits.
			// truncate mantissa to nbits.
			u32 := math.Float32bits(float32(v))
			exp := uint8(0x000000ff & ((u32 << 1) >> 24))
			man := uint16(((1 << (nbits + 1)) - 1) & (u32 >> (23 - nbits - 1)))
			man++
			man = man >> 1
			if man&1<<nbits != 0 {
				man = (1 << nbits) - 1
			}
			if v < 0 {
				man |= 1 << (nbits + 1)
			}
			w.w.grow(3)
			w.writeU8(exp)
			w.writeU16(man)
		}
	}
}

func (w *WBuffer) WriteBool(v bool) {
	var o uint8
	if v {
		o = 1
	}
	w.w.grow(1)
	w.writeU8(o)
}

func (w *WBuffer) WriteStdString(v string) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("string", rvers.StreamerInfo)
	w.WriteString(v)
	_, err := w.SetHeader(hdr)
	if err != nil {
		w.SetErr(err)
	}
}

func (w *WBuffer) WriteString(v string) {
	if w.err != nil {
		return
	}
	l := len(v)
	if l < 255 {
		w.w.grow(1 + l)
		w.writeU8(uint8(l))
		if l > 0 {
			w.write([]byte(v))
		}
		return
	}
	w.w.grow(1 + 4 + l)
	w.writeU8(255)
	w.writeU32(uint32(l))
	w.write([]byte(v))
}

func (w *WBuffer) writeString(v string) {
	if w.err != nil {
		return
	}
	l := len(v)
	if l < 255 {
		w.writeU8(uint8(l))
		w.write([]byte(v))
		return
	}
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

func (w *WBuffer) WriteArrayBool(v []bool) {
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

func (w *WBuffer) WriteArrayI8(v []int8) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v))
	for _, v := range v {
		w.writeI8(v)
	}
}

func (w *WBuffer) WriteArrayU8(v []uint8) {
	if w.err != nil {
		return
	}
	w.w.grow(len(v))
	for _, v := range v {
		w.writeU8(v)
	}
}

func (w *WBuffer) WriteArrayF16(v []root.Float16, elm StreamerElement) {
	for _, v := range v {
		w.WriteF16(v, elm)
	}
}

func (w *WBuffer) WriteArrayD32(v []root.Double32, elm StreamerElement) {
	for _, v := range v {
		w.WriteD32(v, elm)
	}
}

func (w *WBuffer) strlen(v string) int {
	l := len(v)
	if l < 255 {
		return 1 + l
	}
	return 1 + 4 + l
}

func (w *WBuffer) WriteArrayString(v []string) {
	if w.err != nil {
		return
	}
	n := 0
	for _, v := range v {
		n += w.strlen(v)
	}
	w.w.grow(n)

	for _, v := range v {
		w.writeString(v)
	}
}

func (w *WBuffer) WriteStdVectorStrs(v []string) {
	if w.err != nil {
		return
	}

	hdr := w.WriteHeader("vector<string>", rvers.StreamerInfo)
	w.WriteI32(int32(len(v)))

	n := 0
	for _, v := range v {
		n += w.strlen(v)
	}
	w.w.grow(n)

	for _, v := range v {
		w.writeString(v)
	}
	_, _ = w.SetHeader(hdr)
}

func (w *WBuffer) WriteStdBitset(v []uint8) {
	n := len(v)
	w.w.grow(n)
	for i := range v {
		w.writeU8(v[n-1-i])
	}
}

var (
	_ StreamerInfoContext = (*WBuffer)(nil)
)
