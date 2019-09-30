// DO NOT EDIT; automatically generated by brio-gen

package briotest

import (
	"encoding/binary"
	"math"
)

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Hist) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.Name)))
	data = append(data, buf[:8]...)
	data = append(data, []byte(o.Name)...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.Data.X))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], uint64(o.i))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], uint64(o.u))
	data = append(data, buf[:8]...)
	data = append(data, byte(o.i8))
	binary.LittleEndian.PutUint16(buf[:2], uint16(o.i16))
	data = append(data, buf[:2]...)
	binary.LittleEndian.PutUint32(buf[:4], uint32(o.i32))
	data = append(data, buf[:4]...)
	binary.LittleEndian.PutUint64(buf[:8], uint64(o.i64))
	data = append(data, buf[:8]...)
	data = append(data, byte(o.u8))
	binary.LittleEndian.PutUint16(buf[:2], uint16(o.u16))
	data = append(data, buf[:2]...)
	binary.LittleEndian.PutUint32(buf[:4], uint32(o.u32))
	data = append(data, buf[:4]...)
	binary.LittleEndian.PutUint64(buf[:8], uint64(o.u64))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint32(buf[:4], math.Float32bits(o.f32))
	data = append(data, buf[:4]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.f64))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:4], math.Float32bits(real(o.c64)))
	data = append(data, buf[:4]...)
	binary.LittleEndian.PutUint64(buf[:4], math.Float32bits(imag(o.c64)))
	data = append(data, buf[:4]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(real(o.c128)))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(imag(o.c128)))
	data = append(data, buf[:8]...)
	switch o.b {
	case false:
		data = append(data, uint8(0))
	default:
		data = append(data, uint8(1))
	}
	for i := range o.arrI8 {
		o := &o.arrI8[i]
		data = append(data, byte(o))
	}
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.sliF64)))
	data = append(data, buf[:8]...)
	for i := range o.sliF64 {
		o := &o.sliF64[i]
		binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o))
		data = append(data, buf[:8]...)
	}
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
	binary.LittleEndian.PutUint64(buf[:8], uint64(len(o.sliPs)))
	data = append(data, buf[:8]...)
	for i := range o.sliPs {
		o := o.sliPs[i]
		{
			v := *o
			{
				sub, err := v.MarshalBinary()
				if err != nil {
					return nil, err
				}
				binary.LittleEndian.PutUint64(buf[:8], uint64(len(sub)))
				data = append(data, buf[:8]...)
				data = append(data, sub...)
			}
		}
	}
	{
		v := *o.ptr
		binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(v))
		data = append(data, buf[:8]...)
	}
	data = append(data, byte(o.myu8))
	binary.LittleEndian.PutUint16(buf[:2], uint16(o.myu16))
	data = append(data, buf[:2]...)
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Hist) UnmarshalBinary(data []byte) (err error) {
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		data = data[8:]
		o.Name = string(data[:n])
		data = data[n:]
	}
	o.Data.X = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	o.i = int(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	o.u = uint(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	o.i8 = int8(data[0])
	data = data[1:]
	o.i16 = int16(binary.LittleEndian.Uint16(data[:2]))
	data = data[2:]
	o.i32 = int32(binary.LittleEndian.Uint32(data[:4]))
	data = data[4:]
	o.i64 = int64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	o.u8 = uint8(data[0])
	data = data[1:]
	o.u16 = uint16(binary.LittleEndian.Uint16(data[:2]))
	data = data[2:]
	o.u32 = uint32(binary.LittleEndian.Uint32(data[:4]))
	data = data[4:]
	o.u64 = uint64(binary.LittleEndian.Uint64(data[:8]))
	data = data[8:]
	o.f32 = float32(math.Float32frombits(binary.LittleEndian.Uint32(data[:4])))
	data = data[4:]
	o.f64 = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	o.c64 = complex64(complex(math.Float32frombits(binary.LittleEndian.Uint32(data[:4])), math.Float32frombits(binary.LittleEndian.Uint32(data[4:8]))))
	data = data[8:]
	o.c128 = complex128(complex(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])), math.Float64frombits(binary.LittleEndian.Uint64(data[8:16]))))
	data = data[16:]
	switch data[i] {
	case 0:
		o.b = false
	default:
		o.b = true
	}
	data = data[1:]
	for i := range o.arrI8 {
		o.arrI8[i] = int8(data[0])
		data = data[1:]
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.sliF64 = make([]float64, n)
		data = data[8:]
		for i := range o.sliF64 {
			o.sliF64[i] = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
			data = data[8:]
		}
	}
	{
		n := int(binary.LittleEndian.Uint64(data[:8]))
		o.bins = make([]Bin, n)
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
		o.sliPs = make([]*Bin, n)
		data = data[8:]
		for i := range o.sliPs {
			var oi Bin
			{
				var v Bin
				{
					n := int(binary.LittleEndian.Uint64(data[:8]))
					data = data[8:]
					err = v.UnmarshalBinary(data[:n])
					if err != nil {
						return err
					}
					data = data[n:]
				}
				oi = &v

			}
			o.sliPs[i] = oi
		}
	}
	{
		var v float64
		v = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
		data = data[8:]
		o.ptr = &v

	}
	o.myu8 = U8(data[0])
	data = data[1:]
	o.myu16 = U16(binary.LittleEndian.Uint16(data[:2]))
	data = data[2:]
	return err
}

// MarshalBinary implements encoding.BinaryMarshaler
func (o *Bin) MarshalBinary() (data []byte, err error) {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.x))
	data = append(data, buf[:8]...)
	binary.LittleEndian.PutUint64(buf[:8], math.Float64bits(o.y))
	data = append(data, buf[:8]...)
	return data, err
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler
func (o *Bin) UnmarshalBinary(data []byte) (err error) {
	o.x = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	o.y = float64(math.Float64frombits(binary.LittleEndian.Uint64(data[:8])))
	data = data[8:]
	return err
}
