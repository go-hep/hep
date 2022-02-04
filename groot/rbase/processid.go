// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rbase

import (
	"fmt"
	"reflect"

	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/root"
	"go-hep.org/x/hep/groot/rtypes"
	"go-hep.org/x/hep/groot/rvers"
)

// ProcessID is ROOT's way to provide a process identifier object.
type ProcessID struct {
	named Named

	objs map[uint32]root.Object
}

func (*ProcessID) Class() string {
	return "TProcessID"
}

func (pid *ProcessID) UID() uint32 {
	return pid.named.UID()
}

func (*ProcessID) RVersion() int16 {
	return rvers.ProcessID
}

// Name returns the name of the instance
func (pid *ProcessID) Name() string {
	return pid.named.Name()
}

// Title returns the title of the instance
func (pid *ProcessID) Title() string {
	return pid.named.Title()
}

func (pid *ProcessID) SetName(name string)   { pid.named.SetName(name) }
func (pid *ProcessID) SetTitle(title string) { pid.named.SetTitle(title) }

func (pid *ProcessID) String() string {
	return fmt.Sprintf("%s{Name: %s, Title: %s}", pid.Class(), pid.Name(), pid.Title())
}

func (pid *ProcessID) UnmarshalROOT(r *rbytes.RBuffer) error {
	if r.Err() != nil {
		return r.Err()
	}

	beg := r.Pos()
	vers, pos, bcnt := r.ReadVersion(pid.Class())
	if vers > rvers.ProcessID {
		panic(fmt.Errorf("rbase: invalid %s version=%d > %d", pid.Class(), vers, rvers.ProcessID))
	}

	r.ReadObject(&pid.named)

	r.CheckByteCount(pos, bcnt, beg, pid.Class())
	return r.Err()
}

func (pid *ProcessID) MarshalROOT(w *rbytes.WBuffer) (int, error) {
	if w.Err() != nil {
		return 0, w.Err()
	}

	pos := w.WriteVersion(pid.RVersion())
	w.WriteObject(&pid.named)

	return w.SetByteCount(pos, pid.Class())
}

func init() {
	f := func() reflect.Value {
		o := &ProcessID{
			named: *NewNamed("", ""),
		}
		return reflect.ValueOf(o)
	}
	rtypes.Factory.Add("TProcessID", f)
}

var (
	_ root.Object        = (*ProcessID)(nil)
	_ root.UIDer         = (*ProcessID)(nil)
	_ root.Named         = (*ProcessID)(nil)
	_ rbytes.Marshaler   = (*ProcessID)(nil)
	_ rbytes.Unmarshaler = (*ProcessID)(nil)
)
