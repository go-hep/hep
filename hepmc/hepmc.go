// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package hepmc is a pure Go implementation of the C++ HepMC-2 library.
package hepmc // import "go-hep.org/x/hep/hepmc"

import (
	"errors"
	"fmt"
	"io"
	"sort"

	"go-hep.org/x/hep/fmom"
)

var (
	errNilVtx      = errors.New("hepmc: nil Vertex")
	errNilParticle = errors.New("hepmc: nil Particle")
)

// Delete deletes an event and allows memory to be reclaimed by the garbage collector
func Delete(evt *Event) error {
	var err error
	if evt == nil {
		return err
	}
	if evt.SignalVertex != nil {
		evt.SignalVertex.Event = nil
	}
	evt.SignalVertex = nil
	evt.Beams[0] = nil
	evt.Beams[1] = nil

	for _, p := range evt.Particles {
		p.ProdVertex = nil
		p.EndVertex = nil
		p.Flow.Particle = nil
		delete(evt.Particles, p.Barcode)
	}
	for _, vtx := range evt.Vertices {
		vtx.Event = nil
		vtx.ParticlesIn = nil
		vtx.ParticlesOut = nil
		delete(evt.Vertices, vtx.Barcode)
	}

	evt.Particles = nil
	evt.Vertices = nil
	return err
}

// Event represents a record for MC generators (for use at any stage of generation)
//
// This type is intended as both a "container class" ( to store a MC
// event for interface between MC generators and detector simulation )
// and also as a "work in progress class" ( that could be used inside
// a generator and modified as the event is built ).
type Event struct {
	SignalProcessID int     // id of the signal process
	EventNumber     int     // event number
	Mpi             int     // number of multi particle interactions
	Scale           float64 // energy scale,
	AlphaQCD        float64 // QCD coupling, see hep-ph/0109068
	AlphaQED        float64 // QED coupling, see hep-ph/0109068

	SignalVertex *Vertex      // signal vertex
	Beams        [2]*Particle // incoming beams
	Weights      Weights      // weights for this event. first weight is used by default for hit and miss
	RandomStates []int64      // container of random number generator states

	Vertices  map[int]*Vertex
	Particles map[int]*Particle

	CrossSection *CrossSection
	HeavyIon     *HeavyIon
	PdfInfo      *PdfInfo
	MomentumUnit MomentumUnit
	LengthUnit   LengthUnit

	bcparts int // barcode suggestions for particles
	bcverts int // barcode suggestions for vertices
}

// Delete prepares this event for GC-reclaim
func (evt *Event) Delete() error {
	return Delete(evt)
}

// AddVertex adds a vertex to this event
func (evt *Event) AddVertex(vtx *Vertex) error {
	if vtx == nil {
		return errNilVtx
	}
	//TODO(sbinet): warn and remove from previous event
	//if vtx.Event != nil && vtx.Event != evt {
	//}
	return vtx.setParentEvent(evt)
}

