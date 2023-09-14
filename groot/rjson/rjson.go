// Copyright ©2023 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rjson contains tools to marshal ROOT objects to JSON.
package rjson // import "go-hep.org/x/hep/groot/rjson"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
)

func Marshal(o root.Object) ([]byte, error) {
	if o, ok := o.(rbytes.RSlicer); ok {
		buf := new(bytes.Buffer)
		enc := newEncoder(buf)
		err := enc.Encode(o)
		return buf.Bytes(), err
	}

	panic(fmt.Errorf("not implemented for %T", o))
}

type encoder struct {
	w io.Writer
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{w: w}
}

func (enc *encoder) encode(v any) error {
	if v, ok := v.(rbytes.RSlicer); ok {
		return enc.Encode(v)
	}

	var (
		err error
		w   = enc.w
		rv  = reflect.Indirect(reflect.ValueOf(v))
	)

	switch rv.Kind() {
	case reflect.Slice:
		_, err = w.Write([]byte("["))
		if err != nil {
			return err
		}
		for i := 0; i < rv.Len(); i++ {
			if i > 0 {
				_, err := w.Write([]byte(","))
				if err != nil {
					return err
				}
			}
			err = enc.encode(rv.Index(i).Interface())
			if err != nil {
				return err
			}
		}
		_, err = w.Write([]byte("]"))
		if err != nil {
			return err
		}
	default:
		raw, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = w.Write(raw)
		if err != nil {
			return err
		}
	}

	return nil
}

func (enc *encoder) Encode(v rbytes.RSlicer) error {
	mbrs := v.RMembers()
	name := v.(root.Object).Class()
	w := enc.w
	_, err := fmt.Fprintf(w, "{%q: %q", "_typename", name)
	if err != nil {
		return err
	}
	for _, m := range mbrs {
		_, err := fmt.Fprintf(w, ", %q: ", m.Name)
		if err != nil {
			return err
		}
		err = enc.encode(m.Value)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "}")
	if err != nil {
		return err
	}
	return nil
}
