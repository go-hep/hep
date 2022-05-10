// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"go/types"
	"strings"

	"go-hep.org/x/hep/groot/rmeta"
)

// Generator is the interface used to generate ROOT related code.
type Generator interface {
	// Generate generates code for a given type.
	Generate(typ string) error

	// Format formats the Go generated code.
	Format() ([]byte, error)
}

func gotype2RMeta(t types.Type) rmeta.Enum {
	switch ut := t.Underlying().(type) {
	case *types.Basic:
		switch ut.Kind() {
		case types.Bool:
			return rmeta.Bool
		case types.Uint8:
			return rmeta.Uint8
		case types.Uint16:
			return rmeta.Uint16
		case types.Uint32, types.Uint:
			return rmeta.Uint32
		case types.Uint64:
			return rmeta.Uint64
		case types.Int8:
			return rmeta.Int8
		case types.Int16:
			return rmeta.Int16
		case types.Int32, types.Int:
			return rmeta.Int32
		case types.Int64:
			return rmeta.Int64
		case types.Float32:
			return rmeta.Float32
		case types.Float64:
			return rmeta.Float64
		case types.String:
			return rmeta.TString
		}
	case *types.Struct:
		return rmeta.Any
	case *types.Slice:
		return rmeta.STL
	case *types.Array:
		return rmeta.OffsetL + gotype2RMeta(ut.Elem())
	}
	return -1
}

// GoName2Cxx translates a fully-qualified Go type name to a C++ one.
// e.g.:
//   - go-hep.org/x/hep/hbook.H1D -> go_hep_org::x::hep::hbook::H1D
func GoName2Cxx(name string) string {
	repl := strings.NewReplacer(
		"-", "_",
		"/", "::",
		".", "_",
	)
	i := strings.LastIndex(name, ".")
	if i > 0 {
		name = name[:i] + "::" + name[i+1:]
	}
	return repl.Replace(name)
}

// Typename returns a language dependent typename, usually encoded inside a
// StreamerInfo's title.
func Typename(name, title string) (string, bool) {
	if title == "" {
		return name, false
	}
	i := strings.Index(title, ";")
	if i <= 0 {
		return name, false
	}
	lang := title[:i]
	title = strings.TrimSpace(title[i+1:])
	switch lang {
	case "Go":
		return title, GoName2Cxx(title) == name
	default:
		return title, false
	}
}
