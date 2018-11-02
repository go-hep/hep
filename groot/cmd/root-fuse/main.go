// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !windows

// Command root-fuse mounts the contents of a ROOT file as a local directory.
//
// Usage:
//
//   $> root-fuse [OPTIONS] <ROOT file> <mount-dir>
//
// Example:
//
//   $> root-fuse ./testdata/simple.root /mnt/dir
//   $> root-fuse -v ./testdata/simple.root /mnt/dir
//   $> root-fuse root://eospublic.cern.ch:1094//eos/opendata/atlas/OutreachDatasets/2016-07-29/MC/mc_173045.DYtautauM08to15.root /mnt/dir
//
// Options:
//   -v	enable verbose mode
//
package main // import "go-hep.org/x/hep/groot/cmd/root-fuse"

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	"go-hep.org/x/hep/groot"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `root-fuse mounts the contents of a ROOT file as a local directory.

Usage:

  $> root-fuse [OPTIONS] <ROOT file> <mount-dir>

Example:

  $> root-fuse ./testdata/simple.root /mnt/dir
  $> root-fuse root://server.example.root/data/simple.root /mnt/dir
  $> root-fuse -v ./testdata/simple.root /mnt/dir

Options:
`)
		flag.PrintDefaults()
	}
}

func main() {
	log.SetPrefix("root-fuse: ")
	log.SetFlags(0)

	dbg := flag.Bool("v", false, "enable verbose mode")

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		log.Printf("missing directory operands")
		os.Exit(2)
	}

	src := flag.Arg(0)
	dst := flag.Arg(1)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	ready := make(chan struct{})
	go run(dst, src, *dbg, ready, c)
	<-ready
	<-c
}

func run(target, fname string, verbose bool, ready chan struct{}, c chan os.Signal) {
	err := os.MkdirAll(target, 0755)
	if err != nil {
		log.Fatal(err)
	}

	f, err := groot.Open(fname)
	if err != nil {
		log.Fatalf("could not open ROOT file %q: %v", fname, err)
	}
	defer f.Close()

	fs := NewFS(f)
	if fs == nil {
		log.Fatalf("could not create FUSE FS")
	}

	nfs := pathfs.NewPathNodeFs(fs, nil)

	server, _, err := nodefs.MountRoot(target, nfs.Root(), &nodefs.Options{
		Debug: verbose,
	})

	if err != nil {
		log.Fatalf("mount failed: %v", err)
	}

	select {
	case ready <- struct{}{}:
	default:
	}

	go server.Serve()

	select {
	case v := <-c:
		if verbose {
			log.Printf("unmounting %q...", target)
		}
		err := server.Unmount()
		if err != nil {
			log.Printf("could not unmount %q: %v", target, err)
			return
		}
		if verbose {
			log.Printf("unmounting %q... [done]", target)
		}
		c <- v
	}
}
