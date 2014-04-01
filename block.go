package rio

import (
	"bytes"
	"encoding"
	"reflect"
)

type Block interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler

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

func (blk *blockImpl) MarshalBinary() ([]byte, error) {
	var err error
	var buf bytes.Buffer

	err = bwrite(&buf, blk.rv.Interface())
	return buf.Bytes(), err
}

func (blk *blockImpl) UnmarshalBinary(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = bread(buf, blk.rv.Interface())
	data = data[len(data)-buf.Len():]
	return err
}
