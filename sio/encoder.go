// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"encoding/binary"
	"reflect"
)

// Encoder encodes values into a SIO stream.
// Encoder provides a nice API to deal with errors that may occur during encoding.
type Encoder struct {
	w   Writer
	err error
}

// NewEncoder creates a new Encoder writing to the provided sio.Writer.
func NewEncoder(w Writer) *Encoder {
	return &Encoder{w: w}
}

// Encode writes the next value to the output sio stream.
func (enc *Encoder) Encode(data interface{}) {
	if enc.err != nil {
		return
	}
	enc.err = marshal(enc.w, data)
}

// Tag tags a pointer, assigning it a unique identifier, so links between values
// (inside a given sio record) can be stored.
func (enc *Encoder) Tag(ptr interface{}) {
	if enc.err != nil {
		return
	}
	enc.err = enc.w.Tag(ptr)
}

// Pointer marks a (pointer to a) pointer, assigning it a unique identifier,
// so links between values (inside a given SIO record) can be stored.
func (enc *Encoder) Pointer(ptr interface{}) {
	if enc.err != nil {
		return
	}
	enc.err = enc.w.Pointer(ptr)
}

// Err returns the first encountered error while decoding, if any.
func (dec *Encoder) Err() error {
	return dec.err
}

// marshal marshals ptr to a stream of bytes.
// If ptr implements Codec, use it.
func marshal(w Writer, ptr interface{}) error {
	if ptr, ok := ptr.(Marshaler); ok {
		return ptr.MarshalSio(w)
	}
	return bwrite(w, ptr)
}

func bwrite(w Writer, data interface{}) error {

	bo := binary.BigEndian
	rrv := reflect.ValueOf(data)
	rv := reflect.Indirect(rrv)
	// fmt.Printf("::: [%v] :::...\n", rrv.Type())
	// defer fmt.Printf("### [%v] [done]\n", rrv.Type())

	switch rv.Type().Kind() {
	case reflect.Struct:
		//fmt.Printf(">>> struct: [%v]...\n", rv.Type())
		for i, n := 0, rv.NumField(); i < n; i++ {
			//fmt.Printf(">>> i=%d [%v] (%v)...\n", i, rv.Field(i).Type(), rv.Type().Name())
			err := marshal(w, rv.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			//fmt.Printf(">>> i=%d [%v] (%v)...[done]\n", i, rv.Field(i).Type(), rv.Type().Name())
		}
		//fmt.Printf(">>> struct: [%v]...[done]\n", rv.Type())
		return nil
	case reflect.String:
		str := rv.String()
		sz := uint32(len(str))
		// fmt.Printf("++++> (%d) [%s]\n", sz, string(str))
		err := bwrite(w, &sz)
		if err != nil {
			return err
		}
		bstr := []byte(str)
		bstr = append(bstr, make([]byte, align4U32(sz)-sz)...)
		_, err = w.Write(bstr)
		if err != nil {
			return err
		}
		// fmt.Printf("<++++ (%d) [%s]\n", sz, string(str))
		return nil

	case reflect.Slice:
		// fmt.Printf(">>> slice: [%v|%v]...\n", rv.Type(), rv.Type().Elem().Kind())
		sz := uint32(rv.Len())
		// fmt.Printf(">>> slice: %d [%v]\n", sz, rv.Type())
		err := bwrite(w, &sz)
		if err != nil {
			return err
		}
		for i := 0; i < int(sz); i++ {
			err = marshal(w, rv.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		// fmt.Printf(">>> slice: [%v]... [done] (%v)\n", rv.Type(), rv.Interface())
		return err

	case reflect.Map:
		//fmt.Printf(">>> map: [%v]...\n", rv.Type())
		sz := uint32(rv.Len())
		err := bwrite(w, &sz)
		if err != nil {
			return err
		}
		//fmt.Printf(">>> map: %d [%v]\n", sz, rv.Type())
		for _, kv := range rv.MapKeys() {
			vv := rv.MapIndex(kv)
			err = marshal(w, kv.Interface())
			if err != nil {
				return err
			}
			err = marshal(w, vv.Interface())
			if err != nil {
				return err
			}
			//fmt.Printf("m - %d: {%v} - {%v}\n", i, kv.Elem().Interface(), vv.Elem().Interface())
		}
		return nil

	default:
		//fmt.Printf(">>> binary - [%v]...\n", rv.Type())
		err := binary.Write(w, bo, data)
		//fmt.Printf(">>> binary - [%v]... [done]\n", rv.Type())
		return err
	}
}