// Print prints the event to w in a human-readable format
func (evt *Event) Print(w io.Writer) error {
	var err error
	const liner = ("________________________________________" +
		"________________________________________")

	_, err = fmt.Fprintf(w, "%s\n", liner)
	if err != nil {
		return err
	}

	sigVtx := 0
	if evt.SignalVertex != nil {
		sigVtx = evt.SignalVertex.Barcode
	}
	_, err = fmt.Fprintf(
		w,
		"GenEvent: #%04d ID=%5d SignalProcessGenVertex Barcode: %d\n",
		evt.EventNumber,
		evt.SignalProcessID,
		sigVtx,
	)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(
		w,
		" Momentum units:%8s     Position units:%8s\n",
		evt.MomentumUnit.String(),
		evt.LengthUnit.String(),
	)
	if err != nil {
		return err
	}

	if evt.CrossSection != nil {
		_, err = fmt.Fprintf(
			w,
			" Cross Section: %e +/- %e\n",
			evt.CrossSection.Value,
			evt.CrossSection.Error,
		)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fprintf(
		w,
		" Entries this event: %d vertices, %d particles.\n",
		len(evt.Vertices),
		len(evt.Particles),
	)
	if err != nil {
		return err
	}

	if evt.Beams[0] != nil && evt.Beams[1] != nil {
		_, err = fmt.Fprintf(
			w,
			" Beam Particle barcodes: %d %d \n",
			evt.Beams[0].Barcode,
			evt.Beams[1].Barcode,
		)
	} else {
		_, err = fmt.Fprintf(
			w,
			" Beam Particles are not defined.\n",
		)
	}
	if err != nil {
		return err
	}

	// random state
	_, err = fmt.Fprintf(
		w,
		" RndmState(%d)=",
		len(evt.RandomStates),
	)
	if err != nil {
		return err
	}

	for _, rnd := range evt.RandomStates {
		_, err = fmt.Fprintf(w, "%d ", rnd)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "\n")
	if err != nil {
		return err
	}

	// weights
	_, err = fmt.Fprintf(
		w,
		" Wgts(%d)=",
		len(evt.Weights.Map),
	)
	if err != nil {
		return err
	}

	for n := range evt.Weights.Map {
		_, err = fmt.Fprintf(w, "(%s,%f) ", n, evt.Weights.At(n))
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(
		w,
		" EventScale %7.5f [energy] \t alphaQCD=%8.8f\t alphaQED=%8.8f\n",
		evt.Scale,
		evt.AlphaQCD,
		evt.AlphaQED,
	)
	if err != nil {
		return err
	}

	// print a legend to describe the particle info
	_, err = fmt.Fprintf(
		w,
		"                                    GenParticle Legend\n"+
			"        Barcode   PDG ID      ( Px,       Py,       Pz,     E )"+
			" Stat  DecayVtx\n")
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(w, "%s\n", liner)
	if err != nil {
		return err
	}

	// print all vertices
	barcodes := make([]int, 0, len(evt.Vertices))
	for bc := range evt.Vertices {
		barcodes = append(barcodes, -bc)
	}
	sort.Ints(barcodes)
	for _, bc := range barcodes {
		vtx := evt.Vertices[-bc]
		err = vtx.Print(w)
		if err != nil {
			return err
		}
	}
	_, err = fmt.Fprintf(w, "%s\n", liner)
	if err != nil {
		return err
	}
	return err
}

func (evt *Event) removeVertex(bc int) {
	// TODO(sbinet) remove barcode from suggested vtx-barcodes ?
	delete(evt.Vertices, bc)
}

func (evt *Event) removeParticle(bc int) {
	// TODO(sbinet) remove barcode from suggested particle-barcodes ?
	delete(evt.Particles, bc)
}

func (evt *Event) vBarcode() int {
	// TODO(sbinet) make goroutine-safe
	evt.bcverts++
	return -evt.bcverts
}

func (evt *Event) pBarcode() int {
	// TODO(sbinet) make goroutine-safe
	evt.bcparts++
	return evt.bcparts
}

// Particle represents a generator particle within an event coming in/out of a vertex
//
// Particle is the basic building block of the event record
type Particle struct {
	Momentum      fmom.PxPyPzE // momentum vector
	PdgID         int64        // id according to PDG convention
	Status        int          // status code as defined for HEPEVT
	Flow          Flow         // flow of this particle
	Polarization  Polarization // polarization of this particle
	ProdVertex    *Vertex      // pointer to production vertex (nil if vacuum or beam)
	EndVertex     *Vertex      // pointer to decay vertex (nil if not-decayed)
	Barcode       int          // unique identifier in the event
	GeneratedMass float64      // mass of this particle when it was generated
}

func (p *Particle) dump(w io.Writer) error {
	var err error
	_, err = fmt.Fprintf(
		w,
		" %9d%9d %+9.2e,%+9.2e,%+9.2e,%+9.2e",
		p.Barcode,
		p.PdgID,
		p.Momentum.Px(),
		p.Momentum.Py(),
		p.Momentum.Pz(),
		p.Momentum.E(),
	)
	if err != nil {
		return err
	}

	switch {
	case p.EndVertex != nil && p.EndVertex.Barcode != 0:
		_, err = fmt.Fprintf(w, "%4d %9d", p.Status, p.EndVertex.Barcode)
	case p.EndVertex == nil:
		_, err = fmt.Fprintf(w, "%4d", p.Status)
	default:
		_, err = fmt.Fprintf(w, "%4d %p", p.Status, p.EndVertex)
	}
	return err
}

// Particles is a []*Particle sorted by increasing-barcodes
type Particles []*Particle

func (ps Particles) Len() int {
	return len(ps)
}
func (ps Particles) Less(i, j int) bool {
	return ps[i].Barcode < ps[j].Barcode
}
func (ps Particles) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// Vertices is a []*Vertex sorted by increasing-barcodes
type Vertices []*Vertex

func (ps Vertices) Len() int {
	return len(ps)
}
func (ps Vertices) Less(i, j int) bool {
	return ps[i].Barcode < ps[j].Barcode
}
func (ps Vertices) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

// Vertex represents a generator vertex within an event
// A vertex is indirectly (via particle "edges") linked to other
//
//	vertices ("nodes") to form a composite "graph"
type Vertex struct {
	Position     fmom.PxPyPzE // 4-vector of vertex [mm]
	ParticlesIn  []*Particle  // all incoming particles
	ParticlesOut []*Particle  // all outgoing particles
	ID           int          // vertex id
	Weights      Weights      // weights for this vertex
	Event        *Event       // pointer to event owning this vertex
	Barcode      int          // unique identifier in the event
}

func (vtx *Vertex) setParentEvent(evt *Event) error {
	var err error
	origEvt := vtx.Event
	vtx.Event = evt
	// if orig_evt == evt {
	// 	return err
	// }
	if evt != nil {
		if vtx.Barcode == 0 {
			vtx.Barcode = evt.vBarcode()
		}
		evt.Vertices[vtx.Barcode] = vtx
	}
	if origEvt != nil && origEvt != evt {
		origEvt.removeVertex(vtx.Barcode)
	}
	// we also need to loop over all the particles which are owned by
	// this vertex and remove their barcodes from the old event.
	for _, p := range vtx.ParticlesIn {
		if p.ProdVertex == nil {
			if evt != nil {
				evt.Particles[p.Barcode] = p
			}
			if origEvt != nil && origEvt != evt {
				origEvt.removeParticle(p.Barcode)
			}
		}
	}

	for _, p := range vtx.ParticlesOut {
		if evt != nil {
			evt.Particles[p.Barcode] = p
		}
		if origEvt != nil && origEvt != evt {
			origEvt.removeParticle(p.Barcode)
		}
	}
	return err
}

// AddParticleIn adds a particle to the list of in-coming particles to this vertex
func (vtx *Vertex) AddParticleIn(p *Particle) error {
	var err error
	if p == nil {
		return errNilParticle
	}
	// if p had a decay vertex, remove it from that vertex's list
	if p.EndVertex != nil {
		err = p.EndVertex.removeParticleIn(p)
		if err != nil {
			return err
		}
	}
	// make sure we don't add it twice...
	err = vtx.removeParticleIn(p)
	if err != nil {
		return err
	}
	if p.Barcode == 0 {
		p.Barcode = vtx.Event.pBarcode()
	}
	p.EndVertex = vtx
	vtx.ParticlesIn = append(vtx.ParticlesIn, p)
	vtx.Event.Particles[p.Barcode] = p
	return err
}

// AddParticleOut adds a particle to the list of out-going particles to this vertex
func (vtx *Vertex) AddParticleOut(p *Particle) error {
	var err error
	if p == nil {
		return errNilParticle
	}
	// if p had a production vertex, remove it from that vertex's list
	if p.ProdVertex != nil {
		err = p.ProdVertex.removeParticleOut(p)
		if err != nil {
			return err
		}
	}
	// make sure we don't add it twice...
	err = vtx.removeParticleOut(p)
	if err != nil {
		return err
	}
	if p.Barcode == 0 {
		p.Barcode = vtx.Event.pBarcode()
	}
	p.ProdVertex = vtx
	vtx.ParticlesOut = append(vtx.ParticlesOut, p)
	vtx.Event.Particles[p.Barcode] = p
	return err
}

func (vtx *Vertex) removeParticleIn(p *Particle) error {
	var err error
	nparts := len(vtx.ParticlesIn)
	switch nparts {
	case 0:
		//FIXME: logical error ?
		return err
	}
	idx := -1
	for i, pp := range vtx.ParticlesIn {
		if pp == p {
			idx = i
			break
		}
	}
	if idx >= 0 {
		copy(vtx.ParticlesIn[idx:], vtx.ParticlesIn[idx+1:])
		vtx.ParticlesIn[len(vtx.ParticlesIn)-1] = nil
		vtx.ParticlesIn = vtx.ParticlesIn[:len(vtx.ParticlesIn)-1]
	}
	return err
}

func (vtx *Vertex) removeParticleOut(p *Particle) error {
	var err error
	nparts := len(vtx.ParticlesOut)
	switch nparts {
	case 0:
		//FIXME: logical error ?
		return err
	}
	idx := -1
	for i, pp := range vtx.ParticlesOut {
		if pp == p {
			idx = i
			break
		}
	}
	if idx >= 0 {
		copy(vtx.ParticlesOut[idx:], vtx.ParticlesOut[idx+1:])
		n := len(vtx.ParticlesOut)
		vtx.ParticlesOut[n-1] = nil
		vtx.ParticlesOut = vtx.ParticlesOut[:n-1]
	}
	return err
}

// Print prints the vertex to w in a human-readable format
func (vtx *Vertex) Print(w io.Writer) error {
	var (
		err  error
		zero fmom.PxPyPzE
	)
	if vtx.Barcode != 0 {
		if vtx.Position != zero {
			_, err = fmt.Fprintf(
				w,
				"Vertex:%9d ID:%5d (X,cT)=%+9.2e,%+9.2e,%+9.2e,%+9.2e\n",
				vtx.Barcode,
				vtx.ID,
				vtx.Position.X(),
				vtx.Position.Y(),
				vtx.Position.Z(),
				vtx.Position.T(),
			)
		} else {
			_, err = fmt.Fprintf(
				w,
				"GenVertex:%9d ID:%5d (X,cT):0\n",
				vtx.Barcode,
				vtx.ID,
			)
		}
	} else {
		// if the vertex doesn't have a unique barcode assigned, then
		// we print its memory address instead.. so that the
		// print out gives us a unique tag for the particle.
		if vtx.Position != zero {
			_, err = fmt.Fprintf(
				w,
				"Vertex:%p ID:%5d (X,cT)=%+9.2e,%+9.2e,%+9.2e,%+9.2e\n",
				vtx,
				vtx.ID,
				vtx.Position.X(),
				vtx.Position.Y(),
				vtx.Position.Z(),
				vtx.Position.T(),
			)

		} else {
			_, err = fmt.Fprintf(
				w,
				"GenVertex:%9d ID:%5d (X,cT):0\n",
				vtx.Barcode,
				vtx.ID,
			)
		}
	}
	if err != nil {
		return err
	}

	// print the weights if any
	if len(vtx.Weights.Slice) > 0 {
		_, err = fmt.Fprintf(w, " Wgts(%d)=", len(vtx.Weights.Slice))
		if err != nil {
			return err
		}
		for _, weight := range vtx.Weights.Slice {
			_, err = fmt.Fprintf(w, "%e ", weight)
			if err != nil {
				return err
			}
		}
		_, err = fmt.Fprintf(w, "\n")
		if err != nil {
			return err
		}
	}
	// sort incoming particles by barcode
	sort.Sort(Particles(vtx.ParticlesIn))

	// print out all incoming particles
	for i, p := range vtx.ParticlesIn {
		if i == 0 {
			_, err = fmt.Fprintf(w, " I:%2d", len(vtx.ParticlesIn))
		} else {
			_, err = fmt.Fprintf(w, "     ")
		}
		if err != nil {
			return err
		}
		err = p.dump(w)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "\n")
		if err != nil {
			return err
		}
	}

	// sort outgoing particles by barcode
	sort.Sort(Particles(vtx.ParticlesOut))

	// print out all outgoing particles
	for i, p := range vtx.ParticlesOut {
		if i == 0 {
			_, err = fmt.Fprintf(w, " O:%2d", len(vtx.ParticlesOut))
		} else {
			_, err = fmt.Fprintf(w, "     ")
		}
		if err != nil {
			return err
		}
		err = p.dump(w)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "\n")
		if err != nil {
			return err
		}
	}
	return err
}

