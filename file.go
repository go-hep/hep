package rootio

import (
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
	nbytesname  int64
	units       byte
	compression int32
	seekinfo    int64 // pointer to TStreamerInfo
	nbytesinfo  int32 // sizeof(TStreamerInfo)
	uuid        [18]byte

	keys []Key
}

// Open opens the named ROOT file for reading. If successful, methods on the
// returned file can be used for reading; the associated file descriptor
// has mode os.O_RDONLY.
func Open(path string) (*File, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %q (%q)", path, err.Error())
	}

	f := &File{Reader: fd, id: path}

	return f, f.readHeader()
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

	// Header

	err = f.readBin(&f.magic)
	if err != nil {
		return err
	}

	if string(f.magic[:]) != "root" {
		return fmt.Errorf("%q is not a root file", f.id)
	}

	err = f.readInt32(&f.version)
	if err != nil {
		return err
	}

	err = f.readInt32(&f.begin)
	if err != nil {
		return err
	}

	err = f.readPtr(&f.end)
	if err != nil {
		return err
	}

	err = f.readPtr(&f.seekfree)
	if err != nil {
		return err
	}

	err = f.readPtr(&f.nbytesfree)
	if err != nil {
		return err
	}

	err = f.readInt32(&f.nfree)
	if err != nil {
		return err
	}

	err = f.readPtr(&f.nbytesname)
	if err != nil {
		return err
	}

	err = f.readBin(&f.units)
	if err != nil {
		return err
	}

	err = f.readInt32(&f.compression)
	if err != nil {
		return err
	}

	err = f.readPtr(&f.seekinfo)
	if err != nil {
		return err
	}

	err = f.readInt32(&f.nbytesinfo)
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
	f.keys = append(f.keys, Key{f: f})
	key := &(f.keys[len(f.keys)-1])
	return key.Read()
}

func (f *File) readBin(v interface{}) error {
	return B.Read(f, E, v)
}

func (f *File) Map() {
	for _, k := range f.keys {
		if k.classname == "TBasket" {
			//b := k.AsBasket()
			fmt.Printf("%8s %60s %6v %6v %f\n", k.classname, k.name, k.bytes-k.keylen, k.objlen, float64(k.objlen)/float64(k.bytes-k.keylen))
		} else {
			//println(k.classname, k.name, k.title)
			fmt.Printf("%8s %60s %6v %6v %f\n", k.classname, k.name, k.bytes-k.keylen, k.objlen, float64(k.objlen)/float64(k.bytes-k.keylen))
		}
	}

}

func (f *File) readString(s *string) error {
	var err error
	var length byte
	var buf [256]byte

	err = f.readBin(&length)
	if err != nil {
		return err
	}

	if length != 0 {
		err = f.readBin(buf[:length])
		if err != nil {
			return err
		}
		*s = string(buf[:length])
	}
	return err
}

func (f *File) readInt16(v interface{}) error {
	var err error
	var d int16
	err = f.readBin(&d)
	if err != nil {
		return err
	}

	switch uv := v.(type) {
	case *int32:
		*uv = int32(d)
	case *int64:
		*uv = int64(d)
	default:
		panic("Unknown type")
	}

	return err
}

func (f *File) readInt32(v interface{}) error {
	var err error
	switch uv := v.(type) {
	case *int32:
		err = f.readBin(v)
	case *int64:
		var d int32
		err = f.readBin(&d)
		*uv = int64(d)
	default:
		panic("Unknown type")
	}
	return err
}

func (f *File) readPtr(v interface{}) error {
	var err error
	if f.version > 1000000 {
		err = f.readBin(v)
	} else {
		err = f.readInt32(v)
	}
	return err
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
	for _, k := range f.keys {
		k.f = nil
	}
	f.keys = nil
	return f.Reader.Close()
}

// Keys returns the list of keys this File contains
func (f *File) Keys() []Key {
	return f.keys
}

// EOF
