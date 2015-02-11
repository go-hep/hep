package rio

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestEmptyRWRecord(t *testing.T) {
	wrec := rioRecord{
		Header: rioHeader{
			Len:   1,
			Frame: recFrame,
		},
		Options: 2,
		CLen:    3,
		XLen:    4,
		Name:    "rio-record",
	}

	buf := new(bytes.Buffer)

	err := gob.NewEncoder(buf).Encode(&wrec)
	if err != nil {
		t.Fatalf("error encoding record: %v\n", err)
	}

	var rrec rioRecord
	err = gob.NewDecoder(buf).Decode(&rrec)
	if err != nil {
		t.Fatalf("error decoding record: %v\n", err)
	}

	size := rioAlign(buf.Len())
	if size != buf.Len() {
		t.Fatalf("buffer not 4-byte aligned. want=%d. got=%d\n",
			size, buf.Len(),
		)
	}

	if !reflect.DeepEqual(wrec, rrec) {
		t.Fatalf("error:\nwrec=%#v\nrrec=%#v\n", wrec, rrec)
	}
}

func TestReader(t *testing.T) {

	{
		rbuf := bytes.NewReader(magic[:])
		r := NewReader(rbuf)
		if r == nil {
			t.Fatalf("error creating new rio Reader")
		}
	}
	{
		rbuf := new(bytes.Buffer)
		r := NewReader(rbuf)
		if r != nil {
			t.Fatalf("NewReader should have failed")
		}
	}
}

func TestRW(t *testing.T) {
	const nmax = 10
	wbuf := new(bytes.Buffer)
	w := NewWriter(wbuf)
	if w == nil {
		t.Fatalf("error creating new rio Writer")
	}

	wrec := w.Record("data")
	err := wrec.Connect("event", &event{})
	if err != nil {
		t.Fatalf("error connecting block: %v", err)
	}
	wblk := wrec.Block("event")

	for i := 0; i < nmax; i++ {
		data := event{
			runnbr: int64(i),
			evtnbr: int64(1000 + i),
			id:     fmt.Sprintf("id-%04d", i),
			eles: []electron{
				electron{[4]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3)}},
				electron{[4]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3)}},
			},
			muons: []muon{
				muon{[4]float64{float64(-i), float64(-i - 1), float64(-i - 2), float64(-i - 3)}},
				muon{[4]float64{float64(-i), float64(-i - 1), float64(-i - 2), float64(-i - 3)}},
				muon{[4]float64{float64(-i), float64(-i - 1), float64(-i - 2), float64(-i - 3)}},
			},
		}
		if wblk.raw.Version != data.RioVersion() {
			t.Fatalf("error rio-version. want=%d. got=%d", data.RioVersion(), wblk.raw.Version)
		}

		err := wblk.Write(&data)
		if err != nil {
			t.Fatalf("error writing data[%d]: %v\n", i, err)
		}

		err = wrec.Write()
		if err != nil {
			t.Fatalf("error writing record[%d]: %v\n", i, err)
		}
	}
	err = w.Close()
	if err != nil {
		t.Fatalf("error closing writer: %v\n", err)
	}

	r := NewReader(wbuf)
	if r == nil {
		t.Fatalf("error creating new rio Reader")
	}

	rrec := r.Record("data")
	err = rrec.Connect("event", &event{})
	if err != nil {
		t.Fatalf("error connecting block: %v", err)
	}
	rblk := rrec.Block("event")

	for i := 0; i < nmax; i++ {
		err := rrec.Read()
		if err != nil {
			t.Fatalf("error loading record[%d]: %v\nbuf: %v\nraw: %#v\n", i, err,
				wbuf.Bytes(),
				rblk.raw,
			)
		}

		var data event
		err = rblk.Read(&data)
		if err != nil {
			t.Fatalf("error reading data[%d]: %v\n", i, err)
		}

		if rblk.raw.Version != data.RioVersion() {
			t.Fatalf("error rio-version. want=%d. got=%d", data.RioVersion(), rblk.raw.Version)
		}

		want := event{
			runnbr: int64(i),
			evtnbr: int64(1000 + i),
			id:     fmt.Sprintf("id-%04d", i),
			eles: []electron{
				electron{[4]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3)}},
				electron{[4]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3)}},
			},
			muons: []muon{
				muon{[4]float64{float64(-i), float64(-i - 1), float64(-i - 2), float64(-i - 3)}},
				muon{[4]float64{float64(-i), float64(-i - 1), float64(-i - 2), float64(-i - 3)}},
				muon{[4]float64{float64(-i), float64(-i - 1), float64(-i - 2), float64(-i - 3)}},
			},
		}

		if !reflect.DeepEqual(data, want) {
			t.Fatalf("error data[%d].\nwant=%#v\ngot =%#v\n", i, want, data)
		}
	}

	err = r.Close()
	if err != nil {
		t.Fatalf("error closing reading: %v\n", err)
	}
}

