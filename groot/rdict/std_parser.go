// Copyright 2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"

	"go-hep.org/x/hep/groot/rmeta"
)

func parseStdSet(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "set", "std::set":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::set template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 1, 2, 3:
		// std::set<
		//  K,
		//  Cmp=std::less<K>,
		//  Alloc=std::allocator<K>,
		// >
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::set template %q", tmpl))
	}
}

func parseStdUnorderedSet(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "unordered_set", "std::unordered_set":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::unordered_set template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 1, 2, 3, 4:
		// std::unordered_set<
		//   K,
		//   Hash=std::hash<K>,
		//   Eq=std::equal_to<K>,
		//   Alloc=std::allocator<K>,
		// >
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::unordered_set template %q", tmpl))
	}
}

func parseStdList(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "list", "std::list":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::list template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 1, 2:
		// std::set<
		//  T,
		//  Allocator=std::allocator<T>,
		// >
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::list template %q", tmpl))
	}
}

func parseStdDeque(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "deque", "std::deque":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::deque template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 1, 2:
		// std::deque<
		//  T,
		//  Allocator=std::allocator<T>,
		// >
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::deque template %q", tmpl))
	}
}

func parseStdVector(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "vector", "std::vector":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::vector template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 1, 2:
		// std::vector<
		//  T,
		//  Allocator=std::allocator<T>,
		// >
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::vector template %q", tmpl))
	}
}

func parseStdMap(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "map", "std::map":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::map template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 1, 2, 3, 4:
		// std::map<
		//  K,
		//  V,
		//  Cmp=std::less<K>,
		//  Allocator=std::allocator<std::pair<const K,V>>,
		// >
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::map template %q", tmpl))
	}
}

func parseStdUnorderedMap(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "unordered_map", "std::unordered_map":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::unordered_map template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 1, 2, 3, 4, 5:
		// std::unordered_map<
		//  K,
		//  V,
		//  Hash=std::hash<K>,
		//  Eq=std::equal_to<K>,
		//  Allocator=std::allocator<std::pair<K,V>>,
		// >
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::unordered_map template %q", tmpl))
	}
}

func parseStdBitset(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	if len(cxx.Args) > 1 {
		panic(fmt.Errorf("invalid std::bitset template %q", tmpl))
	}
	stdArgNonEmpty(tmpl, cxx.Args)
	return cxx.Args
}

func parseStdPair(tmpl string) []string {
	cxx := rmeta.CxxTemplateFrom(tmpl)
	switch cxx.Name {
	case "pair", "std::pair":
		// ok.
	default:
		panic(fmt.Errorf("rdict: invalid std::pair template %q", tmpl))
	}
	switch len(cxx.Args) {
	case 2:
		// std::pair<T1,T2>
		stdArgNonEmpty(tmpl, cxx.Args)
		return cxx.Args
	default:
		panic(fmt.Errorf("rdict: invalid std::pair template %q", tmpl))
	}
}

func stdArgNonEmpty(typename string, vs []string) {
	for _, v := range vs {
		if v == "" {
			panic(fmt.Errorf("rdict: invalid empty type argument %q", typename))
		}
	}
}
