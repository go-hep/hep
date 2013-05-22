package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-hep/hepmc"
)

func main() {
	var err error
	var r io.Reader
	var w io.Writer = os.Stdout

	switch len(os.Args) {
	case 1:
		r = os.Stdin
	case 2:
		r, err = os.Open(os.Args[1])
		if err != nil {
			fmt.Printf("**error: %v\n", err)
			os.Exit(1)
		}
	default:
	}

	dec := hepmc.NewDecoder(r)
	if dec == nil {
		fmt.Printf("**error: nil decoder\n")
		os.Exit(1)
	}

	for {
		var evt hepmc.Event
		err = dec.Decode(&evt)
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("**error: %v\n", err)
			os.Exit(1)
		}
		err = evt.Print(w)
		if err != nil {
			fmt.Printf("**error: %v\n", err)
			os.Exit(1)
		}
	}
}

// EOF