// event holds data to be serialized
type event struct {
	runnbr int64
	evtnbr int64
	id     string
	eles   []electron
	muons  []muon
}

func (evt *event) RioEncode(w io.Writer) error {
	err := binary.Write(w, Endian, evt.runnbr)
	if err != nil {
		return err
	}

	err = binary.Write(w, Endian, evt.evtnbr)
	if err != nil {
		return err
	}

	err = binary.Write(w, Endian, int64(len(evt.eles)))
	if err != nil {
		return err
	}
	for _, ele := range evt.eles {
		err = ele.RioEncode(w)
		if err != nil {
			return err
		}
	}

	err = binary.Write(w, Endian, int64(len(evt.muons)))
	if err != nil {
		return err
	}
	for _, muon := range evt.muons {
		err = muon.RioEncode(w)
		if err != nil {
			return err
		}
	}

	err = binary.Write(w, Endian, int64(len(evt.id)))
	if err != nil {
		return err
	}

	err = binary.Write(w, Endian, []byte(evt.id))
	if err != nil {
		return err
	}

	return err
}

func (evt *event) RioDecode(r io.Reader) error {
	err := binary.Read(r, Endian, &evt.runnbr)
	if err != nil {
		return err
	}

	err = binary.Read(r, Endian, &evt.evtnbr)
	if err != nil {
		return err
	}

	neles := int64(0)
	err = binary.Read(r, Endian, &neles)
	if err != nil {
		return err
	}

	evt.eles = make([]electron, int(neles))
	for i := range evt.eles {
		ele := &evt.eles[i]
		err = ele.RioDecode(r)
		if err != nil {
			return err
		}
	}

	nmuons := int64(0)
	err = binary.Read(r, Endian, &nmuons)
	if err != nil {
		return err
	}

	evt.muons = make([]muon, int(nmuons))
	for i := range evt.muons {
		muon := &evt.muons[i]
		err = muon.RioDecode(r)
		if err != nil {
			return err
		}
	}

	var nid int64
	err = binary.Read(r, Endian, &nid)
	if err != nil {
		return err
	}

	bid := make([]byte, int(nid))
	err = binary.Read(r, Endian, &bid)
	if err != nil {
		return err
	}
	evt.id = string(bid)

	return err
}

func (evt *event) RioVersion() Version {
	return Version(42)
}

type electron struct {
	p4 [4]float64
}

func (ele *electron) RioEncode(w io.Writer) error {
	return binary.Write(w, Endian, ele.p4)
}

func (ele *electron) RioDecode(r io.Reader) error {
	return binary.Read(r, Endian, &ele.p4)
}

type muon struct {
	p4 [4]float64
}

func (muon *muon) RioEncode(w io.Writer) error {
	return binary.Write(w, Endian, muon.p4)
}

func (muon *muon) RioDecode(r io.Reader) error {
	return binary.Read(r, Endian, &muon.p4)
}
