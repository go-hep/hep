package fwk

type Base struct {
	Type string
	Name string
}

func (c Base) CompName() string {
	return c.Name
}

func (c Base) CompType() string {
	return c.Type
}

// EOF
