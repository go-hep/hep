package rootio

import (
	"encoding/binary"
	"io"
)

type rootDecoder struct {
	r io.Reader
}

func (dec *rootDecoder) readString(s *string) error {
	var err error
	var length byte
	var buf [256]byte

	err = dec.readBin(&length)
	if err != nil {
		return err
	}

	if length != 0 {
		err = dec.readBin(buf[:length])
		if err != nil {
			return err
		}
		*s = string(buf[:length])
	}
	return err

}

func (dec *rootDecoder) readBin(v interface{}) error {
	return binary.Read(dec.r, E, v)
}

func (dec *rootDecoder) readInt16(v interface{}) error {
	var err error
	var d int16
	err = dec.readBin(&d)
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

func (dec *rootDecoder) readInt32(v interface{}) error {
	var err error
	switch uv := v.(type) {
	case *int32:
		err = dec.readBin(v)
	case *int64:
		var d int32
		err = dec.readBin(&d)
		*uv = int64(d)
	default:
		panic("Unknown type")
	}
	return err
}

func (dec *rootDecoder) readVersion() (version int16, position, bytecount int32, err error) {

	err = dec.readBin(&bytecount)
	if err != nil {
		return
	}

	err = dec.readBin(&version)
	if err != nil {
		return
	}

	var id int32
	err = dec.readBin(&id)
	if err != nil {
		return
	}

	var bits int32
	err = dec.readBin(&bits)
	if err != nil {
		return
	}

	//FIXME: hack
	var trash [8]byte
	err = dec.readBin(&trash)
	if err != nil {
		return
	}
	//fmt.Printf("## data = %#v\n", trash[:])

	return version, position, bytecount, err
}

func (dec *rootDecoder) readTNamed() (name, title string, err error) {

	/*
		// FIXME: handle kIsOnHeap || kIsReferenced
		// if (bits & kIsReferenced) == 0 {
		// 	var x int16
		// 	err = dec.readBin(&x)
		// 	if err != nil {
		// 		return
		// 	}
		// }
	*/
	err = dec.readString(&name)
	if err != nil {
		return name, title, err
	}

	err = dec.readString(&title)
	if err != nil {
		return name, title, err
	}

	return name, title, err
}

func (dec *rootDecoder) readTAttLine() (color, style, width int16, err error) {
	err = dec.readBin(&color)
	if err != nil {
		return
	}

	err = dec.readBin(&style)
	if err != nil {
		return
	}

	err = dec.readBin(&width)
	if err != nil {
		return
	}

	return
}

func (dec *rootDecoder) readTAttFill() (color, style int16, err error) {
	err = dec.readBin(&color)
	if err != nil {
		return
	}

	err = dec.readBin(&style)
	if err != nil {
		return
	}
	return
}

func (dec *rootDecoder) readTAttMarker() (color, style int16, width float32, err error) {
	err = dec.readBin(&color)
	if err != nil {
		return
	}

	err = dec.readBin(&style)
	if err != nil {
		return
	}

	err = dec.readBin(&width)
	if err != nil {
		return
	}
	return
}

// EOF
