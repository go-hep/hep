// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore

// list-groot-types lists all the groot types registered with the rtypes factory.
package main

import (
	"log"
	"sort"
	"strings"

	"go-hep.org/x/hep/groot/rtypes"
	_ "go-hep.org/x/hep/groot/ztypes"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	var (
		croot []string
		groot []string
	)

	for _, k := range rtypes.Factory.Keys() {
		switch {
		case strings.HasPrefix(k, "*") || strings.Contains(k, "."):
			groot = append(groot, k)
		default:
			croot = append(croot, k)
		}
	}

	sort.Strings(croot)
	sort.Strings(groot)

	log.Printf(strings.Repeat("*", 80))
	log.Printf("C++/ROOT types: %d", len(croot))
	for _, k := range croot {
		log.Printf("%s", k)
	}

	log.Printf(strings.Repeat("*", 80))
	log.Printf("groot types: %d", len(groot))
	for _, k := range groot {
		log.Printf("%s", k)
	}

	log.Printf("%d known types", rtypes.Factory.Len())
	log.Printf(strings.Repeat("*", 80))
}
