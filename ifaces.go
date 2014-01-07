package rootio

// Class represents a ROOT class.
// Class instances are created by a ClassFactory.

type Class interface {
	// GetCheckSum gets the check sum for this ROOT class
	//CheckSum() int

	// Members returns the list of members for this ROOT class
	Members() []Member

	// Version returns the version number for this ROOT class
	Version() int

	// Name returns the ROOT class name for this ROOT class
	Name() string
}

// Member represents a single member of a ROOT class
type Member interface {
	// GetArrayDim returns the dimension of the array (if any)
	ArrayDim() int

	// GetComment returns the comment associated with this member
	Comment() string

	// Name returns the name of this member
	Name() string

	// Type returns the class of this member
	Type() Class

	// GetValue returns the value of this member
	//GetValue(o Object) reflect.Value
}

// Object represents a ROOT object
type Object interface {
	// Class returns the ROOT class of this object
	Class() string

	// Name returns the name of this ROOT object
	Name() string

	// Title returns the title of this ROOT object
	Title() string
}

// ClassFactory creates ROOT classes
type ClassFactory interface {
	Create(name string) Class
}

// Directory describes a ROOT directory structure in memory.
type Directory interface {
	// return pointer to object identified by namecycle
	//   namecycle has the format name;cycle
	//   name  = * is illegal, cycle = * is illegal
	//   cycle = "" or cycle = 9999 ==> apply to a memory object
	//
	//   examples:
	//     foo   : get object named foo in memory
	//             if object is not in memory, try with highest cycle from file
	//     foo;1 : get cycle 1 of foo on file
	Get(namecycle string) (Object, error)

	// return whether an object identified by namecycle exists in directory
	//   namecycle has the format name;cycle
	//   name  = * is illegal, cycle = * is illegal
	//   cycle = "" or cycle = 9999 ==> apply to a memory object
	//
	//   examples:
	//     foo   : get object named foo in memory
	//             if object is not in memory, try with highest cycle from file
	//     foo;1 : get cycle 1 of foo on file
	Has(namecycle string) bool
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
type ROOTUnmarshaler interface {
	UnmarshalROOT(data []byte) error
}

// ROOTMarshaler is the interface implemented by an object that can
// marshal itself into a ROOT buffer
type ROOTMarshaler interface {
	MarshalROOT() (data []byte, err error)
}

// EOF
