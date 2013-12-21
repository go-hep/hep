// root-ls dumps the content of a ROOT file
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"

	"github.com/go-hep/rootio"
)

var g_fname = flag.String("f", "", "path to the ROOT to inspect")
var g_prof = flag.String("profile", "", "filename of cpuprofile")

func main() {
	flag.Parse()

	if *g_prof != "" {
		f, err := os.Create(*g_prof)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *g_fname == "" {
		fmt.Fprintf(os.Stderr, "**error** you need to give a ROOT file\n")
		flag.Usage()
		os.Exit(1)
	}

	f, err := rootio.Open(*g_fname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error** %v\n", err)
		os.Exit(1)
	}

	for _, k := range f.Keys() {
		fmt.Printf("%-8s %-40s %s\n", k.ClassName(), k.Name(), k.Title())
	}
}
