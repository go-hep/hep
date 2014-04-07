rio
===

[![Build Status](https://drone.io/github.com/go-hep/rio/status.png)](https://drone.io/github.com/go-hep/rio/latest)

`rio` is a Record-oriented I/O library, modeled after SIO (Serial I/O).

## Installation

```sh
$ go get github.com/go-hep/rio
```

## Documentation

The documentation is browsable at godoc.org:
 http://godoc.org/github.com/go-hep/rio

## Example

```go
package main

import (
	"fmt"
	"io"

	"github.com/go-hep/rio"
)

type LCRunHeader struct {
	RunNbr   int32
	Detector string
	Descr    string
	SubDets  []string
}

func main() {
	fname := "c_sim.slcio"
	fmt.Printf(">>> opening [%s]\n", fname)

	f, err := rio.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var runhdr LCRunHeader
	runhdr.RunNbr = 42
	rec := f.Record("LCRunHeader")
	rec.SetUnpack(true)
	rec.Connect("RunHeader", &runhdr)

	fmt.Printf("::: [%s]\n", f.Name())
	fmt.Printf("::: [%s]\n", f.FileName())
	for {
		fmt.Printf("-----------------------------------------------\n")
		var rec *rio.Record

		rec, err = f.ReadRecord()
		fmt.Printf("***\n")
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}
		if rec == nil {
			fmt.Printf("** no record!\n")
			break
		}

		fmt.Printf(">> record: %s\n", rec.Name())
		if rec.Name() == "LCRunHeader" {
			fmt.Printf("runnbr: %d\n", runhdr.RunNbr)
			fmt.Printf("det:    %q\n", runhdr.Detector)
			dets := "["
			for i, det := range runhdr.SubDets {
				dets += fmt.Sprintf("%q", det)
				if i+1 != len(runhdr.SubDets) {
					dets += ", "
				}
			}
			dets += "]"
			fmt.Printf("dets:   %s\n", dets)
		}
	}
}
```

## TODO

- implement read/write pointers
- implement buffered i/o
- handle big files (ie: file offsets as `int64`)

## Bibliography

- `SIO`: http://www-sldnt.slac.stanford.edu/nld/new/Docs/FileFormats/sio.pdf
