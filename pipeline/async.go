package pipeline

import (
	"github.com/cnmax/gologging-ext/core"
	"github.com/cnmax/gologging-ext/internal"
)

type Async struct {
	writer core.Writer
	ch     chan *core.Entry
}

func NewAsync(writer core.Writer) *Async {
	a := &Async{
		writer: writer,
		ch:     make(chan *core.Entry, 1000),
	}

	go a.run()
	return a
}

func (a *Async) Async() bool {
	return true
}

func (a *Async) Write(e *core.Entry) error {
	select {
	case a.ch <- e:
		return nil
	default:
		return &core.BackpressureError{
			Capacity:   cap(a.ch),
			CurrentLen: len(a.ch),
		}
	}
}

func (a *Async) run() {
	for e := range a.ch {
		func() {
			defer func() {
				if r := recover(); r != nil {
					internal.Error("async", "run panic: %v", r)
				}
			}()

			if err := a.writer.Write(e); err != nil {
				internal.Error("async", "write error: %v", err)
			}
		}()
	}
}
