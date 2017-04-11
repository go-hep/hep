// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"fmt"

	"go-hep.org/x/hep/sio"
)

func Open(fname string) (*Reader, error) {
	f, err := sio.Open(fname)
	if err != nil {
		return nil, err
	}
	r := &Reader{
		f: f,
	}
	r.init()
	return r, nil
}

type Reader struct {
	f    *sio.Stream
	idx  IOIndex
	rnd  RandomAccess
	rhdr RunHeader
	ehdr EventHeader
	evt  Event
	err  error
}

func (r *Reader) Close() error {
	return r.f.Close()
}

func (r *Reader) init() error {
	var (
		err error
		rec *sio.Record
	)

	//r.f.Record("LCIOIndex", &r.idx) // FIXME(sbinet)
	rec = r.f.Record("LCIORandomAccess")
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect("RandomAccess", &r.rnd)
		if err != nil {
			return err
		}
	}

	rec = r.f.Record("LCRunHeader")
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect("RunHeader", &r.rhdr)
		if err != nil {
			return err
		}
	}

	rec = r.f.Record("LCEventHeader")
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect("EventHeader", &r.ehdr)
		if err != nil {
			return err
		}
	}

	rec = r.f.Record("LCEvent")
	if rec != nil {
		rec.SetUnpack(true)
	}

	return nil
}

func (r *Reader) Next() bool {
	if r.err != nil {
		return false
	}
	rec, err := r.f.ReadRecord()
	if err != nil {
		r.err = err
		return false
	}
	// log.Printf(">>> %q", rec.Name())
	switch rec.Name() {
	case "LCRunHeader":
		if !rec.Unpack() {
			r.err = fmt.Errorf("lcio: expected record %q to unpack", rec.Name())
			return false
		}
		return r.Next()

	case "LCEventHeader":
		if !rec.Unpack() {
			r.err = fmt.Errorf("lcio: expected record %q to unpack", rec.Name())
			return false
		}
		err = r.remap()
		if err != nil {
			r.err = err
			return false
		}
		return r.Next()

	case "LCEvent":
		if !rec.Unpack() {
			r.err = fmt.Errorf("lcio: expected record %q to unpack", rec.Name())
		}
	}
	return true
}

func (r *Reader) remap() error {
	var err error
	rec := r.f.Record("LCEvent")
	r.evt.Collections = make(map[string]interface{}, len(r.ehdr.Blocks))
	r.evt.Names = make([]string, 0, len(r.ehdr.Blocks))
	for _, blk := range r.ehdr.Blocks {
		ptr := typeFrom(blk.Type)
		if ptr == nil {
			continue
		}
		err = rec.Connect(blk.Name, ptr)
		if err != nil {
			r.err = err
			return err
		}
		r.evt.Collections[blk.Name] = ptr
		r.evt.Names = append(r.evt.Names, blk.Name)
	}
	r.evt.RunNumber = r.ehdr.RunNumber
	r.evt.EventNumber = r.ehdr.EventNumber
	r.evt.TimeStamp = r.ehdr.TimeStamp
	r.evt.Detector = r.ehdr.Detector
	r.evt.Params = r.ehdr.Params
	return err
}

func (r *Reader) Event() Event {
	return r.evt
}

func (r *Reader) Err() error {
	return r.err
}
