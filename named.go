package rootio

import "bytes"

// The TNamed class is the base class for all named ROOT classes
// A TNamed contains the essential elements (name, title)
// to identify a derived object in containers, directories and files.
// Most member functions defined in this base class are in general
// overridden by the derived classes.
type named struct {
	name  string
	title string
}

// Name returns the name of the instance
func (n *named) Name() string {
	return n.name
}

// Title returns the title of the instance
func (n *named) Title() string {
	return n.title
}

func (n *named) Class() string {
	return "TNamed"
}

func (n *named) UnmarshalROOT(data []byte) error {
	var err error
	dec := rootDecoder{r: bytes.NewBuffer(data)}

	err = dec.readString(&n.name)
	if err != nil {
		return err
	}

	err = dec.readString(&n.title)
	if err != nil {
		return err
	}

	return err
}

// testing interfaces
var _ Object = (*named)(nil)
var _ ROOTUnmarshaler = (*named)(nil)

// EOF
