// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main // import "go-hep.org/x/hep/sio/cmd/sio-ls-records"

import (
	"flag"
	"fmt"
	"os"

	"go-hep.org/x/hep/sio"
)

func main() {
	var fname string

	flag.Parse()

	if flag.NArg() > 0 {
		fname = flag.Arg(0)
	}

	fmt.Printf("::: inspecting file [%s]...\n", fname)
	if fname == "" {
		flag.Usage()
		os.Exit(1)
	}

	f, err := sio.Open(fname)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	_, _ = f.ReadRecord()

	_, err = f.Seek(0, 0)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
	}

	for _, rec := range f.Records() {
		fmt.Printf(" -> %v\n", rec.Name())
	}

	fmt.Printf("::: inspecting file [%s]... [done]\n", fname)
}
