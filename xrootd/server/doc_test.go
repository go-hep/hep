// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server // import "go-hep.org/x/hep/xrootd/server"

import (
	"log"
	"net"
)

func Example() {
	addr := "0.0.0.0:1094"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("could not listen on %q: %v", addr, err)
	}

	srv := New(Default(), func(err error) {
		log.Printf("an error occured: %v", err)
	})

	log.Printf("listening on %v...", listener.Addr())

	if err = srv.Serve(listener); err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}
