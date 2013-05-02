#include "c-hepevt.h"

#ifndef HEPEVT_EntriesAllocation
#define HEPEVT_EntriesAllocation 10000
#endif  // HEPEVT_EntriesAllocation

#ifdef __cplusplus
const unsigned int hepevt_bytes_allocation = 
  sizeof(long int) * ( 2 + 6 * HEPEVT_EntriesAllocation )
  + sizeof(double) * ( 9 * HEPEVT_EntriesAllocation );
#else 
#define hepevt_bytes_allocation                               \
  (sizeof(long int) * ( 2 + 6 * HEPEVT_EntriesAllocation )  \
   + sizeof(double) * ( 9 * HEPEVT_EntriesAllocation ))
#endif

#ifdef _WIN32 // Platform: Windows MS Visual C++
struct HEPEVT_DEF{
        char data[hepevt_bytes_allocation];
    };
#ifdef __cplusplus
extern "C"
#endif
HEPEVT_DEF HEPEVT;
#define hepevt HEPEVT

#else
#ifdef __cplusplus
extern "C" {
#endif
  extern struct {
    char data[hepevt_bytes_allocation];
  } hepevt_ __attribute__((weak));
#ifdef __cplusplus
} /* extern "C" */
#endif
#define hepevt hepevt_
#endif // Platform

