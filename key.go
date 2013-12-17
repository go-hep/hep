package rootio

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"os"

	//"github.com/kr/pretty"
)

type Key struct {
	f *File // underlying file

	bytes    int32
	version  int16
	objlen   int32
	datetime int32
	keylen   int32
	cycle    int16

	seekkey  int64
	seekpdir int64

	classname, name, title string

	pdat int64 // Pointer to everything after the above.

	read bool
	data []byte
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
	} else {
		// TODO..
	}
	return []byte{}
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
	f := k.uf

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
