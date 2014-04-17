package heppdt

type location int

//  PID digits (base 10) are: n nr nl nq1 nq2 nq3 nj
//  The location enum provides a convenient index into the PID.
const (
	_ location = iota
	nj
	nq3
	nq2
	nq1
	nl
	nr
	n
	n8
	n9
	n10
)

/// constituent quarks
type Quarks struct {
	Nq1 int16
	Nq2 int16
	Nq3 int16
}

// Particle Identification number
// In the standard numbering scheme, the PID digits (base 10) are:
//           +/- n nr nl nq1 nq2 nq3 nj
// It is expected that any 7 digit number used as a PID will adhere to
// the Monte Carlo numbering scheme documented by the PDG.
// Note that particles not already explicitly defined
// can be expressed within this numbering scheme.
type PID int

/*
// Check for QBall or any exotic particle with electric charge beyond the qqq scheme
bool ParticleID::isQBall( ) const
{
    // Ad-hoc numbering for such particles is 100xxxx0,
    // where xxxx is the charge in tenths.
    if( extraBits() != 1 ) { return false; }
    if( digit(n) != 0 )  { return false; }
    if( digit(nr) != 0 )  { return false; }
    // check the core number
    if( (abspid()/10)%10000 == 0 )  { return false; }
    // these particles have spin zero for now
    if( digit(nj) != 0 )  { return false; }
    return true;
}
*/

// IsQBall checks for QBall or any exotic particle with electric charge
// beyond the qqq scheme.
// Ad-hoc numbering for such particles is 100xxxx0, where xxxx is the
// charge in tenths.
func (pid PID) IsQBall() bool {
	return false
}
