package fads

import (
	"math"
	"math/rand"
	"reflect"
	"sort"
	"sync"

	"github.com/go-hep/fwk"
	"github.com/go-hep/random"
)

type etaphiBin struct {
	eta float64
	phi float64
}

type EtaPhiBin struct {
	EtaBins []float64
	PhiBins []float64
}

type EtaPhiGrid struct {
	eta []float64
	phi map[float64][]float64
}

func NewEtaPhiGrid(bins []EtaPhiBin) EtaPhiGrid {

	neta := 0
	for _, bin := range bins {
		nn := len(bin.EtaBins)
		if nn > neta {
			neta = nn
		}
	}

	grid := EtaPhiGrid{
		eta: make([]float64, 0, neta),
		phi: make(map[float64][]float64, neta),
	}

	for _, bin := range bins {
		for _, eta := range bin.EtaBins {
			phibins, ok := grid.phi[eta]
			if !ok {
				phibins = make([]float64, 0, len(bin.PhiBins))
				grid.phi[eta] = phibins
				grid.eta = append(grid.eta, eta)
			}
			sort.Float64s(phibins)
			for _, phi := range bin.PhiBins {
				i := sort.SearchFloat64s(phibins, phi)
				if !(i < len(phibins) && phibins[i] == phi) {
					phibins = append(phibins, phi)
					sort.Float64s(phibins)
				}
			}
			grid.phi[eta] = phibins
		}
	}

	sort.Float64s(grid.eta)

	return grid
}

// EtaPhiIndex returns the eta/phi bin indices corresponding to a given eta/phi pair.
// EtaPhiIndex returns false if no bin contains this eta/phi pair.
func (grid *EtaPhiGrid) EtaPhiIndex(eta, phi float64) (int, int, bool) {
	// find eta bin
	etaidx := sort.SearchFloat64s(grid.eta, eta)
	if !(etaidx < len(grid.eta)) {
		return -etaidx, 0, false
	}
	etabin := grid.eta[etaidx]
	// special case of lowest-edge: test whether eta is inside acceptance
	if etaidx == 0 && etabin > eta {
		return -1, 0, false
	}

	phibins := grid.phi[etabin]

	// find phi bin
	phiidx := sort.SearchFloat64s(phibins, phi)
	if !(phiidx < len(phibins)) {
		return etaidx, -phiidx, false
	}
	//phibin := phibins[phiidx]

	return etaidx, phiidx, true
}

// EtaPhiBin returns the eta/phi bin center corresponding to a given eta/phi index pair.
func (grid *EtaPhiGrid) EtaPhiBin(ieta, iphi int) (float64, float64, bool) {
	var etac float64
	var phic float64

	if ieta > len(grid.eta) || ieta <= 0 {
		return etac, phic, false
	}
	eta := grid.eta[ieta]
	etac = 0.5 * (grid.eta[ieta-1] + eta)

	phibins := grid.phi[eta]
	if iphi > len(phibins) || iphi <= 0 {
		return etac, phic, false
	}

	phi := phibins[iphi]
	phic = 0.5 * (phibins[iphi-1] + phi)

	return etac, phic, true
}

type etwData struct {
	Ene        float64
	Time       float64
	WeightTime float64
}

func (etw *etwData) Add(ene, t float64) {
	etw.Ene += ene
	sqrt := math.Sqrt(ene)
	etw.Time += sqrt * t
	etw.WeightTime += sqrt
}

type caloTrack struct {
	ECal etwData
	HCal etwData
}

type caloTower struct {
	Eta   float64
	Phi   float64
	Edges [4]float64

	ECal etwData
	HCal etwData

	TrackHits  int
	PhotonHits int
}

type EneFrac struct {
	ECal float64
	HCal float64
}

type Calorimeter struct {
	fwk.TaskBase

	efrac map[int]EneFrac
	bins  EtaPhiGrid

	etaphibins []EtaPhiBin
	ecalres    func(eta, ene float64) float64
	hcalres    func(eta, ene float64) float64

	particles   string
	tracks      string
	towers      string
	photons     string
	eflowtracks string
	eflowtowers string

	seed int64
	src  rand.Source

	gauss random.Dist
	dmu   sync.Mutex
}

