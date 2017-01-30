// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
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

func (b *Basket) Class() string {
	return "TBasket"
}

func (b *Basket) UnmarshalROOT(r *RBuffer) error {
	if err := b.Key.UnmarshalROOT(r); err != nil {
		return err
	}

	if b.Class() != "TBasket" {
		return fmt.Errorf("rootio.Basket: Key is not a Basket")
	}

	b.Version = r.ReadU16()
	b.Buffersize = r.ReadI32()
	b.Evbuffersize = r.ReadI32()

	if b.Evbuffersize < 0 {
		err := fmt.Errorf("rootio.Basket: incorrect Evbuffersize [%v]", b.Evbuffersize)
		b.Evbuffersize = 0
		return err
	}

	b.Nevbuf = r.ReadI32()
	b.Last = r.ReadI32()
	b.Flag = r.ReadU8()
	if b.Last > b.Buffersize {
		b.Buffersize = b.Last
	}

	return r.Err()
}

func init() {
	f := func() reflect.Value {
		o := &Basket{}
		return reflect.ValueOf(o)
	}
	Factory.add("TBasket", f)
	Factory.add("*rootio.Basket", f)
}

var _ Object = (*Basket)(nil)
var _ Named = (*Basket)(nil)
var _ ROOTUnmarshaler = (*Basket)(nil)
