// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package hepmc

import (
	"bufio"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"
)

type rstream struct {
	tokens tokens
	err    error
}

// Decoder decodes a hepmc Event from a stream
type Decoder struct {
	stream     chan rstream
	seenEvtHdr bool
	ftype      hepmcFileType

	sigProcBc int // barcode of signal vertex
	bp1       int // barcode of beam1
	bp2       int // barcode of beam2
}

// NewDecoder returns a new hepmc Decoder that reads from the io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	dec := &Decoder{
		stream: make(chan rstream),
	}
	go dec.readlines(bufio.NewReader(r))

	return dec
}

func (dec *Decoder) readlines(r *bufio.Reader) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		dec.stream <- rstream{
			tokens: newtokens(strings.Split(s.Text(), " ")),
			err:    nil,
		}
	}

	err := s.Err()
	if err == nil {
		err = io.EOF
	}

	dec.stream <- rstream{
		err: err,
	}
}

func (dec *Decoder) readline() (tokens, error) {
	state := <-dec.stream
	return state.tokens, state.err
}

// Decode reads the next value from the stream and stores it into evt.
func (dec *Decoder) Decode(evt *Event) error {
	var err error

	// search for event listing key
	if !dec.seenEvtHdr {
		err = dec.findFileType()
		if err != nil {
			return err
		}
		dec.seenEvtHdr = true
	}

	dec.sigProcBc = 0
	dec.bp1 = 0
	dec.bp2 = 0

	// test to be sure the next entry is of type 'E'
	// FIXME: implement more of the logic from HepMC::IO_GenEvent
	tokens, err := dec.readline()
	if err != nil {
		return err
	}

	peek := tokens.toks[0]
	if peek[0] != 'E' {
		err = dec.findEndKey(tokens)
		if err == io.EOF {
			return err
		}
		if err != nil {
			err = fmt.Errorf("hepmc.decode: invalid file (expected 'E' got '%v'. line=%q)", string(peek[0]), tokens)
			return err
		}
	}

	nVtx := 0
loop:
	for {
		switch tokens.at(0)[0] {
		case 'E':
			// call appropriate decoder method
			switch dec.ftype {
			case hepmcGenEvent:
				err = dec.decodeEvent(evt, &nVtx, tokens)
			case hepmcASCII:
				err = dec.decodeASCII(evt, &nVtx, tokens)
			case hepmcExtendedASCII:
				err = dec.decodeExtendedASCII(evt, &nVtx, tokens)
			case hepmcASCIIPdt:
				err = fmt.Errorf("hepmc.decode: HepMC::IO_Ascii-PARTICLE_DATA is NOT implemented (yet)")
			case hepmcExtendedASCIIPdt:
				err = fmt.Errorf("hepmc.decode: HepMC::IO_ExtendedAscii-PARTICLE_DATA is NOT implemented (yet)")
			default:
				err = fmt.Errorf("hepmc.decode: unknown file format (%v)", dec.ftype)
			}
			if err != nil {
				return err
			}
		case 'N':
			_ = tokens.next() // header 'N'
			nWeights, err := tokens.int()
			if err != nil {
				return err
			}
			names := make(map[string]int, nWeights)
			for i := 0; i < nWeights; i++ {
				nn, err := strconv.Unquote(tokens.next())
				if err != nil {
					return err
				}
				names[nn] = i
			}

			//fmt.Printf("== weights: %v\n", names)
			evt.Weights.Map = names

		case 'U':
			if dec.ftype == hepmcGenEvent {
				err = dec.decodeUnits(evt, tokens)
				if err != nil {
					return err
				}
			}
		case 'C':
			err = dec.decodeCrossSection(evt, tokens)
			if err != nil {
				return err
			}
		case 'H':
			switch dec.ftype {
			case hepmcGenEvent, hepmcExtendedASCII:
				var hi HeavyIon
				err = dec.decodeHeavyIon(&hi, tokens)
				if err != nil {
					return err
				}
				evt.HeavyIon = &hi
			}
		case 'F':
			switch dec.ftype {
			case hepmcGenEvent, hepmcExtendedASCII:
				var pdf PdfInfo
				err = dec.decodePdfInfo(&pdf, tokens)
				if err != nil {
					return err
				}
				evt.PdfInfo = &pdf
			}
		case 'V', 'P':
			break loop

		default:

			return fmt.Errorf(
				"hepmc.decoder: invalid file (got '%v')",
				peek[0],
			)
		}

		tokens, err = dec.readline()
		if err != nil {
			return err
		}
	}
	// dec.r.Read(peek[0:])
	// end vertices of particles are not connected until
	// after the full event is read
	//  => store the values in a map until then
	pidxToEndVtx := make(map[int]int, nVtx) // particle-idx to end_vtx barcode

	// decode the vertices
	for i := 0; i < nVtx; i++ {
		if i != 0 {
			tokens, err = dec.readline()
			if err != nil {
				return err
			}
		}
		vtx := &Vertex{}
		vtx.Event = evt
		err = dec.decodeVertex(evt, vtx, pidxToEndVtx, tokens)
		if err != nil {
			return err
		}
		evt.Vertices[vtx.Barcode] = vtx
		sort.Sort(Particles(vtx.ParticlesOut))
	}

	// set the signal process vertex
	if dec.sigProcBc != 0 {
		for _, vtx := range evt.Vertices {
			if vtx.Barcode == dec.sigProcBc {
				evt.SignalVertex = vtx
				break
			}
		}
		if evt.SignalVertex == nil {
			return fmt.Errorf("hepmc.decode: could not find signal vertex (barcode=%d)", dec.sigProcBc)
		}
	}

	// connect particles to their end vertices
	for i, endVtxBc := range pidxToEndVtx {
		p := evt.Particles[i]
		vtx := evt.Vertices[endVtxBc]
		vtx.ParticlesIn = append(vtx.ParticlesIn, p)
		p.EndVertex = vtx
		// also look for the beam particles
		if p.Barcode == dec.bp1 {
			evt.Beams[0] = p
		}
		if p.Barcode == dec.bp2 {
			evt.Beams[1] = p
		}
		sort.Sort(Particles(vtx.ParticlesIn))
	}
	return err
}

