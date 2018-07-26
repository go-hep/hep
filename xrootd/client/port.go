// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"net"
	"strconv"
)

func parseAddr(addr string) string {
	_, _, err := net.SplitHostPort(addr)
	if err == nil {
		return addr
	}
	switch err := err.(type) {
	case *net.AddrError:
		switch err.Err {
		case "missing port in address":
			port, e := net.LookupPort("tcp", "rootd")
			if e != nil {
				return addr + ":1094"
			}
			return addr + ":" + strconv.Itoa(port)
		}
	}
	return addr
}
