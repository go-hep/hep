package rootio

import (
	"bytes"
	"fmt"
	"os"
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

	named named // name+title of this directory
	file  *File // pointer to current file in memory
	keys  []Key
}

// recordSize returns the size of the directory header in bytes
func (dir *directory) recordSize(version int32) int64 {
	var nbytes int64
	nbytes += 2 // fVersion
	nbytes += 4 // ctime
	nbytes += 4 // mtime
	nbytes += 4 // nbyteskeys
	nbytes += 4 // nbytesname
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

func (dir *directory) readDirInfo() error {
	var err error
	f := dir.file
	nbytes := int64(f.nbytesname) + dir.recordSize(f.version)

	if nbytes+f.begin > f.end {
		return fmt.Errorf(
			"rootio: file [%v] has an incorrect header length [%v] or incorrect end of file length [%v]",
			f.id,
			f.begin+nbytes,
			f.end,
		)
	}

	data := make([]byte, int(nbytes))
	_, err = f.ReadAt(data, f.begin)
	if err != nil {
		return err
	}

	tobject_sz := 2 /*version*/ + 4 /* fUniqueID */ + 4 /*fBits*/ + 22 /*process-id*/
	err = dir.named.UnmarshalROOT(data[tobject_sz:f.nbytesname])
	if err != nil {
		return err
	}

	err = dir.UnmarshalROOT(data[f.nbytesname:])
	if err != nil {
		return err
	}

	nk := 4 // Key::fNumberOfBytes
	dec := rootDecoder{r: bytes.NewBuffer(data[nk:])}
	var keyversion int16
	err = dec.readBin(&keyversion)
	if err != nil {
		return err
	}

	if keyversion > 1000 {
		// large files
		nk += 2     // Key::fVersion
		nk += 2 * 4 // Key::fObjectSize, Date
		nk += 2 * 2 // Key::fKeyLength, fCycle
		nk += 2 * 8 // Key::fSeekKey, fSeekParentDirectory
	} else {
		nk += 2     // Key::fVersion
		nk += 2 * 4 // Key::fObjectSize, Date
		nk += 2 * 2 // Key::fKeyLength, fCycle
		nk += 2 * 4 // Key::fSeekKey, fSeekParentDirectory
	}

	dec = rootDecoder{r: bytes.NewBuffer(data[nk:])}
	classname := ""
	err = dec.readString(&classname)
	if err != nil {
		return err
	}
	myprintf("class: [%v]\n", classname)

	cname := ""
	err = dec.readString(&cname)
	if err != nil {
		return err
	}
	myprintf("cname: [%v]\n", cname)

	title := ""
	err = dec.readString(&title)
	if err != nil {
		return err
	}
	myprintf("title: [%v]\n", title)

	if dir.nbytesname < 10 || dir.nbytesname > 1000 {
		return fmt.Errorf("rootio: can't read directory info.")
	}

	return err
}

func (dir *directory) readKeys() error {
	var err error
	if dir.seekkeys <= 0 {
		return nil
	}

	_, err = dir.file.Seek(dir.seekkeys, os.SEEK_SET)
	if err != nil {
		return err
	}

	hdr := Key{f: dir.file}
	err = hdr.Read()
	if err != nil {
		return err
	}
	//myprintf("==> hdr: %#v\n", hdr)

	_, err = dir.file.Seek(dir.seekkeys+int64(hdr.keylen), os.SEEK_SET)
	if err != nil {
		return err
	}
	dec := rootDecoder{r: dir.file}

	var nkeys int32
	err = dec.readInt32(&nkeys)
	if err != nil {
		return err
	}

	for i := 0; i < int(nkeys); i++ {
		err = dir.readKey()
		if err != nil {
			return err
		}
	}
	return err
}

// readKey reads a key and appends it to dir.keys
func (dir *directory) readKey() error {
	dir.keys = append(dir.keys, Key{f: dir.file})
	key := &(dir.keys[len(dir.keys)-1])
	return key.Read()
}

func (dir *directory) Class() string {
	return "TDirectory"
}

func (dir *directory) Name() string {
	return dir.named.Name()
}

func (dir *directory) Title() string {
	return dir.named.Title()
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
func (dir *directory) Get(namecycle string) (Object, bool) {
	name, cycle := decodeNameCycle(namecycle)
	for _, k := range dir.keys {
		if k.Name() == name {
			if cycle != 9999 {
				if k.cycle == cycle {
					return &k, true
				} else {
					return nil, false
				}
			}
			return &k, true
		}
	}
	return nil, false
}

func (dir *directory) UnmarshalROOT(data []byte) error {
	var err error
	dec := rootDecoder{r: bytes.NewBuffer(data)}

	var version int16
	err = dec.readBin(&version)
	if err != nil {
		return err
	}
	myprintf("dir-version: %v\n", version)

	var ctime uint32
	err = dec.readBin(&ctime)
	if err != nil {
		return err
	}
	dir.ctime = datime2time(ctime)
	myprintf("dir-ctime: %v\n", dir.ctime)

	var mtime uint32
	err = dec.readBin(&mtime)
	if err != nil {
		return err
	}
	dir.mtime = datime2time(mtime)
	myprintf("dir-mtime: %v\n", dir.mtime)

	err = dec.readInt32(&dir.nbyteskeys)
	if err != nil {
		return err
	}

	err = dec.readInt32(&dir.nbytesname)
	if err != nil {
		return err
	}

	readptr := dec.readInt64
	if version <= 1000 {
		readptr = dec.readInt32
	}
	err = readptr(&dir.seekdir)
	if err != nil {
		return err
	}

	err = readptr(&dir.seekparent)
	if err != nil {
		return err
	}

	err = readptr(&dir.seekkeys)
	if err != nil {
		return err
	}

	return err
}

// testing interfaces
var _ Object = (*directory)(nil)
var _ Directory = (*directory)(nil)
var _ ROOTUnmarshaler = (*directory)(nil)

// EOF
