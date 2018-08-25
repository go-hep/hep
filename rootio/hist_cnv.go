// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rootio

import (
	"go-hep.org/x/hep/hbook"
)

func newH1D() *H1D {
	return &H1D{
		rvers: 2, // FIXME(sbinet): harmonize versions
		th1:   *newH1(),
	}
}

func newH1() *th1 {
	return &th1{
		rvers:     7, // FIXME(sbinet): harmonize versions
		tnamed:    *newNamed("", ""),
		attline:   *newAttLine(),
		attfill:   *newAttFill(),
		attmarker: *newAttMarker(),
		xaxis:     *newAxis("xaxis"),
		yaxis:     *newAxis("yaxis"),
		zaxis:     *newAxis("zaxis"),
		funcs:     *newList(""),
	}
}

func NewH1DFrom(h *hbook.H1D) *H1D {
	var (
		hroot = newH1D()
		nbins = h.Len()
		edges = make([]float64, 0, nbins+1)
		sumw  = make([]float64, 0, nbins+2)
		sumw2 = make([]float64, 0, nbins+2)
		uflow = h.Binning().Underflow()
		oflow = h.Binning().Overflow()
	)

	sumw = append(sumw, uflow.SumW())
	sumw2 = append(sumw2, uflow.SumW2())

	for i, bin := range h.Binning().Bins() {
		if i == 0 {
			edges = append(edges, bin.XMin())
		}
		edges = append(edges, bin.XMax())
		sumw = append(sumw, bin.SumW())
		sumw2 = append(sumw2, bin.SumW2())
	}
	sumw = append(sumw, oflow.SumW())
	sumw2 = append(sumw2, oflow.SumW2())

	hroot.th1.name = h.Name()
	if v, ok := h.Annotation()["title"]; ok {
		hroot.th1.title = v.(string)
	}
	hroot.th1.entries = float64(h.Entries())
	hroot.th1.tsumw = h.SumW()
	hroot.th1.tsumw2 = h.SumW2()
	hroot.th1.tsumwx = h.SumWX()
	hroot.th1.tsumwx2 = h.SumWX2()
	hroot.th1.ncells = len(edges)
	hroot.th1.xaxis.xbins.Data = edges
	hroot.th1.xaxis.nbins = nbins
	hroot.th1.xaxis.xmin = h.XMin()
	hroot.th1.xaxis.xmax = h.XMax()
	hroot.arr.Data = sumw
	hroot.sumw2.Data = sumw2

	return hroot
}

func (h1d *H1D) UnmarshalYODA(raw []byte) error {
	var h hbook.H1D
	err := h.UnmarshalYODA(raw)
	if err != nil {
		return err
	}

	*h1d = *NewH1DFrom(&h)
	return nil
}

/*
// readYODAHeader parses the input buffer and extracts the YODA header line
// from that buffer.
// readYODAHeader returns the associated YODA path and an error if any.
func readYODAHeader(r *bytes.Buffer, hdr string) (string, error) {
	pos := bytes.Index(r.Bytes(), []byte("\n"))
	if pos < 0 {
		return "", fmt.Errorf("rootio: could not find %s line", hdr)
	}
	path := string(r.Next(pos + 1))
	if !strings.HasPrefix(path, hdr+" ") {
		return "", fmt.Errorf("rootio: could not find %s mark", hdr)
	}

	return path[len(hdr)+1 : len(path)-1], nil
}
func splitYODAHeader(raw []byte) (Object, error) {
	raw = raw[len(begYoda):]
	i := bytes.Index(raw, []byte(" "))
	if i == -1 || i >= len(raw) {
		return nil, fmt.Errorf("invalid YODA header (missing space)")
	}

	var o Object

	switch string(raw[:i]) {
	case "HISTO1D":
		o = Factory.Get("TH1D")
	case "HISTO2D":
		o = Factory.Get("TH2D")
	case "PROFILE1D":
		// o = Factory.Get("TProfile")
		return nil, errYODAIgnore
	case "PROFILE2D":
		return nil, errYODAIgnore
	case "SCATTER1D":
		return nil, errYODAIgnore
	case "SCATTER2D":
		o = Factory.Get("TGraphAsymmErrors")
	case "SCATTER3D":
		return nil, errYODAIgnore
	case "COUNTER":
		return nil, errYODAIgnore
	default:
		return nil, fmt.Errorf("unhandled YODA object type %q", string(raw[:i]))
	}

	return o, nil
}
*/
