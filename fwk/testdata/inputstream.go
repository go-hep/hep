// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"fmt"
	"io"

	"go-hep.org/x/hep/fwk"
)

type InputStream struct {
	output string
	R      io.Reader
}

func (stream *InputStream) Connect(ports []fwk.Port) error {
	var err error
	stream.output = ports[0].Name

	return err
}

func (stream *InputStream) Read(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	var data int64
	_, err = fmt.Fscanf(stream.R, "%d\n", &data)
	if err != nil {
		return err
	}

	err = store.Put(stream.output, data)
	if err != nil {
		return err
	}

	return err
}

func (stream *InputStream) Disconnect() error {
	var err error
	return err
}
