package rootio

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"

	//"github.com/kr/pretty"
)

// Key is a key (a label) in a ROOT file
//
//  The Key class includes functions to book space on a file,
//   to create I/O buffers, to fill these buffers
//   to compress/uncompress data buffers.
//
//  Before saving (making persistent) an object on a file, a key must
//  be created. The key structure contains all the information to
//  uniquely identify a persistent object on a file.
//  The Key class is used by ROOT:
//    - to write an object in the Current Directory
//    - to write a new ntuple buffer
type Key struct {
	f *File // underlying file

	bytes    int32 // number of bytes for the compressed object+key
	version  int16 // version of the Key struct
	objlen   int32 // length of uncompressed object
	datetime int32 // Date/Time when the object was written
	keylen   int32 // number of bytes for the Key struct
	cycle    int16 // cycle number of the object

	// address of the object on file (points to Key.bytes)
	// this is a redundant information used to cross-check
	// the data base integrity
	seekkey  int64
	seekpdir int64 // pointer to the directory supporting this object

	classname string // object class name
	name      string // name of the object
	title     string // title of the object

	pdat int64 // Pointer to everything after the above.

	read bool
	data []byte
}

func (k *Key) Class() string {
	return k.classname
}

func (k *Key) Name() string {
	return k.name
}

func (k *Key) Title() string {
	return k.title
}

func (k *Key) load() ([]byte, error) {
	if !k.read {
		if int64(cap(k.data)) < int64(k.objlen) {
			k.data = make([]byte, k.objlen)
		}
		_, err := io.ReadFull(k.f, k.data)
		if err != nil {
			return nil, err
		}
		k.read = true
	}
	return k.data, nil
}

// Note: this contains ZL[src][dst] where src and dst are 3 bytes each.
// Won't bother with this for the moment, since we can cross-check against
// objlen.
const ROOT_HDRSIZE = 9

// Bytes returns the buffer of bytes corresponding to the Key's value
func (k *Key) Bytes() ([]byte, error) {
	if k.isCompressed() {
		// ... therefore it's compressed
		start := k.seekkey + int64(k.keylen) + ROOT_HDRSIZE
		r := io.NewSectionReader(k.f, start, int64(k.bytes)-int64(k.keylen))
		rc, err := zlib.NewReader(r)
		if err != nil {
			panic(err)
		}

		buf := &bytes.Buffer{}
		_, err = io.Copy(buf, rc)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	// ... not compressed
	start := k.seekkey + int64(k.keylen)
	r := io.NewSectionReader(k.f, start, int64(k.bytes))
	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Return Basket data associated with this key, if there is any
func (k *Key) AsBasket() *Basket {
	if k.classname != "TBasket" {
		panic("rootio.Key: Key is not a basket!")
	}
	b := &Basket{}
	_, err := k.f.Seek(int64(k.pdat), os.SEEK_SET)
	if err != nil {
		panic(fmt.Errorf("rootio.Key: %v", err))
	}
	err = k.f.readBin(b)
	if err != nil {
		panic(fmt.Errorf("rootio.Key: %v", err))
	}
	return b
}

func (k *Key) isCompressed() bool {
	return k.objlen != k.bytes-k.keylen
}

func (k *Key) Read() error {
	var err error
	f := k.f

	key_offset := f.Tell()

	err = f.readBin(&k.bytes)
	if err != nil {
		return err
	}

	if k.bytes < 0 {
		//fmt.Println("Jumping gap: ", k.bytes)
		k.classname = "[GAP]"
		_, err = f.Seek(int64(-k.bytes)-4, os.SEEK_CUR)
		return err
	}
	err = f.readBin(&k.version)
	if err != nil {
		return err
	}

	err = f.readBin(&k.objlen)
	if err != nil {
		return err
	}

	err = f.readBin(&k.datetime)
	if err != nil {
		return err
	}

	err = f.readInt16(&k.keylen)
	if err != nil {
		return err
	}

	err = f.readBin(&k.cycle)
	if err != nil {
		return err
	}

	if k.version > 1000 {
		err = f.readBin(&k.seekkey)
		if err != nil {
			return err
		}
	} else {
		err = f.readInt32(&k.seekkey)
		if err != nil {
			return err
		}
	}

	if k.version > 1000 {
		err = f.readBin(&k.seekpdir)
		if err != nil {
			return err
		}
	} else {
		err = f.readInt32(&k.seekpdir)
		if err != nil {
			return err
		}
	}

	if k.seekkey == 0 {
		k.seekkey = key_offset
	}
	if k.seekkey != key_offset {
		//pretty.Printf("%+v", k)
		panic(fmt.Errorf("Consistency failure: key offset %v doesn't match actual offset %v",
			k.seekkey, key_offset))
	}

	err = f.readString(&k.classname)
	if err != nil {
		return err
	}

	err = f.readString(&k.name)
	if err != nil {
		return err
	}

	err = f.readString(&k.title)
	if err != nil {
		return err
	}

	k.pdat = f.Tell()

	_, err = f.Seek(int64(key_offset+int64(k.bytes)), os.SEEK_SET)
	return err
}

// testing interfaces

var _ Object = (*Key)(nil)
