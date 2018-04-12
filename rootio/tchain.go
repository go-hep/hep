package rootio

type tchain struct {
	trees []Tree
}

func (tchain) Class() string {
	return "TChain"
}

// Chain returns a tchain that is the concatenation of all the input Trees.
func Chain(trees ...Tree) tchain {
	var t tchain
	t.trees = append(t.trees, trees...)
	return t
}

func (t tchain) Entries() int64 {
	var v int64 = 0
	for i := range t.trees {
		v = v + t.trees[i].Entries()
	}
	return v

}

// TotBytes return the total number of bytes before compression
func (t tchain) TotBytes() int64 {
	var v int64 = 0
	for i := range t.trees {
		v = v + t.trees[i].TotBytes()
	}
	return v

}

//Total number of bytes after compression

func (t tchain) ZipBytes() int64 {
	var v int64 = 0
	for i := range t.trees {
		v = v + t.trees[i].ZipBytes()
	}
	return v

}

func (t tchain) Branches() []Branch {

	return t.trees[0].Branches()
	/*var branch []Branch
	for i := range t.trees {
	branch= append(branch,t.trees[i].Branches())
	}


	return branch
	*/
}

func (t tchain) Branch(name string) Branch {

	for _, br := range t.trees[0].Branches() {
		if br.Name() == name {
			return br
		}
	}
	return nil

	/*	var branch Branch
		for i := range t.trees {
			for _, br := range t.trees[i].Branches() {
				if br.Name() == name {
					branch = append(branch, br)
				}
			}
			branch = append(branch, nil)
		}
		return branch
	*/
}

var (
	_ Tree = (*tchain)(nil)
)

func (t tchain) Leaves() []Leaf {

	return t.trees[0].Leaves()
}

func (t tchain) getFile() *File {

	return t.trees[0].getFile()
}

func (t tchain) loadEntry(i int64) error {
	for _, b := range t.trees[0].Branches() {
		err := b.loadEntry(i)
		if err != nil {
			return err
		}
	}
	return nil
}

//Methods of Named interface :

func (t tchain) Name() string {
	return t.trees[0].Name()
}

func (t tchain) Title() string {
	return t.trees[0].Title()
}
