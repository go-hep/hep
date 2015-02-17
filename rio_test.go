// Copyright 2015 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rio

import (
	"bytes"
	"compress/flate"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"reflect"
	"testing"
)

func TestOptions(t *testing.T) {
	for _, kind := range []CompressorKind{CompressNone, CompressZlib, CompressGzip} {
		for _, level := range []int{flate.DefaultCompression, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9} {
			for _, codec := range []int{0, 1, 2} {
				o := NewOptions(kind, level, codec)
				if o.CompressorKind() != kind {
					t.Errorf("invalid CompressorKind. want=%v. got=%v",
						kind, o.CompressorKind(),
					)
				}
				if o.CompressorLevel() != level {
					t.Errorf("invalid CompressorLevel. want=%v. got=%v",
						level, o.CompressorLevel(),
					)
				}
				if o.CompressorCodec() != codec {
					t.Errorf("invalid CompressorCodec. want=%v. got=%v",
						codec, o.CompressorCodec(),
					)
				}
			}
		}
	}
}

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
		rbuf := bytes.NewReader(rioMagic[:])
		r, err := NewReader(rbuf)
		if err != nil || r == nil {
			t.Fatalf("error creating new rio Reader: %v", err)
		}
	}
	{
		rbuf := new(bytes.Buffer)
		r, err := NewReader(rbuf)
		if err == nil || r != nil {
			t.Fatalf("NewReader should have failed")
		}
	}
}

