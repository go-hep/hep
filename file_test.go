package rootio

import (
	"bytes"
	B "encoding/binary"
	"io"
	"os"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

func TestFileReader(t *testing.T) {
	f, err := Open("testdata/small.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

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

	{
		k := getkey("Int64")

		basket := k.AsBasket()
		data, err := k.ReadContents()
		if err != nil {
			t.Fatalf(err.Error())
		}
		buf := bytes.NewBuffer(data)

		if buf.Len() == 0 {
			t.Fatalf("invalid key size")
		}

		if true {
			fd, _ := os.Create("testdata_int64.bytes")
			io.Copy(fd, buf)
			fd.Close()
			// buf has been consumed...
			buf = bytes.NewBuffer(data)
		}

		for i := 0; i < int(basket.Nevbuf); i++ {
			var data int64
			err := B.Read(buf, E, &data)
			if err != nil {
				t.Fatalf("could not read entry [%d]: %v\n", i, err)
			}
			if data != int64(i) {
				t.Fatalf("expected data=%v (got=%v)\n", i, data)
			}
		}
	}

	{
		k := getkey("Float64")

		basket := k.AsBasket()
		data, err := k.ReadContents()
		if err != nil {
			t.Fatalf(err.Error())
		}
		buf := bytes.NewBuffer(data)

		if buf.Len() == 0 {
			t.Fatalf("invalid key size")
		}

		if true {
			fd, _ := os.Create("testdata_float64.bytes")
			io.Copy(fd, buf)
			fd.Close()
			// buf has been consumed...
			buf = bytes.NewBuffer(data)
		}

		for i := 0; i < int(basket.Nevbuf); i++ {
			var data float64
			err := B.Read(buf, E, &data)
			if err != nil {
				t.Fatalf("could not read entry [%d]: %v\n", i, err)
			}
			if data != float64(i) {
				t.Fatalf("expected data=%v (got=%v)\n", i, data)
			}
		}
	}
}
