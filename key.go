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

func (k *Key) ClassName() string {
	return k.classname
}

func (k *Key) Name() string {
	return k.name
}

func (k *Key) Title() string {
	return k.title
}

func (k *Key) Data() []byte {
	if !k.read {
		if int64(cap(k.data)) < int64(k.objlen) {
			k.data = make([]byte, k.objlen)
		}
		io.ReadFull(k.f, k.data) // TODO(pwaller): Error check
	}
	return k.data
}

// Note: this contains ZL[src][dst] where src and dst are 3 bytes each.
// Won't bother with this for the moment, since we can cross-check against
// objlen.
const ROOT_HDRSIZE = 9

func (k *Key) ReadContents() []byte {
	if k.Compressed() {
		// ... therefore it's compressed
		start := k.seekkey + int64(k.keylen) + ROOT_HDRSIZE
		r := io.NewSectionReader(k.f, start, int64(k.bytes)-int64(k.keylen))
		rc, err := zlib.NewReader(r)
		if err != nil {
			panic(err)
		}

		buf := &bytes.Buffer{}
		io.Copy(buf, rc)
		return buf.Bytes()
	}
	// ... not compressed
	start := k.seekkey + int64(k.keylen)
	r := io.NewSectionReader(k.f, start, int64(k.bytes))
	buf := &bytes.Buffer{}
	io.Copy(buf, r)
	return buf.Bytes()
}

// Return Basket data associated with this key, if there is any
func (k *Key) AsBasket() *Basket {
	if k.classname != "TBasket" {
		panic("Key is not a basket!")
	}
	b := &Basket{}
	k.f.Seek(int64(k.pdat), os.SEEK_SET)
	k.f.ReadBin(b)
	return b
}

func (k *Key) Compressed() bool {
	return k.objlen != k.bytes-k.keylen
}

func (k *Key) Read() {
	f := k.f

	key_offset := f.Tell()

	f.ReadBin(&k.bytes)
	if k.bytes < 0 {
		//fmt.Println("Jumping gap: ", k.bytes)
		k.classname = "[GAP]"
		f.Seek(int64(-k.bytes)-4, os.SEEK_CUR)
		return
	}
	f.ReadBin(&k.version)
	f.ReadBin(&k.objlen)
	f.ReadBin(&k.datetime)
	f.ReadInt16(&k.keylen)
	f.ReadBin(&k.cycle)

	if k.version > 1000 {
		f.ReadBin(&k.seekkey)
	} else {
		f.ReadInt32(&k.seekkey)
	}

	if k.version > 1000 {
		f.ReadBin(&k.seekpdir)
	} else {
		f.ReadInt32(&k.seekpdir)
	}

	if k.seekkey == 0 {
		k.seekkey = key_offset
	}
	if k.seekkey != key_offset {
		//pretty.Printf("%+v", k)
		panic(fmt.Errorf("Consistency failure: key offset %v doesn't match actual offset %v",
			k.seekkey, key_offset))
	}

	f.ReadString(&k.classname)
	f.ReadString(&k.name)
	f.ReadString(&k.title)

	k.pdat = f.Tell()

	f.Seek(int64(key_offset+int64(k.bytes)), os.SEEK_SET)
}
