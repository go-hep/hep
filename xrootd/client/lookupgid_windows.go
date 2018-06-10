// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build windows

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"os/user"
)

func lookupGroupID(usr *user.User) (string, error) {
	// Since user.LookupGroupId is not implemented under Windows fallback to the username.
	return usr.Username, nil
}
