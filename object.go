package dao

// Object is the general handle to any dao analysis object.
type Object interface {
	Annotation() Annotation
	Name() string
}
