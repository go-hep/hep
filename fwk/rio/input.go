// Copyright 2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"io"
	"os"
	"reflect"

	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/rio"
	"golang.org/x/xerrors"
)

// InputStreamer reads data from a (set of) rio-stream(s)
type InputStreamer struct {
	Names []string            // input filenames
	r     io.ReadCloser       // underlying input file(s)
	rio   *rio.Reader         // input rio-stream
	scan  *rio.Scanner        // input records-scanner
	ports map[string]fwk.Port // input ports to read/populate
}

func (input *InputStreamer) Connect(ports []fwk.Port) error {
	var err error

	input.ports = make(map[string]fwk.Port, len(ports))

	// FIXME(sbinet): handle multi-reader
	// FIXME(sbinet): handle local/remote files, protocols
	input.r, err = os.Open(input.Names[0])
	if err != nil {
		return err
	}

	input.rio, err = rio.NewReader(input.r)
	if err != nil {
		return err
	}

	recnames := make([]rio.Selector, 0, len(input.ports))
	for _, port := range ports {
		input.ports[port.Name] = port
		rec := input.rio.Record(port.Name)
		err = rec.Connect(port.Name, reflect.New(port.Type))
		if err != nil {
			return err
		}
		recnames = append(recnames, rio.Selector{Name: port.Name, Unpack: true})
	}

	input.scan = rio.NewScanner(input.rio)
	input.scan.Select(recnames)
	return err
}

func (input *InputStreamer) Read(ctx fwk.Context) error {
	store := ctx.Store()
	recs := make(map[string]struct{}, len(input.ports))
	for i := 0; i < len(input.ports); i++ {
		if !input.scan.Scan() {
			err := input.scan.Err()
			if err == nil {
				return io.EOF
			}
		}
		rec := input.scan.Record()
		blk := rec.Block(rec.Name())
		obj := reflect.New(input.ports[rec.Name()].Type).Elem()
		err := blk.Read(obj.Addr().Interface())
		if err != nil {
			return xerrors.Errorf("block-read error: %w", err)
		}
		err = store.Put(rec.Name(), obj.Interface())
		if err != nil {
			return xerrors.Errorf("store-put error: %w", err)
		}
		recs[rec.Name()] = struct{}{}
	}

	if len(recs) != len(input.ports) {
		return xerrors.Errorf("fwk.rio: expected inputs: %d. got: %d.", len(input.ports), len(recs))
	}

	return nil
}

func (input *InputStreamer) Disconnect() error {
	var err error
	// make sure we don't leak filedescriptors
	defer input.r.Close()

	err = input.rio.Close()
	if err != nil {
		return err
	}

	err = input.r.Close()
	if err != nil {
		return err
	}

	return err
}
