package lhef_test

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/go-hep/lhef"
)

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
		fmt.Printf("===[%d]===\n", i)
		evt, err := dec.Decode()
		if err == io.EOF {
			fmt.Printf("** EOF **\n")
			break
		}
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("evt: %v\n", *evt)
	}
}

// EOF
