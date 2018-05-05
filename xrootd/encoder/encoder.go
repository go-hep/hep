// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package encoder // import "go-hep.org/x/hep/xrootd/encoder"

import (
	"encoding/binary"
	"reflect"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd/protocol"
)

// MarshalRequest encodes request body alongside with request and stream ids
func MarshalRequest(requestID uint16, streamID protocol.StreamID, requestBody interface{}) ([]byte, error) {
	requestHeader := make([]byte, 4)

	copy(requestHeader, streamID[:])
	binary.BigEndian.PutUint16(requestHeader[2:], requestID)

	b, err := Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	return append(requestHeader, b...), nil
}

// Marshal encodes structure to the bytes following the XRootd protocol specification.
// Fields are encoded in the same order as they are specified in the struct definition.
// Each field is encoded in network byte order (BigEndian) without any alignment or padding.
// Slices and arrays of uint8 are binary copied without further encoding.
// Supported types are: uint8, uint16, int32, int64, slices and arrays of uint8.
func Marshal(x interface{}) ([]byte, error) {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.Errorf("Cannot marshal %s, struct is expected", v.Kind())
	}

	dataSize, err := calculateSizeForMarshaling(v)
	if err != nil {
		return nil, err
	}

	data := make([]byte, dataSize)
	pos := 0
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldSize := 0
		switch field.Kind() {
		case reflect.Uint8:
			fieldSize = 1
			data[pos] = uint8(field.Uint())
		case reflect.Uint16:
			fieldSize = 2
			binary.BigEndian.PutUint16(data[pos:pos+fieldSize], uint16(field.Uint()))
		case reflect.Int32:
			fieldSize = 4
			binary.BigEndian.PutUint32(data[pos:pos+fieldSize], uint32(field.Int()))
		case reflect.Int64:
			fieldSize = 8
			binary.BigEndian.PutUint64(data[pos:pos+fieldSize], uint64(field.Int()))
		case reflect.Slice:
			elementKind := field.Type().Elem().Kind()
			if elementKind != reflect.Uint8 {
				return nil, errors.Errorf("Cannot marshal slice of %s, only slices of uint8 are supported", elementKind)
			}

			fieldSize = field.Len()
			reflect.Copy(reflect.ValueOf(data[pos:pos+fieldSize]), field)
		case reflect.Array:
			elementKind := field.Type().Elem().Kind()
			if elementKind != reflect.Uint8 {
				return nil, errors.Errorf("Cannot marshal array of %s, only arrays of uint8 are supported", elementKind)
			}

			fieldSize = field.Len()
			reflect.Copy(reflect.ValueOf(data[pos:pos+fieldSize]), field)

		default:
			return nil, errors.Errorf("Cannot marshal kind %s", field.Kind())
		}
		pos += fieldSize
	}

	return data, nil
}

// Unmarshal decodes data from byte slice following the XRootd protocol specification.
// Fields are decoded in the same order as they are specified in the struct definition.
// Each field is decoded from network byte order (BigEndian) without any alignment or padding.
// Slices and arrays of uint8 are binary copied without further decoding.
// Since the length of the slice is unknown, all bytes to the end of data is copied to it.
// Supported types are: uint8, uint16, int32, int64, slices and arrays of uint8.
func Unmarshal(data []byte, x interface{}) (err error) {
	v := reflect.ValueOf(x)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return errors.Errorf("Cannot unmarshal %s, struct is expected", v.Kind())
	}

	pos := 0

	for i := 0; i < v.NumField() && err == nil; i++ {
		field := v.Field(i)
		fieldSize := 0
		switch field.Kind() {
		case reflect.Uint8:
			fieldSize = 1
			field.SetUint(uint64(data[pos]))
		case reflect.Uint16:
			fieldSize = 2
			var value = binary.BigEndian.Uint16(data[pos : pos+fieldSize])
			field.SetUint(uint64(value))
		case reflect.Int32:
			fieldSize = 4
			var value = int32(binary.BigEndian.Uint32(data[pos : pos+fieldSize]))
			field.SetInt(int64(value))
		case reflect.Int64:
			fieldSize = 8
			var value = int64(binary.BigEndian.Uint64(data[pos : pos+fieldSize]))
			field.SetInt(value)
		case reflect.Slice:
			elementKind := field.Type().Elem().Kind()
			if elementKind != reflect.Uint8 {
				return errors.Errorf("Cannot unmarshal slice of %s, only slices of uint8 are supported", elementKind)
			}

			bytes := data[pos:]
			fieldSize = len(bytes)
			field.SetBytes(bytes)
		case reflect.Array:
			elementKind := field.Type().Elem().Kind()
			if elementKind != reflect.Uint8 {
				return errors.Errorf("Cannot unmarshal array of %s, only arrays of uint8 are supported", elementKind)
			}

			fieldSize = field.Len()
			reflect.Copy(field, reflect.ValueOf(data[pos:pos+fieldSize]))
		default:
			err = errors.Errorf("Cannot unmarshal kind %s", field.Kind())
		}
		pos += fieldSize
	}
	return
}

func calculateSizeForMarshaling(v reflect.Value) (size int, err error) {
	for i := 0; i < v.NumField() && err == nil; i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.Uint8:
			size++
		case reflect.Uint16:
			size += 2
		case reflect.Int32:
			size += 4
		case reflect.Int64:
			size += 8
		case reflect.Array:
			size += field.Len()
		case reflect.Slice:
			size += field.Len()
		default:
			err = errors.Errorf("Cannot marshal kind %s", field.Kind())
		}
	}
	return
}
