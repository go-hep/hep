// Copyright 2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package heppdt

import (
	"math"
)

type location int

//  PID digits (base 10) are: n Nr Nl Nq1 Nq2 Nq3 Nj
//  The location enum provides a convenient index into the PID.
const (
	_ location = iota
	Nj
	Nq3
	Nq2
	Nq1
	Nl
	Nr
	N
	N8
	N9
	N10
)

// Quarks describes a given quark mixture
type Quarks struct {
	Nq1 int16
	Nq2 int16
	Nq3 int16
}

// Particle Identification number
// In the standard numbering scheme, the PID digits (base 10) are:
//           +/- n Nr Nl Nq1 Nq2 Nq3 Nj
// It is expected that any 7 digit number used as a PID will adhere to
// the Monte Carlo numbering scheme documented by the PDG.
// Note that particles not already explicitly defined
// can be expressed within this numbering scheme.
type PID int

// ExtraBits returns everything beyoind the 7th digit
// (e.g. outside the numbering scheme)
func (pid PID) ExtraBits() int {
	return pid.AbsPID() / 10000000
}

// FundamentalID returns the first 2 digits if this is a "fundamental"
// particle.
// ID==100 is a special case (internal generator ID's are 81-100)
// Also, 101 and 102 are now used for geantinos
func (pid PID) FundamentalID() int {
	if pid.Digit(N10) == 1 && pid.Digit(N9) == 0 {
		return 0
	}

	if pid.Digit(Nq2) == 2 && pid.Digit(Nq1) == 0 {
		return pid.AbsPID() % 10000
	} else if pid.AbsPID() <= 102 {
		return pid.AbsPID()
	} else {
		return 0
	}
}

// Digit splits the PID into constituent integers
func (pid PID) Digit(loc location) int {
	//  PID digits (base 10) are: n Nr Nl Nq1 Nq2 Nq3 Nj
	//  the location enum provides a convenient index into the PID
	num := int(math.Pow(10.0, float64(loc-1)))
	return int(pid.AbsPID()/num) % 10
}

// findQ
func (pid PID) findQ(q int) bool {
	if pid.IsDyon() {
		return false
	}
	if pid.IsRhadron() {
		iz := 7
		for i := 6; i > 1; i-- {
			if pid.Digit(location(i)) == 0 {
				iz = i
			} else if i == iz-1 {
				// ignore squark or gluino
			} else {
				if pid.Digit(location(i)) == q {
					return true
				}
			}
		}
		return false
	}
	if pid.Digit(Nq3) == q || pid.Digit(Nq2) == q || pid.Digit(Nq1) == q {
		return true
	}
	if pid.IsPentaquark() {
		if pid.Digit(Nl) == q || pid.Digit(Nr) == q {
			return true
		}
	}
	return false
}

// AbsPID returns the absolute value of the particle ID
func (pid PID) AbsPID() int {
	id := int(pid)
	if id >= 0 {
		return id
	}
	return -int(id)
}

// IsValid returns whether PID is a valid particle ID
func (pid PID) IsValid() bool {
	if pid.ExtraBits() > 0 {
		switch {
		case pid.IsNucleus():
			return true
		case pid.IsQBall():
			return true
		default:
			return false
		}
	}

	switch {
	case pid.IsSUSY():
		return true
	case pid.IsRhadron():
		return true
	case pid.IsDyon():
		return true
	case pid.IsMeson():
		return true
	case pid.IsBaryon():
		return true
	case pid.IsDiQuark():
		return true
	case pid.FundamentalID() > 0:
		return true
	case pid.IsPentaquark():
		return true
	}

	return false
}

// IsMeson returns whether this is a valid meson ID
func (pid PID) IsMeson() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.AbsPID() <= 100:
		return false
	case pid.FundamentalID() <= 100 && pid.FundamentalID() > 0:
		return false
	}
	apid := pid.AbsPID()
	id := int(pid)

	switch {
	case apid == 130 || apid == 310 || apid == 210:
		return true

	case apid == 150 || apid == 350 || apid == 510 || apid == 530:
		// EvtGen odd number
		return true
	case id == 110 || id == 990 || id == 9990:
		// pomeron, etc...
		return true
	case pid.Digit(Nj) > 0 && pid.Digit(Nq3) > 0 && pid.Digit(Nq2) > 0 && pid.Digit(Nq1) == 0:
		// check for illegal antiparticles
		switch {
		case pid.Digit(Nq3) == pid.Digit(Nq2) && id < 0:
			return false
		default:
			return true
		}
	}
	return false
}

