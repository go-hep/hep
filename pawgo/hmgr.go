// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"strings"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
)

type histMgr struct {
	msg  *log.Logger
	hmap map[string]hbook.Histogram
}

func newHistMgr(msg *log.Logger) *histMgr {
	return &histMgr{
		msg:  msg,
		hmap: make(map[string]hbook.Histogram),
	}
}

func (mgr *histMgr) find(fmgr *fileMgr, path string) (hbook.Histogram, error) {
	var err error
	const prefix = "/file/id/"
	if !strings.HasPrefix(path, prefix) {
		return nil, fmt.Errorf("invalid path [%s] (missing prefix [%s])", path, prefix)
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
		return nil, fmt.Errorf("invalid path [%s] (missing file-id and histo-name)", path)
	}

	fid := toks[0]

	r, ok := fmgr.rfds[fid]
	if !ok {
		return nil, fmt.Errorf("unknown file-id [%s]", fid)
	}

	hname := strings.Join(toks[1:], "/")

	switch r.typ(hname) {
	case "*go-hep.org/x/hep/hbook.H1D":
		var h1 hbook.H1D
		err = r.read(hname, &h1)
		if err != nil {
			return nil, err
		}
		return &h1, nil

	case "*go-hep.org/x/hep/hbook.H2D":
		var h2 hbook.H2D
		err = r.read(hname, &h2)
		if err != nil {
			return nil, err
		}
		return &h2, nil

	default:
		return nil, fmt.Errorf("%q not an histogram (%s)", path, r.typ(hname))
	}
}

func (mgr *histMgr) open(fmgr *fileMgr, hid, path string) error {
	h, err := mgr.find(fmgr, path)
	if err != nil {
		return err
	}
	mgr.hmap[hid] = h
	return nil
}

func (mgr *histMgr) plot(fmgr *fileMgr, wmgr *winMgr, hid string) (*window, error) {
	var (
		h   hbook.Histogram
		err error
	)
	if strings.HasPrefix(hid, "/file/id/") {
		// directly plot from file
		h, err = mgr.find(fmgr, hid)
		if err != nil {
			return nil, err
		}
	} else {
		var ok bool
		h, ok = mgr.hmap[hid]
		if !ok {
			return nil, fmt.Errorf("unknown histogram [id=%s]", hid)
		}
	}

	switch h := h.(type) {
	case *hbook.H1D:
		return mgr.plotH1D(wmgr, h)
	case *hbook.H2D:
		return mgr.plotH2D(wmgr, h)
	}

	return nil, fmt.Errorf("unknown histogram type %T [id=%s]", h, hid)
}

func (mgr *histMgr) plotH1D(wmgr *winMgr, h *hbook.H1D) (*window, error) {
	fmt.Fprintf(
		mgr.msg.Writer(),
		"== h1d: name=%q\nentries=%d\nmean=%+8.3f\nRMS= %+8.3f\n",
		h.Name(), h.Entries(), h.XMean(), h.XRMS(),
	)

	p := hplot.New()
	p.Title.Text = h.Name()
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	hh := hplot.NewH1D(h)
	hh.Infos.Style = hplot.HInfoSummary

	p.Add(hh)
	p.Add(hplot.NewGrid())

	win := wmgr.newPlot(p)
	return win, nil
}

func (mgr *histMgr) plotH2D(wmgr *winMgr, h *hbook.H2D) (*window, error) {
	fmt.Fprintf(
		mgr.msg.Writer(),
		"== h2d: name=%q\nentries=%d\nxmean=%+8.3f\nxRMS= %+8.3f\nymean=%+8.3f\nyRMS= %+8.3f\n",
		h.Name(), h.Entries(),
		h.XMean(), h.XRMS(),
		h.YMean(), h.YRMS(),
	)

	p := hplot.New()
	p.Title.Text = h.Name()
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	hh := hplot.NewH2D(h, nil)
	hh.Infos.Style = hplot.HInfoNone

	p.Add(hh)
	p.Add(hplot.NewGrid())

	win := wmgr.newPlot(p)
	return win, nil
}
