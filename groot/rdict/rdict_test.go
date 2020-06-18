// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"math"
	"reflect"
	"testing"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rmeta"
)

func TestElementGetRange(t *testing.T) {
	for _, tc := range []struct {
		name               string
		title              string
		rtype              rmeta.Enum
		xmin, xmax, factor float64
	}{
		{
			name:  "empty",
			title: "",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-d32",
			title: "var/d",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-f64",
			title: "var/D",
			rtype: rmeta.Float64,
		},
		{
			name:  "normal-f64-ndims",
			title: "var[10][20][30]/D",
			rtype: rmeta.Float64,
		},
		{
			name:  "normal-1d",
			title: "var[3]/d",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-2d",
			title: "var[3][4]/d",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-3d",
			title: "var[3][4][5]/d",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-with-brackets",
			title: "From [tleft,tright+10 ns]",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-with-brackets-2",
			title: "Bias voltage [V]",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-with-brackets-3",
			title: "Bias voltage [0, 100]",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-with-brackets-4",
			title: "Bias/voltage [0, 100]",
			rtype: rmeta.Double32,
		},
		{
			name:  "normal-with-brackets-5",
			title: "Bias voltage [0]",
			rtype: rmeta.Double32,
		},
		{
			name:   "range",
			title:  "[ 0 , 100 ]",
			rtype:  rmeta.Double32,
			xmin:   0,
			xmax:   100,
			factor: float64(0xffffffff) / 100,
		},
		{
			name:   "range-ndim",
			title:  "var[3]/d[ 0 , 100 ]",
			rtype:  rmeta.Double32,
			xmin:   0,
			xmax:   100,
			factor: float64(0xffffffff) / 100,
		},
		{
			name:   "range-nbits",
			title:  "[ 10 , 100, 30 ]",
			rtype:  rmeta.Double32,
			xmin:   10,
			xmax:   100,
			factor: float64(1<<30) / 90,
		},
		{
			name:   "range-nbits-1d",
			title:  "var[3]/d[ 10 , 100, 30 ]",
			rtype:  rmeta.Double32,
			xmin:   10,
			xmax:   100,
			factor: float64(1<<30) / 90,
		},
		{
			name:   "range-pi",
			title:  "[ -pi , pi ]",
			rtype:  rmeta.Double32,
			xmin:   -math.Pi,
			xmax:   +math.Pi,
			factor: float64(0xffffffff) / (2 * math.Pi),
		},
		{
			name:   "range-pi/2",
			title:  "[ -pi/2 , 2pi ]",
			rtype:  rmeta.Double32,
			xmin:   -math.Pi / 2,
			xmax:   2 * math.Pi,
			factor: float64(0xffffffff) / (2*math.Pi + math.Pi/2),
		},
		{
			name:   "range-twopi/4",
			title:  "[ -pi/4 , twopi ]",
			rtype:  rmeta.Double32,
			xmin:   -math.Pi / 4,
			xmax:   2 * math.Pi,
			factor: float64(0xffffffff) / (2*math.Pi + math.Pi/4),
		},
		{
			name:   "range-2pi",
			title:  "[ -2*pi , 2*pi ]",
			rtype:  rmeta.Double32,
			xmin:   -2 * math.Pi,
			xmax:   +2 * math.Pi,
			factor: float64(0xffffffff) / (4 * math.Pi),
		},
		{
			name:  "float32-15bits",
			title: "[ 0 , 0 , 15 ]",
			rtype: rmeta.Double32,
		},
		{
			name:  "float32-14bits",
			title: "[ 0 , 0 , 14 ]",
			rtype: rmeta.Double32,
			xmin:  float64(14) + 0.1,
		},
		{
			name:  "float32-3bits",
			title: "[ 10 , 10 , 3 ]",
			rtype: rmeta.Double32,
			xmin:  float64(3) + 0.1,
			xmax:  10,
		},
		{
			name:  "float32-2bits",
			title: "[ 0 , 0 , 2 ]",
			rtype: rmeta.Double32,
			xmin:  float64(2) + 0.1,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			elmt := Element{
				Name: *rbase.NewNamed(tc.name, tc.title),
				Type: tc.rtype,
			}
			elmt.parse()
			if got, want := elmt.XMin, tc.xmin; got != want {
				t.Fatalf("invalid xmin: got=%v, want=%v", got, want)
			}
			if got, want := elmt.XMax, tc.xmax; got != want {
				t.Fatalf("invalid xmax: got=%v, want=%v", got, want)
			}
			if got, want := elmt.Factor, tc.factor; got != want {
				t.Fatalf("invalid factor: got=%v, want=%v", got, want)
			}
		})
	}
}

