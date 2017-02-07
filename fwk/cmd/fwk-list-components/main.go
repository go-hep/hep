package main

import (
	"fmt"

	"github.com/go-hep/fwk"
)

func main() {
	comps := fwk.Registry()
	fmt.Printf("::: components... (%d)\n", len(comps))
	for i, c := range comps {
		fmt.Printf("[%04d/%04d] %s\n", i, len(comps), c)
	}
}
