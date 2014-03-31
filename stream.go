package rio

import (
	"bytes"
	"compress/flate"
	"compress/zlib"
	"fmt"
	"io"
	"os"
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
		name:   fname,
		f:      f,
		bufcap: 8 * 1024,
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
		name:   fname,
		f:      f,
		bufcap: 8 * 1024,
	}

	return stream, err
}

// Stream manages operations of a single RIO stream.
type Stream struct {
	name string   // stream name
	f    *os.File // file handle

	record string // record name being read
	block  string // block name being read

	recpos  int64 // start position of last record read
	complvl int   // compression level

	buf    []byte // buffer of raw data
	bufcap int    // capacity of scratch buffer
}

// Close closes a stream and the underlying file
func (stream *Stream) Close() error {
	return stream.f.Close()
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

// ReadRecord reads the next record
func (stream *Stream) ReadRecord() (*Record, error) {
	var err error
	var record *Record

	stream.recpos = -1

	requested := false
	// loop over records until a requested one turns up
	for !requested {

		stream.recpos = stream.CurPos()
		fmt.Printf(">>> recpos=%d\n", stream.recpos)

		// interpret: 1) length of the record header
		//            2) record marker
		var rechdr recordHeader
		err = stream.read(&rechdr)
		if err != nil {
			return nil, err
		}

		//fmt.Printf(">>> buf=%v\n", buf[:])
		fmt.Printf(">>> hdr=%v\n", rechdr)
		fmt.Printf(">>> buftyp=0x%08x (0x%08x)\n", rechdr.BufType, g_mark_record)

		if rechdr.BufType != g_mark_record {
			return nil, ErrStreamNoRecMarker
		}

		curpos := stream.CurPos()
		fmt.Printf(">>> pos --0: %d (%d)\n", curpos, rechdr.HdrLen-8)
		var recdata recordData
		err = stream.read(&recdata)
		if err != nil {
			return nil, err
		}
		fmt.Printf(">>> rec=%v\n", recdata)
		buf := make([]byte, recdata.NameLen+((4-(recdata.NameLen&g_align))&g_align))
		err = stream.read(buf)
		if err != nil {
			return nil, err
		}
		recname := string(buf)
		fmt.Printf(">>> name=[%s]\n", recname)
		fmt.Printf(">>> pos --1: %d [%d]\n", stream.CurPos(), recdata.NameLen)
		// FIXME:
		// *record = Mgr.Record(recname)
		// requested = *record != nil && (*record).Unpack()
		requested = true

		// if the record is not interesting, go to next record.
		// skip over any padding bytes inserted to make the next record header
		// start on a 4-bytes boundary in the file
		if !requested {
			recdata.DataLen += (4 - (recdata.DataLen & g_align)) & g_align
			curpos, err = stream.Seek(int64(recdata.DataLen), 1)
			if curpos != int64(recdata.DataLen+rechdr.HdrLen)+stream.recpos {
				fmt.Printf("pos: %d\nrec: %d\nlen: %d\n", curpos, recdata.DataLen, stream.recpos)
				return nil, io.EOF
			}
			if err != nil {
				return nil, err
			}
			continue
		}

		record = &Record{
			name: recname,
		}

		// extract the compression bit from the options word
		compress := (recdata.Options & g_opt_compress) != 0
		if !compress {
			// read the rest of the record data.
			// note that uncompressed data is *ALWAYS* aligned to a 4-bytes boundary
			// in the file, so no pad skipping is necessary
			buf := make([]byte, recdata.DataLen)
			err = stream.read(buf)
			if err != nil {
				return nil, err
			}
			record.buf = buf

		} else {
			// read the compressed record data
			cbuf := make([]byte, recdata.DataLen)
			err = stream.read(cbuf)
			if err != nil {
				return nil, err
			}

			// handle padding bytes that may have been inserted to make the next
			// record header start on a 4-bytes boundary in the file.
			padlen := (4 - (recdata.DataLen & g_align)) & g_align
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
			buf := make([]byte, recdata.UCmpLen)
			nb, err := unzip.Read(buf)
			unzip.Close()
			if err != nil {
				return nil, err
			}
			if nb != len(buf) {
				return nil, io.EOF
			}
			record.buf = buf
		}
	}
	return record, err
}

func (stream *Stream) read(data interface{}) error {
	return bread(stream.f, data)
}

func (stream *Stream) write(data interface{}) error {
	return bwrite(stream.f, data)
}

// EOF
