// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Command xrd-cp copies files and directories from a remote xrootd server
// to local storage.
//
// Usage:
//
//  $> xrd-cp [OPTIONS] <src-1> [<src-2> [...]] <dst>
//
// Example:
//
//  $> xrd-cp root://server.example.com/some/file1.txt .
//  $> xrd-cp root://gopher@server.example.com/some/file1.txt .
//  $> xrd-cp root://server.example.com/some/file1.txt foo.txt
//  $> xrd-cp root://server.example.com/some/file1.txt - > foo.txt
//  $> xrd-cp -r root://server.example.com/some/dir .
//  $> xrd-cp -r root://server.example.com/some/dir outdir
//
// Options:
//   -r	copy directories recursively
//   -v	enable verbose mode
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	stdpath "path"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/xrootd"
	"go-hep.org/x/hep/xrootd/xrdfs"
	"go-hep.org/x/hep/xrootd/xrdio"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `xrd-cp copies files and directories from a remote xrootd server to local storage.

Usage:

 $> xrd-cp [OPTIONS] <src-1> [<src-2> [...]] <dst>

Example:

 $> xrd-cp root://server.example.com/some/file1.txt .
 $> xrd-cp root://gopher@server.example.com/some/file1.txt .
 $> xrd-cp root://server.example.com/some/file1.txt foo.txt
 $> xrd-cp root://server.example.com/some/file1.txt - > foo.txt
 $> xrd-cp -r root://server.example.com/some/dir .
 $> xrd-cp -r root://server.example.com/some/dir outdir

Options:
`)
		flag.PrintDefaults()
	}
}

func main() {
	log.SetPrefix("xrd-cp: ")
	log.SetFlags(0)

	var (
		recFlag     = flag.Bool("r", false, "copy directories recursively")
		verboseFlag = flag.Bool("v", false, "enable verbose mode")
	)

	flag.Parse()

	switch n := flag.NArg(); n {
	case 0:
		flag.Usage()
		log.Fatalf("missing file operand")
	case 1:
		flag.Usage()
		log.Fatalf("missing destination file operand after %q", flag.Arg(0))
	case 2:
		err := xrdcopy(flag.Arg(1), flag.Arg(0), *recFlag, *verboseFlag)
		if err != nil {
			log.Fatalf("could not copy %q to %q: %v", flag.Arg(0), flag.Arg(1), err)
		}
	default:
		dst := flag.Arg(flag.NArg() - 1)
		for _, src := range flag.Args()[:flag.NArg()-1] {
			err := xrdcopy(dst, src, *recFlag, *verboseFlag)
			if err != nil {
				log.Fatalf("could not copy %q to %q: %v", src, dst, err)
			}
		}
	}
}

func xrdcopy(dst, srcPath string, recursive, verbose bool) error {
	cli, src, err := xrdremote(srcPath)
	if err != nil {
		return err
	}
	defer cli.Close()

	ctx := context.Background()

	fs := cli.FS()
	var jobs jobs
	var addDir func(root, src string) error

	addDir = func(root, src string) error {
		fi, err := fs.Stat(ctx, src)
		if err != nil {
			return errors.WithMessage(err, "could not stat remote src")
		}
		switch {
		case fi.IsDir():
			if !recursive {
				return errors.Errorf("xrd-cp: -r not specified; omitting directory %q", src)
			}
			dst := stdpath.Join(root, stdpath.Base(src))
			err = os.MkdirAll(dst, 0755)
			if err != nil {
				return errors.WithMessage(err, "could not create output directory")
			}

			ents, err := fs.Dirlist(ctx, src)
			if err != nil {
				return errors.WithMessage(err, "could not list directory")
			}
			for _, e := range ents {
				err = addDir(dst, stdpath.Join(src, e.Name()))
				if err != nil {
					return err
				}
			}
		default:
			jobs.add(job{
				fs:  fs,
				src: src,
				dst: stdpath.Join(root, stdpath.Base(src)),
			})
		}
		return nil
	}

	fiSrc, err := fs.Stat(ctx, src)
	if err != nil {
		return errors.WithMessage(err, "could not stat remote src")
	}

	fiDst, errDst := os.Stat(dst)
	switch {
	case fiSrc.IsDir():
		switch {
		case errDst != nil && os.IsNotExist(errDst):
			err = os.MkdirAll(dst, 0755)
			if err != nil {
				return errors.WithMessage(err, "could not create output directory")
			}
			ents, err := fs.Dirlist(ctx, src)
			if err != nil {
				return errors.WithMessage(err, "could not list directory")
			}
			for _, e := range ents {
				err = addDir(dst, stdpath.Join(src, e.Name()))
				if err != nil {
					return err
				}
			}

		case errDst != nil:
			return errors.WithMessage(errDst, "could not stat local dst")
		case fiDst.IsDir():
			err = addDir(dst, src)
			if err != nil {
				return err
			}
		}

	default:
		switch {
		case errDst != nil && os.IsNotExist(errDst):
			// ok... dst will be the output file.
		case errDst != nil:
			return errors.WithMessage(errDst, "could not stat local dst")
		case fiDst.IsDir():
			dst = stdpath.Join(dst, stdpath.Base(src))
		}

		jobs.add(job{
			fs:  fs,
			src: src,
			dst: dst,
		})
	}

	n, err := jobs.run(ctx)
	if verbose {
		log.Printf("transferred %d bytes", n)
	}
	return err
}

func xrdremote(name string) (client *xrootd.Client, path string, err error) {
	url, err := xrdio.Parse(name)
	if err != nil {
		return nil, "", errors.WithStack(err)
	}

	path = url.Path
	client, err = xrootd.NewClient(context.Background(), url.Addr, url.User)
	return client, path, err
}

type job struct {
	fs  xrdfs.FileSystem
	src string
	dst string
}

func (j job) run(ctx context.Context) (int, error) {
	var (
		o   io.WriteCloser
		err error
	)
	switch j.dst {
	case "-", "":
		o = os.Stdout
	case ".":
		j.dst = stdpath.Base(j.src)
		fallthrough
	default:
		o, err = os.Create(j.dst)
		if err != nil {
			return 0, errors.WithMessage(err, "could not create output file")
		}
	}
	defer o.Close()

	f, err := xrdio.OpenFrom(j.fs, j.src)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// TODO(sbinet): make buffer a field of job to reduce memory pressure.
	// TODO(sbinet): use clever heuristics for buffer size?
	n, err := io.CopyBuffer(o, f, make([]byte, 16*1024*1024))
	if err != nil {
		return int(n), errors.WithMessage(err, "could not copy to output file")
	}

	err = o.Close()
	if err != nil {
		return int(n), errors.WithMessage(err, "could not close output file")
	}

	return int(n), nil
}

type jobs struct {
	slice []job
}

func (js *jobs) add(j job) {
	js.slice = append(js.slice, j)
}

func (js *jobs) run(ctx context.Context) (int, error) {
	var n int
	for _, j := range js.slice {
		nn, err := j.run(ctx)
		n += nn
		if err != nil {
			return n, err
		}
	}
	return n, nil
}
