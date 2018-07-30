// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lcio

import (
	"compress/flate"

	"go-hep.org/x/hep/sio"
)

// Create creates a new LCIO writer, saving data in file fname.
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

// Writer provides a way to write LCIO RunHeaders and Events to an
// output SIO stream.
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

// Close closes the underlying output stream and makes it unavailable for
// further I/O operations.
// Close will synchronize and commit to disk any lingering data before closing
// the output stream.
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

	rec = w.f.Record(Records.Index)
	if rec != nil {
		err = rec.Connect(Blocks.Index, &w.data.idx)
		if err != nil {
			return err
		}
	}
	w.recs.idx = rec

	rec = w.f.Record(Records.RandomAccess)
	if rec != nil {
		err = rec.Connect(Blocks.RandomAccess, &w.data.rnd)
		if err != nil {
			return err
		}
	}
	w.recs.rnd = rec

	rec = w.f.Record(Records.RunHeader)
	if rec != nil {
		err = rec.Connect(Blocks.RunHeader, &w.data.rhdr)
		if err != nil {
			return err
		}
	}
	w.recs.rhdr = rec

	rec = w.f.Record(Records.EventHeader)
	if rec != nil {
		err = rec.Connect(Blocks.EventHeader, &w.data.ehdr)
		if err != nil {
			return err
		}
	}
	w.recs.ehdr = rec

	rec = w.f.Record(Records.Event)
	w.recs.evt = rec

	w.SetCompressionLevel(w.clvl)
	return nil
}

// SetCompressionLevel sets the compression level to lvl.
// lvl must be a compress/flate compression value.
// SetCompressionLevel must be called before WriteRunHeader or WriteEvent.
func (w *Writer) SetCompressionLevel(lvl int) {
	w.clvl = lvl
	compress := w.clvl != flate.NoCompression
	w.f.SetCompressionLevel(lvl)
	for _, rec := range []*sio.Record{
		w.recs.idx,
		w.recs.rnd,
		w.recs.rhdr,
		w.recs.ehdr,
		w.recs.evt,
	} {
		if rec == nil {
			continue
		}
		rec.SetCompress(compress)
	}
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
	w.data.ehdr.Blocks = make([]BlockDescr, len(evt.names))
	for i, n := range evt.names {
		w.data.ehdr.Blocks[i] = BlockDescr{Name: n, Type: typeName(evt.colls[n])}
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
