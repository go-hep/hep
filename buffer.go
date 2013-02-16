package root

import (
	"bytes"
	B "encoding/binary"
	"reflect"
)

type Basket struct {
	Version      uint16
	Buffersize   int32 // length in bytes
	Evbuffersize int32 // length in int_t or fixed length of each entry
	Nevbuf       int32 // number of entries in basket
	Last         int32 // pointer to last used byte in basket
	Flag         byte
}

func (k *Key) DecodeVector(in *bytes.Buffer, dst interface{}) int {
	// Discard three int16s (like 40 00 00 0e 00 09)
	x := in.Next(6)
	_ = x // sometimes we want to look at this.

	var n int32
	err := B.Read(in, B.BigEndian, &n)
	if err != nil {
		panic(err)
	}

	err = B.Read(in, B.BigEndian, reflect.ValueOf(dst).Slice(0, int(n)).Interface())
	if err != nil {
		panic(err)
	}
	return int(n)
}
