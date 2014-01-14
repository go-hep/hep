package rootio

import (
	"bytes"
	"fmt"
	"reflect"
)

// A Tree object is a list of Branch.
//   To Create a TTree object one must:
//    - Create the TTree header via the TTree constructor
//    - Call the TBranch constructor for every branch.
//
//   To Fill this object, use member function Fill with no parameters
//     The Fill function loops on all defined TBranch
type Tree struct {
	f *File // underlying file

	named named

	entries  int64 // Number of entries
	totbytes int64 // Total number of bytes in all branches before compression
	zipbytes int64 // Total number of bytes in all branches after  compression

	branches []Object // list of branches
	leaves   []Object // direct pointers to individual branch leaves
}

func (tree *Tree) Class() string {
	return "TTree" //tree.classname
}

func (tree *Tree) Name() string {
	return tree.named.Name()
}

func (tree *Tree) Title() string {
	return tree.named.Title()
}

func (tree *Tree) Entries() int64 {
	return tree.entries
}

func (tree *Tree) TotBytes() int64 {
	return tree.totbytes
}

func (tree *Tree) ZipBytes() int64 {
	return tree.zipbytes
}

// ROOTUnmarshaler is the interface implemented by an object that can
// unmarshal itself from a ROOT buffer
func (tree *Tree) UnmarshalROOT(data []byte) error {
	var err error
	dec := rootDecoder{r: bytes.NewBuffer(data)}

	vers, pos, bcnt, err := dec.readVersion()
	if err != nil {
		println(vers, pos, bcnt)
		return err
	}
	//fmt.Printf(">>> version: %v\n", vers)

	name, title, err := dec.readTNamed()
	if err != nil {
		return err
	}

	tree.name = name
	tree.title = title

	_, _, _, err = dec.readTAttLine()
	if err != nil {
		return err
	}

	_, _, err = dec.readTAttFill()
	if err != nil {
		return err
	}

	_, _, _, err = dec.readTAttMarker()
	if err != nil {
		return err
	}

	if vers < 16 {
		return fmt.Errorf(
			"rootio.Tree: tree [%s] with version [%v] is not supported (too old)",
			tree.name,
			vers,
		)
	}

	// FIXME: hack. where do these 18 bytes come from ?
	var trash [18]byte
	err = dec.readBin(&trash)
	if err != nil {
		return err
	}

	//fmt.Printf("### data = %v\n", dec.data.Bytes()[:64])
	err = dec.readBin(&tree.entries)
	if err != nil {
		return err
	}

	err = dec.readBin(&tree.totbytes)
	if err != nil {
		return err
	}

	err = dec.readBin(&tree.zipbytes)
	if err != nil {
		return err
	}

	return err
}

func init() {
	f := func() reflect.Value {
		o := &Tree{}
		return reflect.ValueOf(o)
	}
	Factory.db["TTree"] = f
	Factory.db["*rootio.Tree"] = f
}

// testing interfaces
var _ Object = (*Tree)(nil)
var _ ROOTUnmarshaler = (*Tree)(nil)

// EOF
