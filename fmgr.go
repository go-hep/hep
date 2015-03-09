package main

import (
	"fmt"
	"io"
	"os"

	"github.com/go-hep/rio"
)

type rfile struct {
	id  int
	n   string
	r   io.ReadCloser
	rio *rio.Reader
}

func (r *rfile) open(fname string) error {
	var err error

	r.n = fname
	r.r, err = os.Open(fname)
	if err != nil {
		return err
	}

	r.rio, err = rio.NewReader(r.r)
	if err != nil {
		return err
	}

	return err
}

func (r *rfile) ls() error {
	var err error

	f, err := os.Open(r.n)
	if err != nil {
		return err
	}
	defer f.Close()

	rr, err := rio.NewReader(f)
	if err != nil {
		return err
	}
	defer rr.Close()

	// FIXME(sbinet) instead of crawling through the whole file,
	//               use rio.StreamHdr when available.
	recs := make(map[string]struct{})
	scan := rio.NewScanner(rr)
	for scan.Scan() {
		err = scan.Err()
		if err != nil {
			break
		}
		n := scan.Record().Name()
		recs[n] = struct{}{}
	}
	fmt.Printf("/file/id/%d name=%s\n", r.id, r.n)
	for k := range recs {
		fmt.Printf("  %s\n", k)
	}

	err = scan.Err()
	if err != nil {
		return err
	}

	return err
}

func (r *rfile) read(name string, ptr interface{}) error {
	var err error

	f, err := os.Open(r.n)
	if err != nil {
		return err
	}
	defer f.Close()

	rr, err := rio.NewReader(f)
	if err != nil {
		return err
	}
	defer rr.Close()

	scan := rio.NewScanner(rr)
	scan.Select([]rio.Selector{
		{Name: name, Unpack: true},
	})

	for scan.Scan() {
		err = scan.Err()
		if err != nil {
			break
		}
		rec := scan.Record()
		n := rec.Name()
		if n != name {
			continue
		}

		blk := rec.Block(name)
		err = blk.Read(ptr)
		if err != nil {
			return err
		}
		return err
	}
	err = scan.Err()
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return err
	}

	return fmt.Errorf("no record [%s] in file [id=%d name=%s]", name, r.id, r.n)
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
