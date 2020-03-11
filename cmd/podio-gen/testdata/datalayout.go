// Automatically generated. DO NOT EDIT.

package podio

// SimpleStruct
type SimpleStruct struct {
	X int32
	Y int32
	Z int32
	P [4]int32
}

// NotSoSimpleStruct
type NotSoSimpleStruct struct {
	Data SimpleStruct
}

// ex2::NamespaceStruct
type ex2_NamespaceStruct struct {
	X int32
	Y int32
}

// ex2::NamespaceInNamespaceStruct
type ex2_NamespaceInNamespaceStruct struct {
	Data ex2_NamespaceStruct
}

// EventInfo
// Event info
type EventInfo struct {
	Number int32 // event number
}

// ExampleHit
// Example Hit
type ExampleHit struct {
	CellID uint64  // cellID
	X      float64 // x-coordinate
	Y      float64 // y-coordinate
	Z      float64 // z-coordinate
	Energy float64 // measured energy deposit
}

// ExampleMC
// Example MC-particle
type ExampleMC struct {
	Energy    float64      // energy
	PDG       int32        // PDG code
	Parents   []*ExampleMC // parents
	Daughters []*ExampleMC // daughters
}

// ExampleCluster
// Cluster
type ExampleCluster struct {
	Energy   float64           // cluster energy
	Hits     []*ExampleHit     // hits contained in the cluster
	Clusters []*ExampleCluster // sub clusters used to create this cluster
}

// ExampleReferencingType
// Referencing Type
type ExampleReferencingType struct {
	Clusters []*ExampleCluster         // some refs to Clusters
	Refs     []*ExampleReferencingType // refs into same type
}

// ExampleWithVectorMember
// Type with a vector member
type ExampleWithVectorMember struct {
	Count []int32 // various ADC counts
}

// ExampleWithOneRelation
// Type with one relation member
type ExampleWithOneRelation struct {
	Cluster *ExampleCluster // a particular cluster
}

// ExampleWithComponent
// Type with one component
type ExampleWithComponent struct {
	Component NotSoSimpleStruct // a component
}

// ExampleForCyclicDependency1
// Type for cyclic dependency
type ExampleForCyclicDependency1 struct {
	Ref *ExampleForCyclicDependency2 // a ref
}

// ExampleForCyclicDependency2
// Type for cyclic dependency
type ExampleForCyclicDependency2 struct {
	Ref *ExampleForCyclicDependency1 // a ref
}

// ExampleWithString
// Type with a string
type ExampleWithString struct {
	TheString string // the string
}

// ex42::ExampleWithNamespace
// Type with namespace and namespaced member
type ex42_ExampleWithNamespace struct {
	Data ex2_NamespaceStruct // a component
}

// ex42::ExampleWithARelation
// Type with namespace and namespaced relation
type ex42_ExampleWithARelation struct {
	Number float32                      // just a number
	Ref    *ex42_ExampleWithNamespace   // a ref in a namespace
	Refs   []*ex42_ExampleWithNamespace // multiple refs in a namespace
}

// ExampleWithArray
// Datatype with an array member
type ExampleWithArray struct {
	ArrayStruct       NotSoSimpleStruct      // component that contains an array
	MyArray           [4]int32               // array-member without space to test regex
	AnotherArray2     [4]int32               // array-member with space to test regex
	Snail_case_array  [4]int32               // snail case to test regex
	Snail_case_Array3 [4]int32               // mixing things up for regex
	StructArray       [4]ex2_NamespaceStruct // an array containing structs
}
