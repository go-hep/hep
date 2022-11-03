// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootcnv

import (
	"errors"
	"fmt"
	"io"

	"go-hep.org/x/hep/groot/rtree"
	"go-hep.org/x/hep/hepmc"
)

type rstream struct {
	evt hepmc.Event
	err error
}

type FlatTreeReader struct {
	r *rtree.Reader

	evt   event
	rvars []rtree.ReadVar

	evts chan rstream
}

func NewFlatTreeReader(t rtree.Tree, opts ...rtree.ReadOption) (*FlatTreeReader, error) {
	r := FlatTreeReader{
		evts: make(chan rstream),
	}

	r.rvars = rtree.ReadVarsFromStruct(&r.evt)

	rr, err := rtree.NewReader(t, r.rvars, opts...)
	if err != nil {
		return nil, fmt.Errorf("hepmc: could not create ROOT Tree reader: %w", err)
	}
	r.r = rr

	go func() {
		defer close(r.evts)
		err := rr.Read(func(ctx rtree.RCtx) error {
			defer r.evt.reset()
			var o rstream
			err := r.evt.write(&o.evt)
			if err != nil {
				return err
			}
			r.evts <- o
			return nil
		})
		if err != nil {
			r.evts <- rstream{err: err}
			return
		}
		r.evts <- rstream{err: io.EOF}
	}()

	return &r, nil
}

func (r *FlatTreeReader) Close() error {
	return r.r.Close()
}

func (r *FlatTreeReader) Read(evt *hepmc.Event) error {
	o := <-r.evts
	if o.err != nil {
		if errors.Is(o.err, io.EOF) {
			return io.EOF
		}
		return fmt.Errorf("hepmc: could not read event from ROOT: %w", o.err)
	}
	*evt = o.evt
	return nil
}
