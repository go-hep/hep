// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rmeta provides tools to interoperate with ROOT Meta.
package rmeta // import "go-hep.org/x/hep/groot/rmeta"

import (
	"strings"
)

// CxxTemplateArgsOf extracts the typenames of a C++ templated typename.
// Ex:
//  std::map<K,V> -> []string{"K", "V"}
//  std::vector<T> -> []string{"T"}
//  Foo<T1,T2,std::map<K,V>> -> []string{"T1", "T2", "std::map<K,V>"}
func CxxTemplateArgsOf(typename string) []string {
	name := strings.TrimSpace(typename)
	name = name[strings.Index(name, "<")+1:] // drop heading 'xxx<'
	name = name[:len(name)-1]                // drop trailing '>'
	name = strings.TrimSpace(name)

	switch strings.Count(name, ",") {
	case 0:
		return []string{name}
	case 1:
		// easy case of std::map<K,V> where none of K or V are templated.
		i := strings.Index(name, ",")
		k := strings.TrimSpace(name[:i])
		v := strings.TrimSpace(name[i+1:])
		return []string{k, v}
	default:
		var (
			types []string
			bldr  strings.Builder
			tmpl  int
		)
		for _, s := range name {
			switch s {
			case '<':
				tmpl++
				bldr.WriteRune(s)
			case '>':
				tmpl--
				bldr.WriteRune(s)
			case ',':
				switch {
				case tmpl > 0:
					bldr.WriteRune(s)
				default:
					types = append(types, strings.TrimSpace(bldr.String()))
					bldr.Reset()
				}
			default:
				bldr.WriteRune(s)

			}
		}
		types = append(types, strings.TrimSpace(bldr.String()))
		return types
	}
}
