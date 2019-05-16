// Copyright 2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv_test

import (
	"database/sql"
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
