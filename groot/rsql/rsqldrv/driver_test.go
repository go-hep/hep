// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv_test

import (
	"database/sql"
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/xwb1989/sqlparser"
	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rsql/rsqldrv"
	"go-hep.org/x/hep/groot/rtree"
)

func TestOpenDB(t *testing.T) {
	f, err := groot.Open("../../testdata/simple.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := o.(rtree.Tree)

	db := rsqldrv.OpenDB(rtree.FileOf(tree))
	defer db.Close()

	type data struct {
		i32 int32
		f32 float32
		str string
	}

	rows, err := db.Query("SELECT * FROM tree")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var (
		got  []data
		want = []data{
			{1, 1.1, "uno"},
			{2, 2.2, "dos"},
			{3, 3.3, "tres"},
			{4, 4.4, "quatro"},
		}
	)

	for rows.Next() {
		var v data
		err = rows.Scan(&v.i32, &v.f32, &v.str)
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, v)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid select\ngot = %#v\nwant= %#v", got, want)
	}
}

func TestOpenWithConnector(t *testing.T) {
	f, err := groot.Open("../../testdata/simple.root")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		t.Fatal(err)
	}

	tree := o.(rtree.Tree)

	db := sql.OpenDB(rsqldrv.Connector(rtree.FileOf(tree)))
	defer db.Close()

	type data struct {
		i32 int32
		f32 float32
		str string
	}

	rows, err := db.Query("SELECT * FROM tree")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	var (
		got  []data
		want = []data{
			{1, 1.1, "uno"},
			{2, 2.2, "dos"},
			{3, 3.3, "tres"},
			{4, 4.4, "quatro"},
		}
	)

	for rows.Next() {
		var v data
		err = rows.Scan(&v.i32, &v.f32, &v.str)
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, v)
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("invalid select\ngot = %#v\nwant= %#v", got, want)
	}
}

