// fwk provides a set of tools to process High Energy Physics events data.
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
// A fwk application processes data concurrently, a 2 different levels:
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
package fwk
