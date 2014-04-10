package dao_test

import (
	"bytes"
	"encoding/gob"
	"reflect"
	"sync"
	"testing"

	"github.com/go-hep/dao"
)

func TestH1D(t *testing.T) {
	h1 := dao.NewH1D(100, 0., 100.)
	if h1 == nil {
		t.Errorf("nil pointer to H1D")
	}

	h1.Annotation()["name"] = "h1"

	n := h1.Name()
	if n != "h1" {
		t.Errorf("expected H1D.Name() == %q (got %q)\n",
			n, "h1")
	}
	nbins := h1.Axis().Bins()
	if nbins != 100 {
		t.Errorf("expected H1D.Axis.Bins() == %v (got %v)\n",
			100, nbins,
		)
	}
	low := h1.Axis().LowerEdge()
	if low != 0. {
		t.Errorf("expected H1D.Axis.LowerEdge() == %v (got %v)\n",
			0., low,
		)
	}
	up := h1.Axis().UpperEdge()
	if up != 100. {
		t.Errorf("expected H1D.Axis.UpperEdge() == %v (got %v)\n",
			100., up,
		)
	}

	for idx := 0; idx < nbins; idx++ {
		size := h1.Axis().BinWidth(idx)
		if size != 1. {
			t.Errorf("expected H1D.Axis.BinWidth(%v) == %v (got %v)\n",
				idx, 1., size,
			)
		}
	}
}

func BenchmarkH1DSTFillConst(b *testing.B) {
	b.StopTimer()
	h1 := dao.NewH1D(100, 0., 100.)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		h1.Fill(10., 1.)
	}
}

func BenchmarkH1DFillFlat(b *testing.B) {
	b.StopTimer()
	h1 := dao.NewH1D(100, 0., 100.)
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		h1.Fill(rnd()*100., 1.)
	}
}

func BenchmarkH1DFillFlatGo(b *testing.B) {
	b.StopTimer()
	h1 := dao.NewH1D(100, 0., 100.)
	wg := new(sync.WaitGroup)
	//wg.Add(b.N)
	b.StartTimer()

	// throttle...
	q := make(chan struct{}, 1000)
	for i := 0; i < b.N; i++ {
		q <- struct{}{}
		go func() {
			wg.Add(1)
			h1.Fill(rnd()*100., 1.)
			<-q
			wg.Done()
		}()
	}
	wg.Wait()
}

func st_process_evts(n int, hists []*dao.H1D, process func(hists []*dao.H1D)) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			process(hists)
			wg.Done()
		}()
	}
	wg.Wait()
}

func st_process_evts_const(hists []*dao.H1D) {
	for _, h := range hists {
		h.Fill(10., 1.)
	}
}
func BenchmarkNH1DFillConst(b *testing.B) {
	b.StopTimer()
	hists := make([]*dao.H1D, 100)
	for i := 0; i < 100; i++ {
		hists[i] = dao.NewH1D(100, 0., 100.)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		st_process_evts(100, hists, st_process_evts_const)
	}
}

func st_process_evts_flat(hists []*dao.H1D) {
	for _, h := range hists {
		h.Fill(rnd()*100., 1.)
	}
}

func BenchmarkNH1DFillFlat(b *testing.B) {
	b.StopTimer()
	hists := make([]*dao.H1D, 100)
	for i := 0; i < 100; i++ {
		hists[i] = dao.NewH1D(100, 0., 100.)
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		st_process_evts(100, hists, st_process_evts_flat)
	}
}

func TestH1DSerialization(t *testing.T) {
	const nentries = 50
	href := dao.NewH1D(100, 0., 100.)
	for i := 0; i < nentries; i++ {
		href.Fill(rnd()*100., 1.)
	}
	href.Annotation()["title"] = "histo title"
	href.Annotation()["name"] = "histo name"

	// test gob.GobDecode/gob.GobEncode interface
	func() {
		buf := new(bytes.Buffer)
		enc := gob.NewEncoder(buf)
		err := enc.Encode(href)
		if err != nil {
			t.Fatalf("could not serialize histogram: %v", err)
		}

		var hnew dao.H1D
		buf = bytes.NewBuffer(buf.Bytes())
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&hnew)
		if err != nil {
			t.Fatalf("could not deserialize histogram: %v", err)
		}

		if !reflect.DeepEqual(href, &hnew) {
			t.Fatalf("ref=%v\nnew=%v\n", href, &hnew)
		}
	}()

	// test rio.Binary(Un)Marshaler
	func() {
		buf := new(bytes.Buffer)
		err := href.MarshalBinary(buf)
		if err != nil {
			t.Fatalf("could not serialize histogram: %v", err)
		}

		var hnew dao.H1D
		buf = bytes.NewBuffer(buf.Bytes())
		err = hnew.UnmarshalBinary(buf)
		if err != nil {
			t.Fatalf("could not deserialize histogram: %v", err)
		}

		if !reflect.DeepEqual(href, &hnew) {
			t.Fatalf("ref=%v\nnew=%v\n", href, &hnew)
		}
	}()

}

// EOF
