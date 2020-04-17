// Copyright ©2017 The go-hep Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package fwk provides a set of tools to process High Energy Physics events data.
// fwk is a components-based framework, a-la Gaudi, with builtin support for concurrency.
//
// A fwk application consists of a set of components (fwk.Task) which are:
//  - (optionally) configured
//  - started
//  - given the chance to process each event
//  - stopped
//
// Helper components (fwk.Svc) can provide additional features (such as a
// whiteboard/event-store service, a data-flow service, ...) but do not
// typically take (directly) part of the event processing.
//
// Typically, users will implement fwk.Tasks, ie:
//
//   type MyTask struct {
//     fwk.TaskBase
//   }
//
//   // Configure is called once, after having read the properties
//   // from the data-cards.
//   func (tsk *MyTask) Configure(ctx fwk.Context) error { return nil }
//
//   // StartTask is called once (sequentially), just before
//   // the main event-loop processing.
//   func (tsk *MyTask) StartTask(ctx fwk.Context) error { return nil }
//
//   // Process is called for each event, (quite) possibly concurrently.
//   func (tsk *MyTask) Process(ctx fwk.Context)   error { return nil }
//
//   // StopTask is called once (sequentially), just after the
//   // main event-loop processing finished.
//   func (tsk *MyTask) StopTask(ctx fwk.Context)  error { return nil }
//
// A fwk application processes data and leverages concurrency at
// two different levels:
//  - event-level concurrency: multiple events are processed concurrently
//    at any given time, during the event loop;
//  - task-level concurrency: during the event loop, multiple tasks are
//    executing concurrently.
//
// To ensure the proper self-consistency of the global processed event,
// components need to express their data dependencies (input(s)) as well
// as the data they produce (output(s)) for downstream components.
// This is achieved by the concept of a fwk.Port.
// A fwk.Port consists of a pair { Name string; Type reflect.Type }
// where 'Name' is the unique location in the event-store,
// and 'Type' the expected 'go' type of the data at that event-store location.
//
// fwk.Ports can be either INPUT ports or OUTPUT ports.
// Components declare INPUT ports and OUTPUT ports during the 'Configure' stage
// of a fwk application, like so:
//
//  t := reflect.TypeOf([]Electron{})
//  err = component.DeclInPort("Electrons", t)
//  err = component.DeclOutPort("ReScaledElectrons", t)
//
// Then, during the event processing, one gets and puts data from/to the store
// like so:
//
//   func (tsk *MyTask) Process(ctx fwk.Context) error {
//      var err error
//
//      // retrieve the store associated with this event / region-of-interest
//      store := ctx.Store()
//
//      v, err := store.Get("Electrons")
//      if err != nil {
//         return err
//      }
//      eles := v.([]Electron) // type-cast to the correct (underlying) type
//
//      // create output collection
//      out := make([]Electron, 0, len(eles))
//
//      // make sure the collection will be put in the store
//      defer func() {
//         err = store.Put("ReScaledElectrons", out)
//      }()
//
//      // ... do some massaging with 'eles' and 'out'
//
//      return err
//   }
package fwk // import "go-hep.org/x/hep/fwk"
