// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sio_test

import (
	"io"
	"os"
	"reflect"
	"testing"

	"go-hep.org/x/hep/sio"
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

type T1 struct {
	Name string
	T2   *T2
	T3   *T2
	T4   *T2
	T5   *T5
	T6   *T2
	T7   *T2
}

func (t1 *T1) MarshalSio(w sio.Writer) error {
	var err error
	err = sio.Marshal(w, t1.Name)
	if err != nil {
		return err
	}

	err = w.Pointer(&t1.T2)
	if err != nil {
		return err
	}

	err = w.Pointer(&t1.T3)
	if err != nil {
		return err
	}

	err = w.Pointer(&t1.T4)
	if err != nil {
		return err
	}

	err = w.Pointer(&t1.T5)
	if err != nil {
		return err
	}

	err = w.Pointer(&t1.T6)
	if err != nil {
		return err
	}

	err = w.Pointer(&t1.T7)
	if err != nil {
		return err
	}

	err = w.Tag(t1)
	if err != nil {
		return err
	}

	return err
}

func (t1 *T1) UnmarshalSio(r sio.Reader) error {
	var err error
	err = sio.Unmarshal(r, &t1.Name)
	if err != nil {
		return err
	}

	err = r.Pointer(&t1.T2)
	if err != nil {
		return err
	}

	err = r.Pointer(&t1.T3)
	if err != nil {
		return err
	}

	err = r.Pointer(&t1.T4)
	if err != nil {
		return err
	}

	err = r.Pointer(&t1.T5)
	if err != nil {
		return err
	}

	err = r.Pointer(&t1.T6)
	if err != nil {
		return err
	}

	err = r.Pointer(&t1.T7)
	if err != nil {
		return err
	}

	err = r.Tag(t1)
	if err != nil {
		return err
	}

	return nil
}

type T2 struct {
	Name string
}

func (t2 *T2) MarshalSio(w sio.Writer) error {
	var err error
	err = sio.Marshal(w, t2.Name)
	if err != nil {
		return err
	}

	err = w.Tag(t2)
	if err != nil {
		return err
	}

	return err
}

func (t2 *T2) UnmarshalSio(r sio.Reader) error {
	var err error
	err = sio.Unmarshal(r, &t2.Name)
	if err != nil {
		return err
	}

	err = r.Tag(t2)
	if err != nil {
		return err
	}

	return nil
}

type T5 struct {
	Name string
}

func (t2 *T5) MarshalSio(w sio.Writer) error {
	var err error
	err = sio.Marshal(w, t2.Name)
	if err != nil {
		return err
	}
	// no ptag
	return err
}

func (t2 *T5) UnmarshalSio(r sio.Reader) error {
	var err error
	err = sio.Unmarshal(r, &t2.Name)
	if err != nil {
		return err
	}
	// no ptag
	return nil
}

var (
	_ sio.Codec = (*T1)(nil)
	_ sio.Codec = (*T2)(nil)
	_ sio.Codec = (*T5)(nil)
)

func TestPointerStream(t *testing.T) {
	const name = "testdata/ptr.sio"
	func() {
		f, err := sio.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		t7 := T2{Name: "t7"}
		t6 := T2{Name: "t6"}
		t5 := T5{Name: "t5"}
		t4 := T2{Name: "t4"}
		t3 := T2{Name: "t3"}
		t2 := T2{Name: "t2"}
		t1 := T1{
			Name: "t1",
			T2:   &t2, T3: &t3, T4: &t4,
			T5: &t5, T6: &t6, T7: &t7,
		}
		rec := f.Record("Data")
		rec.SetUnpack(true)

		for _, v := range []struct {
			n   string
			ptr interface{}
		}{
			{"T1", &t1},
			{"T2", &t2},
			{"T3", &t3},
			{"T4", &t4},
			{"T5", &t5},
			{"T6", &t6},
			// {"T7", &t7}, // drop it
		} {
			err = rec.Connect(v.n, v.ptr)
			if err != nil {
				t.Fatalf("error connecting %q: %v", v.n, err)
			}
		}

		err = f.WriteRecord(rec)
		if err != nil {
			t.Fatalf("error writing record: %v", err)
		}
		err = f.Sync()
		if err != nil {
			t.Fatalf("error flushing record: %v", err)
		}
	}()

	func() {
		f, err := sio.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		var (
			t1 T1
			t2 T2
			t3 T2
			t4 T2
			t5 T5
			// t6 T2
			t7 T2
		)

		rec := f.Record("Data")
		rec.SetUnpack(true)

		for _, v := range []struct {
			n   string
			ptr interface{}
		}{
			{"T1", &t1},
			{"T2", &t2},
			{"T3", &t3},
			{"T4", &t4},
			{"T5", &t5},
			// {"T6",&t6}, // drop it
			{"T7", &t7},
		} {
			err = rec.Connect(v.n, v.ptr)
			if err != nil {
				t.Fatalf("error connecting %q: %v", v.n, err)
			}
		}

		rec, err = f.ReadRecord()
		if err != nil {
			t.Fatalf("error reading record: %v", err)
		}
		if !rec.Unpack() {
			t.Fatalf("error unpacking record")
		}

		if t1.Name != "t1" {
			t.Errorf("t1.Name = %q", t1.Name)
		}

		if t2.Name != "t2" {
			t.Errorf("t2.Name = %q", t2.Name)
		}

		if t3.Name != "t3" {
			t.Errorf("t3.Name = %q", t3.Name)
		}

		if t4.Name != "t4" {
			t.Errorf("t4.Name = %q", t4.Name)
		}

		if t5.Name != "t5" {
			t.Errorf("t5.Name = %q", t5.Name)
		}

		if t7.Name != "" {
			t.Errorf("t7.Name = %q", t7.Name)
		}

		for _, v := range []struct {
			n   string
			ptr *T2
		}{
			{"t2", t1.T2},
			{"t3", t1.T3},
			{"t4", t1.T4},
		} {
			if v.ptr == nil {
				t.Fatalf("t1.%s == nil", v.n)
			}

			if got, want := v.ptr.Name, v.n; got != want {
				t.Fatalf("t1.%s.Name=%q. want=%q", v.n, got, want)
			}
		}
		if t1.T5 != nil {
			t.Fatalf("t1.T5 = %v. want=nil", t1.T5)
		}

		if t1.T6 != nil {
			t.Fatalf("t1.T6 = %v. want=nil", t1.T6)
		}

		if t1.T7 != nil {
			t.Fatalf("t1.T7 = %v. want=nil", t1.T7)
		}
	}()
}

type C1 struct {
	Name string
	C2   *C2
}

func (c1 *C1) MarshalSio(w sio.Writer) error {
	var err error
	err = sio.Marshal(w, c1.Name)
	if err != nil {
		return err
	}

	err = w.Pointer(&c1.C2)
	if err != nil {
		return err
	}

	err = w.Tag(c1)
	if err != nil {
		return err
	}

	return err
}

func (c1 *C1) UnmarshalSio(r sio.Reader) error {
	var err error
	err = sio.Unmarshal(r, &c1.Name)
	if err != nil {
		return err
	}

	err = r.Pointer(&c1.C2)
	if err != nil {
		return err
	}

	err = r.Tag(c1)
	if err != nil {
		return err
	}

	return nil
}

type C2 struct {
	Name string
	C3   *C3
}

func (c2 *C2) MarshalSio(w sio.Writer) error {
	var err error
	err = sio.Marshal(w, c2.Name)
	if err != nil {
		return err
	}

	err = w.Pointer(&c2.C3)
	if err != nil {
		return err
	}

	err = w.Tag(c2)
	if err != nil {
		return err
	}

	return err
}

func (c2 *C2) UnmarshalSio(r sio.Reader) error {
	var err error
	err = sio.Unmarshal(r, &c2.Name)
	if err != nil {
		return err
	}

	err = r.Pointer(&c2.C3)
	if err != nil {
		return err
	}

	err = r.Tag(c2)
	if err != nil {
		return err
	}

	return nil
}

type C3 struct {
	Name string
	C1   *C1
}

func (c3 *C3) MarshalSio(w sio.Writer) error {
	var err error
	err = sio.Marshal(w, c3.Name)
	if err != nil {
		return err
	}

	err = w.Pointer(&c3.C1)
	if err != nil {
		return err
	}

	err = w.Tag(c3)
	if err != nil {
		return err
	}

	return err
}

func (c3 *C3) UnmarshalSio(r sio.Reader) error {
	var err error
	err = sio.Unmarshal(r, &c3.Name)
	if err != nil {
		return err
	}

	err = r.Pointer(&c3.C1)
	if err != nil {
		return err
	}

	err = r.Tag(c3)
	if err != nil {
		return err
	}

	return nil
}

var (
	_ sio.Codec = (*C1)(nil)
	_ sio.Codec = (*C2)(nil)
	_ sio.Codec = (*C3)(nil)
)

func TestPointerCycleStream(t *testing.T) {
	const name = "testdata/cycle-ptr.sio"
	func() {
		f, err := sio.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		c3 := C3{Name: "c3"}
		c2 := C2{Name: "c2"}
		c1 := C1{Name: "c1", C2: &c2}
		c2.C3 = &c3
		c3.C1 = &c1
		rec := f.Record("Data")
		rec.SetUnpack(true)

		for _, v := range []struct {
			n   string
			ptr interface{}
		}{
			{"C1", &c1},
			{"C2", &c2},
			{"C3", &c3},
		} {
			err = rec.Connect(v.n, v.ptr)
			if err != nil {
				t.Fatalf("error connecting %q: %v", v.n, err)
			}
		}

		err = f.WriteRecord(rec)
		if err != nil {
			t.Fatalf("error writing record: %v", err)
		}
		err = f.Sync()
		if err != nil {
			t.Fatalf("error flushing record: %v", err)
		}
	}()

	func() {
		f, err := sio.Open(name)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		var (
			c1 C1
			c2 C2
			c3 C3
		)

		rec := f.Record("Data")
		rec.SetUnpack(true)

		for _, v := range []struct {
			n   string
			ptr interface{}
		}{
			{"C1", &c1},
			{"C2", &c2},
			{"C3", &c3},
		} {
			err = rec.Connect(v.n, v.ptr)
			if err != nil {
				t.Fatalf("error connecting %q: %v", v.n, err)
			}
		}

		rec, err = f.ReadRecord()
		if err != nil {
			t.Fatalf("error reading record: %v", err)
		}
		if !rec.Unpack() {
			t.Fatalf("error unpacking record")
		}

		if c1.Name != "c1" {
			t.Errorf("c1.Name = %q", c1.Name)
		}

		if c2.Name != "c2" {
			t.Errorf("c2.Name = %q", c2.Name)
		}

		if c3.Name != "c3" {
			t.Errorf("c3.Name = %q", c3.Name)
		}

		switch {
		case c1.C2 == nil:
			t.Errorf("c1.C2 == nil")
		case c1.C2.Name != "c2":
			t.Errorf("c1.C2.Name = %q", c1.C2.Name)
		case c1.C2 != &c2:
			t.Errorf("c1.C2 = %v", c1.C2)
		}

		switch {
		case c2.C3 == nil:
			t.Errorf("c2.C3 == nil")
		case c2.C3.Name != "c3":
			t.Errorf("c2.C3.Name = %q", c2.C3.Name)
		case c2.C3 != &c3:
			t.Errorf("c2.C3 = %v", c2.C3)
		}

		switch {
		case c3.C1 == nil:
			t.Errorf("c3.C1 == nil")
		case c3.C1.Name != "c1":
			t.Errorf("c3.C1.Name = %q", c3.C1.Name)
		case c3.C1 != &c1:
			t.Errorf("c3.C1 = %v", c3.C1)
		}
	}()
}
