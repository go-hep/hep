# Copyright 2017 The go-hep Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

import ROOT

f = ROOT.TFile.Open("gauss-h1.root","RECREATE")
histos = []
for t in [(ROOT.TH1D, "h1d"), (ROOT.TH1F, "h1f")]:
    cls = t[0]
    name = t[1]
    h = cls(name, name, 10, -4, 4)
    h.StatOverflows(True)
    h.Sumw2()
    with open("gauss-1d-data.dat") as ff:
        for l in ff.readlines():
            x, w = l.split()
            h.Fill(float(x),float(w))
            pass
        pass
    histos.append(h)
    pass
f.Write()
f.Close()
