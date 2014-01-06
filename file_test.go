package rootio

import (
	"bytes"
	B "encoding/binary"
	"io"
	"os"
	"reflect"
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

	for _, table := range []struct {
		n string
		t reflect.Type
	}{
		{"Int32", reflect.TypeOf(int32(0))},
		{"Int64", reflect.TypeOf(int64(0))},
		{"UInt32", reflect.TypeOf(uint32(0))},
		{"UInt64", reflect.TypeOf(uint64(0))},
		{"Float32", reflect.TypeOf(float32(0))},
		{"Float64", reflect.TypeOf(float64(0))},

		{"ArrayInt32", reflect.TypeOf([10]int32{})},
		{"ArrayInt64", reflect.TypeOf([10]int64{})},
		{"ArrayUInt32", reflect.TypeOf([10]uint32{})},
		{"ArrayUInt64", reflect.TypeOf([10]uint64{})},
		{"ArrayFloat32", reflect.TypeOf([10]float32{})},
		{"ArrayFloat64", reflect.TypeOf([10]float64{})},
	} {
		k := getkey(table.n)

		basket := k.AsBasket()
		data, err := k.Bytes()
		if err != nil {
			t.Fatalf(err.Error())
		}
		buf := bytes.NewBuffer(data)

		if buf.Len() == 0 {
			t.Fatalf("invalid key size")
		}

		if true {
			fd, _ := os.Create("testdata/read_" + table.n + ".bytes")
			io.Copy(fd, buf)
			fd.Close()
			// buf has been consumed...
			buf = bytes.NewBuffer(data)
		}

		for i := 0; i < int(basket.Nevbuf); i++ {
			data := reflect.New(table.t)
			err := B.Read(buf, E, data.Interface())
			if err != nil {
				t.Fatalf("could not read entry [%d]: %v\n", i, err)
			}
			switch table.t.Kind() {
			case reflect.Array:
				for jj := 0; jj < table.t.Len(); jj++ {
					vref := reflect.ValueOf(i).Convert(table.t.Elem())
					vchk := reflect.ValueOf(data.Elem().Index(jj).Interface()).Convert(table.t.Elem())
					if !reflect.DeepEqual(vref, vchk) {
						t.Fatalf("%s: expected data[%d]=%v (got=%v)\n", table.n, jj, vref.Interface(), vchk.Interface())
					}
				}
			default:
				vref := reflect.ValueOf(i).Convert(table.t)
				vchk := reflect.ValueOf(data.Elem().Interface()).Convert(table.t)
				if !reflect.DeepEqual(vref, vchk) {
					t.Fatalf("%s: expected data=%v (got=%v)\n", table.n, vref.Interface(), vchk.Interface())
				}
			}
		}
	}

}
