lhef2hepmc
==========

``lhef2hepmc`` is a simple `LHEF` -> `HEPMC` converter program.
It is a pure ``Go`` re-implementation of the `C++` converter from
`Rivet`.

## Installation

```sh
$ go get github.com/go-hep/lhef2hepmc
```

## Example

```sh
$ lhef2hepmc -i in.lhef -o out.hepmc
$ lhef2hepmc < in.lhef > out.hepmc
```

