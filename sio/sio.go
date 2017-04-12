// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
)

// Reader is the interface that wraps the basic io.Reader interface
// and adds SIO pointer tagging capabilities.
type Reader interface {
	io.Reader

	Version() uint32
	Tag(ptr interface{}) error
	Pointer(ptr interface{}) error
}

// Writer is the interface that wraps the basic io.Writer interface
// and adds SIO pointer tagging capabilities.
type Writer interface {
	io.Writer

	Version() uint32
	Tag(ptr interface{}) error
	Pointer(ptr interface{}) error
}

// Marshaler is the interface implemented by an object that can marshal
// itself into a binary, sio-compatible, form.
type Marshaler interface {
	MarshalSio(w Writer) error
}

// Unmarshaler is the interface implemented by an object that can
// unmarshal a binary, sio-compatible, representation of itself.
type Unmarshaler interface {
	UnmarshalSio(r Reader) error
}

// Code is the interface implemented by an object that can
// unmarshal and marshal itself from and to a binary, sio-compatible, form.
type Codec interface {
	Marshaler
	Unmarshaler
}

type reader struct {
	buf *bytes.Buffer
	ver uint32
	ptr map[uint32]interface{}
	tag map[interface{}]uint32
}

func newReader(data []byte) *reader {
	return &reader{
		buf: bytes.NewBuffer(data),
		ptr: make(map[uint32]interface{}),
		tag: make(map[interface{}]uint32),
	}
}

func (r *reader) Read(data []byte) (int, error) {
	return r.buf.Read(data)
}

func (r *reader) Bytes() []byte {
	return r.buf.Bytes()
}

func (r *reader) Len() int {
	return r.buf.Len()
}

func (r *reader) Next(n int) []byte {
	return r.buf.Next(n)
}

func (r *reader) Version() uint32 {
	min := r.ver & uint32(0x0000ffff)
	maj := (r.ver & uint32(0xffff0000)) >> 16
	return maj*1000 + min
}

func (r *reader) Tag(ptr interface{}) error {
	var id uint32
	err := binary.Read(r.buf, binary.BigEndian, &id)
	if err != nil {
		return err
	}
	if id == ptagMarker {
		return nil
	}
	r.tag[ptr] = id
	return nil
}

func (r *reader) Pointer(ptr interface{}) error {
	rptr := reflect.ValueOf(ptr)
	if rptr.Kind() != reflect.Ptr || rptr.Elem().Kind() != reflect.Ptr {
		panic(fmt.Errorf("sio: Reader.Pointer expects a pointer to pointer"))
	}

	var pid uint32
	err := binary.Read(r, binary.BigEndian, &pid)
	if err != nil {
		return err
	}
	if pid == pntrMarker {
		return nil
	}

	r.ptr[pid] = ptr
	return nil
}

func (r *reader) relocate() {
ptrloop:
	for pid, ptr := range r.ptr {
		rptr := reflect.ValueOf(ptr)
		for tag, tid := range r.tag {
			if tid == pid {
				rtag := reflect.ValueOf(tag)
				rptr.Elem().Set(rtag)
				continue ptrloop
			}
		}
	}
}

type writer struct {
	buf *bytes.Buffer
	ver uint32
	ids uint32
	ptr map[uint32]interface{}
	tag map[interface{}]uint32
}

func newWriter() *writer {
	return &writer{
		buf: new(bytes.Buffer),
		ptr: make(map[uint32]interface{}),
		tag: make(map[interface{}]uint32),
	}
}

func newWriterFrom(w *writer) *writer {
	return &writer{
		buf: new(bytes.Buffer),
		ver: w.ver,
		ids: w.ids,
		ptr: w.ptr,
		tag: w.tag,
	}
}

func (w *writer) Write(data []byte) (int, error) {
	return w.buf.Write(data)
}

func (w *writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *writer) Len() int {
	return w.buf.Len()
}

func (w *writer) Version() uint32 {
	min := w.ver & uint32(0x0000ffff)
	maj := (w.ver & uint32(0xffff0000)) >> 16
	return maj*1000 + min
}

func (w *writer) Tag(ptr interface{}) error {
	var id uint32 = ptagMarker
	if _, ok := w.tag[ptr]; !ok {
		err := w.genID()
		if err != nil {
			return err
		}
		w.tag[ptr] = w.ids
	}
	id = w.tag[ptr]
	err := binary.Write(w.buf, binary.BigEndian, &id)
	if err != nil {
		return err
	}
	return nil
}

func (w *writer) Pointer(ptr interface{}) error {
	ptr = reflect.ValueOf(ptr).Elem().Interface()
	var id uint32 = pntrMarker
	if _, ok := w.tag[ptr]; !ok {
		err := w.genID()
		if err != nil {
			return err
		}
		w.tag[ptr] = w.ids
	}
	id = w.tag[ptr]
	err := binary.Write(w.buf, binary.BigEndian, &id)
	if err != nil {
		return err
	}
	return nil
}

func (w *writer) genID() error {
	if w.ids+1 == math.MaxUint32 {
		return errPointerIDOverflow
	}
	w.ids++
	return nil
}

var _ Reader = (*reader)(nil)
var _ Writer = (*writer)(nil)
