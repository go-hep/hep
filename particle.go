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
