// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command xrd-ls lists directory contents on a remote xrootd server.
//
// Usage:
//
//  $> xrd-ls [OPTIONS] <dir-1> [<dir-2> [...]]
//
// Example:
//
//  $> xrd-ls root://server.example.com/some/dir
//  $> xrd-ls -l root://server.example.com/some/dir
//  $> xrd-ls -R root://server.example.com/some/dir
//  $> xrd-ls -l -R root://server.example.com/some/dir
//
// Options:
//   -R	list subdirectories recursively
//   -l	use a long listing format
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"text/tabwriter"

	"github.com/pkg/errors"
	xrdclient "go-hep.org/x/hep/xrootd/client"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdio"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `xrd-ls lists directory contents on a remote xrootd server.

Usage:

 $> xrd-ls [OPTIONS] <dir-1> [<dir-2> [...]]

Example:

 $> xrd-ls root://server.example.com/some/dir
 $> xrd-ls -l root://server.example.com/some/dir
 $> xrd-ls -R root://server.example.com/some/dir
 $> xrd-ls -l -R root://server.example.com/some/dir

Options:
`)
		flag.PrintDefaults()
	}
}

func main() {
	log.SetPrefix("xrd-ls: ")
	log.SetFlags(0)

	var (
		recFlag  = flag.Bool("R", false, "list subdirectories recursively")
		longFlag = flag.Bool("l", false, "use a long listing format")
	)

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		log.Fatalf("missing directory operand")
	}

	for i, dir := range flag.Args() {
		if i > 0 {
			// separate consecutive files by an empty line
			fmt.Printf("\n")
		}
		err := xrdls(dir, *longFlag, *recFlag)
		if err != nil {
			log.Fatalf("could not list %q content: %v", dir, err)
		}
	}
}

func xrdls(name string, long, recursive bool) error {
	url, err := xrdio.Parse(name)
	if err != nil {
		return errors.Errorf("could not parse %q: %v", name, err)
	}

	ctx := context.Background()

	c, err := xrdclient.NewClient(ctx, url.Addr, url.User)
	if err != nil {
		return errors.Errorf("could not create client: %v", err)
	}
	defer c.Close()

	fs := c.FS()

	fi, err := fs.Stat(ctx, url.Path)
	if err != nil {
		return errors.Errorf("could not stat %q: %v", url.Path, err)
	}
	err = display(ctx, fs, url.Path, fi, long, recursive)
	if err != nil {
		return err
	}

	return nil
}

func display(ctx context.Context, fs xrdfs.FileSystem, root string, fi os.FileInfo, long, recursive bool) error {
	if !fi.IsDir() {
		format(os.Stdout, root, fi, long)
		return nil
	}

	end := ""
	if recursive {
		end = ":"
	}

	fmt.Printf("%s%s\n", path.Join(root, fi.Name()), end)
	dir := path.Join(root, fi.Name())
	if long {
		fmt.Printf("total %d\n", fi.Size())
	}
	ents, err := fs.Dirlist(ctx, dir)
	if err != nil {
		return errors.Errorf("could not list dir %q: %v", dir, err)
	}
	o := tabwriter.NewWriter(os.Stdout, 8, 4, 0, ' ', tabwriter.AlignRight)
	for _, e := range ents {
		format(o, root, e, long)
	}
	o.Flush()
	if recursive {
		for _, e := range ents {
			if !e.IsDir() {
				continue
			}
			// make an empty line before going into a subdirectory.
			fmt.Printf("\n")
			err := display(ctx, fs, dir, e, long, recursive)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func format(o io.Writer, root string, fi os.FileInfo, long bool) {
	if !long {
		fmt.Fprintf(o, "%s\n", path.Join(root, fi.Name()))
		return
	}

	fmt.Fprintf(o, "%v\t %d\t %s\t %s\n", fi.Mode(), fi.Size(), fi.ModTime().Format("Jan 02 15:04"), fi.Name())
}
