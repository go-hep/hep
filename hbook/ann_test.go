// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hbook

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAnnotationYODA(t *testing.T) {
	ann1 := make(Annotation)
	ann1["title"] = "my title"
	ann1["name"] = "h1d"
	ann1["meta"] = "data"

	got1, err := ann1.MarshalYODA()
	if err != nil {
		t.Fatalf("could not marshal annotation: %+v", err)
	}

	ann2 := make(Annotation)
	err = ann2.UnmarshalYODA(got1)
	if err != nil {
		t.Fatalf("could not unmarshal annotation: %+v", err)
	}

	if !reflect.DeepEqual(ann1, ann2) {
		t.Fatalf("r/w roundtrip failed:\ngot= %#v\nwant=%#v", ann2, ann1)
	}
}

func TestAnnotationClone(t *testing.T) {
	ann := make(Annotation)
	ann["title"] = "my title"
	ann["name"] = "h1d"
	ann["meta"] = []float64{1.1, 2.1, 3.1, 4.1}

	clo := ann.clone()

	if got, want := clo, ann; !reflect.DeepEqual(got, want) {
		t.Fatalf("cloning failed:\ngot= %#v\nwant=%#v", got, want)
	}

	got, err := clo.MarshalYODA()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	want, err := ann.MarshalYODA()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("ann differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(got),
			),
		)
	}

	// test the 2 values are indeed decoupled.
	delete(clo, "meta")
	clo["data"] = 42.42

	got, err = ann.MarshalYODA()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}

	if !bytes.Equal(got, want) {
		t.Fatalf("ann differ:\n%s\n",
			cmp.Diff(
				string(want),
				string(got),
			),
		)
	}
}
