// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package xrdio

import (
	"fmt"
	"net"
	"strings"
)

// URL stores an absolute reference to a XRootD path.
type URL struct {
	Addr string // address (host [:port]) of the server
	User string // user name to use to log in
	Path string // path to the remote file or directory
}

// Parse parses name into an xrootd URL structure.
func Parse(name string) (URL, error) {
	var (
		user string
		addr string
		path string
		err  error
	)

	idx := strings.Index(name, "://")
	switch idx {
	case -1:
		path = name
	default:
		uri := name[idx+len("://"):]
		tok := strings.SplitN(uri, "/", 2)
		user, addr, err = parseUA(tok[0])
		if err != nil {
			return URL{}, fmt.Errorf("could not parse URI %q: %w", name, err)
		}
		path = "/" + tok[1]
	}

	if strings.HasPrefix(path, "//") {
		path = path[1:]
	}

	return URL{Addr: addr, User: user, Path: path}, nil
}

func parseUA(s string) (user, addr string, err error) {
	switch {
	case strings.Contains(s, "@"):
		toks := strings.SplitN(s, "@", 2)
		user = parseUser(toks[0])
		addr = toks[1]
	default:
		addr = s
	}

	switch {
	case strings.HasPrefix(addr, "["): // IPv6 literal
		idx := strings.LastIndex(addr, "]")
		col := strings.Index(addr[idx+1:], ":")
		if col >= 0 {
			_, _, err = net.SplitHostPort(addr)
		}
	case strings.Contains(addr, ":"):
		_, _, err = net.SplitHostPort(addr)
	}

	if err != nil {
		return "", "", fmt.Errorf("could not extract host+port from URI: %w", err)
	}

	return user, addr, nil
}

func parseUser(s string) string {
	idx := strings.Index(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}
