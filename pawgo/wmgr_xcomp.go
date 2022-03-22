// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build cross_compile

package main

import (
	"log"

	"go-hep.org/x/hep/hplot"
)

type winMgr struct {
	msg *log.Logger
}

func newWinMgr(msg *log.Logger) *winMgr {
	return &winMgr{
		msg: msg,
	}
}

func (wmgr *winMgr) newPlot(p *hplot.Plot) *window {
	return nil
}

func (wmgr *winMgr) Close() error {
	return nil
}

type window struct {
}
