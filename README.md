hepmc
=====

``hepmc`` is a pure ``Go`` implementation of the ``C++`` ``HepMC``
library.

## Installation

```sh
$ go get github.com/go-hep/hepmc
```

## Documentation

Doc is on [godoc](http://godoc.org/github.com/go-hep/hepmc)

## Example

```go
package main

import (
  "fmt"
  "io"
  "os"
  
  "github.com/go-hep/hepmc"
)

func main() {
  f, err := os.Open("test.hepmc")
  if err != nil { panic(err) }
  defer f.Close()
  
  dec := hepmc.NewDecoder(f)
  if dec == nil { panic("nil decoder")}
  
  for {
      var evt hepmc.Event
      err = dec.Decode(&evt)
      if err == io.EOF {
          break
      }
      if err != nil {
          panic(err)
      }
      
      fmt.Printf("==evt: %d\n", evt.EventNumber)
      fmt.Printf("  #parts: %d\n", len(evt.Particles))
      fmt.Printf("  #verts: %d\n", len(evt.Vertices))
  }
}
```


