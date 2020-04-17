// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fsm // import "go-hep.org/x/hep/fwk/fsm"

import (
	"fmt"
)

type State int

const (
	Undefined State = iota
	Configuring
	Configured
	Starting
	Started
	Running
	Stopping
	Stopped
	Offline
)

func (state State) String() string {
	switch state {
	case Undefined:
		return "UNDEFINED"
	case Configuring:
		return "CONFIGURING"
	case Configured:
		return "CONFIGURED"
	case Starting:
		return "STARTING"
	case Started:
		return "STARTED"
	case Running:
		return "RUNNING"
	case Stopping:
		return "STOPPING"
	case Stopped:
		return "STOPPED"
	case Offline:
		return "OFFLINE"

	default:
		panic(fmt.Errorf("invalid fsm.State value %d", int(state)))
	}
}
