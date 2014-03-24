package fads

type SimTrackerHit struct {
	CellID     int32
	Pos        [3]float64 // position of this hit
	Dep        float64    // energy deposit of the hit [GeV]
	Time       float64    // time of the hit in [ns]
	P          [3]float64 // 3-momentum of the particle at the hit position in [GeV]
	PathLength float64    // path length of the particle in the sensitive material that resulted in this hit.
	//Mc *McParticle

}

type RawTrackerData struct {
	CellID int32   // cell id
	ChanID int32   // channel id
	Time   int32   // time
	ADC    []int16 // measured ADC values
}

type TrackerData struct {
	CellID int32     // cell id
	Time   float64   // time
	Charge []float64 // calibrated ADC values
}

type TrackerPulse struct {
	CellID   int32        // cell id
	Time     float64      // time of the pulse
	Charge   float64      // integrated charge of the pulse
	Quality  int32        // quality bit flag of the pulse
	Cov      []float64    // covariance matrix of the charge and time measurements. Stored as lower triangle matrix.
	CorrData *TrackerData // tracker data used to create the pulse.
}

type Track struct {
	Type       int32
	Chi2       float32
	DEdx       float32
	DEdxErr    float32
	Radius     float32 // radius of inner-most hit
	SubDetHits []int32
	Tracks     []Track
	Hits       []TrackerHit
	States     []*TrackState
}

type TrackerHit struct {
}

type TrackState struct {
}

// EOF
