// Copyright 2016 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package csvdriver

import (
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-hep/csvutil"
)

func (conn *csvConn) importCSV() error {
	tbl, err := csvutil.Open(conn.cfg.File)
	if err != nil {
		return err
	}
	defer tbl.Close()
	tbl.Reader.Comma = conn.cfg.Comma
	tbl.Reader.Comment = conn.cfg.Comment

	schema, err := inferSchemaFromTable(tbl)
	if err != nil {
		return err
	}

	tx, err := conn.Begin()
	if err != nil {
		log.Fatalf("tx-err: %v\n", err)
		return err
	}
	defer tx.Commit()

	_, err = conn.Exec("create table csv ("+schema.Decl()+")", nil)
	if err != nil {
		log.Fatalf("create-err: %v\n", err)
		return err
	}

	rows, err := tbl.ReadRows(0, -1)
	if err != nil {
		return err
	}
	defer rows.Close()

	vargs, pargs := schema.Args()
	def := schema.Def()
	insert := "insert into csv values(" + def + ");"
	for rows.Next() {
		err = rows.Scan(pargs...)
		if err != nil {
			return err
		}
		for i, arg := range pargs {
			vargs[i] = reflect.ValueOf(arg).Elem().Interface()
		}
		_, err = conn.Exec(insert, vargs)
		if err != nil {
			return err
		}
	}

	err = rows.Err()
	if err == io.EOF {
		err = nil
	}
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func inferSchemaFromTable(tbl *csvutil.Table) (schemaType, error) {
	rows, err := tbl.ReadRows(0, 1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, rows.Err()
	}

	return inferSchemaFromFields(rows.Fields())
}

func inferSchemaFromFields(fields []string) (schemaType, error) {
	schema := make(schemaType, len(fields))
	for i, field := range fields {
		var err error
		_, err = strconv.ParseInt(field, 10, 64)
		if err == nil {
			schema[i] = reflect.ValueOf(int64(0))
			continue
		}

		_, err = strconv.ParseFloat(field, 64)
		if err == nil {
			schema[i] = reflect.ValueOf(float64(0))
			continue
		}

		schema[i] = reflect.ValueOf("")
	}
	return schema, nil
}

type schemaType []reflect.Value

func (st *schemaType) Decl() string {
	o := make([]string, 0, len(*st))
	for i, v := range *st {
		t := v.Type().Kind().String()
		o = append(o, fmt.Sprintf("var%d %s", i+1, t))
	}
	return strings.Join(o, ", ")
}

func (st *schemaType) Args() ([]driver.Value, []interface{}) {
	vargs := make([]driver.Value, len(*st))
	pargs := make([]interface{}, len(*st))
	for i, v := range *st {
		ptr := reflect.New(v.Type())
		vargs[i] = ptr.Elem().Interface()
		pargs[i] = ptr.Interface()
	}
	return vargs, pargs
}

func (st *schemaType) Def() string {
	o := make([]string, len(*st))
	for i := range *st {
		o[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(o, ", ")
}
