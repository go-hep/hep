// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import "go-hep.org/x/hep/sio"

func Create(fname string) (*Writer, error) {
	f, err := sio.Create(fname)
	if err != nil {
		return nil, err
	}
	w := &Writer{
		f: f,
	}
	w.init()
	return w, nil
}

type Writer struct {
	f    *sio.Stream
	idx  Index
	rnd  RandomAccess
	rhdr RunHeader
	ehdr EventHeader
	evt  Event
	err  error
}

func (w *Writer) Close() error {
	return w.f.Close()
}

func (w *Writer) init() error {
	var (
		err error
		rec *sio.Record
	)

	rec = w.f.Record("LCIOIndex")
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect("LCIOIndex", &w.idx)
		if err != nil {
			return err
		}
	}

	rec = w.f.Record("LCIORandomAccess")
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect("LCIORandomAccess", &w.rnd)
		if err != nil {
			return err
		}
	}

	rec = w.f.Record("LCRunHeader")
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect("RunHeader", &w.rhdr)
		if err != nil {
			return err
		}
	}

	rec = w.f.Record("LCEventHeader")
	if rec != nil {
		rec.SetUnpack(true)
		err = rec.Connect("EventHeader", &w.ehdr)
		if err != nil {
			return err
		}
	}

	rec = w.f.Record("LCEvent")
	if rec != nil {
		rec.SetUnpack(true)
	}

	return nil
}
