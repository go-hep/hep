package rio

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

type Block interface {
	BinaryMarshaler
	BinaryUnmarshaler

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

type blockImpl struct {
	rv      reflect.Value
	rt      reflect.Type
	version uint32
	name    string
}

func (blk *blockImpl) Name() string {
	return blk.name
}

func (blk *blockImpl) Version() uint32 {
	return blk.version
}

func (blk *blockImpl) MarshalBinary(buf *bytes.Buffer) error {
	var err error
	err = bwrite(buf, blk.version)
	if err != nil {
		return err
	}
	err = bwrite(buf, blk.name)
	if err != nil {
		return err
	}
	err = bwrite(buf, blk.rv.Interface())
	return err
}

func (blk *blockImpl) UnmarshalBinary(buf *bytes.Buffer) error {
	var err error
	err = bread(buf, blk.rv.Interface())
	return err
}
