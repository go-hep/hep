package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-hep/hbook"
	"github.com/gonum/plot"
	"github.com/gonum/plot/vg/draw"
	"github.com/gonum/plot/vg/vgx11"
)

type histMgr struct {
	h1ds map[int]*hbook.H1D
}

func newHistMgr() histMgr {
	return histMgr{
		h1ds: make(map[int]*hbook.H1D),
	}
}

func (mgr *histMgr) openH1D(fmgr *fileMgr, hid int, path string) error {
	var err error
	const prefix = "/file/id/"
	if !strings.HasPrefix(path, prefix) {
		return fmt.Errorf("invalid path [%s] (missing prefix [%s])", path, prefix)
	}

	var toks []string
	for _, tok := range strings.Split(path[len(prefix):], "/") {
		tok = strings.TrimSpace(tok)
		if tok == "" {
			continue
		}
		toks = append(toks, tok)
	}

	if len(toks) < 2 {
		return fmt.Errorf("invalid path [%s] (missing file-id and histo-name)", path)
	}

	fid, err := strconv.Atoi(toks[0])
	if err != nil {
		return err
	}

	r, ok := fmgr.rfds[fid]
	if !ok {
		return fmt.Errorf("unknown file-id [%d]", fid)
	}

	hname := toks[1]

	var h1d hbook.H1D
	err = r.read(hname, &h1d)
	if err != nil {
		return err
	}

	mgr.h1ds[hid] = &h1d
	return err
}

func (mgr *histMgr) plotH1D(hid int) error {
	var err error
	h, ok := mgr.h1ds[hid]
	if !ok {
		return fmt.Errorf("unknown H1D [id=%d]", hid)
	}

	fmt.Printf("== h1d: name=%q\nentries=%d\nmean=%+8.3f\nRMS= %+8.3f\n",
		h.Name(), h.Entries(), h.Mean(), h.RMS(),
	)

	p, err := plot.New()
	if err != nil {
		return err
	}
	p.Title.Text = h.Name()
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	hh, err := NewH1D(h)
	if err != nil {
		return err
	}
	p.Add(hh)

	cnv, err := vgx11.New(4*96, 4*96, "paw")
	if err != nil {
		return err
	}

	p.Draw(draw.New(cnv))
	cnv.Paint()

	return err
}
