// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"

	"go-hep.org/x/hep/rootio/cmd/root-srv/server"
)

var (
	addrFlag = flag.String("addr", ":8080", "server address:port")
)

func main() {
	flag.Parse()

	server.Init()
	log.Printf("server listening on %s", *addrFlag)
	log.Fatal(http.ListenAndServe(*addrFlag, nil))
}
