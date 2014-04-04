package dao

type Annotation map[string]interface{}

type Histogram interface {
	Annotation() Annotation
	Name() string
	Rank() int
	Axis(int) Axis
	Entries() int64
}

// EOF