#ifdef __cplusplus
extern "C" {
#endif

/* a few locals -------------------------------------------------------------
 */
static int s_sizeof_int = 4;
static int s_sizeof_real = sizeof(double);
static int s_max_number_entries = 4000;

/* utils functions ----------------------------------------------------------
 */
static
double
byte_num_to_double(unsigned int b)
{
  if (b>= hepevt_bytes_allocation) {
    fprintf(stderr,
            "c-hepevt: requested hepevt data exceeds allocation\n");
    return 0;
  }
  if (s_sizeof_real == sizeof(float)) {
    float *f = (float*)&hepevt.data[b];
    return (double)(*f);
  }
  if (s_sizeof_real == sizeof(double)) {
    double *d = (double*)&hepevt.data[b];
    return *d;
  }
  fprintf(stderr,
          "c-hepevt: illegal floating point number length [%i].\n",
          s_sizeof_real);
  return 0;
}

static
int
byte_num_to_int(unsigned int b)
{
  if (b>= hepevt_bytes_allocation) {
    fprintf(stderr,
            "c-hepevt: requested hepevt data exceeds allocation\n");
    return 0;
  }
  if (s_sizeof_int == sizeof(short int)) {
    short int *si = (short int*)&hepevt.data[b];
    return (int)(*si);
  }
  if (s_sizeof_int == sizeof(long int)) {
    long int *li = (long int*)&hepevt.data[b];
    return *li;
  }
  if (s_sizeof_int == sizeof(int)) {
    int *li = (int*)&hepevt.data[b];
    return *li;
  }
  fprintf(stderr,
          "c-hepevt: illegal integer number length [%i].\n",
          s_sizeof_int);
  return 0;
}

static
void
write_byte_num_d(double data, unsigned int b)
{
  if (b>= hepevt_bytes_allocation) {
    fprintf(stderr,
            "c-hepevt: requested hepevt data exceeds allocation\n");
    return;
  }

  if (s_sizeof_real == sizeof(float)) {
    float *f = (float*)&hepevt.data[b];
    *f = (float)data;
    return;
  }
  if (s_sizeof_real == sizeof(double)) {
    double *d = (double*)&hepevt.data[b];
    *d = (double)data;
    return;
  }
  fprintf(stderr,
          "c-hepevt: illegal floating point number length [%i].\n",
          s_sizeof_real);
  return;
}

static
void
write_byte_num_i(int data, unsigned int b)
{
  if (b>= hepevt_bytes_allocation) {
    fprintf(stderr,
            "c-hepevt: requested hepevt data exceeds allocation\n");
    return;
  }

  if (s_sizeof_int == sizeof(short int)) {
    short int *i = (short int*)&hepevt.data[b];
    *i = (short int)data;
    return;
  }
  if (s_sizeof_int == sizeof(long int)) {
    long int *i = (long int*)&hepevt.data[b];
    *i = (int)data;
    return;
  }
  if (s_sizeof_int == sizeof(int)) {
    int *i = (int*)&hepevt.data[b];
    *i = (int)data;
    return;
  }
  fprintf(stderr,
          "c-hepevt: illegal integer number length [%i].\n",
          s_sizeof_int);
  return;
}

static
void
print_legend(FILE *f)
{
  fprintf(f,
          "%4s %4s %4s %5s   %10s, %9s, %9s, %9s, %10s\n",
          "Indx","Stat","Par-","chil-",
          "(  P_x","P_y","P_z","Energy","M ) ");
  fprintf(f,
          "%9s %4s %4s    %10s, %9s, %9s, %9s) %9s\n",
          "ID ","ents","dren",
          "Prod (   X","Y","Z","cT", "[mm]");
}

/* HEPEVT C-API -------------------------------------------------------------
 */

/** write information from HEPEVT common block
 */
void hepevt_print_hepevt(FILE *f)
{
  if (NULL == f) {
    f = stdout;
  }
  fprintf(f, 
          "________________________________________"
          "________________________________________\n");
  fprintf(f, 
          "***** HEPEVT Common Event#: %i, %i particles (max %i) *****",
          hepevt_event_number(),
          hepevt_number_entries(),
          hepevt_max_number_entries()
          );
  if (hepevt_is_double_precision()) {
    fprintf(f, " Double Precision");
  } else {
    fprintf(f, " Single Precision");
  }
  fprintf(f, 
          "\n%d-byte integers, "
          "%d-byte floating point numbers, "
          "%d-allocated entries.\n",
          s_sizeof_int,
          s_sizeof_real,
          hepevt_max_number_entries());
  print_legend(f);
  fprintf(f, 
          "________________________________________"
          "________________________________________\n");
  int i = 1;
  for (i = 1; i <= hepevt_number_entries(); ++i) {
    hepevt_print_particle(i, f);
  }
  fprintf(f, 
          "________________________________________"
          "________________________________________\n");
  fflush(f);
  return;
}

/** write particle information
 */
void hepevt_print_particle(int i, FILE *f)
{
  /// dumps the content HEPEVT particle entry i   (Width is 120)
  /// here i is the C array index (i.e. it starts at 0 ... whereas the
  /// fortran array index starts at 1) So if there's 100 particles, the
  /// last valid index is 100-1=99
  if (f == NULL) {
    f = stdout;
  }
  fprintf(f, 
          "%4d %+4d %4d %4d    (%9.3g, %9.3g, %9.3g, %9.3g, %9.3g)\n",
          i, hepevt_status_code(i), 
          hepevt_first_parent(i), 
          hepevt_first_child(i),
          hepevt_px(i), hepevt_py(i), hepevt_pz(i), 
          hepevt_e(i),  hepevt_m(i) );
  fprintf(f,
          "%+9d %4d %4d    (%9.3g, %9.3g, %9.3g, %9.3g)\n",
          hepevt_pdg_id(i), 
          hepevt_last_parent(i), 
          hepevt_last_child(i), 
          hepevt_x(i), hepevt_y(i), hepevt_z(i), hepevt_t(i) );
  return;
}

/** true if common block uses double
 */
int hepevt_is_double_precision()
{
  if (sizeof(double) == s_sizeof_real) {
    return 1;
  }
  return 0;
}

/** check for problems with HEPEVT common block
 */
int
hepevt_check_hepevt_consistency(FILE *f)
{
  if (f == NULL) {
    f = stdout;
  }
  const char *hdr = "\n\n\t*** WARNING Inconsistent HEPEVT input, Event %10d ***\n";
  int is_consistent = 1;
  int i = 1;
  for (i = 1; i <= hepevt_number_entries(); ++i) {
    // 1. check its mothers
    int moth1 = hepevt_first_parent(i);
    int moth2 = hepevt_last_parent(i);
    if ( moth2 < moth1 ) {
      if (is_consistent) {
        is_consistent = 0;
        fprintf(f, hdr, hepevt_event_number());
        print_legend(f);
      }
      fprintf(f,
              "Inconsistent entry %d first parent > last parent\n",
              i);
      hepevt_print_particle(i, f);
    }
    int m = 0;
    for (m = moth1; m <= moth2 && m != 0; ++m) {
      if (m > hepevt_number_entries() || m < 0) {
        if (is_consistent) {
          fprintf(f, hdr, hepevt_event_number());
          is_consistent = 0;
          print_legend(f);
        }
        fprintf(f, "Inconsistent entry %d mother points out of range\n",
                i);
        hepevt_print_particle(i, f);
      }
      int mchild1 = hepevt_first_child(m);
      int mchild2 = hepevt_last_child(m);
      // we dont consider null pointers as inconsistent
      if (mchild1==0 && mchild2==0) {
        continue;
      }
      if (i<mchild1 || i>mchild2) {
        if (is_consistent) {
          fprintf(f, hdr, hepevt_event_number());
          is_consistent = 0;
          print_legend(f);
        }
        fprintf(f,
                "Inconsistent mother-daughter relationship between %d & "
                "%d (try !trust_mother)\n",
                i, m);
        hepevt_print_particle(i, f);
        hepevt_print_particle(m, f);
      }
    }
    // 2. check its daughters
    int dau1 = hepevt_first_child(i);
    int dau2 = hepevt_last_child(i);
    if (dau2 < dau1) {
      if (is_consistent) {
        fprintf(f, hdr, hepevt_event_number());
        is_consistent = 0;
        print_legend(f);
      }
      fprintf(f,
              "Inconsistent entry %d first child > last child\n",
              i);
      hepevt_print_particle(i, f);
    }
    int d = dau1;
    for (d = dau1; d<=dau2 && d!=0; ++d) {
      if (d > hepevt_number_entries() || d < 0) {
        if (is_consistent) {
          fprintf(f, hdr, hepevt_event_number());
          is_consistent = 0;
          print_legend(f);
        }
        fprintf(f,
                "Inconsistent entry %d child points out of range\n",
                i);
        hepevt_print_particle(i, f);
      }
      int d_moth1 = hepevt_first_parent(d);
      int d_moth2 = hepevt_last_parent(d);
      // we dont consider null pointers as inconsistent
      if (d_moth1==0 && d_moth2==0) {
        continue;
      }
      if (i<d_moth1 || i>d_moth2) {
        if (is_consistent) {
          fprintf(f, hdr, hepevt_event_number());
          is_consistent = 0;
          print_legend(f);
        }
        fprintf(f,
                "Inconsistent mother-daughter relationship between "
                "%d & %d (try trust_mothers)\n",
                i, d);
        hepevt_print_particle(i, f);
        hepevt_print_particle(d, f);
      }
    }
  } // loop over entries

  if (!is_consistent) {
    fprintf
      (f,
       "\n"
       );
  }
  return is_consistent;
}

/** set all entries in HEPEVT to zero
 */
void
hepevt_zero_everything()
{
  hepevt_set_event_number(0);
  hepevt_set_number_entries(0);
  int i = 1;
  for (i = 1; i<=s_max_number_entries; ++i) {
    hepevt_set_status_code(i, 0);
    hepevt_set_pdg_id(i, 0);
    hepevt_set_parents( i, 0, 0 );
    hepevt_set_children( i, 0, 0 );
    hepevt_set_momentum( i, 0, 0, 0, 0 );
    hepevt_set_mass( i, 0 );
    hepevt_set_position( i, 0, 0, 0, 0 );
  }
}

/* ---------- access methods ------------------------------------------------
 */

/** event number
 */
int
hepevt_event_number()
{
  return byte_num_to_int(0);
}

/** number of entries in current event
 */
int
hepevt_number_entries()
{
  int nhep = byte_num_to_int( 1*s_sizeof_int );
  return (nhep <= s_max_number_entries 
          ? nhep
          : s_max_number_entries);
}

/** status code
 */
int
hepevt_status_code(int idx)
{
  return byte_num_to_int( (2+idx-1) * s_sizeof_int );
}

/** PDG particle id
 */
int
hepevt_pdg_id(int idx)
{
  return byte_num_to_int( (2+s_max_number_entries+idx-1) 
                          * s_sizeof_int );
}

/** index of 1st mother
 */
int
hepevt_first_parent(int idx)
{
  int parent = byte_num_to_int( (2+2*s_max_number_entries+2*(idx-1)) 
                                * s_sizeof_int ); 
  return ( parent > 0 && parent <= hepevt_number_entries() ) 
    ? parent 
    : 0;
}

/** index of last mother
 */
int
hepevt_last_parent(int idx)
{
  // Returns the Index of the LAST parent in the HEPEVT record
  // for particle with Index index.
  // If there is only one parent, the last parent is forced to 
  // be the same as the first parent.
  // If there are no parents for this particle, both the first_parent
  // and the last_parent with return 0.
  // Error checking is done to ensure the parent is always
  // within range ( 0 <= parent <= nhep )
  //
  int firstparent = hepevt_first_parent(idx);
  int parent = byte_num_to_int( (2+2*s_max_number_entries+2*(idx-1)+1) 
                                * s_sizeof_int ); 
  return ( parent > firstparent && parent <= hepevt_number_entries() ) 
    ? parent 
    : firstparent; 
}

/** number of parents
 */
int
hepevt_number_parents(int idx)
{
  int firstparent = hepevt_first_parent(idx);
  return ( firstparent>0 ) ? 
    ( 1+hepevt_last_parent(idx)-firstparent ) : 0;
}

/** index of 1st daughter
 */
int
hepevt_first_child(int idx)
{
  int child = byte_num_to_int( (2+4*s_max_number_entries+2*(idx-1)) 
                               * s_sizeof_int ); 
  return ( child > 0 && child <= hepevt_number_entries() ) 
    ? child 
    : 0; 
}

/** index of last daughter
 */
int
hepevt_last_child(int idx)
{
  // Returns the Index of the LAST child in the HEPEVT record
  // for particle with Index index.
  // If there is only one child, the last child is forced to 
  // be the same as the first child.
  // If there are no children for this particle, both the first_child
  // and the last_child with return 0.
  // Error checking is done to ensure the child is always
  // within range ( 0 <= parent <= nhep )
  //
  int firstchild = hepevt_first_child(idx);
  int child = byte_num_to_int( (2+4*s_max_number_entries+2*(idx-1)+1) 
                               * s_sizeof_int ); 
  return ( child > firstchild && child <= hepevt_number_entries() ) 
    ? child 
    : firstchild;
}

/** number of children
 */
int
hepevt_number_children(int idx)
{
  int firstchild = hepevt_first_child(idx);
  return ( firstchild>0 ) 
    ? ( 1+hepevt_last_child(idx)-firstchild ) 
    : 0;
}

/** X momentum
 */
double
hepevt_px(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + (5*(idx-1)+0) * s_sizeof_real );
}

