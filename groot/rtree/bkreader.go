// Copyright ©2020 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rtree

import (
	"fmt"
	"io"
	"runtime"

	"go-hep.org/x/hep/groot/riofs"
)

// bkreader is a read-ahead basket reader.
type bkreader struct {
	f     *riofs.File
	spans []rspan

	beg    int64         // first event to process
	end    int64         // last-1 event to process (ie: [beg,end) half-open interval of entries to process)
	ready  chan bkReq    // baskets ready to be handed to the reader
	reuse  chan bkReq    // baskets to reuse for input reading
	exit   chan struct{} // closes when finished
	n      int           // number of in-flight baskets
	cur    *rbasket      // current buffer being served
	closed chan struct{} // channel is closed when the async reader shuts down

	name string
}

type bkReq struct {
	bkt *rbasket
	err error
}

func newBkReader(b Branch, n int, beg, end int64) *bkreader {
	if n < 0 {
		n = runtime.NumCPU() + 1
	}
	if n == 0 {
		n = 1
	}
	base := asBranch(b)
	if m := len(base.basketSeek); n > m && m != 0 {
		n = m
	}
	bkr := &bkreader{
		f:      b.getTree().f,
		spans:  make([]rspan, len(base.basketSeek)),
		beg:    beg,
		end:    end,
		ready:  make(chan bkReq, n),
		reuse:  make(chan bkReq, n),
		exit:   make(chan struct{}),
		n:      n,
		closed: make(chan struct{}),
		name:   b.Name(),
	}

	for i, seek := range base.basketSeek {
		bkr.spans[i] = rspan{
			pos: seek,
			sz:  base.basketBytes[i],
			beg: base.basketEntry[i],
			end: base.basketEntry[i+1],
		}
	}

	for i := 0; i < n; i++ {
		bkr.reuse <- bkReq{bkt: new(rbasket), err: nil}
	}

	switch {
	case base.entries == base.basketEntry[len(base.basketSeek)]:
		// ok, normal case.
	default: // recover baskets
		var beg int64
		if len(bkr.spans) > 0 {
			beg = bkr.spans[len(bkr.spans)-1].end
		}
		for i := range base.baskets {
			bkt := &base.baskets[i]
			span := rspan{
				pos: 0,
				beg: beg,
				end: beg + int64(bkt.nevbuf),
				bkt: bkt,
			}
			bkr.spans = append(bkr.spans, span)
			beg = span.end
		}
	}

	all := rspan{
		beg: bkr.spans[0].beg,
		end: bkr.spans[len(bkr.spans)-1].end,
	}

	if beg >= all.end {
		// empty range: no need to run prefetch loop.
		defer close(bkr.closed)
		defer close(bkr.ready)
		return bkr
	}

	ibeg, iend := bkr.findBaskets(beg, end)
	if ibeg < 0 || iend < 0 {
		panic(fmt.Errorf(
			"rtree: could not find basket index for span [%d, %d): [%d, %d), spans: %#v",
			beg, end, ibeg, iend, bkr.spans,
		))
	}

	go bkr.run(base.entryOffsetLen, ibeg, iend)

	return bkr
}

func (bkr *bkreader) findBaskets(beg, end int64) (int, int) {
	var (
		ibeg = -1
		iend = -1
	)
	for i, v := range bkr.spans {
		if ibeg < 0 && (v.beg <= beg && beg < v.end) {
			ibeg = i
		}
		if iend < 0 && (end <= v.end) {
			iend = i + 1
		}
	}
	if iend < 0 {
		iend = len(bkr.spans)
	}

	return ibeg, iend
}

func (bkr *bkreader) run(eoff, beg, end int) {
	defer close(bkr.closed)
	defer close(bkr.ready)
	for i, span := range bkr.spans[beg:end] {
		select {
		case tok := <-bkr.reuse:
			tok.err = tok.bkt.inflate(bkr.name, beg+i, span, eoff, bkr.f)
			bkr.ready <- tok
		case <-bkr.exit:
			return
		}
	}
}

func (bkr *bkreader) read() (*rbasket, error) {
	if bkr.cur != nil {
		bkr.cur.reset()
		bkr.reuse <- bkReq{bkt: bkr.cur, err: nil}
		bkr.cur = nil
	}
	tok, ok := <-bkr.ready
	if !ok {
		return nil, io.EOF
	}
	bkr.cur = tok.bkt

	return bkr.cur, tok.err
}

func (bkr *bkreader) close() {
	select {
	case <-bkr.closed:
	case bkr.exit <- struct{}{}:
		<-bkr.closed
	}
}

type rspan struct {
	beg int64 // first entry
	end int64 // last entry

	pos int64 // span location on-disk
	sz  int32 // basket size

	// for recovered baskets
	bkt *Basket
}
