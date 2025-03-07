// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
)

// type arrayMode int32

type elemDescr struct {
	otype  rmeta.Enum
	ntype  rmeta.Enum
	offset int // actually an index to the struct's field or to array's element
	length int
	elem   rbytes.StreamerElement
	method []int
	oclass string
	nclass string
	mbr    any // member streamer
}

type streamerConfig struct {
	si     *StreamerInfo
	eid    int // element ID
	descr  *elemDescr
	offset int // offset/index within object. negative if no offset to be applied.
	length int // number of elements for fixed-length arrays

	count func() int // optional func to give the length of ROOT's C var-len arrays.
}

func (cfg *streamerConfig) counter(recv any) int {
	if cfg.count != nil {
		return cfg.count()
	}
	return int(reflect.ValueOf(recv).Elem().FieldByIndex(cfg.descr.method).Int())
}

func (cfg *streamerConfig) adjust(recv any) any {
	if cfg == nil || cfg.offset < 0 {
		return recv
	}
	rv := reflect.ValueOf(recv).Elem()
	switch rv.Kind() {
	case reflect.Struct:
		return rv.Field(cfg.offset).Addr().Interface()
	case reflect.Array, reflect.Slice:
		return rv.Index(cfg.offset).Addr().Interface()
	default:
		return recv
	}
}
