package dal

// Object is the general handle to any dal analysis object.
type Object interface {
	Annotation() Annotation
	Name() string
}