func TestParseStdContainers(t *testing.T) {
	for _, tc := range []struct {
		name   string
		parse  func(string) []string
		want   []string
		panics string
	}{
		// std::vector
		{
			name:  "vector<int>",
			parse: parseStdVector,
			want:  []string{"int"},
		},
		{
			name:  "std::vector<int>",
			parse: parseStdVector,
			want:  []string{"int"},
		},
		{
			name:  "vector<vector<int>>",
			parse: parseStdVector,
			want:  []string{"vector<int>"},
		},
		{
			name:  "vector<int,allocator<int>>",
			parse: parseStdVector,
			want:  []string{"int", "allocator<int>"},
		},
		{
			name:  "vector<map<int,long int>>",
			parse: parseStdVector,
			want:  []string{"map<int,long int>"},
		},
		{
			name:   "vector<int",
			parse:  parseStdVector,
			panics: `invalid std::vector container name (missing '>'): "vector<int"`,
		},
		{
			name:   "xvector<int>",
			parse:  parseStdVector,
			panics: `invalid std::vector container name (missing 'vector<'): "xvector<int>"`,
		},
		{
			name:   "vector<>",
			parse:  parseStdVector,
			panics: `invalid std::vector container name (missing element type): "vector<>"`,
		},
		{
			name:   "vector<t1,t2,t3>",
			parse:  parseStdVector,
			panics: `invalid std::vector template "vector<t1,t2,t3>"`,
		},
		// std::map
		{
			name:  "map< int , int >",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "map<int,int>",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "std::map<int,int>",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "map<int,int>",
			parse: parseStdMap,
			want:  []string{"int", "int"},
		},
		{
			name:  "map<int,string>",
			parse: parseStdMap,
			want:  []string{"int", "string"},
		},
		{
			name:  "map<int,vector<int>>",
			parse: parseStdMap,
			want:  []string{"int", "vector<int>"},
		},
		{
			name:  "map<int,vector<int> >",
			parse: parseStdMap,
			want:  []string{"int", "vector<int>"},
		},
		{
			name:  "map<int,map<string,int> >",
			parse: parseStdMap,
			want:  []string{"int", "map<string,int>"},
		},
		{
			name:  "map<map<string,int>, int>",
			parse: parseStdMap,
			want:  []string{"map<string,int>", "int"},
		},
		{
			name:  "map<map<string,int>, map<int,string>>",
			parse: parseStdMap,
			want:  []string{"map<string,int>", "map<int,string>"},
		},
		{
			name:  "map<long int,long int>",
			parse: parseStdMap,
			want:  []string{"long int", "long int"},
		},
		{
			name:  "map<long int, vector<long int>, allocator<pair<const long int, vector<long int>>>",
			parse: parseStdMap,
			want:  []string{"long int", "vector<long int>", "allocator<pair<const long int, vector<long int>>"},
		},
		{
			name:   "map<k,v",
			parse:  parseStdMap,
			panics: `invalid std::map container name (missing '>'): "map<k,v"`,
		},
		{
			name:   "map<k,v,a,XXX>",
			parse:  parseStdMap,
			panics: `invalid std::map template "map<k,v,a,XXX>"`,
		},
		{
			name:   "map<>",
			parse:  parseStdMap,
			panics: `invalid std::map template "map<>"`,
		},
		{
			name:   "xmap<k,v>",
			parse:  parseStdMap,
			panics: `invalid std::map container name (missing 'map<'): "xmap<k,v>"`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if tc.panics != "" {
				defer func() {
					err := recover()
					if err == nil {
						t.Fatalf("expected a panic (%s)", tc.panics)
					}
					if got, want := err.(error).Error(), tc.panics; got != want {
						t.Fatalf("invalid panic message: got=%s, want=%s", got, want)
					}
				}()
			}
			got := tc.parse(tc.name)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got=%q, want=%q", got, tc.want)
			}
		})
	}
}

