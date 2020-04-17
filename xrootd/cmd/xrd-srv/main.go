// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command xrd-srv serves data from a local filesystem over the XRootD protocol.
package main // import "go-hep.org/x/hep/xrootd/cmd/xrd-srv"

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go-hep.org/x/hep/xrootd"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `xrd-srv serves data from a local filesystem over the XRootD protocol. 

Usage:

 $> xrd-srv [OPTIONS] <base-dir>

Example:

 $> xrd-srv /tmp
 $> xrd-srv -addr=0.0.0.0:1094 /tmp

Options:
`)
		flag.PrintDefaults()
	}
}

func main() {
	log.SetPrefix("xrd-srv: ")
	log.SetFlags(0)

	addr := flag.String("addr", "0.0.0.0:1094", "listen to the provided address")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing base dir operand")
	}

	baseDir := flag.Arg(0)

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("could not listen on %q: %v", *addr, err)
	}

	srv := xrootd.NewServer(xrootd.NewFSHandler(baseDir), func(err error) {
		log.Printf("an error occured: %v", err)
	})

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go func() {
		log.Printf("listening on %v...", listener.Addr())
		if err = srv.Serve(listener); err != nil {
			log.Fatalf("could not serve: %v", err)
		}
	}()

	<-ch
	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("could not shutdown: %v", err)
	}
}
