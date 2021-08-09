// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdml

import (
	"bytes"
	"encoding/xml"
	"os"
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
			raw, err := os.ReadFile(tc.fname)
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

func TestReadPosition(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Position
	}{
		{
			raw:  `<position name="box_position" x="25.0" y="50.0" z="75.0" unit="cm"/>`,
			want: Position{Name: "box_position", X: "25.0", Y: "50.0", Z: "75.0", Unit: "cm"},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Position
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

func TestReadRotation(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Rotation
	}{
		{
			raw:  `<rotation name="identity" x="0" y="0" z="0" unit="deg"/>`,
			want: Rotation{Name: "identity", X: "0", Y: "0", Z: "0", Unit: "deg"},
		},
		{
			raw:  `<rotation name="rot-z" x="0" y="0" z="30" unit="deg"/>`,
			want: Rotation{Name: "rot-z", X: "0", Y: "0", Z: "30", Unit: "deg"},
		},
		{
			raw:  `<rotation name="rot-x" x="10" unit="deg"/>`,
			want: Rotation{Name: "rot-x", X: "10", Y: "", Z: "", Unit: "deg"},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Rotation
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

func TestReadScale(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Scale
	}{
		{
			raw:  `<scale name="identity" x="1" y="1" z="1"/>`,
			want: Scale{Name: "identity", X: "1", Y: "1", Z: "1"},
		},
		{
			raw:  `<scale name="reflection" x="-1" y="-1" z="1"/>`,
			want: Scale{Name: "reflection", X: "-1", Y: "-1", Z: "1"},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Scale
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

func TestReadMatrix(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Matrix
	}{
		{
			raw:  `<matrix name="m1" coldim="3" values="0.4 9 126 8.5 7 21 34.6 7 9"/>`,
			want: Matrix{Name: "m1", Cols: 3, Values: "0.4 9 126 8.5 7 21 34.6 7 9"},
		},
		{
			// FIXME(sbinet): customize Matrix.UnmarshalXML to post-process lines into []float64 or []string
			raw: `<matrix name="m2" coldim="3" values=" 0.4 9 126
 8.5 7 21
34.6 7  9"/>`,
			want: Matrix{Name: "m2", Cols: 3, Values: " 0.4 9 126\n 8.5 7 21\n34.6 7  9"},
		},
		{
			raw:  `<matrix name="m3" coldim="1" values="1 2 3 4"/>`,
			want: Matrix{Name: "m3", Cols: 1, Values: "1 2 3 4"},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Matrix
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

func TestReadIsotope(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Isotope
	}{
		{
			raw: `
			<isotope name="U235" Z="92" N="235">
				<atom type="A" value="235.01"/>
			</isotope>`,
			want: Isotope{Name: "U235", Z: 92, N: 235, Atom: Atom{Type: "A", Value: 235.01}},
		},
		{
			raw: `
			<isotope name="U238" Z="92" N="238">
				<atom type="A" value="235.03"/>
			</isotope>`,
			want: Isotope{Name: "U238", Z: 92, N: 238, Atom: Atom{Type: "A", Value: 235.03}},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Isotope
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

func TestReadElement(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Element
	}{
		{
			raw: `
			<element name="Oxygen" Z="8" formula="O">
				<atom value="16"/>
			</element>`,
			want: Element{Name: "Oxygen", Z: 8, Formula: "O", Atom: Atom{Type: "", Value: 16}},
		},
		{
			raw: `
			<element name="enriched_uranium">
				<fraction ref="U235" n="0.9" />
				<fraction ref="U238" n="0.1" />
			</element>`,
			want: Element{
				Name: "enriched_uranium",
				Fractions: []Fraction{
					{Ref: "U235", N: 0.9},
					{Ref: "U238", N: 0.1},
				},
			},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Element
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

func TestReadMaterial(t *testing.T) {
	for _, tc := range []struct {
		raw  string
		want Material
	}{
		{
			raw: `
			<material name="Al" Z="13.0">
				<D    value="2.70"/>
				<atom value="26.98"/>
			</material>`,
			want: Material{Name: "Al", Z: 13.0, Density: Density{Value: 2.7}, Atom: Atom{Value: 26.98}},
		},
		{
			raw: `
			<material name="Water" formula="H2O">
				<D    value="1.0"/>
				<composite n="2" ref="Hydrogen"/>
				<composite n="1" ref="Oxygen"/>
			</material>`,
			want: Material{
				Name:    "Water",
				Formula: "H2O",
				Density: Density{Value: 1},
				Composites: []Composite{
					{N: 2, Ref: "Hydrogen"},
					{N: 1, Ref: "Oxygen"},
				},
			},
		},
		{
			raw: `
			<material name="Air" formula="air">
				<D    value="0.0012899999999999999"/>
				<fraction n="0.7" ref="Nitrogen"/>
				<fraction n="0.3" ref="Oxygen"/>
			</material>`,
			want: Material{
				Name:    "Air",
				Formula: "air",
				Density: Density{Value: 0.0012899999999999999},
				Fractions: []Fraction{
					{N: 0.7, Ref: "Nitrogen"},
					{N: 0.3, Ref: "Oxygen"},
				},
			},
		},
	} {
		t.Run(tc.want.Name, func(t *testing.T) {
			var v Material
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