/** Y momentum
 */
double
hepevt_py(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + (5*(idx-1)+1) * s_sizeof_real );
}

/** Z momentum
 */
double
hepevt_pz(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + (5*(idx-1)+2) * s_sizeof_real );
}

/** Energy
 */
double
hepevt_e(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + (5*(idx-1)+3) * s_sizeof_real );
}

/** generated mass
 */
double
hepevt_m(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + (5*(idx-1)+4) * s_sizeof_real );
}

/** X production vertex
 */
double
hepevt_x(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + ( 5*s_max_number_entries
                                 + (4*(idx-1)+0) ) *s_sizeof_real );
}

/** Y production vertex
 */
double
hepevt_y(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + ( 5*s_max_number_entries
                                 + (4*(idx-1)+1) ) *s_sizeof_real );
}

/** Z production vertex
 */
double
hepevt_z(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + ( 5*s_max_number_entries
                                 + (4*(idx-1)+2) ) *s_sizeof_real );
}

/** production time
 */
double
hepevt_t(int idx)
{
  return byte_num_to_double( (2+6*s_max_number_entries)*s_sizeof_int
                             + ( 5*s_max_number_entries
                                 + (4*(idx-1)+3) ) *s_sizeof_real );
}

/* set methods ------------------------------------------------------------
 */

