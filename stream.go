package rio

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"fmt"
	"io"
	"os"
	"unsafe"
)

// Open opens and connects a RIO stream to a file for reading
func Open(fname string) (*Stream, error) {
	var stream *Stream
	var err error

	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	stream = &Stream{
		name: fname,
		f:    f,
		recs: make(map[string]*Record),
	}

	return stream, err
}

// Create opens and connects a RIO stream to a file for writing
func Create(fname string) (*Stream, error) {
	var stream *Stream
	var err error

	f, err := os.Create(fname)
	if err != nil {
		return nil, err
	}

	stream = &Stream{
		name: fname,
		f:    f,
		recs: make(map[string]*Record),
	}

	return stream, err
}

// Stream manages operations of a single RIO stream.
type Stream struct {
	name string   // stream name
	f    *os.File // file handle

	recpos  int64 // start position of last record read
	complvl int   // compression level

	recs map[string]*Record // records to read/write
}

// Fd returns the integer Unix file descriptor referencing the underlying open file.
func (stream *Stream) Fd() uintptr {
	return stream.f.Fd()
}

// Close closes a stream and the underlying file
func (stream *Stream) Close() error {
	return stream.f.Close()
}

// Stat returns the FileInfo structure describing underlying file. If there is an
// error, it will be of type *os.PathError.
func (stream *Stream) Stat() (os.FileInfo, error) {
	return stream.f.Stat()
}

// Sync commits the current contents of the stream to stable storage.
func (stream *Stream) Sync() error {
	return stream.f.Sync()
}

// Name returns the stream name
func (stream *Stream) Name() string {
	return stream.name
}

// FileName returns the name of the file connected to that stream
func (stream *Stream) FileName() string {
	return stream.f.Name()
}

// Mode returns the stream mode (as os.FileMode)
func (stream *Stream) Mode() (os.FileMode, error) {
	var mode os.FileMode
	fi, err := stream.f.Stat()
	if err != nil {
		return mode, err
	}

	return fi.Mode(), nil
}

// SetCompressionLevel sets the (zlib) compression level
func (stream *Stream) SetCompressionLevel(lvl int) {
	if lvl < 0 {
		stream.complvl = flate.DefaultCompression
	} else if lvl > 9 {
		stream.complvl = flate.BestCompression
	} else {
		stream.complvl = lvl
	}
}

// CurPos returns the current position in the file
//  -1 if error
func (stream *Stream) CurPos() int64 {
	pos, err := stream.f.Seek(0, 1)
	if err != nil {
		return -1
	}
	return pos
}

// Seek sets the offset for the next Read or Write on the stream to offset,
// interpreted according to whence:  0 means relative to the origin of the
// file, 1 means relative to the current offset, and 2 means relative to
// the end. It returns the new offset and an error, if any.
func (stream *Stream) Seek(offset int64, whence int) (int64, error) {
	return stream.f.Seek(offset, whence)
}

// Record adds a Record to the list of records to read/write or
// returns the Record with that name.
func (stream *Stream) Record(name string) *Record {
	rec, dup := stream.recs[name]
	if dup {
		return rec
	}
	rec = &Record{
		name:   name,
		unpack: false,
		blocks: make(map[string]Block),
	}
	stream.recs[name] = rec
	return stream.recs[name]
}

// HasRecord returns whether a Record with name n has been added to this Stream
func (stream *Stream) HasRecord(n string) bool {
	_, ok := stream.recs[n]
	return ok
}

// DelRecord removes the Record with name n from this Stream.
// DelRecord is a no-op if such a Record was not known to the Stream.
func (stream *Stream) DelRecord(n string) {
	delete(stream.recs, n)
}

// Records returns the list of Records currently attached to this Stream.
func (stream *Stream) Records() []*Record {
	recs := make([]*Record, 0, len(stream.recs))
	for _, rec := range stream.recs {
		recs = append(recs, rec)
	}
	return recs
}

func (stream *Stream) dump() {
	fmt.Printf("=========== stream [%s] ============\n", stream.name)
	fmt.Printf("::: records: (%d)\n", len(stream.recs))
	for k, rec := range stream.recs {
		fmt.Printf("::: %s: %v\n", k, rec)
	}
	return
}

