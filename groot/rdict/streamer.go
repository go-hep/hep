// Copyright 2018 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rdict

import (
	"fmt"
	"reflect"
	"strings"

	"go-hep.org/x/hep/groot/rbase"
	"go-hep.org/x/hep/groot/rbytes"
	"go-hep.org/x/hep/groot/rcont"
	"go-hep.org/x/hep/groot/rmeta"
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
		named:  *rbase.NewNamed(typ.Name(), typ.Name()),
		objarr: rcont.NewObjArray(),
	}
	switch typ.Kind() {
	case reflect.Struct:
		si.elems = make([]rbytes.StreamerElement, 0, typ.NumField())
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
			panic(fmt.Errorf("rdict: invalid struct array field %#v", field))
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
			panic(fmt.Errorf("rdict: invalid struct slice field %#v", field))
		}

	default:
		panic(fmt.Errorf("rdict: invalid struct field %#v", field))
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

// StreamerInfo returns the named StreamerInfo.
// If version is negative, the latest version should be returned.
func (store *streamerStoreImpl) StreamerInfo(name string, version int) (rbytes.StreamerInfo, error) {
	return store.ctx.StreamerInfo(name, version)
}

func (store *streamerStoreImpl) addStreamer(si rbytes.StreamerInfo) {
	store.db[si.Name()] = si
}

func nameOf(field reflect.StructField) string {
	tag, ok := field.Tag.Lookup("groot")
	if !ok {
		return field.Name
	}
	if i := strings.Index(tag, "["); i > 0 {
		tag = tag[:i]
	}
	return tag
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
	_ streamerStore              = (*streamerStoreImpl)(nil)
	_ rbytes.StreamerInfoContext = (*streamerStoreImpl)(nil)
)
