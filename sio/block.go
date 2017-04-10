// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"bytes"
	"reflect"
)

// Marshaler is the interface implemented by an object that can marshal
// itself into a binary, sio-compatible, form.
type Marshaler interface {
	MarshalSio(buf *bytes.Buffer) error
}

// Unmarshaler is the interface implemented by an object that can
// unmarshal a binary, sio-compatible, representation of itself.
type Unmarshaler interface {
	UnmarshalSio(buf *bytes.Buffer) error
}

// Code is the interface implemented by an object that can
// unmarshal and marshal itself from and to a binary, sio-compatible, form.
type Codec interface {
	Marshaler
	Unmarshaler
}

// Block is the interface implemented by an object that can be
// stored to (and loaded from) an SIO stream.
type Block interface {
	Codec

	Name() string
	Version() uint32
}

// blockHeader describes the on-disk block data (header part)
type blockHeader struct {
	Len uint32 // length of this block
	Typ uint32 // block marker
}

// blockData describes the on-disk block data (payload part)
type blockData struct {
	Version uint32 // version of this block
	NameLen uint32 // length of the block name
}

// genericBlock provides a generic, reflect-based Block implementation
type genericBlock struct {
	rv      reflect.Value
	rt      reflect.Type
	version uint32
	name    string
}

func (blk *genericBlock) Name() string {
	return blk.name
}

func (blk *genericBlock) Version() uint32 {
	return blk.version
}

func (blk *genericBlock) MarshalSio(buf *bytes.Buffer) error {
	var err error
	err = bwrite(buf, blk.rv.Interface())
	return err
}

func (blk *genericBlock) UnmarshalSio(buf *bytes.Buffer) error {
	var err error
	err = bread(buf, blk.rv.Interface())
	return err
}

// userBlock adapts a user-provided Codec implementation into a Block one.
type userBlock struct {
	version uint32
	name    string
	blk     Codec
}

func (blk *userBlock) Name() string {
	return blk.name
}

func (blk *userBlock) Version() uint32 {
	return blk.version
}

func (blk *userBlock) MarshalSio(buf *bytes.Buffer) error {
	return blk.blk.MarshalSio(buf)
}

func (blk *userBlock) UnmarshalSio(buf *bytes.Buffer) error {
	return blk.blk.UnmarshalSio(buf)
}