/** set event number
 */
void
hepevt_set_event_number(int evtno)
{ 
  write_byte_num_i( evtno, 0 ); 
}

/** set number of entries in HEPEVT
 */
void
hepevt_set_number_entries(int noentries)
{ 
  write_byte_num_i( noentries, 1*s_sizeof_int ); 
}

/** set particle status
 */
void
hepevt_set_status_code(int idx, int status)
{
  if ( idx <= 0 || idx > s_max_number_entries ) {
    return;
  }
  write_byte_num_i( status, (2+idx-1) * s_sizeof_int );
}

/** set particle ID
 */
void
hepevt_set_pdg_id(int idx, int id)
{
  if ( idx <= 0 || idx > s_max_number_entries ) {
    return;
  }
  write_byte_num_i( id, (2+s_max_number_entries+idx-1) *s_sizeof_int );
}

/** define parents of a particle
 */
void
hepevt_set_parents(int index, int firstparent, int lastparent)
{
  if ( index <= 0 || index > s_max_number_entries ) {
    return;
  }
  write_byte_num_i( firstparent, 
                    (2+2*s_max_number_entries+2*(index-1)) 
                    * s_sizeof_int );
  write_byte_num_i( lastparent, 
                    (2+2*s_max_number_entries+2*(index-1)+1) 
				    * s_sizeof_int );
}

