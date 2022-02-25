// Copyright ©2022 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"io"
	"os"
	"sync"

	"golang.org/x/sync/errgroup"
)

type reader interface {
	io.Reader
	io.ReaderAt
	io.Closer
}

type store interface {
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Closer

	Name() string
}

type rcache struct {
	r reader
	o store

	mu  sync.RWMutex
	sps spans
}

func rcacheOf(r reader) (*rcache, error) {
	f, err := os.CreateTemp("", "riofs-remote-")
	if err != nil {
		return nil, err
	}

	return &rcache{r: r, o: f}, nil
}

func (r *rcache) Close() error {
	e1 := r.r.Close()
	e2 := r.o.Close()
	_ = os.RemoveAll(r.o.Name())
	if e1 != nil {
		return e1
	}
	return e2
}

func (r *rcache) Read(p []byte) (int, error) {
	return r.r.Read(p)
}

func (r *rcache) ReadAt(p []byte, off int64) (int, error) {
	sp := span{off: off, len: int64(len(p))}
	oo := r.split(sp)
	if len(oo) == 0 {
		return r.o.ReadAt(p, off)
	}
	var (
		grp errgroup.Group
		ii  int64
	)
	for i := range oo {
		spa := oo[i]
		beg := ii
		end := ii + spa.len
		ii = end
		grp.Go(func() error {
			return r.fetch(p[beg:end], spa)
		})
	}

	err := grp.Wait()
	if err != nil {
		return 0, err
	}

	return r.o.ReadAt(p, off)
}

func (r *rcache) split(sp span) []span {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return split(sp, r.sps)
}

func (r *rcache) fetch(p []byte, sp span) error {
	_, err := r.r.ReadAt(p, sp.off)
	if err != nil {
		return err
	}
	_, err = r.o.WriteAt(p, sp.off)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sps.add(sp)
	return nil
}
