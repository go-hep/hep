// Copyright 2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

// list-groot-sinfos lists all the known StreamerInfos known to groot.
package main

import (
	"log"

	"go-hep.org/x/hep/groot/rdict"
	_ "go-hep.org/x/hep/groot/ztypes"
)

func main() {
	log.SetPrefix("")
	log.SetFlags(0)

	sinfos := rdict.StreamerInfos.Values()
	for _, si := range sinfos {
		log.Printf("version=%03d, checksum=0x%08x, name=%q", si.ClassVersion(), si.CheckSum(), si.Name())
	}

	log.Printf("%d known streamers", len(sinfos))
}
