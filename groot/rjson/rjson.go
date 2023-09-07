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
	"log"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rdict"
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

	vers := -1
	if o, ok := o.(interface{ RVersion() int }); ok {
		vers = o.RVersion()
	}
	log.Printf("--> marshal %T", o)
	si, err := rdict.StreamerInfos.StreamerInfo(o.Class(), vers)
	if err != nil {
		return nil, fmt.Errorf("could not find streamer info for %T: %w", o, err)
	}

	w := new(bytes.Buffer)
	w.WriteByte('{')
	for i, elt := range si.Elements() {
		log.Printf("elt: %+v", elt)
		switch elt := elt.(type) {
		case *rdict.StreamerBase:
			log.Printf(">>> %+v", elt)
		}
		if i > 0 {
			w.WriteByte(',')
		}
		fmt.Fprintf(w, "%q: %q", "_typename", elt.Name())
	}
	w.WriteByte('}')

	rv := reflect.Indirect(reflect.ValueOf(o))
	for i := 0; i < rv.NumField(); i++ {
		rf := rv.Field(i)
		ft := rv.Type().Field(i)
		var (
			vv   any
			name = ft.Name
		)
		if ft.IsExported() {
			vv = rf.Interface()
		}
		log.Printf("field[%d]: %q %+v (%s)", i, name, vv, ft.Type.Name())
	}

	return w.Bytes(), nil
}

type encoder struct {
	w io.Writer
}

func newEncoder(w io.Writer) *encoder {
	return &encoder{w: w}
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
		if v, ok := m.Value.(rbytes.RSlicer); ok {
			err := enc.Encode(v)
			if err != nil {
				return err
			}
			continue
		}
		if v := reflect.Indirect(reflect.ValueOf(m.Value)); v.Kind() == reflect.Slice && v.Len() == 0 {
			_, err = w.Write([]byte("[]"))
			if err != nil {
				return err
			}
			continue
		}
		v, err := json.Marshal(m.Value)
		if err != nil {
			return err
		}
		_, err = w.Write(v)
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
