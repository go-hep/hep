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

    t = ROOT.TTree("tree", "my tree title", splitlevel)

    i32 = carray("i", [0])
    i64 = carray("l", [0])
    u32 = carray("I", [0])
    u64 = carray("L", [0])
    f32 = carray("f", [0.])
    f64 = carray("d", [0.])
    s06 = carray("b", [0]*7)

    arr_i32 = carray("i", [0]*ARRAYSZ)
    arr_i64 = carray("l", [0]*ARRAYSZ)
    arr_u32 = carray("I", [0]*ARRAYSZ)
    arr_u64 = carray("L", [0]*ARRAYSZ)
    arr_f32 = carray("f", [0]*ARRAYSZ)
    arr_f64 = carray("d", [0]*ARRAYSZ)

    n = carray("i", [0])
    sli_i32 = carray("i", [0]*ARRAYSZ)
    sli_i64 = carray("l", [0]*ARRAYSZ)
    sli_u32 = carray("I", [0]*ARRAYSZ)
    sli_u64 = carray("L", [0]*ARRAYSZ)
    sli_f32 = carray("f", [0]*ARRAYSZ)
    sli_f64 = carray("d", [0]*ARRAYSZ)

    t.Branch("Int32", i32, "Int32/I")
    t.Branch("Int64", i64, "Int64/L")
    t.Branch("UInt32", u32, "UInt32/i")
    t.Branch("UInt64", u64, "UInt64/l")
    t.Branch("Float32", f32, "Float32/F")
    t.Branch("Float64", f64, "Float64/D")
    t.Branch("Str", s06, "Str/C")

    t.Branch("ArrayInt32", arr_i32, "ArrayInt32[10]/I")
    t.Branch("ArrayInt64", arr_i64, "ArrayInt64[10]/L")
    t.Branch("ArrayUInt32", arr_u32, "ArrayInt32[10]/i")
    t.Branch("ArrayUInt64", arr_u64, "ArrayInt64[10]/l")
    t.Branch("ArrayFloat32", arr_f32, "ArrayFloat32[10]/F")
    t.Branch("ArrayFloat64", arr_f64, "ArrayFloat64[10]/D")

    t.Branch("N", n, "N/I")
    t.Branch("SliceInt32", sli_i32, "SliceInt32[N]/I")
    t.Branch("SliceInt64", sli_i64, "SliceInt64[N]/L")
    t.Branch("SliceUInt32", sli_u32, "SliceInt32[N]/i")
    t.Branch("SliceUInt64", sli_u64, "SliceInt64[N]/l")
    t.Branch("SliceFloat32", sli_f32, "SliceFloat32[N]/F")
    t.Branch("SliceFloat64", sli_f64, "SliceFloat64[N]/D")


    for i in range(EVTMAX):
        #print ">>>",i
        i32[0] = i
        i64[0] = i
        u32[0] = i
        u64[0] = i

        f32[0] = i
        f64[0] = i
        for ii,vv in enumerate(bytes("evt-%03d" % i,"ascii")):
            s06[ii]=vv
            pass
        
        for jj in range(ARRAYSZ):
            arr_i32[jj] = i
            arr_i64[jj] = i
            arr_u32[jj] = i
            arr_u64[jj] = i
            arr_f32[jj] = i
            arr_f64[jj] = i

        n[0] = i % 10
        for jj in range(ARRAYSZ):
            sli_i32[jj] = i
            sli_i64[jj] = i
            sli_u32[jj] = i
            sli_u64[jj] = i
            sli_f32[jj] = i
            sli_f64[jj] = i


        t.Fill()
        pass

    f.Write()
    f.Close()

if __name__ == "__main__":
    main()
    
    
