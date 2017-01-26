// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"bytes"
	"fmt"
	"reflect"
)

type Basket struct {
	*Key

	Version      uint16
	Buffersize   int32 // length in bytes
	Evbuffersize int32 // length in int_t or fixed length of each entry
	Nevbuf       int32 // number of entries in basket
	Last         int32 // pointer to last used byte in basket
	Flag         byte
}

func (b *Basket) UnmarshalROOT(data *bytes.Buffer) error {

	//fmt.Printf("+++++ Basket.UnmarshalROOT ++++\n")
	err := b.Key.UnmarshalROOT(data)
	if err != nil {
		return err
	}

	if b.Class() != "TBasket" {
		return fmt.Errorf("rootio.Basket: Key is not a Basket")
	}

	dec := newDecoder(data)
	err = dec.readBin(&b.Version)
	if err != nil {
		return err
	}

	err = dec.readInt32(&b.Buffersize)
	if err != nil {
		return err
	}

	err = dec.readInt32(&b.Evbuffersize)
	if err != nil {
		return err
	}

	if b.Evbuffersize < 0 {
		err = fmt.Errorf("rootio.Basket: incorrect Evbuffersize [%v]", b.Evbuffersize)
		b.Evbuffersize = 0
		return err
	}

	err = dec.readInt32(&b.Nevbuf)
	if err != nil {
		return err
	}

	err = dec.readInt32(&b.Last)
	if err != nil {
		return err
	}

	err = dec.readBin(&b.Flag)
	if err != nil {
		return err
	}

	if b.Last > b.Buffersize {
		b.Buffersize = b.Last
	}

	return err
}

func init() {
	f := func() reflect.Value {
		o := &Basket{}
		return reflect.ValueOf(o)
	}
	Factory.db["TBasket"] = f
	Factory.db["*rootio.Basket"] = f
}

var _ Object = (*Key)(nil)
var _ ROOTUnmarshaler = (*Key)(nil)
