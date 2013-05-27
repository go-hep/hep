package main

import (
	"flag"
	"fmt"
	"io"
	"math"
	"os"

	"github.com/go-hep/hepmc"
	"github.com/go-hep/lhef"
)

var (
	ifname = flag.String("i", "", "path to LHEF input file (default: STDIN)")
	ofname = flag.String("o", "", "path to HEPMC output file (default: STDOUT)")

	// in case IDWTUP == +/-4, one has to keep track of the accumulated
	// weights and event numbers to evaluate the cross section on-the-fly.
	// The last evaluation is the one used.
	// Better to be sure that crossSection() is never used to fill the
	// histograms, but only in the finalization stage, by reweighting the
	// histograms with crossSection()/sumOfWeights()
	acc_weight         = 0.0
	acc_weight_squared = 0.0
	evt_nbr            = 0
)

func main() {
	flag.Parse()

	var r io.Reader
	if *ifname == "" {
		r = os.Stdin
	} else {
		f, err := os.Open(*ifname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error: %v\n", err)
			os.Exit(1)
		}
		r = f
		defer f.Close()
	}

	var w io.Writer
	if *ofname == "" {
		w = os.Stdout
	} else {
		f, err := os.Create(*ofname)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error: %v\n", err)
			os.Exit(1)
		}
		w = f
		defer f.Close()
	}

	dec, err := lhef.NewDecoder(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "**error: %v\n", err)
		os.Exit(1)
	}

	enc := hepmc.NewEncoder(w)
	if enc == nil {
		fmt.Fprintf(os.Stderr, "**error: nil hepmc.Encoder\n")
		os.Exit(1)
	}
	defer enc.Close()

	for ievt := 0; ; ievt++ {
		lhevt, err := dec.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error at evt #%d: %v\n", ievt, err)
		}

		evt := hepmc.Event{
			EventNumber: ievt + 1,
			Particles:   make(map[int]*hepmc.Particle),
			Vertices:    make(map[int]*hepmc.Vertex),
			Weights:     hepmc.NewWeights(),
		}

		// define the units
		evt.MomentumUnit = hepmc.GEV
		evt.LengthUnit = hepmc.MM

		weight := lhevt.XWGTUP
		evt.Weights.Add("0", weight)

		xsecval := -1.0
		xsecerr := -1.0
		switch math.Abs(float64(dec.Run.IDWTUP)) {
		case 3:
			xsecval = dec.Run.XSECUP[0]
			xsecerr = dec.Run.XSECUP[1]

		case 4:
			acc_weight += weight
			acc_weight_squared += weight * weight
			evt_nbr += 1
			xsecval = acc_weight / float64(evt_nbr)
			xsecerr2 := (acc_weight_squared/float64(evt_nbr) - xsecval*xsecval) / float64(evt_nbr)

			if xsecerr2 < 0 {
				fmt.Fprintf(os.Stderr, "WARNING: xsecerr^2 < 0. forcing to zero. (%f)\n", xsecerr2)
				xsecerr2 = 0.
			}
			xsecerr = math.Sqrt(xsecerr2)

		default:
			fmt.Fprintf(
				os.Stderr,
				"**error: IDWTUP=%v value not handled yet.\n",
				dec.Run.IDWTUP,
			)
			os.Exit(1)
		}

		evt.CrossSection = &hepmc.CrossSection{Value: xsecval, Error: xsecerr}
		vtx := hepmc.Vertex{
			Event:   &evt,
			Barcode: -1,
		}
		p1 := hepmc.Particle{
			Momentum: hepmc.FourVector{
				0, 0,
				dec.Run.EBMUP[0],
				dec.Run.EBMUP[0],
			},
			PdgId:   dec.Run.IDBMUP[0],
			Status:  4,
			Barcode: 1,
		}
		p2 := hepmc.Particle{
			Momentum: hepmc.FourVector{
				0, 0,
				dec.Run.EBMUP[1],
				dec.Run.EBMUP[1],
			},
			PdgId:   dec.Run.IDBMUP[1],
			Status:  4,
			Barcode: 2,
		}
		err = vtx.AddParticleIn(&p1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error at evt #%d: %v\n", ievt, err)
			os.Exit(1)
		}
		err = vtx.AddParticleIn(&p2)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error at evt #%d: %v\n", ievt, err)
			os.Exit(1)
		}
		evt.Beams[0] = &p1
		evt.Beams[1] = &p2

		nmax := 2
		imax := int(lhevt.NUP)
		for i := 0; i < imax; i++ {
			if lhevt.ISTUP[i] != 1 {
				continue
			}
			nmax += 1
			vtx.AddParticleOut(&hepmc.Particle{
				Momentum: hepmc.FourVector{
					lhevt.PUP[i][0],
					lhevt.PUP[i][1],
					lhevt.PUP[i][2],
					lhevt.PUP[i][3],
				},
				GeneratedMass: lhevt.PUP[i][4],
				PdgId:         lhevt.IDUP[i],
				Status:        1,
				Barcode:       3 + i,
			})
		}
		err = evt.AddVertex(&vtx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error at evt #%d: %v\n", ievt, err)
			os.Exit(1)
		}

		nparts := len(evt.Particles)
		if nmax != nparts {
			fmt.Fprintf(os.Stderr, "**error at evt #%d: LHEF/HEPMC inconsistency (LHEF particles: %d, HEPMC particles: %d)\n", ievt, nmax, nparts)
			for _, p := range evt.Particles {
				fmt.Fprintf(os.Stderr, "part: %v\n", p)
			}
			fmt.Fprintf(os.Stderr, "p1: %v\n", p1)
			fmt.Fprintf(os.Stderr, "p2: %v\n", p2)
			os.Exit(1)
		}
		if len(evt.Vertices) != 1 {
			fmt.Fprintf(os.Stderr, "**error at evt #%d: inconsistent number of vertices in HEPMC (got %d, expected 1)\n", ievt, len(evt.Vertices))
			os.Exit(1)
		}
		//fmt.Fprintf(os.Stderr, "nparts: %v\nnverts: %v\n", len(evt.Particles), len(evt.Vertices))

		err = enc.Encode(&evt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "**error at evt #%d: %v\n", ievt, err)
			os.Exit(1)
		}
	}
}

// EOF
