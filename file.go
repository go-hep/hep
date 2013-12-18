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

func Open(path string) (*File, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Unable to open %q (%q)", path, err.Error())
	}

	f := &File{Reader: fd, id: path}

	return f, f.ReadHeader()
}

func (f *File) ReadHeader() (err error) {

	var stage string

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Error reading file named %q while %s (%q)",
				f.id, stage, r.(error).Error())
		}
	}()

	stage = "reading header"

	// Header

	f.ReadBin(&f.magic)

	if string(f.magic[:]) != "root" {
		return fmt.Errorf("%q is not a root file", f.id)
	}

	f.ReadInt32(&f.version)

	f.ReadInt32(&f.begin)
	f.ReadPtr(&f.end)

	f.ReadPtr(&f.seekfree)
	f.ReadPtr(&f.nbytesfree)

	f.ReadInt32(&f.nfree)
	f.ReadPtr(&f.nbytesname)
	f.ReadBin(&f.units)
	f.ReadInt32(&f.compression)
	f.ReadPtr(&f.seekinfo)
	f.ReadInt32(&f.nbytesinfo)
	f.ReadBin(&f.uuid)

	stage = "reading keys"

	// Contents of file

	f.Seek(int64(f.begin), os.SEEK_SET)

	for f.Tell() < f.end {
		f.ReadKey()
	}

	return nil
}

// Read a key and append it to f.keys
func (f *File) ReadKey() {
	f.keys = append(f.keys, Key{f: f})
	key := &(f.keys[len(f.keys)-1])
	key.Read()
}

func (f *File) ReadBin(v interface{}) {
	err := B.Read(f, E, v)
	if err != nil {
		panic(err)
	}
}

var buf [256]byte

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

func (f *File) ReadString(s *string) {
	var length byte

	f.ReadBin(&length)

	if length != 0 {
		f.ReadBin(buf[:length])
		*s = string(buf[:length])
	}
}

func (f *File) ReadInt16(v interface{}) {
	var d int16
	f.ReadBin(&d)

	switch uv := v.(type) {
	case *int32:
		*uv = int32(d)
	case *int64:
		*uv = int64(d)
	default:
		panic("Unknown type")
	}
}

func (f *File) ReadInt32(v interface{}) {
	switch uv := v.(type) {
	case *int32:
		f.ReadBin(v)
	case *int64:
		var d int32
		f.ReadBin(&d)
		*uv = int64(d)
	default:
		panic("Unknown type")
	}
}

func (f *File) ReadPtr(v interface{}) {
	if f.version > 1000000 {
		f.ReadBin(v)
	} else {
		f.ReadInt32(v)
	}
}

func (f *File) Tell() int64 {
	where, err := f.Seek(0, os.SEEK_CUR)
	if err != nil {
		panic(err)
	}
	return where
}

func (f *File) Close() error {
	for _, k := range f.keys {
		k.f = nil
	}
	f.keys = nil
	return f.Reader.Close()
}
