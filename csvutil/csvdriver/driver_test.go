// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvdriver_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"go-hep.org/x/hep/csvutil/csvdriver"
)

func TestOpen(t *testing.T) {
	db, err := csvdriver.Conn{
		File:    "testdata/simple.csv",
		Comment: '#',
		Comma:   ';',
	}.Open()

	if err != nil {
		t.Errorf("error opening CSV file: %v\n", err)
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		t.Errorf("error starting tx: %v\n", err)
		return
	}
	defer tx.Commit()

	var done = make(chan error)
	go func() {
		done <- db.Ping()
	}()

	select {
	case <-time.After(2 * time.Second):
		t.Fatalf("ping timeout")
	case err := <-done:
		if err != nil {
			t.Fatalf("error pinging db: %v\n", err)
		}
	}

	rows, err := tx.Query("select var1, var2, var3 from csv order by id();")
	if err != nil {
		t.Errorf("error querying db: %v\n", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			i int64
			f float64
			s string
		)
		err = rows.Scan(&i, &f, &s)
		if err != nil {
			t.Errorf("error scanning db: %v\n", err)
			return
		}
		fmt.Printf("i=%v f=%v s=%q\n", i, f, s)
	}

	err = rows.Close()
	if err != nil {
		t.Errorf("error closing rows: %v\n", err)
		return
	}

	err = db.Close()
	if err != nil {
		t.Errorf("error closing db: %v\n", err)
		return
	}
}

func TestOpenName(t *testing.T) {
	db, err := sql.Open("csv", "testdata/simple-noheaders.csv")
	if err != nil {
		t.Errorf("error opening CSV file: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Errorf("error pinging db: %v\n", err)
	}
}

func TestQL(t *testing.T) {
	db, err := sql.Open("ql", "memory://out-create-ql.csv")
	if err != nil {
		t.Fatalf("error creating CSV-QL file: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("error pinging db: %v\n", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("error starting transaction: %v\n", err)
	}
	defer tx.Commit()

	_, err = tx.Exec("create table csv (var1 int64, var2 float64, var3 string);")
	if err != nil {
		t.Fatalf("error creating table: %v\n", err)
	}

	for i := 0; i < 10; i++ {
		f := float64(i)
		s := fmt.Sprintf("str-%d", i)
		_, err = tx.Exec("insert into csv values($1,$2,$3);", i, f, s)
		if err != nil {
			t.Fatalf("error inserting row %d: %v\n", i+1, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		t.Fatalf("error committing transaction: %v\n", err)
	}
}

func TestCreate(t *testing.T) {
	const fname = "testdata/out-create.csv"
	defer os.Remove(fname)

	db, err := csvdriver.Create(fname)
	if err != nil {
		t.Fatalf("error creating CSV file: %v\n", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("error pinging db: %v\n", err)
	}

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("error starting transaction: %v\n", err)
	}
	defer tx.Commit()

	_, err = tx.Exec("create table csv (var1 int64, var2 float64, var3 string);")
	if err != nil {
		t.Fatalf("error creating table: %v\n", err)
	}

	for i := 0; i < 10; i++ {
		f := float64(i)
		s := fmt.Sprintf("str-%d", i)
		_, err = tx.Exec("insert into csv values($1,$2,$3);", i, f, s)
		if err != nil {
			t.Fatalf("error inserting row %d: %v\n", i+1, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		t.Fatalf("error committing transaction: %v\n", err)
	}
}