func (dec *Decoder) findFileType() error {

	for {
		tokens, err := dec.readline()
		if err != nil {
			if err == io.EOF {
				err = io.ErrUnexpectedEOF
			}
			return err
		}
		if len(tokens.toks) <= 0 {
			continue
		}
		line := tokens.next()
		switch line {
		case "":
			// no-op

		case startGenEvent:
			dec.ftype = hepmcGenEvent
			return nil

		case startASCII:
			dec.ftype = hepmcASCII
			return nil

		case startExtendedASCII:
			dec.ftype = hepmcExtendedASCII
			return nil

		case startPdt:
			dec.ftype = hepmcASCIIPdt
			return nil

		case startExtendedASCIIPdt:
			dec.ftype = hepmcExtendedASCIIPdt
			return nil
		}
	}
}

func (dec *Decoder) findEndKey(tokens tokens) error {
	var err error = io.EOF
	line := tokens.next()
	if line[0] != 'H' {
		err = fmt.Errorf("hepmc.decode: not an end-key (line=%q)", line)
		return err
	}

	var ftype hepmcFileType
	switch line {
	default:
		err = fmt.Errorf("hepmc.decode: invalid file type (value=%q)", line)
		return err

	case endGenEvent:
		ftype = hepmcGenEvent

	case endASCII:
		ftype = hepmcASCII

	case endExtendedASCII:
		ftype = hepmcExtendedASCII

	case endPdt:
		ftype = hepmcASCIIPdt

	case endExtendedASCIIPdt:
		ftype = hepmcExtendedASCIIPdt
	}

	if ftype != dec.ftype {
		err = fmt.Errorf(
			"hepmc.decode: file type changed from %v to %v",
			dec.ftype, ftype,
		)
	}
	return err
}

