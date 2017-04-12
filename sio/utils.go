// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio

import (
	"encoding/binary"
	"reflect"
)

// align4U32 returns sz adjusted to align at 4-byte boundaries
func align4U32(sz uint32) uint32 {
	return sz + (4-(sz&alignLen))&alignLen
}

// align4I32 returns sz adjusted to align at 4-byte boundaries
func align4I32(sz int32) int32 {
	return sz + (4-(sz&int32(alignLen)))&int32(alignLen)
}

// align4I64 returns sz adjusted to align at 4-byte boundaries
func align4I64(sz int64) int64 {
	return sz + (4-(sz&int64(alignLen)))&int64(alignLen)
}

// Unmarshal unmarshals a stream of bytes into ptr.
// If ptr implements Codec, use it.
func Unmarshal(r Reader, ptr interface{}) error {
	if ptr, ok := ptr.(Unmarshaler); ok {
		return ptr.UnmarshalSio(r)
	}
	return bread(r, ptr)
}

func bread(r Reader, data interface{}) error {
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
			err := Unmarshal(r, rv.Field(i).Addr().Interface())
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
		for i := 0; i < int(sz); i++ {
			err = Unmarshal(r, slice.Index(i).Addr().Interface())
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
			err = Unmarshal(r, kv.Interface())
			if err != nil {
				return err
			}
			vv := reflect.New(vt)
			err = Unmarshal(r, vv.Interface())
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
	panic("not reachable")
}

// Marshal marshals ptr to a stream of bytes.
// If ptr implements Codec, use it.
func Marshal(w Writer, ptr interface{}) error {
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
			err := Marshal(w, rv.Field(i).Addr().Interface())
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
			err = Marshal(w, rv.Index(i).Addr().Interface())
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
			err = Marshal(w, kv.Interface())
			if err != nil {
				return err
			}
			err = Marshal(w, vv.Interface())
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
	panic("not reachable")

}
