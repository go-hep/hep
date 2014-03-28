package rio

type Block interface {
	Name() string
	Xfer(stream *Stream, op Operation, version int) error
	Version() uint
}