// IsBaryon returns whether this is a valid baryon id
func (pid PID) IsBaryon() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.AbsPID() <= 100:
		return false
	case pid.FundamentalID() <= 100 && pid.FundamentalID() > 0:
		return false
	case pid.AbsPID() == 2110 || pid.AbsPID() == 2210:
		return true
	case pid.Digit(Nj) > 0 && pid.Digit(Nq3) > 0 && pid.Digit(Nq2) > 0 && pid.Digit(Nq1) > 0:
		return true
	}
	return false
}

// IsDiQuark returns whether this is a valid diquark id
func (pid PID) IsDiQuark() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.AbsPID() <= 100:
		return false
	case pid.FundamentalID() <= 100 && pid.FundamentalID() > 0:
		return false
	case pid.Digit(Nj) > 0 && pid.Digit(Nq3) == 0 && pid.Digit(Nq2) > 0 && pid.Digit(Nq1) > 0:
		// EvtGen uses the diquarks for quark pairs, so for instance,
		// 5501 is a valid "diquark" for EvtGen
		// if pid.Digit(Nj) == 1 && pid.Digit(Nq2) == pid.Digit(Nq1) { 	// illegal
		//   return false
		// } else {
		return true
		// }
	}
	return false
}

// IsHadron returns whether this is a valid hadron id
func (pid PID) IsHadron() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.IsMeson():
		return true
	case pid.IsBaryon():
		return true
	case pid.IsPentaquark():
		return true
	}
	return false
}

// IsLepton returns whether this is a valid lepton id
func (pid PID) IsLepton() bool {
	if pid.ExtraBits() > 0 {
		return false
	}
	if fid := pid.FundamentalID(); fid >= 11 && fid <= 18 {
		return true
	}
	return false
}

// IsNucleus returns whether this is a valid nucleus id.
// This implements the 2006 Monte Carlon nuclear code scheme.
// Ion numbers are +/- 10LZZZAAAI.
// AAA is A - total baryon number
// ZZZ is Z - total charge
// L is the total number of strange quarks.
// I is the isomer number, with I=0 corresponding to the ground state.
func (pid PID) IsNucleus() bool {
	// a proton can also be a hydrogen nucleus
	if pid.AbsPID() == 2212 {
		return true
	}
	// new standard: +/- 10LZZZAAAI
	if pid.Digit(N10) == 1 && pid.Digit(N9) == 0 {
		// charge should always be less than or equal to baryon number
		if pid.A() >= pid.Z() {
			return true
		}
	}
	return false
}

// IsPentaquark returns whether this is a valid pentaquark id
func (pid PID) IsPentaquark() bool {
	// a pentaquark is of the form 9abcdej,
	// where j is the spin and a, b, c, d, and e are quarks
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.Digit(N) != 9:
		return false
	case pid.Digit(Nr) == 9 || pid.Digit(Nr) == 0:
		return false
	case pid.Digit(Nj) == 9 || pid.Digit(Nl) == 0:
		return false

	case pid.Digit(Nq1) == 0:
		return false
	case pid.Digit(Nq2) == 0:
		return false
	case pid.Digit(Nq3) == 0:
		return false
	case pid.Digit(Nj) == 0:
		return false

	case pid.Digit(Nq2) > pid.Digit(Nq1):
		return false
	case pid.Digit(Nq1) > pid.Digit(Nl):
		return false
	case pid.Digit(Nl) > pid.Digit(Nr):
		return false
	}

	return true
}

// IsSUSY returns whether this is a valid SUSY particle id
func (pid PID) IsSUSY() bool {
	// fundamental SUSY particles have n =1 or =2
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.Digit(N) != 1 && pid.Digit(N) != 2:
		return false
	case pid.Digit(Nr) != 0:
		return false
	case pid.FundamentalID() == 0:
		return false
	}

	return true
}

// IsRhadron returns whether this is a valid R-hadron particle id
func (pid PID) IsRhadron() bool {

	// an R-hadron is of the form 10abcdj,
	// where j is the spin and a, b, c, and d are quarks or gluons
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.Digit(N) != 1:
		return false
	case pid.Digit(Nr) != 0:
		return false
	case pid.IsSUSY():
		return false

	case pid.Digit(Nq2) == 0:
		return false // All R-hadrons have a least 3 core digits
	case pid.Digit(Nq3) == 0:
		return false // All R-hadrons have a least 3 core digits
	case pid.Digit(Nj) == 0:
		return false // All R-hadrons have a least 3 core digits
	}

	return true
}

