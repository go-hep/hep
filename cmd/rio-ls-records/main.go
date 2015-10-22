package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/go-hep/rio"
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

	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	r, err := rio.NewReader(f)
	if err != nil {
		fmt.Printf("*** error creating rio.Reader: %v\n", err)
		os.Exit(2)
	}

	scan := rio.NewScanner(r)
	for scan.Scan() {
		// scans through the whole stream
		err = scan.Err()
		if err != nil {
			break
		}
		rec := scan.Record()
		fmt.Printf(" -> %v\n", rec.Name())
	}
	err = scan.Err()
	if err != nil {
		fmt.Printf("*** error: %v\n", err)
		os.Exit(2)
	}

	fmt.Printf("::: inspecting file [%s]... [done]\n", fname)
}
