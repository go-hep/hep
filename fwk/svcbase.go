// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"go-hep.org/x/hep/fwk/fsm"
)

// SvcBase provides a base implementation for fwk.Svc
type SvcBase struct {
	t   string
	n   string
	mgr App
}

// NewSvc creates a new SvcBase of type typ and name name,
// managed by the fwk.App mgr.
func NewSvc(typ, name string, mgr App) SvcBase {
	return SvcBase{
		t:   typ,
		n:   name,
		mgr: mgr,
	}
}

// Type returns the fully qualified type of the underlying service.
// e.g. "go-hep.org/x/hep/fwk/testdata.svc1"
func (svc *SvcBase) Type() string {
	return svc.t
}

// Name returns the name of the underlying service.
// e.g. "my-service"
func (svc *SvcBase) Name() string {
	return svc.n
}

// DeclProp declares this service has a property named n,
// and takes a pointer to the associated value.
func (svc *SvcBase) DeclProp(n string, ptr interface{}) error {
	return svc.mgr.DeclProp(svc, n, ptr)
}

// SetProp sets the property name n with the value v.
func (svc *SvcBase) SetProp(name string, value interface{}) error {
	return svc.mgr.SetProp(svc, name, value)
}

// GetProp returns the value of the property named n.
func (svc *SvcBase) GetProp(name string) (interface{}, error) {
	return svc.mgr.GetProp(svc, name)
}

// FSMState returns the current state of the FSM
func (svc *SvcBase) FSMState() fsm.State {
	return svc.mgr.FSMState()
}
