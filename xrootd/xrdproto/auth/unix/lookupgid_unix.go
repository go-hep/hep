// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !windows
// +build !windows

package unix // import "go-hep.org/x/hep/xrootd/xrdproto/auth/unix"

import (
	"os/user"
)

func lookupGroupID(usr *user.User) (string, error) {
	group, err := user.LookupGroupId(usr.Gid)
	if err != nil {
		return "", err
	}
	return group.Name, nil
}
