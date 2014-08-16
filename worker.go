package fwk

type worker struct {
	slot  int
	keys  []string
	store datastore
	ctxs  []context
	msg   msgstream

	evts  <-chan int64
	quit  <-chan struct{}
	errch chan<- error
}

func (wrk *worker) run(tsks []Task) {
	for {
		select {
		case ievt, ok := <-wrk.evts:
			if !ok {
				return
			}
			wrk.msg.Infof(">>> running evt=%d...\n", ievt)
			err := wrk.store.reset(wrk.keys)
			if err != nil {
				wrk.errch <- err
				return
			}
			errch := make(chan error, len(tsks))
			for i, tsk := range tsks {
				go func(i int, tsk Task) {
					//wrk.msg.Infof(">>> running [%s]...\n", tsk.Name())
					ctx := wrk.ctxs[i]
					ctx.id = ievt
					errch <- tsk.Process(ctx)
					// FIXME(sbinet) dont be so eager to flush...
					ctx.msg.flush()
				}(i, tsk)
			}
			for i := 0; i < len(tsks); i++ {
				err := <-errch
				if err != nil {
					close(errch)
					wrk.msg.flush()
					wrk.errch <- err
					return
				}
			}
			wrk.msg.flush()

		case <-wrk.quit:
			return
		}
	}
	return
}
