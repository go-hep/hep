package slha_test

import (
	"fmt"
	"math"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/go-hep/slha"
)

func TestDecode(t *testing.T) {
	const fname = "testdata/sps1a.spc"
	f, err := os.Open(fname)
	if err != nil {
		t.Fatalf("error opening file [%s]: %v\n", fname, err)
	}
	defer f.Close()

	var data slha.SLHA
	err = slha.NewDecoder(f).Decode(&data)
	if err != nil {
		t.Fatalf("error decoding file [%s]: %v\n", fname, err)
	}
}

func TestEncode(t *testing.T) {
	const fname = "testdata/sps1a.spc"
	f, err := os.Open(fname)
	if err != nil {
		t.Fatalf("error opening file [%s]: %v\n", fname, err)
	}
	defer f.Close()

	var data slha.SLHA
	err = slha.NewDecoder(f).Decode(&data)
	if err != nil {
		t.Fatalf("error decoding file [%s]: %v\n", fname, err)
	}

	const ofname = "testdata/write.out"
	out, err := os.Create(ofname)
	if err != nil {
		t.Fatalf("error creating file [%s]: %v\n", ofname, err)
	}
	defer out.Close()

	err = slha.NewEncoder(out).Encode(&data)
	if err != nil {
		t.Fatalf("error encoding file [%s]: %v\n", ofname, err)
	}

	err = out.Close()
	if err != nil {
		t.Fatalf("error closing file [%s]: %v\n", ofname, err)
	}
}

func TestRW(t *testing.T) {
	for _, fname := range []string{
		"testdata/sps1a.spc",
		"testdata/ex1-snowmass-point-1a.slha",
		"testdata/slha1.txt",
		"testdata/slha2.txt",
	} {
		f, err := os.Open(fname)
		if err != nil {
			t.Fatalf("error opening file [%s]: %v\n", fname, err)
		}
		defer f.Close()

		var wdata slha.SLHA
		err = slha.NewDecoder(f).Decode(&wdata)
		if err != nil {
			t.Fatalf("error decoding file [%s]: %v\n", fname, err)
		}
		f.Close()

		ofname := fname + ".out"
		f, err = os.Create(ofname)
		if err != nil {
			t.Fatalf("error creating file [%s]: %v\n", ofname, err)
		}
		defer f.Close()

		err = slha.NewEncoder(f).Encode(&wdata)
		if err != nil {
			t.Fatalf("error encoding file [%s]: %v\n", ofname, err)
		}

		err = f.Close()
		if err != nil {
			t.Fatalf("error closing file [%s]: %v\n", ofname, err)
		}

		f, err = os.Open(ofname)
		if err != nil {
			t.Fatalf("error re-opening file [%s]: %v\n", ofname, err)
		}
		defer f.Close()

		var rdata slha.SLHA
		err = slha.NewDecoder(f).Decode(&rdata)
		if err != nil {
			t.Fatalf("error re-decoding file [%s]: %v\n", ofname, err)
		}
		f.Close()

		if ok, str := compareSLHA(rdata, wdata); !ok {
			t.Fatalf("error - SLHA data differ - file %s\n%s\n", ofname, str)
		} else {
			os.Remove(ofname)
		}
	}
}