func TestQuery(t *testing.T) {
	db, err := sql.Open("root", "../../testdata/simple.root")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	type data struct {
		i32 int32
		f32 float32
		str string
	}
	for _, tc := range []struct {
		query string
		args  []interface{}
		cols  []string
		want  []data
	}{
		{
			query: `SELECT * FROM tree`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one+10, two+20, "--"+three+"--") FROM tree`,
			cols:  []string{"", "", ""},
			want: []data{
				{11, 21.1, "--uno--"},
				{12, 22.2, "--dos--"},
				{13, 23.3, "--tres--"},
				{14, 24.4, "--quatro--"},
			},
		},
		{
			query: `SELECT (one+?, two+?, ?+three+"--") FROM tree`,
			cols:  []string{"", "", ""},
			args:  []interface{}{int32(10), 20.0, "++"},
			want: []data{
				{11, 21.1, "++uno--"},
				{12, 22.2, "++dos--"},
				{13, 23.3, "++tres--"},
				{14, 24.4, "++quatro--"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (one <= 2)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (2 >= one)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (one <= ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []interface{}{int32(2)},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (one <= ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []interface{}{2},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (two <= ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []interface{}{2.2},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (two <= ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []interface{}{3},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (three != ? && three != ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []interface{}{"tres", "quatro"},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (one > -1 && two > -1 && three != "N/A")`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (one >= 3)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (one > 2)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (two >= 2.2 || two < 10)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (two >= 2.2 && two < 4)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE ((two+2) >= 4.2)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE ((2+two) >= 4.2)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (2+two >= 4.2)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE ((2*two) >= 4.2)`,
			cols:  []string{"one", "two", "three"},
			want: []data{
				{2, 2.2, "dos"},
				{3, 3.3, "tres"},
				{4, 4.4, "quatro"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (three = ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []interface{}{"quatro"},
			want: []data{
				{4, 4.4, "quatro"},
			},
		},
	} {
		t.Run(tc.query, func(t *testing.T) {
			_, err := sqlparser.Parse(tc.query)
			if err != nil {
				t.Fatalf("could not parse query %q: %v", tc.query, err)
			}

			rows, err := db.Query(tc.query, tc.args...)
			if err != nil {
				t.Fatal(err)
			}
			defer rows.Close()

			cols, err := rows.Columns()
			if err != nil {
				t.Fatal(err)
			}

			if got, want := cols, tc.cols; !reflect.DeepEqual(got, want) {
				t.Fatalf("invalid columns\ngot= %q\nwant=%q", got, want)
			}

			var got []data
			for rows.Next() {
				var v data
				err = rows.Scan(&v.i32, &v.f32, &v.str)
				if err != nil {
					t.Fatal(err)
				}
				got = append(got, v)
			}

			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("invalid select\ngot = %#v\nwant= %#v", got, tc.want)
			}
		})
	}
}

func TestFlatTree(t *testing.T) {
	type event struct {
		I32    int32       `groot:"Int32"`
		I64    int64       `groot:"Int64"`
		U32    uint32      `groot:"UInt32"`
		U64    uint64      `groot:"UInt64"`
		F32    float32     `groot:"Float32"`
		F64    float64     `groot:"Float64"`
		Str    string      `groot:"Str"`
		ArrI32 [10]int32   `groot:"ArrayInt32"`
		ArrI64 [10]int64   `groot:"ArrayInt64"`
		ArrU32 [10]uint32  `groot:"ArrayUInt32"`
		ArrU64 [10]uint64  `groot:"ArrayUInt64"`
		ArrF32 [10]float32 `groot:"ArrayFloat32"`
		ArrF64 [10]float64 `groot:"ArrayFloat64"`
		N      int32       `groot:"N"`
		SliI32 []int32     `groot:"SliceInt32"`
		SliI64 []int64     `groot:"SliceInt64"`
		SliU32 []uint32    `groot:"SliceUInt32"`
		SliU64 []uint64    `groot:"SliceUInt64"`
		SliF32 []float32   `groot:"SliceFloat32"`
		SliF64 []float64   `groot:"SliceFloat64"`
	}

	want := func(i int64) (data event) {
		data.I32 = int32(i)
		data.I64 = int64(i)
		data.U32 = uint32(i)
		data.U64 = uint64(i)
		data.F32 = float32(i)
		data.F64 = float64(i)
		data.Str = fmt.Sprintf("evt-%03d", i)
		for ii := range data.ArrI32 {
			data.ArrI32[ii] = int32(i)
			data.ArrI64[ii] = int64(i)
			data.ArrU32[ii] = uint32(i)
			data.ArrU64[ii] = uint64(i)
			data.ArrF32[ii] = float32(i)
			data.ArrF64[ii] = float64(i)
		}
		data.N = int32(i) % 10
		data.SliI32 = make([]int32, int(data.N))
		data.SliI64 = make([]int64, int(data.N))
		data.SliU32 = make([]uint32, int(data.N))
		data.SliU64 = make([]uint64, int(data.N))
		data.SliF32 = make([]float32, int(data.N))
		data.SliF64 = make([]float64, int(data.N))
		for ii := 0; ii < int(data.N); ii++ {
			data.SliI32[ii] = int32(i)
			data.SliI64[ii] = int64(i)
			data.SliU32[ii] = uint32(i)
			data.SliU64[ii] = uint64(i)
			data.SliF32[ii] = float32(i)
			data.SliF64[ii] = float64(i)
		}
		return data
	}

	db, err := sql.Open("root", "../../testdata/small-flat-tree.root")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tree")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	cols, err := rows.ColumnTypes()
	if err != nil {
		t.Fatal(err)
	}

	for i, want := range []struct {
		name        string
		hasNullable bool
		hasLength   bool

		nullable bool
		length   int64
		scanType reflect.Type
	}{
		{"Int32", true, false, false, 0, reflect.ValueOf(int32(0)).Type()},
		{"Int64", true, false, false, 0, reflect.ValueOf(int64(0)).Type()},
		{"UInt32", true, false, false, 0, reflect.ValueOf(uint32(0)).Type()},
		{"UInt64", true, false, false, 0, reflect.ValueOf(uint64(0)).Type()},
		{"Float32", true, false, false, 0, reflect.ValueOf(float32(0)).Type()},
		{"Float64", true, false, false, 0, reflect.ValueOf(float64(0)).Type()},
		{"Str", true, false, false, 0, reflect.ValueOf("").Type()},
		{"ArrayInt32", true, true, false, 10, reflect.ValueOf(int32(0)).Type()},
		{"ArrayInt64", true, true, false, 10, reflect.ValueOf(int64(0)).Type()},
		{"ArrayUInt32", true, true, false, 10, reflect.ValueOf(uint32(0)).Type()},
		{"ArrayUInt64", true, true, false, 10, reflect.ValueOf(uint64(0)).Type()},
		{"ArrayFloat32", true, true, false, 10, reflect.ValueOf(float32(0)).Type()},
		{"ArrayFloat64", true, true, false, 10, reflect.ValueOf(float64(0)).Type()},
		{"N", true, false, false, 0, reflect.ValueOf(int32(0)).Type()},
		{"SliceInt32", true, true, true, math.MaxInt64, reflect.ValueOf(int32(0)).Type()},
		{"SliceInt64", true, true, true, math.MaxInt64, reflect.ValueOf(int64(0)).Type()},
		{"SliceUInt32", true, true, true, math.MaxInt64, reflect.ValueOf(uint32(0)).Type()},
		{"SliceUInt64", true, true, true, math.MaxInt64, reflect.ValueOf(uint64(0)).Type()},
		{"SliceFloat32", true, true, true, math.MaxInt64, reflect.ValueOf(float32(0)).Type()},
		{"SliceFloat64", true, true, true, math.MaxInt64, reflect.ValueOf(float64(0)).Type()},
	} {
		got := cols[i]
		if got.Name() != want.name {
			t.Fatalf("col[%d]: invalid name. got=%q, want=%q", i, got.Name(), want.name)
		}

		nullable, hasNullable := got.Nullable()
		if hasNullable != want.hasNullable {
			t.Fatalf("col[%d]: invalid nullable state. got=%v, want=%v", i, hasNullable, want.hasNullable)
		}
		if nullable != want.nullable {
			t.Fatalf("col[%d]: invalid nullable. got=%v, want=%v", i, nullable, want.nullable)
		}

		length, hasLength := got.Length()
		if hasLength != want.hasLength {
			t.Fatalf("col[%d]: invalid length state. got=%v, want=%v", i, hasLength, want.hasLength)
		}
		if length != want.length {
			t.Fatalf("col[%d]: invalid length. got=%v, want=%v", i, length, want.length)
		}

		if got, want := got.ScanType(), want.scanType; !reflect.DeepEqual(got, want) {
			t.Fatalf("col[%d]: invalid type. got=%v, want=%v", i, got, want)
		}
	}

	var i int64
	for rows.Next() {
		var v event
		err = rows.Scan(
			&v.I32, &v.I64,
			&v.U32, &v.U64,
			&v.F32, &v.F64,
			&v.Str,
			&v.ArrI32, &v.ArrI64,
			&v.ArrU32, &v.ArrU64,
			&v.ArrF32, &v.ArrF64,
			&v.N,
			&v.SliI32, &v.SliI64,
			&v.SliU32, &v.SliU64,
			&v.SliF32, &v.SliF64,
		)
		if err != nil {
			t.Fatal(err)
		}

		if got, want := v, want(i); !reflect.DeepEqual(got, want) {
			t.Fatalf("invalid row[%d]:\ngot= %#v\nwant=%#v\n", i, got, want)
		}
		i++
	}
}
