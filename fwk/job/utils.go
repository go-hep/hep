// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package job

import (
	"fmt"
	"strings"

	"go-hep.org/x/hep/fwk"
)

// MsgLevel returns the fwk.Level according to the given lvl string value.
// MsgLevel panics if no fwk.Level value corresponds to the lvl string value.
// Valid values are: "DEBUG", "INFO", "WARNING"|"WARN" and "ERROR"|"ERR".
func MsgLevel(lvl string) fwk.Level {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return fwk.LvlDebug
	case "INFO":
		return fwk.LvlInfo
	case "WARNING", "WARN":
		return fwk.LvlWarning
	case "ERROR", "ERR":
		return fwk.LvlError
	default:
		panic(fmt.Errorf("fwk.MsgLevel: invalid fwk.Level string %q", lvl))
	}
}
