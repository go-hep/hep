package hepevt

// Event is the Go representation of the FORTRAN-77 HEPEVT common block:
//
//   PARAMETER (NMXHEP=2000) 
//   COMMON/HEPEVT/NEVHEP,NHEP,ISTHEP(NMXHEP),IDHEP(NMXHEP), 
//   &       JMOHEP(2,NMXHEP),JDAHEP(2,NMXHEP),PHEP(5,NMXHEP),VHEP(4,NMXHEP)
type Event struct {
	Nevhep int          // event number (or some special meaning, see doc for details)
	Nhep   int          // actual number of entries in current event
	Isthep []int        // status code for n'th entry
	Idhep  []int        // particle identifier according to PDG
	Jmohep [][2]int     // index of 1st and 2nd mother
	Jdahep [][2]int     // index of 1st and 2nd daughter
	Phep   [][5]float64 // particle 5-vector (px,py,pz,e,m)
	Vhep   [][4]float64 // vertex 4-vector (x,y,z,t)
}

type Particle struct {
	Status    int32      // status code (see hepevt doc)
	Id        int32      // barcode
	Mothers   [2]int32   // indices of 1st and 2nd mothers
	Daughters [2]int32   // indices of 1st and 2nd mothers
	P         [5]float64 // (px,py,pz,e,m)
	V         [4]float64 // vertex position (x,y,z,t)
}

// EOF
