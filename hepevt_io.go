package hepevt

/*
 #include <stdlib.h>
*/
import "C"

import (
	"fmt"
)

//const nmxhep = 4000

// I/O ----------------------------------------------------------------------

//////////////////////////////////////////////////////////////////////////
//
//      PARAMETER (NMXHEP=2000) 
//      COMMON/HEPEVT/NEVHEP,NHEP,ISTHEP(NMXHEP),IDHEP(NMXHEP), 
//     &        JMOHEP(2,NMXHEP),JDAHEP(2,NMXHEP),PHEP(5,NMXHEP),VHEP(4,NMXHEP)
/**********************************************************/
/*           D E S C R I P T I O N :                      */
/*--------------------------------------------------------*/
/* NEVHEP          - event number (or some special meaning*/
/*                    (see documentation for details)     */
/* NHEP            - actual number of entries in current  */
/*                    event.                              */
/* ISTHEP[IHEP]    - status code for IHEP'th entry - see  */
/*                    documentation for details           */
/* IDHEP [IHEP]    - IHEP'th particle identifier according*/
/*                    to PDG.                             */
/* JMOHEP[IHEP][0] - pointer to position of 1st mother    */
/* JMOHEP[IHEP][1] - pointer to position of 2nd mother    */
/* JDAHEP[IHEP][0] - pointer to position of 1st daughter  */
/* JDAHEP[IHEP][1] - pointer to position of 2nd daughter  */
/* PHEP  [IHEP][0] - X momentum                           */
/* PHEP  [IHEP][1] - Y momentum                           */
/* PHEP  [IHEP][2] - Z momentum                           */
/* PHEP  [IHEP][3] - Energy                               */
/* PHEP  [IHEP][4] - Mass                                 */
/* VHEP  [IHEP][0] - X vertex                             */
/* VHEP  [IHEP][1] - Y vertex                             */
/* VHEP  [IHEP][2] - Z vertex                             */
/* VHEP  [IHEP][3] - production time                      */
/*========================================================*/
type Event struct {
	Nevhep int
	Nhep   int
	Isthep []int
	Idhep  []int
	Jmohep [][2]int
	Jdahep [][2]int
	Phep   [][5]float64
	Vhep   [][4]float64
}

type Encoder interface {
	Encode(v interface{}) error
}

// the global event, mapped onto the HEPEVT common block
var g_evt Event

func GetEvent() *Event {
	evt := &g_evt
	evt.Nevhep = EventNumber()
	evt.Nhep = NumberEntries()
	if len(evt.Isthep) > evt.Nhep {
		evt.Isthep = evt.Isthep[:evt.Nhep]
		evt.Idhep = evt.Idhep[:evt.Nhep]
		evt.Jmohep = evt.Jmohep[:evt.Nhep]
		evt.Jdahep = evt.Jdahep[:evt.Nhep]
		evt.Phep = evt.Phep[:evt.Nhep]
		evt.Vhep = evt.Vhep[:evt.Nhep]
	} else {
		sz := evt.Nhep - len(evt.Isthep)
		evt.Isthep = append(evt.Isthep, make([]int, sz)...)
		evt.Idhep = append(evt.Idhep, make([]int, sz)...)
		evt.Jmohep = append(evt.Jmohep, make([][2]int, sz)...)
		evt.Jdahep = append(evt.Jdahep, make([][2]int, sz)...)
		evt.Phep = append(evt.Phep, make([][5]float64, sz)...)
		evt.Vhep = append(evt.Vhep, make([][4]float64, sz)...)
	}

	for i := 0; i != evt.Nhep; i++ {
		evt.Isthep[i] = StatusCode(i)
		evt.Idhep[i] = PdgId(i)
		// -1 to correct for fortran index value
		evt.Jmohep[i][0] = FirstParent(i) - 1
		//println("::firstparent#",i,evt.Jmohep[i][0])
		// -1 to correct for fortran index value
		evt.Jmohep[i][1] = LastParent(i) - 1
		// -1 to correct for fortran index value
		evt.Jdahep[i][0] = FirstChild(i) - 1
		// -1 to correct for fortran index value
		evt.Jdahep[i][1] = LastChild(i) - 1
		evt.Phep[i][0] = Px(i)
		evt.Phep[i][1] = Py(i)
		evt.Phep[i][2] = Pz(i)
		evt.Phep[i][3] = E(i)
		evt.Phep[i][4] = M(i)
		evt.Vhep[i][0] = X(i)
		evt.Vhep[i][1] = Y(i)
		evt.Vhep[i][2] = Z(i)
		evt.Vhep[i][3] = T(i)
	}
	return evt
}

func SetEvent(evt *Event) {
	SetEventNumber(evt.Nevhep)
	SetNumberEntries(evt.Nhep)
	for i := 0; i < evt.Nhep; i++ {
		SetStatusCode(i, evt.Isthep[i])
		SetPdgId(i, evt.Idhep[i])
		firstparent := evt.Jmohep[i][0] + 1
		lastparent := evt.Jmohep[i][1] + 1
		// if firstparent != lastparent {
		// 	firstparent += 1
		// 	lastparent += 1
		// }
		SetParents(i, firstparent, lastparent)
		firstchild := evt.Jdahep[i][0] + 1
		lastchild := evt.Jdahep[i][1] + 1
		// if firstchild != lastchild || true {
		// 	firstchild += 1
		// 	lastchild += 1
		// }
		SetChildren(i, firstchild, lastchild)
		SetMomentum(i,
			evt.Phep[i][0], evt.Phep[i][1], evt.Phep[i][2], evt.Phep[i][3])
		SetMass(i, evt.Phep[i][4])
		SetPosition(i,
			evt.Vhep[i][0], evt.Vhep[i][1], evt.Vhep[i][2], evt.Vhep[i][3])
	}
}
func WriteTo(enc Encoder) {
	evt := GetEvent()
	//fmt.Printf("--- %v ---\n",evt)

	err := enc.Encode(evt)
	if err != nil {
		fmt.Printf("** error **! %s\n", err)
	}

}
