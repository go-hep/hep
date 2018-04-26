// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat // import "go-hep.org/x/hep/xrootd/requests/stat"

import (
	"bytes"
	"strconv"

	"github.com/pkg/errors"
)

const RequestID uint16 = 3017

type Request struct {
	Options    byte
	Reserved   [11]byte
	FileHandle [4]byte
	PathLength int32
	Path       []byte
}

type Flag int32

const (
	FlagFile     Flag = 0
	FlagXset     Flag = 1
	FlagIsDir    Flag = 2
	FlagOther    Flag = 4
	FlagOffline  Flag = 8
	FlagReadable Flag = 16
	FlagWritable Flag = 32
	FlagPoscpend Flag = 64
	FlagBkpexist Flag = 128
)

type Response struct {
	ID               int64
	Size             int64
	Flags            Flag
	ModificationTime int64
}

func NewRequest(path string) Request {
	var pathBytes = make([]byte, len(path))
	copy(pathBytes, path)

	return Request{0, [11]byte{}, [4]byte{}, int32(len(path)), pathBytes}
}

func ParseReponsee(data []byte) (*Response, error) {
	dataParts := bytes.Split(data, []byte(" "))
	if len(dataParts) != 4 {
		return nil, errors.Errorf("Not enough fields in stat response: %s", data)
	}
	id, err := strconv.ParseInt(string(dataParts[0]), 10, 64)
	if err != nil {
		return nil, err
	}

	size, err := strconv.ParseInt(string(dataParts[1]), 10, 64)
	if err != nil {
		return nil, err
	}

	flags64, err := strconv.ParseInt(string(dataParts[2]), 10, 32)
	if err != nil {
		return nil, err
	}
	flags := Flag(flags64)

	modificationTime, err := strconv.ParseInt(string(dataParts[3][:len(dataParts[3])-1]), 10, 64)
	if err != nil {
		return nil, err
	}

	return &Response{id, size, flags, modificationTime}, nil
}
