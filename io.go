package fwk

// InputStreamer reads data from the underlying io.Reader
// and puts it into fwk's Context
type InputStreamer interface {
	Read(ctx Context) error
}

// OutputStreamer gets data from the Context
// and writes it to the underlying io.Writer
type OutputStreamer interface {
	Write(ctx Context) error
}

// EOF