func (dec *Decoder) decodeUnits(evt *Event, tokens tokens) error {
	var err error
	peek := tokens.next()
	if peek[0] != 'U' {
		return fmt.Errorf("hepmc.decode: expected 'U'. got '%s'", string(peek[0]))
	}

	momUnit := tokens.next()
	evt.MomentumUnit, err = MomentumUnitFromString(momUnit)
	if err != nil {
		return err
	}

	lenUnit := tokens.next()
	evt.LengthUnit, err = LengthUnitFromString(lenUnit)
	if err != nil {
		return err
	}
	return err
}

func (dec *Decoder) decodeCrossSection(evt *Event, tokens tokens) error {
	var err error
	peek := tokens.next()
	if peek[0] != 'C' {
		return fmt.Errorf("hepmc.decode: expected 'C'. got '%s'", string(peek[0]))
	}
	var x CrossSection
	x.Value, err = tokens.float64()
	if err != nil {
		return err
	}
	x.Error, err = tokens.float64()
	if err != nil {
		return err
	}
	evt.CrossSection = &x
	return err
}

func (dec *Decoder) decodeEvent(evt *Event, nVtx *int, tokens tokens) error {
	var (
		err       error
		evtNbr    int
		mpi       int
		scale     float64
		aqcd      float64
		aqed      float64
		sigProcID int
		nRndm     int
		nWeights  int
	)

	peek := tokens.next()
	if peek[0] != 'E' {
		return fmt.Errorf("hepmc.decode: expected 'E'. got '%s'", string(peek[0]))
	}

	evtNbr, err = tokens.int()
	if err != nil {
		return err
	}

	mpi, err = tokens.int()
	if err != nil {
		return err
	}

	scale, err = tokens.float64()
	if err != nil {
		return err
	}

	aqcd, err = tokens.float64()
	if err != nil {
		return err
	}

	aqed, err = tokens.float64()
	if err != nil {
		return err
	}

	sigProcID, err = tokens.int()
	if err != nil {
		return err
	}

	dec.sigProcBc, err = tokens.int()
	if err != nil {
		return err
	}

	*nVtx, err = tokens.int()
	if err != nil {
		return err
	}

	dec.bp1, err = tokens.int()
	if err != nil {
		return err
	}

	dec.bp2, err = tokens.int()
	if err != nil {
		return err
	}

	nRndm, err = tokens.int()
	if err != nil {
		return err
	}

	rndmStates := make([]int64, nRndm)
	for i := 0; i < nRndm; i++ {
		rndmStates[i], err = tokens.int64()
		if err != nil {
			return err
		}
	}

	nWeights, err = tokens.int()
	if err != nil {
		return err
	}

	weights := make([]float64, nWeights)
	for i := 0; i < nWeights; i++ {
		weights[i], err = tokens.float64()
		if err != nil {
			return err
		}
	}

	// fill infos gathered so far
	evt.SignalProcessID = sigProcID
	evt.EventNumber = evtNbr
	evt.Mpi = mpi
	if evt.Weights.Slice == nil {
		evt.Weights = NewWeights()
	}
	evt.Weights.Slice = weights
	evt.RandomStates = rndmStates
	evt.Scale = scale
	evt.AlphaQCD = aqcd
	evt.AlphaQED = aqed

	evt.Vertices = make(map[int]*Vertex, *nVtx)
	evt.Particles = make(map[int]*Particle, *nVtx*2)
	return err
}

