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

func TestFileDirectory(t *testing.T) {
	f, err := Open("testdata/small.root")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer f.Close()

	pretty.DefaultConfig.IncludeUnexported = true

	for _, table := range []struct {
		name     string
		expected bool
	}{
		{"Int32", true},
		{"Int32;0", true},
		{"Int32;1", true}, //FIXME: currently, cycle is just ignored.
		{"Int32_nope", false},
		{"Int32_nope;0", false},

		{"Int64", true},
		{"Int64;0", true},
		{"Int64_nope", false},
		{"Int64_nope;0", false},

		{"Float64", true},
		{"Float64;0", true},
		{"Float64_nope", false},
		{"Float64_nope;0", false},

		{"ArrayFloat64", true},
		{"ArrayFloat64;0", true},
		{"ArrayFloat64_nope", false},
		{"ArrayFloat64_nope;0", false},

		{"tree", true},
	} {
		_, err := f.Get(table.name)
		ok := err == nil
		if ok != table.expected {
			t.Fatalf("%s: expected key to exist=%v (got=%v, err=%v)", table.name, table.expected, ok, err)
		}
	}

	for _, table := range []struct {
		name     string
		expected string
	}{
		{"Int32", "TBasket"},
		{"Int32;0", "TBasket"},

		{"Int64", "TBasket"},
		{"Int64;0", "TBasket"},

		{"Float64", "TBasket"},
		{"Float64;0", "TBasket"},

		{"ArrayFloat64", "TBasket"},
		{"ArrayFloat64;0", "TBasket"},

		{"tree", "TTree"},
	} {
		k, err := f.Get(table.name)
		if err != nil {
			t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
		}

		if k.Class() != table.expected {
			t.Fatalf("%s: expected key with class=%s (got=%v)", table.name, table.expected, k.Class())
		}
	}

	for _, table := range []struct {
		name     string
		expected string
	}{
		{"Int32", "Int32"},
		{"Int32;0", "Int32"},

		{"Int64", "Int64"},
		{"Int64;0", "Int64"},

		{"Float64", "Float64"},
		{"Float64;0", "Float64"},

		{"ArrayFloat64", "ArrayFloat64"},
		{"ArrayFloat64;0", "ArrayFloat64"},

		{"tree", "tree"},
	} {
		k, err := f.Get(table.name)
		if err != nil {
			t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
		}

		if k.Name() != table.expected {
			t.Fatalf("%s: expected key with name=%s (got=%v)", table.name, table.expected, k.Name())
		}
	}

	for _, table := range []struct {
		name     string
		expected string
	}{
		{"Int32", "tree"},
		{"Int32;0", "tree"},

		{"Int64", "tree"},
		{"Int64;0", "tree"},

		{"Float64", "tree"},
		{"Float64;0", "tree"},

		{"ArrayFloat64", "tree"},
		{"ArrayFloat64;0", "tree"},

		{"tree", "my tree title"},
	} {
		k, err := f.Get(table.name)
		if err != nil {
			t.Fatalf("%s: expected key to exist! (got %v)", table.name, err)
		}

		if k.Title() != table.expected {
			t.Fatalf("%s: expected key with title=%s (got=%v)", table.name, table.expected, k.Title())
		}
	}
}

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
