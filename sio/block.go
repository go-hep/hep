// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"reflect"
)

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

func (blk *genericBlock) MarshalSio(w Writer) error {
	var err error
	err = bwrite(w, blk.rv.Interface())
	return err
}

func (blk *genericBlock) UnmarshalSio(r Reader) error {
	var err error
	err = bread(r, blk.rv.Interface())
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

func (blk *userBlock) MarshalSio(w Writer) error {
	return blk.blk.MarshalSio(w)
}

func (blk *userBlock) UnmarshalSio(r Reader) error {
	return blk.blk.UnmarshalSio(r)
}
