// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"bytes"
	"reflect"
)

type BinaryMarshaler interface {
	MarshalBinary(buf *bytes.Buffer) error
}

type BinaryUnmarshaler interface {
	UnmarshalBinary(buf *bytes.Buffer) error
}

type BinaryCodec interface {
	BinaryMarshaler
	BinaryUnmarshaler
}

type Block interface {
	BinaryCodec

	Name() string
	// Xfer(stream *Stream, op Operation, version int) error
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

func (blk *genericBlock) MarshalBinary(buf *bytes.Buffer) error {
	var err error
	err = bwrite(buf, blk.rv.Interface())
	return err
}

func (blk *genericBlock) UnmarshalBinary(buf *bytes.Buffer) error {
	var err error
	err = bread(buf, blk.rv.Interface())
	return err
}

type mBlockImpl struct {
	version uint32
	name    string
	blk     BinaryCodec
}

func (blk *mBlockImpl) Name() string {
	return blk.name
}

func (blk *mBlockImpl) Version() uint32 {
	return blk.version
}

func (blk *mBlockImpl) MarshalBinary(buf *bytes.Buffer) error {
	return blk.blk.MarshalBinary(buf)
}

func (blk *mBlockImpl) UnmarshalBinary(buf *bytes.Buffer) error {
	return blk.blk.UnmarshalBinary(buf)
}