func compareSLHA(a, b slha.SLHA) (bool, string) {
	str := make([]string, 0)
	ok := true
	if len(a.Blocks) != len(b.Blocks) {
		str = append(str,
			fmt.Sprintf("ref - #block: %d\n", len(a.Blocks)),
			fmt.Sprintf("chk - #block: %d\n", len(b.Blocks)),
		)
		ok = false
	} else {
		for i := range a.Blocks {
			ablock := &a.Blocks[i]
			bblock := &b.Blocks[i]

			if math.IsNaN(ablock.Q) {
				ablock.Q = -999
			}
			if math.IsNaN(bblock.Q) {
				bblock.Q = -999
			}

			if ablock.Name != bblock.Name {
				ok = false
				str = append(str,
					fmt.Sprintf("ref - block[%d] n=%q\n", i, ablock.Name),
					fmt.Sprintf("chk - block[%d] n=%q\n", i, bblock.Name),
				)
			}

			if ablock.Comment != bblock.Comment {
				ok = false
				str = append(str,
					fmt.Sprintf("ref - block[%d] n=%q comm=%q\n", i, ablock.Name, ablock.Comment),
					fmt.Sprintf("chk - block[%d] n=%q comm=%q\n", i, bblock.Name, bblock.Comment),
				)

			}

			if ablock.Q != bblock.Q {
				ok = false
				str = append(str,
					fmt.Sprintf("ref - block[%d] n=%q Q=%v\n", i, ablock.Name, ablock.Q),
					fmt.Sprintf("chk - block[%d] n=%q Q=%v\n", i, bblock.Name, bblock.Q),
				)
			}

			if len(ablock.Data) != len(bblock.Data) {
				str = append(str,
					fmt.Sprintf("ref - block[%d] n=%q #entries=%d\n", len(ablock.Data)),
					fmt.Sprintf("chk - block[%d] n=%q #entries=%d\n", len(bblock.Data)),
				)
				ok = false
			} else {
				for j := range ablock.Data {
					aa := &ablock.Data[j]
					bb := &bblock.Data[j]

					if !reflect.DeepEqual(aa.Index, bb.Index) {
						ok = false
						str = append(str,
							fmt.Sprintf("ref - block[%d][%d] index=%v\n", i, j, aa.Index.Index()),
							fmt.Sprintf("chk - block[%d][%d] index=%v\n", i, j, bb.Index.Index()),
						)
					}

					va := aa.Value.Interface()
					vb := bb.Value.Interface()
					if !reflect.DeepEqual(va, vb) {
						ok = false
						str = append(str,
							fmt.Sprintf("ref - block[%d][%d] v=%#v\n", i, j, va),
							fmt.Sprintf("chk - block[%d][%d] v=%#v\n", i, j, vb),
						)

					}

					ca := aa.Value.Comment()
					cb := bb.Value.Comment()
					if !reflect.DeepEqual(ca, cb) {
						ok = false
						str = append(str,
							fmt.Sprintf("ref - block[%d][%d] comm=%q\n", i, j, ca),
							fmt.Sprintf("chk - block[%d][%d] comm=%q\n", i, j, cb),
						)

					}

				}
			}
		}
	}

	if len(a.Particles) != len(b.Particles) {
		str = append(str,
			fmt.Sprintf("ref - #parts: %d\n", len(a.Particles)),
			fmt.Sprintf("chk - #parts: %d\n", len(b.Particles)),
		)
		ok = false
	} else {
		for i := range a.Particles {
			apart := &a.Particles[i]
			bpart := &b.Particles[i]

			if apart.PdgID != bpart.PdgID {
				ok = false
				str = append(str,
					fmt.Sprintf("ref - part[%d] pdgid=%d\n", i, apart.PdgID),
					fmt.Sprintf("chk - part[%d] pdgid=%q\n", i, bpart.PdgID),
				)
			}

			if math.IsNaN(apart.Width) {
				apart.Width = -999
			}
			if math.IsNaN(bpart.Width) {
				bpart.Width = -999
			}

			if math.IsNaN(apart.Mass) {
				apart.Mass = -999
			}
			if math.IsNaN(bpart.Mass) {
				bpart.Mass = -999
			}

			if !reflect.DeepEqual(apart, bpart) {
				ok = false
				str = append(str,
					fmt.Sprintf("ref - part[%d] pdgid=%d\n", i, apart.PdgID),
					fmt.Sprintf("chk - part[%d] pdgid=%d\n", i, bpart.PdgID),
				)
			}
		}
	}

	return ok, strings.Join(str, "\n")
}
