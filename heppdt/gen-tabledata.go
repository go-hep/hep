// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build ignore

// gen-tabledata generates tabledata.go
package main

import (
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Create("tabledata.go")
	if err != nil {
		log.Fatalf(
			"could not create 'tabledata.go': %v\n",
			err,
		)
	}
	defer f.Close()

	for _, fname := range []string{
		"tabledata.header",
		"tabledata.tbl",
		"tabledata.footer",
	} {
		src, err := os.Open(fname)
		if err != nil {
			log.Fatalf(
				"could not open data file [%s]: %v\n",
				fname,
				err,
			)
		}
		defer src.Close()
		_, err = io.Copy(f, src)
		if err != nil {
			log.Fatalf(
				"error copying content of [%s]: %v\n",
				fname,
				err,
			)
		}
		err = src.Close()
		if err != nil {
			log.Fatalf(
				"error closing file [%s]: %v\n",
				fname,
				err,
			)
		}
	}
	err = f.Close()
	if err != nil {
		log.Fatalf("error closing [%s]: %v\n",
			f.Name(),
			err,
		)
	}
}