// IsDyon returns whether this is a valid Dyon (magnetic monopole) id
func (pid PID) IsDyon() bool {
	// Magnetic monopoles and Dyons are assumed to have one unit of
	// Dirac monopole charge and a variable integer number xyz units
	// of electric charge.
	//
	// Codes 411xyz0 are then used when the magnetic and electrical
	// charge sign agree and 412xyz0 when they disagree,
	// with the overall sign of the particle set by the magnetic charge.
	// For now no spin information is provided.

	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.Digit(N) != 4:
		return false
	case pid.Digit(Nr) != 1:
		return false
	case pid.Digit(Nl) != 1 && pid.Digit(Nl) != 2:
		return false
	case pid.Digit(Nq3) == 0:
		return false // all Dyons have at least 1 core digit
	case pid.Digit(Nj) != 0:
		return false // dyons have spin zero for now
	}

	return true
}

// IsQBall checks for QBall or any exotic particle with electric charge
// beyond the qqq scheme.
// Ad-hoc numbering for such particles is 100xxxx0, where xxxx is the
// charge in tenths.
func (pid PID) IsQBall() bool {
	// Ad-hoc numbering for such particles is 100xxxx0,
	// where xxxx is the charge in tenths.
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.Digit(N) != 1 && pid.Digit(N) != 2:
		return false
	case pid.Digit(Nr) != 0:
		return false
	case pid.FundamentalID() == 0:
		return false
	}
	return true
}

// HasUp returns whether this particle contains an up quark
func (pid PID) HasUp() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.FundamentalID() > 0:
		return false
	}
	return pid.findQ(2)
}

// HasDown returns whether this particle contains a down quark
func (pid PID) HasDown() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.FundamentalID() > 0:
		return false
	}
	return pid.findQ(1)
}

// HasStrange returns whether this particle contains a strange quark
func (pid PID) HasStrange() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.FundamentalID() > 0:
		return false
	}
	return pid.findQ(3)
}

// HasCharm returns whether this particle contains a charm quark
func (pid PID) HasCharm() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.FundamentalID() > 0:
		return false
	}
	return pid.findQ(4)
}

// HasBottom returns whether this particle contains a bottom quark
func (pid PID) HasBottom() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.FundamentalID() > 0:
		return false
	}
	return pid.findQ(5)
}

// HasTop returns whether this particle contains a top quark
func (pid PID) HasTop() bool {
	switch {
	case pid.ExtraBits() > 0:
		return false
	case pid.FundamentalID() > 0:
		return false
	}
	return pid.findQ(6)
}

// A returns A if this is a nucleus
func (pid PID) A() int {
	// a proton can also be a hydrogen nucleus
	switch {
	case pid.AbsPID() == 2212:
		return 1
	case pid.Digit(N10) != 1 || pid.Digit(N9) != 0:
		return 0
	}
	return (pid.AbsPID() / 10) % 1000
}

// Z returns Z if this is a nucleus
func (pid PID) Z() int {
	// a proton can also be a hydrogen nucleus
	switch {
	case pid.AbsPID() == 2212:
		return 1
	case pid.Digit(N10) != 1 || pid.Digit(N9) != 0:
		return 0
	}
	return (pid.AbsPID() / 10000) % 1000
}

// Lambda returns lambda if this is a nucleus
func (pid PID) Lambda() int {

	// a proton can also be a hydrogen nucleus
	if pid.AbsPID() == 2212 {
		return 0
	}

	if !pid.IsNucleus() {
		return 0
	}

	return pid.Digit(N8)
}

// JSpin returns 2J+1, where J is the total spin
func (pid PID) JSpin() int {
	fid := pid.FundamentalID()
	if fid > 0 && fid <= 100 {
		switch {
		case fid > 0 && fid < 7:
			return 2

		case fid == 9:
			return 3

		case fid > 10 && fid < 17:
			return 2

		case fid > 20 && fid < 25:
			return 3
		}
		return 0
	} else if pid.ExtraBits() > 0 {
		return 0
	}
	return pid.AbsPID() % 10
}

