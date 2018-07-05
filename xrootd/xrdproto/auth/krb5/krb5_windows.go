// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build windows

package krb5

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const configPath = `c:\winnt\krb5.ini`

func cachePath() string {
	// FIXME: ask people for the popular windows krb5 client and use its cache path.
	// MIT krb5 lacks fresh windows builds IIUC.
	// As for now, hope that either KRB5CCNAME or %TEMP%\krb5cc_{uid} will work.
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
