// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestDump(t *testing.T) {
	f, err := os.Open("../testdata/small.hepmc")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	out := new(bytes.Buffer)
	err = dump(out, f)
	if err != nil {
		t.Fatal(err)
	}

	want, err := ioutil.ReadFile("testdata/small.hepmc.ref")
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out.Bytes(), want) {
		t.Fatalf("dump error.\ngot:\n%s\nwant:\n%s\n", out.Bytes(), want)
	}
}

func TestDumpFail(t *testing.T) {
	ref, err := ioutil.ReadFile("../testdata/small.hepmc")
	if err != nil {
		t.Fatal(err)
	}

	for _, i := range []int{
		0, 1, 2, 3, 4, 5,
		64, 128, 256,
	} {
		raw := ref[:i]
		t.Run("input", func(t *testing.T) {
			r := bytes.NewReader(raw)
			out := new(bytes.Buffer)
			err := dump(out, r)
			if err == nil {
				t.Fatalf("expected a failure\n%q\n", raw)
			}
		})
	}

	want, err := ioutil.ReadFile("testdata/small.hepmc.ref")
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < len(want)-1; i++ {
		nbytes := 0
		t.Run("output", func(t *testing.T) {
			r := bytes.NewReader(ref)
			out := newWriter(nbytes)
			err := dump(out, r)
			if err == nil {
				t.Fatalf("expected a failure")
			}
		})
	}
}

type failWriter struct {
	n int
	c int
}

func newWriter(n int) *failWriter {
	return &failWriter{n: n, c: 0}
}

func (w *failWriter) Write(data []byte) (int, error) {
	for range data {
		w.c++
		if w.c >= w.n {
			return 0, io.ErrUnexpectedEOF
		}
	}
	return len(data), nil
}
