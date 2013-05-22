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

// now we build the graph, which will look like
//                       p7                         #
// p1                   /                           #
//   \v1__p3      p5---v4                           #
//         \_v3_/       \                           #
//         /    \        p8                         #
//    v2__p4     \                                  #
//   /            p6                                #
// p2                                               #
//                                                  #
package main

import (
	"github.com/go-hep/hepmc"
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
		SignalProcessId: 20,
		EventNumber:     1,
		Particles:       make(map[int]*hepmc.Particle),
		Vertices:        make(map[int]*hepmc.Vertex),
	}

	// define the units
	evt.MomentumUnit = hepmc.GEV
	evt.LengthUnit = hepmc.MM

	// create vertex 1 and 2, together with their in-particles
	v1 := &hepmc.Vertex{Barcode: -1}
	err = evt.AddVertex(v1)
	handle_err(err)

	err = v1.AddParticleIn(&hepmc.Particle{
		Momentum: hepmc.FourVector{0, 0, 7000, 7000},
		PdgId:    2212,
		Status:   3,
		Barcode:  1,
	})
	handle_err(err)

	v2 := &hepmc.Vertex{Barcode: -2}
	err = evt.AddVertex(v2)
	handle_err(err)

	err = v2.AddParticleIn(&hepmc.Particle{
		Momentum: hepmc.FourVector{0, 0, -7000, 7000},
		PdgId:    2212,
		Status:   3,
		Barcode:  2,
	})
	handle_err(err)

}

// EOF
