// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"io"

	"go-hep.org/x/hep/fwk"
)

// UI wraps a fwk.Scripter and panics when an error occurs
type UI struct {
	ui fwk.Scripter
}

// Configure configures the underlying fwk.App.
// Configure panics if an error occurs.
func (ui UI) Configure() {
	err := ui.ui.Configure()
	if err != nil {
		panic(err)
	}
}

// Start starts the underlying fwk.App.
// Start panics if an error occurs.
func (ui UI) Start() {
	err := ui.ui.Start()
	if err != nil {
		panic(err)
	}
}

// Run runs the event-loop of the underlying fwk.App.
// Run panics if an error different than io.EOF occurs.
func (ui UI) Run(evtmax int64) {
	err := ui.ui.Run(evtmax)
	if err != nil && err != io.EOF {
		panic(err)
	}
}

// Stop stops the underlying fwk.App.
// Stopt panics if an error occurs.
func (ui UI) Stop() {
	err := ui.ui.Stop()
	if err != nil {
		panic(err)
	}
}

// Shutdown shuts the underlying fwk.App down.
// Shutdown panics if an error occurs.
func (ui UI) Shutdown() {
	err := ui.ui.Shutdown()
	if err != nil {
		panic(err)
	}
}