func (dec *Decoder) decodeVertex(evt *Event, vtx *Vertex, pidxToEndVtx map[int]int, tokens tokens) error {

	var err error
	peek := tokens.next()
	if peek[0] != 'V' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'V', got '%v') line=%q",
			peek[0],
			tokens,
		)
	}

	orphans := 0
	nPartsOut := 0
	nWeights := 0

	vtx.Barcode, err = tokens.int()
	if err != nil {
		return err
	}

	vtx.ID, err = tokens.int()
	if err != nil {
		return err
	}

	vtx.Position[0], err = tokens.float64()
	if err != nil {
		return err
	}

	vtx.Position[1], err = tokens.float64()
	if err != nil {
		return err
	}

	vtx.Position[2], err = tokens.float64()
	if err != nil {
		return err
	}

	vtx.Position[3], err = tokens.float64()
	if err != nil {
		return err
	}

	orphans, err = tokens.int()
	if err != nil {
		return err
	}

	nPartsOut, err = tokens.int()
	if err != nil {
		return err
	}

	nWeights, err = tokens.int()
	if err != nil {
		return err
	}

	// FIXME: reuse buffers ?
	vtx.Weights.Slice = make([]float64, nWeights)
	for i := 0; i < nWeights; i++ {
		vtx.Weights.Slice[i], err = tokens.float64()
		if err != nil {
			return err
		}
	}

	// read and create the associated particles
	// outgoing particles are added to their production vertices immediately.
	// incoming particles are added to a map and handled later.
	for i := 0; i < orphans; i++ {
		p := &Particle{}
		tokens, err = dec.readline()
		if err != nil {
			return err
		}
		err = dec.decodeParticle(evt, p, pidxToEndVtx, tokens)
		if err != nil {
			return err
		}
		evt.Particles[p.Barcode] = p
	}
	// FIXME: reuse buffers ?
	vtx.ParticlesOut = make([]*Particle, nPartsOut)
	for i := 0; i < nPartsOut; i++ {
		p := &Particle{ProdVertex: vtx}
		tokens, err = dec.readline()
		if err != nil {
			return err
		}
		err = dec.decodeParticle(evt, p, pidxToEndVtx, tokens)
		if err != nil {
			return err
		}
		evt.Particles[p.Barcode] = p
		vtx.ParticlesOut[i] = p
	}
	return err
}

func (dec *Decoder) decodeParticle(evt *Event, p *Particle, pidxToEndVtx map[int]int, tokens tokens) error {

	var err error
	peek := tokens.next()
	if peek[0] != 'P' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'P', got '%v')",
			peek[0],
		)
	}
	endBc := 0

	p.Barcode, err = tokens.int()
	if err != nil {
		return err
	}

	p.PdgID, err = tokens.int64()
	if err != nil {
		return err
	}

	p.Momentum[0], err = tokens.float64()
	if err != nil {
		return err
	}

	p.Momentum[1], err = tokens.float64()
	if err != nil {
		return err
	}

	p.Momentum[2], err = tokens.float64()
	if err != nil {
		return err
	}

	p.Momentum[3], err = tokens.float64()
	if err != nil {
		return err
	}

	if dec.ftype != hepmcASCII {
		p.GeneratedMass, err = tokens.float64()
		if err != nil {
			return err
		}
	}

	p.Status, err = tokens.int()
	if err != nil {
		return err
	}

	p.Polarization.Theta, err = tokens.float64()
	if err != nil {
		return err
	}

	p.Polarization.Phi, err = tokens.float64()
	if err != nil {
		return err
	}

	endBc, err = tokens.int()
	if err != nil {
		return err
	}

	err = dec.decodeFlow(&p.Flow, &tokens)
	if err != nil {
		return err
	}
	p.Flow.Particle = p

	// all particles are connected to their end vertex separately
	// after all particles and vertices have been created
	if endBc != 0 {
		pidxToEndVtx[p.Barcode] = endBc
	}
	return err
}

func (dec *Decoder) decodeFlow(flow *Flow, tokens *tokens) error {
	nFlow, err := tokens.int()
	if err != nil {
		return err
	}
	flow.Icode = make(map[int]int, nFlow)
	for i := 0; i < nFlow; i++ {
		k, err := tokens.int()
		if err != nil {
			return err
		}
		v, err := tokens.int()
		if err != nil {
			return err
		}
		flow.Icode[k] = v
	}
	return err
}

