package testdata

import (
	"fmt"
	"io"

	"github.com/go-hep/fwk"
)

type OutputStream struct {
	input string

	W io.Writer
}

func (out *OutputStream) Connect(ports []fwk.Port) error {
	var err error
	out.input = ports[0].Name
	return err
}

func (out *OutputStream) Write(ctx fwk.Context) error {
	var err error
	store := ctx.Store()
	v, err := store.Get(out.input)
	if err != nil {
		return err
	}

	data := v.(int64)
	_, err = out.W.Write([]byte(fmt.Sprintf("%d\n", data)))
	if err != nil {
		return err
	}

	return err
}

func (out *OutputStream) Disconnect() error {
	var err error
	w, ok := out.W.(io.Closer)
	if !ok {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return err
}
