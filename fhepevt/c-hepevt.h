#ifndef CHEPEVT_HEPEVT_H
#define CHEPEVT_HEPEVT_H 1

#include <stdio.h>

#ifdef __cplusplus
extern "C" {
#endif

/** write information from HEPEVT common block
 */
void
hepevt_print_hepevt(FILE *f);

/** write particle information
 */
void
hepevt_print_particle(int index, FILE *f);

/** true if common block uses double
 */
int
hepevt_is_double_precision();

/** check for problems with HEPEVT common block
 */
int
hepevt_check_hepevt_consistency(FILE *f);

/** set all entries in HEPEVT to zero
 */
void
hepevt_zero_everything();

/* ---------- access methods ------------------------------------------------
 */

/** event number
 */
int
hepevt_event_number();

/** number of entries in current event
 */
int
hepevt_number_entries();

/** status code
 */
int
hepevt_status_code(int idx);

/** PDG particle id
 */
int
hepevt_pdg_id(int idx);

/** index of 1st mother
 */
int
hepevt_first_parent(int idx);

/** index of last mother
 */
int
hepevt_last_parent(int idx);

/** number of parents
 */
int
hepevt_number_parents(int idx);

/** index of 1st daughter
 */
int
hepevt_first_child(int idx);

/** index of last daughter
 */
int
hepevt_last_child(int idx);

/** number of children
 */
int
hepevt_number_children(int idx);

/** X momentum
 */
double
hepevt_px(int idx);

/** Y momentum
 */
double
hepevt_py(int idx);

/** Z momentum
 */
double
hepevt_pz(int idx);

/** Energy
 */
double
hepevt_e(int idx);

/** generated mass
 */
double
hepevt_m(int idx);

/** X production vertex
 */
double
hepevt_x(int idx);

/** Y production vertex
 */
double
hepevt_y(int idx);

/** Z production vertex
 */
double
hepevt_z(int idx);

/** production time
 */
double
hepevt_t(int idx);

/* set methods ------------------------------------------------------------
 */

/** set event number
 */
void
hepevt_set_event_number(int evtno);

/** set number of entries in HEPEVT
 */
void
hepevt_set_number_entries(int noentries);

/** set particle status
 */
void
hepevt_set_status_code(int index, int status);

/** set particle ID
 */
void
hepevt_set_pdg_id(int index, int id);

/** define parents of a particle
 */
void
hepevt_set_parents(int index, int firstparent, int lastparent);

/** define children of a particle
 */
void
hepevt_set_children(int index, int firstchild, int lastchild);

/** set particle momentum
 */
void
hepevt_set_momentum(int index, 
                    double px, double py, double pz, double e);

/** set particle mass
 */
void
hepevt_set_mass(int index, double mass);

/** set particle production vertex
 */
void
hepevt_set_position(int index, double x, double y, double z, double t);

/* HEPEVT floorplan -------------------------------------------------------
 */

/** size of integer in bytes
 */
int
hepevt_sizeof_int();

/** size of real in bytes
 */
int
hepevt_sizeof_real();

/** size of common block
 */
int 
hepevt_max_number_entries();

/** define size of integer (in bytes)
 */
void
hepevt_set_sizeof_int(int);

/** define size of real (in bytes)
 */
void
hepevt_set_sizeof_real(int);

/** define size of common block
 */
void
hepevt_set_max_number_entries(int);

  
#ifdef __cplusplus
}
#endif

#endif /* !CHEPEVT_HEPEVT_H */
