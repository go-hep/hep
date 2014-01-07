package rootio

import (
	"bufio"
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

type File struct {
	Reader
	id string //non-root, identifies filename, etc.

	magic   [4]byte
	version int32
	begin   int32

	// Remainder of record is variable length, 4 or 8 bytes per pointer
	end         int64
	seekfree    int64 // first available record
	nbytesfree  int64 // total bytes available
	nfree       int32 // total free bytes
	nbytesname  int64 // number of bytes in TNamed at creation time
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

	dec := rootDecoder{r: bufio.NewReader(f)}

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
		var end int32
		err = dec.readBin(&end)
		if err != nil {
			return err
		}
		f.end = int64(end)

		var seekfree int32
		err = dec.readBin(&seekfree)
		if err != nil {
			return err
		}
		f.seekfree = int64(seekfree)

		err = dec.readBin(&f.nbytesfree)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.nfree)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.nbytesname)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.units)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.compression)
		if err != nil {
			return err
		}

		var seekinfo int32
		err = dec.readBin(&seekinfo)
		if err != nil {
			return err
		}
		f.seekinfo = int64(seekinfo)

		err = dec.readBin(&f.nbytesinfo)
		if err != nil {
			return err
		}

	} else { // large files
		err = dec.readBin(&f.end)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.seekfree)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.nbytesfree)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.nfree)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.nbytesname)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.units)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.compression)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.seekinfo)
		if err != nil {
			return err
		}

		err = dec.readBin(&f.nbytesinfo)
		if err != nil {
			return err
		}
	}

	err = dec.readBin(&f.uuid)
	if err != nil {
		return err
	}

	err = f.readBin(&f.uuid)
	if err != nil {
		return err
	}

	stage = "reading keys"

	// Contents of file

	_, err = f.Seek(int64(f.begin), os.SEEK_SET)
	if err != nil {
		return err
	}

	for f.Tell() < f.end {
		err = f.readKey()
		if err != nil {
			return err
		}
	}

	return err
}

// readKey reads a key and appends it to f.keys
func (f *File) readKey() error {
	f.root.keys = append(f.root.keys, Key{f: f})
	key := &(f.root.keys[len(f.root.keys)-1])
	return key.Read()
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

// Has returns whether an object identified by namecycle exists in directory
//   namecycle has the format name;cycle
//   name  = * is illegal, cycle = * is illegal
//   cycle = "" or cycle = 9999 ==> apply to a memory object
//
//   examples:
//     foo   : get object named foo in memory
//             if object is not in memory, try with highest cycle from file
//     foo;1 : get cycle 1 of foo on file
func (f *File) Has(namecycle string) bool {
	return f.root.Has(namecycle)
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
func (f *File) Get(namecycle string) (Object, error) {
	return f.root.Get(namecycle)
}

// testing interfaces
//var _ Object = (*File)(nil)
var _ Directory = (*File)(nil)

// EOF
