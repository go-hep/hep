// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"context"
)

type ctxType struct {
	id    int64
	slot  int
	store Store
	msg   msgstream
	mgr   App

	ctx context.Context
}

func (ctx ctxType) ID() int64 {
	return ctx.id
}

func (ctx ctxType) Slot() int {
	return ctx.slot
}

func (ctx ctxType) Store() Store {
	return ctx.store
}

func (ctx ctxType) Msg() MsgStream {
	return ctx.msg
}

func (ctx ctxType) Svc(n string) (Svc, error) {
	if ctx.mgr == nil {
		return nil, Errorf("fwk: no fwk.App available to this Context")
	}

	svc := ctx.mgr.GetSvc(n)
	if svc == nil {
		return nil, Errorf("fwk: no such service [%s]", n)
	}
	return svc, nil
}
