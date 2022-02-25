// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"sync"

	"go-hep.org/x/hep/groot/riofs"
)

type preader struct {
	r reader
	n int // number of workers
}

var (
	_ riofs.Reader = (*preader)(nil)
)

const blkSize = 1 * 1024 * 1024 // TODO(sbinet): adjust size for multiple payloads?

func (r *preader) Close() error {
	return r.r.Close()
}

func (r *preader) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *preader) ReadAt(p []byte, off int64) (int, error) {
	switch sz := len(p); {
	default:
		return r.r.ReadAt(p, off)
	case sz > blkSize:
		return r.pread(p, off)
	}
}

func (r *preader) pread(p []byte, off int64) (int, error) {
	nblks := len(p) / blkSize
	sps := make([]span, 0, nblks+1)
	beg := off
	end := off + int64(len(p))
	for beg < end {
		len := int64(blkSize)
		if beg+len > end {
			len = end - beg
		}
		sps = append(sps, span{
			off: beg,
			len: len,
		})
		beg += blkSize
	}
	out := make([]pread, len(sps))
	wrk := make(chan int, r.n)
	go func() {
		defer close(wrk)
		for i := range sps {
			wrk <- i
		}
	}()

	var wg sync.WaitGroup
	wg.Add(r.n)
	for i := 0; i < r.n; i++ {
		go func() {
			defer wg.Done()
			for i := range wrk {
				spn := sps[i]
				beg := int64(i) * blkSize
				end := beg + spn.len
				out[i].n, out[i].err = r.r.ReadAt(p[beg:end], spn.off)
			}
		}()
	}

	wg.Wait()

	var (
		n   int
		err error
	)
	for _, o := range out {
		n += o.n
		if o.err != nil {
			err = o.err
		}
	}
	return n, err
}

type pread struct {
	n   int
	err error
}
