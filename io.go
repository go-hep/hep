package hepmc

// An Encoder writes and encodes hepmc.Event objects into an output stream.
type Encoder interface {
	Encode(evt *Event) error
}

// A Decoder reads and decodes hepmc.Event objects from an input stream.
type Decoder interface {
	Decode(evt *Event) error
}

// EOF
