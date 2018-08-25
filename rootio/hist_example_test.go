package rootio_test

import (
	"fmt"
	"log"
	"os"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hbook/rootcnv"
	"go-hep.org/x/hep/rootio"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func ExampleCreate_histo1D() {
	const fname = "h1d_example.root"
	defer os.Remove(fname)

	f, err := rootio.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	const npoints = 10000

	// Create a normal distribution.
	dist := distuv.Normal{
		Mu:    0,
		Sigma: 1,
		Src:   rand.New(rand.NewSource(0)),
	}

	// Draw some random values from the standard
	// normal distribution.
	h := hbook.NewH1D(20, -4, +4)
	for i := 0; i < npoints; i++ {
		v := dist.Rand()
		h.Fill(v, 1)
	}
	h.Fill(-10, 1) // fill underflow
	h.Fill(-20, 2)
	h.Fill(+10, 3) // fill overflow

	fmt.Printf("original histo:\n")
	fmt.Printf("w-mean:    %.7f\n", h.XMean())
	fmt.Printf("w-rms:     %.7f\n", h.XRMS())

	hroot := rootio.NewH1DFrom(h)

	err = f.Put("h1", hroot)
	if err != nil {
		log.Fatal(err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing ROOT file: %v", err)
	}

	r, err := rootio.Open(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	robj, err := r.Get("h1")
	if err != nil {
		log.Fatal(err)
	}

	hr, err := rootcnv.H1D(robj.(rootio.H1))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nhisto read back:\n")
	fmt.Printf("r-mean:    %.7f\n", hr.XMean())
	fmt.Printf("r-rms:     %.7f\n", hr.XRMS())

	// Output:
	// original histo:
	// w-mean:    0.0023919
	// w-rms:     1.0628679
	//
	// histo read back:
	// r-mean:    0.0023919
	// r-rms:     1.0628679
}
