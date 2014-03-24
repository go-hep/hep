package fads

import (
	"time"
)

type SimCaloHit struct {
	CellID int32
	Ene    float64
	Pos    [3]float64
}

type RawCaloHit struct {
	CellID int32
	Ampl   int32     // amplitude
	Time   time.Time // time stamp
}

type CaloHit struct {
	CellID int32
	Ene    float64 // energy of the hit
	EneErr float64 // error on the hit energy
	Time   float64 // time of the hit in [ns]
}

type Cluster struct {
	Type       byte         // flagword defining the type of the cluster
	Ene        float64      // energy of the cluster energy
	EneErr     float64      // error on the energy of the cluster
	Pos        [3]float64   // position of the cluster
	ErrPos     [6]float64   // covariance matrix of the position
	Theta      float64      // intrinsic directrion of cluster at position: theta
	Phi        float64      // intrinsic direction of cluster at position: phi
	ErrDir     [3]float64   // covariance matrix of the direction
	Shape      []float64    // shape parameters
	PIDs       []ParticleID // particle IDs sorted by their probability
	Clusters   []Cluster    // clusters combined to this cluster
	Hits       []CaloHit    // hits combined to form this cluster
	Weights    []float64
	SubDetEnes []float64 // subdetectors energies
}

// EOF
