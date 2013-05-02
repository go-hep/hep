package hepmc

type Encoder interface {
	Encode(evt *Event) error
}

type Decoder interface {
	Decode(evt *Event) error
}

// EOF
