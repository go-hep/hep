//An update of go-hep/hep#44


package rootio


type tchain struct {
	trees []Tree
}

//Class returns the ROOT class of the argument.
func (tchain) Class() string {
	return "TChain"
}

// Chain returns a tchain that is the concatenation of all the input Trees.
func Chain(trees ...Tree) tchain {
	var t tchain
	t.trees = append(t.trees, trees...)
	return t
}


//Entries returns the total number of entries 
func (t tchain) Entries() int64 {
	var v int64 = 0
	for i := range t.trees {
		v = v + t.trees[i].Entries()
	}
	return v

}
	

// TotBytes return the total number of bytes before compression.
func (t tchain) TotBytes() int64 {
	var v int64 = 0
	for i := range t.trees {
		v = v + t.trees[i].TotBytes()
	}
	return v

}

//ZipBytes returns the total number of bytes after compression.
func (t tchain) ZipBytes() int64 {
	var v int64 = 0
	for i := range t.trees {
		v = v + t.trees[i].ZipBytes()
	}
	return v

}


//Branches returns the list of branches.
func (t tchain) Branches() []Branch {

	return t.trees[0].Branches()
	
}


//Branch returns the branch whose name is the argument.
func (t tchain) Branch(name string) Branch {

	for _, br := range t.trees[0].Branches() {
		if br.Name() == name {
			return br
		}
	}
	return nil
}

var (
	_ Tree = (*tchain)(nil)
)


//Leaves returns direct pointers to individual branch leaves.
func (t tchain) Leaves() []Leaf {

	return t.trees[0].Leaves()
}


//getFile returns the underlying file.
func (t tchain) getFile() *File {

	return t.trees[0].getFile()
}


//loadEntry returns an error if there is a problem during the loading
func (t tchain) loadEntry(i int64) error {
	for _, b := range t.trees[0].Branches() {
		err := b.loadEntry(i)
		if err != nil {
			return err
		}
	}
	return nil
}

//Name returns the name of the ROOT objet in the argument.
func (t tchain) Name() string {
	return t.trees[0].Name()
}


//Title returns the title of the ROOT object in the argument
func (t tchain) Title() string {
	return t.trees[0].Title()
}
