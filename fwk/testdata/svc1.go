// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testdata

import (
	"reflect"

	"go-hep.org/x/hep/fwk"
)

type MyInt int
type MyStruct struct {
	I int
}

type svc1 struct {
	fwk.SvcBase

	i MyInt
	s MyStruct
}

func (svc *svc1) Configure(ctx fwk.Context) error {
	var err error

	return err
}

func (svc *svc1) StartSvc(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("-- start svc --\n")
	return err
}

func (svc *svc1) StopSvc(ctx fwk.Context) error {
	var err error
	msg := ctx.Msg()
	msg.Infof("-- stop svc --\n")
	return err
}

func newsvc1(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error
	svc := &svc1{
		SvcBase: fwk.NewSvc(typ, name, mgr),
	}

	err = svc.DeclProp("Int", &svc.i)
	if err != nil {
		return nil, err
	}

	err = svc.DeclProp("Struct", &svc.s)
	if err != nil {
		return nil, err
	}

	return svc, err
}

func init() {
	fwk.Register(reflect.TypeOf(svc1{}), newsvc1)
}
