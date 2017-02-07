package main

import (
	"flag"
	"fmt"
	"os"

	"go-hep.org/x/hep/slha"
)

func handle(err error) {
	if err != nil {
		printf("**error: %v\n", err)
		panic(err)
	}
}

func printf(format string, args ...interface{}) (int, error) {
	return fmt.Fprintf(os.Stderr, format, args...)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, " $ %s <path-to-SLHA-file>\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() <= 0 {
		printf("**error** need an input file name\n")
		flag.Usage()
		os.Exit(1)
	}

	fname := flag.Arg(0)

	f, err := os.Open(fname)
	if err != nil {
		printf("could not open file [%s]: %v\n", fname, err)
		os.Exit(1)
	}
	defer f.Close()

	data, err := slha.Decode(f)
	if err != nil {
		printf("could not decode file [%s]: %v\n", fname, err)
		os.Exit(1)
	}

	spinfo := data.Blocks.Get("SPINFO")
	if spinfo != nil {
		value, err := spinfo.Get(1)
		handle(err)
		fmt.Printf("spinfo: %s -- %q\n", value.Interface(), value.Comment())
	}

	modsel := data.Blocks.Get("MODSEL")
	if modsel != nil {
		value, err := modsel.Get(1)
		handle(err)
		fmt.Printf("modsel: %d -- %q\n", value.Interface(), value.Comment())
	}

	mass := data.Blocks.Get("MASS")
	if mass != nil {
		value, err := mass.Get(5)
		handle(err)
		fmt.Printf("mass[pdgid=5]: %v -- %q\n", value.Interface(), value.Comment())
	}

	nmix := data.Blocks.Get("NMIX")
	if nmix != nil {
		value, err := nmix.Get(1, 2)
		handle(err)
		fmt.Printf("nmix[1,2] = %v -- %q\n", value.Interface(), value.Comment())
	}
}

// Output:
// spinfo: SOFTSUSY -- "spectrum calculator"
// modsel: 1 -- "sugra"
// mass[pdgid=5]: 4.88991651 -- "b-quark pole mass calculated from mb(mb)_Msbar"
// nmix[1,2] = -0.0531103553 -- "N_12"
