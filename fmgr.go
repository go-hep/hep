package main

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"

	"github.com/go-hep/rio"
)

type fileType interface {
	io.Reader
	io.Seeker
	io.Closer
}

type rfile struct {
	id  int
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

func (r *rfile) ls() error {
	var err error

	fmt.Printf("/file/id/%d name=%s\n", r.id, r.n)
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 0, '\t', 0)
	for _, k := range r.rio.Keys() {
		fmt.Fprintf(w, " \t- %s\t(type=%q)\n", k.Name, k.Blocks[0].Type)
	}
	w.Flush()
	fmt.Printf("\n")

	return err
}

func (r *rfile) read(name string, ptr interface{}) error {
	var err error

	if !r.rio.Has(name) {
		return fmt.Errorf("no record [%s] in file [id=%d name=%s]", name, r.id, r.n)
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
	id  int
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
	rfds map[int]rfile
	wfds map[int]wfile
}

func newFileMgr() fileMgr {
	return fileMgr{
		rfds: make(map[int]rfile),
		wfds: make(map[int]wfile),
	}
}

func (mgr *fileMgr) open(id int, fname string) error {
	var err error
	r, dup := mgr.rfds[id]
	if dup {
		return fmt.Errorf("paw: file [id=%d name=%s] already open", id, r.n)
	}

	r.id = id
	err = r.open(fname)
	if err != nil {
		return err
	}

	mgr.rfds[id] = r
	return nil
}

func (mgr *fileMgr) close(id int) error {
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

	return fmt.Errorf("paw: unknown file [id=%d]", id)
}

func (mgr *fileMgr) ls(id int) error {
	r, ok := mgr.rfds[id]
	if !ok {
		return fmt.Errorf("paw: unknown file [id=%d]", id)
	}

	err := r.ls()
	if err != nil {
		return err
	}

	return err
}

func (mgr *fileMgr) create(id int, fname string) error {
	var err error
	w, dup := mgr.wfds[id]
	if dup {
		return fmt.Errorf("paw: file [id=%d name=%s] already open", id, w.n)
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
			fmt.Printf("error closing file [%s]: %v\n", k, e)
			if err != nil {
				err = e
			}
		}
	}

	for k, w := range mgr.wfds {
		e := w.close()
		if e != nil {
			fmt.Printf("error closing file [%s]: %v\n", k, e)
			if err != nil {
				err = e
			}
		}
	}

	return err
}
