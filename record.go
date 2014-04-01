package rio

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

// recordHeader describes the on-disk record (header part)
type recordHeader struct {
	HdrLen  uint32
	BufType uint32
}

// recordData describes the on-disk record (payload part)
type recordData struct {
	Options uint32
	DataLen uint32 // length of compressed record data
	UCmpLen uint32 // length of uncompressed record data
	NameLen uint32 // length of record name
}

// Record manages blocks of data
type Record struct {
	name   string           // record name
	unpack bool             // whether to unpack incoming records
	blocks map[string]Block // connected blocks
}

// Name returns the name of this record
func (rec *Record) Name() string {
	return rec.name
}

// Unpack returns whether to unpack incoming records
func (rec *Record) Unpack() bool {
	return rec.unpack
}

// SetUnpack sets whether to unpack incoming records
func (rec *Record) SetUnpack(unpack bool) {
	rec.unpack = unpack
}

// Connect connects a Block to this Record (for reading or writing)
func (rec *Record) Connect(name string, ptr interface{}) error {
	var err error
	_, dup := rec.blocks[name]
	if dup {
		return fmt.Errorf("rio.Record: Block name [%s] already connected", name)
	}
	var block Block
	switch ptr.(type) {
	case Block:
		block = ptr.(Block)
	default:
		rt := reflect.TypeOf(ptr)
		if rt.Kind() != reflect.Ptr {
			return fmt.Errorf("rio: Connect needs a pointer to a block of data")
		}
		block = &blockImpl{
			rt:      rt,
			rv:      reflect.ValueOf(ptr),
			version: 0,
			name:    rt.Name(),
		}
	}
	rec.blocks[name] = block
	return err
}

// read reads a record
func (rec *Record) read(buf *bytes.Buffer) error {
	var err error
	//fmt.Printf("::: reading record [%s]... [%d]\n", rec.name, buf.Len())
	// loop until data has been depleted
	for buf.Len() > 0 {
		// read block header
		var hdr blockHeader
		err = bread(buf, &hdr)
		if err != nil {
			return err
		}
		if hdr.Typ != g_mark_block {
			fmt.Printf("*** err record[%s]: noblockmarker\n", rec.name)
			return ErrRecordNoBlockMarker
		}

		var data blockData
		err = bread(buf, &data)
		if err != nil {
			return err
		}
		var cbuf bytes.Buffer
		nlen := align4(data.NameLen)
		n, err := io.CopyN(&cbuf, buf, int64(nlen))
		if err != nil {
			fmt.Printf(">>> err:%v\n", err)
			return err
		}
		if n != int64(nlen) {
			return fmt.Errorf("rio: read too few bytes (got=%d. expected=%d)", n, nlen)
		}
		name := string(cbuf.Bytes()[:data.NameLen])
		blk, ok := rec.blocks[name]
		if !ok {
			fmt.Printf("*** no block [%s]. draining buffer!\n", name)
			// drain the whole buffer
			buf.Next(buf.Len())
			continue
		}
		//fmt.Printf("### %q\n", string(buf.Bytes()))
		err = blk.UnmarshalBinary(buf)
		if err != nil {
			fmt.Printf("*** error unmarshaling record=%q block=%q: %v\n", rec.name, name, err)
			return err
		}
		//fmt.Printf(">>> read record=%q block=%q (buf=%d)\n", rec.name, name, buf.Len())

		// check whether there is still something to be read.
		// if there is, check whether there is a block-marker
		if buf.Len() > 0 {
			rest := buf.Bytes()
			idx := bytes.Index(rest, g_mark_block_b)
			if idx > 0 {
				buf.Next(idx - 8 /*sizeof blockHeader*/)
			} else {
				buf.Next(buf.Len())
			}
		}
	}
	//fmt.Printf("::: reading record [%s]... [done]\n", rec.name)
	return err
}

// EOF
