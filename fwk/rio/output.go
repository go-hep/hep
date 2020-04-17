// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"fmt"
	"io"
	"os"
	"reflect"

	"go-hep.org/x/hep/fwk"
	"go-hep.org/x/hep/rio"
)

// OutputStreamer writes data to a rio-stream.
type OutputStreamer struct {
	Name  string         // output filename
	w     io.WriteCloser // underlying output file
	rio   *rio.Writer    // output rio-stream
	recs  []*rio.Record  // list of connected records to write out
	ports []fwk.Port
}

func (o *OutputStreamer) Connect(ports []fwk.Port) error {
	var err error

	o.ports = make([]fwk.Port, len(ports))
	copy(o.ports, ports)

	// FIXME(sbinet): handle local/remote files, protocols
	o.w, err = os.Create(o.Name)
	if err != nil {
		return err
	}

	o.rio, err = rio.NewWriter(o.w)
	if err != nil {
		return err
	}

	for _, port := range o.ports {
		rec := o.rio.Record(port.Name)
		err = rec.Connect(port.Name, reflect.New(port.Type))
		if err != nil {
			return err
		}
		o.recs = append(o.recs, rec)
	}

	return err
}

func (o *OutputStreamer) Disconnect() error {
	// make sure we don't leak filedescriptors
	defer o.w.Close()

	err := o.rio.Close()
	if err != nil {
		return err
	}

	err = o.w.Close()
	if err != nil {
		return err
	}

	return err
}

func (o *OutputStreamer) Write(ctx fwk.Context) error {
	var err error
	store := ctx.Store()

	for i, rec := range o.recs {
		port := o.ports[i]

		n := rec.Name()
		blk := rec.Block(n)
		obj, err := store.Get(n)
		if err != nil {
			return err
		}

		rt := reflect.TypeOf(obj)
		if rt != port.Type {
			return fmt.Errorf("record[%s]: got type=%q. want type=%q.",
				rec.Name(),
				rt.Name(),
				port.Type,
			)
		}

		err = blk.Write(obj)
		if err != nil {
			return err
		}

		err = rec.Write()
		if err != nil {
			return err
		}
	}
	return err
}
