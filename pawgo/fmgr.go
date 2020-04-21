// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"text/tabwriter"

	"go-hep.org/x/hep/rio"
)

type fileType interface {
	io.Reader
	io.Seeker
	io.Closer
}

type rfile struct {
	id  string
	n   string
	r   fileType
	rio *rio.File
}

func (r *rfile) open(fname string) error {
	var err error

	r.n = fname
	r.r, err = os.Open(fname)
	if err != nil {
		return err
	}

	r.rio, err = rio.Open(r.r)
	if err != nil {
		return err
	}

	return err
}

func (r *rfile) ls(o io.Writer) error {
	var err error

	fmt.Fprintf(o, "/file/id/%s name=%s\n", r.id, r.n)
	w := tabwriter.NewWriter(o, 0, 8, 0, '\t', 0)
	for _, k := range r.rio.Keys() {
		fmt.Fprintf(w, " \t- %s\t(type=%q)\n", k.Name, k.Blocks[0].Type)
	}
	w.Flush()
	fmt.Fprintf(o, "\n")

	return err
}

func (r *rfile) typ(name string) string {
	for _, k := range r.rio.Keys() {
		if k.Name == name {
			return k.Blocks[0].Type
		}
	}

	return ""
}

func (r *rfile) read(name string, ptr interface{}) error {
	var err error

	// FIXME(sbinet): when/if "rio" gets the concept of directories,
	// handle this there.
	if !r.rio.Has(name) {
		return fmt.Errorf("no record [%s] in file [id=%s name=%s]", name, r.id, r.n)
	}

	err = r.rio.Get(name, ptr)
	if err != nil {
		return err
	}

	return err
}

func (r *rfile) close() error {
	defer r.r.Close()
	err := r.rio.Close()
	if err != nil {
		return err
	}
	return r.r.Close()
}

type wfile struct {
	id  string
	n   string
	w   io.WriteCloser
	rio *rio.Writer
}

func (w *wfile) create(fname string) error {
	var err error

	w.n = fname
	w.w, err = os.Create(fname)
	if err != nil {
		return err
	}

	w.rio, err = rio.NewWriter(w.w)
	if err != nil {
		return err
	}

	return err
}

func (w *wfile) close() error {
	defer w.w.Close()
	err := w.rio.Close()
	if err != nil {
		return err
	}
	return w.w.Close()
}

type fileMgr struct {
	msg  *log.Logger
	rfds map[string]rfile
	wfds map[string]wfile
}

func newFileMgr(msg *log.Logger) *fileMgr {
	return &fileMgr{
		msg:  msg,
		rfds: make(map[string]rfile),
		wfds: make(map[string]wfile),
	}
}

func (mgr *fileMgr) open(id string, fname string) error {
	var err error
	r, dup := mgr.rfds[id]
	if dup {
		return fmt.Errorf("paw: file [id=%s name=%s] already open", id, r.n)
	}

	r.id = id
	err = r.open(fname)
	if err != nil {
		return err
	}

	mgr.rfds[id] = r
	return nil
}

func (mgr *fileMgr) close(id string) error {
	r, ok := mgr.rfds[id]
	if ok {
		delete(mgr.rfds, id)
		return r.close()
	}

	w, ok := mgr.wfds[id]
	if ok {
		delete(mgr.wfds, id)
		return w.close()
	}

	return fmt.Errorf("paw: unknown file [id=%s]", id)
}

func (mgr *fileMgr) ls(id string) error {
	if id == "" {
		// list all
		for id := range mgr.rfds {
			err := mgr.ls(id)
			if err != nil {
				return err
			}
		}
		return nil
	}

	r, ok := mgr.rfds[id]
	if !ok {
		return fmt.Errorf("paw: unknown file [id=%s]", id)
	}

	err := r.ls(mgr.msg.Writer())
	if err != nil {
		return err
	}

	return err
}

func (mgr *fileMgr) create(id string, fname string) error {
	var err error
	w, dup := mgr.wfds[id]
	if dup {
		return fmt.Errorf("paw: file [id=%s name=%s] already open", id, w.n)
	}

	w.id = id
	err = w.create(fname)
	if err != nil {
		return err
	}

	mgr.wfds[id] = w
	return nil
}

func (mgr *fileMgr) Close() error {
	var err error
	for k, r := range mgr.rfds {
		e := r.close()
		if e != nil {
			mgr.msg.Printf("error closing file [%s]: %v\n", k, e)
			if err != nil {
				err = e
			}
		}
	}

	for k, w := range mgr.wfds {
		e := w.close()
		if e != nil {
			mgr.msg.Printf("error closing file [%s]: %v\n", k, e)
			if err != nil {
				err = e
			}
		}
	}

	return err
}
