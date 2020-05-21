// Copyright ©2015 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"reflect"
)

// rioAlignU32 returns sz adjusted to align at 4-byte boundaries
func rioAlignU32(sz uint32) uint32 {
	return sz + (4-(sz&gAlign))&gAlign
}

// rioAlign returns sz adjusted to align at 4-byte boundaries
func rioAlign(sz int) int {
	return sz + (4-(sz&int(gAlign)))&int(gAlign)
}

func nameFromType(rt reflect.Type) string {
	if rt == nil {
		return "interface"
	}
	// Default to printed representation for unnamed types
	name := rt.String()

	// But for named types (or pointers to them), qualify with import path.
	// Dereference one pointer looking for a named type.
	star := ""
	if rt.Name() == "" {
		pt := rt
		if pt.Kind() == reflect.Ptr {
			star = "*"
			rt = pt.Elem()
		}
	}

	if rt.Name() != "" {
		switch rt.PkgPath() {
		case "":
			name = star + rt.Name()
		default:
			name = star + rt.PkgPath() + "." + rt.Name()
		}
	}

	return name
}

// EOF
