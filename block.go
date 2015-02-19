// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bytes"
)

// Block manages and desribes a block of data
type Block struct {
	raw rioBlock
}

func newBlock(name string, version Version) Block {
	block := Block{
		raw: rioBlock{
			Header: rioHeader{
				Len:   0,
				Frame: blkFrame,
			},
			Version: version,
			Name:    name,
		},
	}

	return block
}

// Name returns the name of this block
func (blk *Block) Name() string {
	return blk.raw.Name
}

// RioVersion returns the rio-binary version of the block
func (blk *Block) RioVersion() Version {
	return blk.raw.Version
}

// Write writes data to the Writer, in the rio format
func (blk *Block) Write(data interface{}) error {
	var err error

	buf := new(bytes.Buffer) // FIXME(sbinet): use a sync.Pool
	enc := encoder{w: buf}
	err = enc.Encode(data)
	if err != nil {
		return err
	}

	blk.raw.Data = buf.Bytes()
	blk.raw.Header.Len = uint32(len(blk.raw.Data))
	return nil
}

// Read reads data from the Reader, in the rio format
func (blk *Block) Read(data interface{}) error {
	var err error
	buf := bytes.NewReader(blk.raw.Data) // FIXME(sbinet): use a sync.Pool
	dec := decoder{r: buf}
	err = dec.Decode(data)
	if err != nil {
		return err
	}
	return err
}
