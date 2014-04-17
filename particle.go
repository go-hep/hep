package heppdt

// Particle holds informations on a particle as per the PDG booklet
type Particle struct {
	ID          PID           // particle ID
	Name        string        // particle name
	PDG         int           // PDG code of the particle
	Mass        float64       // particle mass in GeV
	Charge      float64       // electrical charge
	ColorCharge float64       // color charge
	Spin        SpinState     // spin state
	Quarks      []Constituent // constituents
	Resonance   Resonance     // resonance
}

// IsStable returns whether this particle is stable
func (p *Particle) IsStable() bool {
	res := &p.Resonance
	if res.Width.Value == -1. {
		return false
	}
	if res.Width.Value > 0 || res.Lifetime().Value > 0 {
		return false
	}
	return true
}
