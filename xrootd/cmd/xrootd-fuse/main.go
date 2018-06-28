// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command xrootd-fuse mounts the directory contents of a remote xrootd server to a local directory.
//
// Usage:
//
//  $> xrootd-fuse [OPTIONS] <remote-dir> <local-dir>
//
// Example:
//
//  $> xrootd-fuse root://server.example.com/some/dir /mnt
//  $> xrootd-fuse -v root://server.example.com/some/dir /mnt
//
// Options:
//   -v	enable verbose mode
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"go-hep.org/x/hep/xrootd/client"
	"go-hep.org/x/hep/xrootd/xrdfuse"
	"go-hep.org/x/hep/xrootd/xrdio"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `xrootd-fuse mounts the directory contents of a remote xrootd server to a local directory.

Usage:

 $> xrootd-fuse [OPTIONS] <remote-dir> <local-dir>

Example:

 $> xrootd-fuse root://server.example.com/some/dir /mnt
 $> xrootd-fuse -v root://server.example.com/some/dir /mnt

Options:
`)
		flag.PrintDefaults()
	}
}

func main() {
	log.SetPrefix("xrootd-fuse: ")
	log.SetFlags(0)

	verbose := flag.Bool("v", false, "enable verbose mode")

	flag.Parse()

	if flag.NArg() != 2 {
		flag.Usage()
		log.Fatalf("missing directory operands")
	}

	url, err := xrdio.Parse(flag.Arg(0))

	c, err := client.NewClient(context.Background(), url.Addr, url.User)
	if err != nil {
		log.Fatalf("could not create client: %v", err)
	}

	fs := xrdfuse.NewFS(c, url.Path, func(e error) {
		log.Println(e)
	})
	nfs := pathfs.NewPathNodeFs(fs, nil)
	server, _, err := nodefs.MountRoot(flag.Arg(1), nfs.Root(), &nodefs.Options{
		Debug: *verbose,
	})
	if err != nil {
		log.Fatalf("could not mount: %v", err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	go server.Serve()

	<-ch
	err = server.Unmount()
	if err != nil {
		log.Fatalf("could not unmount: %v", err)
	}
}
