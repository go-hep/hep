// Code generated by brio-gen; DO NOT EDIT.

package hbook

import (
	"encoding/binary"
	"math"
)

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Dist0D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], uint64(o.N))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.SumW))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.SumW2))
	data = append(data, buf[:8]...)
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Dist0D) UnmarshalBinary(data []byte) (err error) {
	o.N = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	o.SumW = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	o.SumW2 = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Dist1D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.Dist.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.Stats.SumWX))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.Stats.SumWX2))
	data = append(data, buf[:8]...)
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Dist1D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Dist.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	o.Stats.SumWX = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	o.Stats.SumWX2 = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	_ = data
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Dist2D) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	{
		sub, err := o.X.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	{
		sub, err := o.Y.MarshalBinary()
		if err != nil {
			return nil, err
		}
		binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
		data = append(data, buf[:8]...)
		data = append(data, sub...)
	}
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.Stats.SumWXY))
	data = append(data, buf[:8]...)
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Dist2D) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.X.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		err = o.Y.UnmarshalBinary(data[:n])
		if err != nil {
			return err
		}
		data = data[n:]
	}
	o.Stats.SumWXY = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	_ = data
	return err
}
