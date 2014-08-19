package fwk

type StreamControl struct {
	Ports []Port
	Ctx   chan Context
	Err   chan error
	Quit  chan struct{}
}

// InputStreamer reads data from the underlying io.Reader
// and puts it into fwk's Context
type InputStreamer interface {
	Connect(ports []Port) error
	Read(ctx Context) error
	Disconnect() error
}

// OutputStreamer gets data from the Context
// and writes it to the underlying io.Writer
type OutputStreamer interface {
	Connect(ports []Port) error
	Write(ctx Context) error
	Disconnect() error
}

// EOF
