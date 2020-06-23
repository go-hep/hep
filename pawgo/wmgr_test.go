// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"testing"

	"gioui.org/io/key"
	"gioui.org/io/router"
	"gioui.org/io/system"
	"gioui.org/op"
	"go-hep.org/x/hep/hplot"
)

func TestPlot(t *testing.T) {
	wmgr := newWinMgr(nil)
	defer wmgr.Close()

	p := hplot.New()
	p.Title.Text = "my plot"
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	win := wmgr.newPlot(p)
	<-win.ready

	rc := win.handle(system.FrameEvent{
		Frame: func(frame *op.Ops) {},
		Queue: new(router.Router),
	})
	if got, want := rc, winContinue; got != want {
		t.Fatalf("invalid window state: got=%v, want=%v", got, want)
	}

	rc = win.handle(key.Event{Name: "Q"})
	if got, want := rc, winContinue; got != want {
		t.Fatalf("invalid window state: got=%v, want=%v", got, want)
	}

	rc = win.handle(system.DestroyEvent{Err: nil})
	if got, want := rc, winStop; got != want {
		t.Fatalf("invalid window state: got=%v, want=%v", got, want)
	}

	err := wmgr.Close()
	if err != nil {
		t.Fatalf("could not close wmgr: %+v", err)
	}

	wmgr.wg.Wait()
}
