package hbook

// Object is the general handle to any hbook data analysis object.
type Object interface {
	Annotation() Annotation
	Name() string
}
