// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// root2arrow converts the content of a ROOT TTree to an ARROW file.
//
//  Usage of root2arrow:
//    -o string
//      	path to output ARROW file name (default "output.data")
//    -stream
//      	enable ARROW stream (default is to create an ARROW file)
//    -t string
//      	name of the tree to convert (default "tree")
//
//
//  $> root2arrow -o foo.data -t tree ../../groot/testdata/simple.root
//  $> arrow-ls ./foo.data
//  version: V4
//  schema:
//    fields: 3
//      - one: type=int32
//      - two: type=float32
//      - three: type=utf8
//  records: 1
//  $> arrow-cat ./foo.data
//  version: V4
//  record 1/1...
//    col[0] "one": [1 2 3 4]
//    col[1] "two": [1.1 2.2 3.3 4.4]
//    col[2] "three": ["uno" "dos" "tres" "quatro"]
//
package main // import "go-hep.org/x/hep/cmd/root2arrow"

import (
	"flag"
	"io"
	"log"
	"os"

	"github.com/apache/arrow/go/arrow/array"
	"github.com/apache/arrow/go/arrow/ipc"
	"github.com/apache/arrow/go/arrow/memory"
	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rarrow"
	_ "go-hep.org/x/hep/groot/riofs/plugin/http"
	_ "go-hep.org/x/hep/groot/riofs/plugin/xrootd"
	"go-hep.org/x/hep/groot/rtree"
)

func main() {
	log.SetPrefix("root2arrow: ")
	log.SetFlags(0)

	oname := flag.String("o", "output.data", "path to output ARROW file name")
	tname := flag.String("t", "tree", "name of the tree to convert")
	stream := flag.Bool("stream", false, "enable ARROW stream (default is to create an ARROW file)")

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		log.Fatalf("missing input ROOT filename argument")
	}
	fname := flag.Arg(0)

	err := process(*oname, fname, *tname, *stream)
	if err != nil {
		log.Fatal(err)
	}
}

func process(oname, fname, tname string, stream bool) error {
	f, err := groot.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()

	obj, err := f.Get(tname)
	if err != nil {
		return err
	}

	tree, ok := obj.(rtree.Tree)
	if !ok {
		return errors.Errorf("object %q in file %q is not a rtree.Tree", tname, fname)
	}

	mem := memory.NewGoAllocator()

	r := rarrow.NewRecordReader(tree, rarrow.WithAllocator(mem))
	defer r.Release()

	var o *os.File

	switch oname {
	case "":
		o = os.Stdout
	default:
		o, err = os.Create(oname)
		if err != nil {
			return err
		}
		defer o.Close()
	}

	switch {
	case stream:
		err = processStream(o, r, mem)
	default:
		err = processFile(o, r, mem)
	}

	return err
}

func processStream(o io.Writer, r array.RecordReader, mem memory.Allocator) error {
	var err error
	w := ipc.NewWriter(o, ipc.WithSchema(r.Schema()), ipc.WithAllocator(mem))
	defer w.Close()

	i := 0
	for r.Next() {
		rec := r.Record()
		err = w.Write(rec)
		if err != nil {
			return errors.Wrapf(err, "could not write record[%d]", i)
		}
		i++
	}

	err = w.Close()
	if err != nil {
		return errors.Wrap(err, "could not close Arrow stream writer")
	}

	return nil
}

func processFile(o *os.File, r array.RecordReader, mem memory.Allocator) error {
	w, err := ipc.NewFileWriter(o, ipc.WithSchema(r.Schema()), ipc.WithAllocator(mem))
	if err != nil {
		return errors.Wrap(err, "could not create Arrow file writer")
	}
	defer w.Close()

	i := 0
	for r.Next() {
		rec := r.Record()
		err = w.Write(rec)
		if err != nil {
			return errors.Wrapf(err, "could not write record[%d]", i)
		}
		i++
	}

	err = w.Close()
	if err != nil {
		return errors.Wrap(err, "could not close Arrow file writer")
	}

	err = o.Sync()
	if err != nil {
		return errors.Wrap(err, "could not sync data to disk")
	}

	err = o.Close()
	if err != nil {
		return errors.Wrap(err, "could not close output file")
	}

	return nil
}
