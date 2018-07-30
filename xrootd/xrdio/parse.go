// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdio

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// Parse parses name into an xrootd URL structure.
func Parse(name string) (URL, error) {
	urn, err := url.Parse(name)
	if err != nil {
		return URL{}, errors.WithStack(err)
	}

	host := urn.Hostname()
	port := urn.Port()

	path := urn.Path
	if strings.HasPrefix(path, "//") {
		path = path[1:]
	}

	user := ""
	if urn.User != nil {
		user = urn.User.Username()
	}

	addr := host
	if port != "" {
		addr += ":" + port
	}

	return URL{Addr: addr, User: user, Path: path}, nil
}

// URL stores an absolute reference to a XRootD path.
type URL struct {
	Addr string // address (host [:port]) of the server
	User string // user name to use to log in
	Path string // path to the remote file or directory
}
