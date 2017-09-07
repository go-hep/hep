// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"strings"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
)

type histMgr struct {
	h1ds map[string]*hbook.H1D
}

func newHistMgr() *histMgr {
	return &histMgr{
		h1ds: make(map[string]*hbook.H1D),
	}
}

func (mgr *histMgr) openH1D(fmgr *fileMgr, hid, path string) error {
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

	fid := toks[0]

	r, ok := fmgr.rfds[fid]
	if !ok {
		return fmt.Errorf("unknown file-id [%s]", fid)
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

func (mgr *histMgr) plotH1D(wmgr *winMgr, hid string) error {
	var err error
	h, ok := mgr.h1ds[hid]
	if !ok {
		return fmt.Errorf("unknown H1D [id=%s]", hid)
	}

	fmt.Printf("== h1d: name=%q\nentries=%d\nmean=%+8.3f\nRMS= %+8.3f\n",
		h.Name(), h.Entries(), h.XMean(), h.XRMS(),
	)

	p, err := hplot.New()
	if err != nil {
		return err
	}
	p.Title.Text = h.Name()
	p.X.Label.Text = "x"
	p.Y.Label.Text = "y"

	hh, err := hplot.NewH1D(h)
	if err != nil {
		return err
	}
	hh.Infos.Style = hplot.HInfoSummary

	p.Add(hh)

	err = wmgr.newPlot(p)
	if err != nil {
		return err
	}

	return err
}
