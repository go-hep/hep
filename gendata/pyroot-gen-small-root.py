#!/usr/bin/env python2
import ROOT
from array import array as carray

ARRAYSZ = 10
EVTMAX = 100

def main():
    fname = "test-small.root"

    compress = True
    netopt = 0
    splitlevel = 32
    bufsiz = 32000
    
    f = ROOT.TFile.Open(fname, "recreate", "small event file", compress, netopt)
    if not f:
        raise SystemError()

    t = ROOT.TTree("tree", "tree", splitlevel)

    i32 = carray("i", [0])
    i64 = carray("l", [0])
    u32 = carray("I", [0])
    u64 = carray("L", [0])
    f32 = carray("f", [0.])
    f64 = carray("d", [0.])

    arr_i32 = carray("i", [0]*ARRAYSZ)
    arr_i64 = carray("l", [0]*ARRAYSZ)
    arr_u32 = carray("I", [0]*ARRAYSZ)
    arr_u64 = carray("L", [0]*ARRAYSZ)
    arr_f32 = carray("f", [0]*ARRAYSZ)
    arr_f64 = carray("d", [0]*ARRAYSZ)

    t.Branch("Int32", i32, "Int32/I")
    t.Branch("Int64", i64, "Int64/L")
    t.Branch("UInt32", u32, "UInt32/i")
    t.Branch("UInt64", u64, "UInt64/l")
    t.Branch("Float32", f32, "Float32/F")
    t.Branch("Float64", f64, "Float64/D")

    t.Branch("ArrayInt32", arr_i32, "Int32[10]/I")
    t.Branch("ArrayInt64", arr_i64, "Int64[10]/L")
    t.Branch("ArrayUInt32", arr_u32, "Int32[10]/i")
    t.Branch("ArrayUInt64", arr_u64, "Int64[10]/l")
    t.Branch("ArrayFloat32", arr_f32, "Float32[10]/F")
    t.Branch("ArrayFloat64", arr_f64, "Float64[10]/D")

    for i in xrange(EVTMAX):
        #print ">>>",i
        i32[0] = i
        i64[0] = i
        u32[0] = i
        u64[0] = i

        f32[0] = i
        f64[0] = i
        
        for jj in range(ARRAYSZ):
            arr_i32[jj] = i
            arr_i64[jj] = i
            arr_u32[jj] = i
            arr_u64[jj] = i
            arr_f32[jj] = i
            arr_f64[jj] = i
        t.Fill()
        pass

    f.Write()
    f.Close()

if __name__ == "__main__":
    main()
    
    