func TestGenChecksum(t *testing.T) {
	sbt := func(n string, t rmeta.Enum, et string) *StreamerBasicType {
		return &StreamerBasicType{
			StreamerElement: Element{
				Name:  *rbase.NewNamed(n, ""),
				Type:  t,
				EName: et,
			}.New(),
		}
	}
	sbsli := func(n, t string, typ rmeta.Enum, et string) *StreamerBasicType {
		return &StreamerBasicType{
			StreamerElement: Element{
				Name:  *rbase.NewNamed(n, t),
				Type:  typ + rmeta.OffsetP,
				EName: et + "*",
			}.New(),
		}
	}
	sbarr := func(n string, t rmeta.Enum, et string, i int) *StreamerBasicType {
		return &StreamerBasicType{
			StreamerElement: Element{
				Name:   *rbase.NewNamed(n, ""),
				Type:   t + rmeta.OffsetL,
				ArrDim: 1,
				ArrLen: int32(i),
				MaxIdx: [5]int32{int32(i)},
				EName:  et,
			}.New(),
		}
	}
	tstr := func(n string) *StreamerString {
		return &StreamerString{
			StreamerElement: Element{
				Name:  *rbase.NewNamed(n, ""),
				Type:  rmeta.TString,
				EName: "TString",
			}.New(),
		}
	}
	stlstr := func(n string) *StreamerSTLstring {
		return &StreamerSTLstring{
			StreamerSTL: StreamerSTL{
				StreamerElement: Element{
					Name:  *rbase.NewNamed(n, ""),
					Type:  rmeta.TString,
					EName: "string",
				}.New(),
			},
		}
	}
	stlvec := func(n, et string) *StreamerSTL {
		return &StreamerSTL{
			StreamerElement: Element{
				Name:  *rbase.NewNamed(n, ""),
				Type:  rmeta.STL,
				EName: "vector<" + et + ">",
			}.New(),
		}
	}
	soa := func(n, et string) *StreamerObjectAny {
		return &StreamerObjectAny{
			StreamerElement: Element{
				Name:  *rbase.NewNamed(n, ""),
				Type:  rmeta.Any,
				EName: et,
			}.New(),
		}
	}

	for _, tc := range []struct {
		name  string
		elems []rbytes.StreamerElement
		want  uint32
	}{
		{
			name: "P3",
			elems: []rbytes.StreamerElement{
				sbt("Px", rmeta.Int32, "int"),
				sbt("Py", rmeta.Float64, "double"),
				sbt("Pz", rmeta.Int32, "int"),
			},
			want: 1678002455, // obtained w/ 6.20/04
		},
		{
			name: "ArrF64",
			elems: []rbytes.StreamerElement{
				sbarr("Arr", rmeta.Float64, "double", 10),
			},
			want: 1711917547, // obtained w/ 6.20/04
		},
		{
			name: "SliF64",
			elems: []rbytes.StreamerElement{
				sbt("N", rmeta.Int32, "int"),
				sbsli("Sli", "[N]", rmeta.Float64, "double"),
			},
			want: 193076120, // obtained w/ 6.20/04
		},
		{
			name: "StlVecF64",
			elems: []rbytes.StreamerElement{
				stlvec("Stl", "double"),
			},
			want: 2364618348, // obtained w/ 6.20/04
		},
		{
			name: "Event",
			elems: []rbytes.StreamerElement{
				tstr("Beg"),
				sbt("I16", rmeta.Int16, "short"),
				sbt("I32", rmeta.Int32, "int"),
				sbt("I64", rmeta.Int64, "long"),
				sbt("U16", rmeta.Uint16, "unsigned short"),
				sbt("U32", rmeta.Uint32, "unsigned int"),
				sbt("U64", rmeta.Uint64, "unsigned long"),
				sbt("F32", rmeta.Float32, "float"),
				sbt("F64", rmeta.Float64, "double"),
				tstr("Str"),
				soa("P3", "P3"),
				sbarr("ArrayI16", rmeta.Int16, "short", 10),
				sbarr("ArrayI32", rmeta.Int32, "int", 10),
				sbarr("ArrayI64", rmeta.Int64, "long", 10),
				sbarr("ArrayU16", rmeta.Uint16, "unsigned short", 10),
				sbarr("ArrayU32", rmeta.Uint32, "unsigned int", 10),
				sbarr("ArrayU64", rmeta.Uint64, "unsigned long", 10),
				sbarr("ArrayF32", rmeta.Float32, "float", 10),
				sbarr("ArrayF64", rmeta.Float64, "double", 10),
				sbt("N", rmeta.Int32, "int"),
				sbsli("SliceI16", "[N]", rmeta.Int16, "short"),
				sbsli("SliceI32", "[N]", rmeta.Int32, "int"),
				sbsli("SliceI64", "[N]", rmeta.Int64, "long"),
				sbsli("SliceU16", "[N]", rmeta.Uint16, "unsigned short"),
				sbsli("SliceU32", "[N]", rmeta.Uint32, "unsigned int"),
				sbsli("SliceU64", "[N]", rmeta.Uint64, "unsigned long"),
				sbsli("SliceF32", "[N]", rmeta.Float32, "float"),
				sbsli("SliceF64", "[N]", rmeta.Float64, "double"),
				stlstr("StdStr"),
				stlvec("StlVecI16", "short"),
				stlvec("StlVecI32", "int"),
				stlvec("StlVecI64", "long"),
				stlvec("StlVecU16", "unsigned short"),
				stlvec("StlVecU32", "unsigned int"),
				stlvec("StlVecU64", "unsigned long"),
				stlvec("StlVecF32", "float"),
				stlvec("StlVecF64", "double"),
				stlvec("StlVecStr", "string"),
				tstr("End"),
			},
			want: 1123173915, // obtained w/ 6.20/04
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			chksum := genChecksum(tc.name, tc.elems)
			if got, want := chksum, tc.want; got != want {
				t.Fatalf("invalid checksum: got=%d, want=%d", got, want)
			}
		})
	}
}