// LSpin returns the orbital angular momentum.
// Valid for mesons only
func (pid PID) LSpin() int {
	if !pid.IsMeson() {
		return 0
	}

	tent := (pid.AbsPID() / 1000000) % 10
	if tent == 9 {
		return 0
	}

	Nl := (pid.AbsPID() / 10000) % 10
	js := pid.AbsPID() % 10

	if Nl == 0 && js == 3 {
		return 0
	} else if Nl == 0 && js == 5 {
		return 1
	} else if Nl == 0 && js == 7 {
		return 2
	} else if Nl == 0 && js == 9 {
		return 3
	} else if Nl == 0 && js == 1 {
		return 0
	} else if Nl == 1 && js == 3 {
		return 1
	} else if Nl == 1 && js == 5 {
		return 2
	} else if Nl == 1 && js == 7 {
		return 3
	} else if Nl == 1 && js == 9 {
		return 4
	} else if Nl == 2 && js == 3 {
		return 1
	} else if Nl == 2 && js == 5 {
		return 2
	} else if Nl == 2 && js == 7 {
		return 3
	} else if Nl == 2 && js == 9 {
		return 4
	} else if Nl == 1 && js == 1 {
		return 1
	} else if Nl == 3 && js == 3 {
		return 2
	} else if Nl == 3 && js == 5 {
		return 3
	} else if Nl == 3 && js == 7 {
		return 4
	} else if Nl == 3 && js == 9 {
		return 5
	}
	// default to zero
	return 0
}

// SSpin returns the spin. Valid for mesons only
func (pid PID) SSpin() int {
	if !pid.IsMeson() {
		return 0
	}

	tent := (pid.AbsPID() / 1000000) % 10
	if tent == 9 {
		return 0
	}

	Nl := (pid.AbsPID() / 10000) % 10
	js := pid.AbsPID() % 10

	if Nl == 0 && js >= 3 {
		return 1
	} else if Nl == 0 && js == 1 {
		return 0
	} else if Nl == 1 && js >= 3 {
		return 0
	} else if Nl == 2 && js >= 3 {
		return 1
	} else if Nl == 1 && js == 1 {
		return 1
	} else if Nl == 3 && js >= 3 {
		return 1
	}
	// default to zero
	return 0
}

var ch100 = [100]int{
	-1, 2, -1, 2, -1, 2, -1, 2, 0, 0,
	-3, 0, -3, 0, -3, 0, -3, 0, 0, 0,
	0, 0, 0, 3, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 3, 0, 0, 3, 0, 0, 0,
	0, -1, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 6, 3, 6, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
}

// threeCharge
func (pid PID) threeCharge() int {
	var charge int
	q1 := pid.Digit(Nq1)
	q2 := pid.Digit(Nq2)
	q3 := pid.Digit(Nq3)
	ida := pid.AbsPID()
	fid := pid.FundamentalID()

	if ida == 0 { // illegal
		return 0
	} else if pid.ExtraBits() > 0 {
		if pid.IsNucleus() { // ion
			return 3 * pid.Z()
		} else if pid.IsQBall() { // QBall
			charge = 3 * ((ida / 10) % 10000)
		} else { // not an ion
			return 0
		}
	} else if pid.IsDyon() { // Dyon
		charge = 3 * ((ida / 10) % 1000)
		// this is half right
		// the charge sign will be changed below if pid < 0
		if pid.Digit(Nl) == 2 {
			charge = -charge
		}
	} else if fid > 0 && fid <= 100 { // use table
		charge = ch100[fid-1]
		if ida == 1000017 || ida == 1000018 {
			charge = 0
		}
		if ida == 1000034 || ida == 1000052 {
			charge = 0
		}
		if ida == 1000053 || ida == 1000054 {
			charge = 0
		}
		if ida == 5100061 || ida == 5100062 {
			charge = 6
		}
	} else if pid.Digit(Nj) == 0 { // KL, Ks, or undefined
		return 0
	} else if (q1 == 0) || (pid.IsRhadron() && (q1 == 9)) { // meson			// mesons
		if q2 == 3 || q2 == 5 {
			charge = ch100[q3-1] - ch100[q2-1]
		} else {
			charge = ch100[q2-1] - ch100[q3-1]
		}
	} else if q3 == 0 { // diquarks
		charge = ch100[q2-1] + ch100[q1-1]
	} else if pid.IsBaryon() || (pid.IsRhadron() && (pid.Digit(Nl) == 9)) { // baryon 			// baryons
		charge = ch100[q3-1] + ch100[q2-1] + ch100[q1-1]
	}
	if charge == 0 {
		return 0
	} else if int(pid) < 0 {
		charge = -charge
	}
	return charge

}

const onethird = 1. / 3.0
const onethirtith = 1. / 30.0

// Charge returns the actual charge which might be fractional
func (pid PID) Charge() float64 {
	c := pid.threeCharge()
	if pid.IsQBall() {
		return float64(c) * onethirtith
	}
	return float64(c) * onethird
}

// Quarks returns a list of 3 constituent quarks
func (pid PID) Quarks() Quarks {
	return Quarks{
		Nq1: int16(pid.Digit(Nq1)),
		Nq2: int16(pid.Digit(Nq2)),
		Nq3: int16(pid.Digit(Nq3)),
	}
}
