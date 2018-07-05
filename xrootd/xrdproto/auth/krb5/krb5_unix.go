// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

package krb5

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// FIXME(sbinet): may be overwritten by $KRB5_CONFIG
// FIXME(sbinet): Linux puts it at:     /etc/krb5.conf
//                others may put it at: /etc/krb5/krb5.conf

const configPath = "/etc/krb5.conf"

func cachePath() string {
	if v := os.Getenv("KRB5CCNAME"); v != "" {
		if strings.HasPrefix(v, "FILE:") {
			v = string(v[len("FILE:"):])
		}
		return v
	}

	usr, err := user.Current()
	if err != nil {
		return ""
	}

	v := filepath.Join(os.TempDir(), fmt.Sprintf("krb5cc_%s", usr.Uid))
	return v
}