// HeavyIon holds additional information for heavy-ion collisions
type HeavyIon struct {
	NCollHard         int     // number of hard scatterings
	NPartProj         int     // number of projectile participants
	NPartTarg         int     // number of target participants
	NColl             int     // number of NN (nucleon-nucleon) collisions
	NNwColl           int     // Number of N-Nwounded collisions
	NwNColl           int     // Number of Nwounded-N collisons
	NwNwColl          int     // Number of Nwounded-Nwounded collisions
	SpectatorNeutrons int     // Number of spectators neutrons
	SpectatorProtons  int     // Number of spectators protons
	ImpactParameter   float32 // Impact Parameter(fm) of collision
	EventPlaneAngle   float32 // Azimuthal angle of event plane
	Eccentricity      float32 // eccentricity of participating nucleons in the transverse plane (as in phobos nucl-ex/0510031)
	SigmaInelNN       float32 // nucleon-nucleon inelastic (including diffractive) cross-section
}

// CrossSection is used to store the generated cross section.
// This type is meant to be used to pass, on an event by event basis,
// the current best guess of the total cross section.
// It is expected that the final cross section will be stored elsewhere.
type CrossSection struct {
	Value float64 // value of the cross-section (in pb)
	Error float64 // error on the value of the cross-section (in pb)
	//IsSet bool
}

