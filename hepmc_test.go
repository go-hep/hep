package hepmc_test

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/go-hep/hepmc"
)

const fname = "testdata/test.hepmc"

//const fname = "testdata/one.hepmc"

func TestEventRW(t *testing.T) {

	for _, table := range []struct {
		fname    string
		outfname string
	}{
		{fname, "out.hepmc"},
		{"out.hepmc", "rb.out.hepmc"},
	} {
		test_evt_rw(t, table.fname, table.outfname)
	}

	// clean-up
	for _, fname := range []string{"out.hepmc", "rb.out.hepmc"} {
		err := os.Remove(fname)
		if err != nil {
			t.Fatal(err)
		}
		err = os.Remove(fname + ".todiff")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func test_evt_rw(t *testing.T, fname, outfname string) {
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
	evts := make([]*hepmc.Event, 6)
	for ievt := 0; ievt < NEVTS; ievt++ {
		var evt hepmc.Event
		err = dec.Decode(&evt)
		if err != nil {
			if err == io.EOF && ievt == 6 {
				break
			}
			t.Fatal(err)
		}
		evts[ievt] = &evt
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
			t.Fatal(err)
		}
	}
	err = enc.Close()
	if err != nil {
		t.Fatal(err)
	}

	// test output files
	cmd := exec.Command(
		"/bin/sh", "-c",
		fmt.Sprintf("(cat %s | sort >| %s.todiff)", outfname, outfname),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

	cmd = exec.Command(
		"diff", "-urN",
		"testdata/ref.hepmc.todiff",
		fmt.Sprintf("%s.todiff", outfname),
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatal(err)
	}

}

// EOF
