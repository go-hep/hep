package rio_test

import (
	"io"
	"reflect"
	"testing"

	"github.com/go-hep/rio"
)

func TestStreamOpen(t *testing.T) {
	const fname = "testdata/c_sim.slcio"
	f, err := rio.Open(fname)
	if err != nil {
		t.Fatalf("could not open [%s]: %v", fname)
	}
	defer f.Close()

	if f.Name() != fname {
		t.Fatalf("rio.Stream.Name: expected [%s]. got [%s]", fname, f.Name())
	}

	if f.FileName() != fname {
		t.Fatalf("rio.Stream.FileName: expected [%s]. got [%s]", fname, f.FileName())
	}

	fi, err := f.Mode()
	if err != nil {
		t.Fatalf("could not retrieve stream mode: %v", err)
	}

	if !fi.IsRegular() {
		t.Fatalf("rio.Stream.Mode: expected regular file")
	}

	if f.CurPos() != 0 {
		t.Fatalf("expected pos=%v. got=%v", 0, f.CurPos())
	}
}

func TestReadLcio(t *testing.T) {
	const fname = "testdata/c_sim.slcio"

	f, err := rio.Open(fname)
	if err != nil {
		t.Fatalf("could not open [%s]: %v", fname)
	}
	defer f.Close()

	type LCRunHeader struct {
		RunNbr   int32
		Detector string
		Descr    string
		SubDets  []string
		//Params   LCParameters
	}
	
	var runhdr LCRunHeader
	runhdr.RunNbr = 42

	rec := f.Record("LCRunHeader")
	if !f.HasRecord("LCRunHeader") {
		t.Fatalf("expected stream to have LCRunHeader record")
	}
	if rec.Unpack() {
		t.Fatalf("expected record to NOT unpack by default")
	}
	if rec.Name() != "LCRunHeader" {
		t.Fatalf("expected record name=[%s]. got=[%s]", "LCRunHeader", rec.Name())
	}

	rec.SetUnpack(true)
	if !rec.Unpack() {
		t.Fatalf("expected record to unpack now")
	}

	err = rec.Connect("RunHeader", &runhdr)
	if err != nil {
		t.Fatalf("error connecting [RunHeader]: %v", err)
	}

	for	nrecs := 0; nrecs < 100; nrecs++ {
		rec, err := f.ReadRecord()
		if err != nil {
			if err == io.EOF && nrecs == 10 {
				break
			}
			t.Fatalf("error reading record: %v (nrecs=%d)", err, nrecs)
		}

		if rec == nil {
			t.Fatalf("got nil record! (nrecs=%d)", nrecs)
		}

		if rec.Name() != "LCRunHeader" {
			t.Fatalf("expected record name=[%s]. got=[%s]. (nrecs=%d)",
				"LCRunHeader",
				rec.Name(),
				nrecs,
			)
		}

		if int(runhdr.RunNbr) != nrecs {
			t.Fatalf("expected runnbr=%d. got=%d.", nrecs, runhdr.RunNbr)
		}
		if runhdr.Detector != "D09TileHcal" {
			t.Fatalf("expected detector=[%s]. got=[%s]. (nrecs=%d)",
				"D09TileHcal",
				runhdr.Detector,
				nrecs,
			)
		}
		subdets := []string{"ECAL007", "TPC4711"}
		if !reflect.DeepEqual(runhdr.SubDets, subdets) {
			t.Fatalf("expected subdets=%v. got=%v (nrecs=%d)",
				subdets,
				runhdr.SubDets,
			)
		}
	}
}