func TestRW(t *testing.T) {
	const nmax = 100
	makeles := func(i int) []electron {
		eles := make([]electron, 0, nmax)
		for j := 0; j < nmax; j++ {
			eles = append(
				eles,
				electron{[4]float64{float64(i), float64(i + 1), float64(i + 2), float64(i + 3)}},
			)
		}
		return eles
	}
	makemuons := func(i int) []muon {
		muons := make([]muon, 0, nmax)
		for j := 0; j < nmax; j++ {
			muons = append(
				muons,
				muon{[4]float64{float64(-i), float64(-i - 1), float64(-i - 2), float64(-i - 3)}},
			)
		}
		return muons
	}

	for ii, test := range []struct {
		lvl   int
		ckind CompressorKind
	}{
		{
			lvl:   0,
			ckind: CompressNone,
		},
		{
			lvl:   1,
			ckind: CompressNone,
		},

		// flate
		{
			lvl:   flate.NoCompression,
			ckind: CompressFlate,
		},
		{
			lvl:   flate.DefaultCompression,
			ckind: CompressFlate,
		},
		{
			lvl:   flate.BestCompression,
			ckind: CompressFlate,
		},
		{
			lvl:   flate.BestSpeed,
			ckind: CompressFlate,
		},

		// zlib
		{
			lvl:   flate.NoCompression,
			ckind: CompressZlib,
		},
		{
			lvl:   flate.DefaultCompression,
			ckind: CompressZlib,
		},
		{
			lvl:   flate.BestCompression,
			ckind: CompressZlib,
		},
		{
			lvl:   flate.BestSpeed,
			ckind: CompressZlib,
		},

		// gzip
		{
			lvl:   flate.NoCompression,
			ckind: CompressGzip,
		},
		{
			lvl:   flate.DefaultCompression,
			ckind: CompressGzip,
		},
		{
			lvl:   flate.BestCompression,
			ckind: CompressGzip,
		},
		{
			lvl:   flate.BestSpeed,
			ckind: CompressGzip,
		},
	} {
		wbuf := new(bytes.Buffer)
		w, err := NewWriter(wbuf)
		if w == nil || err != nil {
			t.Fatalf("test[%d]: error creating new rio Writer: %v", ii, err)
		}

		err = w.SetCompressor(test.ckind, test.lvl)
		if err != nil {
			t.Fatalf("test[%d]: error setting compressor (%#v): %v", ii, test, err)
		}

		wrec := w.Record("data")
		err = wrec.Connect("event", &event{})
		if err != nil {
			t.Fatalf("test[%d]: error connecting block: %v", ii, err)
		}
		wblk := wrec.Block("event")

		for i := 0; i < nmax; i++ {
			data := event{
				runnbr: int64(i),
				evtnbr: int64(1000 + i),
				id:     fmt.Sprintf("id-%04d", i),
				eles:   makeles(i),
				muons:  makemuons(i),
			}
			if wblk.raw.Version != data.RioVersion() {
				t.Fatalf("test[%d]: error rio-version. want=%d. got=%d", ii, data.RioVersion(), wblk.raw.Version)
			}

			err := wblk.Write(&data)
			if err != nil {
				t.Fatalf("test[%d]: error writing data[%d]: %v\n", ii, i, err)
			}

			err = wrec.Write()
			if err != nil {
				t.Fatalf("test[%d]: error writing record[%d]: %v\n", ii, i, err)
			}
		}
		err = w.Close()
		if err != nil {
			t.Fatalf("test[%d]: error closing writer: %v\n", ii, err)
		}

		// fmt.Printf("::: kind: %7q lvl: %2d: size=%8d\n", test.ckind, test.lvl, wbuf.Len())

		r, err := NewReader(wbuf)
		if err != nil {
			t.Fatalf("test[%d]: error creating new rio Reader: %v", ii, err)
		}

		rrec := r.Record("data")
		err = rrec.Connect("event", &event{})
		if err != nil {
			t.Fatalf("test[%d]: error connecting block: %v", ii, err)
		}
		rblk := rrec.Block("event")

		for i := 0; i < nmax; i++ {
			err := rrec.Read()
			if err != nil {
				t.Fatalf("test[%d]: error loading record[%d]: %v\nbuf: %v\nraw: %#v\n", ii, i, err,
					wbuf.Bytes(),
					rblk.raw,
				)
			}

			var data event
			err = rblk.Read(&data)
			if err != nil {
				t.Fatalf("test[%d]: error reading data[%d]: %v\n", ii, i, err)
			}

			if rblk.raw.Version != data.RioVersion() {
				t.Fatalf("test[%d]: error rio-version. want=%d. got=%d", ii, data.RioVersion(), rblk.raw.Version)
			}

			want := event{
				runnbr: int64(i),
				evtnbr: int64(1000 + i),
				id:     fmt.Sprintf("id-%04d", i),
				eles:   makeles(i),
				muons:  makemuons(i),
			}

			if !reflect.DeepEqual(data, want) {
				t.Fatalf("test[%d]: error data[%d].\nwant=%#v\ngot =%#v\n", ii, i, want, data)
			}
		}

		err = r.Close()
		if err != nil {
			t.Fatalf("test[%d]: error closing reading: %v\n", ii, err)
		}
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

func (evt *event) RioMarshal(w io.Writer) error {
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
		err = ele.RioMarshal(w)
		if err != nil {
			return err
		}
	}

	err = binary.Write(w, Endian, int64(len(evt.muons)))
	if err != nil {
		return err
	}
	for _, muon := range evt.muons {
		err = muon.RioMarshal(w)
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

func (evt *event) RioUnmarshal(r io.Reader) error {
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
		err = ele.RioUnmarshal(r)
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
		err = muon.RioUnmarshal(r)
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

func (ele *electron) RioMarshal(w io.Writer) error {
	return binary.Write(w, Endian, ele.p4)
}

func (ele *electron) RioUnmarshal(r io.Reader) error {
	return binary.Read(r, Endian, &ele.p4)
}

type muon struct {
	p4 [4]float64
}

func (muon *muon) RioMarshal(w io.Writer) error {
	return binary.Write(w, Endian, muon.p4)
}

func (muon *muon) RioUnmarshal(r io.Reader) error {
	return binary.Read(r, Endian, &muon.p4)
}