func (tsk *Calorimeter) Configure(ctx fwk.Context) error {
	var err error

	err = tsk.DeclInPort(tsk.particles, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclInPort(tsk.tracks, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.towers, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.photons, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.eflowtracks, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	err = tsk.DeclOutPort(tsk.eflowtowers, reflect.TypeOf([]Candidate{}))
	if err != nil {
		return err
	}

	tsk.src = rand.NewSource(tsk.seed)
	tsk.gauss = random.Gauss(0, 1, &tsk.src)
	return err
}

func (tsk *Calorimeter) StartTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Calorimeter) StopTask(ctx fwk.Context) error {
	var err error

	return err
}

func (tsk *Calorimeter) Process(ctx fwk.Context) error {
	var err error

	store := ctx.Store()
	msg := ctx.Msg()

	v, err := store.Get(tsk.particles)
	if err != nil {
		return err
	}

	parts := v.([]Candidate)
	msg.Debugf(">>> particles: %v\n", len(parts))

	v, err = store.Get(tsk.tracks)
	if err != nil {
		return err
	}
	tracks := v.([]Candidate)
	msg.Debugf(">>> tracks: %v\n", len(tracks))

	towers := make([]Candidate, 0, len(tracks))
	defer func() {
		err = store.Put(tsk.towers, towers)
	}()

	photons := make([]Candidate, 0, len(tracks))
	defer func() {
		err = store.Put(tsk.photons, photons)
	}()

	eflowtracks := make([]Candidate, 0, len(tracks))
	defer func() {
		err = store.Put(tsk.eflowtracks, eflowtracks)
	}()

	eflowtowers := make([]Candidate, 0, len(tracks))
	defer func() {
		err = store.Put(tsk.eflowtowers, eflowtowers)
	}()

	hits := make(map[int64][]int64)
	twrecal := make([]float64, 0, len(parts))
	twrhcal := make([]float64, 0, len(parts))
	trkecal := make([]float64, 0, len(tracks))
	trkhcal := make([]float64, 0, len(tracks))

	// process particles
	for i := range parts {
		part := &parts[i]
		// msg.Debugf("part[%d]=%#v\n", i, part.Pos)

		abspid := part.Pid
		if abspid < 0 {
			abspid = -abspid
		}

		frac, ok := tsk.efrac[int(abspid)]
		if !ok {
			frac = tsk.efrac[0]
		}

		twrecal = append(twrecal, frac.ECal)
		twrhcal = append(twrhcal, frac.HCal)

		if frac.ECal < 1e-9 && frac.HCal < 1e-9 {
			// msg.Debugf("part[%d]= not enough ECal|HCal\n", i)
			continue
		}

		// find eta/phi bin
		etabin, phibin, ok := tsk.bins.EtaPhiIndex(part.Pos.Eta(), part.Pos.Phi())
		if !ok {
			// msg.Debugf("part[%d]= not in acceptance\n", i)
			continue
		}

		flags := int64(0)
		if abspid == 11 || abspid == 22 {
			flags |= (int64(1)) << 1
		}

		// make tower hit:
		// {16-bits: eta bin-id} {16-bits: phi bin-id} {8-bits: flags}
		// {24-bits: particle number}
		hit := (int64(etabin) << 48) | (int64(phibin) << 32) | (int64(flags) << 24) | int64(i)
		towerid := hit >> 32
		hits[towerid] = append(hits[towerid], hit)
	}

	// process tracks
	for i := range tracks {

		track := &tracks[i]
		// msg.Debugf("track[%d]=%#v\n", i, track.Pos)
		abspid := track.Pid
		if abspid < 0 {
			abspid = -abspid
		}

		frac, ok := tsk.efrac[int(abspid)]
		if !ok {
			frac = tsk.efrac[0]
		}

		trkecal = append(trkecal, frac.ECal)
		trkhcal = append(trkhcal, frac.HCal)

		if frac.ECal < 1e-9 && frac.HCal < 1e-9 {
			// msg.Debugf("track[%d]= not enough ECal|HCal\n", i)
			continue
		}

		// find eta/phi bin
		etabin, phibin, ok := tsk.bins.EtaPhiIndex(track.Pos.Eta(), track.Pos.Phi())
		if !ok {
			// msg.Debugf("track[%d]= not in acceptance\n", i)
			continue
		}

		flags := int64(1)

		// make tower hit:
		// {16-bits: eta bin-id} {16-bits: phi bin-id} {8-bits: flags}
		// {24-bits: track number}
		hit := (int64(etabin) << 48) | (int64(phibin) << 32) | (int64(flags) << 24) | int64(i)
		towerid := hit >> 32
		hits[towerid] = append(hits[towerid], hit)
	}

	nhits := 0
	twrhits := make([]int64, 0, len(hits))
	for towerid, hits := range hits {
		// hits are sorted first by eta bin-id, then phi bin-id,
		// then flags and then by particle or track number
		sort.Sort(int64Slice(hits))
		twrhits = append(twrhits, towerid)
		nhits += len(hits)
	}

	// hits are sorted first by eta bin-id, then phi bin-id,
	// then flags and then by particle or track number
	sort.Sort(int64Slice(twrhits))

	msg.Debugf("tower-hits: %d (%d)\n", nhits, len(twrhits))

	// process hits
	for _, towerid := range twrhits {
		iphi := (towerid >> 00) & 0x000000000000FFFF
		ieta := (towerid >> 16) & 0x000000000000FFFF

		// get eta/phi of tower's center
		eta, phi, ok := tsk.bins.EtaPhiBin(int(ieta), int(iphi))
		if !ok {
			return fwk.Errorf("calorimeter: no valid eta/phi bin (ieta=%d iphi=%d)", ieta, iphi)
		}

		etabins := tsk.bins.eta
		phibins := tsk.bins.phi[etabins[ieta]]

		calotower := caloTower{
			Eta: eta,
			Phi: phi,
			Edges: [4]float64{
				etabins[ieta-1],
				etabins[ieta],
				phibins[iphi-1],
				phibins[iphi],
			},
		}

		calotrk := caloTrack{}

		var tower Candidate
		twrtrks := make([]Candidate, 0, len(tracks))

		for _, hit := range hits[towerid] {
			flags := (hit >> 24) & 0x00000000000000FF
			n := hit & 0x0000000000FFFFFF
			// etaphi := hit >> 32

			switch {
			case (flags & 1) != 0: // track hits
				calotower.TrackHits++
				track := &tracks[n]
				ene := track.Mom.E()
				t := track.Pos.T()
				calotrk.ECal.Add(ene*trkecal[n], t)
				calotrk.HCal.Add(ene*trkhcal[n], t)
				twrtrks = append(twrtrks, *track)

			default:
				if (flags & 2) != 0 { // photon hits
					calotower.PhotonHits++
				}

				part := &parts[n]
				ene := part.Mom.E()
				t := part.Pos.T()
				calotower.ECal.Add(ene*twrecal[n], t)
				calotower.HCal.Add(ene*twrhcal[n], t)
				tower.Add(part)
			}
			// msg.Debugf("hit=0x%x >> flags=0x%x, n=%d\n", hit, flags, n)
		}

		ecalSigma := tsk.ecalres(calotower.Eta, calotower.ECal.Ene)
		ecalEne := tsk.lognormal(calotower.ECal.Ene, ecalSigma)
		ecalTime := 0.0
		if calotower.ECal.WeightTime >= 1e-9 {
			ecalTime = calotower.ECal.Time / calotower.ECal.WeightTime
		}

		hcalSigma := tsk.hcalres(calotower.Eta, calotower.HCal.Ene)
		hcalEne := tsk.lognormal(calotower.HCal.Ene, hcalSigma)
		hcalTime := 0.0
		if calotower.HCal.WeightTime >= 1e-9 {
			hcalTime = calotower.HCal.Time / calotower.HCal.WeightTime
		}

		ene := ecalEne + hcalEne
		esqrt := math.Sqrt(ecalEne)
		hsqrt := math.Sqrt(hcalEne)
		time := (esqrt*ecalTime + hsqrt*hcalTime) / (esqrt + hsqrt)

		eta = rand.Float64()*(calotower.Edges[1]-calotower.Edges[0]) + calotower.Edges[0]
		phi = rand.Float64()*(calotower.Edges[3]-calotower.Edges[2]) + calotower.Edges[2]

		pt := ene / math.Cosh(eta)

		tower.Pos = newPtEtaPhiE(1, eta, phi, time)
		tower.Mom = newPtEtaPhiE(pt, eta, phi, ene)
		tower.Eem = ecalEne
		tower.Ehad = hcalEne

		tower.Edges = calotower.Edges

		if ene > 0 {
			if calotower.PhotonHits > 0 && calotower.TrackHits == 0 {
				photons = append(photons, tower)
			}
			towers = append(towers, tower)
		}
	}

	return err
}

func (tsk *Calorimeter) lognormal(mean, sigma float64) float64 {
	if mean <= 0 {
		return 0
	}

	b := math.Sqrt(math.Log(1 + (sigma*sigma)/(mean*mean)))
	a := math.Log(mean) - 0.5*b*b

	tsk.dmu.Lock()
	bgauss := tsk.gauss()
	tsk.dmu.Unlock()
	return math.Exp(a + bgauss)
}

func newCalorimeter(typ, name string, mgr fwk.App) (fwk.Component, error) {
	var err error

	tsk := &Calorimeter{
		TaskBase: fwk.NewTask(typ, name, mgr),

		bins:    NewEtaPhiGrid(nil),
		efrac:   make(map[int]EneFrac),
		ecalres: func(eta, ene float64) float64 { return 0 },
		hcalres: func(eta, ene float64) float64 { return 0 },

		particles:   "/fads/particles",
		tracks:      "/fads/tracks",
		towers:      "/fads/towers",
		photons:     "/fads/photons",
		eflowtracks: "/fads/eflowtracks",
		eflowtowers: "/fads/eflowtowers",

		seed: 1234,
	}

	// --

	err = tsk.DeclProp("EtaPhiBins", &tsk.bins)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EnergyFraction", &tsk.efrac)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("ECalResolution", &tsk.ecalres)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("HCalResolution", &tsk.hcalres)
	if err != nil {
		return nil, err
	}

	// --

	err = tsk.DeclProp("Particles", &tsk.particles)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Tracks", &tsk.tracks)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Towers", &tsk.towers)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Photons", &tsk.photons)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EFlowTracks", &tsk.eflowtracks)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("EFlowTowers", &tsk.eflowtowers)
	if err != nil {
		return nil, err
	}

	err = tsk.DeclProp("Seed", &tsk.seed)
	if err != nil {
		return nil, err
	}

	return tsk, err
}

func init() {
	fwk.Register(reflect.TypeOf(Calorimeter{}), newCalorimeter)
}
