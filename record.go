package rio

// recordHeader describes the on-disk record (header part)
type recordHeader struct {
	HdrLen  uint32
	BufType uint32
}

// recordData describes the on-disk record (payload part)
type recordData struct {
	Options uint32
	DataLen uint32 // length of compressed record data
	UCmpLen uint32 // length of uncompressed record data
	NameLen uint32 // length of record name
}

// Record manages blocks of data
type Record struct {
	name   string           // record name
	buf    []byte           // record payload
	unpack bool             // whether to unpack incoming records
	blocks map[string]Block // connected blocks
}

// Name returns the name of this record
func (rec *Record) Name() string {
	return rec.name
}

// Unpack returns whether to unpack incoming records
func (rec *Record) Unpack() bool {
	return rec.unpack
}

// EOF
