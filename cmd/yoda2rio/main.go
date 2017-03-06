// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// yoda2rio converts YODA files containing hbook-like values (H1D, H2D, P1D, ...)
// into rio files.
//
// Example:
//
//  $> yoda2rio rivet.yoda >| rivet.rio
package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/rio"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("yoda2rio: ")
	log.SetOutput(os.Stderr)

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: yoda2rio [options] <file1.yoda> [<file2.yoda> [...]]

ex:
 $ yoda2rio rivet.yoda >| rivet.rio
`)
	}

	flag.Parse()

	if flag.NArg() < 1 {
		log.Printf("missing input file name")
		flag.Usage()
		flag.PrintDefaults()
	}

	o, err := rio.NewWriter(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
	defer o.Close()

	for _, fname := range flag.Args() {
		convert(o, fname)
	}
}

func convert(w *rio.Writer, fname string) {
	r, err := os.Open(fname)
	if err != nil {
		log.Fatalf("error opening file [%s]: %v\n", fname, err)
	}
	defer r.Close()

	vs, err := yodaSlice(r)
	if err != nil {
		log.Fatalf("error decoding YODA file [%s]: %v\n", fname, err)
	}

	for _, v := range vs {
		err = w.WriteValue(v.Name(), v.Value())
		if err != nil {
			log.Fatalf("error writing %q from YODA file [%s]: %v\n", v.Name(), fname, err)
		}
	}
}

var (
	yodaHeader = []byte("BEGIN YODA_")
	yodaFooter = []byte("END YODA_")
)

func yodaSlice(r io.Reader) ([]Yoda, error) {
	var (
		err   error
		o     []Yoda
		block = make([]byte, 0, 1024)
		rt    reflect.Type
		name  string
	)
	scan := bufio.NewScanner(r)
	for scan.Scan() {
		raw := scan.Bytes()
		switch {
		case bytes.HasPrefix(raw, yodaHeader):
			rt, name, err = splitHeader(raw)
			if err != nil {
				log.Fatalf("error parsing YODA header: %v", err)
			}
			block = block[:0]
			block = append(block, raw...)
			block = append(block, '\n')

		default:
			block = append(block, raw...)
			block = append(block, '\n')

		case bytes.HasPrefix(raw, yodaFooter):
			block = append(block, raw...)
			block = append(block, '\n')

			v := reflect.New(rt).Elem()
			err = v.Addr().Interface().(unmarshalYoda).UnmarshalYODA(block)
			if err != nil {
				log.Fatalf("error unmarshaling YODA %q (type=%v): %v\n%v\n===\n", name, rt.Name(), err, string(block))
			}
			o = append(o, &yoda{name: name, ptr: v.Addr().Interface()})
		}
	}
	err = scan.Err()
	if err != nil {
		return nil, err
	}
	return o, nil
}

func splitHeader(raw []byte) (reflect.Type, string, error) {
	raw = raw[len(yodaHeader):]
	i := bytes.Index(raw, []byte(" "))
	if i == -1 || i >= len(raw) {
		return nil, "", fmt.Errorf("invalid YODA header (missing space)")
	}

	var rt reflect.Type

	switch string(raw[:i]) {
	case "HISTO1D":
		rt = reflect.TypeOf((*hbook.H1D)(nil)).Elem()
	case "HISTO2D":
		rt = reflect.TypeOf((*hbook.H2D)(nil)).Elem()
	case "PROFILE1D":
		rt = reflect.TypeOf((*hbook.P1D)(nil)).Elem()
	case "SCATTER2D":
		rt = reflect.TypeOf((*hbook.S2D)(nil)).Elem()
	default:
		log.Fatalf("unhandled YODA object type %q", string(raw[:i]))
	}

	name := raw[i+2:] // +2 to also remove the leading '/'
	return rt, strings.TrimSpace(string(name)), nil
}

type Yoda interface {
	Name() string
	Value() interface{}
}

type unmarshalYoda interface {
	UnmarshalYODA([]byte) error
}

type yoda struct {
	name string
	ptr  interface{}
}

func (y *yoda) Name() string {
	return y.name
}

func (y *yoda) Value() interface{} {
	return y.ptr
}
