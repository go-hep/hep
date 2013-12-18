package main

import (
	"fmt"

	"github.com/go-hep/croot"
)

type Event struct {
	I int64
	F float64
}

func main() {
	const fname = "test-small.root"
	const evtmax = 100
	const splitlevel = 32
	const bufsiz = 32000
	const compress = 1
	const netopt = 0

	f, err := croot.OpenFile(fname, "recreate", "small event file", compress, netopt)
	if err != nil {
		panic(err.Error())
	}

	// create a tree
	tree := croot.NewTree("tree", "tree", splitlevel)

	e := Event{}

	_, err = tree.Branch2("Int64", &e.I, "Int64/L", bufsiz)
	if err != nil {
		panic(err.Error())
	}

	_, err = tree.Branch2("Float64", &e.F, "Float64/D", bufsiz)
	if err != nil {
		panic(err.Error())
	}

	// fill some events with random numbers
	for iev := int64(0); iev != evtmax; iev++ {
		if iev%1000 == 0 {
			fmt.Printf(":: processing event %d...\n", iev)
		}

		e.I = iev
		e.F = float64(iev)

		_, err = tree.Fill()
		if err != nil {
			panic(err.Error())
		}
	}
	f.Write("", 0, 0)
	f.Close("")

}

// EOF
