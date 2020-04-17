# Copyright ©2017 The go-hep Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

from array import array as carray
import ROOT

f = ROOT.TFile.Open("gauss-h2.root", "RECREATE")
for t in [
        (ROOT.TH2F, "h2f", (3,0,3, 3,0,3)),
        (ROOT.TH2D, "h2d", (3,0,3, 3,0,3)),
        (ROOT.TH2F, "h2f-var", ( 
            (3, carray("f", [0.0, 1.5, 2.0, 3.0])),
            (3, carray("f", [0.0, 1.5, 2.0, 3.0])))),
        (ROOT.TH2D, "h2d-var", (
            (3, carray("d", [0.0, 1.5, 2.0, 3.0])),
            (3, carray("d", [0.0, 1.5, 2.0, 3.0]))))
        ]:
    cls, name, args = t
    if len(args) == 6:
        h = cls(name, name, args[0], args[1], args[2], args[3], args[4], args[5])
    elif len(args) == 2:
        h = cls(name, name, args[0][0], args[0][1], args[1][0], args[1][1])
    else:
        raise ValueError("invalid number of arguments %d" % len(args))
    h.StatOverflows(True)
    h.Sumw2()
    with open("gauss-2d-data.dat") as ff:
        for l in ff.readlines():
            x, y, w = l.split()
            h.Fill(float(x),float(y),float(w))
            pass
        pass
    h.Fill(+5,+5,101) # NE 
    h.Fill(+0,+5,102) # N
    h.Fill(-5,+5,103) # NW
    h.Fill(-5,+0,104) # W
    h.Fill(-5,-5,105) # SW
    h.Fill(+0,-5,106) # S
    h.Fill(+5,-5,107) # SE
    h.Fill(+5,+0,108) # E

    histos.append(h)
    pass
f.Write()
f.Close()
