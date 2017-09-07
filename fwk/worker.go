// Copyright 2017 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fwk

import (
	"fmt"

	nctx "golang.org/x/net/context"
)

type workercontrol struct {
	evts   chan context
	done   chan struct{}
	errc   chan error
	runctx nctx.Context
}

type worker struct {
	slot int
	keys []string
	//store datastore
	ctxs []context
	msg  msgstream

	evts   <-chan context
	done   chan<- struct{}
	errc   chan<- error
	runctx nctx.Context
}

func newWorker(i int, app *appmgr, ctrl *workercontrol) *worker {
	wrk := &worker{
		slot:   i,
		keys:   app.dflow.keys(),
		ctxs:   make([]context, len(app.tsks)),
		msg:    newMsgStream(fmt.Sprintf("%s-worker-%03d", app.name, i), app.msg.lvl, nil),
		evts:   ctrl.evts,
		done:   ctrl.done,
		errc:   ctrl.errc,
		runctx: ctrl.runctx,
	}
	for j, tsk := range app.tsks {
		wrk.ctxs[j] = context{
			id:   -1,
			slot: i,
			msg:  newMsgStream(tsk.Name(), app.msg.lvl, nil),
			mgr:  nil, // nobody's supposed to access mgr's state during event-loop
		}
	}

	go wrk.run(app.tsks)

	return wrk
}

func (wrk *worker) run(tsks []Task) {
	defer func() {
		wrk.done <- struct{}{}
	}()

	for {
		select {
		case ievt, ok := <-wrk.evts:
			if !ok {
				return
			}
			wrk.msg.Debugf(">>> running evt=%d...\n", ievt.ID())

			evtstore := ievt.store.(*datastore)
			evtctx, evtCancel := nctx.WithCancel(wrk.runctx)
			evt := taskrunner{
				ievt:   ievt.ID(),
				errc:   make(chan error, len(tsks)),
				evtctx: evtctx,
			}
			for i, tsk := range tsks {
				ctx := wrk.ctxs[i]
				ctx.store = evtstore
				ctx.ctx = evtctx
				go evt.run(i, ctx, tsk)
			}
			ndone := 0
		errloop:
			for {
				select {
				case err, ok := <-evt.errc:
					if !ok {
						return
					}
					ndone++
					if err != nil {
						evtCancel()
						evtstore.close()
						wrk.msg.flush()

						wrk.errc <- err
						return
					}
					if ndone == len(tsks) {
						break errloop
					}
				case <-evtctx.Done():
					evtstore.close()
					wrk.msg.flush()
					return
				}
			}
			err := evtstore.reset(wrk.keys)
			evtstore.close()
			wrk.msg.flush()

			if err != nil {
				wrk.errc <- err
				return
			}
		case <-wrk.runctx.Done():
			//wrk.store.close()
			return
		}
	}
}

type taskrunner struct {
	errc   chan error
	evtctx nctx.Context

	ievt int64
}

func (run taskrunner) run(i int, ctx context, tsk Task) {
	ctx.id = run.ievt
	select {
	case run.errc <- tsk.Process(ctx):
		// FIXME(sbinet) dont be so eager to flush...
		ctx.msg.flush()
	case <-run.evtctx.Done():
		ctx.msg.flush()
	}
}
