csvutil
=======

[![GoDoc](https://godoc.org/github.com/go-hep/csvutil?status.svg)](https://godoc.org/github.com/go-hep/csvutil)

`csvutil` is a set of types and funcs to deal with CSV data files in a somewhat convenient way.

## Installation

```sh
$> go get github.com/go-hep/csvutil
```

## Documentation

Documentation is available on [godoc](https://godoc.org):

[godoc.org/github.com/go-hep/csvutil](https://godoc.org/github.com/go-hep/csvutil)

## Example

```go
package main

import (
	"io"
	"log"
	"os"

	"github.com/go-hep/csvutil"
)

func main() {
	fname := "testdata/simple.csv"
	tbl, err := csvutil.Open(fname)
	if err != nil {
		log.Fatalf("could not open %s: %v\n", fname, err)
	}
	defer tbl.Close()
	tbl.Reader.Comma = ';'
	tbl.Reader.Comment = '#'

	rows, err := tbl.ReadRows(0, 10)
	if err != nil {
		log.Fatalf("could read rows [0, 10): %v\n", err)
	}
	defer rows.Close()

	irow := 0
	for rows.Next() {
		data := struct {
			I int
			F float64
			S string
		}{}
		err = rows.Scan(&data)
		if err != nil {
			log.Fatalf("error reading row %d: %v\n", irow, err)
		}
	}
	err = rows.Err()
	if err != nil && err != io.EOF {
		log.Fatalf("error: %v\n", err)
	}
}
```
