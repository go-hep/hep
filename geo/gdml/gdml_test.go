// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdml

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"reflect"
	"testing"
)

func TestReadSchema(t *testing.T) {
	t.Skip("boo")

	for _, tc := range []struct {
		fname string
		err   error
		want  Schema
	}{
		{
			fname: "testdata/test.gdml",
		},
	} {
		t.Run(tc.fname, func(t *testing.T) {
			raw, err := ioutil.ReadFile(tc.fname)
			if err != nil {
				t.Fatal(err)
			}
			var schema Schema
			dec := xml.NewDecoder(bytes.NewReader(raw))
			err = dec.Decode(&schema)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(schema, tc.want) {
				t.Fatal(err)
			}
		})
	}
}

func TestReadConstant(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Constant
	}{
		{
			raw:  `<constant name="length" value="6.25"/>`,
			want: Constant{Name: "length", Value: 6.25},
		},
		{
			raw:  `<constant name="mass" value="6"/>`,
			want: Constant{Name: "mass", Value: 6},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Constant
			err := xml.NewDecoder(bytes.NewReader([]byte(tc.raw))).Decode(&v)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(v, tc.want) {
				t.Fatalf("got = %#v\nwant= %#v", v, tc.want)
			}
		})
	}
}

func TestReadQuantity(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Quantity
	}{
		{
			raw:  `<quantity  name="W_Density" type="density" value="1" unit="g/cm3"/>`,
			want: Quantity{Name: "W_Density", Type: "density", Value: 1, Unit: "g/cm3"},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Quantity
			err := xml.NewDecoder(bytes.NewReader([]byte(tc.raw))).Decode(&v)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(v, tc.want) {
				t.Fatalf("got = %#v\nwant= %#v", v, tc.want)
			}
		})
	}
}

func TestReadVariable(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Variable
	}{
		{
			raw:  `<variable name="x" value="6"/>`,
			want: Variable{Name: "x", Value: "6"},
		},
		{
			raw:  `<variable name="y" value="x/2"/>`,
			want: Variable{Name: "y", Value: "x/2"},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Variable
			err := xml.NewDecoder(bytes.NewReader([]byte(tc.raw))).Decode(&v)
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(v, tc.want) {
				t.Fatalf("got = %#v\nwant= %#v", v, tc.want)
			}
		})
	}
}
