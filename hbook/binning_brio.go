// Code generated by brio-gen; DO NOT EDIT.

package hbook

import (
	"encoding/binary"
	"math"
)

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Range) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.Min))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.Max))
	data = append(data, buf[:8]...)
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Range) UnmarshalBinary(data []byte) (err error) {
	o.Min = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	o.Max = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Binning1D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.Bins)))
	data = append(data, buf[:8]...)
	for i := range o.Bins {
		o := &o.Bins[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	{
		sub, err := o.Dist.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	for i := range o.Outflows {
		o := &o.Outflows[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	{
		sub, err := o.XRange.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Binning1D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.Bins = make([]Bin1D, n)
		data = data[8:]
		for i := range o.Bins {
			oi := &o.Bins[i]
			{
				n := int(binary.LittleEndian.Uint64(data[:8]))
				data = data[8:]
				err = oi.UnmarshalBinary(data[:n])
				if err != nil {
					return err
				}
				data = data[n:]
			}
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Dist.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	for i := range o.Outflows {
		oi := &o.Outflows[i]
		{
			n := int(binary.LittleEndian.Uint64(data[:8]))
			data = data[8:]
			err = oi.UnmarshalBinary(data[:n])
			if err != nil {
				return err
			}
			data = data[n:]
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.XRange.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *binningP1D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.bins)))
	data = append(data, buf[:8]...)
	for i := range o.bins {
		o := &o.bins[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	{
		sub, err := o.dist.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	for i := range o.outflows {
		o := &o.outflows[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	{
		sub, err := o.xrange.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.xstep))
	data = append(data, buf[:8]...)
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *binningP1D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.bins = make([]BinP1D, n)
		data = data[8:]
		for i := range o.bins {
			oi := &o.bins[i]
			{
				n := int(binary.LittleEndian.Uint64(data[:8]))
				data = data[8:]
				err = oi.UnmarshalBinary(data[:n])
				if err != nil {
					return err
				}
				data = data[n:]
			}
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.dist.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	for i := range o.outflows {
		oi := &o.outflows[i]
		{
			n := int(binary.LittleEndian.Uint64(data[:8]))
			data = data[8:]
			err = oi.UnmarshalBinary(data[:n])
			if err != nil {
				return err
			}
			data = data[n:]
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.xrange.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	o.xstep = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Bin1D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.Range.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.Dist.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Bin1D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Range.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Dist.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *BinP1D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.xrange.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.dist.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *BinP1D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.xrange.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.dist.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Binning2D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.Bins)))
	data = append(data, buf[:8]...)
	for i := range o.Bins {
		o := &o.Bins[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	{
		sub, err := o.Dist.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	for i := range o.Outflows {
		o := &o.Outflows[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	{
		sub, err := o.XRange.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.YRange.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	binary.LittleEndian.PutUint64(buf[:8], uint64(o.Nx))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], uint64(o.Ny))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.XEdges)))
	data = append(data, buf[:8]...)
	for i := range o.XEdges {
		o := &o.XEdges[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.YEdges)))
	data = append(data, buf[:8]...)
	for i := range o.YEdges {
		o := &o.YEdges[i]
		{
			sub, err := o.MarshalBinary()
			if err != nil {
				return nil, err
			}
			binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
			data = append(data, buf[:8]...)
			data = append(data, sub...)
		}
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Binning2D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.Bins = make([]Bin2D, n)
		data = data[8:]
		for i := range o.Bins {
			oi := &o.Bins[i]
			{
				n := int(binary.LittleEndian.Uint64(data[:8]))
				data = data[8:]
				err = oi.UnmarshalBinary(data[:n])
				if err != nil {
					return err
				}
				data = data[n:]
			}
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Dist.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	for i := range o.Outflows {
		oi := &o.Outflows[i]
		{
			n := int(binary.LittleEndian.Uint64(data[:8]))
			data = data[8:]
			err = oi.UnmarshalBinary(data[:n])
			if err != nil {
				return err
			}
			data = data[n:]
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.XRange.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.YRange.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	o.Nx = int(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	o.Ny = int(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.XEdges = make([]Bin1D, n)
		data = data[8:]
		for i := range o.XEdges {
			oi := &o.XEdges[i]
			{
				n := int(binary.LittleEndian.Uint64(data[:8]))
				data = data[8:]
				err = oi.UnmarshalBinary(data[:n])
				if err != nil {
					return err
				}
				data = data[n:]
			}
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.YEdges = make([]Bin1D, n)
		data = data[8:]
		for i := range o.YEdges {
			oi := &o.YEdges[i]
			{
				n := int(binary.LittleEndian.Uint64(data[:8]))
				data = data[8:]
				err = oi.UnmarshalBinary(data[:n])
				if err != nil {
					return err
				}
				data = data[n:]
			}
		}
	}
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Bin2D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.XRange.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.YRange.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.Dist.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Bin2D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.XRange.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.YRange.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Dist.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	_ = data
	return err
}
