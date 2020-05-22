// Copyright ©2019 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rsqldrv_test

import (
	"database/sql"
	"fmt"
	"log"

	"go-hep.org/x/hep/groot"
	"go-hep.org/x/hep/groot/rsql/rsqldrv"
	"go-hep.org/x/hep/groot/rtree"
)

func ExampleOpen() {
	db, err := sql.Open("root", "../../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tree")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type data struct {
		i32 int32
		f32 float32
		str string
	}

	n := 0
	for rows.Next() {
		var v data
		err := rows.Scan(&v.i32, &v.f32, &v.str)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("row[%d]: (%v, %v, %q)\n", n, v.i32, v.f32, v.str)
		n++
	}

	// Output:
	// row[0]: (1, 1.1, "uno")
	// row[1]: (2, 2.2, "dos")
	// row[2]: (3, 3.3, "tres")
	// row[3]: (4, 4.4, "quatro")
}

func ExampleOpen_tuple() {
	db, err := sql.Open("root", "../../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT (one, two, three) FROM tree")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type data struct {
		i32 int32
		f32 float32
		str string
	}

	n := 0
	for rows.Next() {
		var v data
		err := rows.Scan(&v.i32, &v.f32, &v.str)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("row[%d]: (%v, %v, %q)\n", n, v.i32, v.f32, v.str)
		n++
	}

	// Output:
	// row[0]: (1, 1.1, "uno")
	// row[1]: (2, 2.2, "dos")
	// row[2]: (3, 3.3, "tres")
	// row[3]: (4, 4.4, "quatro")
}

func ExampleOpen_whereStmt() {
	db, err := sql.Open("root", "../../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT (one, two, three) FROM tree WHERE (one>2 && two < 5)")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type data struct {
		i32 int32
		f32 float32
		str string
	}

	n := 0
	for rows.Next() {
		var v data
		err := rows.Scan(&v.i32, &v.f32, &v.str)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("row[%d]: (%v, %v, %q)\n", n, v.i32, v.f32, v.str)
		n++
	}

	// Output:
	// row[0]: (3, 3.3, "tres")
	// row[1]: (4, 4.4, "quatro")
}

func ExampleOpen_connector() {
	f, err := groot.Open("../../testdata/simple.root")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	o, err := f.Get("tree")
	if err != nil {
		log.Fatal(err)
	}

	tree := o.(rtree.Tree)

	db := sql.OpenDB(rsqldrv.Connector(rtree.FileOf(tree)))
	defer db.Close()

	rows, err := db.Query("SELECT * FROM tree")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	type data struct {
		i32 int32
		f32 float32
		str string
	}

	n := 0
	for rows.Next() {
		var v data
		err = rows.Scan(&v.i32, &v.f32, &v.str)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("row[%d]: (%v, %v, %q)\n", n, v.i32, v.f32, v.str)
		n++
	}

	// Output:
	// row[0]: (1, 1.1, "uno")
	// row[1]: (2, 2.2, "dos")
	// row[2]: (3, 3.3, "tres")
	// row[3]: (4, 4.4, "quatro")
}
