package hepmc_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/go-hep/hepmc"
)

//const fname = "testdata/test.hepmc"
//const fname = "testdata/one.hepmc"

func TestEventRW(t *testing.T) {

	for _, table := range []struct {
		fname    string
		outfname string
		nevts    int
	}{
		{"testdata/small.hepmc", "out.small.hepmc", 1},

		{"testdata/test.hepmc", "out.hepmc", 6},
		{"out.hepmc", "rb.out.hepmc", 6},
	} {
		test_evt_rw(t, table.fname, table.outfname, table.nevts)
	}

	// clean-up
	for _, fname := range []string{"out.small.hepmc", "out.hepmc", "rb.out.hepmc"} {
		err := os.Remove(fname)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func test_evt_rw(t *testing.T, fname, outfname string, nevts int) {
	f, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	dec := hepmc.NewDecoder(f)
	if dec == nil {
		t.Fatal(fmt.Errorf("hepmc.decoder: nil decoder"))
	}

	const NEVTS = 10
	evts := make([]*hepmc.Event, 0, NEVTS)
	for ievt := 0; ievt < NEVTS; ievt++ {
		var evt hepmc.Event
		err = dec.Decode(&evt)
		if err != nil {
			if err == io.EOF && ievt == nevts {
				break
			}
			t.Fatalf("file: %s. ievt=%d err=%v\n", fname, ievt, err)
		}
		evts = append(evts, &evt)
		defer hepmc.Delete(&evt)
	}

	o, err := os.Create(outfname)
	if err != nil {
		t.Fatal(err)
	}
	defer o.Close()

	enc := hepmc.NewEncoder(o)
	if enc == nil {
		t.Fatal(fmt.Errorf("hepmc.encoder: nil encoder"))
	}
	for _, evt := range evts {
		err = enc.Encode(evt)
		if err != nil {
			t.Fatalf("file: %s. err=%v\n", fname, err)
		}
	}
	err = enc.Close()
	if err != nil {
		t.Fatalf("file: %s. err=%v\n", fname, err)
	}

	// test output files
	cmd := exec.Command("diff", "-urN", fname, outfname)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("file: %s. err=%v\n", fname, err)
	}

}

// EOF
