package hepmc

import (
	"bufio"
	//"bytes"
	"fmt"
	"io"
	"sort"
)

// Decoder decodes a hepmc Event from a stream
type Decoder struct {
	r *bufio.Reader
	//bbuf io.Reader
	//tbuf *bytes.Buffer
	seenEvtHdr bool
	ftype      hepmcFileType

	sigProcBc int // barcode of signal vertex
	bp1       int // barcode of beam1
	bp2       int // barcode of beam2
}

// NewDecoder returns a new hepmc Decoder that reads from the io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	//tbuf := bytes.NewBuffer(nil)
	//tr := io.TeeReader(r, tbuf)
	if rr, ok := r.(*bufio.Reader); ok {
		return &Decoder{r: rr}
	}
	//return &Decoder{r: bufio.NewReader(tr), tbuf: tbuf}
	return &Decoder{r: bufio.NewReader(r)}
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
	peek, err := dec.r.Peek(1)
	if err != nil {
		return err
	}

	if peek[0] != 'E' {
		err = dec.findEndKey()
		if err != nil {
			err = fmt.Errorf("hepmc.decode: invalid file (expected 'E' got '%v')", string(peek[0]))
			return err
		}
		// are we at the end of the file ?
		_, err = dec.r.Peek(1)
		if err != nil {
			return err
		}
		err = fmt.Errorf("hepmc.decode: end key not found")
		return err
	}

	nVtx := 0
	readingEvtHdr := true
	for readingEvtHdr {
		peek, err = dec.r.Peek(1)
		if err != nil {
			return err
		}
		//fmt.Printf("--> '%v'...\n", string(peek[0]))
		switch peek[0] {
		case 'E':
			// call appropriate decoder method
			switch dec.ftype {
			case hepmcGenEvent:
				err = dec.decodeEvent(evt, &nVtx)
			case hepmcASCII:
				err = dec.decodeASCII(evt, &nVtx)
			case hepmcExtendedASCII:
				err = dec.decodeExtendedASCII(evt, &nVtx)
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
			nWeights := 0
			_, err = fmt.Fscanf(dec.r, "N %d", &nWeights)
			if err != nil {
				return err
			}
			names := make(map[string]int, nWeights)
			for i := 0; i < nWeights; i++ {
				nn := ""
				_, err = fmt.Fscanf(dec.r, " %q", &nn)
				if err != nil {
					return err
				}
				names[nn] = i
			}

			_, err = dec.r.ReadString('\n')
			if err != nil {
				return err
			}
			//fmt.Printf("== weights: %v\n", names)
			evt.Weights.Map = names

		case 'U':
			if dec.ftype == hepmcGenEvent {
				err = dec.decodeUnits(evt)
				if err != nil {
					return err
				}
			}
		case 'C':
			err = dec.decodeCrossSection(evt)
			if err != nil {
				return err
			}
		case 'H':
			switch dec.ftype {
			case hepmcGenEvent, hepmcExtendedASCII:
				var hi HeavyIon
				err = dec.decodeHeavyIon(&hi)
				if err != nil {
					return err
				}
				evt.HeavyIon = &hi
			}
		case 'F':
			switch dec.ftype {
			case hepmcGenEvent, hepmcExtendedASCII:
				var pdf PdfInfo
				err = dec.decodePdfInfo(&pdf)
				if err != nil {
					return err
				}
				evt.PdfInfo = &pdf
			}
		case 'V', 'P':
			readingEvtHdr = false

		default:

			return fmt.Errorf(
				"hepmc.decoder: invalid file (got '%v')",
				peek[0],
			)
		}
	}
	// dec.r.Read(peek[0:])

	// end vertices of particles are not connected until
	// after the full event is read
	//  => store the values in a map until then
	pidxToEndVtx := make(map[int]int, nVtx) // particle-idx to end_vtx barcode

	// decode the vertices
	for i := 0; i < nVtx; i++ {
		vtx := &Vertex{}
		vtx.Event = evt
		err = dec.decodeVertex(evt, vtx, pidxToEndVtx)
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
		line, err := dec.r.ReadString('\n')
		if err != nil {
			return err
		}
		if line == "" {
			continue
		}
		line = line[:len(line)-1]
		//fmt.Printf("--> %q\n", line)
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

	err := fmt.Errorf("hepmc.ascii: invalid input file")
	return err

}

func (dec *Decoder) findEndKey() error {

	peek, err := dec.r.Peek(1)
	if peek[0] != 'H' {
		err = fmt.Errorf("hepmc.decode: not an end-key (%v)", string(peek))
		return err
	}
	line, err := dec.r.ReadString('\n')
	if err != nil {
		return err
	}
	if line != "" {
		line = line[:len(line)-1]
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

func (dec *Decoder) decodeUnits(evt *Event) error {
	var err error
	momUnit := ""
	lenUnit := ""
	_, err = fmt.Fscanf(dec.r, "U %s %s\n", &momUnit, &lenUnit)
	if err != nil {
		return err
	}
	evt.MomentumUnit, err = MomentumUnitFromString(momUnit)
	if err != nil {
		return err
	}
	evt.LengthUnit, err = LengthUnitFromString(lenUnit)
	if err != nil {
		return err
	}
	return err
}

func (dec *Decoder) decodeCrossSection(evt *Event) error {
	var err error
	var x CrossSection
	_, err = fmt.Fscanf(dec.r, "C %e %e\n", &x.Value, &x.Error)
	if err != nil {
		return err
	}
	evt.CrossSection = &x
	return err
}

func (dec *Decoder) decodeEvent(evt *Event, nVtx *int) error {
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

	_, err = fmt.Fscanf(
		dec.r,
		"E %d %d %e %e %e %d %d %d %d %d %d",
		&evtNbr,
		&mpi,
		&scale,
		&aqcd,
		&aqed,
		&sigProcID,
		&dec.sigProcBc,
		nVtx,
		&dec.bp1,
		&dec.bp2,
		&nRndm,
	)
	if err != nil {
		return err
	}
	rndmStates := make([]int64, nRndm)
	for i := 0; i < nRndm; i++ {
		_, err = fmt.Fscanf(dec.r, " %d", &rndmStates[i])
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fscanf(dec.r, " %d", &nWeights)
	if err != nil {
		return err
	}
	weights := make([]float64, nWeights)
	for i := 0; i < nWeights; i++ {
		_, err = fmt.Fscanf(dec.r, " %e", &weights[i])
		if err != nil {
			return err
		}
	}

	_, err = dec.r.ReadString('\n')
	if err != nil {
		return err
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

func (dec *Decoder) decodeVertex(evt *Event, vtx *Vertex, pidxToEndVtx map[int]int) error {

	var err error
	peek, err := dec.r.Peek(1)
	if err != nil {
		return err
	}
	if peek[0] != 'V' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'V', got '%v')",
			peek[0],
		)
	}

	orphans := 0
	nPartsOut := 0
	nWeights := 0

	_, err = fmt.Fscanf(
		dec.r,
		"V %d %d %e %e %e %e %d %d %d",
		&vtx.Barcode,
		&vtx.ID,
		&vtx.Position[0], &vtx.Position[1], &vtx.Position[2], &vtx.Position[3],
		&orphans,
		&nPartsOut,
		&nWeights,
	)
	if err != nil {
		return err
	}
	// FIXME: reuse buffers ?
	vtx.Weights.Slice = make([]float64, nWeights)
	for i := 0; i < nWeights; i++ {
		_, err = fmt.Fscanf(dec.r, " %e", &vtx.Weights.Slice[i])
		if err != nil {
			return err
		}
	}
	_, err = dec.r.ReadString('\n')
	if err != nil {
		return err
	}

	// read and create the associated particles
	// outgoing particles are added to their production vertices immediately.
	// incoming particles are added to a map and handled later.
	for i := 0; i < orphans; i++ {
		p := &Particle{}
		err = dec.decodeParticle(evt, p, pidxToEndVtx)
		if err != nil {
			return err
		}
		evt.Particles[p.Barcode] = p
	}
	// FIXME: reuse buffers ?
	vtx.ParticlesOut = make([]*Particle, nPartsOut)
	for i := 0; i < nPartsOut; i++ {
		p := &Particle{ProdVertex: vtx}
		err = dec.decodeParticle(evt, p, pidxToEndVtx)
		if err != nil {
			return err
		}
		evt.Particles[p.Barcode] = p
		vtx.ParticlesOut[i] = p
	}
	return err
}

func (dec *Decoder) decodeParticle(evt *Event, p *Particle, pidxToEndVtx map[int]int) error {

	var err error
	peek, err := dec.r.Peek(1)
	if err != nil {
		return err
	}
	if peek[0] != 'P' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'P', got '%v')",
			peek[0],
		)
	}

	endBc := 0
	_, err = fmt.Fscanf(
		dec.r,
		"P %d %d %e %e %e %e",
		&p.Barcode,
		&p.PdgID,
		&p.Momentum[0], &p.Momentum[1], &p.Momentum[2], &p.Momentum[3],
	)
	if err != nil {
		return err
	}
	if dec.ftype != hepmcASCII {
		_, err = fmt.Fscanf(
			dec.r, " %e", &p.GeneratedMass)
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fscanf(
		dec.r,
		"%d %e %e %d",
		&p.Status,
		&p.Polarization.Theta,
		&p.Polarization.Phi,
		&endBc,
	)
	if err != nil {
		return nil
	}

	err = dec.decodeFlow(&p.Flow)
	if err != nil {
		return err
	}
	p.Flow.Particle = p

	//fmt.Printf(">>> flow-sz: %d == %v\n", len(p.Flow.Icode), p.Flow.Icode)
	_, err = dec.r.ReadString('\n')
	if err != nil {
		return err
	}

	// all particles are connected to their end vertex separately
	// after all particles and vertices have been created
	if endBc != 0 {
		pidxToEndVtx[p.Barcode] = endBc
	}
	return err
}

func (dec *Decoder) decodeFlow(flow *Flow) error {
	var err error
	nFlow := 0
	_, err = fmt.Fscanf(dec.r, "%d", &nFlow)
	if err != nil {
		return err
	}
	//fmt.Printf("flow-sz: %d\n", n_flow)
	flow.Icode = make(map[int]int, nFlow)
	for i := 0; i < nFlow; i++ {
		k := 0
		v := 0
		_, err = fmt.Fscanf(dec.r, " %d %d", &k, &v)
		if err != nil {
			return err
		}
		flow.Icode[k] = v
		//fmt.Printf("  %d -> %d\n", k, flow.Icode[k])
	}
	return err
}

func (dec *Decoder) decodeHeavyIon(hi *HeavyIon) error {
	var err error
	peek, err := dec.r.Peek(1)
	if err != nil {
		return err
	}
	if peek[0] != 'H' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'H', got '%v')",
			peek[0],
		)
	}
	_, err = fmt.Fscanf(
		dec.r,
		"H %d %d %d %d %d %d %d %d %d %e %e %e %e\n",
		&hi.NCollHard,
		&hi.NPartProj,
		&hi.NPartTarg,
		&hi.NColl,
		&hi.NNwColl,
		&hi.NwNColl,
		&hi.NwNwColl,
		&hi.SpectatorNeutrons,
		&hi.SpectatorProtons,
		&hi.ImpactParameter,
		&hi.EventPlaneAngle,
		&hi.Eccentricity,
		&hi.SigmaInelNN,
	)
	return err
}

