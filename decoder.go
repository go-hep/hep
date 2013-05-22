package hepmc

import (
	"bufio"
	//"bytes"
	"fmt"
	"io"
)

type Decoder struct {
	r *bufio.Reader
	//bbuf io.Reader
	//tbuf *bytes.Buffer
	seen_evt_hdr bool
	ftype        hepmc_ftype

	sig_proc_bc int // barcode of signal vertex
	bp1         int // barcode of beam1
	bp2         int // barcode of beam2
}

func NewDecoder(r io.Reader) *Decoder {
	//tbuf := bytes.NewBuffer(nil)
	//tr := io.TeeReader(r, tbuf)
	if rr, ok := r.(*bufio.Reader); ok {
		return &Decoder{r: rr}
	}
	//return &Decoder{r: bufio.NewReader(tr), tbuf: tbuf}
	return &Decoder{r: bufio.NewReader(r)}
}

func (dec *Decoder) Decode(evt *Event) error {
	var err error

	// search for event listing key
	if !dec.seen_evt_hdr {
		err = dec.find_file_type()
		if err != nil {
			return err
		}
		dec.seen_evt_hdr = true
	}

	dec.sig_proc_bc = 0
	dec.bp1 = 0
	dec.bp2 = 0

	// test to be sure the next entry is of type 'E'
	// FIXME: implement more of the logic from HepMC::IO_GenEvent
	peek, err := dec.r.Peek(1)
	if err != nil {
		return err
	}

	if peek[0] != 'E' {
		err = dec.find_end_key()
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

	n_vtx := 0
	reading_evt_hdr := true
	for reading_evt_hdr {
		peek, err = dec.r.Peek(1)
		if err != nil {
			return err
		}
		//fmt.Printf("--> '%v'...\n", string(peek[0]))
		switch peek[0] {
		case 'E':
			// call appropriate decoder method
			switch dec.ftype {
			case hepmc_genevent:
				err = dec.decode_genevent(evt, &n_vtx)
			case hepmc_ascii:
				err = dec.decode_ascii(evt, &n_vtx)
			case hepmc_extendedascii:
				err = dec.decode_extendedascii(evt, &n_vtx)
			case hepmc_ascii_pdt:
				// nop
			case hepmc_extendedascii_pdt:
				// nop
			default:
				panic("unreachable")
			}
			if err != nil {
				return err
			}
		case 'N':
			n_weights := 0
			_, err = fmt.Fscanf(dec.r, "N %d", &n_weights)
			if err != nil {
				return err
			}
			names := make(map[string]int, n_weights)
			for i := 0; i < n_weights; i++ {
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
			if dec.ftype == hepmc_genevent {
				err = dec.decode_units(evt)
				if err != nil {
					return err
				}
			}
		case 'C':
			err = dec.decode_cross_section(evt)
			if err != nil {
				return err
			}
		case 'H':
			switch dec.ftype {
			case hepmc_genevent, hepmc_extendedascii:
				var hi HeavyIon
				err = dec.decode_heavy_ion(&hi)
				if err != nil {
					return err
				}
				evt.HeavyIon = &hi
			}
		case 'F':
			switch dec.ftype {
			case hepmc_genevent, hepmc_extendedascii:
				var pdf PdfInfo
				err = dec.decode_pdf_info(&pdf)
				if err != nil {
					return err
				}
				evt.PdfInfo = &pdf
			}
		case 'V', 'P':
			reading_evt_hdr = false

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
	pidx_to_end_vtx := make(map[int]int, n_vtx) // particle-idx to end_vtx barcode

	// decode the vertices
	for i := 0; i < n_vtx; i++ {
		vtx := &Vertex{}
		vtx.Event = evt
		err = dec.decode_vertex(evt, vtx, pidx_to_end_vtx)
		if err != nil {
			return err
		}
		evt.Vertices[vtx.Barcode] = vtx
	}

	// set the signal process vertex
	if dec.sig_proc_bc != 0 {
		for _, vtx := range evt.Vertices {
			if vtx.Barcode == dec.sig_proc_bc {
				evt.SignalVertex = vtx
				break
			}
		}
		if evt.SignalVertex == nil {
			return fmt.Errorf("hepmc.decode: could not find signal vertex (barcode=%d)", dec.sig_proc_bc)
		}
	}

	// connect particles to their end vertices
	for i, end_vtx_bc := range pidx_to_end_vtx {
		p := evt.Particles[i]
		vtx := evt.Vertices[end_vtx_bc]
		vtx.ParticlesIn = append(vtx.ParticlesIn, p)
		p.EndVertex = vtx
		// also look for the beam particles
		if p.Barcode == dec.bp1 {
			evt.Beams[0] = p
		}
		if p.Barcode == dec.bp2 {
			evt.Beams[1] = p
		}
	}
	return err
}

func (dec *Decoder) find_file_type() error {

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

		case genevent_start:
			dec.ftype = hepmc_genevent
			return nil

		case ascii_start:
			dec.ftype = hepmc_ascii
			return nil

		case extendedascii_start:
			dec.ftype = hepmc_extendedascii
			return nil

		case pdt_start:
			dec.ftype = hepmc_ascii_pdt
			return nil

		case extendedascii_pdt_start:
			dec.ftype = hepmc_extendedascii_pdt
			return nil
		}
	}

	err := fmt.Errorf("hepmc.ascii: invalid input file")
	return err

}

func (dec *Decoder) find_end_key() error {

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

	var ftype hepmc_ftype
	switch line {
	default:
		err = fmt.Errorf("hepmc.decode: invalid file type (value=%q)", line)
		return err

	case genevent_end:
		ftype = hepmc_genevent

	case ascii_end:
		ftype = hepmc_ascii

	case extendedascii_end:
		ftype = hepmc_extendedascii

	case pdt_end:
		ftype = hepmc_ascii_pdt

	case extendedascii_pdt_end:
		ftype = hepmc_extendedascii_pdt
	}

	if ftype != dec.ftype {
		err = fmt.Errorf(
			"hepmc.decode: file type changed from %v to %v",
			dec.ftype, ftype,
		)
	}
	return err
}

func (dec *Decoder) decode_units(evt *Event) error {
	var err error
	mom_unit := ""
	len_unit := ""
	_, err = fmt.Fscanf(dec.r, "U %s %s\n", &mom_unit, &len_unit)
	if err != nil {
		return err
	}
	evt.MomentumUnit, err = MomentumUnitFromString(mom_unit)
	if err != nil {
		return err
	}
	evt.LengthUnit, err = LengthUnitFromString(len_unit)
	if err != nil {
		return err
	}
	return err
}

func (dec *Decoder) decode_cross_section(evt *Event) error {
	var err error
	var x CrossSection
	_, err = fmt.Fscanf(dec.r, "C %e %e\n", &x.Value, &x.Error)
	if err != nil {
		return err
	}
	evt.CrossSection = &x
	return err
}

func (dec *Decoder) decode_genevent(evt *Event, n_vtx *int) error {
	var (
		err         error
		evt_nbr     int
		mpi         int
		scale       float64
		a_qcd       float64
		a_qed       float64
		sig_proc_id int
		n_rndm      int
		n_weights   int
	)

	_, err = fmt.Fscanf(
		dec.r,
		"E %d %d %e %e %e %d %d %d %d %d %d",
		&evt_nbr,
		&mpi,
		&scale,
		&a_qcd,
		&a_qed,
		&sig_proc_id,
		&dec.sig_proc_bc,
		n_vtx,
		&dec.bp1,
		&dec.bp2,
		&n_rndm,
	)
	if err != nil {
		return err
	}
	rndm_states := make([]int64, n_rndm)
	for i := 0; i < n_rndm; i++ {
		_, err = fmt.Fscanf(dec.r, " %d", &rndm_states[i])
		if err != nil {
			return err
		}
	}

	_, err = fmt.Fscanf(dec.r, " %d", &n_weights)
	if err != nil {
		return err
	}
	weights := make([]float64, n_weights)
	for i := 0; i < n_weights; i++ {
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
	evt.SignalProcessId = sig_proc_id
	evt.EventNumber = evt_nbr
	evt.Mpi = mpi
	if evt.Weights.Slice == nil {
		evt.Weights = NewWeights()
	}
	evt.Weights.Slice = weights
	evt.RandomStates = rndm_states
	evt.Scale = scale
	evt.AlphaQCD = a_qcd
	evt.AlphaQED = a_qed

	evt.Vertices = make(map[int]*Vertex, *n_vtx)
	evt.Particles = make(map[int]*Particle, *n_vtx*2)
	return err
}

func (dec *Decoder) decode_vertex(
	evt *Event,
	vtx *Vertex,
	pidx_to_end_vtx map[int]int) error {

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
	n_parts_out := 0
	n_weights := 0

	_, err = fmt.Fscanf(
		dec.r,
		"V %d %d %e %e %e %e %d %d %d",
		&vtx.Barcode,
		&vtx.Id,
		&vtx.Position[0], &vtx.Position[1], &vtx.Position[2], &vtx.Position[3],
		&orphans,
		&n_parts_out,
		&n_weights,
	)
	if err != nil {
		return err
	}
	// FIXME: reuse buffers ?
	vtx.Weights.Slice = make([]float64, n_weights)
	for i := 0; i < n_weights; i++ {
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
		err = dec.decode_particle(evt, p, pidx_to_end_vtx)
		if err != nil {
			return err
		}
		evt.Particles[p.Barcode] = p
	}
	// FIXME: reuse buffers ?
	vtx.ParticlesOut = make([]*Particle, n_parts_out)
	for i := 0; i < n_parts_out; i++ {
		p := &Particle{ProdVertex: vtx}
		err = dec.decode_particle(evt, p, pidx_to_end_vtx)
		if err != nil {
			return err
		}
		evt.Particles[p.Barcode] = p
		vtx.ParticlesOut[i] = p
	}
	return err
}

func (dec *Decoder) decode_particle(
	evt *Event,
	p *Particle,
	pidx_to_end_vtx map[int]int) error {

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

	end_bc := 0
	_, err = fmt.Fscanf(
		dec.r,
		"P %d %d %e %e %e %e",
		&p.Barcode,
		&p.PdgId,
		&p.Momentum[0], &p.Momentum[1], &p.Momentum[2], &p.Momentum[3],
	)
	if err != nil {
		return err
	}
	if dec.ftype != hepmc_ascii {
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
		&end_bc,
	)
	if err != nil {
		return nil
	}

	err = dec.decode_flow(&p.Flow)
	if err != nil {
		return err
	}
	//fmt.Printf(">>> flow-sz: %d == %v\n", len(p.Flow.Icode), p.Flow.Icode)
	_, err = dec.r.ReadString('\n')
	if err != nil {
		return err
	}

	// all particles are connected to their end vertex separately
	// after all particles and vertices have been created
	if end_bc != 0 {
		pidx_to_end_vtx[p.Barcode] = end_bc
	}
	return err
}

func (dec *Decoder) decode_flow(flow *Flow) error {
	var err error
	n_flow := 0
	_, err = fmt.Fscanf(dec.r, "%d", &n_flow)
	if err != nil {
		return err
	}
	//fmt.Printf("flow-sz: %d\n", n_flow)
	flow.Icode = make(map[int]int, n_flow)
	for i := 0; i < n_flow; i++ {
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

func (dec *Decoder) decode_heavy_ion(hi *HeavyIon) error {
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
		&hi.Ncoll_hard,
		&hi.Npart_proj,
		&hi.Npart_targ,
		&hi.Ncoll,
		&hi.N_Nwounded_collisions,
		&hi.Nwounded_N_collisions,
		&hi.Nwounded_Nwounded_collisions,
		&hi.Spectator_neutrons,
		&hi.Spectator_protons,
		&hi.Impact_parameter,
		&hi.Event_plane_angle,
		&hi.Eccentricity,
		&hi.Sigma_inel_NN,
	)
	return err
}

func (dec *Decoder) decode_pdf_info(pdf *PdfInfo) error {
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
		&pdf.Id1,
		&pdf.Id2,
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

func (dec *Decoder) decode_ascii(evt *Event, n_vtx *int) error {
	var err error
	return err
}

func (dec *Decoder) decode_extendedascii(evt *Event, n_vtx *int) error {
	var err error
	return err
}

// EOF
