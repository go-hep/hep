// Copyright ©2017 The go-hep Authors. All rights reserved.
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
	err = r.init()
	if err != nil {
		return nil, err
	}
	return r, nil
}

type Reader struct {
	f    *sio.Stream
	idx  Index
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

	rec = r.f.Record(Records.Index)
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect(Blocks.Index, &r.idx)
		if err != nil {
			return err
		}
	}

	rec = r.f.Record(Records.RandomAccess)
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect(Blocks.RandomAccess, &r.rnd)
		if err != nil {
			return err
		}
	}

	rec = r.f.Record(Records.RunHeader)
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect(Blocks.RunHeader, &r.rhdr)
		if err != nil {
			return err
		}
	}

	rec = r.f.Record(Records.EventHeader)
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect(Blocks.EventHeader, &r.ehdr)
		if err != nil {
			return err
		}
	}

	rec = r.f.Record(Records.Event)
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
	case Records.Index:
		if !rec.Unpack() {
			r.err = fmt.Errorf("lcio: expected record %q to unpack", rec.Name())
			return false
		}
		return r.Next()

	case Records.RandomAccess:
		if !rec.Unpack() {
			r.err = fmt.Errorf("lcio: expected record %q to unpack", rec.Name())
			return false
		}
		return r.Next()

	case Records.RunHeader:
		if !rec.Unpack() {
			r.err = fmt.Errorf("lcio: expected record %q to unpack", rec.Name())
			return false
		}
		return r.Next()

	case Records.EventHeader:
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

	case Records.Event:
		if !rec.Unpack() {
			r.err = fmt.Errorf("lcio: expected record %q to unpack", rec.Name())
		}
	}
	return true
}

func (r *Reader) remap() error {
	var err error
	rec := r.f.Record(Records.Event)
	r.evt.colls = nil
	r.evt.names = nil
	if len(r.ehdr.Blocks) > 0 {
		r.evt.colls = make(map[string]interface{}, len(r.ehdr.Blocks))
		r.evt.names = make([]string, 0, len(r.ehdr.Blocks))
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
			r.evt.colls[blk.Name] = ptr
			r.evt.names = append(r.evt.names, blk.Name)
		}
	}
	r.evt.RunNumber = r.ehdr.RunNumber
	r.evt.EventNumber = r.ehdr.EventNumber
	r.evt.TimeStamp = r.ehdr.TimeStamp
	r.evt.Detector = r.ehdr.Detector
	r.evt.Params = r.ehdr.Params
	return err
}

func (r *Reader) RunHeader() RunHeader {
	return r.rhdr
}

func (r *Reader) EventHeader() EventHeader {
	return r.ehdr
}

func (r *Reader) Event() Event {
	return r.evt
}

func (r *Reader) Err() error {
	return r.err
}