// PdfInfo holds informations about the partons distribution functions
type PdfInfo struct {
	ID1      int     // flavour code of first parton
	ID2      int     // flavour code of second parton
	LHAPdf1  int     // LHA PDF id of first parton
	LHAPdf2  int     // LHA PDF id of second parton
	X1       float64 // fraction of beam momentum carried by first parton ("beam side")
	X2       float64 // fraction of beam momentum carried by second parton ("target side")
	ScalePDF float64 // Q-scale used in evaluation of PDF's   (in GeV)
	Pdf1     float64 // PDF (id1, x1, Q)
	Pdf2     float64 // PDF (id2, x2, Q)
}

// Flow represents a particle's flow and keeps track of an arbitrary number of flow patterns within a graph (i.e. color flow, charge flow, lepton number flow,...)
//
// Flow patterns are coded with an integer, in the same manner as in Herwig.
// Note: 0 is NOT allowed as code index nor as flow code since it
// is used to indicate null.
//
// This class can be used to keep track of flow patterns within
// a graph. An example is color flow. If we have two quarks going through
// an s-channel gluon to form two more quarks:
//
//	\q1       /q3   then we can keep track of the color flow with the
//	 \_______/      HepMC::Flow class as follows:
//	 /   g   \.
//	/q2       \q4
//
// lets say the color flows from q2-->g-->q3  and q1-->g-->q4
// the individual colors are unimportant, but the flow pattern is.
// We can capture this flow by assigning the first pattern (q2-->g-->q3)
// a unique (arbitrary) flow code 678 and the second pattern (q1-->g-->q4)
// flow code 269  ( you can ask HepMC::Flow to choose
// a unique code for you using Flow::set_unique_icode() ).
// these codes with the particles as follows:
//
//	q2->flow().set_icode(1,678);
//	g->flow().set_icode(1,678);
//	q3->flow().set_icode(1,678);
//	q1->flow().set_icode(1,269);
//	g->flow().set_icode(2,269);
//	q4->flow().set_icode(1,269);
//
// later on if we wish to know the color partner of q1 we can ask for a list
// of all particles connected via this code to q1 which do have less than
// 2 color partners using:
//
//	vector<GenParticle*> result=q1->dangling_connected_partners(q1->icode(1),1,2);
//
// this will return a list containing q1 and q4.
//
//	vector<GenParticle*> result=q1->connected_partners(q1->icode(1),1,2);
//
// would return a list containing q1, g, and q4.
type Flow struct {
	Particle *Particle   // the particle this flow describes
	Icode    map[int]int // flow patterns as (code_index, icode)
}

// Polarization holds informations about a particle's polarization
type Polarization struct {
	Theta float64 // polar angle of polarization in radians [0, math.Pi)
	Phi   float64 // azimuthal angle of polarization in radians [0, 2*math.Pi)
}

// Weights holds informations about the event's and vertices' generation weights.
type Weights struct {
	Slice []float64      // the slice of weight values
	Map   map[string]int // the map of name->index-in-the-slice
}

// Add adds a new weight with name n and value v.
func (w *Weights) Add(n string, v float64) error {
	_, ok := w.Map[n]
	if ok {
		return fmt.Errorf("hepmc.Weights.Add: name [%s] already in container", n)
	}
	idx := len(w.Slice)
	w.Map[n] = idx
	w.Slice = append(w.Slice, v)
	return nil
}

// At returns the weight's value named n.
func (w Weights) At(n string) float64 {
	idx, ok := w.Map[n]
	if ok {
		return w.Slice[idx]
	}
	panic("hepmc.Weights.At: invalid name [" + n + "]")
}

// NewWeights creates a new set of weights.
func NewWeights() Weights {
	return Weights{
		Slice: make([]float64, 0, 1),
		Map:   make(map[string]int),
	}
}

// EOF
