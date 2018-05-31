// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdenc // import "go-hep.org/x/hep/xrootd/internal/xrdenc"

import "encoding/binary"

// Encoder encodes values to a buffer according to the XRootD protocol.
type Encoder struct {
	buf []byte
}

func (enc *Encoder) Bytes() []byte { return enc.buf }

func (enc *Encoder) WriteU8(v uint8) {
	enc.buf = append(enc.buf, v)
}

func (enc *Encoder) WriteU16(v uint16) {
	var buf [2]byte
	binary.BigEndian.PutUint16(buf[:], v)
	enc.buf = append(enc.buf, buf[:]...)
}

func (enc *Encoder) WriteI32(v int32) {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], uint32(v))
	enc.buf = append(enc.buf, buf[:]...)
}

func (enc *Encoder) WriteI64(v int64) {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(v))
	enc.buf = append(enc.buf, buf[:]...)
}

func (enc *Encoder) WriteLen(n int) {
	enc.WriteI32(int32(n))
}

func (enc *Encoder) WriteBytes(vs []byte) {
	enc.buf = append(enc.buf, vs...)
}

func (enc *Encoder) WriteStr(str string) {
	enc.WriteLen(len(str))
	enc.WriteBytes([]byte(str))
}

func (enc *Encoder) WriteReserved(n int) {
	enc.buf = append(enc.buf, make([]byte, n)...)
}

// Decoder decodes values from a buffer according to the XRootD protocol.
type Decoder struct {
	buf []byte
	pos int
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{buf: data}
}

func (dec *Decoder) ReadU8() uint8 {
	o := dec.buf[dec.pos]
	dec.pos++
	return o
}

func (dec *Decoder) ReadU16() uint16 {
	beg := dec.pos
	end := dec.pos + 2
	dec.pos += 2
	o := binary.BigEndian.Uint16(dec.buf[beg:end])
	return o
}

func (dec *Decoder) ReadI32() int32 {
	beg := dec.pos
	end := dec.pos + 4
	dec.pos += 4
	o := binary.BigEndian.Uint32(dec.buf[beg:end])
	return int32(o)
}

func (dec *Decoder) ReadI64() int64 {
	beg := dec.pos
	end := dec.pos + 8
	dec.pos += 8
	o := binary.BigEndian.Uint64(dec.buf[beg:end])
	return int64(o)
}

func (dec *Decoder) ReadLen() int {
	return int(dec.ReadI32())
}

func (dec *Decoder) ReadBytes(data []byte) {
	n := len(data)
	beg := dec.pos
	end := dec.pos + n
	copy(data, dec.buf[beg:end])
	dec.pos += n
	return
}

func (dec *Decoder) ReadStr() string {
	n := dec.ReadLen()
	beg := dec.pos
	end := dec.pos + n
	dec.pos += n
	return string(dec.buf[beg:end])
}

func (dec *Decoder) Skip(n int) {
	dec.pos += n
}