func (dec *Decoder) decodePdfInfo(pdf *PdfInfo) error {
	var err error
	peek, err := dec.r.Peek(1)
	if err != nil {
		return err
	}
	if peek[0] != 'F' {
		return fmt.Errorf(
			"hepmc.decode: invalid file (expected 'F', got '%v')",
			peek[0],
		)
	}
	_, err = fmt.Fscanf(
		dec.r,
		"F %d %d %e %e %e %e %e %d %d\n",
		&pdf.ID1,
		&pdf.ID2,
		&pdf.X1,
		&pdf.X2,
		&pdf.ScalePDF,
		&pdf.Pdf1,
		&pdf.Pdf2,
		&pdf.LHAPdf1,
		&pdf.LHAPdf2,
	)
	return err
}

func (dec *Decoder) decodeASCII(evt *Event, nVtx *int) error {
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

	_, err = fmt.Fscanf(
		dec.r,
		"E %d %e %e %e %d %d %d %d",
		&evtNbr,
		&scale,
		&aqcd,
		&aqed,
		&sigProcID,
		&dec.sigProcBc,
		nVtx,
		&nRndm,
	)
	if err != nil {
		return err
	}
	rndmStates := make([]int64, nRndm)
	for i := 0; i < nRndm; i++ {
		_, err = fmt.Fscanf(dec.r, " %d", &rndmStates[i])
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fscanf(dec.r, " %d", &nWeights)
	if err != nil {
		return err
	}
	weights := make([]float64, nWeights)
	for i := 0; i < nWeights; i++ {
		_, err = fmt.Fscanf(dec.r, " %e", &weights[i])
		if err != nil {
			return err
		}
	}

	_, err = dec.r.ReadString('\n')
	if err != nil {
		return err
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

func (dec *Decoder) decodeExtendedASCII(evt *Event, nVtx *int) error {
	return dec.decodeEvent(evt, nVtx)
}

// EOF
