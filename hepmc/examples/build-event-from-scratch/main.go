// In this example we will place the following event into HepMC "by hand"
//
//     name status pdg_id  parent Px       Py    Pz       Energy      Mass
//  1  !p+!    3   2212    0,0    0.000    0.000 7000.000 7000.000    0.938
//  2  !p+!    3   2212    0,0    0.000    0.000-7000.000 7000.000    0.938
//=========================================================================
//  3  !d!     3      1    1,1    0.750   -1.569   32.191   32.238    0.000
//  4  !u~!    3     -2    2,2   -3.047  -19.000  -54.629   57.920    0.000
//  5  !W-!    3    -24    1,2    1.517   -20.68  -20.605   85.925   80.799
//  6  !gamma! 1     22    1,2   -3.813    0.113   -1.833    4.233    0.000
//  7  !d!     1      1    5,5   -2.445   28.816    6.082   29.552    0.010
//  8  !u~!    1     -2    5,5    3.962  -49.498  -26.687   56.373    0.006
//
// now we build the graph, which will look like
//  #                       p7                         #
//  # p1                   /                           #
//  #   \v1__p3      p5---v4                           #
//  #         \_v3_/       \                           #
//  #         /    \        p8                         #
//  #    v2__p4     \                                  #
//  #   /            p6                                #
//  # p2                                               #
//  #                                                  #
//
package main

import (
	"os"

	"go-hep.org/x/hep/fmom"
	"go-hep.org/x/hep/hepmc"
)

func handle_err(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error

	// first create the event container, with signal process 20, event number 1
	evt := hepmc.Event{
		SignalProcessID: 20,
		EventNumber:     1,
		Particles:       make(map[int]*hepmc.Particle),
		Vertices:        make(map[int]*hepmc.Vertex),
	}
	defer hepmc.Delete(&evt)

	// define the units
	evt.MomentumUnit = hepmc.GEV
	evt.LengthUnit = hepmc.MM

	// create vertex 1 and 2, together with their in-particles
	v1 := &hepmc.Vertex{}
	err = evt.AddVertex(v1)
	handle_err(err)

	err = v1.AddParticleIn(&hepmc.Particle{
		Momentum: fmom.PxPyPzE{0, 0, 7000, 7000},
		PdgID:    2212,
		Status:   3,
	})
	handle_err(err)

	v2 := &hepmc.Vertex{}
	err = evt.AddVertex(v2)
	handle_err(err)

	err = v2.AddParticleIn(&hepmc.Particle{
		Momentum: fmom.PxPyPzE{0, 0, -7000, 7000},
		PdgID:    2212,
		Status:   3,
		//Barcode:  2,
	})
	handle_err(err)

	// create the outgoing particles of v1 and v2
	p3 := &hepmc.Particle{
		Momentum: fmom.PxPyPzE{.750, -1.569, 32.191, 32.238},
		PdgID:    1,
		Status:   3,
		// Barcode: 3,
	}
	err = v1.AddParticleOut(p3)
	handle_err(err)

	p4 := &hepmc.Particle{
		Momentum: fmom.PxPyPzE{-3.047, -19., -54.629, 57.920},
		PdgID:    -2,
		Status:   3,
		// Barcode: 4,
	}
	err = v2.AddParticleOut(p4)
	handle_err(err)

	// create v3
	v3 := &hepmc.Vertex{}
	err = evt.AddVertex(v3)
	handle_err(err)

	err = v3.AddParticleIn(p3)
	handle_err(err)

	err = v3.AddParticleIn(p4)
	handle_err(err)

	err = v3.AddParticleOut(&hepmc.Particle{
		Momentum: fmom.PxPyPzE{-3.813, 0.113, -1.833, 4.233},
		PdgID:    22,
		Status:   1,
	})
	handle_err(err)

	p5 := &hepmc.Particle{
		Momentum: fmom.PxPyPzE{1.517, -20.68, -20.605, 85.925},
		PdgID:    -24,
		Status:   3,
	}
	err = v3.AddParticleOut(p5)
	handle_err(err)

	// create v4
	v4 := &hepmc.Vertex{
		Position: fmom.PxPyPzE{0.12, -0.3, 0.05, 0.004},
	}
	err = evt.AddVertex(v4)
	handle_err(err)

	err = v4.AddParticleIn(p5)
	handle_err(err)

	err = v4.AddParticleOut(&hepmc.Particle{
		Momentum: fmom.PxPyPzE{-2.445, 28.816, 6.082, 29.552},
		PdgID:    1,
		Status:   1,
	})
	handle_err(err)

	err = v4.AddParticleOut(&hepmc.Particle{
		Momentum: fmom.PxPyPzE{3.962, -49.498, -26.687, 56.373},
		PdgID:    -2,
		Status:   1,
	})
	handle_err(err)

	evt.SignalVertex = v3

	err = evt.Print(os.Stdout)
	handle_err(err)

	if len(os.Args) > 1 {
		fname := os.Args[1]
		out, err := os.Create(fname)
		handle_err(err)
		defer out.Close()

		enc := hepmc.NewEncoder(out)
		err = enc.Encode(&evt)
		handle_err(err)

		err = enc.Close()
		handle_err(err)

		err = out.Close()
		handle_err(err)
	}
}

// EOF
