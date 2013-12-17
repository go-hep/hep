package rootio

import (
	"bytes"
	"io"
	"log"
	"os"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestFileReader(t *testing.T) {
	f, err := Open("test-small.root")
	if err != nil {
		t.Fatal(err.Error())
	}

	pretty.DefaultConfig.IncludeUnexported = true

	//f.Map()
	//return

	getkey := func(n string) Key {
		var k Key
		for _, k = range f.keys {
			if k.name == n {
				return k
			}
		}
		t.Fatalf("could not find key [%s]", n)
		return Key{}
	}

	k := getkey("mcevt_weight")

	basket := k.AsBasket()
	buf := bytes.NewBuffer(k.ReadContents())

	if buf.Len() == 0 {
		t.Fatalf("invalid key size")
	}

	fd, _ := os.Create("mcevt.bytes")
	io.Copy(fd, buf)
	fd.Close()

	var contents [1024]uint32
	for i := 0; i < int(basket.Nevbuf); i++ {
		n := k.DecodeVector(buf, contents[:])
		log.Print("Event: ", contents[:n])
	}

}
