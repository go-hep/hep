package rio

import (
	"io"
	"os"
	"reflect"

	"github.com/go-hep/fwk"
	"github.com/go-hep/rio"
)

// OutputStreamer writes data to a rio-stream.
type OutputStreamer struct {
	Name string         // output filename
	w    io.WriteCloser // underlying output file
	rio  *rio.Writer    // output rio-stream
	recs []*rio.Record  // list of connected records to write out
}

func (o *OutputStreamer) Connect(ports []fwk.Port) error {
	var err error

	// FIXME(sbinet): handle local/remote files, protocols
	o.w, err = os.Create(o.Name)
	if err != nil {
		return err
	}

	o.rio, err = rio.NewWriter(o.w)
	if err != nil {
		return err
	}

	for _, port := range ports {
		rec := o.rio.Record(port.Name)
		err = rec.Connect(port.Name, reflect.New(port.Type))
		if err != nil {
			return err
		}
		o.recs = append(o.recs, rec)
	}

	return err
}

func (o *OutputStreamer) Disconnect() error {
	defer o.w.Close()

	err := o.rio.Close()
	if err != nil {
		return err
	}

	err = o.w.Close()
	if err != nil {
		return err
	}

	return err
}

func (o *OutputStreamer) Write(ctx fwk.Context) error {
	var err error
	store := ctx.Store()

	for _, rec := range o.recs {
		n := rec.Name()
		blk := rec.Block(n)
		obj, err := store.Get(n)
		if err != nil {
			return err
		}
		err = blk.Write(obj)
		if err != nil {
			return err
		}

		err = rec.Write()
		if err != nil {
			return err
		}
	}
	return err
}
