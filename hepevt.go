// interface to the FORTRAN HEPEVT common block
package hepevt

/*
 #include "c-hepevt.h"

 #include <stdio.h>
 #include <stdlib.h>
*/
import "C"

import (
	"os"
	"unsafe"
)

// write information from HEPEVT common block
func PrintHepevt(f *os.File) {
	if f == nil {
		f = os.Stdout
	}
	c_fd := C.int(f.Fd())
	_ = f.Sync()
	c_mode := C.CString("a")
	defer C.free(unsafe.Pointer(c_mode))
	c_f := C.fdopen(c_fd, c_mode)
	C.fflush(c_f)
	C.hepevt_print_hepevt(c_f)
	C.fflush(c_f)
	_ = f.Sync()
}

// write particle information
func PrintParticle(idx int, f *os.File) {
	if f == nil {
		f = os.Stdout
	}
	c_fd := C.int(f.Fd())
	_ = f.Sync()
	c_mode := C.CString("a")
	defer C.free(unsafe.Pointer(c_mode))
	c_f := C.fdopen(c_fd, c_mode)
	C.fflush(c_f)
	C.hepevt_print_particle(C.int(idx+1), c_f)
	C.fflush(c_f)
	_ = f.Sync()
}

// true if common block uses double
func IsDoublePrecision() bool {
	c := C.hepevt_is_double_precision()
	if c != C.int(0) {
		return true
	}
	return false
}

// check for problems with HEPEVT common block
func CheckHepevtConsistency(f *os.File) bool {
	
	if f == nil {
		f = os.Stdout
	}
	c_fd := C.int(f.Fd())
	_ = f.Sync()
	c_mode := C.CString("a")
	defer C.free(unsafe.Pointer(c_mode))
	c_f := C.fdopen(c_fd, c_mode)
	C.fflush(c_f)

	o := C.hepevt_check_hepevt_consistency(c_f)
	C.fflush(c_f)
	_ = f.Sync()

	if o != C.int(0) {
		return true
	}
	return false
}

// set all entries in HEPEVT to zero
func ZeroEverything() {
	C.hepevt_zero_everything()
}

// access methods ------------------------------------------------------------

// event number
func EventNumber() int {
	return int(C.hepevt_event_number())
}

// number of entries in current event
func NumberEntries() int {
	return int(C.hepevt_number_entries())
}

// status code
func StatusCode(idx int) int {
	return int(C.hepevt_status_code(C.int(idx+1)))
}

// PDG particle id
func PdgId(idx int) int {
	return int(C.hepevt_pdg_id(C.int(idx+1)))
}

// index of 1st mother
func FirstParent(idx int) int {
	v := int(C.hepevt_first_parent(C.int(idx+1)))
	//println("::firstparent",idx,v)
	return v
}

// index of last mother
func LastParent(idx int) int {
	v := int(C.hepevt_last_parent(C.int(idx+1)))
	//println("::lastparent ",idx,v)
	return v
}

// number of parents
func NumberParents(idx int) int {
	return int(C.hepevt_number_parents(C.int(idx+1)))
}

// index of 1st daughter
func FirstChild(idx int) int {
	v := int(C.hepevt_first_child(C.int(idx+1)))
	//println("::firstchild:",idx,v)
	return v
}

// index of last daughter
func LastChild(idx int) int {
	v := int(C.hepevt_last_child(C.int(idx+1)))
	//println("::lastchild: ",idx,v)
	return v
}

// number of children
func NumberChildren(idx int) int {
	return int(C.hepevt_number_children(C.int(idx+1)))
}

// X momentum
func Px(idx int) float64 {
	return float64(C.hepevt_px(C.int(idx+1)))
}

// Y momentum
func Py(idx int) float64 {
	return float64(C.hepevt_py(C.int(idx+1)))
}

// Z momentum
func Pz(idx int) float64 {
	return float64(C.hepevt_pz(C.int(idx+1)))
}

// Energy
func E(idx int) float64 {
	return float64(C.hepevt_e(C.int(idx+1)))
}

// Generated mass
func M(idx int) float64 {
	return float64(C.hepevt_m(C.int(idx+1)))
}

// X production vertex
func X(idx int) float64 {
	return float64(C.hepevt_x(C.int(idx+1)))
}

// Y production vertex
func Y(idx int) float64 {
	return float64(C.hepevt_y(C.int(idx+1)))
}

// Z production vertex
func Z(idx int) float64 {
	return float64(C.hepevt_z(C.int(idx+1)))
}

// production time
func T(idx int) float64 {
	return float64(C.hepevt_t(C.int(idx+1)))
}

// set methods ---------------------------------------------------------------

// FIXME: make sure the massaging of indices is consistent...
//         - do we expect a C-based index or a FORTRAN one ?
//         - for index ?
//         - for firstchild/lastchild ?

// set event number
func SetEventNumber(evtno int) {
	C.hepevt_set_event_number(C.int(evtno))
}

// set number of entries in HEPEVT
func SetNumberEntries(nentries int) {
	C.hepevt_set_number_entries(C.int(nentries))
}

// set particle status code
func SetStatusCode(index, status int) {
	C.hepevt_set_status_code(C.int(index+1), C.int(status))
}

// set particle PDG-id
func SetPdgId(index, id int) {
	C.hepevt_set_pdg_id(C.int(index+1), C.int(id))
}

// define parents of a particle
func SetParents(index, firstparent, lastparent int) {
	C.hepevt_set_parents(
		C.int(index+1),
		C.int(firstparent),
		C.int(lastparent))
}

// define children of a particle
func SetChildren(index, firstchild, lastchild int) {
	C.hepevt_set_children(
		C.int(index+1),
		C.int(firstchild),
		C.int(lastchild))
}

// set particle momentum
func SetMomentum(index int, px, py, pz, e float64) {
	C.hepevt_set_momentum(
		C.int(index+1),
		C.double(px), C.double(py),
		C.double(pz), C.double(e))
}

// set particle mass
func SetMass(index int, m float64) {
	C.hepevt_set_mass(C.int(index+1), C.double(m))
}

// set particle production vertex
func SetPosition(index int, x, y, z, t float64) {
	C.hepevt_set_position(
		C.int(index+1),
		C.double(x), C.double(y),
		C.double(z), C.double(t))
}

// HEPEVT floorplan ---------------------------------------------------------

// size of integer in bytes
func SizeofInt() int {
	return int(C.hepevt_sizeof_int())
}

// size of real in bytes
func SizeofReal() int {
	return int(C.hepevt_sizeof_real())
}

// size of common block
func MaxNumberEntries() int {
	return int(C.hepevt_max_number_entries())
}

// define size of integer
func SetSizeofInt(sz int) {
	C.hepevt_set_sizeof_int(C.int(sz))
}

// define size of real
func SetSizeofReal(sz int) {
	C.hepevt_set_sizeof_real(C.int(sz))
}

// define size of common block
func SetMaxNumberEntries(sz int) {
	C.hepevt_set_max_number_entries(C.int(sz))
}

//func 