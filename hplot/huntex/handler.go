// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package huntex

import (
	"log"
	"reflect"
	"unsafe"

	"gonum.org/v1/plot"
)

type Handler struct {
	p   *plot.Plot
	set map[*string]hctx
	rep replace
}

type hctx struct {
	old string
}

func NewHandler(p *plot.Plot) *Handler {
	h := &Handler{
		p:   p,
		set: make(map[*string]hctx),
	}
	h.run()
	return h
}

func (h *Handler) run() {
	//	h.untex(&h.p.Title.Text)
	//	h.untex(&h.p.X.Label.Text)
	//	h.untex(&h.p.Y.Label.Text)
	h.reflect(reflect.ValueOf(h.p))
}

func (h *Handler) reflect(v reflect.Value) {
	v = reflect.Indirect(v)
	switch v.Kind() {
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			h.reflect(e)
		}

	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			e := v.Index(i)
			h.reflect(e)
		}

	case reflect.String:
		ptr := (*string)(unsafe.Pointer(v.UnsafeAddr()))
		//ptr := v.Addr().Interface().(*string)
		h.untex(ptr)

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if !f.CanAddr() {
				log.Printf("discarding (addr=%v,set=%v) %#v (%T)", f.CanAddr(), f.CanSet(), v.Type().Field(i), v.Interface())
				continue
			}
			h.reflect(f)
		}
	}
}

func (h *Handler) untex(s *string) {
	v := h.rep.replace(*s)
	h.set[s] = hctx{old: *s}
	*s = v
}

func (h *Handler) Close() error {
	for ptr, ctx := range h.set {
		*ptr = ctx.old
	}
	return nil
}
