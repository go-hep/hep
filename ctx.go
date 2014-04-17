package fwk

type context struct {
	id    int64
	slot  int
	store Store
	msg   msgstream
}

func (ctx context) Id() int64 {
	return ctx.id
}

func (ctx context) Slot() int {
	return ctx.slot
}

func (ctx context) Store() Store {
	return ctx.store
}

func (ctx context) Msg() MsgStream {
	return ctx.msg
}

// EOF
