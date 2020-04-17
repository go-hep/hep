// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// rio-ls-records displays the list of records stored in a given rio file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/rio"
)

func main() {

	log.SetFlags(0)
	log.SetPrefix("rio-ls-records: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: rio-ls-records [options] <file.rio> [<file2.rio> [<file3.rio> [...]]]
	
ex:
 $ rio-ls-records file.rio
 `,
		)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		log.Printf("missing filename argument\n")
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}

	for _, fname := range flag.Args() {
		inspect(fname)
	}
}

func inspect(fname string) {
	log.Printf("inspecting file [%s]...\n", fname)
	if fname == "" {
		flag.Usage()
		flag.PrintDefaults()
		os.Exit(1)
	}

	rtypes := metaData(fname)

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r, err := rio.NewReader(f)
	if err != nil {
		log.Fatalf("error creating rio.Reader: %v\n", err)
	}

	scan := rio.NewScanner(r)
	for scan.Scan() {
		// scans through the whole stream
		err = scan.Err()
		if err != nil {
			break
		}
		rec := scan.Record()
		if rec.Name() == rio.MetaRecord {
			continue
		}
		rtype := rtypes[rec.Name()]
		if rtype != "" {
			rtype = "type=" + rtype
		}
		fmt.Printf(" -> %-20s%s\n", rec.Name(), rtype)
	}
	err = scan.Err()
	if err != nil {
		log.Fatalf("error during file scan: %v\n", err)
	}

	log.Printf("inspecting file [%s]... [done]\n", fname)
}

func metaData(fname string) map[string]string {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var rtypes = make(map[string]string)
	r, err := rio.NewReader(f)
	if err != nil {
		return rtypes
	}

	scan := rio.NewScanner(r)
	scan.Select([]rio.Selector{{Name: rio.MetaRecord, Unpack: true}})
	if !scan.Scan() {
		log.Fatal(scan.Err())
	}
	rec := scan.Record()
	if rec == nil {
		return rtypes
	}

	blk := rec.Block(rio.MetaRecord)
	if blk == nil {
		return rtypes
	}

	var meta rio.Metadata
	err = blk.Read(&meta)
	if err != nil {
		return rtypes
	}

	for _, mrec := range meta.Records {
		mblk := mrec.Blocks[0]
		rtypes[mblk.Name] = mblk.Type
	}

	return rtypes
}
