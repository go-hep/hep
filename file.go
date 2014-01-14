package rootio

import (
	"bytes"
	B "encoding/binary"
	"fmt"
	"io"
	"os"
)

const LargeFileBoundary = 0x7FFFFFFF

var E = B.BigEndian

type Reader interface {
	io.Reader
	io.ReaderAt
	io.Seeker
	io.Closer
}

// A ROOT file is a suite of consecutive data records (TKey's) with
// the following format (see also the TKey class). If the key is
// located past the 32 bit file limit (> 2 GB) then some fields will
// be 8 instead of 4 bytes:
//    1->4            Nbytes    = Length of compressed object (in bytes)
//    5->6            Version   = TKey version identifier
//    7->10           ObjLen    = Length of uncompressed object
//    11->14          Datime    = Date and time when object was written to file
//    15->16          KeyLen    = Length of the key structure (in bytes)
//    17->18          Cycle     = Cycle of key
//    19->22 [19->26] SeekKey   = Pointer to record itself (consistency check)
//    23->26 [27->34] SeekPdir  = Pointer to directory header
//    27->27 [35->35] lname     = Number of bytes in the class name
//    28->.. [36->..] ClassName = Object Class Name
//    ..->..          lname     = Number of bytes in the object name
//    ..->..          Name      = lName bytes with the name of the object
//    ..->..          lTitle    = Number of bytes in the object title
//    ..->..          Title     = Title of the object
//    ----->          DATA      = Data bytes associated to the object
//
// The first data record starts at byte fBEGIN (currently set to kBEGIN).
// Bytes 1->kBEGIN contain the file description, when fVersion >= 1000000
// it is a large file (> 2 GB) and the offsets will be 8 bytes long and
// fUnits will be set to 8:
//    1->4            "root"      = Root file identifier
//    5->8            fVersion    = File format version
//    9->12           fBEGIN      = Pointer to first data record
//    13->16 [13->20] fEND        = Pointer to first free word at the EOF
//    17->20 [21->28] fSeekFree   = Pointer to FREE data record
//    21->24 [29->32] fNbytesFree = Number of bytes in FREE data record
//    25->28 [33->36] nfree       = Number of free data records
//    29->32 [37->40] fNbytesName = Number of bytes in TNamed at creation time
//    33->33 [41->41] fUnits      = Number of bytes for file pointers
//    34->37 [42->45] fCompress   = Compression level and algorithm
//    38->41 [46->53] fSeekInfo   = Pointer to TStreamerInfo record
//    42->45 [54->57] fNbytesInfo = Number of bytes in TStreamerInfo record
//    46->63 [58->75] fUUID       = Universal Unique ID
type File struct {
	Reader
	id string //non-root, identifies filename, etc.

	magic   [4]byte
	version int32
	begin   int64

	// Remainder of record is variable length, 4 or 8 bytes per pointer
	end         int64
	seekfree    int64 // first available record
	nbytesfree  int32 // total bytes available
	nfree       int32 // total free bytes
	nbytesname  int32 // number of bytes in TNamed at creation time
	units       byte
	compression int32
	seekinfo    int64 // pointer to TStreamerInfo
	nbytesinfo  int32 // sizeof(TStreamerInfo)
	uuid        [18]byte

	root directory // root directory of this file
}

// Open opens the named ROOT file for reading. If successful, methods on the
// returned file can be used for reading; the associated file descriptor
// has mode os.O_RDONLY.
func Open(path string) (*File, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %q (%q)", path, err.Error())
	}

	f := &File{
		Reader: fd,
		id:     path,
	}
	f.root = directory{file: f}

	err = f.readHeader()
	if err != nil {
		return nil, err
	}

	return f, nil
}

func (f *File) readHeader() (err error) {

	var stage string

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error reading file named %q while %s (%q)",
				f.id, stage, r.(error).Error())
		}
	}()

	stage = "reading header"

	buf := make([]byte, 64)
	_, err = f.ReadAt(buf, 0)
	if err != nil {
		return err
	}

	dec := NewDecoder(bytes.NewBuffer(buf))

	// Header

	err = dec.readBin(&f.magic)
	if err != nil {
		return err
	}

	if string(f.magic[:]) != "root" {
		return fmt.Errorf("%q is not a root file", f.id)
	}

	err = dec.readInt32(&f.version)
	if err != nil {
		return err
	}

	err = dec.readInt32(&f.begin)
	if err != nil {
		return err
	}

	if f.version < 1000000 { // small file
		err = dec.readInt32(&f.end)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.seekfree)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nbytesfree)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nfree)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nbytesname)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.units)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.compression)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.seekinfo)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nbytesinfo)
		if err != nil {
			return err
		}

	} else { // large files
		err = dec.readInt64(&f.end)
		if err != nil {
			return err
		}

		err = dec.readInt64(&f.seekfree)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nbytesfree)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nfree)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nbytesname)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.units)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.compression)
		if err != nil {
			return err
		}

		err = dec.readInt64(&f.seekinfo)
		if err != nil {
			return err
		}

		err = dec.readInt32(&f.nbytesinfo)
		if err != nil {
			return err
		}
	}

	err = dec.readBin(&f.uuid)
	if err != nil {
		return err
	}

	stage = "read directory info"
	err = f.root.readDirInfo()
	if err != nil {
		return err
	}

	stage = "read streamerinfos"
	err = f.readStreamerInfo()
	if err != nil {
		return err
	}

	stage = "read keys"
	err = f.root.readKeys()
	if err != nil {
		return err
	}

	return err
}

func (f *File) Map() {
	for _, k := range f.root.keys {
		if k.classname == "TBasket" {
			//b := k.AsBasket()
			fmt.Printf("%8s %60s %6v %6v %f\n", k.classname, k.name, k.bytes-k.keylen, k.objlen, float64(k.objlen)/float64(k.bytes-k.keylen))
		} else {
			//println(k.classname, k.name, k.title)
			fmt.Printf("%8s %60s %6v %6v %f\n", k.classname, k.name, k.bytes-k.keylen, k.objlen, float64(k.objlen)/float64(k.bytes-k.keylen))
		}
	}

}

func (f *File) Tell() int64 {
	where, err := f.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	return where
}

// Close closes the File, rendering it unusable for I/O. It returns an
// error, if any.
func (f *File) Close() error {
	for _, k := range f.root.keys {
		k.f = nil
	}
	f.root.keys = nil
	f.root.file = nil
	return f.Reader.Close()
}

// Keys returns the list of keys this File contains
func (f *File) Keys() []Key {
	return f.root.keys
}

func (f *File) Name() string {
	return f.root.Name()
}

func (f *File) Title() string {
	return f.root.Title()
}

func (f *File) Class() string {
	return "TFile"
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
func (f *File) Get(namecycle string) (Object, bool) {
	return f.root.Get(namecycle)
}

// testing interfaces
var _ Object = (*File)(nil)
var _ Directory = (*File)(nil)

// EOF
