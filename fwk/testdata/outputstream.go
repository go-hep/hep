// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"fmt"
	"io"

	"go-hep.org/x/hep/fwk"
)

type OutputStream struct {
	input string

	W io.Writer
}

func (out *OutputStream) Connect(ports []fwk.Port) error {
	var err error
	out.input = ports[0].Name
	return err
}

func (out *OutputStream) Write(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	v, err := store.Get(out.input)
	if err != nil {
		return err
	}

	data := v.(int64)
	_, err = out.W.Write([]byte(fmt.Sprintf("%d\n", data)))
	if err != nil {
		return err
	}

	return err
}

func (out *OutputStream) Disconnect() error {
	var err error

	return err
}
