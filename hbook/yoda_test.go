package hbook

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestReadYODAHeader(t *testing.T) {
	const mark = "BEGIN YODA_HISTO1D"
	for _, tc := range []struct {
		str  string
		want string
		err  error
	}{
		{
			str:  "BEGIN YODA_HISTO1D /name\n",
			want: "/name",
		},
		{
			str:  "BEGIN YODA_HISTO1D /name with whitespace\n",
			want: "/name with whitespace",
		},
		{
			str:  "BEGIN YODA /name",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s line", mark),
		},
		{
			str:  "BEGIN YODA /name\n",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
		{
			str:  "\nBEGIN YODA /name",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
		{
			str:  "\nBEGIN YODA /name\n",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
		{
			str:  " BEGIN YODA /name\n",
			want: "",
			err:  fmt.Errorf("hbook: could not find %s mark", mark),
		},
	} {
		t.Run(tc.want, func(t *testing.T) {
			v, err := readYODAHeader(bytes.NewBuffer([]byte(tc.str)), mark)
			if err == nil && tc.err != nil {
				t.Fatalf("got err=nil, want=%v", tc.err.Error())
			}
			if err != nil && tc.err == nil {
				t.Fatalf("got=%v, want=nil", err.Error())
			}
			if err != nil && tc.err != nil {
				if !reflect.DeepEqual(err.Error(), tc.err.Error()) { // FIXME(sbinet): use proper error comparison w/ Go1.13
					t.Fatalf("got error=%v, want=%v", err, tc.err)
				}
			}
			if v != tc.want {
				t.Fatalf("got: %q, want: %q", v, tc.want)
			}
		})
	}
}
