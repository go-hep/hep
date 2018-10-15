// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package host contains the implementation for the "host" security provider.
package host // import "go-hep.org/x/hep/xrootd/xrdproto/auth/host"

import (
	"os"

	"go-hep.org/x/hep/xrootd/xrdproto/auth"
)

// Default is a host security provider configured from the current local host.
// If the credentials could not be correctly configured, Default will be nil.
var Default auth.Auther

func init() {
	host, err := os.Hostname()
	if err != nil {
		return
	}
	Default = &Auth{Hostname: host}
}

// Auth implements the host security provider.
type Auth struct {
	Hostname string
}

// Provider implements auth.Auther
func (*Auth) Provider() string {
	return "host"
}

// Type indicates the host authentication protocol is used.
var Type = [4]byte{'h', 'o', 's', 't'}

// Request implements auth.Auther
func (a *Auth) Request(params []string) (*auth.Request, error) {
	return &auth.Request{Type: Type, Credentials: "host\000" + a.Hostname + "\000"}, nil
}

var (
	_ auth.Auther = (*Auth)(nil)
)
