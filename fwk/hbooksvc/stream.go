// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbooksvc

import (
	"io"
	"os"

	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/rio"
)

// Mode describes the open-mode of a stream
type Mode int

const (
	Read  Mode = Mode(os.O_RDONLY)
	Write Mode = Mode(os.O_WRONLY)
)

// Stream defines an input or output hbook stream
type Stream struct {
	Name string // input|output file name
	Mode Mode   // read|write
}

type istream struct {
	name  string // stream name
	fname string // file name
	f     io.ReadCloser
	r     *rio.Reader
	objs  []fwk.Hist
}

func (stream *istream) close() error {
	defer stream.f.Close() // do not leak file descriptors
	err := stream.r.Close()
	if err != nil {
		return err
	}

	err = stream.f.Close()
	if err != nil {
		return err
	}

	return err
}

func (stream *istream) read(name string, ptr interface{}) error {
	var err error

	seekr, ok := stream.f.(io.Seeker)
	if !ok {
		return fwk.Errorf("hbooksvc: input stream [%s] is not seek-able", stream.name)
	}

	pos, err := seekr.Seek(0, 1)
	if err != nil {
		return err
	}
	defer seekr.Seek(pos, 0)

	_, err = seekr.Seek(0, 0)
	if err != nil {
		return err
	}

	r := seekr.(io.Reader)
	rr, err := rio.NewReader(r)
	if err != nil {
		return err
	}
	defer rr.Close()

	scan := rio.NewScanner(rr)
	scan.Select([]rio.Selector{{Name: name, Unpack: true}})
	if !scan.Scan() {
		return scan.Err()
	}
	rec := scan.Record()
	if rec == nil {
		return fwk.Errorf("hbooksvc: could not find record [%s] in stream [%s]", name, stream.name)
	}
	blk := rec.Block(name)
	if blk == nil {
		return fwk.Errorf(
			"hbooksvc: could not get block [%s] from record [%s] in stream [%s]",
			name, name, stream.name,
		)
	}
	err = blk.Read(ptr)
	if err != nil {
		return fwk.Errorf(
			"hbooksvc: could not read data from block [%s] from record [%s] in stream [%s]: %v",
			name, name, stream.name, err,
		)
	}
	return err
}

type ostream struct {
	name  string // stream name
	fname string // file name
	f     io.WriteCloser
	w     *rio.Writer
	objs  []fwk.Hist
}

func (stream *ostream) write() error {
	for i := range stream.objs {
		obj := stream.objs[i]
		name := string(obj.Name())
		rec := stream.w.Record(name)
		err := rec.Connect(name, obj.Value())
		if err != nil {
			return fwk.Errorf(
				"error writing object [%s] to stream [%s]: %v",
				name, stream.name, err,
			)
		}

		blk := rec.Block(name)
		err = blk.Write(obj.Value())
		if err != nil {
			return fwk.Errorf(
				"error writing object [%s] to stream [%s]: %v",
				name, stream.name, err,
			)
		}

		err = rec.Write()
		if err != nil {
			return fwk.Errorf(
				"error writing object [%s] to stream [%s]: %v",
				name, stream.name, err,
			)
		}
	}

	return nil
}

func (stream *ostream) close() error {
	defer stream.f.Close() // do not leak file descriptors
	err := stream.w.Close()
	if err != nil {
		return err
	}

	err = stream.f.Close()
	if err != nil {
		return err
	}

	return err
}
