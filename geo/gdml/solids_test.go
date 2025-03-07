// Copyright ©2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gdml

import (
	"bytes"
	"encoding/xml"
	"reflect"
	"testing"
)

func TestReadSolids(t *testing.T) {
	for _, tc := range []struct {
		name string
		raw  string
		want any
	}{
		{
			name: "box",
			raw:  `<box name="mybox" x="30" y="32" z="33" lunit="mm"/>`,
			want: Box{Name: "mybox", X: "30", Y: "32", Z: "33", LUnit: "mm"},
		},
		{
			name: "cone",
			raw:  `<cone name="mycone" rmin1="10" rmax1="15" rmin2="20" rmax2="25" z="30" startphi="1" deltaphi="4" aunit="rad" lunit="mm"/>`,
			want: Cone{
				Name:  "mycone",
				RMin1: "10", RMax1: "15", RMin2: "20", RMax2: "25", Z: "30", StartPhi: "1", DPhi: "4",
				AUnit: "rad",
				LUnit: "mm",
			},
		},
		{
			name: "ellipsoid",
			raw:  `<ellipsoid name="myellipsoid" ax="10" by="15" cz="20" zcut1="2" zcut2="4" lunit="mm"/>`,
			want: Ellipsoid{Name: "myellipsoid", Ax: "10", By: "15", Cz: "20", ZCut1: "2", ZCut2: "4", LUnit: "mm"},
		},
		{
			name: "elliptical-tube",
			raw:  `<eltube name="myeltube" dx="10" dy="15" dz="20" lunit="mm"/>`,
			want: EllipticalTube{Name: "myeltube", Dx: "10", Dy: "15", Dz: "20", LUnit: "mm"},
		},
		{
			name: "elliptical-cone",
			raw:  `<elcone name="myelcone" dx="10" dy="15" zmax="2" zcut="1.5" lunit="mm"/>`,
			want: EllipticalCone{Name: "myelcone", Dx: "10", Dy: "15", ZMax: "2", ZCut: "1.5", LUnit: "mm"},
		},
		{
			name: "orb",
			raw:  `<orb name="myorb" r="10" lunit="mm"/>`,
			want: Orb{Name: "myorb", R: "10", LUnit: "mm"},
		},
		{
			name: "paraboloid",
			raw:  `<paraboloid name="parab" rlo="10" rhi="15" dz="20" aunit="rad" lunit="mm"/>`,
			want: Paraboloid{Name: "parab", Rlo: "10", Rhi: "15", Dz: "20", AUnit: "rad", LUnit: "mm"},
		},
		{
			name: "parallelepiped",
			raw:  `<para name="para" x="10" y="11" z="12" alpha="1" theta="2" phi="3" aunit="rad" lunit="mm"/>`,
			want: Parallelepiped{Name: "para", X: "10", Y: "11", Z: "12", Alpha: "1", Theta: "2", Phi: "3", AUnit: "rad", LUnit: "mm"},
		},
		{
			name: "polycone",
			raw: `
			<polycone name="polycone" startphi="1" deltaphi="4" aunit="rad" lunit="mm">
				<zplane rmin="1" rmax="9" z="10" />
				<zplane rmin="3" rmax="5" z="12" />
			</polycone>`,
			want: PolyCone{
				Name:     "polycone",
				StartPhi: "1", DPhi: "4", AUnit: "rad", LUnit: "mm",
				ZPlanes: []ZPlane{
					{XMLName: xml.Name{Local: "zplane"}, RMin: "1", RMax: "9", Z: "10"},
					{XMLName: xml.Name{Local: "zplane"}, RMin: "3", RMax: "5", Z: "12"},
				},
			},
		},
		{
			name: "generic-polycone",
			raw: `
			<genericPolycone name="polycone" startphi="1" deltaphi="4" aunit="rad" lunit="mm">
				<rzpoint r="1" z="5" />
				<rzpoint r="3" z="10" />
				<rzpoint r="1" z="12" />
			</genericPolycone>`,
			want: GenericPolyCone{
				Name:     "polycone",
				StartPhi: "1", DPhi: "4", AUnit: "rad", LUnit: "mm",
				RZPoints: []RZPoint{
					{XMLName: xml.Name{Local: "rzpoint"}, R: "1", Z: "5"},
					{XMLName: xml.Name{Local: "rzpoint"}, R: "3", Z: "10"},
					{XMLName: xml.Name{Local: "rzpoint"}, R: "1", Z: "12"},
				},
			},
		},
		{
			name: "polyhedron",
			raw: `
			<polyhedra name="polyhedra" startphi="1" deltaphi="4" numsides="10" aunit="rad" lunit="mm">
				<zplane rmin="1" rmax="9" z="10" />
				<zplane rmin="3" rmax="5" z="12" />
			</polyhedra>`,
			want: PolyHedra{
				Name:     "polyhedra",
				StartPhi: "1", DPhi: "4", NumSides: "10", AUnit: "rad", LUnit: "mm",
				ZPlanes: []ZPlane{
					{XMLName: xml.Name{Local: "zplane"}, RMin: "1", RMax: "9", Z: "10"},
					{XMLName: xml.Name{Local: "zplane"}, RMin: "3", RMax: "5", Z: "12"},
				},
			},
		},
		{
			name: "generic-polyhedra",
			raw: `
			<genericPolyhedra name="polyhedra" startphi="1" deltaphi="4" numsides="10" aunit="rad" lunit="mm">
				<rzpoint r="1" z="5" />
				<rzpoint r="3" z="10" />
				<rzpoint r="1" z="12" />
			</genericPolyhedra>`,
			want: GenericPolyHedra{
				Name:     "polyhedra",
				StartPhi: "1", DPhi: "4", NumSides: "10", AUnit: "rad", LUnit: "mm",
				RZPoints: []RZPoint{
					{XMLName: xml.Name{Local: "rzpoint"}, R: "1", Z: "5"},
					{XMLName: xml.Name{Local: "rzpoint"}, R: "3", Z: "10"},
					{XMLName: xml.Name{Local: "rzpoint"}, R: "1", Z: "12"},
				},
			},
		},
		{
			name: "sphere",
			raw:  `<sphere name="sphere" rmin="1" rmax="4" deltaphi="1" deltatheta="2" aunit="rad" lunit="mm"/>`,
			want: Sphere{Name: "sphere", RMin: "1", RMax: "4", DPhi: "1", DTheta: "2", AUnit: "rad", LUnit: "mm"},
		},
		{
			name: "torus",
			raw:  `<torus name="torus" rmin="1" rmax="4" rtor="2" deltaphi="3" startphi="1" aunit="rad" lunit="mm"/>`,
			want: Torus{Name: "torus", RMin: "1", RMax: "4", Rtor: "2", DPhi: "3", StartPhi: "1", AUnit: "rad", LUnit: "mm"},
		},
		{
			name: "trapezoid",
			raw:  `<trd name="trd" x1="9" x2="8" y1="6" y2="5" z="10" lunit="mm"/>`,
			want: Trapezoid{Name: "trd", X1: "9", X2: "8", Y1: "6", Y2: "5", Z: "10", LUnit: "mm"},
		},
		{
			name: "trap",
			raw: `<trap name="trap" 
				z="10" theta="1" phi="2" y1="15" x1="10" x2="10.1" alpha1="1.1"
				y2="15.1" x3="10.3" x4="10.4" alpha2="1.2"
				aunit="rad" lunit="mm"
			/>`,
			want: GeneralTrapezoid{
				Name: "trap", AUnit: "rad", LUnit: "mm",
				Z: "10", Theta: "1", Phi: "2",
				Y1: "15", X1: "10",
				X2:     "10.1",
				Alpha1: "1.1",
				Y2:     "15.1", X3: "10.3", X4: "10.4",
				Alpha2: "1.2",
			},
		},
		{
			name: "hype",
			raw:  `<hype name="hype" rmin="1" rmax="2" z="20" inst="3" outst="4" lunit="mm"/>`,
			want: HyperbolicTube{Name: "hype", RMin: "1", RMax: "2", Z: "20", InStereo: "3", OutStereo: "4", LUnit: "mm"},
		},
		{
			name: "cuttube",
			raw: `<cutTube name="cuttube" z="20" rmin="1" rmax="5" startphi="1" deltaphi="4"
				lowX="15" lowY="16" lowZ="17"
				highX="20" highY="21" highZ="22"
				aunit="rad" lunit="mm"
			/>`,
			want: CutTube{
				Name: "cuttube", AUnit: "rad", LUnit: "mm",
				Z: "20", RMin: "1", RMax: "5",
				StartPhi: "1", DPhi: "4",
				LowX: "15", LowY: "16", LowZ: "17",
				HighX: "20", HighY: "21", HighZ: "22",
			},
		},
		{
			name: "tube",
			raw:  `<tube name="tube" z="20" rmin="1" rmax="5" startphi="1" deltaphi="4" aunit="rad" lunit="mm" />`,
			want: Tube{
				Name: "tube", AUnit: "rad", LUnit: "mm",
				Z: "20", RMin: "1", RMax: "5",
				StartPhi: "1", DPhi: "4",
			},
		},
		{
			name: "twistedbox",
			raw:  `<twistedbox name="twistbox" PhiTwist="1" x="30" y="32" z="33" aunit="rad" lunit="mm"/>`,
			want: TwistedBox{Name: "twistbox", PhiTwist: "1", X: "30", Y: "32", Z: "33", AUnit: "rad", LUnit: "mm"},
		},
		{
			name: "twistedtrapezoid",
			raw:  `<twistedtrd name="twistedtrd" PhiTwist="1" x1="9" x2="8" y1="6" y2="5" z="10" aunit="rad" lunit="mm"/>`,
			want: TwistedTrapezoid{Name: "twistedtrd", PhiTwist: "1", X1: "9", X2: "8", Y1: "6", Y2: "5", Z: "10", AUnit: "rad", LUnit: "mm"},
		},
		{
			name: "twistedtrap",
			raw: `<twistedtrap name="twistedtrap"
				PhiTwist="1"
				z="10" theta="1" phi="2" y1="15" x1="10" x2="10.1"
				y2="15.1" x3="10.3" x4="10.4" Alph="1.2"
				aunit="rad" lunit="mm"
			/>`,
			want: TwistedGeneralTrapezoid{
				Name: "twistedtrap", AUnit: "rad", LUnit: "mm", PhiTwist: "1",
				Z: "10", Theta: "1", Phi: "2",
				Y1: "15", X1: "10",
				X2: "10.1",
				Y2: "15.1", X3: "10.3", X4: "10.4",
				Alpha: "1.2",
			},
		},
		{
			name: "twistedtube",
			raw: `<twistedtube name="twisttube" endinnerrad="1" endouterrad="2" zlen="3"
				twistedangle="4"
				phi="5"
				midinnerrad="6" midouterrad="7"
				nseg="8" totphi="9"
				aunit="rad" lunit="mm"
			/>`,
			want: TwistedTube{
				Name: "twisttube", AUnit: "rad", LUnit: "mm",
				EndInnerRad: "1", EndOuterRad: "2",
				ZLen:         "3",
				TwistedAngle: "4",
				Phi:          "5",
				MidInnerRad:  "6", MidOuterRad: "7",
				NSeg: "8", TotPhi: "9",
			},
		},
		{
			name: "xtru",
			raw: `<xtru name="xtru" lunit="mm" >
				<twoDimVertex x="3" y="9" />
				<twoDimVertex x="1" y="5" />
				<twoDimVertex x="2" y="4" />
				<section zOrder="1" zPosition="2" xOffset="5" yOffset="3" scalingFactor="3" />
				<section zOrder="2" zPosition="5" xOffset="3" yOffset="5" scalingFactor="1" />
			</xtru>`,
			want: Extruded{
				Name: "xtru", LUnit: "mm",
				Vertices: []Vertex2D{
					{XMLName: xml.Name{Local: "twoDimVertex"}, X: "3", Y: "9"},
					{XMLName: xml.Name{Local: "twoDimVertex"}, X: "1", Y: "5"},
					{XMLName: xml.Name{Local: "twoDimVertex"}, X: "2", Y: "4"},
				},
				Sections: []Section{
					{XMLName: xml.Name{Local: "section"}, ZOrder: "1", ZPos: "2", XOff: "5", YOff: "3", Fact: "3"},
					{XMLName: xml.Name{Local: "section"}, ZOrder: "2", ZPos: "5", XOff: "3", YOff: "5", Fact: "1"},
				},
			},
		},
		{
			name: "arb8",
			raw: `<arb8 name="arb8" lunit="mm"
				v1x="11" v1y="12"
				v2x="13" v2y="14"
				v3x="15" v3y="16"
				v4x="17" v4y="18"
				v5x="19" v5y="20"
				v6x="21" v6y="22"
				v7x="23" v7y="24"
				v8x="25" v8y="26"
				dz="27"
			/>`,
			want: ArbitraryTrapezoid{
				Name: "arb8", LUnit: "mm",
				V1x: "11", V1y: "12",
				V2x: "13", V2y: "14",
				V3x: "15", V3y: "16",
				V4x: "17", V4y: "18",
				V5x: "19", V5y: "20",
				V6x: "21", V6y: "22",
				V7x: "23", V7y: "24",
				V8x: "25", V8y: "26",
				Dz: "27",
			},
		},
		{
			name: "tesselated",
			raw: `<tesselated name="pyramid">
				<triangular vertex1="v1" vertex2="v2" vertex3="v6" type="ABSOLUTE" />
				<triangular vertex1="v2" vertex2="v3" vertex3="v6" type="ABSOLUTE" />
				<triangular vertex1="v3" vertex2="v4" vertex3="v5" type="ABSOLUTE" />
				<triangular vertex1="v4" vertex2="v1" vertex3="v5" type="ABSOLUTE" />
				<triangular vertex1="v1" vertex2="v6" vertex3="v5" type="ABSOLUTE" />
				<triangular vertex1="v6" vertex2="v3" vertex3="v5" type="ABSOLUTE" />
				<quadrangular vertex1="v4" vertex2="v3" vertex3="v2" vertex4="v1" type="ABSOLUTE" />
			</tesselated>`,
			want: Tesselated{
				Name: "pyramid",
				Tris: []Triangular{
					{XMLName: xml.Name{Local: "triangular"}, Vtx1: "v1", Vtx2: "v2", Vtx3: "v6", Type: "ABSOLUTE"},
					{XMLName: xml.Name{Local: "triangular"}, Vtx1: "v2", Vtx2: "v3", Vtx3: "v6", Type: "ABSOLUTE"},
					{XMLName: xml.Name{Local: "triangular"}, Vtx1: "v3", Vtx2: "v4", Vtx3: "v5", Type: "ABSOLUTE"},
					{XMLName: xml.Name{Local: "triangular"}, Vtx1: "v4", Vtx2: "v1", Vtx3: "v5", Type: "ABSOLUTE"},
					{XMLName: xml.Name{Local: "triangular"}, Vtx1: "v1", Vtx2: "v6", Vtx3: "v5", Type: "ABSOLUTE"},
					{XMLName: xml.Name{Local: "triangular"}, Vtx1: "v6", Vtx2: "v3", Vtx3: "v5", Type: "ABSOLUTE"},
				},
				Quads: []Quadrangular{
					{XMLName: xml.Name{Local: "quadrangular"}, Vtx1: "v4", Vtx2: "v3", Vtx3: "v2", Vtx4: "v1", Type: "ABSOLUTE"},
				},
			},
		},
		{
			name: "tetrahedron",
			raw:  `<tet name="halfpyramid" vertex1="v1" vertex2="v2" vertex3="v3" vertex4="v4" />`,
			want: TetraHedron{
				Name: "halfpyramid",
				Vtx1: "v1", Vtx2: "v2", Vtx3: "v3", Vtx4: "v4",
			},
		},
		{
			name: "scaled-tube",
			raw: `<scaledSolid name="scaled-tube">
				<solidref ref="my-tube"/>
				<scale name="tube_scale" x="1" y="2" z="3"/>
			</scaledSolid>`,
			want: ScaledSolid{
				Name:  "scaled-tube",
				Ref:   SolidRef{XMLName: xml.Name{Local: "solidref"}, Ref: "my-tube"},
				Scale: Scale{Name: "tube_scale", X: "1", Y: "2", Z: "3"},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var v1 = reflect.New(reflect.TypeOf(tc.want)).Elem()
			err := xml.NewDecoder(bytes.NewReader([]byte(tc.raw))).Decode(v1.Addr().Interface())
			if err != nil {
				t.Fatal(err)
			}

			want := reflect.New(reflect.TypeOf(tc.want)).Elem()
			want.Set(reflect.ValueOf(tc.want))
			field := want.FieldByName("XMLName")
			if field != (reflect.Value{}) {
				field.Set(v1.FieldByName("XMLName"))
			}
			if !reflect.DeepEqual(v1.Interface(), want.Interface()) {
				t.Fatalf("error:\ngot = %#v\nwant= %#v", v1.Interface(), want.Interface())
			}

			out := new(bytes.Buffer)
			err = xml.NewEncoder(out).Encode(v1.Interface())
			if err != nil {
				t.Fatal(err)
			}

			raw := out.String()
			var v2 = reflect.New(reflect.TypeOf(tc.want)).Elem()
			err = xml.NewDecoder(bytes.NewReader([]byte(raw))).Decode(v2.Addr().Interface())
			if err != nil {
				t.Fatal(err)
			}

			if !reflect.DeepEqual(v2.Interface(), want.Interface()) {
				t.Fatalf("error:\ngot = %#v\nwant= %#v", v2.Interface(), want.Interface())
			}
		})
	}
}
