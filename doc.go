// Package rootio provides a pure-go read-access to ROOT files.
// rootio might, with time, provide write-access too.
//
// A typical usage is as follow:
//
//   f, err := rootio.Open("ntup.root")
//   t := f.Get("tree").(*rootio.Tree)
//   fmt.Printf("entries= %v\n", t.Entries())
package rootio

// EOF