// ReadRecord reads the next record
func (stream *Stream) ReadRecord() (*Record, error) {
	var err error
	var record *Record

	// fmt.Printf("~~~ Read()... ~~~~~~~~~~~~~~~~~~\n")
	// defer fmt.Printf("~~~ Read()... ~~~~~~~~~~~~~~~~~~ [done]\n")

	stream.recpos = -1

	requested := false
	// loop over records until a requested one turns up
	for !requested {

		stream.recpos = stream.CurPos()
		// fmt.Printf(">>> recpos=%d <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<\n", stream.recpos)

		// interpret: 1) length of the record header
		//            2) record marker
		var rechdr recordHeader
		err = stream.read(&rechdr)
		if err != nil {
			return nil, err
		}
		//fmt.Printf(">>> buf=%v\n", buf[:])
		//fmt.Printf(">>> rechdr=%v\n", rechdr)

		if rechdr.Typ != g_mark_record {
			return nil, ErrStreamNoRecMarker
		}

		curpos := stream.CurPos()
		// fmt.Printf(">>> pos --0: %d (%d)\n", curpos, rechdr.Len-8)
		var recdata recordData
		err = stream.read(&recdata)
		if err != nil {
			return nil, err
		}
		// fmt.Printf(">>> rec=%v\n", recdata)
		buf := make([]byte, align4(recdata.NameLen))
		_, err = stream.f.Read(buf)
		if err != nil {
			return nil, err
		}
		recname := string(buf[:recdata.NameLen])
		// fmt.Printf(">>> name=[%s]\n", recname)
		// fmt.Printf(">>> pos --1: %d [%d]\n", stream.CurPos(), recdata.NameLen)
		record = stream.Record(recname)
		record.options = recdata.Options
		requested = record != nil && record.Unpack()

		// if the record is not interesting, go to next record.
		// skip over any padding bytes inserted to make the next record header
		// start on a 4-bytes boundary in the file
		if !requested {
			recdata.DataLen = align4(recdata.DataLen)
			curpos, err = stream.Seek(int64(recdata.DataLen), 1)
			if curpos != int64(recdata.DataLen+rechdr.Len)+stream.recpos {
				//fmt.Printf("pos: %d\nrec: %d\nlen: %d\n", curpos, recdata.DataLen, stream.recpos)
				return nil, io.EOF
			}
			if err != nil {
				return nil, err
			}
			continue
		}

		// extract the compression bit from the options word
		compress := record.Compress()
		if !compress {
			// read the rest of the record data.
			// note that uncompressed data is *ALWAYS* aligned to a 4-bytes boundary
			// in the file, so no pad skipping is necessary
			buf = make([]byte, recdata.DataLen)
			_, err = stream.f.Read(buf)
			if err != nil {
				return nil, err
			}

		} else {
			// read the compressed record data
			cbuf := make([]byte, recdata.DataLen)
			_, err = stream.f.Read(cbuf)
			if err != nil {
				return nil, err
			}

			// handle padding bytes that may have been inserted to make the next
			// record header start on a 4-bytes boundary in the file.
			padlen := align4(recdata.DataLen) - recdata.DataLen
			if padlen > 0 {
				curpos, err = stream.Seek(int64(padlen), 1)
				if err != nil {
					return nil, err
				}
			}

			unzip, err := zlib.NewReader(bytes.NewBuffer(cbuf))
			if err != nil {
				return nil, err
			}
			buf = make([]byte, recdata.UCmpLen)
			nb, err := unzip.Read(buf)
			unzip.Close()
			if err != nil {
				return nil, err
			}
			if nb != len(buf) {
				return nil, io.EOF
			}
			//stream.recpos = recstart
		}
		recbuf := bytes.NewBuffer(buf)
		//fmt.Printf("::: recbuf: %d buf:%d\n", recbuf.Len(), len(buf))
		err = record.read(recbuf)
		if err != nil {
			return record, err
		}
	}
	return record, err
}

func (stream *Stream) WriteRecord(record *Record) error {
	var err error
	// fmt.Printf("~~~ Write(%v)...\n", record.Name())
	// defer fmt.Printf("~~~ Write(%v)... [done]\n", record.Name())

	rechdr := recordHeader{
		Len: 0,
		Typ: g_mark_record,
	}
	recdata := recordData{
		Options: record.options,
		DataLen: 0,
		UCmpLen: 0,
		NameLen: uint32(len(record.name)),
	}

	rechdr.Len = uint32(unsafe.Sizeof(rechdr)) + uint32(unsafe.Sizeof(recdata)) +
		uint32(recdata.NameLen)

	var buf bytes.Buffer
	err = record.write(&buf)
	if err != nil {
		return err
	}

	ucmplen := uint32(buf.Len())
	recdata.UCmpLen = ucmplen
	recdata.DataLen = ucmplen

	if record.Compress() {
		var b bytes.Buffer
		zip, err := zlib.NewWriterLevel(&b, stream.complvl)
		if err != nil {
			return err
		}
		_, err = zip.Write(buf.Bytes())
		if err != nil {
			return err
		}
		err = zip.Close()
		if err != nil {
			return err
		}
		recdata.DataLen = align4(uint32(b.Len()))

		buf = b
	}

	err = stream.write(&rechdr)
	if err != nil {
		return err
	}

	err = stream.write(&recdata)
	if err != nil {
		return err
	}

	_, err = stream.f.Write([]byte(record.name))
	if err != nil {
		return err
	}

	padlen := align4(recdata.NameLen) - recdata.NameLen
	if padlen > 0 {
		_, err = stream.f.Write(make([]byte, int(padlen)))
		if err != nil {
			return err
		}
	}

	n := int64(buf.Len())
	w, err := io.Copy(stream.f, &buf)
	if err != nil {
		return err
	}

	if n != w {
		return fmt.Errorf("rio: written to few bytes (%d). expected (%d)", w, n)
	}

	return err
}

func (stream *Stream) read(data interface{}) error {
	return bread(stream.f, data)
}

func (stream *Stream) write(data interface{}) error {
	return bwrite(stream.f, data)
}

// EOF
