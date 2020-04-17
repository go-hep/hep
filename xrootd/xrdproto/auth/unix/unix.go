// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package unix contains the implementation of unix security provider.
package unix // import "go-hep.org/x/hep/xrootd/xrdproto/auth/unix"

import (
	"os/user"

	"go-hep.org/x/hep/xrootd/xrdproto/auth"
)

// Default is an unix security provider configured from current username and group.
// If the credentials could not be correctly configured, Default will be nil.
var Default auth.Auther

func init() {
	u, err := user.Current()
	if err != nil {
		return
	}
	g, err := lookupGroupID(u)
	if err != nil {
		return
	}
	Default = &Auth{User: u.Username, Group: g}
}

// Auth implements unix security provider.
type Auth struct {
	User  string
	Group string
}

// Provider implements auth.Auther
func (*Auth) Provider() string {
	return "unix"
}

// Type indicates that unix authentication protocol is used.
var Type = [4]byte{'u', 'n', 'i', 'x'}

// Request implements auth.Auther
func (a *Auth) Request(params []string) (*auth.Request, error) {
	return &auth.Request{Type: Type, Credentials: "unix\000" + a.User + " " + a.Group + "\000"}, nil
}

var (
	_ auth.Auther = (*Auth)(nil)
)
