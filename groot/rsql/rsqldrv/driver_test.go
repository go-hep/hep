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
	"go-hep.org/x/hep/groot/root"
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
		args  []any
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
			args:  []any{int32(10), 20.0, "++"},
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
			args:  []any{int32(2)},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (one <= ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []any{2},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (two <= ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []any{2.2},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (two <= ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []any{3},
			want: []data{
				{1, 1.1, "uno"},
				{2, 2.2, "dos"},
			},
		},
		{
			query: `SELECT (one, two, three) FROM tree WHERE (three != ? && three != ?)`,
			cols:  []string{"one", "two", "three"},
			args:  []any{"tres", "quatro"},
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
			args:  []any{"quatro"},
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

type eventData struct {
	B      bool              `groot:"B"`
	Str    string            `groot:"Str"`
	I8     int8              `groot:"I8"`
	I16    int16             `groot:"I16"`
	I32    int32             `groot:"I32"`
	I64    int64             `groot:"I64"`
	U8     uint8             `groot:"U8"`
	U16    uint16            `groot:"U16"`
	U32    uint32            `groot:"U32"`
	U64    uint64            `groot:"U64"`
	F32    float32           `groot:"F32"`
	F64    float64           `groot:"F64"`
	D16    root.Float16      `groot:"D16"`
	D32    root.Double32     `groot:"D32"`
	ArrBs  [10]bool          `groot:"ArrBs[10]"`
	ArrI8  [10]int8          `groot:"ArrI8[10]"`
	ArrI16 [10]int16         `groot:"ArrI16[10]"`
	ArrI32 [10]int32         `groot:"ArrI32[10]"`
	ArrI64 [10]int64         `groot:"ArrI64[10]"`
	ArrU8  [10]uint8         `groot:"ArrU8[10]"`
	ArrU16 [10]uint16        `groot:"ArrU16[10]"`
	ArrU32 [10]uint32        `groot:"ArrU32[10]"`
	ArrU64 [10]uint64        `groot:"ArrU64[10]"`
	ArrF32 [10]float32       `groot:"ArrF32[10]"`
	ArrF64 [10]float64       `groot:"ArrF64[10]"`
	ArrD16 [10]root.Float16  `groot:"ArrD16[10]"`
	ArrD32 [10]root.Double32 `groot:"ArrD32[10]"`
	N      int32             `groot:"N"`
	SliBs  []bool            `groot:"SliBs[N]"`
	SliI8  []int8            `groot:"SliI8[N]"`
	SliI16 []int16           `groot:"SliI16[N]"`
	SliI32 []int32           `groot:"SliI32[N]"`
	SliI64 []int64           `groot:"SliI64[N]"`
	SliU8  []uint8           `groot:"SliU8[N]"`
	SliU16 []uint16          `groot:"SliU16[N]"`
	SliU32 []uint32          `groot:"SliU32[N]"`
	SliU64 []uint64          `groot:"SliU64[N]"`
	SliF32 []float32         `groot:"SliF32[N]"`
	SliF64 []float64         `groot:"SliF64[N]"`
	SliD16 []root.Float16    `groot:"SliD16[N]"`
	SliD32 []root.Double32   `groot:"SliD32[N]"`
}

func (eventData) want(i int64) (data eventData) {
	data.B = i%2 == 0
	data.Str = fmt.Sprintf("str-%d", i)
	data.I8 = int8(-i)
	data.I16 = int16(-i)
	data.I32 = int32(-i)
	data.I64 = int64(-i)
	data.U8 = uint8(i)
	data.U16 = uint16(i)
	data.U32 = uint32(i)
	data.U64 = uint64(i)
	data.F32 = float32(i)
	data.F64 = float64(i)
	data.D16 = root.Float16(i)
	data.D32 = root.Double32(i)
	for ii := range data.ArrI32 {
		data.ArrBs[ii] = ii == int(i)
		data.ArrI8[ii] = int8(-i)
		data.ArrI16[ii] = int16(-i)
		data.ArrI32[ii] = int32(-i)
		data.ArrI64[ii] = int64(-i)
		data.ArrU8[ii] = uint8(i)
		data.ArrU16[ii] = uint16(i)
		data.ArrU32[ii] = uint32(i)
		data.ArrU64[ii] = uint64(i)
		data.ArrF32[ii] = float32(i)
		data.ArrF64[ii] = float64(i)
		data.ArrD16[ii] = root.Float16(i)
		data.ArrD32[ii] = root.Double32(i)
	}
	data.N = int32(i) % 10
	data.SliBs = make([]bool, int(data.N))
	data.SliI8 = make([]int8, int(data.N))
	data.SliI16 = make([]int16, int(data.N))
	data.SliI32 = make([]int32, int(data.N))
	data.SliI64 = make([]int64, int(data.N))
	data.SliU8 = make([]uint8, int(data.N))
	data.SliU16 = make([]uint16, int(data.N))
	data.SliU32 = make([]uint32, int(data.N))
	data.SliU64 = make([]uint64, int(data.N))
	data.SliF32 = make([]float32, int(data.N))
	data.SliF64 = make([]float64, int(data.N))
	data.SliD16 = make([]root.Float16, int(data.N))
	data.SliD32 = make([]root.Double32, int(data.N))
	for ii := range int(data.N) {
		data.SliBs[ii] = (ii + 1) == int(i)
		data.SliI8[ii] = int8(-i)
		data.SliI16[ii] = int16(-i)
		data.SliI32[ii] = int32(-i)
		data.SliI64[ii] = int64(-i)
		data.SliU8[ii] = uint8(i)
		data.SliU16[ii] = uint16(i)
		data.SliU32[ii] = uint32(i)
		data.SliU64[ii] = uint64(i)
		data.SliF32[ii] = float32(i)
		data.SliF64[ii] = float64(i)
		data.SliD16[ii] = root.Float16(i)
		data.SliD32[ii] = root.Double32(i)
	}
	return data
}

func TestFlatTree(t *testing.T) {
	db, err := sql.Open("root", "../../testdata/x-flat-tree.root")
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
		{"B", true, false, false, 0, reflect.ValueOf(false).Type()},
		{"Str", true, false, false, 0, reflect.ValueOf("").Type()},
		{"I8", true, false, false, 0, reflect.ValueOf(int8(0)).Type()},
		{"I16", true, false, false, 0, reflect.ValueOf(int16(0)).Type()},
		{"I32", true, false, false, 0, reflect.ValueOf(int32(0)).Type()},
		{"I64", true, false, false, 0, reflect.ValueOf(int64(0)).Type()},
		{"U8", true, false, false, 0, reflect.ValueOf(uint8(0)).Type()},
		{"U16", true, false, false, 0, reflect.ValueOf(uint16(0)).Type()},
		{"U32", true, false, false, 0, reflect.ValueOf(uint32(0)).Type()},
		{"U64", true, false, false, 0, reflect.ValueOf(uint64(0)).Type()},
		{"F32", true, false, false, 0, reflect.ValueOf(float32(0)).Type()},
		{"F64", true, false, false, 0, reflect.ValueOf(float64(0)).Type()},
		{"D16", true, false, false, 0, reflect.ValueOf(root.Float16(0)).Type()},
		{"D32", true, false, false, 0, reflect.ValueOf(root.Double32(0)).Type()},
		{"ArrBs", true, true, false, 10, reflect.ValueOf(false).Type()},
		{"ArrI8", true, true, false, 10, reflect.ValueOf(int8(0)).Type()},
		{"ArrI16", true, true, false, 10, reflect.ValueOf(int16(0)).Type()},
		{"ArrI32", true, true, false, 10, reflect.ValueOf(int32(0)).Type()},
		{"ArrI64", true, true, false, 10, reflect.ValueOf(int64(0)).Type()},
		{"ArrU8", true, true, false, 10, reflect.ValueOf(uint8(0)).Type()},
		{"ArrU16", true, true, false, 10, reflect.ValueOf(uint16(0)).Type()},
		{"ArrU32", true, true, false, 10, reflect.ValueOf(uint32(0)).Type()},
		{"ArrU64", true, true, false, 10, reflect.ValueOf(uint64(0)).Type()},
		{"ArrF32", true, true, false, 10, reflect.ValueOf(float32(0)).Type()},
		{"ArrF64", true, true, false, 10, reflect.ValueOf(float64(0)).Type()},
		{"ArrD16", true, true, false, 10, reflect.ValueOf(root.Float16(0)).Type()},
		{"ArrD32", true, true, false, 10, reflect.ValueOf(root.Double32(0)).Type()},
		{"N", true, false, false, 0, reflect.ValueOf(int32(0)).Type()},
		{"SliBs", true, true, true, math.MaxInt64, reflect.ValueOf(false).Type()},
		{"SliI8", true, true, true, math.MaxInt64, reflect.ValueOf(int8(0)).Type()},
		{"SliI16", true, true, true, math.MaxInt64, reflect.ValueOf(int16(0)).Type()},
		{"SliI32", true, true, true, math.MaxInt64, reflect.ValueOf(int32(0)).Type()},
		{"SliI64", true, true, true, math.MaxInt64, reflect.ValueOf(int64(0)).Type()},
		{"SliU8", true, true, true, math.MaxInt64, reflect.ValueOf(uint8(0)).Type()},
		{"SliU16", true, true, true, math.MaxInt64, reflect.ValueOf(uint16(0)).Type()},
		{"SliU32", true, true, true, math.MaxInt64, reflect.ValueOf(uint32(0)).Type()},
		{"SliU64", true, true, true, math.MaxInt64, reflect.ValueOf(uint64(0)).Type()},
		{"SliF32", true, true, true, math.MaxInt64, reflect.ValueOf(float32(0)).Type()},
		{"SliF64", true, true, true, math.MaxInt64, reflect.ValueOf(float64(0)).Type()},
		{"SliD16", true, true, true, math.MaxInt64, reflect.ValueOf(root.Float16(0)).Type()},
		{"SliD32", true, true, true, math.MaxInt64, reflect.ValueOf(root.Double32(0)).Type()},
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

	var (
		want = eventData{}.want
		i    int64
	)
	for rows.Next() {
		var v eventData
		err = rows.Scan(
			&v.B,
			&v.Str,
			&v.I8, &v.I16, &v.I32, &v.I64,
			&v.U8, &v.U16, &v.U32, &v.U64,
			&v.F32, &v.F64,
			&v.D16, &v.D32,
			&v.ArrBs,
			&v.ArrI8, &v.ArrI16, &v.ArrI32, &v.ArrI64,
			&v.ArrU8, &v.ArrU16, &v.ArrU32, &v.ArrU64,
			&v.ArrF32, &v.ArrF64,
			&v.ArrD16, &v.ArrD32,
			&v.N,
			&v.SliBs,
			&v.SliI8, &v.SliI16, &v.SliI32, &v.SliI64,
			&v.SliU8, &v.SliU16, &v.SliU32, &v.SliU64,
			&v.SliF32, &v.SliF64,
			&v.SliD16, &v.SliD32,
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
