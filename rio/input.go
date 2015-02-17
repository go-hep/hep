// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"io"
	"os"
	"reflect"

	"github.com/go-hep/fwk"
	"github.com/go-hep/rio"
)

// InputStream reads data from a (set of) rio-stream(s)
type InputStream struct {
	Names []string      // input filenames
	r     io.ReadCloser // underlying input file(s)
	rio   *rio.Reader   // input rio-stream
	scan  *rio.Scanner  // input records-scanner
	recs  []*rio.Record // list of connected records to read in
	ports []fwk.Port
}

func (input *InputStream) Connect(ports []fwk.Port) error {
	var err error

	input.ports = make([]fwk.Port, len(ports))
	copy(input.ports, ports)

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
	for _, port := range input.ports {
		rec := input.rio.Record(port.Name)
		err = rec.Connect(port.Name, reflect.New(port.Type))
		if err != nil {
			return err
		}
		input.recs = append(input.recs, rec)
		recnames = append(recnames, rio.Selector{Name: port.Name, Unpack: true})
	}

	input.scan = rio.NewScanner(input.rio)
	input.scan.Select(recnames)
	return err
}

func (input *InputStream) Read(ctx fwk.Context) error {
	var err error
	return err
}

func (input *InputStream) Disconnect() error {
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