/** define children of a particle
 */
void
hepevt_set_children(int index, int firstchild, int lastchild)
{
  if ( index <= 0 || index > s_max_number_entries ) {
    return;
  }
  write_byte_num_i( firstchild, (2+4*s_max_number_entries+2*(index-1)) 
                    *s_sizeof_int );
  write_byte_num_i( lastchild, (2+4*s_max_number_entries+2*(index-1)+1) 
				    *s_sizeof_int );
}

/** set particle momentum
 */
void
hepevt_set_momentum(int index, 
                    double px, double py, double pz, double e)
{
  if ( index <= 0 || index > s_max_number_entries ) {
    return;
  }
  write_byte_num_d( px, (2+6*s_max_number_entries) *s_sizeof_int
                  + (5*(index-1)+0) *s_sizeof_real );
  write_byte_num_d( py, (2+6*s_max_number_entries)*s_sizeof_int
                  + (5*(index-1)+1) *s_sizeof_real );
  write_byte_num_d( pz, (2+6*s_max_number_entries)*s_sizeof_int
                  + (5*(index-1)+2) *s_sizeof_real );
  write_byte_num_d( e,  (2+6*s_max_number_entries)*s_sizeof_int
                  + (5*(index-1)+3) *s_sizeof_real );
}

/** set particle mass
 */
void
hepevt_set_mass(int index, double mass)
{
  if ( index <= 0 || index > s_max_number_entries ) {
    return;
  }
  write_byte_num_d( mass, (2+6*s_max_number_entries)*s_sizeof_int
                    + (5*(index-1)+4) *s_sizeof_real );
}

/** set particle production vertex
 */
void
hepevt_set_position(int index, double x, double y, double z, double t)
{
  if ( index <= 0 || index > s_max_number_entries ) {
    return;
  }
  write_byte_num_d( x, (2+6*s_max_number_entries)*s_sizeof_int
                    + ( 5*s_max_number_entries
                        + (4*(index-1)+0) ) *s_sizeof_real );
  write_byte_num_d( y, (2+6*s_max_number_entries)*s_sizeof_int
                    + ( 5*s_max_number_entries
                        + (4*(index-1)+1) ) *s_sizeof_real );
  write_byte_num_d( z, (2+6*s_max_number_entries)*s_sizeof_int
                    + ( 5*s_max_number_entries
                        + (4*(index-1)+2) ) *s_sizeof_real );
  write_byte_num_d( t, (2+6*s_max_number_entries)*s_sizeof_int
                    + ( 5*s_max_number_entries
                        + (4*(index-1)+3) ) *s_sizeof_real );
}

/* HEPEVT floorplan -------------------------------------------------------
 */

/** size of integer in bytes
 */
int
hepevt_sizeof_int()
{
  return s_sizeof_int;
}

/** size of real in bytes
 */
int
hepevt_sizeof_real()
{
  return s_sizeof_real;
}

/** size of common block
 */
int 
hepevt_max_number_entries()
{
  return s_max_number_entries;
}

/** define size of integer (in bytes)
 */
void
hepevt_set_sizeof_int(int sz)
{
  if (sz != sizeof(short int) && 
      sz != sizeof(long int)  &&
      sz != sizeof(int)) {
    fprintf(stderr,
            "c-hepevt is not able to handle integers "
            "of size other than 2 or 4. "
            "you requested [%d]\n",
            sz);
    return;
  }
  s_sizeof_int = sz;
}

/** define size of real (in bytes)
 */
void
hepevt_set_sizeof_real(int sz)
{
  if (sz != sizeof(float) && 
      sz != sizeof(double)) {
    fprintf(stderr,
            "c-hepevt is not able to handle floating point numbers "
            "of size other than 4 or 8. "
            "you requested [%d]\n",
            sz);
    return;
  }
  s_sizeof_real = sz;
}

/** define size of common block
 */
void
hepevt_set_max_number_entries(int entries)
{
  s_max_number_entries = entries;
}


#ifdef __cplusplus
} // extern "C"
#endif
