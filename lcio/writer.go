// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"compress/flate"

	"go-hep.org/x/hep/sio"
)

func Create(fname string) (*Writer, error) {
	f, err := sio.Create(fname)
	if err != nil {
		return nil, err
	}
	w := &Writer{
		f:    f,
		clvl: flate.DefaultCompression,
	}
	w.err = w.init()
	return w, w.err
}

type Writer struct {
	f    *sio.Stream
	clvl int
	recs struct {
		idx  *sio.Record
		rnd  *sio.Record
		rhdr *sio.Record
		ehdr *sio.Record
		evt  *sio.Record
	}
	data struct {
		idx  Index
		rnd  RandomAccess
		rhdr RunHeader
		ehdr EventHeader
	}
	closed bool
	err    error
}

func (w *Writer) Close() error {
	if w.closed {
		return w.err
	}
	w.err = w.f.Sync()
	if w.err != nil {
		return w.err
	}
	w.err = w.f.Close()
	w.closed = true
	return w.err
}

func (w *Writer) init() error {
	var (
		err error
		rec *sio.Record
	)

	compress := w.clvl != flate.NoCompression
	if compress {
		w.f.SetCompressionLevel(w.clvl)
	}

	rec = w.f.Record(Records.Index)
	if rec != nil {
		rec.SetUnpack(true)
		rec.SetCompress(compress)
		err = rec.Connect(Blocks.Index, &w.data.idx)
		if err != nil {
			return err
		}
	}
	w.recs.idx = rec

	rec = w.f.Record(Records.RandomAccess)
	if rec != nil {
		rec.SetUnpack(true)
		rec.SetCompress(compress)
		err = rec.Connect(Blocks.RandomAccess, &w.data.rnd)
		if err != nil {
			return err
		}
	}
	w.recs.rnd = rec

	rec = w.f.Record(Records.RunHeader)
	if rec != nil {
		rec.SetUnpack(true)
		rec.SetCompress(compress)
		err = rec.Connect(Blocks.RunHeader, &w.data.rhdr)
		if err != nil {
			return err
		}
	}
	w.recs.rhdr = rec

	rec = w.f.Record(Records.EventHeader)
	if rec != nil {
		rec.SetUnpack(true)
		rec.SetCompress(compress)
		err = rec.Connect(Blocks.EventHeader, &w.data.ehdr)
		if err != nil {
			return err
		}
	}
	w.recs.ehdr = rec

	rec = w.f.Record(Records.Event)
	if rec != nil {
		rec.SetUnpack(true)
	}
	w.recs.evt = rec

	return nil
}

func (w *Writer) SetCompressionLevel(lvl int) {
	w.clvl = lvl
	w.f.SetCompressionLevel(lvl)
}

func (w *Writer) WriteRunHeader(run *RunHeader) error {
	if w.err != nil {
		return w.err
	}
	w.data.rhdr = *run
	w.data.rhdr.SubDetectors = make([]string, len(run.SubDetectors))
	copy(w.data.rhdr.SubDetectors, run.SubDetectors)

	w.err = w.f.WriteRecord(w.recs.rhdr)
	return w.err
}

func (w *Writer) WriteEvent(evt *Event) error {
	if w.err != nil {
		return w.err
	}
	w.data.ehdr.RunNumber = evt.RunNumber
	w.data.ehdr.EventNumber = evt.EventNumber
	w.data.ehdr.TimeStamp = evt.TimeStamp
	w.data.ehdr.Detector = evt.Detector
	w.data.ehdr.Blocks = make([]Block, len(evt.names))
	for i, n := range evt.names {
		w.data.ehdr.Blocks[i] = Block{Name: n, Type: typeName(evt.colls[n])}
	}
	w.data.ehdr.Params = evt.Params

	w.err = w.f.WriteRecord(w.recs.ehdr)
	if w.err != nil {
		return w.err
	}

	w.recs.evt.Disconnect()

	for _, name := range evt.names {
		coll := evt.colls[name]
		w.err = w.recs.evt.Connect(name, coll)
		if w.err != nil {
			return w.err
		}
	}

	w.err = w.f.WriteRecord(w.recs.evt)
	if w.err != nil {
		return w.err
	}

	return w.err
}
