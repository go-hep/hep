// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
	"go-hep.org/x/hep/groot/internal/rmeta"
	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
)

// StreamerOf generates a StreamerInfo from a reflect.Type.
func StreamerOf(ctx rbytes.StreamerInfoContext, typ reflect.Type) rbytes.StreamerInfo {
	bldr := newStreamerBuilder(ctx, typ)
	return bldr.genStreamer(typ)
}

type streamerStore interface {
	rbytes.StreamerInfoContext
	addStreamer(si rbytes.StreamerInfo)
}

type streamerBuilder struct {
	ctx streamerStore
	typ reflect.Type
}

func newStreamerBuilder(ctx rbytes.StreamerInfoContext, typ reflect.Type) *streamerBuilder {
	return &streamerBuilder{ctx: newStreamerStore(ctx), typ: typ}
}

func (bld *streamerBuilder) genStreamer(typ reflect.Type) rbytes.StreamerInfo {
	si := &StreamerInfo{
		named: *rbase.NewNamed(typ.Name(), typ.Name()),
	}
	switch typ.Kind() {
	case reflect.Struct:
		for i := 0; i < typ.NumField(); i++ {
			ft := typ.Field(i)
			si.elems = append(si.elems, bld.genField(ft))
		}
	}
	return si
}

func (bld *streamerBuilder) genField(field reflect.StructField) rbytes.StreamerElement {
	switch field.Type.Kind() {
	case reflect.Bool:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::bool",
			},
		}
	case reflect.Int8:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::int8",
			},
		}
	case reflect.Int16:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::int16",
			},
		}
	case reflect.Int32:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::int32",
			},
		}
	case reflect.Int64:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::int64",
			},
		}
	case reflect.Uint8:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::uint8",
			},
		}
	case reflect.Uint16:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::uint16",
			},
		}
	case reflect.Uint32:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::uint32",
			},
		}
	case reflect.Uint64:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::uint64",
			},
		}
	case reflect.Float32:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::float32",
			},
		}
	case reflect.Float64:
		return &StreamerBasicType{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.GoType2ROOTEnum[field.Type],
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::float64",
			},
		}
	case reflect.String:
		return &StreamerString{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.TString,
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  "golang::string",
			},
		}
	case reflect.Struct:
		return &StreamerObjectAny{
			StreamerElement{
				named:  *rbase.NewNamed(nameOf(field), ""),
				etype:  rmeta.Any,
				esize:  int32(field.Type.Size()),
				offset: offsetOf(field),
				ename:  field.Type.Name(),
			},
		}

	case reflect.Array:
		et := field.Type.Elem()
		switch et.Kind() {
		case reflect.Bool:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::bool",
				},
			}
		case reflect.Int8:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::int8",
				},
			}
		case reflect.Int16:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::int16",
				},
			}
		case reflect.Int32:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::int32",
				},
			}
		case reflect.Int64:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::int64",
				},
			}
		case reflect.Uint8:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::uint8",
				},
			}
		case reflect.Uint16:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::uint16",
				},
			}
		case reflect.Uint32:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::uint32",
				},
			}
		case reflect.Uint64:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::uint64",
				},
			}
		case reflect.Float32:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::float32",
				},
			}
		case reflect.Float64:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.GoType2ROOTEnum[et],
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::float64",
				},
			}
		case reflect.String:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.TString,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::string",
				},
			}
		case reflect.Struct:
			return &StreamerBasicType{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.OffsetL + rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  typenameOf(et),
				},
			}
		default:
			panic(errors.Errorf("rdict: invalid struct array field %#v", field))
		}

	case reflect.Slice:
		et := field.Type.Elem()
		switch et.Kind() {
		case reflect.Bool:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::bool>",
				},
			}
		case reflect.Int8:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::int8>",
				},
			}
		case reflect.Int16:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::int16>",
				},
			}
		case reflect.Int32:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::int32>",
				},
			}
		case reflect.Int64:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::int64>",
				},
			}
		case reflect.Uint8:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::uint8>",
				},
			}
		case reflect.Uint16:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::uint16>",
				},
			}
		case reflect.Uint32:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::uint32>",
				},
			}
		case reflect.Uint64:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::uint64>",
				},
			}
		case reflect.Float32:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::float32>",
				},
			}
		case reflect.Float64:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::float64>",
				},
			}
		case reflect.String:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  "golang::slice<golang::string>",
				},
			}
		case reflect.Struct:
			return &StreamerObjectAny{
				StreamerElement{
					named:  *rbase.NewNamed(nameOf(field), ""),
					etype:  rmeta.Any,
					esize:  int32(field.Type.Size()),
					offset: offsetOf(field),
					ename:  fmt.Sprintf("golang::slice<%s>", typenameOf(et)),
				},
			}

		default:
			panic(errors.Errorf("rdict: invalid struct slice field %#v", field))
		}

	default:
		panic(errors.Errorf("rdict: invalid struct field %#v", field))
	}
}

type streamerStoreImpl struct {
	ctx rbytes.StreamerInfoContext
	db  map[string]rbytes.StreamerInfo
}

func newStreamerStore(ctx rbytes.StreamerInfoContext) streamerStore {
	if ctx, ok := ctx.(streamerStore); ok {
		return ctx
	}

	return &streamerStoreImpl{
		ctx: ctx,
		db:  make(map[string]rbytes.StreamerInfo),
	}
}

func (store *streamerStoreImpl) StreamerInfo(name string) (rbytes.StreamerInfo, error) {
	return store.ctx.StreamerInfo(name)
}

func (store *streamerStoreImpl) addStreamer(si rbytes.StreamerInfo) {
	store.db[si.Name()] = si
}

func nameOf(field reflect.StructField) string {
	tag, ok := field.Tag.Lookup("groot")
	if ok {
		return tag
	}
	return field.Name
}

func offsetOf(field reflect.StructField) int32 {
	// return int32(field.Offset)
	// FIXME(sbinet): it seems ROOT expects 0 here...
	return 0
}

func typenameOf(typ reflect.Type) string {
	return typ.Name()
}

var (
	_ streamerStore = (*streamerStoreImpl)(nil)
)
