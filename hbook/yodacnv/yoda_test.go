// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package yodacnv_test

import (
	"bytes"
	"reflect"
	"testing"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/yodacnv"
)

var (
	rdata []byte
	h1    *hbook.H1D
	h2    *hbook.H2D
	p1    *hbook.P1D
	s2    *hbook.S2D
)

func TestReadWrite(t *testing.T) {
	r := bytes.NewReader(rdata)
	objs, err := yodacnv.Read(r)
	if err != nil {
		t.Fatal(err)
	}

	w := new(bytes.Buffer)
	for _, v := range objs {
		err = yodacnv.Write(w, v.(yodacnv.Marshaler))
		if err != nil {
			t.Fatal(err)
		}
	}

	if !reflect.DeepEqual(w.Bytes(), rdata) {
		t.Fatalf("got:\n%s\nwant:\n%s\n", string(w.Bytes()), string(rdata))
	}
}

func init() {

	add := func(o yodacnv.Marshaler) {
		raw, err := o.MarshalYODA()
		if err != nil {
			panic(err)
		}
		rdata = append(rdata, raw...)
	}

	h1 = hbook.NewH1D(10, -4, 4)
	h1.Annotation()["name"] = "histo-1d"
	h1.Fill(1, 1)
	h1.Fill(2, 1)
	h1.Fill(-3, 1)
	h1.Fill(-4, 1)
	h1.Fill(0, 1)
	h1.Fill(0, 1)
	h1.Fill(10, 1)
	h1.Fill(-10, 1)

	add(h1)

	h2 = hbook.NewH2D(5, -1, 1, 5, -2, +2)
	h2.Annotation()["name"] = "histo-2d"
	h2.Fill(+0.5, +1, 1)
	h2.Fill(-0.5, +1, 1)
	h2.Fill(+0.0, -1, 1)

	add(h2)

	p1 = hbook.NewP1D(10, -4, +4)
	for i := 0; i < 10; i++ {
		v := float64(i)
		p1.Fill(v, v*2, 1)
	}
	p1.Fill(-10, 10, 1)

	add(p1)

	s2 = hbook.NewS2DFromH1D(h1)
	add(s2)
}
