// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"encoding/binary"
	"reflect"
)

// Decoder decodes SIO streams.
// Decoder provides a nice API to deal with errors that may occur during decoding.
type Decoder struct {
	r   Reader
	err error
}

// NewDecoder creates a new Decoder reading from the provided sio.Reader.
func NewDecoder(r Reader) *Decoder {
	return &Decoder{r: r}
}

// Decode reads the next value from the input sio stream and stores it in
// the data, an empty interface value wrapping a pointer to a concrete value.
func (dec *Decoder) Decode(ptr any) {
	if dec.err != nil {
		return
	}

	dec.err = unmarshal(dec.r, ptr)
}

// Tag tags a pointer, assigning it a unique identifier, so links between values
// (inside a given SIO record) can be rebuilt.
func (dec *Decoder) Tag(ptr any) {
	if dec.err != nil {
		return
	}
	dec.err = dec.r.Tag(ptr)
}

// Pointer marks a (pointer to a) pointer, assigning it a unique identifier,
// so links between values (inside a given SIO record) can be rebuilt.
func (dec *Decoder) Pointer(ptr any) {
	if dec.err != nil {
		return
	}
	dec.err = dec.r.Pointer(ptr)
}

// Err returns the first encountered error while decoding, if any.
func (dec *Decoder) Err() error {
	return dec.err
}

// unmarshal unmarshals a stream of bytes into ptr.
// If ptr implements Codec, use it.
func unmarshal(r Reader, ptr any) error {
	if ptr, ok := ptr.(Unmarshaler); ok {
		return ptr.UnmarshalSio(r)
	}
	return bread(r, ptr)
}

func bread(r Reader, data any) error {
	bo := binary.BigEndian
	rv := reflect.ValueOf(data)
	//fmt.Printf("::: [%v] :::...\n", rv.Type())
	//defer fmt.Printf("### [%v] [done]\n", rv.Type())

	switch rv.Type().Kind() {
	case reflect.Ptr:
		rv = rv.Elem()
	}
	switch rv.Type().Kind() {
	case reflect.Struct:
		//fmt.Printf(">>> struct: [%v]...\n", rv.Type())
		for i, n := 0, rv.NumField(); i < n; i++ {
			//fmt.Printf(">>> i=%d [%v] (%v)...\n", i, rv.Field(i).Type(), rv.Type().Name()+"."+rv.Type().Field(i).Name)
			err := unmarshal(r, rv.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			//fmt.Printf(">>> i=%d [%v] (%v)...[done=>%v]\n", i, rv.Field(i).Type(), rv.Type().Name(), rv.Field(i).Interface())
		}
		//fmt.Printf(">>> struct: [%v]...[done]\n", rv.Type())
		return nil
	case reflect.String:
		//fmt.Printf("++++> string...\n")
		sz := uint32(0)
		err := bread(r, &sz)
		if err != nil {
			return err
		}
		strlen := align4U32(sz)
		//fmt.Printf("string [%d=>%d]...\n", sz, strlen)
		str := make([]byte, strlen)
		_, err = r.Read(str)
		if err != nil {
			return err
		}
		rv.SetString(string(str[:sz]))
		//fmt.Printf("<++++ [%s]\n", string(str))
		return nil

	case reflect.Slice:
		//fmt.Printf("<<< slice: [%v|%v]...\n", rv.Type(), rv.Type().Elem().Kind())
		sz := uint32(0)
		err := bread(r, &sz)
		if err != nil {
			return err
		}
		//fmt.Printf("<<< slice: %d [%v]\n", sz, rv.Type())
		slice := reflect.MakeSlice(rv.Type(), int(sz), int(sz))
		for i := range int(sz) {
			err = unmarshal(r, slice.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		rv.Set(slice)
		//fmt.Printf("<<< slice: [%v]... [done] (%#v)\n", rv.Type(), rv.Interface())
		return err

	case reflect.Map:
		//fmt.Printf(">>> map: [%v]...\n", rv.Type())
		sz := uint32(0)
		err := bread(r, &sz)
		if err != nil {
			return err
		}
		//fmt.Printf(">>> map: %d [%v]\n", sz, rv.Type())
		m := reflect.MakeMap(rv.Type())
		kt := rv.Type().Key()
		vt := rv.Type().Elem()
		for i, n := 0, int(sz); i < n; i++ {
			kv := reflect.New(kt)
			err = unmarshal(r, kv.Interface())
			if err != nil {
				return err
			}
			vv := reflect.New(vt)
			err = unmarshal(r, vv.Interface())
			if err != nil {
				return err
			}
			//fmt.Printf("m - %d: {%v} - {%v}\n", i, kv.Elem().Interface(), vv.Elem().Interface())
			m.SetMapIndex(kv.Elem(), vv.Elem())
		}
		rv.Set(m)
		return nil

	default:
		//fmt.Printf(">>> binary - [%v]...\n", rv.Type())
		err := binary.Read(r, bo, data)
		//fmt.Printf(">>> binary - [%v]... [done]\n", rv.Type())
		return err
	}
}
