package rio

import (
	"encoding/binary"
	"io"
	"reflect"
)

// align4 returns sz adjusted to align at 4-byte boundaries
func align4(sz uint32) uint32 {
	return sz + (4-(sz&g_align))&g_align
}

func bread(r io.Reader, data interface{}) error {
	bo := binary.BigEndian
	rv := reflect.ValueOf(data)
	//fmt.Printf("::: [%v] :::...\n", rv.Type())
	//defer fmt.Printf("### [%v] [done]\n", rv.Type())

	switch rv.Type().Kind() {
	case reflect.Ptr:
		rv = rv.Elem()
	}
	switch rv.Type().Kind() {
	case reflect.Struct:
		//fmt.Printf(">>> struct: [%v]...\n", rv.Type())
		for i, n := 0, rv.NumField(); i < n; i++ {
			//fmt.Printf(">>> i=%d [%v] (%v)...\n", i, rv.Field(i).Type(), rv.Type().Name())
			err := bread(r, rv.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			//fmt.Printf(">>> i=%d [%v] (%v)...[done]\n", i, rv.Field(i).Type(), rv.Type().Name())
		}
		//fmt.Printf(">>> struct: [%v]...[done]\n", rv.Type())
		return nil
	case reflect.String:
		sz := uint32(0)
		err := bread(r, &sz)
		if err != nil {
			return err
		}
		strlen := align4(sz)
		//fmt.Printf("string [%d=>%d]...\n", sz, strlen)
		str := make([]byte, strlen)
		err = bread(r, str)
		if err != nil {
			return err
		}
		rv.SetString(string(str[:sz]))
		//fmt.Printf("++++ [%s]\n", string(str))
		return nil

	case reflect.Slice:
		//fmt.Printf(">>> slice: [%v|%v]...\n", rv.Type(), rv.Type().Elem().Kind())
		rve := rv.Type().Elem()
		if kk := rve.Kind(); kk != reflect.String && kk != reflect.Struct {
			return binary.Read(r, bo, data)
		}
		sz := uint32(0)
		err := bread(r, &sz)
		if err != nil {
			return err
		}
		//fmt.Printf(">>> slice: %d [%v]\n", sz, rv.Type())
		slice := reflect.MakeSlice(rv.Type(), int(sz), int(sz))
		for i := 0; i < int(sz); i++ {
			err = bread(r, slice.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		rv.Set(slice)
		//fmt.Printf(">>> slice: [%v]... [done] (%v)\n", rv.Type(), slice.Interface())
		return err

	case reflect.Map:
		//fmt.Printf(">>> map: [%v]...\n", rv.Type())
		sz := uint32(0)
		err := bread(r, &sz)
		if err != nil {
			return err
		}
		//fmt.Printf(">>> map: %d [%v]\n", sz, rv.Type())
		m := reflect.MakeMap(rv.Type())
		kt := rv.Type().Key()
		vt := rv.Type().Elem()
		for i, n := 0, int(sz); i < n; i++ {
			kv := reflect.New(kt)
			err = bread(r, kv.Interface())
			if err != nil {
				return err
			}
			vv := reflect.New(vt)
			err = bread(r, vv.Interface())
			if err != nil {
				return err
			}
			//fmt.Printf("m - %d: {%v} - {%v}\n", i, kv.Elem().Interface(), vv.Elem().Interface())
			m.SetMapIndex(kv.Elem(), vv.Elem())
		}
		rv.Set(m)
		return nil

	default:
		//fmt.Printf(">>> binary - [%v]...\n", rv.Type())
		err := binary.Read(r, bo, data)
		//fmt.Printf(">>> binary - [%v]... [done]\n", rv.Type())
		return err
	}
	panic("not reachable")
}

func bwrite(w io.Writer, data interface{}) error {
	//return binary.Write(w, binary.BigEndian, data)
	bo := binary.BigEndian
	rv := reflect.ValueOf(data)
	//fmt.Printf("::: [%v] :::...\n", rv.Type())
	//defer fmt.Printf("### [%v] [done]\n", rv.Type())

	switch rv.Type().Kind() {
	case reflect.Ptr:
		rv = rv.Elem()
	}
	switch rv.Type().Kind() {
	case reflect.Struct:
		//fmt.Printf(">>> struct: [%v]...\n", rv.Type())
		for i, n := 0, rv.NumField(); i < n; i++ {
			//fmt.Printf(">>> i=%d [%v] (%v)...\n", i, rv.Field(i).Type(), rv.Type().Name())
			err := bwrite(w, rv.Field(i).Addr().Interface())
			if err != nil {
				return err
			}
			//fmt.Printf(">>> i=%d [%v] (%v)...[done]\n", i, rv.Field(i).Type(), rv.Type().Name())
		}
		//fmt.Printf(">>> struct: [%v]...[done]\n", rv.Type())
		return nil
	case reflect.String:
		str := rv.String()
		sz := uint32(len(str))
		err := bwrite(w, &sz)
		if err != nil {
			return err
		}
		bstr := []byte(str)
		bstr = append(bstr, make([]byte, align4(sz)-sz)...)
		_, err = w.Write(bstr)
		if err != nil {
			return err
		}
		//fmt.Printf("<++++ (%d) [%s]\n", sz, string(str))
		return nil

	case reflect.Slice:
		//fmt.Printf(">>> slice: [%v|%v]...\n", rv.Type(), rv.Type().Elem().Kind())
		rve := rv.Type().Elem()
		if kk := rve.Kind(); kk != reflect.String && kk != reflect.Struct {
			return binary.Write(w, bo, data)
		}
		sz := uint32(rv.Len())
		err := bwrite(w, &sz)
		if err != nil {
			return err
		}
		for i := 0; i < int(sz); i++ {
			err = bwrite(w, rv.Index(i).Addr().Interface())
			if err != nil {
				return err
			}
		}
		//fmt.Printf(">>> slice: [%v]... [done] (%v)\n", rv.Type(), slice.Interface())
		return err

	case reflect.Map:
		//fmt.Printf(">>> map: [%v]...\n", rv.Type())
		sz := uint32(rv.Len())
		err := bwrite(w, &sz)
		if err != nil {
			return err
		}
		//fmt.Printf(">>> map: %d [%v]\n", sz, rv.Type())
		for _, kv := range rv.MapKeys() {
			vv := rv.MapIndex(kv)
			err = bwrite(w, kv.Interface())
			if err != nil {
				return err
			}
			err = bwrite(w, vv.Interface())
			if err != nil {
				return err
			}
			//fmt.Printf("m - %d: {%v} - {%v}\n", i, kv.Elem().Interface(), vv.Elem().Interface())
		}
		return nil

	default:
		//fmt.Printf(">>> binary - [%v]...\n", rv.Type())
		err := binary.Write(w, bo, data)
		//fmt.Printf(">>> binary - [%v]... [done]\n", rv.Type())
		return err
	}
	panic("not reachable")

}

// EOF
