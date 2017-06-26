// Package ordsync allows to gather results (or do something else)
// from concurrent goroutines in particular order.
package ordsync

import (
	"sync"
)

// DeferGroup is a chain of Deferred functions
// that are executed in the order they were created
type DeferGroup struct {
	last *Deferred
}

// Wait blocks until last Deferred created by Defer method is done
// First Deferred function is called immediately
func (f *DeferGroup) Wait() {
	if f.last != nil {
		f.last.wait()
	}
}

// Defer returns new Deferred linked to the tail of this DeferGroup chain.
func (f *DeferGroup) Defer() Deferred {
	m := &sync.Mutex{}
	m.Lock()
	ret := Deferred{f.last, m}
	f.last = &ret
	return ret
}

// Deferred represents function that is always executed in the order it was created.
// Deferred is a link in DeferGroup chain. Deferred must be created by DeferGroup.DeferGroup() call.
type Deferred struct {
	prev *Deferred
	m    *sync.Mutex
}

func (d *Deferred) wait() {
	d.m.Lock()
}

// Do runs given function only after previous deferred function is done
// (until that time it blocks)
// IMPORTANT: can be called only once, panics on second call
func (d *Deferred) Do(f func()) {
	if d.m == nil {
		panic("Deferred.Do can be called only once")
	}
	if d.prev != nil {
		d.prev.wait()
	}
	f()
	d.m.Unlock()
	d.m = nil
	d.prev = nil
}
