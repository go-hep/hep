package lhef_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/go-hep/lhef"
)

const r_debug = false
const ifname = "testdata/ttbar.lhe"

func TestLhefReading(t *testing.T) {
	f, err := os.Open(ifname)
	if err != nil {
		t.Error(err)
	}

	dec, err := lhef.NewDecoder(f)
	if err != nil {
		t.Error(err)
	}

	for i := 0; ; i++ {
		if r_debug {
			fmt.Printf("===[%d]===\n", i)
		}
		evt, err := dec.Decode()
		if err == io.EOF {
			if r_debug {
				fmt.Printf("** EOF **\n")
			}
			break
		}
		if err != nil {
			t.Error(err)
		}
		if r_debug {
			fmt.Printf("evt: %v\n", *evt)
		}
	}
}

// EOF