func (dec *Decoder) decodeHeavyIon(hi *HeavyIon, tokens tokens) error {
	var err error
	peek := tokens.next()
	if peek[0] != 'H' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'H', got '%v')",
			peek[0],
		)
	}

	hi.NCollHard, err = tokens.int()
	if err != nil {
		return err
	}

	hi.NPartProj, err = tokens.int()
	if err != nil {
		return err
	}

	hi.NPartTarg, err = tokens.int()
	if err != nil {
		return err
	}

	hi.NColl, err = tokens.int()
	if err != nil {
		return err
	}

	hi.NNwColl, err = tokens.int()
	if err != nil {
		return err
	}

	hi.NwNColl, err = tokens.int()
	if err != nil {
		return err
	}

	hi.NwNwColl, err = tokens.int()
	if err != nil {
		return err
	}

	hi.SpectatorNeutrons, err = tokens.int()
	if err != nil {
		return err
	}

	hi.SpectatorProtons, err = tokens.int()
	if err != nil {
		return err
	}

	hi.ImpactParameter, err = tokens.float32()
	if err != nil {
		return err
	}

	hi.EventPlaneAngle, err = tokens.float32()
	if err != nil {
		return err
	}

	hi.Eccentricity, err = tokens.float32()
	if err != nil {
		return err
	}

	hi.SigmaInelNN, err = tokens.float32()
	if err != nil {
		return err
	}

	return err
}

func (dec *Decoder) decodePdfInfo(pdf *PdfInfo, tokens tokens) error {
	var err error
	peek := tokens.next()
	if peek[0] != 'F' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'F', got '%v')",
			peek[0],
		)
	}

	pdf.ID1, err = tokens.int()
	if err != nil {
		return err
	}

	pdf.ID2, err = tokens.int()
	if err != nil {
		return err
	}

	pdf.X1, err = tokens.float64()
	if err != nil {
		return err
	}

	pdf.X2, err = tokens.float64()
	if err != nil {
		return err
	}

	pdf.ScalePDF, err = tokens.float64()
	if err != nil {
		return err
	}

	pdf.Pdf1, err = tokens.float64()
	if err != nil {
		return err

	}

	pdf.Pdf2, err = tokens.float64()
	if err != nil {
		return err
	}

	pdf.LHAPdf1, err = tokens.int()
	if err != nil {
		return err
	}

	pdf.LHAPdf2, err = tokens.int()
	if err != nil {
		return err
	}

	return err
}

func (dec *Decoder) decodeASCII(evt *Event, nVtx *int, tokens tokens) error {
	var (
		err       error
		evtNbr    int
		mpi       int
		scale     float64
		aqcd      float64
		aqed      float64
		sigProcID int
		nRndm     int
		nWeights  int
	)

	evtNbr, err = tokens.int()
	if err != nil {
		return err
	}

	scale, err = tokens.float64()
	if err != nil {
		return err
	}

	aqcd, err = tokens.float64()
	if err != nil {
		return err
	}

	aqed, err = tokens.float64()
	if err != nil {
		return err
	}

	sigProcID, err = tokens.int()
	if err != nil {
		return err
	}

	dec.sigProcBc, err = tokens.int()
	if err != nil {
		return err
	}

	*nVtx, err = tokens.int()
	if err != nil {
		return err
	}

	nRndm, err = tokens.int()
	if err != nil {
		return err

	}

	rndmStates := make([]int64, nRndm)
	for i := 0; i < nRndm; i++ {
		rndmStates[i], err = tokens.int64()
		if err != nil {
			return err
		}
	}

	nWeights, err = tokens.int()
	if err != nil {
		return err
	}

	weights := make([]float64, nWeights)
	for i := 0; i < nWeights; i++ {
		weights[i], err = tokens.float64()
		if err != nil {
			return err
		}
	}

	// fill infos gathered so far
	evt.SignalProcessID = sigProcID
	evt.EventNumber = evtNbr
	evt.Mpi = mpi
	if evt.Weights.Slice == nil {
		evt.Weights = NewWeights()
	}
	evt.Weights.Slice = weights
	evt.RandomStates = rndmStates
	evt.Scale = scale
	evt.AlphaQCD = aqcd
	evt.AlphaQED = aqed

	evt.Vertices = make(map[int]*Vertex, *nVtx)
	evt.Particles = make(map[int]*Particle, *nVtx*2)
	return err
}

func (dec *Decoder) decodeExtendedASCII(evt *Event, nVtx *int, tokens tokens) error {
	return dec.decodeEvent(evt, nVtx, tokens)
}
