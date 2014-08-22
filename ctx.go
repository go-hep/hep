package fwk

type context struct {
	id    int64
	slot  int
	store Store
	msg   msgstream
	mgr   App
}

func (ctx context) ID() int64 {
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

func (ctx context) Svc(n string) (Svc, error) {
	svc := ctx.mgr.GetSvc(n)
	if svc == nil {
		return nil, Errorf("fwk: no such service [%s]", n)
	}
	return svc, nil
}

// EOF
