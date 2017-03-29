# Copyright 2017 The go-hep Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from array import array as carray
import ROOT

f = ROOT.TFile.Open("gauss-h1.root","RECREATE")
histos = []
for t in [
        (ROOT.TH1D, "h1d", (10,-4,4)),
        (ROOT.TH1F, "h1f", (10,-4,4)),
        (ROOT.TH1F, "h1d-var", (10, carray("d", [
            -4.0, -3.2, -2.4, -1.6, -0.8,  0,
            +0.8, +1.6, +2.4, +3.2, +4.0
        ]))),
        (ROOT.TH1F, "h1f-var", (10, carray("f", [
            -4.0, -3.2, -2.4, -1.6, -0.8,  0,
            +0.8, +1.6, +2.4, +3.2, +4.0
        ])))
        ]:
    cls, name, args = t
    if len(args) == 3:
        h = cls(name, name, args[0], args[1], args[2])
    elif len(args) == 2:
        h = cls(name, name, args[0], args[1])
    else:
        raise ValueError("invalid number of arguments %d" % len(args))
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
