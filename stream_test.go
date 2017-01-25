package sio_test

import (
	//"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/go-hep/sio"
)

type RunHeader struct {
	RunNbr   int32
	Detector string
	Descr    string
	SubDets  []string
	//Params   Parameters

	Ints   []int64
	Floats []float64
}

func TestStreamOpen(t *testing.T) {
	const fname = "testdata/runhdr.sio"
	f, err := sio.Open(fname)
	if err != nil {
		t.Fatalf("could not open [%s]: %v", fname, err)
	}
	defer f.Close()

	if f.Name() != fname {
		t.Fatalf("sio.Stream.Name: expected [%s]. got [%s]", fname, f.Name())
	}

	if f.FileName() != fname {
		t.Fatalf("sio.Stream.FileName: expected [%s]. got [%s]", fname, f.FileName())
	}

	fi, err := f.Mode()
	if err != nil {
		t.Fatalf("could not retrieve stream mode: %v", err)
	}

	if !fi.IsRegular() {
		t.Fatalf("sio.Stream.Mode: expected regular file")
	}

	if f.CurPos() != 0 {
		t.Fatalf("expected pos=%v. got=%v", 0, f.CurPos())
	}
}

func TestStreamCreate(t *testing.T) {
	const fname = "testdata/out.sio"
	defer os.RemoveAll(fname)

	f, err := sio.Create(fname)
	if err != nil {
		t.Fatalf("could not create [%s]: %v", fname, err)
	}

	if f.Name() != fname {
		t.Fatalf("sio.Stream.Name: expected [%s]. got [%s]", fname, f.Name())
	}

	if f.FileName() != fname {
		t.Fatalf("sio.Stream.FileName: expected [%s]. got [%s]", fname, f.FileName())
	}

	fi, err := f.Mode()
	if err != nil {
		t.Fatalf("could not retrieve stream mode: %v", err)
	}

	if !fi.IsRegular() {
		t.Fatalf("sio.Stream.Mode: expected regular file")
	}

	if f.CurPos() != 0 {
		t.Fatalf("expected pos=%v. got=%v", 0, f.CurPos())
	}
}

func TestReadRunHeader(t *testing.T) {
	testReadStream(t, "testdata/runhdr.sio")
}

func TestReadRunHeaderCompr(t *testing.T) {
	t.Skip("sio read compression disabled")
	testReadStream(t, "testdata/runhdr-compr.sio")
}

func TestWriteRunHeader(t *testing.T) {
	const fname = "testdata/out.sio"
	defer os.RemoveAll(fname)
	testWriteStream(t, fname)
}

func TestReadWrite(t *testing.T) {
	const fname = "testdata/rw.sio"
	defer os.RemoveAll(fname)
	testWriteStream(t, fname)
	testReadStream(t, fname)
}

func testReadStream(t *testing.T, fname string) {

	f, err := sio.Open(fname)
	if err != nil {
		t.Fatalf("could not open [%s]: %v", fname, err)
	}
	defer f.Close()

	runhdr := RunHeader{
		RunNbr:   42,
		Detector: "---",
		Descr:    "---",
		SubDets:  []string{},
		Floats:   []float64{},
		Ints:     []int64{},
	}

	rec := f.Record("RioRunHeader")
	if !f.HasRecord("RioRunHeader") {
		t.Fatalf("expected stream to have LCRunHeader record")
	}
	if rec.Unpack() {
		t.Fatalf("expected record to NOT unpack by default")
	}
	if rec.Name() != "RioRunHeader" {
		t.Fatalf("expected record name=[%s]. got=[%s]", "RioRunHeader", rec.Name())
	}

	rec.SetUnpack(true)
	if !rec.Unpack() {
		t.Fatalf("expected record to unpack now")
	}

	err = rec.Connect("RunHeader", &runhdr)
	if err != nil {
		t.Fatalf("error connecting [RunHeader]: %v", err)
	}

	for nrecs := 0; nrecs < 100; nrecs++ {
		//fmt.Printf("::: irec=%d, fname=%q\n", nrecs, fname)
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

		if rec.Name() != "RioRunHeader" {
			t.Fatalf("expected record name=[%s]. got=[%s]. (nrecs=%d)",
				"RioRunHeader",
				rec.Name(),
				nrecs,
			)
		}

		if int(runhdr.RunNbr) != nrecs {
			t.Fatalf("expected runnbr=%d. got=%d.", nrecs, runhdr.RunNbr)
		}
		if runhdr.Detector != "MyDetector" {
			t.Fatalf("expected detector=[%s]. got=[%s]. (nrecs=%d)",
				"MyDetector",
				runhdr.Detector,
				nrecs,
			)
		}
		if runhdr.Descr != "dummy run number" {
			t.Fatalf("expected descr=[%s]. got=[%s]. (nrecs=%d)",
				"dummy run number",
				runhdr.Descr,
				nrecs,
			)
		}
		subdets := []string{"subdet 0", "subdet 1"}
		if !reflect.DeepEqual(runhdr.SubDets, subdets) {
			t.Fatalf("expected subdets=%v. got=%v (nrecs=%d)",
				subdets,
				runhdr.SubDets,
				nrecs,
			)
		}

		floats := []float64{
			float64(nrecs) + 100,
			float64(nrecs) + 200,
			float64(nrecs) + 300,
		}
		if !reflect.DeepEqual(runhdr.Floats, floats) {
			t.Fatalf("expected floats=%v. got=%v (nrecs=%d)",
				floats,
				runhdr.Floats,
				nrecs,
			)
		}

		ints := []int64{
			int64(nrecs) + 100,
			int64(nrecs) + 200,
			int64(nrecs) + 300,
		}
		if !reflect.DeepEqual(runhdr.Ints, ints) {
			t.Fatalf("expected ints=%v. got=%v (nrecs=%d)",
				floats,
				runhdr.Floats,
				nrecs,
			)
		}
	}
}

func testWriteStream(t *testing.T, fname string) {
	f, err := sio.Create(fname)
	if err != nil {
		t.Fatalf("could not create [%s]: %v", fname, err)
	}

	defer f.Close()

	var runhdr RunHeader
	runhdr.RunNbr = 42

	rec := f.Record("RioRunHeader")
	if rec == nil {
		t.Fatalf("could not create record [RioRunHeader]")
	}
	rec.SetUnpack(true)
	if !rec.Unpack() {
		t.Fatalf("expected record to unpack now")
	}

	err = rec.Connect("RunHeader", &runhdr)
	if err != nil {
		t.Fatalf("error connecting [RunHeader]: %v", err)
	}

	for irec := 0; irec < 10; irec++ {
		runhdr = RunHeader{
			RunNbr:   int32(irec),
			Detector: "MyDetector",
			Descr:    "dummy run number",
			SubDets:  []string{"subdet 0", "subdet 1"},
			Floats: []float64{
				float64(irec) + 100,
				float64(irec) + 200,
				float64(irec) + 300,
			},
			Ints: []int64{
				int64(irec) + 100,
				int64(irec) + 200,
				int64(irec) + 300,
			},
		}
		err = f.WriteRecord(rec)
		if err != nil {
			t.Fatalf("error writing record: %v (irec=%d)", err, irec)
		}

		err = f.Sync()
		if err != nil {
			t.Fatalf("error flushing record: %v (irec=%d)", err, irec)
		}
	}
}
