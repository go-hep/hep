package dao

type Annotation map[string]interface{}

// Histogram is an n-dim histogram (with weighted entries)
type Histogram interface {
	Annotation() Annotation
	Name() string
	Rank() int
	Axis(int) Axis
	Entries() int64
}

// EOF
