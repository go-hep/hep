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

	vers    uint16
	bufsize int // length in bytes
	nevsize int // length in int_t or fixed length of each entry
	nevbuf  int // number of entries in basket
	last    int // pointer to last used byte in basket
	flag    byte
}

func (b *Basket) Class() string {
	return "TBasket"
}

func (b *Basket) UnmarshalROOT(r *RBuffer) error {
	if r.err != nil {
		return r.err
	}

	if err := b.Key.UnmarshalROOT(r); err != nil {
		return err
	}

	if b.Class() != "TBasket" {
		return fmt.Errorf("rootio.Basket: Key is not a Basket")
	}

	b.vers = r.ReadU16()
	b.bufsize = int(r.ReadI32())
	b.nevsize = int(r.ReadI32())

	if b.nevsize < 0 {
		r.err = fmt.Errorf("rootio.Basket: incorrect event buffer size [%v]", b.nevsize)
		b.nevsize = 0
		return r.err
	}

	b.nevbuf = int(r.ReadI32())
	b.last = int(r.ReadI32())
	b.flag = r.ReadU8()

	if b.last > b.bufsize {
		b.bufsize = b.last
	}

	fmt.Printf("+++ TBasket: %q %q vers=%d bufsize=%d nevsize=%d nevbuf=%d last=%v flag=0x%x\n",
		b.Name(), b.Title(), b.vers, b.bufsize, b.nevsize, b.nevbuf, b.last, b.flag,
	)

	if b.flag == 0 {
		return r.Err()
	}

	panic("not implemented")

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
