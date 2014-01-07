package rootio

import (
	"fmt"
	"time"
)

type directory struct {
	ctime      time.Time // time of directory's creation
	mtime      time.Time // time of directory's last modification
	nbyteskeys int32     // number of bytes for the keys
	nbytesname int32     // number of bytes in TNamed at creation time
	seekdir    int64     // location of directory on file
	seekparent int64     // location of parent directory on file
	seekkeys   int64     // location of Keys record on file

	file *File // pointer to current file in memory
	keys []Key
}

// recordSize returns the size of the directory header in bytes
func (dir *directory) recordSize(version int32) int64 {
	nbytes := int64(2) // fVersion
	nbytes += 4        // ctime
	nbytes += 4        // mtime
	nbytes += 4        // nbyteskeys
	nbytes += 4        // nbytesname
	if version >= 40000 {
		// assume that the file may be above 2 Gbytes if file version is > 4
		nbytes += 8 // seekdir
		nbytes += 8 // seekparent
		nbytes += 8 // seekkeys
	} else {
		nbytes += 4 // seekdir
		nbytes += 4 // seekparent
		nbytes += 4 // seekkeys
	}
	return nbytes
}

// Has returns whether an object identified by namecycle exists in directory
//   namecycle has the format name;cycle
//   name  = * is illegal, cycle = * is illegal
//   cycle = "" or cycle = 9999 ==> apply to a memory object
//
//   examples:
//     foo   : get object named foo in memory
//             if object is not in memory, try with highest cycle from file
//     foo;1 : get cycle 1 of foo on file
func (dir *directory) Has(namecycle string) bool {
	name, cycle := decodeNameCycle(namecycle)
	for _, k := range dir.keys {
		if k.Name() == name {
			if cycle != 9999 {
				return k.cycle == cycle
			}
			return true
		}
	}
	return false
}

// Get returns the object identified by namecycle
//   namecycle has the format name;cycle
//   name  = * is illegal, cycle = * is illegal
//   cycle = "" or cycle = 9999 ==> apply to a memory object
//
//   examples:
//     foo   : get object named foo in memory
//             if object is not in memory, try with highest cycle from file
//     foo;1 : get cycle 1 of foo on file
func (dir *directory) Get(namecycle string) (Object, error) {
	name, cycle := decodeNameCycle(namecycle)
	for _, k := range dir.keys {
		if k.Name() == name {
			if cycle != 9999 {
				if k.cycle == cycle {
					return &k, nil
				} else {
					return nil, fmt.Errorf("rootio.File: no such key [%s]", namecycle)
				}
			}
			return &k, nil
		}
	}
	return nil, fmt.Errorf("rootio.File: no such key [%s]", namecycle)
}

// testing interfaces
//var _ Object = (*File)(nil)
var _ Directory = (*directory)(nil)

// EOF
