// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package yodacnv provides tools to read/write YODA archive files.
package yodacnv

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"reflect"

	"go-hep.org/x/hep/hbook"
)

var (
	begYoda = []byte("BEGIN YODA_")
	endYoda = []byte("END YODA_")
)

// Read reads a YODA stream and converts the YODA values into their
// go-hep/hbook equivalents.
func Read(r io.Reader) ([]hbook.Object, error) {
	var (
		err   error
		o     []hbook.Object
		block = make([]byte, 0, 1024)
		rt    reflect.Type
	)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		raw := scan.Bytes()
		switch {
		case bytes.HasPrefix(raw, begYoda):
			rt, err = splitHeader(raw)
			if err != nil {
				return nil, fmt.Errorf("yoda: error parsing YODA header (%v)", err)
			}
			block = block[:0]
			block = append(block, raw...)
			block = append(block, '\n')

		default:
			block = append(block, raw...)
			block = append(block, '\n')

		case bytes.HasPrefix(raw, endYoda):
			block = append(block, raw...)
			block = append(block, '\n')

			v := reflect.New(rt).Elem()
			err = v.Addr().Interface().(Unmarshaler).UnmarshalYODA(block)
			if err != nil {
				return nil, err
			}
			o = append(o, v.Addr().Interface().(hbook.Object))
		}
	}
	err = scan.Err()
	if err != nil {
		return nil, err
	}
	return o, nil
}

// Write writes values to a YODA stream.
func Write(w io.Writer, args ...Marshaler) error {
	for _, v := range args {
		raw, err := v.MarshalYODA()
		if err != nil {
			return err
		}
		n, err := w.Write(raw)
		if err != nil {
			return err
		}
		if n < len(raw) {
			return io.ErrShortWrite
		}
	}
	return nil
}

func splitHeader(raw []byte) (reflect.Type, error) {
	raw = raw[len(begYoda):]
	i := bytes.Index(raw, []byte(" "))
	if i == -1 || i >= len(raw) {
		return nil, fmt.Errorf("invalid YODA header (missing space)")
	}

	var rt reflect.Type

	switch string(raw[:i]) {
	case "HISTO1D":
		rt = reflect.TypeOf((*hbook.H1D)(nil)).Elem()
	case "HISTO2D":
		rt = reflect.TypeOf((*hbook.H2D)(nil)).Elem()
	case "PROFILE1D":
		rt = reflect.TypeOf((*hbook.P1D)(nil)).Elem()
	case "SCATTER2D":
		rt = reflect.TypeOf((*hbook.S2D)(nil)).Elem()
	default:
		return nil, fmt.Errorf("unhandled YODA object type %q", string(raw[:i]))
	}

	return rt, nil
}

// Unmarshaler is the interface implemented by an object that can
// unmarshal a YODA representation of itself.
type Unmarshaler interface {
	UnmarshalYODA([]byte) error
}

// Marshaler is the interface implemented by an object that can
// marshal itself into a YODA form.
type Marshaler interface {
	MarshalYODA() ([]byte, error)
}
