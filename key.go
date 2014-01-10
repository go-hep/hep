package rootio

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"time"

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

	bytes    int32     // number of bytes for the compressed object+key
	version  int16     // version of the Key struct
	objlen   int32     // length of uncompressed object
	datetime time.Time // Date/Time when the object was written
	keylen   int32     // number of bytes for the Key struct
	cycle    int16     // cycle number of the object

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

// Value returns the data corresponding to the Key's value
func (k *Key) Value() interface{} {
	var v interface{}

	factory := Factory.Get(k.Class())
	if factory == nil {
		panic(fmt.Errorf("key[%v]: no factory for type [%s]\n", k.Name(), k.Class()))
	}

	vv := factory()
	if vv, ok := vv.Interface().(ROOTUnmarshaler); ok {
		data, err := k.Bytes()
		if err != nil {
			panic(fmt.Errorf("key[%v]: %v", k.Name(), err))
		}
		err = vv.UnmarshalROOT(data)
		if err != nil {
			panic(err)
		}
		v = vv
	} else {
		panic(fmt.Errorf(
			"key[%v]: type [%s] does not satisfy the ROOTUnmarshaler interface",
			k.Name(), k.Class(),
		))
	}

	// FIXME: hack.
	if vv, ok := v.(*Tree); ok {
		vv.f = k.f
	}
	return v
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

	dec := rootDecoder{r: k.f}
	err = dec.readBin(b)
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

	dec := rootDecoder{r: f}
	key_offset := f.Tell()
	myprintf(":: Key.Read (@%v)\n", key_offset)
	err = dec.readInt32(&k.bytes)
	if err != nil {
		return err
	}

	if k.bytes < 0 {
		//fmt.Println("Jumping gap: ", k.bytes)
		k.classname = "[GAP]"
		_, err = dec.r.(io.Seeker).Seek(int64(-k.bytes)-4, os.SEEK_CUR)
		return err
	}
	err = dec.readInt16(&k.version)
	if err != nil {
		return err
	}

	err = dec.readInt32(&k.objlen)
	if err != nil {
		return err
	}

	var datetime uint32
	err = dec.readBin(&datetime)
	if err != nil {
		return err
	}
	k.datetime = datime2time(datetime)

	err = dec.readInt16(&k.keylen)
	if err != nil {
		return err
	}

	err = dec.readInt16(&k.cycle)
	if err != nil {
		return err
	}

	if k.version > 1000 {
		err = dec.readInt64(&k.seekkey)
		if err != nil {
			return err
		}
		err = dec.readInt64(&k.seekpdir)
		if err != nil {
			return err
		}
	} else {
		err = dec.readInt32(&k.seekkey)
		if err != nil {
			return err
		}
		err = dec.readInt32(&k.seekpdir)
		if err != nil {
			return err
		}
	}

	myprintf("--- key ---\n")
	myprintf("key-nbytes:  %v (@%v)\n", k.bytes, key_offset)
	myprintf("key-version: %v\n", k.version)
	myprintf("key-objlen:  %v\n", k.objlen)
	myprintf("key-cdate:   %v\n", k.datetime)
	myprintf("key-keylen:  %v\n", k.keylen)
	myprintf("key-cycle:   %v\n", k.cycle)
	myprintf("key-seekkey: %v\n", k.seekkey)
	myprintf("key-seekpdir:%v\n", k.seekpdir)
	myprintf("key-compress: %v %v %v %v %v\n", k.isCompressed(), k.objlen, k.bytes-k.keylen, k.bytes, k.keylen)

	if k.seekkey == 0 {
		k.seekkey = key_offset
	}

	err = dec.readString(&k.classname)
	if err != nil {
		return err
	}

	err = dec.readString(&k.name)
	if err != nil {
		return err
	}

	err = dec.readString(&k.title)
	if err != nil {
		return err
	}

	myprintf("key-class: [%v]\n", k.classname)
	myprintf("key-name:  [%v]\n", k.name)
	myprintf("key-title: [%v]\n", k.title)
	myprintf("key-descr:  %v (@%v) [%v|%v|%v]\n", k.bytes, key_offset, k.classname, k.name, k.title)

	k.pdat = dec.r.(interface {
		Tell() int64
	}).Tell()

	if k.seekkey != key_offset {
		err = fmt.Errorf(
			"rootio.Key: Consistency failure: key offset %v doesn't match actual offset %v",
			k.seekkey, key_offset,
		)
		//fmt.Printf("*** %v\n", err)
		//panic(err)
	}

	_, err = dec.r.(io.Seeker).Seek(int64(key_offset+int64(k.bytes)), os.SEEK_SET)
	return err
}

// testing interfaces

var _ Object = (*Key)(nil)
